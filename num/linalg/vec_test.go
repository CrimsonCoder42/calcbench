package linalg

import (
	"errors"
	"math"
	"math/big"
	"strconv"
	"testing"
)

func eqSlice(a, b []float64, tol float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !near(a[i], b[i], tol) {
			return false
		}
	}
	return true
}

func TestVecN_workedExamples(t *testing.T) {
	// Hand-checked examples from the doc comments.
	type c struct {
		name string
		got  []float64
		err  error
		want []float64
		tol  float64
	}
	g := func(s []float64, err error) ([]float64, error) { return s, err }
	var cs []c
	add := func(name string, want []float64, tol float64, s []float64, err error) {
		cs = append(cs, c{name, s, err, want, tol})
	}
	s, err := g(VecAdd(nil, []float64{1, 2, 3, 4}, []float64{4, 3, 2, 1}))
	add("VecAdd", []float64{5, 5, 5, 5}, 0, s, err)
	s, err = g(VecSub(nil, []float64{1, 2, 3, 4}, []float64{4, 3, 2, 1}))
	add("VecSub", []float64{-3, -1, 1, 3}, 0, s, err)
	s, err = g(VecScale(nil, []float64{1, -2, 3, 0}, 2))
	add("VecScale", []float64{2, -4, 6, 0}, 0, s, err)
	s, err = g(VecDiv(nil, []float64{2, -4, 6, 0}, 2))
	add("VecDiv", []float64{1, -2, 3, 0}, 0, s, err)
	s, err = g(VecNormalize(nil, []float64{1, 2, 2, 4}))
	add("VecNormalize", []float64{0.2, 0.4, 0.4, 0.8}, ulp1, s, err)
	s, err = g(VecProject(nil, []float64{1, 2, 3, 4}, []float64{0, 0, 0, 2}))
	add("VecProject", []float64{0, 0, 0, 4}, 0, s, err)
	s, err = g(VecLerp(nil, []float64{0, 10}, []float64{4, 20}, 0.5))
	add("VecLerp", []float64{2, 15}, 0, s, err)
	s, err = g(VecMin(nil, []float64{1, 5, -2}, []float64{3, 0, -4}))
	add("VecMin", []float64{1, 0, -4}, 0, s, err)
	s, err = g(VecMax(nil, []float64{1, 5, -2}, []float64{3, 0, -4}))
	add("VecMax", []float64{3, 5, -2}, 0, s, err)
	for _, x := range cs {
		if x.err != nil || !eqSlice(x.got, x.want, x.tol) {
			t.Errorf("%s = %v, %v; want %v", x.name, x.got, x.err, x.want)
		}
	}
	scal := []struct {
		name      string
		got, want float64
		err       error
		tol       float64
	}{
		{"VecDot", first(VecDot([]float64{1, 2, 3, 4}, []float64{4, 3, 2, 1})), 20, nil, 0},
		{"VecDotCompensated", first(VecDotCompensated([]float64{1e16, 1, -1e16}, []float64{1, 1, 1})), 1, nil, 0},
		{"VecNorm", first(VecNorm([]float64{1, 2, 2, 4})), 5, nil, 0},
		{"VecNorm1", first(VecNorm1([]float64{1, -2, 3, -4})), 10, nil, 0},
		{"VecNormInf", first(VecNormInf([]float64{1, -7, 3})), 7, nil, 0},
		{"VecNormP p=2", first(VecNormP([]float64{3, 4}, 2)), 5, nil, ulp1},
		{"VecNormP p=3", first(VecNormP([]float64{1, 1, 1, 1}, 3)), 1.5874010519681994, nil, 2 * ulp1},
		{"VecDistance", first(VecDistance([]float64{1, 1, 1, 1}, []float64{2, 3, 3, 5})), 5, nil, 0},
		{"VecAngleBetween", first(VecAngleBetween([]float64{1, 0, 0, 0}, []float64{1, 1, 1, 1})), math.Pi / 3, nil, 2 * ulp1},
	}
	for _, x := range scal {
		if !near(x.got, x.want, x.tol) {
			t.Errorf("%s = %v, want %v", x.name, x.got, x.want)
		}
	}
	// The plain FMA dot product loses the 1 entirely (the reason VecDotCompensated exists).
	if d, _ := VecDot([]float64{1e16, 1, -1e16}, []float64{1, 1, 1}); d != 0 {
		t.Errorf("VecDot cancellation example = %v, documented as 0", d)
	}
}

