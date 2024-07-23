package log_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/rs/zerolog"

	"cosmossdk.io/log"
)

func TestLoggerOptionStackTrace(t *testing.T) {
	// t.Run("With stack trace", func(t *testing.T) {
	// 	buf := new(bytes.Buffer)
	// 	logger := log.NewLogger(buf, log.TraceOption(true), log.ColorOption(false))
	// 	logger.Error("this log should be displayed", "error", inner())
	// 	if strings.Count(buf.String(), "logger_test.go") != 1 {
	// 		t.Errorf("expected stack trace, but got: %s", buf.String())
	// 	}
	// })

	t.Run("Without stack trace", func(t *testing.T) {
		buf := new(bytes.Buffer)
		logger := log.NewLogger(buf, log.TraceOption(false), log.ColorOption(false))
		logger.Error("this log should be displayed", "error", inner())
		if strings.Count(buf.String(), "logger_test.go") > 0 {
			t.Errorf("unexpected stack trace found: %s", buf.String())
		}
	})
}

func inner() error {
	return errors.New("seems we have an error here")
}

type _MockHook string

func (h _MockHook) Run(e *zerolog.Event, l zerolog.Level, msg string) {
	e.Bool(string(h), true)
}

func TestLoggerOptionHooks(t *testing.T) {
	t.Run("With hooks", func(t *testing.T) {
		buf := new(bytes.Buffer)
		var (
			mockHook1 _MockHook = "mock_message1"
			mockHook2 _MockHook = "mock_message2"
		)
		logger := log.NewLogger(buf, log.HooksOption(mockHook1, mockHook2), log.ColorOption(false))
		logger.Info("hello world")

		output := buf.String()
		if !strings.Contains(output, "mock_message1=true") {
			t.Error("expected output to contain 'mock_message1=true'")
		}
		if !strings.Contains(output, "mock_message2=true") {
			t.Error("expected output to contain 'mock_message2=true'")
		}
	})

	t.Run("Without hooks", func(t *testing.T) {
		buf := new(bytes.Buffer)
		logger := log.NewLogger(buf, log.HooksOption(), log.ColorOption(false))
		logger.Info("hello world")

		output := buf.String()
		if !strings.Contains(output, "hello world") {
			t.Error("expected output to contain 'hello world'")
		}
	})
}
