package events

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"math/big"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

const pingABI = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Ping","type":"event"}]`

// fakeBackend is an in-memory Backend used to exercise Start without a node.
type fakeBackend struct {
	mu sync.Mutex

	head        uint64
	logs        []types.Log
	blockErr    error
	filterErr   error
	subErr      error
	blockCalls  int
	filterCalls int
	subCalls    int
	lastQuery   ethereum.FilterQuery

	onBlock func()
	onSub   func(context.Context, ethereum.FilterQuery, chan<- types.Log) (ethereum.Subscription, error)
}

func (f *fakeBackend) BlockNumber(context.Context) (uint64, error) {
	f.mu.Lock()
	f.blockCalls++
	err := f.blockErr
	f.mu.Unlock()
	if f.onBlock != nil {
		f.onBlock()
	}
	return f.head, err
}

func (f *fakeBackend) FilterLogs(_ context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.filterCalls++
	f.lastQuery = query
	return f.logs, f.filterErr
}

func (f *fakeBackend) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	f.mu.Lock()
	f.subCalls++
	f.lastQuery = query
	err := f.subErr
	f.mu.Unlock()
	if err != nil {
		return nil, err
	}
	if f.onSub != nil {
		return f.onSub(ctx, query, ch)
	}
	return newFakeSubscription(), nil
}

func (f *fakeBackend) counts() (blocks, filters, subs int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.blockCalls, f.filterCalls, f.subCalls
}

type fakeSubscription struct {
	errCh     chan error
	unsubOnce sync.Once
	unsubCh   chan struct{}
}

func newFakeSubscription() *fakeSubscription {
	return &fakeSubscription{errCh: make(chan error, 1), unsubCh: make(chan struct{})}
}

func (s *fakeSubscription) Err() <-chan error { return s.errCh }

func (s *fakeSubscription) Unsubscribe() {
	s.unsubOnce.Do(func() {
		close(s.unsubCh)
		close(s.errCh)
	})
}

func TestNew(t *testing.T) {
	_, err := New(&Config{})
	require.ErrorContains(t, err, "requires URL or Client")

	_, err = New(&Config{URL: "://not-a-url"})
	require.Error(t, err)

	_, err = New(&Config{Client: &fakeBackend{}, Mode: Mode(9)})
	require.ErrorContains(t, err, "invalid Mode")

	w, err := New(&Config{
		Client:     &fakeBackend{},
		Address:    common.HexToAddress("0x01"),
		RetryDelay: 50 * time.Millisecond,
	})
	require.NoError(t, err)
	require.Equal(t, DefaultPollInterval, w.cfg.PollInterval)
	require.Equal(t, DefaultMaxRetryDelay, w.cfg.MaxRetryDelay)
	require.True(t, w.subscribeMode(), "a supplied client is tried with a subscription first")
}

func TestEventDispatchesAndDecodesLog(t *testing.T) {
	parsed, err := abi.JSON(strings.NewReader(pingABI))
	require.NoError(t, err)

	from := common.HexToAddress("0xabc")
	data, err := parsed.Events["Ping"].Inputs.NonIndexed().Pack(big.NewInt(7))
	require.NoError(t, err)

	raw := types.Log{
		Address:     common.HexToAddress("0x01"),
		Topics:      []common.Hash{parsed.Events["Ping"].ID, common.BytesToHash(from.Bytes())},
		Data:        data,
		BlockNumber: 42,
		TxHash:      common.HexToHash("0x99"),
	}

	var got Event
	w, err := New(&Config{
		Client: &fakeBackend{},
		ABI:    &parsed,
		Handler: func(_ context.Context, event Event) error {
			got = event
			return nil
		},
	})
	require.NoError(t, err)

	require.NoError(t, w.Event(context.Background(), raw))
	require.Equal(t, raw, got.Log)
	require.Equal(t, "Ping", got.Name)
	require.Equal(t, big.NewInt(7), got.Data["value"])
	require.Equal(t, from, got.Data["from"])
	require.Equal(t, uint64(42), got.BlockNumber())
	require.Equal(t, raw.TxHash, got.TxHash())
	require.Equal(t, raw.Address, got.Address())
	require.Equal(t, raw.Topics, got.Topics())
	require.Equal(t, data, got.DataBytes())
	require.False(t, got.Removed())
}

func TestEventReportsHandlerError(t *testing.T) {
	handlerErr := errors.New("boom")
	var reported error
	w, err := New(&Config{
		Client:  &fakeBackend{},
		Handler: func(context.Context, Event) error { return handlerErr },
		OnError: func(_ context.Context, err error) { reported = err },
	})
	require.NoError(t, err)

	err = w.Event(context.Background(), types.Log{})
	require.ErrorIs(t, err, handlerErr)
	require.ErrorIs(t, reported, handlerErr)
}

