package ctxstack

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"time"
)

type holder struct {
	mu    sync.RWMutex
	stack []byte
}

type ctxkey struct{}

func WithDeadlineCause(parent context.Context, d time.Time, cause error) (context.Context, context.CancelFunc) {

	ctx, cancel := context.WithDeadlineCause(parent, d, cause)

	v := new(holder)
	ctx = context.WithValue(ctx, ctxkey{}, v)

	stop := context.AfterFunc(ctx, func() {

		if !errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return
		}

		v.mu.Lock()
		var stack [16 * 1024]byte
		n := runtime.Stack(stack[:], true)
		v.stack = stack[:n]
		v.mu.Unlock()
	})

	return ctx, func() {
		stop()
		cancel()
	}
}

func WithDeadline(parent context.Context, d time.Time) (context.Context, context.CancelFunc) {
	return WithDeadlineCause(parent, d, nil)
}

func WithTimeoutCause(parent context.Context, timeout time.Duration, cause error) (context.Context, context.CancelFunc) {
	return WithDeadlineCause(parent, time.Now().Add(timeout), cause)
}

func WithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return WithDeadline(parent, time.Now().Add(timeout))
}

func Stack(ctx context.Context) []byte {

	v, _ := ctx.Value(ctxkey{}).(*holder)
	if v == nil {
		return nil
	}

	v.mu.RLock()
	defer v.mu.RUnlock()
	if v.stack != nil {
		return v.stack
	}

	return nil
}
