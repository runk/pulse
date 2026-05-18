package runner

import (
	"fmt"
	"os"
	"sync"

	"github.com/runk/pulse/internal/check"
)

func Execute(checks []check.Check, concurrency uint16) {
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
}