func TestStartPollsUntilContextCanceled(t *testing.T) {
	contract := common.HexToAddress("0x1234")
	backend := &fakeBackend{
		head: 11,
		logs: []types.Log{{Address: contract, BlockNumber: 11}},
	}
	delivered := make(chan Event, 4)
	w, err := New(&Config{
		Client:       backend,
		Address:      contract,
		Mode:         ModePoll,
		PollInterval: 5 * time.Millisecond,
		Handler: func(_ context.Context, event Event) error {
			delivered <- event
			return nil
		},
	})
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() { errCh <- w.Start(ctx) }()

	select {
	case event := <-delivered:
		require.Equal(t, contract, event.Address())
	case <-time.After(2 * time.Second):
		t.Fatal("no event delivered")
	}
	cancel()
	require.NoError(t, waitStart(errCh))

	_, filters, _ := backend.counts()
	require.NotZero(t, filters)
	require.Equal(t, big.NewInt(11), backend.lastQuery.FromBlock)
	require.Equal(t, big.NewInt(11), backend.lastQuery.ToBlock)
	require.Equal(t, []common.Address{contract}, backend.lastQuery.Addresses)
}

func TestStartSubscriptionStopsOnCancel(t *testing.T) {
	backend := &fakeBackend{}
	var sub *fakeSubscription
	backend.onSub = func(_ context.Context, _ ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
		sub = newFakeSubscription()
		ch <- types.Log{BlockNumber: 7}
		return sub, nil
	}

	delivered := make(chan struct{}, 1)
	w, err := New(&Config{
		Client:       backend,
		Mode:         ModeSubscribe,
		PollInterval: time.Millisecond,
		Handler: func(context.Context, Event) error {
			delivered <- struct{}{}
			return nil
		},
	})
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() { errCh <- w.Start(ctx) }()

	select {
	case <-delivered:
	case <-time.After(2 * time.Second):
		t.Fatal("no subscribed event delivered")
	}
	cancel()
	require.NoError(t, waitStart(errCh))

	_, _, subs := backend.counts()
	require.Equal(t, 1, subs)
	select {
	case <-sub.unsubCh:
	default:
		t.Fatal("subscription was not unsubscribed on cancel")
	}
}

func TestCloseStopsStart(t *testing.T) {
	blocking := make(chan struct{})
	backend := &fakeBackend{}
	backend.onBlock = func() {
		select {
		case <-blocking:
		default:
			close(blocking)
		}
	}

	w, err := New(&Config{Client: backend, Mode: ModePoll})
	require.NoError(t, err)

	errCh := make(chan error, 1)
	go func() { errCh <- w.Start(context.Background()) }()
	select {
	case <-blocking:
	case <-time.After(2 * time.Second):
		t.Fatal("Start did not poll")
	}

	require.NoError(t, w.Close())
	require.NoError(t, waitStart(errCh))
	select {
	case <-w.Done():
	default:
		t.Fatal("Done was not closed")
	}
	require.NoError(t, w.Close(), "Close is idempotent")
	require.NoError(t, w.Start(context.Background()), "a closed watcher does not restart")
}

func TestStartTwice(t *testing.T) {
	blocking := make(chan struct{})
	backend := &fakeBackend{}
	backend.onBlock = func() {
		select {
		case <-blocking:
		default:
			close(blocking)
		}
	}

	w, err := New(&Config{Client: backend, Mode: ModePoll})
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() { errCh <- w.Start(ctx) }()
	select {
	case <-blocking:
	case <-time.After(2 * time.Second):
		t.Fatal("Start did not poll")
	}
	require.ErrorIs(t, w.Start(ctx), ErrAlreadyStarted)
	cancel()
	require.NoError(t, waitStart(errCh))
}

func TestSubscribeFallsBackToPolling(t *testing.T) {
	backend := &fakeBackend{
		head:     3,
		subErr:   errors.New("notifications not supported"),
		blockErr: errors.New("node down"),
	}
	polled := make(chan struct{}, 1)
	backend.onBlock = func() {
		select {
		case polled <- struct{}{}:
		default:
		}
	}

	w, err := New(&Config{Client: backend, RetryDelay: time.Millisecond})
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() { errCh <- w.Start(ctx) }()
	select {
	case <-polled:
	case <-time.After(2 * time.Second):
		t.Fatal("did not fall back to polling")
	}
	cancel()
	require.NoError(t, waitStart(errCh))
	_, filters, subs := backend.counts()
	require.Equal(t, 1, subs)
	require.Zero(t, filters, "BlockNumber failed, so FilterLogs is never reached")
}

