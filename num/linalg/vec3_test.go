package linalg

import (
	"errors"
	"math"
	"testing"
)

func TestVec3_workedExamples(t *testing.T) {
	// Hand-checked examples from the doc comments (Strang 2016 §1.1, §1.2, §4.2, §5.3, §10.6).
	eq := func(name string, got, want Vec3, err error, tol float64) {
		t.Helper()
		if err != nil || !nearV3(got, want, tol) {
			t.Errorf("%s = %v, %v; want %v", name, got, err, want)
		}
	}
	eqF := func(name string, got, want float64, err error, tol float64) {
		t.Helper()
		if err != nil || !near(got, want, tol) {
			t.Errorf("%s = %v, %v; want %v", name, got, err, want)
		}
	}
	r, err := Vec3Add(Vec3{1, 2, 3}, Vec3{4, -5, 6})
	eq("Vec3Add", r, Vec3{5, -3, 9}, err, 0)
	r, err = Vec3Sub(Vec3{1, 2, 3}, Vec3{4, -5, 6})
	eq("Vec3Sub", r, Vec3{-3, 7, -3}, err, 0)
	r, err = Vec3Scale(Vec3{1, -2, 3}, 2)
	eq("Vec3Scale", r, Vec3{2, -4, 6}, err, 0)
	r, err = Vec3Div(Vec3{2, -4, 6}, 2)
	eq("Vec3Div", r, Vec3{1, -2, 3}, err, 0)
	f, err := Vec3Dot(Vec3{1, 2, 3}, Vec3{4, -5, 6})
	eqF("Vec3Dot", f, 12, err, 0)
	r, err = Vec3Cross(Vec3{1, 0, 0}, Vec3{0, 1, 0})
	eq("Vec3Cross x̂×ŷ", r, Vec3{0, 0, 1}, err, 0)
	r, err = Vec3Cross(Vec3{1, 2, 3}, Vec3{4, 5, 6})
	eq("Vec3Cross", r, Vec3{-3, 6, -3}, err, 0)
	f, err = Vec3Norm(Vec3{2, 3, 6})
	eqF("Vec3Norm", f, 7, err, 0)
	r, err = Vec3Normalize(Vec3{2, 3, 6})
	eq("Vec3Normalize", r, Vec3{2.0 / 7, 3.0 / 7, 6.0 / 7}, err, 2*ulp1)
	f, err = Vec3Distance(Vec3{1, 2, 3}, Vec3{4, 6, 15})
	eqF("Vec3Distance", f, 13, err, 0)
	f, err = Vec3AngleBetween(Vec3{1, 0, 0}, Vec3{0, 1, 1})
	eqF("Vec3AngleBetween ⟂", f, math.Pi/2, err, 2*ulp1)
	f, err = Vec3AngleBetween(Vec3{1, 0, 0}, Vec3{1, 1, 0})
	eqF("Vec3AngleBetween 45°", f, math.Pi/4, err, 2*ulp1)
	r, err = Vec3DirectionAngles(Vec3{1, 1, 0})
	eq("Vec3DirectionAngles", r, Vec3{math.Pi / 4, math.Pi / 4, math.Pi / 2}, err, 2*ulp1)
	f, err = Vec3ComponentAlong(Vec3{3, 4, 5}, Vec3{0, 2, 0})
	eqF("Vec3ComponentAlong", f, 4, err, 0)
	r, err = Vec3Project(Vec3{3, 4, 5}, Vec3{0, 2, 0})
	eq("Vec3Project", r, Vec3{0, 4, 0}, err, 0)
	r, err = Vec3Reject(Vec3{3, 4, 5}, Vec3{0, 2, 0})
	eq("Vec3Reject", r, Vec3{3, 0, 5}, err, 0)
	r, err = Vec3Reflect(Vec3{1, -1, 2}, Vec3{0, 3, 0})
	eq("Vec3Reflect", r, Vec3{1, 1, 2}, err, 0)
	f, err = Vec3ParallelogramArea(Vec3{3, 0, 0}, Vec3{1, 2, 0})
	eqF("Vec3ParallelogramArea", f, 6, err, 0)
	f, err = Vec3TriangleArea(Vec3{3, 0, 0}, Vec3{1, 2, 0})
	eqF("Vec3TriangleArea", f, 3, err, 0)
	f, err = Vec3ScalarTripleProduct(Vec3{1, 0, 0}, Vec3{0, 1, 0}, Vec3{0, 0, 1})
	eqF("Vec3ScalarTripleProduct unit", f, 1, err, 0)
	f, err = Vec3ScalarTripleProduct(Vec3{1, 2, 3}, Vec3{4, 5, 6}, Vec3{7, 8, 10})
	eqF("Vec3ScalarTripleProduct", f, -3, err, 0)
	r, err = Vec3VectorTripleProduct(Vec3{1, 0, 0}, Vec3{1, 0, 0}, Vec3{0, 1, 0})
	eq("Vec3VectorTripleProduct", r, Vec3{0, -1, 0}, err, 0)
	r, err = Vec3PlaneNormal(Vec3{0, 0, 0}, Vec3{1, 0, 0}, Vec3{0, 1, 0})
	eq("Vec3PlaneNormal", r, Vec3{0, 0, 1}, err, 0)
	r, err = Vec3Rotate(Vec3{1, 0, 0}, Vec3{0, 0, 1}, math.Pi/2)
	eq("Vec3Rotate", r, Vec3{0, 1, 0}, err, 1e-16)
	r, err = Vec3Lerp(Vec3{0, 0, 0}, Vec3{2, 4, 6}, 0.25)
	eq("Vec3Lerp", r, Vec3{0.5, 1, 1.5}, err, 0)
	r, err = Vec3Min(Vec3{1, 5, -2}, Vec3{3, 0, -4})
	eq("Vec3Min", r, Vec3{1, 0, -4}, err, 0)
	r, err = Vec3Max(Vec3{1, 5, -2}, Vec3{3, 0, -4})
	eq("Vec3Max", r, Vec3{3, 5, -2}, err, 0)
}

