package prob

import (
	"math"
	"testing"
)

// closedForm bundles one parameterisation of a closed-form continuous distribution for the shared checks.
type closedForm struct {
	name     string
	pdf, cdf func(x float64) (float64, error)
	quantile func(u float64) (float64, error)
	mv       func() (float64, float64, error)
	lo, hi   float64 // integration range covering all but a negligible tail
}

// checkClosedForm asserts CDF(Q(u)) = u, ∫ f over [a, b] = F(b) − F(a), and the moments by quadrature.
func checkClosedForm(t *testing.T, d closedForm) {
	t.Helper()
	for _, u := range []float64{1e-10, 0.01, 0.3, 0.5, 0.77, 0.99} {
		x, err := d.quantile(u)
		if err != nil {
			t.Fatalf("%s: Q(%v): %v", d.name, u, err)
		}
		// Rounding x to a double perturbs it by up to 2⁻⁵³ relative, which moves F by a relative amount
		// κ·2⁻⁵³ with κ = |x·f(x)/F(x)| (the condition number of F at x). Where κ is large (Pareto at
		// u = 1e-10: x = 1.5·(1 + 1.25e-11), κ ≈ 8e10, so ≈ 1e-5 relative) no double x can do better, so
		// the tolerance is the algorithm's 1e-13 plus that unavoidable representation error.
		fx, _ := d.pdf(x)
		kappa := math.Abs(x * fx / u)
		if c, _ := d.cdf(x); math.Abs(c-u) > (1e-13+4*kappa*0x1p-53)*u+4e-16 {
			t.Errorf("%s: CDF(Q(%v)) = %v", d.name, u, c)
		}
	}
	f := func(x float64) float64 { v, _ := d.pdf(x); return v }
	a, b := d.lo+(d.hi-d.lo)*0.2, d.lo+(d.hi-d.lo)*0.35
	ca, _ := d.cdf(a)
	cb, _ := d.cdf(b)
	if got := simpson(f, a, b, 2000); math.Abs(got-(cb-ca)) > 1e-11 {
		t.Errorf("%s: ∫f = %v, CDF diff = %v", d.name, got, cb-ca)
	}
	m, v, err := d.mv()
	if err != nil {
		t.Fatalf("%s: mean/variance: %v", d.name, err)
	}
	m1 := simpson(func(x float64) float64 { return x * f(x) }, d.lo, d.hi, 200000)
	m2 := simpson(func(x float64) float64 { return (x - m) * (x - m) * f(x) }, d.lo, d.hi, 200000)
	if math.Abs(m1-m) > 1e-7*(1+math.Abs(m)) || math.Abs(m2-v) > 1e-7*(1+v) {
		t.Errorf("%s: quadrature (%v, %v), closed form (%v, %v)", d.name, m1, m2, m, v)
	}
}

// JKB vol. 2 ch. 22 §1: f(0) = F(0) = e⁻¹ for μ = 0, β = 1.
func TestGumbel(t *testing.T) {
	if d, err := GumbelPDF(0, 0, 1); err != nil || !within(d, math.Exp(-1), 1e-15) {
		t.Errorf("PDF(0) = (%v, %v)", d, err)
	}
	if c, err := GumbelCDF(0, 0, 1); err != nil || !within(c, math.Exp(-1), 1e-15) {
		t.Errorf("CDF(0) = (%v, %v)", c, err)
	}
	if q, err := GumbelQuantile(math.Exp(-1), 0, 1); err != nil || math.Abs(q) > 1e-15 {
		t.Errorf("Q(e⁻¹) = (%v, %v)", q, err)
	}
	if m, v, _ := GumbelMeanVariance(0, 1); !within(m, 0.5772156649015329, 1e-15) || !within(v, math.Pi*math.Pi/6, 1e-15) {
		t.Errorf("MV = (%v, %v)", m, v)
	}
	if d, _ := GumbelPDF(-1e308, 1e308, 1); d != 0 {
		t.Errorf("far tail = %v", d)
	}
	checkClosedForm(t, closedForm{"gumbel",
		func(x float64) (float64, error) { return GumbelPDF(x, 1, 2) },
		func(x float64) (float64, error) { return GumbelCDF(x, 1, 2) },
		func(u float64) (float64, error) { return GumbelQuantile(u, 1, 2) },
		func() (float64, float64, error) { return GumbelMeanVariance(1, 2) }, -15, 80})
	_, err := GumbelPDF(nan, 0, 1)
	requireErr(t, "x NaN", err, ErrNotFinite)
	_, err = GumbelCDF(0, 0, 0)
	requireErr(t, "scale 0", err, ErrOutOfRange)
	_, err = GumbelPDF(0, 0, 1e-310)
	requireErr(t, "pdf overflow", err, ErrOverflow)
	_, err = GumbelQuantile(0, 0, 1)
	requireErr(t, "u = 0", err, ErrOutOfRange)
	_, err = GumbelQuantile(1, 0, 1)
	requireErr(t, "u = 1", err, ErrOutOfRange)
	_, err = GumbelQuantile(2, 0, 1)
	requireErr(t, "u > 1", err, ErrProbability)
	_, err = GumbelQuantile(0.999, 1.7e308, 1e308)
	requireErr(t, "quantile overflow", err, ErrOverflow)
	_, _, err = GumbelMeanVariance(0, 1e200)
	requireErr(t, "mv overflow", err, ErrOverflow)
}

