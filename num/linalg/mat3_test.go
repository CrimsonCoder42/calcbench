package linalg

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

func nearM(a, b []float64, tol float64) bool { return eqSlice(a, b, tol) }

func randM3(r *rand.Rand) Mat3 {
	var m Mat3
	for i := range m {
		m[i] = sym(r) * 4
	}
	return m
}

func randRotation(r *rand.Rand) Mat3 {
	m, _ := Mat3RotationAxisAngle(Vec3{sym(r), sym(r), sym(r) + 2}, sym(r)*math.Pi)
	return m
}

func TestMat2_workedExamples(t *testing.T) {
	// Hand-checked examples from the doc comments (Strang 2016 §2.4, §2.5, §5.1, §6.1).
	if m := Mat2Identity(); m != (Mat2{1, 0, 0, 1}) {
		t.Errorf("Mat2Identity = %v", m)
	}
	if m, err := Mat2Mul(Mat2{1, 2, 3, 4}, Mat2{5, 6, 7, 8}); err != nil || m != (Mat2{19, 22, 43, 50}) {
		t.Errorf("Mat2Mul = %v, %v", m, err)
	}
	if v, err := Mat2MulVec(Mat2{1, 2, 3, 4}, Vec2{5, 6}); err != nil || v != (Vec2{17, 39}) {
		t.Errorf("Mat2MulVec = %v, %v", v, err)
	}
	if m, err := Mat2Transpose(Mat2{1, 2, 3, 4}); err != nil || m != (Mat2{1, 3, 2, 4}) {
		t.Errorf("Mat2Transpose = %v, %v", m, err)
	}
	if d, err := Mat2Det(Mat2{1, 2, 3, 4}); err != nil || d != -2 {
		t.Errorf("Mat2Det = %v, %v", d, err)
	}
	if m, err := Mat2Inverse(Mat2{4, 7, 2, 6}); err != nil || !nearM(m[:], []float64{0.6, -0.7, -0.2, 0.4}, ulp1) {
		t.Errorf("Mat2Inverse = %v, %v", m, err)
	}
	if tr, err := Mat2Trace(Mat2{1, 2, 3, 4}); err != nil || tr != 5 {
		t.Errorf("Mat2Trace = %v, %v", tr, err)
	}
	if m, err := Mat2Rotation(math.Pi / 2); err != nil || !nearM(m[:], []float64{0, -1, 1, 0}, 1e-16) {
		t.Errorf("Mat2Rotation = %v, %v", m, err)
	}
	if lo, hi, err := Mat2SymmetricEigenvalues(Mat2{2, 1, 1, 2}); err != nil || lo != 1 || hi != 3 {
		t.Errorf("Mat2SymmetricEigenvalues = %v, %v, %v", lo, hi, err)
	}
}