func first[T any](v T, _ error) T { return v }

func TestVecN_edgeCases(t *testing.T) {
	if s, err := VecAdd(nil, nil, nil); err != nil || len(s) != 0 {
		t.Errorf("empty add = %v, %v", s, err)
	}
	if d, err := VecDot(nil, nil); err != nil || d != 0 {
		t.Errorf("empty dot = %v, %v", d, err)
	}
	if n, err := VecNorm([]float64{0, 0}); err != nil || n != 0 {
		t.Errorf("zero norm = %v, %v", n, err)
	}
	if n, err := VecNorm([]float64{-3}); err != nil || n != 3 {
		t.Errorf("one element = %v, %v", n, err)
	}
	if n, err := VecNorm([]float64{3e300, 4e300}); err != nil || !near(n, 5e300, ulp1) {
		t.Errorf("no overflow = %v, %v", n, err)
	}
	if n, err := VecNorm([]float64{3e-310, 4e-310}); err != nil || !near(n/5e-310, 1, 1e-5) {
		t.Errorf("subnormal = %v, %v", n, err) // subnormal inputs carry ~5 digits
	}
	if d, err := VecDistance([]float64{1, 2}, []float64{1, 2}); err != nil || d != 0 {
		t.Errorf("zero distance = %v, %v", d, err)
	}
	if n, err := VecNormP([]float64{0, 0}, 3); err != nil || n != 0 {
		t.Errorf("zero p-norm = %v, %v", n, err)
	}
	// dst may be an input for the element-by-element functions.
	a := []float64{1, 2, 3}
	if s, err := VecScale(a, a, 2); err != nil || &s[0] != &a[0] || a[2] != 6 {
		t.Errorf("in place scale = %v, %v", s, err)
	}
}

func TestVecN_errors(t *testing.T) {
	nan, inf, big := math.NaN(), math.Inf(-1), math.MaxFloat64
	one, two := []float64{1}, []float64{1, 2}
	cases := []struct {
		name string
		err  error
		want error
	}{
		{"Add length", second(VecAdd(nil, one, two)), ErrDimension},
		{"Add NaN", second(VecAdd(nil, []float64{nan}, one)), ErrNotFinite},
		{"Add overflow", second(VecAdd(nil, []float64{big}, []float64{big})), ErrOverflow},
		{"Sub length", second(VecSub(nil, one, two)), ErrDimension},
		{"Sub Inf", second(VecSub(nil, one, []float64{inf})), ErrNotFinite},
		{"Scale NaN k", second(VecScale(nil, one, nan)), ErrNotFinite},
		{"Scale overflow", second(VecScale(nil, []float64{big}, 2)), ErrOverflow},
		{"Div zero", second(VecDiv(nil, one, 0)), ErrZeroDivisor},
		{"Div NaN element", second(VecDiv(nil, []float64{nan}, 2)), ErrNotFinite},
		{"Dot length", second(VecDot(one, two)), ErrDimension},
		{"Dot NaN", second(VecDot([]float64{nan}, one)), ErrNotFinite},
		{"Dot overflow", second(VecDot([]float64{big, big}, []float64{big, big})), ErrOverflow},
		{"DotC length", second(VecDotCompensated(one, two)), ErrDimension},
		{"DotC Inf", second(VecDotCompensated(one, []float64{inf})), ErrNotFinite},
		{"DotC overflow", second(VecDotCompensated([]float64{big}, []float64{2})), ErrOverflow},
		{"Norm NaN", second(VecNorm([]float64{1, nan})), ErrNotFinite},
		{"Norm overflow", second(VecNorm([]float64{big, big})), ErrOverflow},
		{"Norm1 overflow", second(VecNorm1([]float64{big, big})), ErrOverflow},
		{"Norm1 Inf", second(VecNorm1([]float64{inf})), ErrNotFinite},
		{"NormInf NaN", second(VecNormInf([]float64{nan})), ErrNotFinite},
		{"NormP p<1", second(VecNormP(one, 0.5)), ErrOutOfRange},
		{"NormP p Inf", second(VecNormP(one, math.Inf(1))), ErrNotFinite},
		{"NormP NaN", second(VecNormP([]float64{nan}, 2)), ErrNotFinite},
		{"NormP overflow", second(VecNormP([]float64{big, big}, 1)), ErrOverflow},
		{"Normalize zero", second(VecNormalize(nil, []float64{0, 0})), ErrZeroVector},
		{"Normalize empty", second(VecNormalize(nil, nil)), ErrZeroVector},
		{"Normalize NaN", second(VecNormalize(nil, []float64{nan})), ErrNotFinite},
		{"Distance length", second(VecDistance(one, two)), ErrDimension},
		{"Distance overflow", second(VecDistance([]float64{big}, []float64{-big})), ErrOverflow},
		{"Distance NaN", second(VecDistance([]float64{nan}, one)), ErrNotFinite},
		{"Angle zero", second(VecAngleBetween([]float64{0}, one)), ErrZeroVector},
		{"Angle length", second(VecAngleBetween(one, two)), ErrDimension},
		{"Angle NaN", second(VecAngleBetween(one, []float64{nan})), ErrNotFinite},
		{"Project zero", second(VecProject(nil, one, []float64{0})), ErrZeroVector},
		{"Project length", second(VecProject(nil, one, two)), ErrDimension},
		{"Project NaN", second(VecProject(nil, []float64{nan}, one)), ErrNotFinite},
		{"Project overflow", second(VecProject(nil, []float64{big, big}, []float64{1, 1})), ErrOverflow},
		{"Lerp length", second(VecLerp(nil, one, two, 0.5)), ErrDimension},
		{"Lerp NaN t", second(VecLerp(nil, one, one, nan)), ErrNotFinite},
		{"Lerp overflow", second(VecLerp(nil, []float64{0}, []float64{big}, 3)), ErrOverflow},
		{"Min length", second(VecMin(nil, one, two)), ErrDimension},
		{"Min NaN", second(VecMin(nil, []float64{nan}, one)), ErrNotFinite},
		{"Max length", second(VecMax(nil, one, two)), ErrDimension},
		{"Max NaN", second(VecMax(nil, one, []float64{nan})), ErrNotFinite},
	}
	for _, c := range cases {
		if !errors.Is(c.err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, c.err, c.want)
		}
	}
}

