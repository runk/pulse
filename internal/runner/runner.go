package runner

import (
	"fmt"
	"os"
	"sync"

	"github.com/runk/pulse/internal/check"
)

func Execute(checks []check.Check, concurrency uint16) {
	fmt.Println(checks)

	if concurrency == 0 {
		// Otherwise we would have channel with no capacity
		concurrency = 1
	}

	sem := make(chan struct{}, concurrency)
	errCh := make(chan error, len(checks))
	wg := sync.WaitGroup{}

	for _, check := range checks {
		sem <- struct{}{}

		wg.Go(func() {
			defer func() { <-sem }()

			if err := check.Value.Run(); err != nil {
				errCh <- err
			}
		})
	}

	wg.Wait()
	close(errCh)

	errored := false
	for err := range errCh {
		errored = true
		// at least 1
		fmt.Printf("Error: %s\n", err)
	}

	if errored {
		fmt.Println("Policy execution completed with errors.")
		os.Exit(1)
	}
}