func TestMat3_workedExamples(t *testing.T) {
	a := Mat3{1, 2, 3, 0, 1, 4, 5, 6, 0}
	if m := Mat3Identity(); m != (Mat3{1, 0, 0, 0, 1, 0, 0, 0, 1}) {
		t.Errorf("Mat3Identity = %v", m)
	}
	if m, err := Mat3Mul(Mat3{1, 2, 0, 0, 1, 0, 0, 0, 2}, Mat3{1, 0, 0, 3, 1, 0, 0, 0, 1}); err != nil || m != (Mat3{7, 2, 0, 3, 1, 0, 0, 0, 2}) {
		t.Errorf("Mat3Mul = %v, %v", m, err)
	}
	if v, err := Mat3MulVec(Mat3{1, 2, 3, 4, 5, 6, 7, 8, 9}, Vec3{1, 0, -1}); err != nil || v != (Vec3{-2, -2, -2}) {
		t.Errorf("Mat3MulVec = %v, %v", v, err)
	}
	if m, err := Mat3Transpose(Mat3{1, 2, 3, 4, 5, 6, 7, 8, 9}); err != nil || m != (Mat3{1, 4, 7, 2, 5, 8, 3, 6, 9}) {
		t.Errorf("Mat3Transpose = %v, %v", m, err)
	}
	if d, err := Mat3Det(a); err != nil || d != 1 {
		t.Errorf("Mat3Det = %v, %v", d, err)
	}
	// Cofactors of small integers are exact; only the final divisions round: 1e−14 is generous.
	if m, err := Mat3Inverse(a); err != nil || !nearM(m[:], []float64{-24, 18, 5, 20, -15, -4, -5, 4, 1}, 1e-14) {
		t.Errorf("Mat3Inverse = %v, %v", m, err)
	}
	if tr, err := Mat3Trace(Mat3{1, 2, 3, 4, 5, 6, 7, 8, 9}); err != nil || tr != 15 {
		t.Errorf("Mat3Trace = %v, %v", tr, err)
	}
	// Each axis rotation turns the next axis into the one after it (right-hand rule).
	rx, _ := Mat3RotationX(math.Pi / 2)
	ry, _ := Mat3RotationY(math.Pi / 2)
	rz, _ := Mat3RotationZ(math.Pi / 2)
	if v, _ := Mat3MulVec(rx, Vec3{0, 1, 0}); !nearV3(v, Vec3{0, 0, 1}, 1e-16) {
		t.Errorf("Rx(π/2)·ŷ = %v", v)
	}
	if v, _ := Mat3MulVec(ry, Vec3{0, 0, 1}); !nearV3(v, Vec3{1, 0, 0}, 1e-16) {
		t.Errorf("Ry(π/2)·ẑ = %v", v)
	}
	if v, _ := Mat3MulVec(rz, Vec3{1, 0, 0}); !nearV3(v, Vec3{0, 1, 0}, 1e-16) {
		t.Errorf("Rz(π/2)·x̂ = %v", v)
	}
	if m, err := Mat3RotationAxisAngle(Vec3{0, 0, 1}, math.Pi/2); err != nil || !nearM(m[:], rz[:], 2*ulp1) { // 1 − cos θ = 2 sin²(θ/2) carries one rounding
		t.Errorf("Mat3RotationAxisAngle(ẑ, π/2) = %v, %v", m, err)
	}
	if ev, err := Mat3SymmetricEigenvalues(Mat3{2, 0, 0, 0, 3, 4, 0, 4, 9}); err != nil || !nearM(ev[:], []float64{1, 2, 11}, 8*ulp1) {
		t.Errorf("Mat3SymmetricEigenvalues = %v, %v", ev, err)
	}
}

// Diebel 2006 §2.4 eqs. (14)–(16) are frame rotations; the active matrices here must be their transposes.
func TestMat3Rotation_diebelTransposes(t *testing.T) {
	a := 0.7
	c, s := math.Cos(a), math.Sin(a)
	r1 := Mat3{1, 0, 0, 0, c, s, 0, -s, c}
	r2 := Mat3{c, 0, -s, 0, 1, 0, s, 0, c}
	r3 := Mat3{c, s, 0, -s, c, 0, 0, 0, 1}
	rx, _ := Mat3RotationX(a)
	ry, _ := Mat3RotationY(a)
	rz, _ := Mat3RotationZ(a)
	for _, p := range []struct {
		name     string
		got, ref Mat3
	}{{"x", rx, r1}, {"y", ry, r2}, {"z", rz, r3}} {
		if transpose3(p.ref) != p.got {
			t.Errorf("R%s is not the transpose of Diebel's matrix", p.name)
		}
	}
}

