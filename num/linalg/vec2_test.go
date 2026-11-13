package linalg

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// Sinks keep benchmark results alive.
var (
	sinkF    float64
	sinkV2   Vec2
	sinkV3   Vec3
	sinkV4   Vec4
	sinkM2   Mat2
	sinkM3   Mat3
	sinkM4   Mat4
	sinkQ    Quat
	sinkS    []float64
	sinkI    int
	sinkErr  error
	sinkCyl  Cylindrical
	sinkSph  Spherical
	sinkLU   LU
	sinkEig3 [3]float64
)

const ulp1 = 0x1p-52 // the spacing of float64 just above 1

// near reports |got − want| ≤ tol·max(1, |want|).
func near(got, want, tol float64) bool {
	return math.Abs(got-want) <= tol*max(1, math.Abs(want))
}

func nearV3(a, b Vec3, tol float64) bool {
	return near(a[0], b[0], tol) && near(a[1], b[1], tol) && near(a[2], b[2], tol)
}

func newRNG() *rand.Rand { return rand.New(rand.NewPCG(1, 2)) }

// sym returns a uniform value in [−1, 1).
func sym(r *rand.Rand) float64 { return 2*r.Float64() - 1 }

func bigF(x float64) *big.Float { return new(big.Float).SetPrec(2000).SetFloat64(x) }

func TestVec2_workedExamples(t *testing.T) {
	// Hand-checked examples from the doc comments (definitions in Strang 2016 §1.1, §1.2, §5.3).
	if r, err := Vec2Add(Vec2{1, 2}, Vec2{3, -5}); err != nil || r != (Vec2{4, -3}) {
		t.Errorf("Vec2Add = %v, %v", r, err)
	}
	if r, err := Vec2Sub(Vec2{1, 2}, Vec2{3, -5}); err != nil || r != (Vec2{-2, 7}) {
		t.Errorf("Vec2Sub = %v, %v", r, err)
	}
	if r, err := Vec2Scale(Vec2{1, -2}, 3); err != nil || r != (Vec2{3, -6}) {
		t.Errorf("Vec2Scale = %v, %v", r, err)
	}
	if r, err := Vec2Div(Vec2{3, -6}, 3); err != nil || r != (Vec2{1, -2}) {
		t.Errorf("Vec2Div = %v, %v", r, err)
	}
	if r, err := Vec2Dot(Vec2{1, 2}, Vec2{3, 4}); err != nil || r != 11 {
		t.Errorf("Vec2Dot = %v, %v", r, err)
	}
	if r, err := Vec2Cross(Vec2{1, 0}, Vec2{0, 1}); err != nil || r != 1 {
		t.Errorf("Vec2Cross(x̂, ŷ) = %v, %v", r, err)
	}
	if r, err := Vec2Cross(Vec2{3, 1}, Vec2{2, 4}); err != nil || r != 10 {
		t.Errorf("Vec2Cross = %v, %v", r, err)
	}
	if r, err := Vec2Cross(Vec2{0, 1}, Vec2{1, 0}); err != nil || r != -1 {
		t.Errorf("Vec2Cross clockwise = %v, %v, want −1", r, err)
	}
	if r, err := Vec2Norm(Vec2{3, 4}); err != nil || r != 5 {
		t.Errorf("Vec2Norm = %v, %v", r, err)
	}
	if r, err := Vec2Normalize(Vec2{3, 4}); err != nil || !near(r[0], 0.6, ulp1) || !near(r[1], 0.8, ulp1) {
		t.Errorf("Vec2Normalize = %v, %v", r, err)
	}
	if r, err := Vec2AngleBetween(Vec2{1, 0}, Vec2{1, 1}); err != nil || !near(r, math.Pi/4, 2*ulp1) {
		t.Errorf("Vec2AngleBetween = %v, %v", r, err)
	}
	if r, err := Vec2Rotate(Vec2{1, 0}, math.Pi/2); err != nil || !near(r[0], 0, 1e-16) || r[1] != 1 {
		t.Errorf("Vec2Rotate = %v, %v", r, err)
	}
}

