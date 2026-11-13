package linalg

import (
	"errors"
	"math"
	"testing"
)

func TestCoordinates_workedExamples(t *testing.T) {
	// Hand-checked from DLMF §1.5(ii) eqs. 1.5.14 and 1.5.16.
	c, err := CartesianToCylindrical(Vec3{1, 1, 5})
	if err != nil || !near(c.R, math.Sqrt2, ulp1) || !near(c.PhiRad, math.Pi/4, ulp1) || c.Z != 5 {
		t.Errorf("CartesianToCylindrical = %+v, %v", c, err)
	}
	v, err := CylindricalToCartesian(Cylindrical{2, math.Pi / 2, 3})
	if err != nil || !nearV3(v, Vec3{0, 2, 3}, 2e-16) {
		t.Errorf("CylindricalToCartesian = %v, %v", v, err)
	}
	s, err := CartesianToSpherical(Vec3{0, 1, 1})
	if err != nil || !near(s.R, math.Sqrt2, ulp1) || !near(s.ThetaRad, math.Pi/4, ulp1) || !near(s.PhiRad, math.Pi/2, ulp1) {
		t.Errorf("CartesianToSpherical = %+v, %v", s, err)
	}
	v, err = SphericalToCartesian(Spherical{2, math.Pi / 2, 0})
	if err != nil || !nearV3(v, Vec3{2, 0, 0}, 2e-16) {
		t.Errorf("SphericalToCartesian = %v, %v", v, err)
	}
}

func TestCoordinates_conventions(t *testing.T) {
	// φ ∈ [0, 2π): the −y axis is 3π/2, not −π/2; θ is measured from +z.
	if c, _ := CartesianToCylindrical(Vec3{0, -1, 0}); !near(c.PhiRad, 3*math.Pi/2, ulp1) {
		t.Errorf("azimuth of −ŷ = %v", c.PhiRad)
	}
	if s, _ := CartesianToSpherical(Vec3{0, 0, -2}); s.ThetaRad != math.Pi || s.R != 2 || s.PhiRad != 0 {
		t.Errorf("south pole = %+v", s)
	}
	// Origin: angles undefined, 0 returned; −0 is folded to +0.
	if s, _ := CartesianToSpherical(Vec3{0, math.Copysign(0, -1), 0}); s != (Spherical{}) || math.Signbit(s.PhiRad) {
		t.Errorf("origin = %+v", s)
	}
	// A tiny negative azimuth that would round up to 2π is folded to 0.
	if c, _ := CartesianToCylindrical(Vec3{1, -1e-300, 0}); c.PhiRad < 0 || c.PhiRad >= 2*math.Pi {
		t.Errorf("azimuth range = %v", c.PhiRad)
	}
	// Near the pole θ keeps its digits (arccos(z/R) would give 0).
	if s, _ := CartesianToSpherical(Vec3{1e-10, 0, 1}); math.Abs(s.ThetaRad-1e-10) > 1e-25 {
		t.Errorf("θ near the pole = %v", s.ThetaRad)
	}
	if s, err := CartesianToSpherical(Vec3{3e300, 4e300, 0}); err != nil || !near(s.R, 5e300, ulp1) {
		t.Errorf("large R = %+v, %v", s, err)
	}
}

func TestCoordinates_errors(t *testing.T) {
	nan, big := math.NaN(), math.MaxFloat64
	cases := []struct {
		name string
		err  error
		want error
	}{
		{"toCyl NaN", second(CartesianToCylindrical(Vec3{nan, 0, 0})), ErrNotFinite},
		{"toCyl overflow", second(CartesianToCylindrical(Vec3{big, big, 0})), ErrOverflow},
		{"fromCyl R<0", second(CylindricalToCartesian(Cylindrical{-1, 0, 0})), ErrOutOfRange},
		{"fromCyl NaN φ", second(CylindricalToCartesian(Cylindrical{1, nan, 0})), ErrNotFinite},
		{"fromCyl NaN R", second(CylindricalToCartesian(Cylindrical{nan, 0, 0})), ErrNotFinite},
		{"fromCyl Inf Z", second(CylindricalToCartesian(Cylindrical{1, 0, math.Inf(1)})), ErrNotFinite},
		{"toSph NaN", second(CartesianToSpherical(Vec3{0, 0, nan})), ErrNotFinite},
		{"toSph overflow", second(CartesianToSpherical(Vec3{big, big, big})), ErrOverflow},
		{"fromSph R<0", second(SphericalToCartesian(Spherical{-1, 0, 0})), ErrOutOfRange},
		{"fromSph NaN θ", second(SphericalToCartesian(Spherical{1, nan, 0})), ErrNotFinite},
		{"fromSph NaN φ", second(SphericalToCartesian(Spherical{1, 0, nan})), ErrNotFinite},
		{"fromSph NaN R", second(SphericalToCartesian(Spherical{nan, 0, 0})), ErrNotFinite},
	}
	for _, c := range cases {
		if !errors.Is(c.err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, c.err, c.want)
		}
	}
}

// Oracle: round trips through both coordinate systems return the point (fixed seed).
func TestCoordinates_roundTrip(t *testing.T) {
	r := newRNG()
	for trial := range 500 {
		v := Vec3{sym(r) * 10, sym(r) * 10, sym(r) * 10}
		c, _ := CartesianToCylindrical(v)
		w, _ := CylindricalToCartesian(c)
		s, _ := CartesianToSpherical(v)
		u, _ := SphericalToCartesian(s)
		// sin/cos/atan2/hypot each round once: a handful of ulps of |v| ≤ 10√3, so 1e−14 absolute.
		if !nearV3(w, v, 1e-14) || !nearV3(u, v, 1e-14) {
			t.Fatalf("trial %d: %v → %v, %v", trial, v, w, u)
		}
		if c.PhiRad < 0 || c.PhiRad >= 2*math.Pi || s.ThetaRad < 0 || s.ThetaRad > math.Pi {
			t.Fatalf("trial %d: angles out of range %+v %+v", trial, c, s)
		}
	}
}

func TestCoordinates_noAllocation(t *testing.T) {
	v := Vec3{1, 2, 3}
	if n := testing.AllocsPerRun(100, func() {
		sinkCyl, _ = CartesianToCylindrical(v)
		sinkV3, _ = CylindricalToCartesian(sinkCyl)
		sinkSph, _ = CartesianToSpherical(v)
		sinkV3, _ = SphericalToCartesian(sinkSph)
	}); n != 0 {
		t.Errorf("allocations = %v, want 0", n)
	}
}

func BenchmarkCartesianToCylindrical(b *testing.B) {
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
			sinkCyl, _ = CartesianToCylindrical(a3)
		}
	})
}

func BenchmarkCylindricalToCartesian(b *testing.B) {
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
			sinkV3, _ = CylindricalToCartesian(Cylindrical{2, 0.7, 3})
		}
	})
}

func BenchmarkCartesianToSpherical(b *testing.B) {
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
			sinkSph, _ = CartesianToSpherical(a3)
		}
	})
}

func BenchmarkSphericalToCartesian(b *testing.B) {
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
			sinkV3, _ = SphericalToCartesian(Spherical{2, 0.7, 1.1})
		}
	})
}

// Review regression (finding 8): the origin with z = −0 has angles 0, not θ = π.
func TestCartesianToSpherical_negativeZeroOrigin(t *testing.T) {
	if s, err := CartesianToSpherical(Vec3{0, 0, math.Copysign(0, -1)}); err != nil || s != (Spherical{}) {
		t.Errorf("origin (0, 0, −0) = %+v, %v", s, err)
	}
}