// JKB vol. 2 ch. 24 §1: f(1) = e⁻¹/2, F(1) = 1 − e⁻¹/2; Q(0.75) = ln 2.
func TestLaplace(t *testing.T) {
	if d, err := LaplacePDF(1, 0, 1); err != nil || !within(d, math.Exp(-1)/2, 1e-15) {
		t.Errorf("PDF(1) = (%v, %v)", d, err)
	}
	if c, _ := LaplaceCDF(1, 0, 1); !within(c, 1-math.Exp(-1)/2, 1e-15) {
		t.Errorf("CDF(1) = %v", c)
	}
	if c, _ := LaplaceCDF(-1, 0, 1); !within(c, math.Exp(-1)/2, 1e-15) {
		t.Errorf("CDF(−1) = %v", c)
	}
	if q, _ := LaplaceQuantile(0.75, 0, 1); !within(q, math.Ln2, 1e-15) {
		t.Errorf("Q(0.75) = %v", q)
	}
	if m, v, _ := LaplaceMeanVariance(3, 2); m != 3 || v != 8 {
		t.Errorf("MV = (%v, %v)", m, v)
	}
	checkClosedForm(t, closedForm{"laplace",
		func(x float64) (float64, error) { return LaplacePDF(x, -1, 0.5) },
		func(x float64) (float64, error) { return LaplaceCDF(x, -1, 0.5) },
		func(u float64) (float64, error) { return LaplaceQuantile(u, -1, 0.5) },
		func() (float64, float64, error) { return LaplaceMeanVariance(-1, 0.5) }, -1 - 40, -1 + 40})
	_, err := LaplacePDF(0, inf, 1)
	requireErr(t, "location Inf", err, ErrNotFinite)
	_, err = LaplaceCDF(0, 0, -1)
	requireErr(t, "scale < 0", err, ErrOutOfRange)
	_, err = LaplacePDF(0, 0, 1e-310)
	requireErr(t, "pdf overflow", err, ErrOverflow)
	_, err = LaplaceQuantile(1, 0, 1)
	requireErr(t, "u = 1", err, ErrOutOfRange)
	_, err = LaplaceQuantile(nan, 0, 1)
	requireErr(t, "u NaN", err, ErrNotFinite)
	_, err = LaplaceQuantile(1e-300, -1.7e308, 1e308)
	requireErr(t, "quantile overflow", err, ErrOverflow)
	_, _, err = LaplaceMeanVariance(0, 1e200)
	requireErr(t, "mv overflow", err, ErrOverflow)
}

