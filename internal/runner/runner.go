package runner

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/runk/pulse/internal/check"
)

func Execute(checks []check.Check, concurrency uint16, timeout uint32) error {
	if timeout < 10 {
		return fmt.Errorf("Timeout should be >= 10ms, got: %dms", timeout)
	}

	if concurrency < 1 {
		return fmt.Errorf("Concurrency should be at least 1, got: %d", concurrency)
	}

	sem := make(chan struct{}, concurrency)
	errCh := make(chan error, len(checks))
	wg := sync.WaitGroup{}

	for _, check := range checks {
		sem <- struct{}{}

		wg.Go(func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)

			defer cancel()
			defer func() { <-sem }()

			if err := check.Value.Run(ctx); err != nil {
				errCh <- err
			}
		})
	}

	wg.Wait()
	close(errCh)

	fmt.Println("")

	errored := false
	for err := range errCh {
		errored = true
		// at least 1
		fmt.Printf("Error: %s\n", err)
	}

	if errored {
		fmt.Println("Policy execution completed with errors.")
		os.Exit(1)
	} else {
		fmt.Println("Policy execution completed - all checks passed.")
		os.Exit(0)
	}

	return nil
}
