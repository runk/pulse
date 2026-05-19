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
		{"equals pass", NumberMatcher{Equals: new(1.0)}, 1, ""},
		{"equals zero pass", NumberMatcher{Equals: new(0.0)}, 0, ""},
		{"equals negative pass", NumberMatcher{Equals: new(-1.0)}, -1, ""},
		{"equals float pass", NumberMatcher{Equals: new(1.5)}, 1.5, ""},
		{"equals float32 pass", NumberMatcher{Equals: new(1.5)}, float32(1.5), ""},
		{"equals fail", NumberMatcher{Equals: new(1.0)}, 2, "expected 1 but got 2"},
		{"equals negative fail", NumberMatcher{Equals: new(-1.0)}, 1, "expected -1 but got 1"},
		{"equals zero fail", NumberMatcher{Equals: new(0.0)}, 1, "expected 0 but got 1"},
		{"not equals pass", NumberMatcher{NotEquals: new(1.0)}, 2, ""},
		{"not equals negative pass", NumberMatcher{NotEquals: new(-1.0)}, 1, ""},
		{"not equals zero pass", NumberMatcher{NotEquals: new(0.0)}, 1, ""},
		{"not equals fail", NumberMatcher{NotEquals: new(1.0)}, 1, "expected 1 not to be 1"},
		{"not equals negative fail", NumberMatcher{NotEquals: new(-1.0)}, -1, "expected -1 not to be -1"},
		{"not equals zero fail", NumberMatcher{NotEquals: new(0.0)}, 0, "expected 0 not to be 0"},
		{"gt pass", NumberMatcher{Gt: new(1.0)}, 2, ""},
		{"gt negative pass", NumberMatcher{Gt: new(-2.0)}, -1, ""},
		{"gt zero pass", NumberMatcher{Gt: new(0.0)}, 0.1, ""},
		{"gt float pass", NumberMatcher{Gt: new(1.5)}, 1.6, ""},
		{"gt fail when equal", NumberMatcher{Gt: new(1.0)}, 1, "expected 1 to be greater than 1"},
		{"gt fail when less", NumberMatcher{Gt: new(1.0)}, 0, "expected 0 to be greater than 1"},
		{"gt zero fail when equal", NumberMatcher{Gt: new(0.0)}, 0, "expected 0 to be greater than 0"},
		{"gt negative fail", NumberMatcher{Gt: new(-1.0)}, -2, "expected -2 to be greater than -1"},
		{"gte pass when greater", NumberMatcher{Gte: new(1.0)}, 2, ""},
		{"gte pass when equal", NumberMatcher{Gte: new(1.0)}, 1, ""},
		{"gte negative pass when equal", NumberMatcher{Gte: new(-1.0)}, -1, ""},
		{"gte zero pass when equal", NumberMatcher{Gte: new(0.0)}, 0, ""},
		{"gte float pass", NumberMatcher{Gte: new(1.5)}, 1.5, ""},
		{"gte fail", NumberMatcher{Gte: new(1.0)}, 0, "expected 0 to be greater than or equal to 1"},
		{"gte negative fail", NumberMatcher{Gte: new(-1.0)}, -2, "expected -2 to be greater than or equal to -1"},
		{"lt pass", NumberMatcher{Lt: new(2.0)}, 1, ""},
		{"lt negative pass", NumberMatcher{Lt: new(-1.0)}, -2, ""},
		{"lt zero pass", NumberMatcher{Lt: new(0.0)}, -0.1, ""},
		{"lt float pass", NumberMatcher{Lt: new(1.5)}, 1.4, ""},
		{"lt fail when equal", NumberMatcher{Lt: new(2.0)}, 2, "expected 2 to be less than 2"},
		{"lt fail when greater", NumberMatcher{Lt: new(2.0)}, 3, "expected 3 to be less than 2"},
		{"lt zero fail when equal", NumberMatcher{Lt: new(0.0)}, 0, "expected 0 to be less than 0"},
		{"lt negative fail", NumberMatcher{Lt: new(-2.0)}, -1, "expected -1 to be less than -2"},
		{"lte pass when less", NumberMatcher{Lte: new(2.0)}, 1, ""},
		{"lte pass when equal", NumberMatcher{Lte: new(2.0)}, 2, ""},
		{"lte negative pass when equal", NumberMatcher{Lte: new(-1.0)}, -1, ""},
		{"lte zero pass when equal", NumberMatcher{Lte: new(0.0)}, 0, ""},
		{"lte float pass", NumberMatcher{Lte: new(1.5)}, 1.5, ""},
		{"lte fail", NumberMatcher{Lte: new(2.0)}, 3, "expected 3 to be less than or equal to 2"},
		{"lte negative fail", NumberMatcher{Lte: new(-2.0)}, -1, "expected -1 to be less than or equal to -2"},
		{"in pass", NumberMatcher{In: []float64{1, 2, 3}}, 2, ""},
		{"in negative pass", NumberMatcher{In: []float64{-3, -2, -1}}, -2, ""},
		{"in zero pass", NumberMatcher{In: []float64{-1, 0, 1}}, 0, ""},
		{"in float pass", NumberMatcher{In: []float64{1.1, 1.2, 1.3}}, 1.2, ""},
		{"in fail", NumberMatcher{In: []float64{1, 2, 3}}, 4, "expected 4 to be one of [1 2 3]"},
		{"in negative fail", NumberMatcher{In: []float64{-3, -2, -1}}, 1, "expected 1 to be one of [-3 -2 -1]"},
		{"not in pass", NumberMatcher{NotIn: []float64{1, 2, 3}}, 4, ""},
		{"not in negative pass", NumberMatcher{NotIn: []float64{-3, -2, -1}}, 1, ""},
		{"not in zero pass", NumberMatcher{NotIn: []float64{-1, 1}}, 0, ""},
		{"not in float pass", NumberMatcher{NotIn: []float64{1.1, 1.2, 1.3}}, 1.4, ""},
		{"not in fail", NumberMatcher{NotIn: []float64{1, 2, 3}}, 2, "expected 2 not to be one of [1 2 3]"},
		{"not in negative fail", NumberMatcher{NotIn: []float64{-3, -2, -1}}, -2, "expected -2 not to be one of [-3 -2 -1]"},
		{"not in zero fail", NumberMatcher{NotIn: []float64{-1, 0, 1}}, 0, "expected 0 not to be one of [-1 0 1]"},
		{"between pass in range", NumberMatcher{Between: []float64{1, 3}}, 2, ""},
		{"between pass lower bound", NumberMatcher{Between: []float64{1, 3}}, 1, ""},
		{"between pass upper bound", NumberMatcher{Between: []float64{1, 3}}, 3, ""},
		{"between negative pass in range", NumberMatcher{Between: []float64{-3, -1}}, -2, ""},
		{"between zero pass in range", NumberMatcher{Between: []float64{-1, 1}}, 0, ""},
		{"between float pass in range", NumberMatcher{Between: []float64{1.1, 1.3}}, 1.2, ""},
		{"between fail below", NumberMatcher{Between: []float64{1, 3}}, 0, "expected 0 to be between 1 and 3"},
		{"between fail above", NumberMatcher{Between: []float64{1, 3}}, 4, "expected 4 to be between 1 and 3"},
		{"between negative fail below", NumberMatcher{Between: []float64{-3, -1}}, -4, "expected -4 to be between -3 and -1"},
		{"between negative fail above", NumberMatcher{Between: []float64{-3, -1}}, 0, "expected 0 to be between -3 and -1"},
		{"between fail invalid length", NumberMatcher{Between: []float64{1}}, 1, "between requires exactly 2 values, got 1"},
		{"between fail inverted bounds", NumberMatcher{Between: []float64{3, 1}}, 2, "between lower bound 3 is greater than upper bound 1"},
		{"input string is not a number", NumberMatcher{Equals: new(1.0)}, "1", "input is not a number"},
		{"input bool is not a number", NumberMatcher{Equals: new(1.0)}, true, "input is not a number"},
		{"input nil is not a number", NumberMatcher{Equals: new(1.0)}, nil, "input is not a number"},
		{"input slice is not a number", NumberMatcher{Equals: new(1.0)}, []float64{1}, "input is not a number"},
		{"input map is not a number", NumberMatcher{Equals: new(1.0)}, map[string]float64{"value": 1}, "input is not a number"},
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