// exactDot returns Σ aᵢbᵢ rounded once from an exact math/big sum.
func exactDot(a, b []float64) float64 {
	s := new(big.Float).SetPrec(4000)
	for i := range a {
		s.Add(s, new(big.Float).SetPrec(4000).Mul(bigF(a[i]), bigF(b[i])))
	}
	f, _ := s.Float64()
	return f
}

// Oracle: VecDotCompensated meets Ogita–Rump–Oishi Proposition 5.5, |res − aᵀb| ≤ u|aᵀb| + γₙ²|a|ᵀ|b|, on
// ill-conditioned dot products built to cancel (fixed seed), where the plain sum is far off.
func TestVecDotCompensated_oracle(t *testing.T) {
	r := newRNG()
	const u = 0x1p-53
	for trial := range 300 {
		n := 2 + r.IntN(200)
		a, b := make([]float64, n), make([]float64, n)
		for i := range n {
			a[i] = sym(r) * math.Ldexp(1, r.IntN(60)-30)
			b[i] = sym(r) * math.Ldexp(1, r.IntN(60)-30)
		}
		// Make the sum cancel: append the negated plain partial sum as a last term.
		s, _ := VecDot(a, b)
		a = append(a, -s)
		b = append(b, 1)
		n++
		got, err := VecDotCompensated(a, b)
		want := exactDot(a, b)
		abs := 0.0
		for i := range a {
			abs += math.Abs(a[i] * b[i])
		}
		g := float64(n) * u / (1 - float64(n)*u)
		if err != nil || math.Abs(got-want) > u*math.Abs(want)+g*g*abs*1.0001 {
			t.Fatalf("trial %d (n = %d): got %v, exact %v", trial, n, got, want)
		}
	}
}

