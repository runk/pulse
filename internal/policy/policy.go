package policy

import (
	"encoding/json"
	"os"

	"github.com/runk/pulse/internal/check"
)

type Policy struct {
	Name   string        `json:"name"`
	Checks []check.Check `json:"checks"`
}

func (p Policy) Validate() error {
	for _, check := range p.Checks {
		err := check.Value.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func ReadPolicy(filename string) (Policy, error) {
	content, err := os.ReadFile(filename)

	if err != nil {
		return Policy{}, err
	}

	var policy Policy

	if err = json.Unmarshal(content, &policy); err != nil {
		return Policy{}, err
	}

	return policy, nil
}
