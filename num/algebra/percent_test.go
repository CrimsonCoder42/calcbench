package algebra

import (
	"errors"
	"math"
	"testing"
)

// closePct compares with a relative tolerance of 1e-12: each function is two or three roundings, so its
// error is a few ulps; the tolerance only absorbs that.
func closePct(got, want float64) bool {
	return math.Abs(got-want) <= 1e-12*math.Max(1, math.Abs(want))
}

// Kirk, Contemporary Mathematics, OpenStax, 2023, §6.1: Example 6.4 (70 % of 3 500 = 2 450; 156 % of 720 =
// 1 123.2), Example 6.5 (35 % of the total is 70, total 200; 10 % of the total is 4 000, total 40 000),
// Example 6.6 (175 is 35 % of 500).
func TestPercent_workedExamples(t *testing.T) {
	if got, err := PercentOf(70, 3500); err != nil || !closePct(got, 2450) {
		t.Errorf("PercentOf(70, 3500) = %v, %v", got, err)
	}
	if got, err := PercentOf(156, 720); err != nil || !closePct(got, 1123.2) {
		t.Errorf("PercentOf(156, 720) = %v, %v", got, err)
	}
	if got, err := TotalFromPercent(70, 35); err != nil || !closePct(got, 200) {
		t.Errorf("TotalFromPercent(70, 35) = %v, %v", got, err)
	}
	if got, err := TotalFromPercent(4000, 10); err != nil || !closePct(got, 40000) {
		t.Errorf("TotalFromPercent(4000, 10) = %v, %v", got, err)
	}
	if got, err := PercentFromPart(175, 500); err != nil || !closePct(got, 35) {
		t.Errorf("PercentFromPart(175, 500) = %v, %v", got, err)
	}
}

// Hand computations: 40 → 50 is 10/40 = +25 %; 50 → 40 is −10/50 = −20 %; −40 → −50 is −10/40 = −25 %;
// 500 ± 10 % = 550 and 450.
func TestPercentChange_cases(t *testing.T) {
	tests := []struct{ from, to, want float64 }{
		{40, 50, 25}, {50, 40, -20}, {-40, -50, -25}, {7, 7, 0}, {1, 0, -100},
	}
	for _, tc := range tests {
		if got, err := PercentChange(tc.from, tc.to); err != nil || !closePct(got, tc.want) {
			t.Errorf("PercentChange(%v, %v) = %v, %v; want %v", tc.from, tc.to, got, err, tc.want)
		}
	}
	if got, err := ApplyPercentChange(500, 10); err != nil || !closePct(got, 550) {
		t.Errorf("ApplyPercentChange(500, 10) = %v, %v", got, err)
	}
	if got, err := ApplyPercentChange(500, -10); err != nil || !closePct(got, 450) {
		t.Errorf("ApplyPercentChange(500, -10) = %v, %v", got, err)
	}
	if got, err := ApplyPercentChange(0, 50); err != nil || got != 0 {
		t.Errorf("ApplyPercentChange(0, 50) = %v, %v", got, err)
	}
}

// Hand computations: |40 − 50|/45 = 22.2222… %; equal values 0 %; 1 and 3 → 2/2 = 100 %.
func TestPercentDifference_cases(t *testing.T) {
	tests := []struct{ x, y, want float64 }{
		{40, 50, 200.0 / 9}, {50, 40, 200.0 / 9}, {5, 5, 0}, {1, 3, 100},
		{math.MaxFloat64, math.MaxFloat64 / 2, 200.0 / 3},
	}
	for _, tc := range tests {
		if got, err := PercentDifference(tc.x, tc.y); err != nil || !closePct(got, tc.want) {
			t.Errorf("PercentDifference(%v, %v) = %v, %v; want %v", tc.x, tc.y, got, err, tc.want)
		}
	}
}