func TestMat3_errors(t *testing.T) {
	nan, big := math.NaN(), math.MaxFloat64
	bad := Mat3{nan}
	cases := []struct {
		name string
		err  error
		want error
	}{
		{"Mat2Mul NaN", second(Mat2Mul(Mat2{nan}, Mat2{})), ErrNotFinite},
		{"Mat2Mul overflow", second(Mat2Mul(Mat2{big, 0, 0, 0}, Mat2{2, 0, 0, 0})), ErrOverflow},
		{"Mat2MulVec NaN v", second(Mat2MulVec(Mat2{}, Vec2{nan, 0})), ErrNotFinite},
		{"Mat2MulVec overflow", second(Mat2MulVec(Mat2{big, big, 0, 0}, Vec2{1, 1})), ErrOverflow},
		{"Mat2Transpose NaN", second(Mat2Transpose(Mat2{nan})), ErrNotFinite},
		{"Mat2Det overflow", second(Mat2Det(Mat2{big, 0, 0, big})), ErrOverflow},
		{"Mat2Inverse singular", second(Mat2Inverse(Mat2{1, 2, 2, 4})), ErrSingular},
		{"Mat2Inverse zero", second(Mat2Inverse(Mat2{})), ErrSingular},
		{"Mat2Inverse overflow", second(Mat2Inverse(Mat2{1e-310, 0, 0, 1e-310})), ErrOverflow},
		{"Mat2Trace overflow", second(Mat2Trace(Mat2{big, 0, 0, big})), ErrOverflow},
		{"Mat2Rotation NaN", second(Mat2Rotation(nan)), ErrNotFinite},
		{"Mat2Eigen not symmetric", second3(Mat2SymmetricEigenvalues(Mat2{1, 2, 3, 4})), ErrNotSymmetric},
		{"Mat2Eigen NaN", second3(Mat2SymmetricEigenvalues(Mat2{nan})), ErrNotFinite},
		{"Mat2Eigen overflow", second3(Mat2SymmetricEigenvalues(Mat2{big, big, big, big})), ErrOverflow},
		{"Mat3Mul NaN", second(Mat3Mul(bad, Mat3{})), ErrNotFinite},
		{"Mat3Mul overflow", second(Mat3Mul(Mat3{big}, Mat3{2})), ErrOverflow},
		{"Mat3MulVec NaN", second(Mat3MulVec(Mat3{}, Vec3{nan})), ErrNotFinite},
		{"Mat3MulVec overflow", second(Mat3MulVec(Mat3{big, big}, Vec3{1, 1})), ErrOverflow},
		{"Mat3Transpose NaN", second(Mat3Transpose(bad)), ErrNotFinite},
		{"Mat3Det NaN", second(Mat3Det(bad)), ErrNotFinite},
		{"Mat3Det overflow", second(Mat3Det(Mat3{big, 0, 0, 0, big, 0, 0, 0, big})), ErrOverflow},
		{"Mat3Inverse singular", second(Mat3Inverse(Mat3{1, 2, 3, 4, 5, 6, 7, 8, 9})), ErrSingular},
		{"Mat3Inverse zero", second(Mat3Inverse(Mat3{})), ErrSingular},
		{"Mat3Inverse NaN", second(Mat3Inverse(bad)), ErrNotFinite},
		{"Mat3Inverse overflow", second(Mat3Inverse(Mat3{1, 0, 0, 0, 1e-310, 0, 0, 0, 1})), ErrOverflow},
		{"Mat3Trace overflow", second(Mat3Trace(Mat3{big, 0, 0, 0, big})), ErrOverflow},
		{"RotationX NaN", second(Mat3RotationX(nan)), ErrNotFinite},
		{"RotationY NaN", second(Mat3RotationY(nan)), ErrNotFinite},
		{"RotationZ NaN", second(Mat3RotationZ(nan)), ErrNotFinite},
		{"AxisAngle zero axis", second(Mat3RotationAxisAngle(Vec3{}, 1)), ErrZeroVector},
		{"AxisAngle NaN", second(Mat3RotationAxisAngle(Vec3{1}, nan)), ErrNotFinite},
		{"Mat3Eigen not symmetric", second(Mat3SymmetricEigenvalues(Mat3{1, 2, 0, 0, 1, 0, 0, 0, 1})), ErrNotSymmetric},
		{"Mat3Eigen NaN", second(Mat3SymmetricEigenvalues(bad)), ErrNotFinite},
		{"Mat3Eigen overflow", second(Mat3SymmetricEigenvalues(Mat3{big, big, big, big, big, big, big, big, big})), ErrOverflow},
	}
	for _, c := range cases {
		if !errors.Is(c.err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, c.err, c.want)
		}
	}
}

func second3[A, B any](_ A, _ B, err error) error { return err }

