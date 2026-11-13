package fitness

import (
	"errors"
	"math"
	"testing"
)

// parity: ledger Pace, source cases 1–3 (1860 s over 3.1 mi = 600 s/mi; zero distance and zero time refused).
func TestPace_parity(t *testing.T) {
	got, err := Pace(1860, 3.1)
	// 1860/3.1 in float64 is within 1e-12 of 600.
	if err != nil || math.Abs(got-600) > 1e-12 {
		t.Errorf("got %v, %v; want 600", got, err)
	}
	if _, err := Pace(1860, 0); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("zero distance: %v", err)
	}
	if _, err := Pace(0, 1); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("zero time: %v", err)
	}
}

func TestPace_cases(t *testing.T) {
	// 5 km in 25 min → 300 s/km (HAND-DERIVED).
	if got, _ := Pace(1500, 5); got != 300 {
		t.Errorf("5 km: %v", got)
	}
	cases := []struct {
		name   string
		sec, d float64
		want   error
	}{
		{"negative time", -1, 1, ErrOutOfRange},
		{"negative distance", 1, -1, ErrOutOfRange},
		{"NaN time", math.NaN(), 1, ErrNotFinite},
		{"Inf distance", 1, math.Inf(1), ErrNotFinite},
		{"-Inf time", math.Inf(-1), 1, ErrNotFinite},
		{"overflow", math.MaxFloat64, 1e-300, ErrDomain},
	}
	for _, c := range cases {
		if _, err := Pace(c.sec, c.d); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
	// Just inside the range on both inputs.
	if _, err := Pace(math.SmallestNonzeroFloat64, 1); err != nil {
		t.Errorf("tiny time: %v", err)
	}
}

func TestPace_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { _, _ = Pace(1860, 3.1) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkPace(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = Pace(1860, 3.1)
	}
}
