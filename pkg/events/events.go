package events

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	// DefaultPollInterval is used for HTTP endpoints when Config.PollInterval
	// is zero.
	DefaultPollInterval = 3 * time.Second
	// DefaultRetryDelay is the initial delay before retrying a failed RPC call.
	DefaultRetryDelay = time.Second
	// DefaultMaxRetryDelay caps the exponential retry backoff.
	DefaultMaxRetryDelay = 30 * time.Second
	// subscriptionBuffer is the size of the channel fed by eth_subscribe.
	subscriptionBuffer = 64
)

// Watcher states returned by (*Watcher).Status.
const (
	StatusIdle    = "idle"    // created, Start not running
	StatusRunning = "running" // Start is active
	StatusStopped = "stopped" // Start returned after cancellation or error
	StatusClosed  = "closed"  // Close/Stop was called
)

// ErrAlreadyStarted is returned by Start when the watcher is already running.
var ErrAlreadyStarted = errors.New("events: watcher already started")

// Mode selects how logs are fetched.
type Mode uint8

const (
	// ModeAuto uses subscriptions for ws/wss/ipc endpoints and polling for
	// http(s) endpoints. It is the zero value and the usual choice.
	ModeAuto Mode = iota
	// ModeSubscribe forces eth_subscribe; it fails on endpoints that do not
	// support notifications.
	ModeSubscribe
	// ModePoll forces eth_getLogs polling, which works over any transport.
	ModePoll
)

// Backend is the subset of the Ethereum RPC API needed by a Watcher.
// *ethclient.Client implements it.
type Backend interface {
	BlockNumber(ctx context.Context) (uint64, error)
	FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error)
	SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)
}

// Handler is called once for each incoming event. Returning an error is
// reported through Config.OnError (or the logger) but does not stop the
// watcher.
type Handler func(ctx context.Context, event Event) error

// ErrorHandler receives transport errors and handler errors.
type ErrorHandler func(ctx context.Context, err error)

// Event is a log delivered by a Watcher. Log always holds the raw log. When
// Config.ABI is set, Name and Data contain the decoded event as well.
type Event struct {
	Log  types.Log
	Name string
	Data map[string]any
}

// Address returns the address that emitted the event.
func (e Event) Address() common.Address { return e.Log.Address }

// Topics returns the raw event topics, with the event signature in Topics[0].
func (e Event) Topics() []common.Hash { return e.Log.Topics }

// DataBytes returns the raw, ABI encoded non-indexed event arguments.
func (e Event) DataBytes() []byte { return e.Log.Data }

// BlockNumber returns the block that included the event.
func (e Event) BlockNumber() uint64 { return e.Log.BlockNumber }

// TxHash returns the hash of the transaction that emitted the event.
func (e Event) TxHash() common.Hash { return e.Log.TxHash }

// Removed reports whether the log was reverted by a chain reorganization.
func (e Event) Removed() bool { return e.Log.Removed }

// Config holds the settings required to watch a contract.
type Config struct {
	// URL is a JSON-RPC endpoint (http, https, ws, wss or an IPC path). It is
	// ignored when Client is set. One of URL and Client is required.
	URL string
	// Client is an existing backend. When set, the watcher does not own it and
	// will not close it.
	Client Backend

	// Address is the contract to watch. The zero address matches every
	// contract.
	Address common.Address
	// Topics narrows the logs using ethereum.FilterQuery topic semantics: a
	// nil position matches anything, a slice is an OR of alternatives.
	Topics [][]common.Hash
	// FromBlock is the first block to scan. Nil starts at the current head.
	FromBlock *big.Int

	// ABI, when set, is used to decode the name and arguments of each log.
	ABI *abi.ABI
	// Handler is invoked for every matching log. A nil handler discards logs.
	Handler Handler
	// Concurrency is the number of handlers that may run at once. The default
	// of 1 dispatches inline and keeps log order. Values above 1 use a bounded
	// worker pool; completion order is then not guaranteed.
	Concurrency int
	// HandlerTimeout limits how long one handler may run. Zero means no limit.
	// The deadline is passed to the handler through its context.
	HandlerTimeout time.Duration
	// OnError receives non-fatal errors. When nil, errors go to Logger.
	OnError ErrorHandler
	// Logger receives errors when OnError is nil. It must implement Printf
	// and Error; WithDefaults installs a NopLogger when it is nil.
	Logger Logger

	// Mode overrides the transport based subscription/poll choice.
	Mode Mode
	// PollInterval is the delay between eth_getLogs calls in polling mode.
	// Zero selects DefaultPollInterval.
	PollInterval time.Duration
	// RetryDelay is the initial backoff after a failed RPC call.
	// Zero selects DefaultRetryDelay.
	RetryDelay time.Duration
	// MaxRetryDelay caps the exponential backoff.
	// Zero selects DefaultMaxRetryDelay.
	MaxRetryDelay time.Duration
}