func TestVec2_errors(t *testing.T) {
	nan, inf, big := math.NaN(), math.Inf(1), math.MaxFloat64
	cases := []struct {
		name string
		err  error
		want error
	}{
		{"Add NaN", second(Vec2Add(Vec2{nan, 0}, Vec2{})), ErrNotFinite},
		{"Add overflow", second(Vec2Add(Vec2{big, 0}, Vec2{big, 0})), ErrOverflow},
		{"Sub Inf", second(Vec2Sub(Vec2{}, Vec2{0, inf})), ErrNotFinite},
		{"Sub overflow", second(Vec2Sub(Vec2{big, 0}, Vec2{-big, 0})), ErrOverflow},
		{"Scale NaN k", second(Vec2Scale(Vec2{1, 1}, nan)), ErrNotFinite},
		{"Scale overflow", second(Vec2Scale(Vec2{big, 0}, 2)), ErrOverflow},
		{"Div by zero", second(Vec2Div(Vec2{1, 1}, 0)), ErrZeroDivisor},
		{"Div overflow", second(Vec2Div(Vec2{big, 0}, 0.5)), ErrOverflow},
		{"Dot overflow", second(Vec2Dot(Vec2{big, 0}, Vec2{2, 0})), ErrOverflow},
		{"Cross Inf", second(Vec2Cross(Vec2{inf, 0}, Vec2{})), ErrNotFinite},
		{"Norm overflow", second(Vec2Norm(Vec2{big, big})), ErrOverflow},
		{"Normalize zero", second(Vec2Normalize(Vec2{})), ErrZeroVector},
		{"Angle zero", second(Vec2AngleBetween(Vec2{1, 0}, Vec2{})), ErrZeroVector},
		{"Rotate NaN angle", second(Vec2Rotate(Vec2{1, 0}, nan)), ErrNotFinite},
	}
	for _, c := range cases {
		if !errors.Is(c.err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, c.err, c.want)
		}
	}
}

// second returns the error of a two-result call.
func second[T any](_ T, err error) error { return err }

// Oracle: Vec2Cross and Vec2Dot are within 2 ulps of the exact value (Jeannerod, Louvet, Muller 2013),
// checked against math/big on inputs chosen to cancel.
func TestVec2Cross_oracleKahan(t *testing.T) {
	r := newRNG()
	for trial := range 2000 {
		a := Vec2{sym(r), sym(r)}
		k := 1 + sym(r)*1e-9
		b := Vec2{a[0] * k, a[1]*k + sym(r)*1e-12} // nearly parallel: heavy cancellation
		got, _ := Vec2Cross(a, b)
		exact := new(big.Float).Sub(new(big.Float).Mul(bigF(a[0]), bigF(b[1])), new(big.Float).Mul(bigF(a[1]), bigF(b[0])))
		want, _ := exact.Float64()
		if math.Abs(got-want) > 2*math.Abs(want)*0x1p-53 {
			t.Fatalf("trial %d: Vec2Cross(%v, %v) = %v, exact %v", trial, a, b, got, want)
		}
		gotDot, _ := Vec2Dot(a, Vec2{b[0], -b[1]})
		exact = new(big.Float).Sub(new(big.Float).Mul(bigF(a[0]), bigF(b[0])), new(big.Float).Mul(bigF(a[1]), bigF(b[1])))
		want, _ = exact.Float64()
		if math.Abs(gotDot-want) > 2*math.Abs(want)*0x1p-53 {
			t.Fatalf("trial %d: Vec2Dot = %v, exact %v", trial, gotDot, want)
		}
	}
}

// Kahan's formula keeps nearly parallel angles: arccos of the normalised dot would return 0 or ~1.5e−8.
func TestVec2AngleBetween_nearlyParallel(t *testing.T) {
	got, err := Vec2AngleBetween(Vec2{1, 0}, Vec2{1, 1e-10})
	want := math.Atan(1e-10)
	if err != nil || math.Abs(got-want) > 1e-25 {
		t.Errorf("angle = %v, %v, want %v", got, err, want)
	}
	got, _ = Vec2AngleBetween(Vec2{1, 0}, Vec2{-1, 1e-10})
	if math.Abs(got-(math.Pi-want)) > 4e-16 {
		t.Errorf("antiparallel angle = %v, want %v", got, math.Pi-want)
	}
}

