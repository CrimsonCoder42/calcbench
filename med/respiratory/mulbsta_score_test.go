package respiratory

import (
	"errors"
	"math"
	"testing"
)

// HAND-DERIVED from Guo 2019's points (PMID 31849894): multilobular 5 + lymphocytes 0.6 → 4 + quit smoker 2 +
// hypertension 2 + age 65 → 2 = 15; all items with an acute smoker = 5+4+4+3+2+2 = 20.
func TestMuLBSTAScore_handExample(t *testing.T) {
	f := MuLBSTAFindings{MultilobularInfiltrates: true, LymphocytesE9PerL: 0.6, Smoking: MuLBSTAQuitSmoker, Hypertension: true, AgeYears: 65}
	if got, err := MuLBSTAScore(f); err != nil || got != 15 {
		t.Errorf("got (%d, %v), want 15", got, err)
	}
	all := MuLBSTAFindings{true, 0.5, true, MuLBSTAAcuteSmoker, true, 70}
	if got, err := MuLBSTAScore(all); err != nil || got != 20 {
		t.Errorf("all: got (%d, %v), want 20", got, err)
	}
	none := MuLBSTAFindings{LymphocytesE9PerL: 1.5, AgeYears: 40}
	if got, err := MuLBSTAScore(none); err != nil || got != 0 {
		t.Errorf("none: got (%d, %v), want 0", got, err)
	}
}

// Edges: lymphocytes ≤ 0.8 (0.8 scores, 0.81 does not), age ≥ 60 (59.9 does not, 60 does).
func TestMuLBSTAScore_edges(t *testing.T) {
	cases := []struct {
		lymph, age float64
		want       int
	}{{0.8, 40, 4}, {0.81, 40, 0}, {1.5, 59.9, 0}, {1.5, 60, 2}}
	for _, c := range cases {
		if got, err := MuLBSTAScore(MuLBSTAFindings{LymphocytesE9PerL: c.lymph, AgeYears: c.age}); err != nil || got != c.want {
			t.Errorf("lymph %v age %v: got (%d, %v), want %d", c.lymph, c.age, got, err, c.want)
		}
	}
}

func TestMuLBSTAScore_errors(t *testing.T) {
	tests := []struct {
		name    string
		f       MuLBSTAFindings
		wantErr error
	}{
		{"negative lymphocytes", MuLBSTAFindings{LymphocytesE9PerL: -0.1, AgeYears: 50}, ErrOutOfRange},
		{"negative age", MuLBSTAFindings{LymphocytesE9PerL: 1, AgeYears: -1}, ErrOutOfRange},
		{"unknown smoking", MuLBSTAFindings{LymphocytesE9PerL: 1, AgeYears: 50, Smoking: 3}, ErrOutOfRange},
		{"NaN lymphocytes", MuLBSTAFindings{LymphocytesE9PerL: math.NaN(), AgeYears: 50}, ErrNotFinite},
		{"+Inf age", MuLBSTAFindings{LymphocytesE9PerL: 1, AgeYears: math.Inf(1)}, ErrNotFinite},
		{"-Inf lymphocytes", MuLBSTAFindings{LymphocytesE9PerL: math.Inf(-1), AgeYears: 50}, ErrNotFinite},
	}
	for _, tc := range tests {
		if _, err := MuLBSTAScore(tc.f); !errors.Is(err, tc.wantErr) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.wantErr)
		}
	}
}

func TestMuLBSTAScore_noAlloc(t *testing.T) {
	f := MuLBSTAFindings{MultilobularInfiltrates: true, LymphocytesE9PerL: 0.6, AgeYears: 65}
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = MuLBSTAScore(f) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkMuLBSTAScore(b *testing.B) {
	b.ReportAllocs()
	f := MuLBSTAFindings{MultilobularInfiltrates: true, LymphocytesE9PerL: 0.6, Smoking: MuLBSTAAcuteSmoker, AgeYears: 65}
	for b.Loop() {
		sinkInt, _ = MuLBSTAScore(f)
	}
}
