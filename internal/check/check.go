package check

import (
	"encoding/json"
	"fmt"
)

type CheckValue interface {
	Type() string
	Run() error
}

type Check struct {
	Value CheckValue
}

func (c *Check) UnmarshalJSON(data []byte) error {

	var header struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal(data, &header); err != nil {
		return err
	}

	var value CheckValue
	switch header.Type {
	case "http":
		value = HTTPCheck{}
	case "dns":
		value = DNSCheck{}
	default:
		return fmt.Errorf("Unsupported check type: '%s'", header.Type)
	}

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	c.Value = value

	return nil
}
