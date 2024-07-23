package log_test

import (
	"bytes"
	"strings"
	"testing"

	"cosmossdk.io/log"
)

func TestFilteredWriter(t *testing.T) {
	level := "consensus:debug,mempool:debug,*:error"
	filter, err := log.ParseLogLevel(level)
	if err != nil {
		t.Fatalf("Unexpected error parsing log level: %v", err)
	}

	buf := new(bytes.Buffer)
	logger := log.NewLogger(buf, log.FilterOption(filter))

	t.Run("Displayed log", func(t *testing.T) {
		logger.Debug("this log line should be displayed", log.ModuleKey, "consensus")
		if !strings.Contains(buf.String(), "this log line should be displayed") {
			t.Errorf("Expected log line was not found in the output, got: %s", buf.String())
		}
		buf.Reset()
	})

	t.Run("Filtered log", func(t *testing.T) {
		logger.Debug("this log line should be filtered", log.ModuleKey, "server")
		if buf.Len() != 0 {
			t.Errorf("Expected empty buffer, but got: %s", buf.String())
		}
	})
}