func TestVec3_errors(t *testing.T) {
	nan, big := math.NaN(), math.MaxFloat64
	z, x := Vec3{}, Vec3{1, 0, 0}
	n := Vec3{nan, 0, 0}
	cases := []struct {
		name string
		err  error
		want error
	}{
		{"Add NaN", second(Vec3Add(n, x)), ErrNotFinite},
		{"Add overflow", second(Vec3Add(Vec3{big, 0, 0}, Vec3{big, 0, 0})), ErrOverflow},
		{"Sub NaN", second(Vec3Sub(x, n)), ErrNotFinite},
		{"Sub overflow", second(Vec3Sub(Vec3{big, 0, 0}, Vec3{-big, 0, 0})), ErrOverflow},
		{"Scale overflow", second(Vec3Scale(Vec3{big, 0, 0}, 3)), ErrOverflow},
		{"Div zero", second(Vec3Div(x, 0)), ErrZeroDivisor},
		{"Div overflow", second(Vec3Div(Vec3{big, 0, 0}, 0.25)), ErrOverflow},
		{"Dot overflow", second(Vec3Dot(Vec3{big, 0, 0}, Vec3{big, 0, 0})), ErrOverflow},
		{"Cross overflow", second(Vec3Cross(Vec3{big, 0, 0}, Vec3{0, big, 0})), ErrOverflow},
		{"Norm overflow", second(Vec3Norm(Vec3{big, big, big})), ErrOverflow},
		{"Normalize zero", second(Vec3Normalize(z)), ErrZeroVector},
		{"Distance overflow", second(Vec3Distance(Vec3{big, 0, 0}, Vec3{-big, 0, 0})), ErrOverflow},
		{"Angle zero", second(Vec3AngleBetween(z, x)), ErrZeroVector},
		{"DirectionAngles zero", second(Vec3DirectionAngles(z)), ErrZeroVector},
		{"ComponentAlong zero axis", second(Vec3ComponentAlong(x, z)), ErrZeroVector},
		{"Project zero", second(Vec3Project(x, z)), ErrZeroVector},
		{"Reject zero", second(Vec3Reject(x, z)), ErrZeroVector},
		{"Reflect zero", second(Vec3Reflect(x, z)), ErrZeroVector},
		{"Reflect NaN", second(Vec3Reflect(n, x)), ErrNotFinite},
		{"Area overflow", second(Vec3ParallelogramArea(Vec3{big, 0, 0}, Vec3{0, big, 0})), ErrOverflow},
		{"Triangle overflow", second(Vec3TriangleArea(Vec3{big, 0, 0}, Vec3{0, big, 0})), ErrOverflow},
		{"Triple overflow", second(Vec3ScalarTripleProduct(Vec3{big, 0, 0}, Vec3{0, big, 0}, Vec3{0, 0, 2})), ErrOverflow},
		{"VectorTriple NaN", second(Vec3VectorTripleProduct(x, x, n)), ErrNotFinite},
		{"PlaneNormal collinear", second(Vec3PlaneNormal(z, x, Vec3{2, 0, 0})), ErrZeroVector},
		{"PlaneNormal coincident", second(Vec3PlaneNormal(x, x, Vec3{0, 1, 0})), ErrZeroVector},
		{"PlaneNormal overflow", second(Vec3PlaneNormal(Vec3{-big, 0, 0}, Vec3{big, 0, 0}, x)), ErrOverflow},
		{"Rotate zero axis", second(Vec3Rotate(x, z, 1)), ErrZeroVector},
		{"Rotate NaN angle", second(Vec3Rotate(x, x, nan)), ErrNotFinite},
		{"Lerp NaN t", second(Vec3Lerp(x, x, nan)), ErrNotFinite},
		{"Lerp overflow", second(Vec3Lerp(Vec3{0, 0, 0}, Vec3{big, 0, 0}, 4)), ErrOverflow},
		{"Min NaN", second(Vec3Min(n, x)), ErrNotFinite},
		{"Max NaN", second(Vec3Max(x, n)), ErrNotFinite},
	}
	for _, c := range cases {
		if !errors.Is(c.err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, c.err, c.want)
		}
	}
}