// Oracle identities with a fixed seed: (A·B)ᵀ = BᵀAᵀ, A·A⁻¹ = I, det(AB) = det A·det B, the 2×2 and
// 3×3 closed forms agree with the general LU functions, rotations are orthonormal with det +1.
func TestMat_identities(t *testing.T) {
	r := newRNG()
	for trial := range 300 {
		a, b := randM3(r), randM3(r)
		ab, _ := Mat3Mul(a, b)
		at, _ := Mat3Transpose(a)
		bt, _ := Mat3Transpose(b)
		btat, _ := Mat3Mul(bt, at)
		abt, _ := Mat3Transpose(ab)
		if abt != btat {
			t.Fatalf("trial %d: (AB)ᵀ ≠ BᵀAᵀ", trial)
		}
		inv, err := Mat3Inverse(a)
		if err != nil {
			continue
		}
		cond, _ := MatConditionEstimate1(a[:], 3)
		ai, _ := Mat3Mul(a, inv)
		if id := Mat3Identity(); !nearM(ai[:], id[:], 1e-14*cond) {
			t.Fatalf("trial %d: A·A⁻¹ = %v (κ ≈ %v)", trial, ai, cond)
		}
		inv2, _ := MatInverse(nil, a[:], 3)
		if !nearM(inv[:], inv2, 1e-14*cond) {
			t.Fatalf("trial %d: cofactor inverse ≠ LU inverse", trial)
		}
		da, _ := Mat3Det(a)
		db, _ := Mat3Det(b)
		dab, _ := Mat3Det(ab)
		if math.Abs(dab-da*db) > 1e-12*(1+math.Abs(da*db))*cond {
			t.Fatalf("trial %d: det(AB) = %v, det A det B = %v", trial, dab, da*db)
		}
		dlu, _ := MatDet(a[:], 3)
		if math.Abs(da-dlu) > 1e-13*cond*(1+math.Abs(da)) {
			t.Fatalf("trial %d: Mat3Det %v ≠ MatDet %v", trial, da, dlu)
		}
		m2 := Mat2{a[0], a[1], a[3], a[4]}
		m2i, err := Mat2Inverse(m2)
		if err == nil {
			p, _ := Mat2Mul(m2, m2i)
			c2, _ := MatConditionEstimate1(m2[:], 2)
			if !nearM(p[:], []float64{1, 0, 0, 1}, 1e-14*c2) {
				t.Fatalf("trial %d: 2×2 A·A⁻¹ = %v", trial, p)
			}
		}
		rot := randRotation(r)
		rt, _ := Mat3Transpose(rot)
		rrt, _ := Mat3Mul(rot, rt)
		dr, _ := Mat3Det(rot)
		if id := Mat3Identity(); !nearM(rrt[:], id[:], 8*ulp1) || !near(dr, 1, 8*ulp1) {
			t.Fatalf("trial %d: rotation not orthonormal: RRᵀ = %v, det %v", trial, rrt, dr)
		}
	}
}

// Oracle: a symmetric matrix built as R·D·Rᵀ from a random rotation R and known D has eigenvalues D.
func TestMatSymmetricEigenvalues_oracle(t *testing.T) {
	r := newRNG()
	for trial := range 500 {
		d := []float64{sym(r) * 10, sym(r) * 10, sym(r) * 10}
		if trial%5 == 0 {
			d[1] = d[0] // repeated eigenvalue
		}
		rot := randRotation(r)
		var a Mat3
		for i := range 3 {
			for j := range 3 {
				a[3*i+j] = rot[3*i]*d[0]*rot[3*j] + rot[3*i+1]*d[1]*rot[3*j+1] + rot[3*i+2]*d[2]*rot[3*j+2]
			}
		}
		// Symmetrise exactly (rounding makes aᵢⱼ and aⱼᵢ differ in the last bit).
		a[3], a[6], a[7] = a[1], a[2], a[5]
		ev, err := Mat3SymmetricEigenvalues(a)
		slices.Sort(d)
		// Distinct eigenvalues: a few ulps of ‖A‖ (≤ 10 here). A repeated eigenvalue makes p³ − q² vanish,
		// and its rounding error enters through a square root: about √ε·‖A‖ (documented on the function).
		tol := 1e-13 * 10
		if trial%5 == 0 {
			tol = 1e-6
		}
		if err != nil || !nearM(ev[:], d, tol) {
			t.Fatalf("trial %d: eigenvalues %v, want %v (%v)", trial, ev, d, err)
		}
		m2 := Mat2{a[0], a[1], a[1], a[4]}
		lo, hi, _ := Mat2SymmetricEigenvalues(m2)
		// λ² − tλ + det = 0 for both roots.
		tr, det := a[0]+a[4], a[0]*a[4]-a[1]*a[1]
		for _, l := range []float64{lo, hi} {
			if math.Abs(l*l-tr*l+det) > 1e-12*(1+l*l+math.Abs(det)) {
				t.Fatalf("trial %d: %v is not a root of the 2×2 characteristic polynomial", trial, l)
			}
		}
	}
}

