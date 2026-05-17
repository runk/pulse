package runner

import (
	"fmt"

	"github.com/runk/pulse/internal/policy"
)

func Execute(checks *[]policy.Check, concurrency uint16) error {
	fmt.Println(checks)
	return nil
}
