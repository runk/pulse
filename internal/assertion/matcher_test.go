package assertion

import (
	"testing"
)

func TestNumberMatcher(t *testing.T) {
	scenarios := []struct {
		name     string
		matcher  NumberMatcher
		input    any
		expected string
	}{
		{"equals pass", NumberMatcher{Equals: floatPtr(1)}, 1, ""},
		{"equals zero pass", NumberMatcher{Equals: floatPtr(0)}, 0, ""},
		{"equals fail", NumberMatcher{Equals: floatPtr(1)}, 2, "expected 1 but got 2"},
		{"not equals pass", NumberMatcher{NotEquals: floatPtr(1)}, 2, ""},
		{"not equals fail", NumberMatcher{NotEquals: floatPtr(1)}, 1, "expected 1 not to be 1"},
		{"gt pass", NumberMatcher{Gt: floatPtr(1)}, 2, ""},
		{"gt fail when equal", NumberMatcher{Gt: floatPtr(1)}, 1, "expected 1 to be greater than 1"},
		{"gt fail when less", NumberMatcher{Gt: floatPtr(1)}, 0, "expected 0 to be greater than 1"},
		{"gte pass when greater", NumberMatcher{Gte: floatPtr(1)}, 2, ""},
		{"gte pass when equal", NumberMatcher{Gte: floatPtr(1)}, 1, ""},
		{"gte fail", NumberMatcher{Gte: floatPtr(1)}, 0, "expected 0 to be greater than or equal to 1"},
		{"lt pass", NumberMatcher{Lt: floatPtr(2)}, 1, ""},
		{"lt fail when equal", NumberMatcher{Lt: floatPtr(2)}, 2, "expected 2 to be less than 2"},
		{"lt fail when greater", NumberMatcher{Lt: floatPtr(2)}, 3, "expected 3 to be less than 2"},
		{"lte pass when less", NumberMatcher{Lte: floatPtr(2)}, 1, ""},
		{"lte pass when equal", NumberMatcher{Lte: floatPtr(2)}, 2, ""},
		{"lte fail", NumberMatcher{Lte: floatPtr(2)}, 3, "expected 3 to be less than or equal to 2"},
		{"in pass", NumberMatcher{In: []float64{1, 2, 3}}, 2, ""},
		{"in fail", NumberMatcher{In: []float64{1, 2, 3}}, 4, "expected 4 to be one of [1 2 3]"},
		{"not in pass", NumberMatcher{NotIn: []float64{1, 2, 3}}, 4, ""},
		{"not in fail", NumberMatcher{NotIn: []float64{1, 2, 3}}, 2, "expected 2 not to be one of [1 2 3]"},
		{"between pass in range", NumberMatcher{Between: []float64{1, 3}}, 2, ""},
		{"between pass lower bound", NumberMatcher{Between: []float64{1, 3}}, 1, ""},
		{"between pass upper bound", NumberMatcher{Between: []float64{1, 3}}, 3, ""},
		{"between fail below", NumberMatcher{Between: []float64{1, 3}}, 0, "expected 0 to be between 1 and 3"},
		{"between fail above", NumberMatcher{Between: []float64{1, 3}}, 4, "expected 4 to be between 1 and 3"},
		{"between fail invalid length", NumberMatcher{Between: []float64{1}}, 1, "between requires exactly 2 values, got 1"},
		{"between fail inverted bounds", NumberMatcher{Between: []float64{3, 1}}, 2, "between lower bound 3 is greater than upper bound 1"},
		{"input is not a number", NumberMatcher{Equals: floatPtr(1)}, "1", "input is not a number"},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			actual := ""
			if err := scenario.matcher.Match(scenario.input); err != nil {
				actual = err.Error()
			}

			if actual != scenario.expected {
				t.Errorf("Want '%s' but got '%s'", scenario.expected, actual)
			}
		})
	}
}

func floatPtr(v float64) *float64 {
	return &v
}
