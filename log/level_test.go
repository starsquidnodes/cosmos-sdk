package log_test

import (
	"strings"
	"testing"

	"cosmossdk.io/log"
)

func TestParseLogLevel(t *testing.T) {
	t.Run("Empty log level", func(t *testing.T) {
		_, err := log.ParseLogLevel("")
		if err == nil {
			t.Error("Expected error for empty log level, but got nil")
		}
		if !strings.Contains(err.Error(), "empty log level") {
			t.Errorf("Expected error message to contain 'empty log level', but got: %v", err)
		}
	})

	t.Run("Invalid log level", func(t *testing.T) {
		level := "consensus:foo,mempool:debug,*:error"
		_, err := log.ParseLogLevel(level)
		if err == nil {
			t.Error("Expected error for invalid log level, but got nil")
		}
		expectedErrMsg := "invalid log level foo in log level list [consensus:foo mempool:debug *:error]"
		if !strings.Contains(err.Error(), expectedErrMsg) {
			t.Errorf("Expected error message to contain '%s', but got: %v", expectedErrMsg, err)
		}
	})

	t.Run("Valid log level", func(t *testing.T) {
		level := "consensus:debug,mempool:debug,*:error"
		filter, err := log.ParseLogLevel(level)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if filter == nil {
			t.Fatal("Expected non-nil filter, but got nil")
		}

		testCases := []struct {
			module string
			level  string
			expect bool
		}{
			{"consensus", "debug", false},
			{"consensus", "info", false},
			{"consensus", "error", false},
			{"mempool", "debug", false},
			{"mempool", "info", false},
			{"mempool", "error", false},
			{"state", "error", false},
			{"server", "panic", false},
			{"server", "debug", true},
			{"state", "debug", true},
			{"state", "info", true},
		}

		for _, tc := range testCases {
			t.Run(tc.module+":"+tc.level, func(t *testing.T) {
				result := filter(tc.module, tc.level)
				if result != tc.expect {
					t.Errorf("Expected filter(%s, %s) to be %v, but got %v", tc.module, tc.level, tc.expect, result)
				}
			})
		}
	})

	t.Run("Simple error level", func(t *testing.T) {
		level := "error"
		filter, err := log.ParseLogLevel(level)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if filter == nil {
			t.Fatal("Expected non-nil filter, but got nil")
		}

		testCases := []struct {
			module string
			level  string
			expect bool
		}{
			{"state", "error", false},
			{"consensus", "error", false},
			{"consensus", "debug", true},
			{"consensus", "info", true},
			{"state", "debug", true},
		}

		for _, tc := range testCases {
			t.Run(tc.module+":"+tc.level, func(t *testing.T) {
				result := filter(tc.module, tc.level)
				if result != tc.expect {
					t.Errorf("Expected filter(%s, %s) to be %v, but got %v", tc.module, tc.level, tc.expect, result)
				}
			})
		}
	})
}
