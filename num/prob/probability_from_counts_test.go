package prob

import (
	"errors"
	"math"
	"testing"
)

// Shared test helpers for the package (declared here, with the package's first concept).

// within reports whether got is within relTol of want, relative to |want| (absolute when want is 0).
func within(got, want, relTol float64) bool {
	if want == 0 {
		return math.Abs(got) <= relTol
	}
	return math.Abs(got-want) <= relTol*math.Abs(want)
}

// Sinks keep benchmarks and allocation checks from being optimised away.
var (
	sinkFloat float64
	sinkInt   int
	sinkSlice []float64
)

// requireNoAlloc fails when f allocates.
func requireNoAlloc(t *testing.T, name string, f func()) {
	t.Helper()
	if a := testing.AllocsPerRun(100, f); a != 0 {
		t.Errorf("%s allocates %v times per call, want 0", name, a)
	}
}

// requireErr fails unless err wraps target.
func requireErr(t *testing.T, call string, err, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Errorf("%s: err = %v, want %v", call, err, target)
	}
}

var nan, inf = math.NaN(), math.Inf(1)

// Ross §2.5 (equally likely outcomes): a fair die, the event "even" is 3 of 6 outcomes. Hand-checked.
func TestProbabilityFromCounts_sourceExample(t *testing.T) {
	if got, err := ProbabilityFromCounts(3, 6); err != nil || got != 0.5 {
		t.Errorf("ProbabilityFromCounts(3, 6) = (%v, %v), want (0.5, nil)", got, err)
	}
}

func TestProbabilityFromCounts_cases(t *testing.T) {
	tests := []struct {
		name             string
		favorable, total int
		want             float64
	}{
		{"none", 0, 6, 0},
		{"all", 6, 6, 1},
		{"one of one", 1, 1, 1},
		{"one third", 1, 3, 1.0 / 3},
		{"max int", math.MaxInt / 2, math.MaxInt, 0.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := ProbabilityFromCounts(tt.favorable, tt.total); err != nil || got != tt.want {
				t.Errorf("ProbabilityFromCounts(%d, %d) = (%v, %v), want (%v, nil)", tt.favorable, tt.total, got, err, tt.want)
			}
		})
	}
}

func TestProbabilityFromCounts_errors(t *testing.T) {
	for _, c := range [][2]int{{0, 0}, {1, -1}, {-1, 5}, {6, 5}, {math.MinInt, 1}} {
		_, err := ProbabilityFromCounts(c[0], c[1])
		requireErr(t, "ProbabilityFromCounts", err, ErrOutOfRange)
	}
}

func TestProbabilityFromCounts_noAlloc(t *testing.T) {
	requireNoAlloc(t, "ProbabilityFromCounts", func() { sinkFloat, _ = ProbabilityFromCounts(3, 6) })
}

func BenchmarkProbabilityFromCounts(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = ProbabilityFromCounts(3, 6)
	}
}
