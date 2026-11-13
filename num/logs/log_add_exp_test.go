package logs

import (
	"errors"
	"math"
	"testing"
)

func nearRel(got, want, tol float64) bool {
	if got == want {
		return true
	}
	return math.Abs(got-want) <= tol*math.Max(1, math.Abs(want))
}

// Hand-checked values: log(1 + 3) = log 4; equal arguments add log 2 (Blanchard, Higham & Higham 2021, §2).
func TestLogAddExp_cases(t *testing.T) {
	inf := math.Inf(-1)
	for _, tt := range []struct {
		name       string
		a, b, want float64
	}{
		{"log1 plus log3", 0, math.Log(3), math.Log(4)},
		{"equal large", 1000, 1000, 1000 + math.Ln2},
		{"equal very negative", -1000, -1000, -1000 + math.Ln2},
		{"far apart", 0, -800, 0},
		{"swapped", math.Log(3), 0, math.Log(4)},
		{"one log-zero", inf, 2.5, 2.5},
		{"both log-zero", inf, inf, inf},
		{"max float", math.MaxFloat64, math.MaxFloat64, math.MaxFloat64},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// Tolerance 4e-16 relative: one log1p and one add.
			got, err := LogAddExp(tt.a, tt.b)
			if err != nil || !nearRel(got, tt.want, 4e-16) {
				t.Errorf("LogAddExp(%v, %v) = %v, %v; want %v", tt.a, tt.b, got, err, tt.want)
			}
		})
	}
}

func TestLogAddExp_errors(t *testing.T) {
	for _, in := range [][2]float64{{math.NaN(), 0}, {0, math.NaN()}, {math.Inf(1), 0}, {0, math.Inf(1)}} {
		if _, err := LogAddExp(in[0], in[1]); !errors.Is(err, ErrNotFinite) {
			t.Errorf("LogAddExp(%v, %v) err = %v; want ErrNotFinite", in[0], in[1], err)
		}
	}
}

func TestLogAddExp2_cases(t *testing.T) {
	inf := math.Inf(-1)
	for _, tt := range []struct {
		name       string
		a, b, want float64
	}{
		{"equal", 3, 3, 4},
		{"0 and 1", 0, 1, math.Log2(3)},
		{"equal huge", 2000, 2000, 2001},
		{"one log-zero", inf, -7, -7},
		{"both log-zero", inf, inf, inf},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LogAddExp2(tt.a, tt.b)
			if err != nil || !nearRel(got, tt.want, 4e-16) {
				t.Errorf("LogAddExp2(%v, %v) = %v, %v; want %v", tt.a, tt.b, got, err, tt.want)
			}
		})
	}
}

func TestLogAddExp2_errors(t *testing.T) {
	for _, in := range [][2]float64{{math.NaN(), 0}, {0, math.Inf(1)}} {
		if _, err := LogAddExp2(in[0], in[1]); !errors.Is(err, ErrNotFinite) {
			t.Errorf("LogAddExp2(%v, %v) err = %v; want ErrNotFinite", in[0], in[1], err)
		}
	}
}

func TestLogAddExp_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { _, _ = LogAddExp(1, 2); _, _ = LogAddExp2(1, 2) }); n != 0 {
		t.Errorf("allocs = %v; want 0", n)
	}
}

func BenchmarkLogAddExp(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = LogAddExp(1.5, -2.25)
	}
}

func BenchmarkLogAddExp2(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = LogAddExp2(1.5, -2.25)
	}
}