// JKB vol. 2 ch. 23 §1: f(0) = ¼; F(ln 3) = ¾; lower tail F(−800) = e⁻⁸⁰⁰ underflows to 0, F(−30) ≈ e⁻³⁰.
func TestLogistic(t *testing.T) {
	if d, err := LogisticPDF(0, 0, 1); err != nil || d != 0.25 {
		t.Errorf("PDF(0) = (%v, %v)", d, err)
	}
	if c, _ := LogisticCDF(math.Log(3), 0, 1); !within(c, 0.75, 1e-15) {
		t.Errorf("CDF(ln 3) = %v", c)
	}
	if c, _ := LogisticCDF(-30, 0, 1); !within(c, math.Exp(-30)/(1+math.Exp(-30)), 1e-15) {
		t.Errorf("CDF(−30) = %v", c)
	}
	if q, _ := LogisticQuantile(0.75, 0, 1); !within(q, math.Log(3), 1e-15) {
		t.Errorf("Q(0.75) = %v", q)
	}
	if m, v, _ := LogisticMeanVariance(0, 1); m != 0 || !within(v, math.Pi*math.Pi/3, 1e-15) {
		t.Errorf("MV = (%v, %v)", m, v)
	}
	checkClosedForm(t, closedForm{"logistic",
		func(x float64) (float64, error) { return LogisticPDF(x, 2, 1.5) },
		func(x float64) (float64, error) { return LogisticCDF(x, 2, 1.5) },
		func(u float64) (float64, error) { return LogisticQuantile(u, 2, 1.5) },
		func() (float64, float64, error) { return LogisticMeanVariance(2, 1.5) }, 2 - 80, 2 + 80})
	_, err := LogisticPDF(nan, 0, 1)
	requireErr(t, "x NaN", err, ErrNotFinite)
	_, err = LogisticCDF(0, 0, 0)
	requireErr(t, "scale 0", err, ErrOutOfRange)
	_, err = LogisticPDF(0, 0, 1e-310)
	requireErr(t, "pdf overflow", err, ErrOverflow)
	_, err = LogisticQuantile(0, 0, 1)
	requireErr(t, "u = 0", err, ErrOutOfRange)
	_, err = LogisticQuantile(0.999, 1.7e308, 1e308)
	requireErr(t, "quantile overflow", err, ErrOverflow)
	_, _, err = LogisticMeanVariance(0, 1e200)
	requireErr(t, "mv overflow", err, ErrOverflow)
}

// JKB vol. 1 ch. 20: x_m = 1, α = 2: f(2) = ¼, F(2) = ¾, Q(¾) = 2; α = 3 → (1.5, 0.75).
func TestPareto(t *testing.T) {
	if d, err := ParetoPDF(2, 1, 2); err != nil || !within(d, 0.25, 1e-15) {
		t.Errorf("PDF(2) = (%v, %v)", d, err)
	}
	if d, _ := ParetoPDF(0.5, 1, 2); d != 0 {
		t.Errorf("below x_m = %v", d)
	}
	if c, _ := ParetoCDF(2, 1, 2); !within(c, 0.75, 1e-15) {
		t.Errorf("CDF(2) = %v", c)
	}
	if q, _ := ParetoQuantile(0.75, 1, 2); !within(q, 2, 1e-15) {
		t.Errorf("Q(¾) = %v", q)
	}
	if q, _ := ParetoQuantile(0, 3, 2); q != 3 {
		t.Errorf("Q(0) = %v", q)
	}
	if m, v, _ := ParetoMeanVariance(1, 3); !within(m, 1.5, 1e-15) || !within(v, 0.75, 1e-15) {
		t.Errorf("MV = (%v, %v)", m, v)
	}
	// Moments by quadrature need a fast tail: α = 8 on [x_m, 60] (tail mass (1.5/60)⁸ ≈ 1.5e-13).
	checkClosedForm(t, closedForm{"pareto",
		func(x float64) (float64, error) { return ParetoPDF(x, 1.5, 8) },
		func(x float64) (float64, error) { return ParetoCDF(x, 1.5, 8) },
		func(u float64) (float64, error) { return ParetoQuantile(u, 1.5, 8) },
		func() (float64, float64, error) { return ParetoMeanVariance(1.5, 8) }, 1.5, 60})
	_, err := ParetoPDF(nan, 1, 2)
	requireErr(t, "x NaN", err, ErrNotFinite)
	_, err = ParetoCDF(2, 0, 2)
	requireErr(t, "scale 0", err, ErrOutOfRange)
	_, err = ParetoPDF(1e-300, 1e-300, 1e300)
	requireErr(t, "pdf overflow", err, ErrOverflow)
	_, err = ParetoQuantile(1, 1, 2)
	requireErr(t, "u = 1", err, ErrOutOfRange)
	_, err = ParetoQuantile(-1, 1, 2)
	requireErr(t, "u < 0", err, ErrProbability)
	_, err = ParetoQuantile(0.99, 1e300, 0.01)
	requireErr(t, "quantile overflow", err, ErrOverflow)
	_, _, err = ParetoMeanVariance(1, 2)
	requireErr(t, "α ≤ 2", err, ErrOutOfRange)
	_, _, err = ParetoMeanVariance(1e300, 2.0000001)
	requireErr(t, "mv overflow", err, ErrOverflow)
}

