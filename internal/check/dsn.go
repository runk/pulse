package check

import (
	"errors"
	"fmt"
	"net"

	"github.com/runk/pulse/internal/assertion"
)

type DNSAssertion struct {
	CNAME *assertion.StringMatcher     `json:"cname,omitempty"`
	MX    *assertion.StringListMatcher `json:"mx,omitempty"`
	TXT   *assertion.StringListMatcher `json:"txt,omitempty"`
	NS    *assertion.StringListMatcher `json:"ns,omitempty"`
	A     *assertion.StringListMatcher `json:"a,omitempty"`
}

type DNSCheck struct {
	Host       string         `json:"host"`
	Assertions []DNSAssertion `json:"assertions"`
}

func (DNSCheck) Type() string { return "dns" }

func (c DNSCheck) Run() error {

	host := c.Host
	errs := []error{}

	for _, assertion := range c.Assertions {
		if assertion.CNAME != nil {
			cname, err := net.LookupCNAME(host)
			if err != nil {
				errs = append(errs, err)
			}
			if err = assertion.CNAME.Match(cname); err != nil {
				errs = append(errs, err)
			}
		}

		checkRecords(host, net.LookupMX, assertion.MX, &errs)
		checkRecords(host, net.LookupTXT, assertion.TXT, &errs)
		checkRecords(host, net.LookupNS, assertion.NS, &errs)
		checkRecords(host, net.LookupAddr, assertion.A, &errs)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (c DNSCheck) Validate() error {
	if c.Host == "" {
		return errors.New("host is required to perform dns check")
	}

	return nil
}

func checkRecords[T any](
	host string,
	lookup func(host string) ([]T, error),
	assert *assertion.StringListMatcher,
	errs *[]error,
) {
	if assert == nil {
		return
	}

	records, err := lookup(host)
	if err != nil {
		*errs = append(*errs, err)
	}

	inputs := make([]string, len(records))

	for i, record := range records {
		switch v := any(record).(type) {
		case string:
			inputs[i] = v
		case *net.MX:
			inputs[i] = v.Host
		case *net.NS:
			inputs[i] = v.Host
		default:
			*errs = append(*errs, fmt.Errorf("unsupported record type: %T", record))
			return
		}
	}

	if err = assert.Match(inputs); err != nil {
		*errs = append(*errs, err)
	}
}
