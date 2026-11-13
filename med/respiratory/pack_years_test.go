package respiratory

import (
	"errors"
	"math"
	"testing"
)

// USPSTF 2021 definition: 1 pack (20 cigarettes) per day for a year is 1 pack-year. Exact comparisons where the
// arithmetic is exact in binary; 1e-12 otherwise.
func TestPackYears_definition(t *testing.T) {
	if got, err := PackYears(1, 1); err != nil || got != 1 {
		t.Errorf("1 pack/day × 1 y = (%v, %v), want 1", got, err)
	}
	if got, err := PackYearsFromCigarettes(20, 1); err != nil || got != 1 {
		t.Errorf("20 cig/day × 1 y = (%v, %v), want 1", got, err)
	}
	if got, err := PackYearsFromCigarettes(10, 30); err != nil || got != 15 {
		t.Errorf("10 cig/day × 30 y = (%v, %v), want 15", got, err)
	}
}

// parity: ledger pack_years, source cases 1–6. The source takes ages started/stopped and non-smoking years;
// the years smoked are that subtraction, done here by the caller as the consumer will. Case 6 (stopped before
// started) is refused by the source; here the caller's subtraction gives a negative yearsSmoked → ErrOutOfRange.
func TestPackYears_parity(t *testing.T) {
	cases := []struct {
		ppd, started, stopped, nonsmoking float64
		want                              float64
		wantErr                           error
	}{
		{1, 18, 50, 0, 32, nil},          // source case 1
		{0.5, 20, 40, 0, 10, nil},        // source case 2
		{1, 20, 40, 0, 20, nil},          // source case 3
		{0.3, 16, 30, 2, 3.6, nil},       // source case 4 (source shows 3.6 at one decimal)
		{2, 15, 30, 0, 30, nil},          // source case 5
		{1, 40, 30, 0, 0, ErrOutOfRange}, // source case 6: source refuses
	}
	for i, c := range cases {
		got, err := PackYears(c.ppd, c.stopped-c.started-c.nonsmoking)
		if c.wantErr != nil {
			if !errors.Is(err, c.wantErr) {
				t.Errorf("case %d: err = %v, want %v", i+1, err, c.wantErr)
			}
			continue
		}
		if err != nil || math.Abs(got-c.want) > 1e-12 {
			t.Errorf("case %d: got (%v, %v), want %v", i+1, got, err, c.want)
		}
	}
}

func TestPackYears_edges(t *testing.T) {
	if got, err := PackYears(0, 10); err != nil || got != 0 {
		t.Errorf("zero packs: (%v, %v), want 0", got, err)
	}
	if got, err := PackYears(1, 0); err != nil || got != 0 {
		t.Errorf("zero years: (%v, %v), want 0", got, err)
	}
	tests := []struct {
		name    string
		a, b    float64
		cig     bool
		wantErr error
	}{
		{"negative packs", -0.1, 10, false, ErrOutOfRange},
		{"negative years", 1, -1, false, ErrOutOfRange},
		{"NaN packs", math.NaN(), 10, false, ErrNotFinite},
		{"+Inf years", 1, math.Inf(1), false, ErrNotFinite},
		{"-Inf packs", math.Inf(-1), 1, false, ErrNotFinite},
		{"overflow", math.MaxFloat64, 2, false, ErrDomain},
		{"negative cigarettes", -1, 10, true, ErrOutOfRange},
		{"NaN cigarettes", math.NaN(), 10, true, ErrNotFinite},
		{"cigarettes, NaN years", 20, math.NaN(), true, ErrNotFinite},
	}
	for _, tc := range tests {
		var err error
		if tc.cig {
			_, err = PackYearsFromCigarettes(tc.a, tc.b)
		} else {
			_, err = PackYears(tc.a, tc.b)
		}
		if !errors.Is(err, tc.wantErr) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.wantErr)
		}
	}
}

func TestPackYears_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = PackYears(1, 20) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = PackYearsFromCigarettes(15, 20) }); n != 0 {
		t.Errorf("cigarettes allocs = %v, want 0", n)
	}
}

func BenchmarkPackYears(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = PackYears(1.5, 20)
	}
}

func BenchmarkPackYearsFromCigarettes(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = PackYearsFromCigarettes(15, 20)
	}
}
