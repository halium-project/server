package util

import (
	"os"
	"strings"
	"testing"
)

// Type of test.
type Type int

// All the availables type of test.
const (
	// Unit tests don't have dependencies, only mocks.
	Unit Type = iota
	// Integration tests check if the clients work well whithin the environment.
	Integration
	// Run the test like an external client and see if the service provided
	// is spec compliant.
	EndToEnd
)

// Is indicates which type of test is run. It also allows to check with the
// "TEST_ENV" env variable if this test should be ran or not and skip it in
// if the test is not available with the current environment.
func TestIs(t *testing.T, testType Type) {
	t.Helper()

	var env Type

	envList := strings.Split(os.Getenv("TEST_ENV"), ",")

	for _, envStr := range envList {
		switch strings.ToLower(envStr) {
		case "unit":
			env = Unit
		case "integration":
			env = Integration
		case "end_to_end":
			env = EndToEnd
		case "":
			t.Log(`env variable "TEST_ENV" empty, fallback on "UNIT"`)
			env = Unit
		default:
			t.Fatal(`invalid value for env variable "TEST_ENV"`)
		}

		if testType == env {
			return
		}
	}

	t.SkipNow()
}