// Oracle identities on random vectors (fixed seed): Lagrange's identity |a×b|² + (a·b)² = |a|²|b|², the
// cross product is perpendicular to both factors, the triple product equals the 3×3 determinant of the
// rows, a × (b × c) = b(a·c) − c(a·b), projection + rejection = v, reflecting twice is the identity,
// rotation agrees with the axis-angle matrix and preserves length.
func TestVec3_identities(t *testing.T) {
	r := newRNG()
	rv := func() Vec3 { return Vec3{sym(r), sym(r), sym(r)} }
	for trial := range 500 {
		a, b, c := rv(), rv(), rv()
		cr, _ := Vec3Cross(a, b)
		d, _ := Vec3Dot(a, b)
		na, _ := Vec3Norm(a)
		nb, _ := Vec3Norm(b)
		nc, _ := Vec3Norm(cr)
		if !near(nc*nc+d*d, na*na*nb*nb, 8*ulp1) {
			t.Fatalf("trial %d: Lagrange identity fails", trial)
		}
		if p, _ := Vec3Dot(cr, a); math.Abs(p) > 8*ulp1 {
			t.Fatalf("trial %d: (a×b)·a = %v", trial, p)
		}
		tp, _ := Vec3ScalarTripleProduct(a, b, c)
		det, _ := Mat3Det(Mat3{a[0], a[1], a[2], b[0], b[1], b[2], c[0], c[1], c[2]})
		if !near(tp, det, 8*ulp1) {
			t.Fatalf("trial %d: triple %v ≠ det %v", trial, tp, det)
		}
		vt, _ := Vec3VectorTripleProduct(a, b, c)
		ac, _ := Vec3Dot(a, c)
		want := Vec3{b[0]*ac - c[0]*d, b[1]*ac - c[1]*d, b[2]*ac - c[2]*d}
		if !nearV3(vt, want, 16*ulp1) {
			t.Fatalf("trial %d: BAC−CAB: %v vs %v", trial, vt, want)
		}
		p, _ := Vec3Project(a, b)
		q, _ := Vec3Reject(a, b)
		if !nearV3(Vec3{p[0] + q[0], p[1] + q[1], p[2] + q[2]}, a, 8*ulp1) {
			t.Fatalf("trial %d: projection + rejection ≠ v", trial)
		}
		m1, _ := Vec3Reflect(a, b)
		m2, _ := Vec3Reflect(m1, b)
		if !nearV3(m2, a, 16*ulp1) {
			t.Fatalf("trial %d: reflecting twice gives %v, want %v", trial, m2, a)
		}
		th := sym(r) * 4
		rot, _ := Vec3Rotate(a, b, th)
		m, _ := Mat3RotationAxisAngle(b, th)
		mv, _ := Mat3MulVec(m, a)
		nr, _ := Vec3Norm(rot)
		if !nearV3(rot, mv, 16*ulp1) || !near(nr, na, 16*ulp1) {
			t.Fatalf("trial %d: Rodrigues %v vs matrix %v", trial, rot, mv)
		}
		ang, _ := Vec3DirectionAngles(a)
		if s := math.Cos(ang[0])*math.Cos(ang[0]) + math.Cos(ang[1])*math.Cos(ang[1]) + math.Cos(ang[2])*math.Cos(ang[2]); !near(s, 1, 8*ulp1) {
			t.Fatalf("trial %d: direction cosines² sum to %v", trial, s)
		}
	}
}