// Row equilibration: a badly scaled but perfectly conditioned matrix inverts (the unscaled cofactor
// determinant 1e−600 would underflow to 0 and report a false ErrSingular).
func TestMat3Inverse_badlyScaled(t *testing.T) {
	m, err := Mat3Inverse(Mat3{1, 0, 0, 0, 1e-300, 0, 0, 0, 1e-300})
	if err != nil || m[0] != 1 || !near(m[4], 1e300, ulp1) || !near(m[8], 1e300, ulp1) {
		t.Errorf("Mat3Inverse = %v, %v", m, err)
	}
	m2, err := Mat2Inverse(Mat2{1e-200, 0, 0, 1e-200})
	if err != nil || !near(m2[0], 1e200, ulp1) {
		t.Errorf("Mat2Inverse = %v, %v", m2, err)
	}
	m4, err := Mat4Inverse(Mat4{1, 0, 0, 0, 0, 1e-300, 0, 0, 0, 0, 1e-300, 0, 0, 0, 0, 1})
	if err != nil || !near(m4[5], 1e300, ulp1) {
		t.Errorf("Mat4Inverse = %v, %v", m4, err)
	}
}

func TestMat2SymmetricEigenvalues_noCancellation(t *testing.T) {
	// [1 1e-9; 1e-9 1e-18·…]: the small eigenvalue is det/big, not big − r.
	lo, hi, _ := Mat2SymmetricEigenvalues(Mat2{1, 1e-5, 1e-5, 2e-10})
	if want := 1e-10 / hi; math.Abs(lo-want) > 1e-24 { // 1e−24 ≈ 1e−14 relative to the 1e−10 eigenvalue; big − r would lose all digits
		t.Errorf("small eigenvalue = %v, want %v", lo, want)
	}
}

func TestMat3_noAllocation(t *testing.T) {
	a, b := Mat3{1, 2, 3, 0, 1, 4, 5, 6, 0}, Mat3{2, 0, 0, 0, 3, 4, 0, 4, 9}
	m2 := Mat2{4, 7, 2, 6}
	v := Vec3{1, 2, 3}
	if n := testing.AllocsPerRun(100, func() {
		sinkM2 = Mat2Identity()
		sinkM2, _ = Mat2Mul(m2, m2)
		sinkV2, _ = Mat2MulVec(m2, Vec2{1, 2})
		sinkM2, _ = Mat2Transpose(m2)
		sinkF, _ = Mat2Det(m2)
		sinkM2, _ = Mat2Inverse(m2)
		sinkF, _ = Mat2Trace(m2)
		sinkM2, _ = Mat2Rotation(1)
		sinkF, _, _ = Mat2SymmetricEigenvalues(Mat2{2, 1, 1, 2})
		sinkM3 = Mat3Identity()
		sinkM3, _ = Mat3Mul(a, b)
		sinkV3, _ = Mat3MulVec(a, v)
		sinkM3, _ = Mat3Transpose(a)
		sinkF, _ = Mat3Det(a)
		sinkM3, _ = Mat3Inverse(a)
		sinkF, _ = Mat3Trace(a)
		sinkM3, _ = Mat3RotationX(1)
		sinkM3, _ = Mat3RotationY(1)
		sinkM3, _ = Mat3RotationZ(1)
		sinkM3, _ = Mat3RotationAxisAngle(v, 1)
		sinkEig3, _ = Mat3SymmetricEigenvalues(b)
	}); n != 0 {
		t.Errorf("allocations = %v, want 0", n)
	}
}

func BenchmarkMat2Identity(b *testing.B) {
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
			sinkM2 = Mat2Identity()
		}
	})
}

func BenchmarkMat2Mul(b *testing.B) {
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
			sinkM2, _ = Mat2Mul(m2a, m2b)
		}
	})
}

func BenchmarkMat2MulVec(b *testing.B) {
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
			sinkV2, _ = Mat2MulVec(m2a, a2)
		}
	})
}

func BenchmarkMat2Transpose(b *testing.B) {
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
			sinkM2, _ = Mat2Transpose(m2a)
		}
	})
}

func BenchmarkMat2Det(b *testing.B) {
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
			sinkF, _ = Mat2Det(m2a)
		}
	})
}

func BenchmarkMat2Inverse(b *testing.B) {
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
			sinkM2, _ = Mat2Inverse(m2a)
		}
	})
}

func BenchmarkMat2Trace(b *testing.B) {
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
			sinkF, _ = Mat2Trace(m2a)
		}
	})
}

func BenchmarkMat2Rotation(b *testing.B) {
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
			sinkM2, _ = Mat2Rotation(0.7)
		}
	})
}

func BenchmarkMat2SymmetricEigenvalues(b *testing.B) {
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
			sinkF, _, _ = Mat2SymmetricEigenvalues(m2s)
		}
	})
}

