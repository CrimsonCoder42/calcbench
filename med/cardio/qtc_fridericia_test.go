package cardio

import (
	"errors"
	"math"
	"testing"
)

// Hand-derived from the cube-root relation: QT 400 ms at RR 0.8 s → 400 / ∛0.8.
// ∛0.8 = 0.928317766722556; 400 / 0.928317766722556 = 430.88693800637673.
func TestQTcFridericia_handExample(t *testing.T) {
	got, err := QTcFridericia(400, 0.8, 40)
	if err != nil || math.Abs(got-430.88693800637673) > 1e-9 {
		t.Errorf("QTcFridericia(400, 0.8, 40) = (%v, %v), want (430.886938…, nil)", got, err)
	}
}

// At RR = 1 s the correction is the identity (the normalisation point).
func TestQTcFridericia_identityAtOneSecond(t *testing.T) {
	got, err := QTcFridericia(380, 1, 30)
	if err != nil || got != 380 {
		t.Errorf("QTcFridericia(380, 1, 30) = (%v, %v), want (380, nil)", got, err)
	}
}

func TestQTcFridericia_rangeEdges(t *testing.T) {
	rrLo, rrHi := 60.0/135, 60.0/51 // pulse 135/min and 51/min
	tests := []struct {
		name        string
		qt, rr, age float64
		wantErr     error
	}{
		{"rr at 135/min", 350, rrLo, 40, nil},
		{"rr at 51/min", 350, rrHi, 40, nil},
		{"rr just faster than 135/min", 350, rrLo - 1e-9, 40, ErrOutOfRange},
		{"rr just slower than 51/min", 350, rrHi + 1e-9, 40, ErrOutOfRange},
		{"age 2", 350, 1, 2, nil},
		{"age 87", 350, 1, 87, nil},
		{"age just below 2", 350, 1, 1.999, ErrOutOfRange},
		{"age just above 87", 350, 1, 87.001, ErrOutOfRange},
		{"qt zero", 0, 1, 40, ErrOutOfRange},
		{"qt negative", -1, 1, 40, ErrOutOfRange},
		{"qt NaN", math.NaN(), 1, 40, ErrNotFinite},
		{"rr +Inf", 350, math.Inf(1), 40, ErrNotFinite},
		{"age -Inf", 350, 1, math.Inf(-1), ErrNotFinite},
	}
	for _, tc := range tests {
		_, err := QTcFridericia(tc.qt, tc.rr, tc.age)
		if !errors.Is(err, tc.wantErr) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.wantErr)
		}
	}
}

func TestQTcFridericia_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = QTcFridericia(400, 0.8, 40) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkQTcFridericia(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = QTcFridericia(400, 0.8, 40)
	}
}