func TestVec3Lerp_endpointsExact(t *testing.T) {
	a, b := Vec3{0.1, -3.3, 7e10}, Vec3{0.7, 2.9, -1e-3}
	if r, _ := Vec3Lerp(a, b, 0); r != a {
		t.Errorf("t = 0: %v", r)
	}
	if r, _ := Vec3Lerp(a, b, 1); r != b {
		t.Errorf("t = 1: %v", r)
	}
	big := math.MaxFloat64
	if r, err := Vec3Lerp(Vec3{-big, 0, 0}, Vec3{big, 0, 0}, 0.5); err != nil || r[0] != 0 {
		t.Errorf("lerp across the range = %v, %v", r, err)
	}
}

func TestVec3Norm_scaling(t *testing.T) {
	if n, err := Vec3Norm(Vec3{2e300, 3e300, 6e300}); err != nil || !near(n, 7e300, 2*ulp1) {
		t.Errorf("big = %v, %v", n, err)
	}
	if n, err := Vec3Norm(Vec3{2e-300, 3e-300, 6e-300}); err != nil || !near(n/7e-300, 1, 2*ulp1) {
		t.Errorf("tiny = %v, %v", n, err)
	}
	if u, err := Vec3Normalize(Vec3{0, 0, 5e-324}); err != nil || u != (Vec3{0, 0, 1}) {
		t.Errorf("smallest subnormal = %v, %v", u, err)
	}
}

func TestVec3AngleBetween_nearlyParallel(t *testing.T) {
	got, _ := Vec3AngleBetween(Vec3{1, 0, 0}, Vec3{1, 1e-9, 0})
	if want := math.Atan(1e-9); math.Abs(got-want) > 1e-24 {
		t.Errorf("angle = %v, want %v", got, want)
	}
}

func TestVec3_noAllocation(t *testing.T) {
	a, b, c := Vec3{1, 2, 3}, Vec3{-2, 0.5, 4}, Vec3{0, 1, 1}
	if n := testing.AllocsPerRun(100, func() {
		sinkV3, _ = Vec3Add(a, b)
		sinkV3, _ = Vec3Sub(a, b)
		sinkV3, _ = Vec3Scale(a, 2)
		sinkV3, _ = Vec3Div(a, 2)
		sinkF, _ = Vec3Dot(a, b)
		sinkV3, _ = Vec3Cross(a, b)
		sinkF, _ = Vec3Norm(a)
		sinkV3, _ = Vec3Normalize(a)
		sinkF, _ = Vec3Distance(a, b)
		sinkF, _ = Vec3AngleBetween(a, b)
		sinkV3, _ = Vec3DirectionAngles(a)
		sinkF, _ = Vec3ComponentAlong(a, b)
		sinkV3, _ = Vec3Project(a, b)
		sinkV3, _ = Vec3Reject(a, b)
		sinkV3, _ = Vec3Reflect(a, b)
		sinkF, _ = Vec3ParallelogramArea(a, b)
		sinkF, _ = Vec3TriangleArea(a, b)
		sinkF, _ = Vec3ScalarTripleProduct(a, b, c)
		sinkV3, _ = Vec3VectorTripleProduct(a, b, c)
		sinkV3, _ = Vec3PlaneNormal(a, b, c)
		sinkV3, _ = Vec3Rotate(a, b, 1)
		sinkV3, _ = Vec3Lerp(a, b, 0.3)
		sinkV3, _ = Vec3Min(a, b)
		sinkV3, _ = Vec3Max(a, b)
	}); n != 0 {
		t.Errorf("allocations = %v, want 0", n)
	}
}

