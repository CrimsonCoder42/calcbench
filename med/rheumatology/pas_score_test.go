package rheumatology

import (
	"errors"
	"math"
	"testing"
)

// HAND-DERIVED from "(HAQ × 3.33 + pain VAS + PtGA VAS)/3" (Anderson 2012 PMC4028066). Tolerance 1e-12:
// three-term sums.
func TestPASScore_formula(t *testing.T) {
	cases := []struct{ haq, pain, global, want float64 }{
		{0, 0, 0, 0},
		{1, 5, 4, (3.33 + 9) / 3},
		{3, 10, 10, (9.99 + 20) / 3},
		{0.5, 2, 3, (1.665 + 5) / 3},
	}
	for _, c := range cases {
		got, err := PASScore(c.haq, c.pain, c.global)
		if err != nil || math.Abs(got-c.want) > 1e-12 {
			t.Errorf("%v: got %v, %v; want %v", c, got, err, c.want)
		}
	}
}

// parity: ledger pas_ra and pas_ii_ra, source case 1: the consumer averages its HAQ (or HAQ-II) items to a
// 0–3 score and applies (HAQ × 3.33 + pain + global)/3; HAQ-II 1.2, pain 6, global 5 → (3.996 + 11)/3 = 14.996/3 ≈ 4.99867.
func TestPASScore_parity(t *testing.T) {
	got, err := PASScore(1.2, 6, 5)
	if err != nil || math.Abs(got-14.996/3) > 1e-12 {
		t.Errorf("got %v, %v; want 14.996/3", got, err)
	}
}

func TestPASScore_errors(t *testing.T) {
	cases := []struct {
		name             string
		haq, pain, globl float64
		want             error
	}{
		{"haq NaN", math.NaN(), 1, 1, ErrNotFinite},
		{"pain +Inf", 1, math.Inf(1), 1, ErrNotFinite},
		{"global -Inf", 1, 1, math.Inf(-1), ErrNotFinite},
		{"haq 3.01", 3.01, 1, 1, ErrOutOfRange},
		{"haq negative", -0.1, 1, 1, ErrOutOfRange},
		{"pain 10.1", 1, 10.1, 1, ErrOutOfRange},
		{"global negative", 1, 1, -1, ErrOutOfRange},
	}
	for _, c := range cases {
		if _, err := PASScore(c.haq, c.pain, c.globl); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

// Anderson 2012 Table 3 / England 2019 Table 1: both sides of each edge.
func TestPASActivityCategory_edges(t *testing.T) {
	cases := []struct {
		v    float64
		want PASActivity
	}{
		{0, PASRemission}, {0.25, PASRemission}, {0.26, PASLow}, {3.70, PASLow}, {3.71, PASModerate},
		{7.99, PASModerate}, {8.0, PASHigh}, {10, PASHigh},
	}
	for _, c := range cases {
		if got, err := PASActivityCategory(c.v); err != nil || got != c.want {
			t.Errorf("%v: got %q, %v; want %q", c.v, got, err, c.want)
		}
	}
	for _, c := range []struct {
		v    float64
		want error
	}{{math.NaN(), ErrNotFinite}, {math.Inf(1), ErrNotFinite}, {-0.01, ErrOutOfRange}, {10.01, ErrOutOfRange}} {
		if _, err := PASActivityCategory(c.v); !errors.Is(err, c.want) {
			t.Errorf("%v: err = %v, want %v", c.v, err, c.want)
		}
	}
}

func TestPASScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloatRheum, _ = PASScore(1, 5, 4) }); n != 0 {
		t.Errorf("PASScore allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { c, _ := PASActivityCategory(4); sinkBool = c == PASModerate }); n != 0 {
		t.Errorf("PASActivityCategory allocs = %v, want 0", n)
	}
}

func BenchmarkPASScore(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkFloatRheum, _ = PASScore(1, 5, 4)
	}
}

func BenchmarkPASActivityCategory(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c, _ := PASActivityCategory(float64(i%11) * 0.9)
		sinkBool = c == PASHigh
	}
}
