package policy

import (
	"encoding/json"
	"os"
)

type Check struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

type Policy struct {
	Name   string  `json:"name"`
	Checks []Check `json:"checks"`
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
