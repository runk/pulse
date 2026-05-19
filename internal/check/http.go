package check

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/runk/pulse/internal/assertion"
)

type HttpAssertion struct {
	StatusCode *assertion.NumberMatcher            `json:"statusCode,omitempty"`
	Body       *assertion.StringMatcher            `json:"body,omitempty"`
	Headers    map[string]*assertion.StringMatcher `json:"headers,omitempty"`
}

type HTTPCheck struct {
	URL        string          `json:"url"`
	Method     string          `json:"method,omitempty"`
	Body       []byte          `json:"body,omitempty"`
	Assertions []HttpAssertion `json:"assertions,omitempty"`
}

func (HTTPCheck) Type() string { return "http" }

func (c *HTTPCheck) Validate() error {
	// URL validation
	if c.URL == "" {
		return errors.New("URL cannot be blank")
	}

	parsed, err := url.Parse(c.URL)
	if err != nil {
		return err
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("Only http and https schemes supported, got: '%s'", parsed.Scheme)
	}

	// Method validation
	c.Method = strings.ToUpper(c.Method)
	if c.Method == "" {
		c.Method = http.MethodGet
	}

	switch c.Method {
	case http.MethodDelete,
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPut,
		http.MethodTrace:
		// valid
	default:
		return fmt.Errorf("Unsupported method: '%s'", c.Method)
	}

	return nil
}

func (c HTTPCheck) Run() error {
	body := bytes.NewReader(c.Body)

	req, err := http.NewRequest(c.Method, c.URL, body)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	status := res.StatusCode

	fmt.Printf("%s: %d\n", c.URL, res.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if status < 200 || status >= 300 {
		return fmt.Errorf("%s returned non 2xx status: %d", c.URL, status)
	}

	errs := []error{}
	for _, assertion := range c.Assertions {
		if assertion.StatusCode != nil {
			err := assertion.StatusCode.Match(status)
			if err != nil {
				errs = append(errs, err)
			}
		}

		if assertion.Body != nil {
			err := assertion.Body.Match(&resBody)
			if err != nil {
				errs = append(errs, err)
			}
		}

		if assertion.Headers != nil {
			for key, matcher := range assertion.Headers {
				err = matcher.Match(res.Header.Get(key))
				if err != nil {
					errs = append(errs, err)
				}
			}
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
