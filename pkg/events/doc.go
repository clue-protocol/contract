// Package events streams event logs emitted by an EVM contract into a caller
// supplied handler.
//
// # Main API
//
// A Watcher is built and run through three functions:
//
//   - New builds a watcher from a *Config (RPC endpoint or Backend, contract
//     address, topic filters and a handler). Unset fields get defaults from
//     Config.WithDefaults; New(nil) uses DefaultConfig.
//   - Event is the dispatch hook called for each incoming log. It optionally
//     decodes the log with Config.ABI and then invokes Config.Handler.
//   - Start runs the watcher until its context is canceled or Close/Stop is
//     called. Status reports idle, running, stopped or closed.
//
// # Transports and delivery
//
// Websocket (ws/wss) and IPC endpoints use eth_subscribe for real-time
// delivery; HTTP endpoints are polled with eth_getLogs. A supplied Backend is
// tried with a subscription first and falls back to polling. Failed RPC calls
// are retried with exponential backoff. Delivery is at-least-once, so a
// reconnect can replay a log and handlers should be idempotent.
//
// # Cancellation and state
//
// Start returns nil when its context is canceled. Close (or its alias Stop)
// cancels an in-flight Start from another goroutine and closes Done
// immediately; Start waits for in-flight handlers before returning. Handler
// contexts are canceled on shutdown, and Config.HandlerTimeout can bound a
// single handler.
//
// # Concurrency and logging
//
// Handlers run sequentially by default, preserving log order. Config.Concurrency
// above 1 uses a bounded worker pool, allowing fresh events to be handled while
// an older handler is still busy, at the cost of completion ordering.
//
// Config.OnError receives non-fatal errors; otherwise they go to Config.Logger,
// which implements Printf and Error. logger.go provides StdLogger, SlogLogger,
// LogAdapter and NopLogger.
//
// # Example
//
//	w, err := events.New(&events.Config{
//	    URL:     "wss://your-node/ws",
//	    Address: common.HexToAddress("0xYourContract"),
//	    Handler: func(ctx context.Context, event events.Event) error {
//	        fmt.Println(event.BlockNumber(), event.TxHash())
//	        return nil
//	    },
//	})
//	if err != nil {
//	    return err
//	}
//	defer w.Close()
//	return w.Start(ctx)
package events
