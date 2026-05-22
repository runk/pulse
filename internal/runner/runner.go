package runner

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/runk/pulse/internal/check"
)

type Result struct {
	Type    string
	Subject string
	Ok      bool
	Message string
}

func Execute(checks []check.Check, results chan Result, concurrency uint16, timeout uint32) error {
	if timeout < 10 {
		return fmt.Errorf("Timeout should be >= 10ms, got: %dms", timeout)
	}

	if concurrency < 1 {
		return fmt.Errorf("Concurrency should be at least 1, got: %d", concurrency)
	}

	sem := make(chan struct{}, concurrency)
	wg := sync.WaitGroup{}

	for _, check := range checks {
		sem <- struct{}{}

		wg.Go(func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)

			defer cancel()
			defer func() { <-sem }()

			if err := check.Value.Run(ctx); err != nil {
				results <- Result{
					Ok:      false,
					Type:    check.Value.Type(),
					Subject: check.Value.Subject(),
					Message: err.Error(),
				}
			} else {
				results <- Result{
					Ok:      true,
					Type:    check.Value.Type(),
					Subject: check.Value.Subject(),
				}
			}
		})
	}

	wg.Wait()

	return nil
}