// DefaultConfig returns a config with all optional values filled in. Callers
// still need to set URL or Client (and usually a Handler).
func DefaultConfig() *Config { return (&Config{}).WithDefaults() }

// WithDefaults fills unset optional fields in place and returns c. New calls
// it automatically, so callers only need it when they want to inspect or
// override the defaults before constructing a watcher.
func (c *Config) WithDefaults() *Config {
	if c.Mode == 0 {
		c.Mode = ModeAuto
	}
	if c.Concurrency == 0 {
		c.Concurrency = 1
	}
	if c.PollInterval == 0 {
		c.PollInterval = DefaultPollInterval
	}
	if c.RetryDelay == 0 {
		c.RetryDelay = DefaultRetryDelay
	}
	if c.MaxRetryDelay == 0 {
		c.MaxRetryDelay = DefaultMaxRetryDelay
	}
	if c.MaxRetryDelay < c.RetryDelay {
		c.MaxRetryDelay = c.RetryDelay
	}
	if c.Logger == nil {
		c.Logger = NopLogger{}
	}
	return c
}

type eventInfo struct {
	event   abi.Event
	indexed abi.Arguments
}

// Watcher streams contract logs to a handler. Create one with New, then call
// Start. A Watcher must not be started twice concurrently.
type Watcher struct {
	cfg           *Config
	backend       Backend
	closeClient   func()
	closeOnce     sync.Once
	doneOnce      sync.Once
	eventsByTopic map[common.Hash]eventInfo

	mu        sync.Mutex
	status    string
	cancel    context.CancelFunc
	stop      chan struct{}
	done      chan struct{}
	lastBlock uint64
	haveLast  bool
}

// New applies Config defaults, validates cfg and returns a ready to start
// Watcher. A nil config is treated as DefaultConfig(). Defaults are written
// back into cfg, so callers can inspect the values that will be used. A client
// is dialed only when cfg.Client is nil; that client is closed by Close.
func New(cfg *Config) (*Watcher, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	cfg.WithDefaults()

	if cfg.Client == nil && strings.TrimSpace(cfg.URL) == "" {
		return nil, errors.New("events: Config requires URL or Client")
	}
	if cfg.Mode > ModePoll {
		return nil, fmt.Errorf("events: invalid Mode %d", cfg.Mode)
	}
	if cfg.PollInterval < 0 || cfg.RetryDelay < 0 || cfg.Concurrency < 0 || cfg.HandlerTimeout < 0 {
		return nil, errors.New("events: intervals, Concurrency and HandlerTimeout must not be negative")
	}

	w := &Watcher{
		cfg:           cfg,
		status:        StatusIdle,
		stop:          make(chan struct{}),
		done:          make(chan struct{}),
		eventsByTopic: make(map[common.Hash]eventInfo),
	}
	if cfg.ABI != nil {
		for _, event := range cfg.ABI.Events {
			info := eventInfo{event: event}
			for _, input := range event.Inputs {
				if input.Indexed {
					info.indexed = append(info.indexed, input)
				}
			}
			w.eventsByTopic[event.ID] = info
		}
	}
	if cfg.Client != nil {
		w.backend = cfg.Client
	} else {
		client, err := ethclient.Dial(cfg.URL)
		if err != nil {
			return nil, fmt.Errorf("events: dial %q: %w", cfg.URL, err)
		}
		w.backend = client
		w.closeClient = client.Close
	}
	return w, nil
}

// Done is closed when Start returns or when Close is called. Note that a
// handler that ignores its context can still delay Start from returning.
func (w *Watcher) Done() <-chan struct{} { return w.done }

// Close stops the watcher, unblocking Start even if its context is not
// canceled. It is safe to call Close more than once and from any goroutine.
func (w *Watcher) Close() error {
	w.closeOnce.Do(func() {
		w.mu.Lock()
		w.status = StatusClosed
		cancel := w.cancel
		w.mu.Unlock()
		close(w.stop)
		if cancel != nil {
			cancel()
		}
		w.markDone()
		if w.closeClient != nil {
			w.closeClient()
		}
	})
	return nil
}

// Status reports the current state of the watcher: StatusIdle, StatusRunning,
// StatusStopped or StatusClosed.
func (w *Watcher) Status() string {
	w.mu.Lock()
	status := strings.Clone(w.status)
	defer w.mu.Unlock()
	return status
}

// Stop is a convenience alias for Close.
func (w *Watcher) Stop() error { return w.Close() }