func BenchmarkMat3Identity(b *testing.B) {
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
			sinkM3 = Mat3Identity()
		}
	})
}

func BenchmarkMat3Mul(b *testing.B) {
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
			sinkM3, _ = Mat3Mul(m3a, m3b)
		}
	})
}

func BenchmarkMat3MulVec(b *testing.B) {
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
			sinkV3, _ = Mat3MulVec(m3a, a3)
		}
	})
}

func BenchmarkMat3Transpose(b *testing.B) {
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
			sinkM3, _ = Mat3Transpose(m3a)
		}
	})
}

func BenchmarkMat3Det(b *testing.B) {
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
			sinkF, _ = Mat3Det(m3a)
		}
	})
}

func BenchmarkMat3Inverse(b *testing.B) {
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
			sinkM3, _ = Mat3Inverse(m3a)
		}
	})
}

func BenchmarkMat3Trace(b *testing.B) {
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
			sinkF, _ = Mat3Trace(m3a)
		}
	})
}

func BenchmarkMat3RotationX(b *testing.B) {
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
			sinkM3, _ = Mat3RotationX(0.7)
		}
	})
}

func BenchmarkMat3RotationY(b *testing.B) {
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
			sinkM3, _ = Mat3RotationY(0.7)
		}
	})
}

func BenchmarkMat3RotationZ(b *testing.B) {
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
			sinkM3, _ = Mat3RotationZ(0.7)
		}
	})
}

func BenchmarkMat3RotationAxisAngle(b *testing.B) {
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
			sinkM3, _ = Mat3RotationAxisAngle(a3, 0.7)
		}
	})
}

func BenchmarkMat3SymmetricEigenvalues(b *testing.B) {
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
			sinkEig3, _ = Mat3SymmetricEigenvalues(m3s)
		}
	})
}

// Review regression (finding 3): column scaling as well as row scaling. With e = 1e−200 the inverse is
// [3 −1 −1; −1e200 1e200 0; −1e200 0 1e200] (MatInverse agrees); rows alone left det underflowing.
func TestMatInverse_columnScaled(t *testing.T) {
	e := 1e-200
	want := []float64{3, -1, -1, -1e200, 1e200, 0, -1e200, 0, 1e200}
	m, err := Mat3Inverse(Mat3{1, e, e, 1, 2 * e, e, 1, e, 2 * e})
	if err != nil {
		t.Fatalf("Mat3Inverse: %v", err)
	}
	for i := range want {
		if !near(m[i], want[i], 1e-14*max(1, math.Abs(want[i]))) {
			t.Fatalf("Mat3Inverse = %v, want %v", m, want)
		}
	}
	m4, err := Mat4Inverse(Mat4{1, e, e, 0, 1, 2 * e, e, 0, 1, e, 2 * e, 0, 0, 0, 0, 1})
	if err != nil || !near(m4[0], 3, 1e-14) || !near(m4[4]/-1e200, 1, 1e-14) {
		t.Errorf("Mat4Inverse = %v, %v", m4, err)
	}
	// Genuine overflow: the inverse of [1 1e−320; 1 2e−320] has entries ±1e320.
	if _, err := Mat2Inverse(Mat2{1, 1e-320, 1, 2e-320}); !errors.Is(err, ErrOverflow) {
		t.Errorf("Mat2Inverse genuine overflow: err = %v", err)
	}
}

// Review regression (finding 9) and the branches the coverage run found untested.
func TestSymmetricEigenvalues_branches(t *testing.T) {
	if lo, hi, err := Mat2SymmetricEigenvalues(Mat2{5e-324, 0, 0, 5e-324}); err != nil || lo != 5e-324 || hi != 5e-324 {
		t.Errorf("subnormal 2×2 = %v, %v, %v", lo, hi, err)
	}
	if lo, hi, err := Mat2SymmetricEigenvalues(Mat2{}); err != nil || lo != 0 || hi != 0 {
		t.Errorf("zero 2×2 = %v, %v, %v", lo, hi, err)
	}
	if ev, err := Mat3SymmetricEigenvalues(Mat3{}); err != nil || ev != [3]float64{} {
		t.Errorf("zero 3×3 = %v, %v", ev, err)
	}
	if ev, err := Mat3SymmetricEigenvalues(Mat3{3, 0, 0, 0, -1, 0, 0, 0, 2}); err != nil || ev != [3]float64{-1, 2, 3} {
		t.Errorf("diagonal 3×3 = %v, %v", ev, err)
	}
}
