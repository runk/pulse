package policy

import (
	"regexp"
	"testing"
)

func TestReadPolicyOk(t *testing.T) {
	policy, err := ReadPolicy("../../example/policy-basic.json")

	if err != nil {
		t.Fatalf("Cannot read policy: %v", err)
	}

	if policy.Name == "" {
		t.Errorf("Policy should not be nil")
	}

	if len(policy.Checks) == 0 {
		t.Fatal("expected at least one check")
	}
	for i, c := range policy.Checks {
		if c.Value.Type() == "" {
			t.Fatalf("check %d: type is empty", i)
		}
	}
}

func TestReadPolicyNoFile(t *testing.T) {
	policy, err := ReadPolicy("../../example/lemon.json")

	if err == nil {
		t.Fatal("Error must be returned")
	}

	re := regexp.MustCompile("no such file or directory")
	if !re.Match([]byte(err.Error())) {
		t.Fatal("It should report about missing file")
	}

	if policy.Name != "" {
		t.Fatal("Empty policy must be returned")
	}

}

func TestReadPolicyCorrupted(t *testing.T) {
	policy, err := ReadPolicy("../../example/policy-corrupted.json")

	if err == nil {
		t.Fatal("Error must be returned")
	}

	re := regexp.MustCompile("unexpected end of JSON input")
	if !re.Match([]byte(err.Error())) {
		t.Fatal("It should report syntax errors")
	}

	if policy.Name != "" {
		t.Fatal("Empty policy must be returned")
	}
}
