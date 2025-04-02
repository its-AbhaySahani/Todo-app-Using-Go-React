package e2e

import (
    "os"
    "testing"
)

func TestMain(m *testing.M) {
    // Set test container ryuk to false for Windows compatibility
    os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")

    // Run tests
    code := m.Run()

    // Exit with test status code
    os.Exit(code)
}