// JKB vol. 1 ch. 18: σ = 1: f(1) = e^(−½), F(1) = 1 − e^(−½), Q(1 − e⁻²) = 2.
func TestRayleigh(t *testing.T) {
	if d, err := RayleighPDF(1, 1); err != nil || !within(d, math.Exp(-0.5), 1e-15) {
		t.Errorf("PDF(1) = (%v, %v)", d, err)
	}
	if c, _ := RayleighCDF(1, 1); !within(c, -math.Expm1(-0.5), 1e-15) {
		t.Errorf("CDF(1) = %v", c)
	}
	if q, _ := RayleighQuantile(-math.Expm1(-2), 1); !within(q, 2, 1e-15) {
		t.Errorf("Q = %v", q)
	}
	if m, v, _ := RayleighMeanVariance(1); !within(m, math.Sqrt(math.Pi/2), 1e-15) || !within(v, (4-math.Pi)/2, 1e-15) {
		t.Errorf("MV = (%v, %v)", m, v)
	}
	checkClosedForm(t, closedForm{"rayleigh",
		func(x float64) (float64, error) { return RayleighPDF(x, 2) },
		func(x float64) (float64, error) { return RayleighCDF(x, 2) },
		func(u float64) (float64, error) { return RayleighQuantile(u, 2) },
		func() (float64, float64, error) { return RayleighMeanVariance(2) }, 0, 30})
	_, err := RayleighPDF(nan, 1)
	requireErr(t, "x NaN", err, ErrNotFinite)
	_, err = RayleighCDF(1, 0)
	requireErr(t, "scale 0", err, ErrOutOfRange)
	_, err = RayleighPDF(1e-310, 1e-310)
	requireErr(t, "pdf overflow", err, ErrOverflow)
	_, err = RayleighQuantile(1, 1)
	requireErr(t, "u = 1", err, ErrOutOfRange)
	_, err = RayleighQuantile(0.5, 1.7e308)
	requireErr(t, "quantile overflow", err, ErrOverflow)
	_, _, err = RayleighMeanVariance(1e200)
	requireErr(t, "mv overflow", err, ErrOverflow)
}

// JKB vol. 2 ch. 26: a = 0, c = 1, b = 2: f(½) = ½, F(½) = ⅛, Q(⅛) = ½, (1, 1/6).
func TestTriangular(t *testing.T) {
	if d, err := TriangularPDF(0.5, 0, 1, 2); err != nil || d != 0.5 {
		t.Errorf("PDF(½) = (%v, %v)", d, err)
	}
	if d, _ := TriangularPDF(1, 0, 1, 2); d != 1 {
		t.Errorf("PDF(mode) = %v", d)
	}
	if c, _ := TriangularCDF(0.5, 0, 1, 2); c != 0.125 {
		t.Errorf("CDF(½) = %v", c)
	}
	if c, _ := TriangularCDF(1.5, 0, 1, 2); c != 0.875 {
		t.Errorf("CDF(1.5) = %v", c)
	}
	if q, _ := TriangularQuantile(0.125, 0, 1, 2); !within(q, 0.5, 1e-15) {
		t.Errorf("Q(⅛) = %v", q)
	}
	if q, _ := TriangularQuantile(1, 0, 1, 2); q != 2 {
		t.Errorf("Q(1) = %v", q)
	}
	if m, v, _ := TriangularMeanVariance(0, 1, 2); !within(m, 1, 1e-15) || !within(v, 1.0/6, 1e-15) {
		t.Errorf("MV = (%v, %v)", m, v)
	}
	for _, c := range [][3]float64{{1, 1, 4}, {1, 4, 4}, {-2, 0.5, 3}} {
		a, m, b := c[0], c[1], c[2]
		checkClosedForm(t, closedForm{"triangular",
			func(x float64) (float64, error) { return TriangularPDF(x, a, m, b) },
			func(x float64) (float64, error) { return TriangularCDF(x, a, m, b) },
			func(u float64) (float64, error) { return TriangularQuantile(u, a, m, b) },
			func() (float64, float64, error) { return TriangularMeanVariance(a, m, b) }, a, b})
	}
	_, err := TriangularPDF(nan, 0, 1, 2)
	requireErr(t, "x NaN", err, ErrNotFinite)
	_, err = TriangularCDF(0, 0, 3, 2)
	requireErr(t, "mode > high", err, ErrOutOfRange)
	_, err = TriangularCDF(0, 2, 2, 2)
	requireErr(t, "low = high", err, ErrOutOfRange)
	_, err = TriangularQuantile(0.5, -1.7e308, 0, 1.7e308)
	requireErr(t, "width overflow", err, ErrOverflow)
	_, err = TriangularPDF(0, 0, 0, 5e-324)
	requireErr(t, "pdf overflow", err, ErrOverflow)
	_, err = TriangularQuantile(1.5, 0, 1, 2)
	requireErr(t, "u > 1", err, ErrProbability)
	_, _, err = TriangularMeanVariance(0, 1e154, 1.5e154)
	if err != nil {
		t.Errorf("large finite mv: %v", err)
	}
	_, _, err = TriangularMeanVariance(0, 1e300, 1.5e300)
	requireErr(t, "mv overflow", err, ErrOverflow)
}

