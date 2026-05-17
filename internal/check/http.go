package check

import (
	"fmt"
	"io"
	"net/http"
)

type HTTPCheck struct {
	URL string `json:"url"`
}

func (HTTPCheck) Type() string { return "http" }
func (c HTTPCheck) Run() error {
	res, err := http.Get(c.URL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	status := res.StatusCode

	fmt.Printf("%s: %d\n", c.URL, res.StatusCode)

	_, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if status < 200 || status >= 300 {
		return fmt.Errorf("%s returned non 2xx status: %d", c.URL, status)
	}

	return nil
}
