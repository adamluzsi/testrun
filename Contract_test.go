package testcase_test

import (
	"testing"

	"github.com/adamluzsi/testcase"
	"github.com/stretchr/testify/require"
)

func TestRunContracts(t *testing.T) {
	t.Run(`when TB is *testing.T`, func(t *testing.T) {
		sT := &RunContractExampleContract{}
		testcase.RunContract(&testing.T{}, sT)
		require.True(t, sT.TestWasCalled)
		require.False(t, sT.BenchmarkWasCalled)
	})

	t.Run(`when TB is *testing.B`, func(t *testing.T) {
		sB := &RunContractExampleContract{}
		testcase.RunContract(&testing.B{}, sB)
		require.False(t, sB.TestWasCalled)
		require.True(t, sB.BenchmarkWasCalled)
	})

	t.Run(`when TB is *testcase.T with *testing.T under the hood`, func(t *testing.T) {
		sT := &RunContractExampleContract{}
		testcase.RunContract(&testcase.T{TB: &testing.T{}}, sT)
		require.True(t, sT.TestWasCalled)
		require.False(t, sT.BenchmarkWasCalled)
	})

	t.Run(`when TB is *testcase.T with *testing.B under the hood`, func(t *testing.T) {
		sT := &RunContractExampleContract{}
		testcase.RunContract(&testcase.T{TB: &testing.B{}}, sT)
		require.False(t, sT.TestWasCalled)
		require.True(t, sT.BenchmarkWasCalled)
	})

	t.Run(`when TB is *testcase.Spec for *testing.T`, func(t *testing.T) {
		s := testcase.NewSpec(t)
		a := &RunContractExampleContract{}
		b := &RunContractExampleContract{}
		testcase.RunContract(s, a, b)
		s.Finish()
		require.True(t, a.TestWasCalled)
		require.False(t, a.BenchmarkWasCalled)
		require.True(t, b.TestWasCalled)
		require.False(t, b.BenchmarkWasCalled)
	})

	t.Run(`when TB is an unknown test runner type`, func(t *testing.T) {
		type UnknownTestRunner struct{}
		require.Panics(t, func() { testcase.RunContract(UnknownTestingTB{}, &RunContractExampleContract{}) })
	})
}

type customTestTB struct {
	testing.TB
	isFatalFCalled bool
}

func (tb *customTestTB) Run(name string, blk func(tb testing.TB)) bool {
	switch tb := tb.TB.(type) {
	case *testing.T:
		return tb.Run(name, func(t *testing.T) { blk(t) })
	case *testing.B:
		return tb.Run(name, func(b *testing.B) { blk(b) })
	default:
		panic("implement me")
	}
}

func (t *customTestTB) Fatalf(format string, args ...interface{}) {
	t.isFatalFCalled = true
	return
}

type RunContractExampleContract struct {
	TestWasCalled      bool
	BenchmarkWasCalled bool
}

func (spec *RunContractExampleContract) Test(t *testing.T) {
	spec.TestWasCalled = true
}

func (spec *RunContractExampleContract) Benchmark(b *testing.B) {
	spec.BenchmarkWasCalled = true
}
