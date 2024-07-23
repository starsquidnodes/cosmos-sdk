package coretesting

import "cosmossdk.io/core/log"

// NewNopLogger returns a new logger that does nothing.
func NewNopLogger() log.Logger {
	// The custom nopLogger is about 3x faster than a zerolog.Nop().
	return nopLogger{}
}

// nopLogger is a Logger that does nothing when called.
// See the "specialized nop logger" benchmark and compare with the "zerolog nop logger" benchmark.
// The custom implementation is about 3x faster.
type nopLogger struct{}

func (nopLogger) Info(string, ...any)    {}
func (nopLogger) Warn(string, ...any)    {}
func (nopLogger) Error(string, ...any)   {}
func (nopLogger) Debug(string, ...any)   {}
func (nopLogger) With(...any) log.Logger { return nopLogger{} }
func (nopLogger) Impl() any              { return nopLogger{} }