// Oracle: VecNorm is within one ulp (plus the square root's half ulp) of √(exact Σ vᵢ²) from math/big,
// and p = 1, 2 of VecNormP agree with VecNorm1 and VecNorm.
func TestVecNorm_oracle(t *testing.T) {
	r := newRNG()
	for trial := range 300 {
		n := 1 + r.IntN(500)
		v := make([]float64, n)
		for i := range v {
			v[i] = sym(r) * math.Ldexp(1, r.IntN(40)-20)
		}
		got, _ := VecNorm(v)
		ss := new(big.Float).SetPrec(4000)
		for _, x := range v {
			ss.Add(ss, new(big.Float).SetPrec(4000).Mul(bigF(x), bigF(x)))
		}
		want, _ := new(big.Float).SetPrec(4000).Sqrt(ss).Float64()
		if !near(got/want, 1, 2*ulp1) {
			t.Fatalf("trial %d: VecNorm = %v, exact %v", trial, got, want)
		}
		p2, _ := VecNormP(v, 2)
		p1, _ := VecNormP(v, 1)
		n1, _ := VecNorm1(v)
		if !near(p2/got, 1, 8*ulp1) || !near(p1/n1, 1, 8*ulp1) {
			t.Fatalf("trial %d: NormP disagrees: %v vs %v, %v vs %v", trial, p2, got, p1, n1)
		}
		u, _ := VecNormalize(nil, v)
		if nu, _ := VecNorm(u); !near(nu, 1, 4*ulp1) {
			t.Fatalf("trial %d: |normalize(v)| = %v", trial, nu)
		}
	}
}

// Oracle: the n-D angle equals the 3D angle for 3-element vectors, and the n-D distance equals the norm of
// the difference.
func TestVecAngleBetween_oracle(t *testing.T) {
	r := newRNG()
	for trial := range 300 {
		a, b := Vec3{sym(r), sym(r), sym(r)}, Vec3{sym(r), sym(r), sym(r)}
		got, _ := VecAngleBetween(a[:], b[:])
		want, _ := Vec3AngleBetween(a, b)
		if math.Abs(got-want) > 8e-16 {
			t.Fatalf("trial %d: %v vs %v", trial, got, want)
		}
		d, _ := VecDistance(a[:], b[:])
		d3, _ := Vec3Distance(a, b)
		if !near(d, d3, 4*ulp1) {
			t.Fatalf("trial %d: distance %v vs %v", trial, d, d3)
		}
		p, _ := VecProject(nil, a[:], b[:])
		p3, _ := Vec3Project(a, b)
		if !nearV3(Vec3(p), p3, 16*ulp1) {
			t.Fatalf("trial %d: project %v vs %v", trial, p, p3)
		}
	}
}

func TestVecN_dstReuseNoAllocation(t *testing.T) {
	a, b := make([]float64, 64), make([]float64, 64)
	for i := range a {
		a[i], b[i] = float64(i)+1, float64(64-i)
	}
	dst := make([]float64, 64)
	if n := testing.AllocsPerRun(50, func() {
		sinkS, _ = VecAdd(dst, a, b)
		sinkS, _ = VecSub(dst, a, b)
		sinkS, _ = VecScale(dst, a, 2)
		sinkS, _ = VecDiv(dst, a, 2)
		sinkS, _ = VecNormalize(dst, a)
		sinkS, _ = VecProject(dst, a, b)
		sinkS, _ = VecLerp(dst, a, b, 0.25)
		sinkS, _ = VecMin(dst, a, b)
		sinkS, _ = VecMax(dst, a, b)
		sinkF, _ = VecDot(a, b)
		sinkF, _ = VecDotCompensated(a, b)
		sinkF, _ = VecNorm(a)
		sinkF, _ = VecNorm1(a)
		sinkF, _ = VecNormInf(a)
		sinkF, _ = VecNormP(a, 3)
		sinkF, _ = VecDistance(a, b)
		sinkF, _ = VecAngleBetween(a, b)
	}); n != 0 {
		t.Errorf("allocations with a reused dst = %v, want 0", n)
	}
}

// benchVectors returns two length-n vectors with mixed magnitudes and a reusable dst.
func benchVectors(n int) (a, b, dst []float64) {
	r := newRNG()
	a, b, dst = make([]float64, n), make([]float64, n), make([]float64, n)
	for i := range a {
		a[i], b[i] = sym(r)*100, sym(r)*100
	}
	return a, b, dst
}

func BenchmarkVecAdd(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		a, bb, dst := benchVectors(n)
		_, _, _ = a, bb, dst
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkS, _ = VecAdd(dst, a, bb)
			}
		})
	}
}

