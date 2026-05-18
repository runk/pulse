package assertion

import (
	"errors"
	"fmt"
)

type Matcher interface {
	Match(input any) error
}

type NumberMatcher struct {
	Equals    float64   `json:"equals,omitempty"`
	NotEquals float64   `json:"notEquals,omitempty"`
	Gt        float64   `json:"gt,omitempty"`
	Gte       float64   `json:"gte,omitempty"`
	Lt        float64   `json:"lt,omitempty"`
	Lte       float64   `json:"lte,omitempty"`
	In        []float64 `json:"in,omitempty"`
	NotIn     []float64 `json:"notIn,omitempty"`
	Between   []float64 `json:"between,omitempty"`
}

func asFloat64(input any) (float64, bool) {
	switch v := input.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	default:
		return 0, false
	}
}

func (m NumberMatcher) Match(input any) error {
	value, ok := asFloat64(input)
	if !ok {
		return errors.New("input is not a number")
	}

	if m.Equals != 0 {
		if value != m.Equals {
			return fmt.Errorf("expected %v but got %v", m.Equals, value)
		}
	}

	if m.NotEquals != 0 {
		if value == m.NotEquals {
			return fmt.Errorf("expected %v not to be %v", value, m.NotEquals)
		}
	}

	return nil
}

type StringMatcher struct {
	Equals    string `json:"equals,omitempty"`
	NotEquals string `json:"notEquals,omitempty"`
	Contains  string `json:"contains,omitempty"`
	Matches   string `json:"matches,omitempty"`
}

func (m StringMatcher) Match(input any) error {
	return nil
}