// Oracle: rotation preserves length and rotating by −θ undoes it.
func TestVec2Rotate_oracle(t *testing.T) {
	r := newRNG()
	for range 500 {
		v := Vec2{sym(r) * 100, sym(r) * 100}
		th := sym(r) * 10
		w, _ := Vec2Rotate(v, th)
		back, _ := Vec2Rotate(w, -th)
		n0, _ := Vec2Norm(v)
		n1, _ := Vec2Norm(w)
		if !near(n1, n0, 4*ulp1) || !near(back[0], v[0], 1e-13) || !near(back[1], v[1], 1e-13) {
			t.Fatalf("rotate %v by %v: %v, back %v", v, th, w, back)
		}
	}
}

func TestVec2Norm_noOverflow(t *testing.T) {
	if n, err := Vec2Norm(Vec2{3e300, 4e300}); err != nil || !near(n, 5e300, ulp1) {
		t.Errorf("Vec2Norm big = %v, %v", n, err)
	}
	if u, err := Vec2Normalize(Vec2{3e-320, 4e-320}); err != nil || !near(u[0], 0.6, 1e-3) {
		t.Errorf("Vec2Normalize subnormal = %v, %v", u, err) // subnormal inputs carry ~3 digits
	}
}

func TestVec2_noAllocation(t *testing.T) {
	a, b := Vec2{1, 2}, Vec2{3, 4}
	if n := testing.AllocsPerRun(100, func() {
		sinkV2, _ = Vec2Add(a, b)
		sinkV2, _ = Vec2Sub(a, b)
		sinkV2, _ = Vec2Scale(a, 2)
		sinkV2, _ = Vec2Div(a, 2)
		sinkF, _ = Vec2Dot(a, b)
		sinkF, _ = Vec2Cross(a, b)
		sinkF, _ = Vec2Norm(a)
		sinkV2, _ = Vec2Normalize(a)
		sinkF, _ = Vec2AngleBetween(a, b)
		sinkV2, _ = Vec2Rotate(a, 1)
	}); n != 0 {
		t.Errorf("allocations = %v, want 0", n)
	}
}

func BenchmarkVec2Add(b *testing.B) {
	b.Run("scalar", func(b *testing.B) {
		a2, b2 := Vec2{3, -4}, Vec2{1.5, 2.5}
		a3, b3, c3 := Vec3{1, -2, 3}, Vec3{0.5, 4, -1}, Vec3{2, 1, 1}
		m2a, m2b, m2s := Mat2{4, 7, 2, 6}, Mat2{1, 2, 3, 4}, Mat2{2, 1, 1, 3}
		m3a, m3b, m3s := Mat3{1, 2, 3, 0, 1, 4, 5, 6, 0}, Mat3{2, 0, 1, 1, 3, 2, 1, 1, 2}, Mat3{2, 1, 0, 1, 3, 1, 0, 1, 4}
		m4a := Mat4{2, 0, 0, 1, 0, 3, 1, 2, 0, 1, 4, 3, 0, 0, 0, 1}
		qa, qb := Quat{0.9, 0.1, -0.3, 0.2}, Quat{0.5, 0.5, 0.5, 0.5}
		m3r, _ := QuatToMat3(qa)
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = a2, b2, a3, b3, c3, m2a, m2b, m2s, m3a, m3b, m3s, m4a, qa, qb, m3r
		b.ReportAllocs()
		for b.Loop() {
			sinkV2, _ = Vec2Add(a2, b2)
		}
	})
}