// Burden, Faires and Burden, Numerical Analysis, 10th ed., 2016, §1.2, Example 1(a): p = 3.000, p* = 3.100,
// relative error 0.3333×10⁻¹ = 3.333… %. Hand: observed 9 against 10 is 10 %; −11 against −10 is 10 %.
func TestPercentError_cases(t *testing.T) {
	tests := []struct{ obs, tru, want float64 }{
		{3.1, 3.0, 10.0 / 3}, {9, 10, 10}, {-11, -10, 10}, {10, 10, 0},
	}
	for _, tc := range tests {
		if got, err := PercentError(tc.obs, tc.tru); err != nil || !closePct(got, tc.want) {
			t.Errorf("PercentError(%v, %v) = %v, %v; want %v", tc.obs, tc.tru, got, err, tc.want)
		}
	}
}

func TestPercent_errors(t *testing.T) {
	nan, inf := math.NaN(), math.Inf(1)
	tests := []struct {
		name string
		err  error
		want error
	}{
		{"PercentOf NaN", second(PercentOf(nan, 1)), ErrNotFinite},
		{"PercentOf overflow", second(PercentOf(1e10, math.MaxFloat64)), ErrOverflow},
		{"PercentFromPart Inf", second(PercentFromPart(inf, 1)), ErrNotFinite},
		{"PercentFromPart zero total", second(PercentFromPart(1, 0)), ErrDegenerate},
		{"PercentFromPart overflow", second(PercentFromPart(math.MaxFloat64, 0.5)), ErrOverflow},
		{"TotalFromPercent zero percent", second(TotalFromPercent(1, 0)), ErrDegenerate},
		{"TotalFromPercent NaN", second(TotalFromPercent(1, nan)), ErrNotFinite},
		{"TotalFromPercent overflow", second(TotalFromPercent(math.MaxFloat64, 1)), ErrOverflow},
		{"PercentChange zero original", second(PercentChange(0, 5)), ErrDegenerate},
		{"PercentChange Inf", second(PercentChange(1, inf)), ErrNotFinite},
		{"PercentChange overflow", second(PercentChange(1e-300, 1e300)), ErrOverflow},
		{"ApplyPercentChange NaN", second(ApplyPercentChange(nan, 1)), ErrNotFinite},
		{"ApplyPercentChange overflow", second(ApplyPercentChange(math.MaxFloat64, 100)), ErrOverflow},
		{"PercentDifference zero", second(PercentDifference(0, 5)), ErrDegenerate},
		{"PercentDifference negative", second(PercentDifference(-1, 5)), ErrDegenerate},
		{"PercentDifference NaN", second(PercentDifference(nan, 5)), ErrNotFinite},
		{"PercentError zero true", second(PercentError(1, 0)), ErrDegenerate},
		{"PercentError Inf", second(PercentError(inf, 1)), ErrNotFinite},
		{"PercentError overflow", second(PercentError(math.MaxFloat64, 1e-10)), ErrOverflow},
	}
	for _, tc := range tests {
		if !errors.Is(tc.err, tc.want) {
			t.Errorf("%s: err = %v, want %v", tc.name, tc.err, tc.want)
		}
	}
}

func second(_ float64, err error) error { return err }

func TestPercent_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() {
		_, _ = PercentOf(70, 3500)
		_, _ = PercentChange(40, 50)
		_, _ = PercentDifference(40, 50)
		_, _ = PercentError(3.1, 3)
	}); a != 0 {
		t.Errorf("percent functions allocate %v", a)
	}
}

func BenchmarkPercentOf(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = PercentOf(70, 3500)
	}
}

func BenchmarkPercentFromPart(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = PercentFromPart(175, 500)
	}
}

func BenchmarkTotalFromPercent(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = TotalFromPercent(70, 35)
	}
}

func BenchmarkPercentChange(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = PercentChange(40, 50)
	}
}

func BenchmarkApplyPercentChange(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = ApplyPercentChange(500, 10)
	}
}

func BenchmarkPercentDifference(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = PercentDifference(40, 50)
	}
}

func BenchmarkPercentError(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = PercentError(3.1, 3)
	}
}
