package runner

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/runk/pulse/internal/policy"
)

func Execute(checks []policy.Check, concurrency uint16) {
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

			runCheck(&check, errCh)
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

func runCheck(check *policy.Check, errCh chan error) {
	var runner func(check *policy.Check) error

	switch check.Type {
	case "http-check":
		runner = runHttpCheck
	}

	if runner == nil {
		errCh <- errors.New("Unsupported check type")
		return
	}

	if err := runner(check); err != nil {
		errCh <- err
	}
}

func runHttpCheck(check *policy.Check) error {

	res, err := http.Get(check.Url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	status := res.StatusCode

	fmt.Printf("%s: %d\n", check.Url, res.StatusCode)

	_, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if status < 200 || status >= 300 {
		return fmt.Errorf("%s returned non 2xx status: %d", check.Url, status)
	}

	return nil
}
