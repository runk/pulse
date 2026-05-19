package assertion

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

type Matcher interface {
	Match(input any) error
}

type NumberMatcher struct {
	Equals    *float64  `json:"equals,omitempty"`
	NotEquals *float64  `json:"notEquals,omitempty"`
	Gt        *float64  `json:"gt,omitempty"`
	Gte       *float64  `json:"gte,omitempty"`
	Lt        *float64  `json:"lt,omitempty"`
	Lte       *float64  `json:"lte,omitempty"`
	In        []float64 `json:"in,omitempty"`
	NotIn     []float64 `json:"notIn,omitempty"`
	Between   []float64 `json:"between,omitempty"`
}

func (m NumberMatcher) Match(input any) error {
	value, ok := asFloat64(input)
	if !ok {
		return errors.New("input is not a number")
	}

	if m.Equals != nil && *m.Equals != value {
		return fmt.Errorf("expected %v but got %v", *m.Equals, value)
	}

	if m.NotEquals != nil && *m.NotEquals == value {
		return fmt.Errorf("expected %v not to be %v", value, *m.NotEquals)
	}

	if m.Gt != nil && value <= *m.Gt {
		return fmt.Errorf("expected %v to be greater than %v", value, *m.Gt)
	}

	if m.Gte != nil && value < *m.Gte {
		return fmt.Errorf("expected %v to be greater than or equal to %v", value, *m.Gte)
	}

	if m.Lt != nil && value >= *m.Lt {
		return fmt.Errorf("expected %v to be less than %v", value, *m.Lt)
	}

	if m.Lte != nil && value > *m.Lte {
		return fmt.Errorf("expected %v to be less than or equal to %v", value, *m.Lte)
	}

	if len(m.In) > 0 && !slices.Contains(m.In, value) {
		return fmt.Errorf("expected %v to be one of %v", value, m.In)
	}

	if len(m.NotIn) > 0 && slices.Contains(m.NotIn, value) {
		return fmt.Errorf("expected %v not to be one of %v", value, m.NotIn)
	}

	if len(m.Between) > 0 {
		if len(m.Between) != 2 {
			return fmt.Errorf("between requires exactly 2 values, got %v", len(m.Between))
		}

		if m.Between[0] > m.Between[1] {
			return fmt.Errorf("between lower bound %v is greater than upper bound %v", m.Between[0], m.Between[1])
		}

		if value < m.Between[0] || value > m.Between[1] {
			return fmt.Errorf("expected %v to be between %v and %v", value, m.Between[0], m.Between[1])
		}
	}

	return nil
}

type StringMatcher struct {
	Equals    *string `json:"equals,omitempty"`
	NotEquals *string `json:"notEquals,omitempty"`
	Contains  *string `json:"contains,omitempty"`
	Matches   *string `json:"matches,omitempty"`
}

func (m StringMatcher) Match(input any) error {
	value, ok := asString(input)

	if !ok {
		return fmt.Errorf("input is not a string, got: %T", input)
	}

	if m.Equals != nil && *m.Equals != value {
		return fmt.Errorf("expected '%s' to be '%s'", value, *m.Equals)
	}

	if m.NotEquals != nil && *m.NotEquals == value {
		return fmt.Errorf("expected '%s' not to be '%s'", value, *m.NotEquals)
	}

	if m.Contains != nil && !strings.Contains(value, *m.Contains) {
		return fmt.Errorf("expected to contain '%s'", *m.Contains)
	}

	if m.Matches != nil {
		re, err := regexp.Compile(*m.Matches)
		if err != nil {
			return fmt.Errorf("regular expression is invalid: %v", err)
		}

		if !re.MatchString(value) {
			return fmt.Errorf("expected to match '%s'", *m.Matches)
		}
	}

	return nil
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

func asString(input any) (string, bool) {
	switch v := input.(type) {
	case string:
		return v, true
	case []byte:
		return string(v), true
	case *[]byte:
		return string(*v), true
	default:
		return "", false
	}
}
