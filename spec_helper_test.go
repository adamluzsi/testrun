package testcase_test

import (
	"github.com/adamluzsi/testcase/internal"
	"github.com/stretchr/testify/require"
	"testing"
)

type CustomTB struct {
	testing.TB
	isFatalFCalled bool
}

func (tb *CustomTB) Run(name string, blk func(tb testing.TB)) bool {
	switch tb := tb.TB.(type) {
	case *testing.T:
		return tb.Run(name, func(t *testing.T) { blk(t) })
	case *testing.B:
		return tb.Run(name, func(b *testing.B) { blk(b) })
	default:
		panic("implement me")
	}
}

func (t *CustomTB) Fatalf(format string, args ...interface{}) {
	t.isFatalFCalled = true
	return
}

func unsupported(tb testing.TB) {
	tb.Skip(`unsupported`)
}

func isFatalFn(stub *internal.StubTB) func(block func()) bool {
	return func(block func()) bool {
		stub.IsFailed = false
		defer func() { stub.IsFailed = false }()
		var finished bool
		internal.InGoroutine(func() {
			block()
			finished = true
		})
		return !finished && stub.Failed()
	}
}

func willFatalWithMessageFn(stub *internal.StubTB) func(tb testing.TB, blk func()) []string {
	isFatal := isFatalFn(stub)
	return func(tb testing.TB, blk func()) []string {
		stub.Logs = nil
		require.True(tb, isFatal(blk))
		return stub.Logs
	}
}
