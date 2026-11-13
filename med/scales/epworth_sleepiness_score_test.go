package scales

import (
	"errors"
	"fmt"
	"testing"
)

// Shared by the package's tests.

// sinks keep benchmark results observable.
var (
	sinkInt       int
	sinkEpworth   EpworthSleepinessCategory
	sinkSeverity  GCSSeverity
	sinkCondition ApgarCondition
)

// errCase is one call that must fail with want.
type errCase struct {
	name string
	call func() error
	want error
}

func runErrCases(t *testing.T, cases []errCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.call(); !errors.Is(err, tc.want) {
				t.Errorf("err = %v, want %v", err, tc.want)
			}
		})
	}
}

// Johns 1991, Table 1: 8 items, each rated 0–3, total 0–24. HAND-DERIVED sums: eight 1s → 8 (inside the
// controls' reported 2–10); eight 3s → 24 (the maximum); eight 0s → 0 (the minimum); 0+1+2+3+0+1+2+3 → 12.
func TestEpworthSleepinessScore_example(t *testing.T) {
	tests := []struct {
		name  string
		items [8]int
		want  int
	}{
		{"all 1", [8]int{1, 1, 1, 1, 1, 1, 1, 1}, 8},
		{"all 3", [8]int{3, 3, 3, 3, 3, 3, 3, 3}, 24},
		{"all 0", [8]int{}, 0},
		{"mixed", [8]int{0, 1, 2, 3, 0, 1, 2, 3}, 12},
	}
	for _, tc := range tests {
		if got, err := EpworthSleepinessScore(tc.items); err != nil || got != tc.want {
			t.Errorf("%s: EpworthSleepinessScore(%v) = (%v, %v), want %v", tc.name, tc.items, got, err, tc.want)
		}
	}
}

// Each item on both sides of its range: −1 and 4 at every index fail; 0 and 3 pass (covered above).
func TestEpworthSleepinessScore_errors(t *testing.T) {
	var cases []errCase
	for i := range 8 {
		for _, bad := range []int{-1, 4} {
			items := [8]int{1, 1, 1, 1, 1, 1, 1, 1}
			items[i] = bad
			cases = append(cases, errCase{
				name: fmt.Sprintf("items[%d] = %d", i, bad),
				call: func() error { _, err := EpworthSleepinessScore(items); return err },
				want: ErrOutOfRange,
			})
		}
	}
	runErrCases(t, cases)
}

func TestEpworthSleepinessScore_noAllocation(t *testing.T) {
	items := [8]int{1, 2, 0, 3, 1, 2, 0, 1}
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = EpworthSleepinessScore(items) }); n != 0 {
		t.Errorf("EpworthSleepinessScore allocates %v times, want 0", n)
	}
}

func BenchmarkEpworthSleepinessScore(b *testing.B) {
	b.ReportAllocs()
	items := [8]int{1, 2, 0, 3, 1, 2, 0, 1}
	for b.Loop() {
		sinkInt, _ = EpworthSleepinessScore(items)
	}
}
