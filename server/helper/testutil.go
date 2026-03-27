package helper

import (
	"os"
	"testing"
)

// VerifyIfSkipIntegrationTests skips the test unless RUN_INTEGRATION_TESTS=yes.
func VerifyIfSkipIntegrationTests(t *testing.T) {
	t.Helper()
	if os.Getenv("RUN_INTEGRATION_TESTS") != "yes" {
		t.Skipf("skipping integration tests")
	}
}

// IsRunningInDocker reports whether the process is running inside a Docker container.
func IsRunningInDocker() bool {
	if _, err := os.Stat("/.dockerenv"); os.IsNotExist(err) {
		return false
	}
	return true
}