func BenchmarkVec2Sub(b *testing.B) {
	b.Run("scalar", func(b *testing.B) {
		a2, b2 := Vec2{3, -4}, Vec2{1.5, 2.5}
		a3, b3, c3 := Vec3{1, -2, 3}, Vec3{0.5, 4, -1}, Vec3{2, 1, 1}
		m2a, m2b, m2s := Mat2{4, 7, 2, 6}, Mat2{1, 2, 3, 4}, Mat2{2, 1, 1, 3}
		m3a, m3b, m3s := Mat3{1, 2, 3, 0, 1, 4, 5, 6, 0}, Mat3{2, 0, 1, 1, 3, 2, 1, 1, 2}, Mat3{2, 1, 0, 1, 3, 1, 0, 1, 4}
		m4a := Mat4{2, 0, 0, 1, 0, 3, 1, 2, 0, 1, 4, 3, 0, 0, 0, 1}
		qa, qb := Quat{0.9, 0.1, -0.3, 0.2}, Quat{0.5, 0.5, 0.5, 0.5}
		m3r, _ := QuatToMat3(qa)
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = a2, b2, a3, b3, c3, m2a, m2b, m2s, m3a, m3b, m3s, m4a, qa, qb, m3r
		b.ReportAllocs()
		for b.Loop() {
			sinkV2, _ = Vec2Sub(a2, b2)
		}
	})
}

func BenchmarkVec2Scale(b *testing.B) {
	b.Run("scalar", func(b *testing.B) {
		a2, b2 := Vec2{3, -4}, Vec2{1.5, 2.5}
		a3, b3, c3 := Vec3{1, -2, 3}, Vec3{0.5, 4, -1}, Vec3{2, 1, 1}
		m2a, m2b, m2s := Mat2{4, 7, 2, 6}, Mat2{1, 2, 3, 4}, Mat2{2, 1, 1, 3}
		m3a, m3b, m3s := Mat3{1, 2, 3, 0, 1, 4, 5, 6, 0}, Mat3{2, 0, 1, 1, 3, 2, 1, 1, 2}, Mat3{2, 1, 0, 1, 3, 1, 0, 1, 4}
		m4a := Mat4{2, 0, 0, 1, 0, 3, 1, 2, 0, 1, 4, 3, 0, 0, 0, 1}
		qa, qb := Quat{0.9, 0.1, -0.3, 0.2}, Quat{0.5, 0.5, 0.5, 0.5}
		m3r, _ := QuatToMat3(qa)
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = a2, b2, a3, b3, c3, m2a, m2b, m2s, m3a, m3b, m3s, m4a, qa, qb, m3r
		b.ReportAllocs()
		for b.Loop() {
			sinkV2, _ = Vec2Scale(a2, 1.5)
		}
	})
}

func BenchmarkVec2Div(b *testing.B) {
	b.Run("scalar", func(b *testing.B) {
		a2, b2 := Vec2{3, -4}, Vec2{1.5, 2.5}
		a3, b3, c3 := Vec3{1, -2, 3}, Vec3{0.5, 4, -1}, Vec3{2, 1, 1}
		m2a, m2b, m2s := Mat2{4, 7, 2, 6}, Mat2{1, 2, 3, 4}, Mat2{2, 1, 1, 3}
		m3a, m3b, m3s := Mat3{1, 2, 3, 0, 1, 4, 5, 6, 0}, Mat3{2, 0, 1, 1, 3, 2, 1, 1, 2}, Mat3{2, 1, 0, 1, 3, 1, 0, 1, 4}
		m4a := Mat4{2, 0, 0, 1, 0, 3, 1, 2, 0, 1, 4, 3, 0, 0, 0, 1}
		qa, qb := Quat{0.9, 0.1, -0.3, 0.2}, Quat{0.5, 0.5, 0.5, 0.5}
		m3r, _ := QuatToMat3(qa)
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = a2, b2, a3, b3, c3, m2a, m2b, m2s, m3a, m3b, m3s, m4a, qa, qb, m3r
		b.ReportAllocs()
		for b.Loop() {
			sinkV2, _ = Vec2Div(a2, 1.5)
		}
	})
}

