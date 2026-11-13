package cardio

import (
	"math"
	"testing"
)

// HAND-DERIVED from the printed coefficients. Tolerance 1e-12 absolute (four products and sums of values < 100).
func TestSubtleAnteriorSTEMI4Variable_values(t *testing.T) {
	got, err := SubtleAnteriorSTEMI4Variable(400, 15, 10, 2)
	// 0.052·400 − 0.151·15 − 0.268·10 + 1.062·2 = 20.8 − 2.265 − 2.68 + 2.124 = 17.979 (below the 18.2 cut-point)
	if err != nil || math.Abs(got-17.979) > 1e-12 {
		t.Errorf("got (%v, %v), want 17.979", got, err)
	}
	got, err = SubtleAnteriorSTEMI4Variable(420, 10, 5, 3)
	// 21.84 − 1.51 − 1.34 + 3.186 = 22.176
	if err != nil || math.Abs(got-22.176) > 1e-12 {
		t.Errorf("got (%v, %v), want 22.176", got, err)
	}
}

func TestSubtleAnteriorSTEMI3Variable_values(t *testing.T) {
	got, err := SubtleAnteriorSTEMI3Variable(400, 10, 2)
	// 1.196·2 + 0.059·400 − 0.326·10 = 2.392 + 23.6 − 3.26 = 22.732
	if err != nil || math.Abs(got-22.732) > 1e-12 {
		t.Errorf("got (%v, %v), want 22.732", got, err)
	}
}

func TestSubtleAnteriorSTEMI_errors(t *testing.T) {
	nan := math.NaN()
	_, err := SubtleAnteriorSTEMI4Variable(nan, 1, 1, 1)
	checkErr(t, "qtc NaN", err, ErrNotFinite)
	_, err = SubtleAnteriorSTEMI4Variable(400, math.Inf(1), 1, 1)
	checkErr(t, "qrs +Inf", err, ErrNotFinite)
	_, err = SubtleAnteriorSTEMI4Variable(400, -0.1, 1, 1)
	checkErr(t, "qrs < 0", err, ErrOutOfRange)
	_, err = SubtleAnteriorSTEMI3Variable(0, 1, 1)
	checkErr(t, "qtc 0", err, ErrOutOfRange)
	_, err = SubtleAnteriorSTEMI3Variable(400, -1, 1)
	checkErr(t, "rv4 < 0", err, ErrOutOfRange)
	_, err = SubtleAnteriorSTEMI3Variable(400, 1, math.Inf(-1))
	checkErr(t, "ste -Inf", err, ErrNotFinite)
	if _, err := SubtleAnteriorSTEMI3Variable(400, 0, 0); err != nil {
		t.Errorf("zero amplitudes: %v", err)
	}
}

func TestSubtleAnteriorSTEMI_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() {
		sinkFloat, _ = SubtleAnteriorSTEMI4Variable(400, 15, 10, 2)
		sinkFloat, _ = SubtleAnteriorSTEMI3Variable(400, 10, 2)
	}); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkSubtleAnteriorSTEMI4Variable(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = SubtleAnteriorSTEMI4Variable(400, 15, 10, 2)
	}
}

func BenchmarkSubtleAnteriorSTEMI3Variable(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = SubtleAnteriorSTEMI3Variable(400, 10, 2)
	}
}