// Event is called by Start for every matching log. It decodes the log when an
// ABI was configured, invokes the handler and reports handler errors through
// OnError (or Logger). The handler error is returned to the caller as well;
// Start deliberately ignores it so one bad event cannot stop the stream.
func (w *Watcher) Event(ctx context.Context, raw types.Log) error {
	event := Event{Log: raw}
	if w.cfg.ABI != nil && len(raw.Topics) > 0 {
		if info, ok := w.eventsByTopic[raw.Topics[0]]; ok {
			values := make(map[string]any)
			if err := info.event.Inputs.UnpackIntoMap(values, raw.Data); err == nil &&
				abi.ParseTopicsIntoMap(values, info.indexed, raw.Topics[1:]) == nil {
				event.Name = info.event.Name
				event.Data = values
			}
		}
	}

	if w.cfg.Handler == nil {
		return nil
	}
	handlerCtx := ctx
	if w.cfg.HandlerTimeout > 0 {
		var cancel context.CancelFunc
		handlerCtx, cancel = context.WithTimeout(ctx, w.cfg.HandlerTimeout)
		defer cancel()
	}
	if err := w.cfg.Handler(handlerCtx, event); err != nil {
		w.reportError(ctx, fmt.Errorf("events: handle %s log in %s: %w",
			eventName(event, raw), raw.TxHash.Hex(), err))
		return err
	}
	return nil
}

// Start watches logs until ctx is canceled, Close is called or an
// unrecoverable error occurs. It returns nil on cancellation. Calling Start
// on an already running watcher returns ErrAlreadyStarted.
func (w *Watcher) Start(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	w.mu.Lock()
	switch w.status {
	case StatusRunning:
		w.mu.Unlock()
		return ErrAlreadyStarted
	case StatusClosed:
		w.mu.Unlock()
		w.markDone()
		return nil
	}
	runCtx, cancel := context.WithCancel(ctx)
	w.status = StatusRunning
	w.cancel = cancel
	w.mu.Unlock()

	dispatch, waitHandlers := w.newDispatcher(runCtx)
	defer func() {
		cancel()
		waitHandlers()
		w.mu.Lock()
		if w.status != StatusClosed {
			w.status = StatusStopped
		}
		w.cancel = nil
		w.mu.Unlock()
		w.markDone()
	}()

	// Close may race with Start; this goroutine makes sure it cancels runCtx.
	go func() {
		select {
		case <-w.stop:
			cancel()
		case <-runCtx.Done():
		}
	}()

	if w.subscribeMode() {
		if err := w.subscribeLoop(runCtx, dispatch); !errors.Is(err, errFallback) {
			return err
		}
	}
	return w.pollLoop(runCtx, dispatch)
}

// errFallback asks Start to switch from subscriptions to polling.
var errFallback = errors.New("events: fall back to polling")

type dispatchFunc func(context.Context, types.Log)

// newDispatcher returns the function used to hand logs to Event and a function
// that drains in-flight handlers. Concurrency 1 keeps the current inline,
// ordered behaviour; larger values use a bounded worker pool.
func (w *Watcher) newDispatcher(ctx context.Context) (dispatchFunc, func()) {
	if w.cfg.Concurrency <= 1 {
		return func(ctx context.Context, raw types.Log) { _ = w.Event(ctx, raw) }, func() {}
	}

	jobs := make(chan types.Log, w.cfg.Concurrency)
	var wg sync.WaitGroup
	for i := 0; i < w.cfg.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for raw := range jobs {
				_ = w.Event(ctx, raw)
			}
		}()
	}
	dispatch := func(ctx context.Context, raw types.Log) {
		select {
		case <-ctx.Done():
			return
		case <-w.stop:
			return
		default:
		}
		select {
		case jobs <- raw:
		case <-ctx.Done():
		case <-w.stop:
		}
	}
	return dispatch, func() {
		close(jobs)
		wg.Wait()
	}
}

func (w *Watcher) subscribeMode() bool {
	switch w.cfg.Mode {
	case ModeSubscribe:
		return true
	case ModePoll:
		return false
	default:
		// Anything that is not plain HTTP (ws, wss, IPC paths) can subscribe.
		return subscribableURL(w.cfg.URL) || w.cfg.Client != nil
	}
}

func subscribableURL(rawURL string) bool {
	u := strings.ToLower(strings.TrimSpace(rawURL))
	return u != "" && !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://")
}

func (w *Watcher) query(from *big.Int) ethereum.FilterQuery {
	q := ethereum.FilterQuery{Topics: w.cfg.Topics, FromBlock: from}
	if w.cfg.Address != (common.Address{}) {
		q.Addresses = []common.Address{w.cfg.Address}
	}
	return q
}