func BenchmarkVec3Add(b *testing.B) {
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
			sinkV3, _ = Vec3Add(a3, b3)
		}
	})
}

func BenchmarkVec3Sub(b *testing.B) {
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
			sinkV3, _ = Vec3Sub(a3, b3)
		}
	})
}

func BenchmarkVec3Scale(b *testing.B) {
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
			sinkV3, _ = Vec3Scale(a3, 1.5)
		}
	})
}

func BenchmarkVec3Div(b *testing.B) {
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
			sinkV3, _ = Vec3Div(a3, 1.5)
		}
	})
}

func BenchmarkVec3Dot(b *testing.B) {
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
			sinkF, _ = Vec3Dot(a3, b3)
		}
	})
}

func BenchmarkVec3Cross(b *testing.B) {
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
			sinkV3, _ = Vec3Cross(a3, b3)
		}
	})
}

func BenchmarkVec3Norm(b *testing.B) {
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
			sinkF, _ = Vec3Norm(a3)
		}
	})
}

func BenchmarkVec3Normalize(b *testing.B) {
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
			sinkV3, _ = Vec3Normalize(a3)
		}
	})
}

func BenchmarkVec3Distance(b *testing.B) {
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
			sinkF, _ = Vec3Distance(a3, b3)
		}
	})
}

func BenchmarkVec3AngleBetween(b *testing.B) {
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
			sinkF, _ = Vec3AngleBetween(a3, b3)
		}
	})
}

func BenchmarkVec3DirectionAngles(b *testing.B) {
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
			sinkV3, _ = Vec3DirectionAngles(a3)
		}
	})
}

func BenchmarkVec3ComponentAlong(b *testing.B) {
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
			sinkF, _ = Vec3ComponentAlong(a3, b3)
		}
	})
}

func BenchmarkVec3Project(b *testing.B) {
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
			sinkV3, _ = Vec3Project(a3, b3)
		}
	})
}

func BenchmarkVec3Reject(b *testing.B) {
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
			sinkV3, _ = Vec3Reject(a3, b3)
		}
	})
}

func BenchmarkVec3Reflect(b *testing.B) {
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
			sinkV3, _ = Vec3Reflect(a3, b3)
		}
	})
}

func BenchmarkVec3ParallelogramArea(b *testing.B) {
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
			sinkF, _ = Vec3ParallelogramArea(a3, b3)
		}
	})
}

func BenchmarkVec3TriangleArea(b *testing.B) {
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
			sinkF, _ = Vec3TriangleArea(a3, b3)
		}
	})
}

func BenchmarkVec3ScalarTripleProduct(b *testing.B) {
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
			sinkF, _ = Vec3ScalarTripleProduct(a3, b3, c3)
		}
	})
}

func BenchmarkVec3VectorTripleProduct(b *testing.B) {
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
			sinkV3, _ = Vec3VectorTripleProduct(a3, b3, c3)
		}
	})
}

func BenchmarkVec3PlaneNormal(b *testing.B) {
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
			sinkV3, _ = Vec3PlaneNormal(a3, b3, c3)
		}
	})
}

func BenchmarkVec3Rotate(b *testing.B) {
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
			sinkV3, _ = Vec3Rotate(a3, b3, 0.7)
		}
	})
}

func BenchmarkVec3Lerp(b *testing.B) {
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
			sinkV3, _ = Vec3Lerp(a3, b3, 0.3)
		}
	})
}

func BenchmarkVec3Min(b *testing.B) {
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
			sinkV3, _ = Vec3Min(a3, b3)
		}
	})
}

func BenchmarkVec3Max(b *testing.B) {
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
			sinkV3, _ = Vec3Max(a3, b3)
		}
	})
}