func BenchmarkVec2Dot(b *testing.B) {
	b.Run("scalar", func(b *testing.B) {
		a2, b2 := Vec2{3, -4}, Vec2{1.5, 2.5}
		a3, b3, c3 := Vec3{1, -2, 3}, Vec3{0.5, 4, -1}, Vec3{2, 1, 1}
		m2a, m2b, m2s := Mat2{4, 7, 2, 6}, Mat2{1, 2, 3, 4}, Mat2{2, 1, 1, 3}
		m3a, m3b, m3s := Mat3{1, 2, 3, 0, 1, 4, 5, 6, 0}, Mat3{2, 0, 1, 1, 3, 2, 1, 1, 2}, Mat3{2, 1, 0, 1, 3, 1, 0, 1, 4}
		m4a := Mat4{2, 0, 0, 1, 0, 3, 1, 2, 0, 1, 4, 3, 0, 0, 0, 1}
		qa, qb := Quat{0.9, 0.1, -0.3, 0.2}, Quat{0.5, 0.5, 0.5, 0.5}
		m3r, _ := QuatToMat3(qa)
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = a2, b2, a3, b3, c3, m2a, m2b, m2s, m3a, m3b, m3s, m4a, qa, qb, m3r
		b.ReportAllocs()
		for b.Loop() {
			sinkF, _ = Vec2Dot(a2, b2)
		}
	})
}

func BenchmarkVec2Cross(b *testing.B) {
	b.Run("scalar", func(b *testing.B) {
		a2, b2 := Vec2{3, -4}, Vec2{1.5, 2.5}
		a3, b3, c3 := Vec3{1, -2, 3}, Vec3{0.5, 4, -1}, Vec3{2, 1, 1}
		m2a, m2b, m2s := Mat2{4, 7, 2, 6}, Mat2{1, 2, 3, 4}, Mat2{2, 1, 1, 3}
		m3a, m3b, m3s := Mat3{1, 2, 3, 0, 1, 4, 5, 6, 0}, Mat3{2, 0, 1, 1, 3, 2, 1, 1, 2}, Mat3{2, 1, 0, 1, 3, 1, 0, 1, 4}
		m4a := Mat4{2, 0, 0, 1, 0, 3, 1, 2, 0, 1, 4, 3, 0, 0, 0, 1}
		qa, qb := Quat{0.9, 0.1, -0.3, 0.2}, Quat{0.5, 0.5, 0.5, 0.5}
		m3r, _ := QuatToMat3(qa)
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = a2, b2, a3, b3, c3, m2a, m2b, m2s, m3a, m3b, m3s, m4a, qa, qb, m3r
		b.ReportAllocs()
		for b.Loop() {
			sinkF, _ = Vec2Cross(a2, b2)
		}
	})
}

func BenchmarkVec2Norm(b *testing.B) {
	b.Run("scalar", func(b *testing.B) {
		a2, b2 := Vec2{3, -4}, Vec2{1.5, 2.5}
		a3, b3, c3 := Vec3{1, -2, 3}, Vec3{0.5, 4, -1}, Vec3{2, 1, 1}
		m2a, m2b, m2s := Mat2{4, 7, 2, 6}, Mat2{1, 2, 3, 4}, Mat2{2, 1, 1, 3}
		m3a, m3b, m3s := Mat3{1, 2, 3, 0, 1, 4, 5, 6, 0}, Mat3{2, 0, 1, 1, 3, 2, 1, 1, 2}, Mat3{2, 1, 0, 1, 3, 1, 0, 1, 4}
		m4a := Mat4{2, 0, 0, 1, 0, 3, 1, 2, 0, 1, 4, 3, 0, 0, 0, 1}
		qa, qb := Quat{0.9, 0.1, -0.3, 0.2}, Quat{0.5, 0.5, 0.5, 0.5}
		m3r, _ := QuatToMat3(qa)
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = a2, b2, a3, b3, c3, m2a, m2b, m2s, m3a, m3b, m3s, m4a, qa, qb, m3r
		b.ReportAllocs()
		for b.Loop() {
			sinkF, _ = Vec2Norm(a2)
		}
	})
}

