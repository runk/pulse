package check

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

type CheckValue interface {
	Type() string
	Run(ctx context.Context) error
	Validate() error
	Subject() string
}

type Check struct {
	Value CheckValue
}

func (c *Check) UnmarshalJSON(data []byte) error {

	var header struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal(data, &header); err != nil {
		return errors.Join(err, errors.New("Cannot read header"))
	}

	var value CheckValue
	switch header.Type {
	case "http":
		value = &HTTPCheck{}
	case "dns":
		value = &DNSCheck{}
	case "tls":
		value = &TLSCheck{}
	default:
		return fmt.Errorf("Unsupported check type: '%s'", header.Type)
	}

	if err := json.Unmarshal(data, value); err != nil {
		return errors.Join(err, fmt.Errorf("Cannot read value of type '%s'", header.Type))
	}

	c.Value = value

	return nil
}