// Review regressions (findings 2 and 6): finite results near the float64 maximum.
func TestVec3_nearMaxFloat(t *testing.T) {
	// a × b = (−1.5e308, 1.5e308, 0): finite, but its length √2·1.5e308 is not; half of it is.
	if a, err := Vec3TriangleArea(Vec3{0, 0, 1e154}, Vec3{1.5e154, 1.5e154, 0}); err != nil || !near(a/(math.Sqrt2*0.75e308), 1, 4*ulp1) {
		t.Errorf("Vec3TriangleArea = %v, %v", a, err)
	}
	// The cross product itself overflows; halving a first keeps the half-area finite (fallback branch).
	if a, err := Vec3TriangleArea(Vec3{1e300, 0, 0}, Vec3{0, 3e8, 0}); err != nil || !near(a/1.5e308, 1, 4*ulp1) {
		t.Errorf("Vec3TriangleArea fallback = %v, %v", a, err)
	}
	if r, err := Vec3Reject(Vec3{1.5e308, 1.5e308, 0}, Vec3{1, 1, 0}); err != nil || math.Abs(r[0]) > 1e293 || math.Abs(r[1]) > 1e293 {
		t.Errorf("Vec3Reject = %v, %v; want ≈ 0", r, err) // rounding of the projection: ~ulp(1.5e308)
	}
	if r, err := Vec3Reflect(Vec3{1e308, 1e308, 0}, Vec3{1, 1, 0}); err != nil || !nearV3(Vec3{r[0] / 1e308, r[1] / 1e308, r[2]}, Vec3{-1, -1, 0}, 4*ulp1) {
		t.Errorf("Vec3Reflect = %v, %v", r, err)
	}
	if p, err := Vec3Project(Vec3{1e308, 1e308, 0}, Vec3{1, 0, 0}); err != nil || p[0] != 1e308 {
		t.Errorf("Vec3Project = %v, %v", p, err)
	}
	if d, err := Vec3Distance(Vec3{1e308, 0, 0}, Vec3{0, 1e308, 0}); err != nil || !near(d/1e308, math.Sqrt2, 4*ulp1) {
		t.Errorf("Vec3Distance √2·1e308 = %v, %v", d, err)
	}
	if d, err := Vec3Distance(Vec3{1.5e308, 0, 0}, Vec3{0, 1.5e308, 0}); !errors.Is(err, ErrOverflow) {
		t.Errorf("Vec3Distance √2·1.5e308 = %v, %v; want ErrOverflow", d, err)
	}
}

// Review regression (finding 5): a product overflows but the determinant is finite.
func TestDet_intermediateOverflow(t *testing.T) {
	if c, err := Vec2Cross(Vec2{1e200, 1e200}, Vec2{1e200, 1e200}); err != nil || c != 0 {
		t.Errorf("Vec2Cross = %v, %v; want 0", c, err)
	}
	if d, err := Vec2Dot(Vec2{1e200, 1e200}, Vec2{1e200, -1e200}); err != nil || d != 0 {
		t.Errorf("Vec2Dot = %v, %v; want 0", d, err)
	}
	if d, err := Mat3Det(Mat3{1e200, 1e200, 0, 1e200, 1e200, 0, 0, 0, 1}); err != nil || d != 0 {
		t.Errorf("Mat3Det singular = %v, %v; want 0", d, err)
	}
	if d, err := Mat3Det(Mat3{1e200, 0, 0, 0, 1e200, 0, 0, 0, 1e-300}); err != nil || !near(d/1e100, 1, 4*ulp1) {
		t.Errorf("Mat3Det scaled = %v, %v; want 1e100", d, err)
	}
	if d, err := Mat4Det(Mat4{1e200, 0, 0, 0, 0, 1e200, 0, 0, 0, 0, 1e-300, 0, 0, 0, 0, 1}); err != nil || !near(d/1e100, 1, 4*ulp1) {
		t.Errorf("Mat4Det scaled = %v, %v; want 1e100", d, err)
	}
	if c, err := Vec3Cross(Vec3{1e200, 1e200, 0}, Vec3{1e200, 1e200, 0}); err != nil || c != (Vec3{}) {
		t.Errorf("Vec3Cross parallel huge = %v, %v; want 0", c, err)
	}
}