func BenchmarkVec2Normalize(b *testing.B) {
	b.Run("scalar", func(b *testing.B) {
		a2, b2 := Vec2{3, -4}, Vec2{1.5, 2.5}
		a3, b3, c3 := Vec3{1, -2, 3}, Vec3{0.5, 4, -1}, Vec3{2, 1, 1}
		m2a, m2b, m2s := Mat2{4, 7, 2, 6}, Mat2{1, 2, 3, 4}, Mat2{2, 1, 1, 3}
		m3a, m3b, m3s := Mat3{1, 2, 3, 0, 1, 4, 5, 6, 0}, Mat3{2, 0, 1, 1, 3, 2, 1, 1, 2}, Mat3{2, 1, 0, 1, 3, 1, 0, 1, 4}
		m4a := Mat4{2, 0, 0, 1, 0, 3, 1, 2, 0, 1, 4, 3, 0, 0, 0, 1}
		qa, qb := Quat{0.9, 0.1, -0.3, 0.2}, Quat{0.5, 0.5, 0.5, 0.5}
		m3r, _ := QuatToMat3(qa)
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = a2, b2, a3, b3, c3, m2a, m2b, m2s, m3a, m3b, m3s, m4a, qa, qb, m3r
		b.ReportAllocs()
		for b.Loop() {
			sinkV2, _ = Vec2Normalize(a2)
		}
	})
}

func BenchmarkVec2AngleBetween(b *testing.B) {
	b.Run("scalar", func(b *testing.B) {
		a2, b2 := Vec2{3, -4}, Vec2{1.5, 2.5}
		a3, b3, c3 := Vec3{1, -2, 3}, Vec3{0.5, 4, -1}, Vec3{2, 1, 1}
		m2a, m2b, m2s := Mat2{4, 7, 2, 6}, Mat2{1, 2, 3, 4}, Mat2{2, 1, 1, 3}
		m3a, m3b, m3s := Mat3{1, 2, 3, 0, 1, 4, 5, 6, 0}, Mat3{2, 0, 1, 1, 3, 2, 1, 1, 2}, Mat3{2, 1, 0, 1, 3, 1, 0, 1, 4}
		m4a := Mat4{2, 0, 0, 1, 0, 3, 1, 2, 0, 1, 4, 3, 0, 0, 0, 1}
		qa, qb := Quat{0.9, 0.1, -0.3, 0.2}, Quat{0.5, 0.5, 0.5, 0.5}
		m3r, _ := QuatToMat3(qa)
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = a2, b2, a3, b3, c3, m2a, m2b, m2s, m3a, m3b, m3s, m4a, qa, qb, m3r
		b.ReportAllocs()
		for b.Loop() {
			sinkF, _ = Vec2AngleBetween(a2, b2)
		}
	})
}

func BenchmarkVec2Rotate(b *testing.B) {
	b.Run("scalar", func(b *testing.B) {
		a2, b2 := Vec2{3, -4}, Vec2{1.5, 2.5}
		a3, b3, c3 := Vec3{1, -2, 3}, Vec3{0.5, 4, -1}, Vec3{2, 1, 1}
		m2a, m2b, m2s := Mat2{4, 7, 2, 6}, Mat2{1, 2, 3, 4}, Mat2{2, 1, 1, 3}
		m3a, m3b, m3s := Mat3{1, 2, 3, 0, 1, 4, 5, 6, 0}, Mat3{2, 0, 1, 1, 3, 2, 1, 1, 2}, Mat3{2, 1, 0, 1, 3, 1, 0, 1, 4}
		m4a := Mat4{2, 0, 0, 1, 0, 3, 1, 2, 0, 1, 4, 3, 0, 0, 0, 1}
		qa, qb := Quat{0.9, 0.1, -0.3, 0.2}, Quat{0.5, 0.5, 0.5, 0.5}
		m3r, _ := QuatToMat3(qa)
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = a2, b2, a3, b3, c3, m2a, m2b, m2s, m3a, m3b, m3s, m4a, qa, qb, m3r
		b.ReportAllocs()
		for b.Loop() {
			sinkV2, _ = Vec2Rotate(a2, 0.7)
		}
	})
}