func TestBackoffIsCapped(t *testing.T) {
	w, err := New(&Config{
		Client:        &fakeBackend{},
		RetryDelay:    time.Second,
		MaxRetryDelay: 8 * time.Second,
	})
	require.NoError(t, err)
	require.Equal(t, time.Second, w.backoff(0))
	require.Equal(t, 2*time.Second, w.backoff(1))
	require.Equal(t, 8*time.Second, w.backoff(3))
	require.Equal(t, 8*time.Second, w.backoff(100))
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	require.Equal(t, ModeAuto, cfg.Mode)
	require.Equal(t, DefaultPollInterval, cfg.PollInterval)
	require.Equal(t, DefaultRetryDelay, cfg.RetryDelay)
	require.Equal(t, DefaultMaxRetryDelay, cfg.MaxRetryDelay)
	require.Equal(t, 1, cfg.Concurrency)
	require.Zero(t, cfg.HandlerTimeout)
	require.NotNil(t, cfg.Logger)
	require.IsType(t, NopLogger{}, cfg.Logger)

	// New(nil) uses the defaults, but still needs a source of logs.
	_, err := New(nil)
	require.ErrorContains(t, err, "requires URL or Client")

	_, err = New(&Config{Client: &fakeBackend{}, Concurrency: -1})
	require.ErrorContains(t, err, "must not be negative")
	_, err = New(&Config{Client: &fakeBackend{}, HandlerTimeout: -time.Second})
	require.ErrorContains(t, err, "must not be negative")

	// Defaults are written back onto the caller's config.
	cfg = &Config{Client: &fakeBackend{}, RetryDelay: time.Minute}
	cfg.WithDefaults()
	require.Equal(t, time.Minute, cfg.MaxRetryDelay, "max backoff is raised to match retry delay")
	w, err := New(cfg)
	require.NoError(t, err)
	require.Same(t, cfg, w.cfg)
	require.Equal(t, DefaultPollInterval, cfg.PollInterval)
	require.Equal(t, StatusIdle, w.Status())
}

func TestStatusTransitions(t *testing.T) {
	polled := make(chan struct{}, 1)
	backend := &fakeBackend{}
	backend.onBlock = func() {
		select {
		case polled <- struct{}{}:
		default:
		}
	}
	w, err := New(&Config{Client: backend, Mode: ModePoll})
	require.NoError(t, err)
	require.Equal(t, StatusIdle, w.Status())

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() { errCh <- w.Start(ctx) }()
	select {
	case <-polled:
	case <-time.After(2 * time.Second):
		t.Fatal("Start did not poll")
	}
	require.Equal(t, StatusRunning, w.Status())
	require.ErrorIs(t, w.Start(ctx), ErrAlreadyStarted)

	cancel()
	require.NoError(t, waitStart(errCh))
	require.Equal(t, StatusStopped, w.Status())

	ctx2, cancel2 := context.WithCancel(context.Background())
	errCh2 := make(chan error, 1)
	go func() { errCh2 <- w.Start(ctx2) }()
	require.Eventually(t, func() bool { return w.Status() == StatusRunning }, time.Second, 5*time.Millisecond)
	require.NoError(t, w.Stop())
	require.NoError(t, waitStart(errCh2))
	require.Equal(t, StatusClosed, w.Status())
	require.NoError(t, w.Start(context.Background()))
	require.Equal(t, StatusClosed, w.Status())
	cancel2()
}

type recordingLogger struct {
	mu   sync.Mutex
	msgs []string
	errs []string
}

func (l *recordingLogger) Printf(format string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.msgs = append(l.msgs, fmt.Sprintf(format, args...))
}

func (l *recordingLogger) Error(args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.errs = append(l.errs, fmt.Sprint(args...))
}

func (l *recordingLogger) messages() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]string(nil), l.msgs...)
}

func (l *recordingLogger) errors() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]string(nil), l.errs...)
}

var (
	_ Logger = (*recordingLogger)(nil)
	_ Logger = LogAdapter{}
	_ Logger = NopLogger{}
	_ Logger = SlogLogger{}
	_ Logger = (*StdLogger)(nil)
)