func BenchmarkVecSub(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		a, bb, dst := benchVectors(n)
		_, _, _ = a, bb, dst
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkS, _ = VecSub(dst, a, bb)
			}
		})
	}
}

func BenchmarkVecScale(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		a, bb, dst := benchVectors(n)
		_, _, _ = a, bb, dst
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkS, _ = VecScale(dst, a, 1.5)
			}
		})
	}
}

func BenchmarkVecDiv(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		a, bb, dst := benchVectors(n)
		_, _, _ = a, bb, dst
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkS, _ = VecDiv(dst, a, 1.5)
			}
		})
	}
}

func BenchmarkVecDot(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		a, bb, dst := benchVectors(n)
		_, _, _ = a, bb, dst
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkF, _ = VecDot(a, bb)
			}
		})
	}
}

func BenchmarkVecDotCompensated(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		a, bb, dst := benchVectors(n)
		_, _, _ = a, bb, dst
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkF, _ = VecDotCompensated(a, bb)
			}
		})
	}
}

func BenchmarkVecNorm(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		a, bb, dst := benchVectors(n)
		_, _, _ = a, bb, dst
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkF, _ = VecNorm(a)
			}
		})
	}
}

func BenchmarkVecNorm1(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		a, bb, dst := benchVectors(n)
		_, _, _ = a, bb, dst
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkF, _ = VecNorm1(a)
			}
		})
	}
}

func BenchmarkVecNormInf(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		a, bb, dst := benchVectors(n)
		_, _, _ = a, bb, dst
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkF, _ = VecNormInf(a)
			}
		})
	}
}

func BenchmarkVecNormP(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		a, bb, dst := benchVectors(n)
		_, _, _ = a, bb, dst
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkF, _ = VecNormP(a, 3)
			}
		})
	}
}

func BenchmarkVecNormalize(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		a, bb, dst := benchVectors(n)
		_, _, _ = a, bb, dst
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkS, _ = VecNormalize(dst, a)
			}
		})
	}
}

func BenchmarkVecDistance(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		a, bb, dst := benchVectors(n)
		_, _, _ = a, bb, dst
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkF, _ = VecDistance(a, bb)
			}
		})
	}
}

func BenchmarkVecAngleBetween(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		a, bb, dst := benchVectors(n)
		_, _, _ = a, bb, dst
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkF, _ = VecAngleBetween(a, bb)
			}
		})
	}
}

func BenchmarkVecProject(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		a, bb, dst := benchVectors(n)
		_, _, _ = a, bb, dst
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkS, _ = VecProject(dst, a, bb)
			}
		})
	}
}

func BenchmarkVecLerp(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		a, bb, dst := benchVectors(n)
		_, _, _ = a, bb, dst
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkS, _ = VecLerp(dst, a, bb, 0.3)
			}
		})
	}
}

func BenchmarkVecMin(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		a, bb, dst := benchVectors(n)
		_, _, _ = a, bb, dst
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkS, _ = VecMin(dst, a, bb)
			}
		})
	}
}

func BenchmarkVecMax(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		a, bb, dst := benchVectors(n)
		_, _, _ = a, bb, dst
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkS, _ = VecMax(dst, a, bb)
			}
		})
	}
}

// Review regressions (findings 1 and 4): subnormal inputs must not produce NaN or a false overflow.
func TestVecAngleBetween_subnormal(t *testing.T) {
	if a, err := VecAngleBetween([]float64{1e-310, 2e-310}, []float64{1, 2}); err != nil || math.IsNaN(a) || a > 1e-5 {
		t.Errorf("parallel subnormal = %v, %v; want ≈ 0", a, err) // a subnormal carries ~5 digits
	}
	if a, err := VecAngleBetween([]float64{5e-324}, []float64{1}); err != nil || a != 0 {
		t.Errorf("smallest subnormal = %v, %v; want 0", a, err)
	}
}

func TestVecProject_subnormalOnto(t *testing.T) {
	if p, err := VecProject(nil, []float64{1, 2}, []float64{1e-310, 1e-310}); err != nil || !eqSlice(p, []float64{1.5, 1.5}, 1e-5) {
		t.Errorf("project onto a subnormal = %v, %v; want (1.5, 1.5)", p, err) // subnormal direction: ~5 digits
	}
}