func (w *Watcher) subscribeLoop(ctx context.Context, dispatch dispatchFunc) error {
	attempt := 0
	for {
		if w.stopped(ctx) {
			return nil
		}
		from := w.cfg.FromBlock
		if w.haveLast {
			// Resume from the last seen block (inclusive). Redelivery is
			// preferable to missing a log when a subscription drops.
			from = new(big.Int).SetUint64(w.lastBlock)
		}
		logs := make(chan types.Log, subscriptionBuffer)
		sub, err := w.backend.SubscribeFilterLogs(ctx, w.query(from), logs)
		if err != nil {
			if w.cfg.Mode == ModeAuto && !isWS(w.cfg.URL) {
				return errFallback
			}
			w.reportError(ctx, fmt.Errorf("events: subscribe: %w", err))
			if !w.wait(ctx, w.backoff(attempt)) {
				return nil
			}
			attempt++
			continue
		}

		attempt = 0
		run := true
		for run {
			select {
			case <-ctx.Done():
				run = false
			case <-w.stop:
				run = false
			case raw, ok := <-logs:
				if !ok {
					run = false
					continue
				}
				w.noteBlock(raw.BlockNumber)
				dispatch(ctx, raw)
			case subErr, ok := <-sub.Err():
				if ok && subErr != nil {
					w.reportError(ctx, fmt.Errorf("events: subscription: %w", subErr))
				}
				run = false
			}
		}
		sub.Unsubscribe()
		if w.stopped(ctx) {
			return nil
		}
		if !w.wait(ctx, w.backoff(attempt)) {
			return nil
		}
		attempt++
	}
}

func (w *Watcher) pollLoop(ctx context.Context, dispatch dispatchFunc) error {
	var next *big.Int
	if w.cfg.FromBlock != nil {
		next = new(big.Int).Set(w.cfg.FromBlock)
	}
	attempt := 0
	for {
		if w.stopped(ctx) {
			return nil
		}

		head, err := w.backend.BlockNumber(ctx)
		if err != nil {
			w.reportError(ctx, fmt.Errorf("events: block number: %w", err))
			if !w.wait(ctx, w.backoff(attempt)) {
				return nil
			}
			attempt++
			continue
		}
		if next == nil {
			// First successful poll starts at the current head.
			next = new(big.Int).SetUint64(head)
		}

		to := new(big.Int).SetUint64(head)
		if next.Cmp(to) <= 0 {
			logs, filterErr := w.backend.FilterLogs(ctx, w.queryRange(next, to))
			if filterErr != nil {
				w.reportError(ctx, fmt.Errorf("events: filter logs: %w", filterErr))
				if !w.wait(ctx, w.backoff(attempt)) {
					return nil
				}
				attempt++
				continue
			}
			for _, raw := range logs {
				if w.stopped(ctx) {
					return nil
				}
				w.noteBlock(raw.BlockNumber)
				dispatch(ctx, raw)
			}
			next = new(big.Int).Add(to, big.NewInt(1))
		}
		attempt = 0
		if !w.wait(ctx, w.cfg.PollInterval) {
			return nil
		}
	}
}

func (w *Watcher) queryRange(from, to *big.Int) ethereum.FilterQuery {
	q := w.query(from)
	q.ToBlock = to
	return q
}

func (w *Watcher) noteBlock(block uint64) {
	if !w.haveLast || block > w.lastBlock {
		w.lastBlock = block
		w.haveLast = true
	}
}

func (w *Watcher) reportError(ctx context.Context, err error) {
	if w.cfg.OnError != nil {
		w.cfg.OnError(ctx, err)
		return
	}
	w.cfg.Logger.Error(err)
}

// wait blocks for d, or until the watcher/context is canceled. It returns
// false when canceled.
func (w *Watcher) wait(ctx context.Context, d time.Duration) bool {
	if d <= 0 {
		select {
		case <-ctx.Done():
			return false
		case <-w.stop:
			return false
		default:
			return true
		}
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-w.stop:
		return false
	case <-timer.C:
		return true
	}
}

func (w *Watcher) backoff(attempt int) time.Duration {
	d := w.cfg.RetryDelay
	for i := 0; i < attempt; i++ {
		if d >= w.cfg.MaxRetryDelay/2 {
			return w.cfg.MaxRetryDelay
		}
		d *= 2
	}
	if d > w.cfg.MaxRetryDelay {
		return w.cfg.MaxRetryDelay
	}
	return d
}

func (w *Watcher) markDone() { w.doneOnce.Do(func() { close(w.done) }) }

func (w *Watcher) stopped(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
	}
	select {
	case <-w.stop:
		return true
	default:
		return false
	}
}

func isWS(rawURL string) bool {
	u := strings.ToLower(rawURL)
	return strings.HasPrefix(u, "ws://") || strings.HasPrefix(u, "wss://")
}

func eventName(event Event, raw types.Log) string {
	if event.Name != "" {
		return event.Name
	}
	if len(raw.Topics) > 0 {
		return raw.Topics[0].Hex()
	}
	return "log"
}