func TestLoggerImplementations(t *testing.T) {
	handler := func(context.Context, Event) error { return errors.New("handler failed") }

	// A Printf-only logger works through LogAdapter.
	var stdlib bytes.Buffer
	w, err := New(&Config{
		Client:  &fakeBackend{},
		Logger:  LogAdapter{Printfer: log.New(&stdlib, "", 0)},
		Handler: handler,
	})
	require.NoError(t, err)
	require.ErrorContains(t, w.Event(context.Background(), types.Log{}), "handler failed")
	require.Contains(t, stdlib.String(), "handler failed")

	// Errors go to the Error method, not Printf.
	logger := &recordingLogger{}
	w, err = New(&Config{Client: &fakeBackend{}, Logger: logger, Handler: handler})
	require.NoError(t, err)
	require.ErrorContains(t, w.Event(context.Background(), types.Log{}), "handler failed")
	require.Len(t, logger.errors(), 1)
	require.Empty(t, logger.messages())

	// slog is usable through the provided adapter.
	var slogBuf bytes.Buffer
	w, err = New(&Config{
		Client:  &fakeBackend{},
		Logger:  SlogLogger{Logger: slog.New(slog.NewTextHandler(&slogBuf, nil))},
		Handler: handler,
	})
	require.NoError(t, err)
	require.ErrorContains(t, w.Event(context.Background(), types.Log{}), "handler failed")
	require.Contains(t, slogBuf.String(), "handler failed")

	// StdLogger tags errors with an ERROR prefix.
	var stdBuf bytes.Buffer
	w, err = New(&Config{
		Client:  &fakeBackend{},
		Logger:  NewStdLogger(&stdBuf),
		Handler: handler,
	})
	require.NoError(t, err)
	require.ErrorContains(t, w.Event(context.Background(), types.Log{}), "handler failed")
	require.Contains(t, stdBuf.String(), "ERROR")
	require.Contains(t, stdBuf.String(), "handler failed")
}

func TestConcurrentHandlers(t *testing.T) {
	const count = 4
	backend := &fakeBackend{}
	backend.onSub = func(_ context.Context, _ ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
		for i := 0; i < count; i++ {
			ch <- types.Log{BlockNumber: uint64(i), TxHash: common.BytesToHash([]byte{byte(i + 1)})}
		}
		return newFakeSubscription(), nil
	}

	started := make(chan struct{}, count)
	release := make(chan struct{})
	var handled atomic.Int64
	w, err := New(&Config{
		Client:      backend,
		Mode:        ModeSubscribe,
		Concurrency: count,
		Handler: func(ctx context.Context, _ Event) error {
			handled.Add(1)
			started <- struct{}{}
			select {
			case <-release:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		},
	})
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() { errCh <- w.Start(ctx) }()
	for i := 0; i < count; i++ {
		select {
		case <-started:
		case <-time.After(2 * time.Second):
			t.Fatalf("only %d handlers started concurrently", i)
		}
	}
	close(release)
	cancel()
	require.NoError(t, waitStart(errCh))
	require.EqualValues(t, count, handled.Load())
}

func TestSequentialHandlersKeepOrder(t *testing.T) {
	hashes := []common.Hash{
		common.BytesToHash([]byte{1}),
		common.BytesToHash([]byte{2}),
		common.BytesToHash([]byte{3}),
	}
	backend := &fakeBackend{head: 3, logs: []types.Log{
		{TxHash: hashes[0]},
		{TxHash: hashes[1]},
		{TxHash: hashes[2]},
	}}

	var mu sync.Mutex
	var order []common.Hash
	w, err := New(&Config{
		Client:       backend,
		Mode:         ModePoll,
		PollInterval: time.Millisecond,
		Handler: func(_ context.Context, event Event) error {
			mu.Lock()
			order = append(order, event.TxHash())
			mu.Unlock()
			return nil
		},
	})
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() { errCh <- w.Start(ctx) }()
	require.Eventually(t, func() bool {
		mu.Lock()
		defer mu.Unlock()
		return len(order) == len(hashes)
	}, 2*time.Second, 5*time.Millisecond)
	cancel()
	require.NoError(t, waitStart(errCh))

	mu.Lock()
	defer mu.Unlock()
	require.Equal(t, hashes, order)
}

func TestHandlerTimeout(t *testing.T) {
	var reported error
	w, err := New(&Config{
		Client:         &fakeBackend{},
		HandlerTimeout: 10 * time.Millisecond,
		Handler: func(ctx context.Context, _ Event) error {
			<-ctx.Done()
			return ctx.Err()
		},
		OnError: func(_ context.Context, err error) { reported = err },
	})
	require.NoError(t, err)

	err = w.Event(context.Background(), types.Log{})
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.ErrorIs(t, reported, context.DeadlineExceeded)
}

func waitStart(errCh <-chan error) error {
	select {
	case err := <-errCh:
		return err
	case <-time.After(2 * time.Second):
		return errors.New("Start did not return")
	}
}