func TestStringMatcher(t *testing.T) {
	scenarios := []struct {
		name     string
		matcher  StringMatcher
		input    any
		expected string
	}{
		{"equals pass", StringMatcher{Equals: new("abc")}, "abc", ""},
		{"equals empty pass", StringMatcher{Equals: new("")}, "", ""},
		{"equals fail", StringMatcher{Equals: new("abc")}, "xyz", "expected 'xyz' to be 'abc'"},
		{"equals case-sensitive fail", StringMatcher{Equals: new("abc")}, "ABC", "expected 'ABC' to be 'abc'"},
		{"equals empty fail", StringMatcher{Equals: new("")}, "abc", "expected 'abc' to be ''"},
		{"not equals pass", StringMatcher{NotEquals: new("abc")}, "xyz", ""},
		{"not equals case-sensitive pass", StringMatcher{NotEquals: new("abc")}, "ABC", ""},
		{"not equals empty pass", StringMatcher{NotEquals: new("")}, "abc", ""},
		{"not equals fail", StringMatcher{NotEquals: new("abc")}, "abc", "expected 'abc' not to be 'abc'"},
		{"not equals empty fail", StringMatcher{NotEquals: new("")}, "", "expected '' not to be ''"},
		{"contains pass", StringMatcher{Contains: new("a")}, "abc", ""},
		{"contains empty pass", StringMatcher{Contains: new("")}, "abc", ""},
		{"contains fail", StringMatcher{Contains: new("a")}, "xyz", "expected to contain 'a'"},
		{"contains case-sensitive fail", StringMatcher{Contains: new("a")}, "ABC", "expected to contain 'a'"},
		{"matches pass", StringMatcher{Matches: new(`abc-\d+`)}, "abc-123", ""},
		{"matches partial pass", StringMatcher{Matches: new(`\d+`)}, "abc-123", ""},
		{"matches empty pass", StringMatcher{Matches: new("")}, "abc", ""},
		{"matches fail", StringMatcher{Matches: new(`abc-\d+`)}, "abc", "expected to match 'abc-\\d+'"},
		{"matches invalid regex", StringMatcher{Matches: new("[")}, "abc", "regular expression is invalid: error parsing regexp: missing closing ]: `[`"},
		{"input int is not a string", StringMatcher{Equals: new("abc")}, 1, "input is not a string, got: int"},
		{"input bool is not a string", StringMatcher{Equals: new("abc")}, true, "input is not a string, got: bool"},
		{"input nil is not a string", StringMatcher{Equals: new("abc")}, nil, "input is not a string, got: <nil>"},
		{"input slice is not a string", StringMatcher{Equals: new("abc")}, []string{"abc"}, "input is not a string, got: []string"},
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
