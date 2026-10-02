// Package main provides a minimal worker entry point for hellnet-worker-template.
// Replace the work function with your actual background job logic
// (e.g., Kafka consumer, SQS polling, scheduled tasks, etc.).
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/guilhermelinosp/hellnet-lib-telemetry/telemetry"
)

// tel holds the observability instance. It is nil when main is not executed
// (e.g. during unit tests), in which case jobs run without instrumentation.
var tel *telemetry.Telemetry

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t, err := telemetry.New(ctx)
	if err != nil {
		log.Fatalf("failed to init telemetry: %v", err)
	}
	tel = t
	defer func() {
		// flush with its own deadline: ctx is already cancelled on shutdown
		flushCtx, flushCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer flushCancel()
		_ = tel.Close(flushCtx)
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})
	go func() {
		defer close(done)
		workerLoop(ctx)
	}()

	tel.Log(ctx).Info("worker started")

	select {
	case <-sig:
		tel.Log(ctx).Info("shutdown signal received, stopping worker")
		cancel()
	case <-done:
		tel.Log(ctx).Info("worker finished")
	}

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		tel.Log(ctx).Error("worker shutdown timed out")
	}
}

func workerLoop(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			if tel != nil {
				tel.Log(ctx).Info("worker loop exiting")
			}
			return
		case t := <-ticker.C:
			runJob(ctx, t)
		}
	}
}

// runJob executes one tick of work with observability (when tel is available).
func runJob(ctx context.Context, t time.Time) {
	do := func(c context.Context) error { return doWork(t) }
	if tel != nil {
		_ = tel.WorkerContext(ctx, "tick", do)
		return
	}
	_ = do(ctx)
}

func doWork(t time.Time) error {
	_, _ = fmt.Printf("tick at %s\n", t.Format(time.RFC3339))
	return nil
}