// JKB vol. 2 ch. 25 (beta(a, 1)): a = 3: f(½) = ¾, F(½) = ⅛, Q(⅛) = ½; a = 2 → (2/3, 1/18).
func TestPowerFunction(t *testing.T) {
	if d, err := PowerFunctionPDF(0.5, 3); err != nil || !within(d, 0.75, 1e-15) {
		t.Errorf("PDF(½) = (%v, %v)", d, err)
	}
	if d, _ := PowerFunctionPDF(0, 1); d != 1 {
		t.Errorf("PDF(0; 1) = %v", d)
	}
	if c, _ := PowerFunctionCDF(0.5, 3); !within(c, 0.125, 1e-15) {
		t.Errorf("CDF(½) = %v", c)
	}
	if q, _ := PowerFunctionQuantile(0.125, 3); !within(q, 0.5, 1e-15) {
		t.Errorf("Q(⅛) = %v", q)
	}
	if m, v, _ := PowerFunctionMeanVariance(2); !within(m, 2.0/3, 1e-15) || !within(v, 1.0/18, 1e-15) {
		t.Errorf("MV = (%v, %v)", m, v)
	}
	b, _ := BetaPDF(0.3, 2.5, 1)
	if d, _ := PowerFunctionPDF(0.3, 2.5); !within(d, b, 1e-14) {
		t.Errorf("vs BetaPDF: %v, %v", d, b)
	}
	checkClosedForm(t, closedForm{"power",
		func(x float64) (float64, error) { return PowerFunctionPDF(x, 2.5) },
		func(x float64) (float64, error) { return PowerFunctionCDF(x, 2.5) },
		func(u float64) (float64, error) { return PowerFunctionQuantile(u, 2.5) },
		func() (float64, float64, error) { return PowerFunctionMeanVariance(2.5) }, 0, 1})
	_, err := PowerFunctionPDF(0, 0.5)
	requireErr(t, "x = 0, a < 1", err, ErrOutOfRange)
	_, err = PowerFunctionCDF(0.5, 0)
	requireErr(t, "a = 0", err, ErrOutOfRange)
	_, err = PowerFunctionPDF(nan, 2)
	requireErr(t, "x NaN", err, ErrNotFinite)
	_, err = PowerFunctionPDF(5e-324, 0.001)
	requireErr(t, "pdf overflow", err, ErrOverflow)
	_, err = PowerFunctionQuantile(2, 2)
	requireErr(t, "u > 1", err, ErrProbability)
	_, _, err = PowerFunctionMeanVariance(inf)
	requireErr(t, "a Inf", err, ErrNotFinite)
}

func TestClosedForm_noAlloc(t *testing.T) {
	requireNoAlloc(t, "GumbelPDF", func() { sinkFloat, _ = GumbelPDF(1, 0, 1) })
	requireNoAlloc(t, "TriangularQuantile", func() { sinkFloat, _ = TriangularQuantile(0.3, 0, 1, 2) })
	requireNoAlloc(t, "ParetoCDF", func() { sinkFloat, _ = ParetoCDF(2, 1, 2) })
}
