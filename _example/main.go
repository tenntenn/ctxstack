package main

import (
	"context"
	"fmt"
	"time"

	"github.com/tenntenn/ctxstack"
)

func main() {
	ctx, cancel := ctxstack.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := f(ctx); err != nil {
		fmt.Println(err, string(ctxstack.Stack(ctx)))
	}
}

func f(ctx context.Context) error {
	ticker := time.NewTicker(2 * time.Second)
	select {
	case <-ticker.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}
