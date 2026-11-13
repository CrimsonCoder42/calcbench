package psychiatry

import (
	"errors"
	"testing"
)

func mdqYes(n int) [13]bool {
	var a [13]bool
	for i := 0; i < n; i++ {
		a[i] = true
	}
	return a
}

// Nagata 2013 / Cybulski 2020: ≥ 7 yes, clustered, moderate or serious problem; each condition on both sides.
func TestMDQScreen_example(t *testing.T) {
	cases := []struct {
		n       int
		same    bool
		problem int
		want    MDQResult
	}{
		{7, true, MDQModerateProblem, MDQResult{7, true}},
		{6, true, MDQSeriousProblem, MDQResult{6, false}},
		{13, true, MDQMinorProblem, MDQResult{13, false}},
		{13, false, MDQSeriousProblem, MDQResult{13, false}},
		{0, false, MDQNoProblem, MDQResult{0, false}},
	}
	for _, c := range cases {
		if got, err := MDQScreen(mdqYes(c.n), c.same, c.problem); err != nil || got != c.want {
			t.Errorf("MDQScreen(%d, %v, %d) = (%+v, %v), want %+v", c.n, c.same, c.problem, got, err, c.want)
		}
	}
	for _, bad := range []int{-1, 4} {
		if _, err := MDQScreen(mdqYes(7), true, bad); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("problem %d: %v, want ErrOutOfRange", bad, err)
		}
	}
}

// parity: ledger mdq_bipolar, source cases 1–6 (source mdq_bipolar_test.go: none, six, seven positive, seven no
// co-occurrence, ten minor impairment, all).
func TestMDQScreen_parity(t *testing.T) {
	cases := []struct {
		n       int
		same    bool
		problem int
		want    MDQResult
	}{
		{0, false, 0, MDQResult{0, false}},  // parity: ledger mdq_bipolar, source case 1
		{6, true, 3, MDQResult{6, false}},   // parity: ledger mdq_bipolar, source case 2
		{7, true, 2, MDQResult{7, true}},    // parity: ledger mdq_bipolar, source case 3
		{7, false, 3, MDQResult{7, false}},  // parity: ledger mdq_bipolar, source case 4
		{10, true, 1, MDQResult{10, false}}, // parity: ledger mdq_bipolar, source case 5
		{13, true, 3, MDQResult{13, true}},  // parity: ledger mdq_bipolar, source case 6
	}
	for i, c := range cases {
		if got, err := MDQScreen(mdqYes(c.n), c.same, c.problem); err != nil || got != c.want {
			t.Errorf("case %d = (%+v, %v), want %+v", i+1, got, err, c.want)
		}
	}
}

func TestMDQScreen_noAlloc(t *testing.T) {
	in := mdqYes(7)
	if n := testing.AllocsPerRun(100, func() { r, _ := MDQScreen(in, true, 2); sinkBool = r.Positive }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkMDQScreen(b *testing.B) {
	b.ReportAllocs()
	in := mdqYes(7)
	for b.Loop() {
		r, _ := MDQScreen(in, true, 2)
		sinkBool = r.Positive
	}
}
