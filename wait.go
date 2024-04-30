package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

// ConnectWaiter provides a means of checking whether something (e.g. a database connection) is ready/available.
type ConnectWaiter interface {
	// Connect is called when we need to wait until a service is ready.
	// It should return when the service is ready, the context expires or an error occurs.
	// A non-nil error should only be returned when readiness can not be checked (not when the context expires).
	Connect(ctx context.Context) error
}

func waitFor(waiter ConnectWaiter, timeout time.Duration, logger *log.Logger) error {
	logger.Printf("waiting %s until %T is ready", timeout, waiter)
	ctx, cancel := context.WithTimeoutCause(
		context.Background(),
		timeout,
		fmt.Errorf("timeout of %s expired", timeout),
	)
	defer cancel()

	err := waiter.Connect(ctx)
	errCtx := ctx.Err()
	return errors.Join(err, errCtx)
}
