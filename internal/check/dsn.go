package check

import "fmt"

type DNSCheck struct {
	Host string `json:"host"`
}

func (DNSCheck) Type() string { return "dns" }

func (c DNSCheck) Run() error {
	fmt.Println("DNS check", c.Host)
	return nil
}

func (c DNSCheck) Validate() error {
	return nil
}
