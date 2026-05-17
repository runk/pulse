package check

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CheckValue interface {
	Type() string
	Run() error
}

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

type DNSCheck struct {
	Host string `json:"host"`
}

func (DNSCheck) Type() string { return "dns" }
func (c DNSCheck) Run() error {
	fmt.Println("DNS check", c.Host)
	return nil
}

type Check struct {
	Value CheckValue
}

func (rc *Check) UnmarshalJSON(data []byte) error {

	var header struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal(data, &header); err != nil {
		return err
	}

	switch header.Type {
	case "http":
		value := HTTPCheck{}
		if err := json.Unmarshal(data, &value); err != nil {
			return err
		}

		rc.Value = value
	default:
		return fmt.Errorf("Unsupported check type: '%s'", header.Type)
	}

	return nil
}
