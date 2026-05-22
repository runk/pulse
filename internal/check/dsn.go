package check

import (
	"context"
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

func (c DNSCheck) Run(ctx context.Context) error {

	host := c.Host
	errs := []error{}

	resolver := net.Resolver{}

	for _, assertion := range c.Assertions {
		if assertion.CNAME != nil {
			cname, err := resolver.LookupCNAME(ctx, host)
			if err != nil {
				errs = append(errs, err)
			}
			if err = assertion.CNAME.Match(cname); err != nil {
				errs = append(errs, err)
			}
		}

		checkRecords(ctx, host, resolver.LookupMX, assertion.MX, &errs)
		checkRecords(ctx, host, resolver.LookupTXT, assertion.TXT, &errs)
		checkRecords(ctx, host, resolver.LookupNS, assertion.NS, &errs)
		checkRecords(ctx, host, resolver.LookupHost, assertion.A, &errs)
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

func (c DNSCheck) Subject() string {
	return fmt.Sprintf("%s", c.Host)
}

func checkRecords[T any](
	ctx context.Context,
	host string,
	lookup func(ctx context.Context, host string) ([]T, error),
	assert *assertion.StringListMatcher,
	errs *[]error,
) {
	if assert == nil {
		return
	}

	records, err := lookup(ctx, host)
	if err != nil {
		*errs = append(*errs, err)
		return
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
			*errs = append(*errs, fmt.Errorf("Unsupported type for record matching: %T", record))
			return
		}
	}

	if err = assert.Match(inputs); err != nil {
		*errs = append(*errs, err)
	}
}
