package linalg

import (
	"errors"
	"math"
	"testing"
)

func TestMat4_workedExamples(t *testing.T) {
	// Hand-checked examples (Strang 2016 §2.4, §5.2, §5.3, §10.6 pp. 495–497).
	id := Mat4Identity()
	if id != (Mat4{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}) {
		t.Errorf("Mat4Identity = %v", id)
	}
	t1, _ := Mat4Translation(Vec3{1, 2, 3})
	t2, _ := Mat4Translation(Vec3{4, 5, 6})
	t3, _ := Mat4Translation(Vec3{5, 7, 9})
	if m, err := Mat4Mul(t1, t2); err != nil || m != t3 {
		t.Errorf("Mat4Mul of translations = %v, %v", m, err)
	}
	if v, err := Mat4MulVec(t1, Vec4{0, 0, 0, 1}); err != nil || v != (Vec4{1, 2, 3, 1}) {
		t.Errorf("Mat4MulVec = %v, %v", v, err)
	}
	if m, err := Mat4Transpose(t1); err != nil || m[12] != 1 || m[13] != 2 || m[14] != 3 || m[3] != 0 {
		t.Errorf("Mat4Transpose = %v, %v", m, err)
	}
	diag := Mat4{1, 0, 0, 0, 0, 2, 0, 0, 0, 0, 3, 0, 0, 0, 0, 4}
	if d, err := Mat4Det(diag); err != nil || d != 24 {
		t.Errorf("Mat4Det diag = %v, %v", d, err)
	}
	if d, err := Mat4Det(t1); err != nil || d != 1 {
		t.Errorf("Mat4Det translation = %v, %v", d, err)
	}
	tinv, _ := Mat4Translation(Vec3{-1, -2, -3})
	if m, err := Mat4Inverse(t1); err != nil || !nearM(m[:], tinv[:], 0) {
		t.Errorf("Mat4Inverse translation = %v, %v", m, err)
	}
	if tr, err := Mat4Trace(diag); err != nil || tr != 10 {
		t.Errorf("Mat4Trace = %v, %v", tr, err)
	}
	s, _ := Mat4Scale(Vec3{2, 3, 4})
	if p, err := Mat4TransformPoint(s, Vec3{1, 1, 1}); err != nil || p != (Vec3{2, 3, 4}) {
		t.Errorf("Mat4Scale point = %v, %v", p, err)
	}
	rz, _ := Mat3RotationZ(math.Pi / 2)
	h, _ := Mat4FromMat3(rz)
	if v, err := Mat4MulVec(h, Vec4{1, 0, 0, 1}); err != nil || !near(v[0], 0, 1e-16) || v[1] != 1 || v[3] != 1 {
		t.Errorf("Mat4FromMat3 = %v, %v", v, err)
	}
	if p, err := Mat4TransformPoint(t1, Vec3{1, 1, 1}); err != nil || p != (Vec3{2, 3, 4}) {
		t.Errorf("Mat4TransformPoint = %v, %v", p, err)
	}
	if d, err := Mat4TransformDirection(t1, Vec3{1, 1, 1}); err != nil || d != (Vec3{1, 1, 1}) {
		t.Errorf("Mat4TransformDirection = %v, %v", d, err)
	}
	// A perspective row: w′ = z, so (2, 4, 2) maps to (1, 2, 1).
	persp := Mat4{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0}
	if p, err := Mat4TransformPoint(persp, Vec3{2, 4, 2}); err != nil || p != (Vec3{1, 2, 1}) {
		t.Errorf("perspective divide = %v, %v", p, err)
	}
}

func TestMat4_errors(t *testing.T) {
	nan, big := math.NaN(), math.MaxFloat64
	bad := Mat4{nan}
	persp := Mat4{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0}
	cases := []struct {
		name string
		err  error
		want error
	}{
		{"Mul NaN", second(Mat4Mul(bad, Mat4{})), ErrNotFinite},
		{"Mul overflow", second(Mat4Mul(Mat4{big}, Mat4{2})), ErrOverflow},
		{"MulVec NaN", second(Mat4MulVec(Mat4{}, Vec4{nan})), ErrNotFinite},
		{"MulVec overflow", second(Mat4MulVec(Mat4{big, big}, Vec4{1, 1})), ErrOverflow},
		{"Transpose NaN", second(Mat4Transpose(bad)), ErrNotFinite},
		{"Det NaN", second(Mat4Det(bad)), ErrNotFinite},
		{"Det overflow", second(Mat4Det(Mat4{big, 0, 0, 0, 0, big, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1})), ErrOverflow},
		{"Inverse singular", second(Mat4Inverse(Mat4{1, 2, 3, 4, 2, 4, 6, 8, 0, 0, 1, 0, 0, 0, 0, 1})), ErrSingular},
		{"Inverse zero", second(Mat4Inverse(Mat4{})), ErrSingular},
		{"Inverse NaN", second(Mat4Inverse(bad)), ErrNotFinite},
		{"Inverse overflow", second(Mat4Inverse(Mat4{1, 0, 0, 0, 0, 1e-310, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1})), ErrOverflow},
		{"Trace overflow", second(Mat4Trace(Mat4{big, 0, 0, 0, 0, big})), ErrOverflow},
		{"Translation NaN", second(Mat4Translation(Vec3{nan})), ErrNotFinite},
		{"Scale NaN", second(Mat4Scale(Vec3{nan})), ErrNotFinite},
		{"FromMat3 NaN", second(Mat4FromMat3(Mat3{nan})), ErrNotFinite},
		{"TransformPoint w = 0", second(Mat4TransformPoint(persp, Vec3{1, 1, 0})), ErrOutOfRange},
		{"TransformPoint NaN", second(Mat4TransformPoint(Mat4Identity(), Vec3{nan})), ErrNotFinite},
		{"TransformPoint overflow", second(Mat4TransformPoint(Mat4{big, 0, 0, big, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}, Vec3{1, 0, 0})), ErrOverflow},
		{"TransformDirection NaN", second(Mat4TransformDirection(bad, Vec3{})), ErrNotFinite},
		{"TransformDirection overflow", second(Mat4TransformDirection(Mat4{big, big}, Vec3{1, 1})), ErrOverflow},
	}
	for _, c := range cases {
		if !errors.Is(c.err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, c.err, c.want)
		}
	}
}

// Oracle (fixed seed): the cofactor inverse and determinant agree with LU (MatInverse, MatDet), A·A⁻¹ = I,
// (AB)ᵀ = BᵀAᵀ.
func TestMat4_oracle(t *testing.T) {
	r := newRNG()
	for trial := range 300 {
		var a, b Mat4
		for i := range a {
			a[i], b[i] = sym(r)*4, sym(r)*4
		}
		d, _ := Mat4Det(a)
		dlu, _ := MatDet(a[:], 4)
		cond, err := MatConditionEstimate1(a[:], 4)
		if err != nil {
			continue
		}
		if math.Abs(d-dlu) > 1e-13*cond*(1+math.Abs(d)) {
			t.Fatalf("trial %d: Mat4Det %v ≠ MatDet %v", trial, d, dlu)
		}
		inv, _ := Mat4Inverse(a)
		ai, _ := Mat4Mul(a, inv)
		if id := Mat4Identity(); !nearM(ai[:], id[:], 1e-14*cond) {
			t.Fatalf("trial %d: A·A⁻¹ = %v", trial, ai)
		}
		inv2, _ := MatInverse(nil, a[:], 4)
		if !nearM(inv[:], inv2, 1e-14*cond) {
			t.Fatalf("trial %d: cofactor ≠ LU inverse", trial)
		}
		ab, _ := Mat4Mul(a, b)
		abt, _ := Mat4Transpose(ab)
		at, _ := Mat4Transpose(a)
		bt, _ := Mat4Transpose(b)
		btat, _ := Mat4Mul(bt, at)
		if abt != btat {
			t.Fatalf("trial %d: (AB)ᵀ ≠ BᵀAᵀ", trial)
		}
	}
}

func TestMat4_noAllocation(t *testing.T) {
	a := Mat4{2, 0, 0, 1, 0, 3, 0, 2, 0, 0, 4, 3, 0, 0, 0, 1}
	v := Vec3{1, 2, 3}
	if n := testing.AllocsPerRun(100, func() {
		sinkM4 = Mat4Identity()
		sinkM4, _ = Mat4Mul(a, a)
		sinkV4, _ = Mat4MulVec(a, Vec4{1, 2, 3, 1})
		sinkM4, _ = Mat4Transpose(a)
		sinkF, _ = Mat4Det(a)
		sinkM4, _ = Mat4Inverse(a)
		sinkF, _ = Mat4Trace(a)
		sinkM4, _ = Mat4Translation(v)
		sinkM4, _ = Mat4Scale(v)
		sinkM4, _ = Mat4FromMat3(Mat3Identity())
		sinkV3, _ = Mat4TransformPoint(a, v)
		sinkV3, _ = Mat4TransformDirection(a, v)
	}); n != 0 {
		t.Errorf("allocations = %v, want 0", n)
	}
}

func BenchmarkMat4Identity(b *testing.B) {
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
			sinkM4 = Mat4Identity()
		}
	})
}

func BenchmarkMat4Mul(b *testing.B) {
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
			sinkM4, _ = Mat4Mul(m4a, m4a)
		}
	})
}

func BenchmarkMat4MulVec(b *testing.B) {
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
			sinkV4, _ = Mat4MulVec(m4a, Vec4{1, 2, 3, 1})
		}
	})
}

func BenchmarkMat4Transpose(b *testing.B) {
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
			sinkM4, _ = Mat4Transpose(m4a)
		}
	})
}

func BenchmarkMat4Det(b *testing.B) {
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
			sinkF, _ = Mat4Det(m4a)
		}
	})
}

func BenchmarkMat4Inverse(b *testing.B) {
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
			sinkM4, _ = Mat4Inverse(m4a)
		}
	})
}

func BenchmarkMat4Trace(b *testing.B) {
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
			sinkF, _ = Mat4Trace(m4a)
		}
	})
}

func BenchmarkMat4Translation(b *testing.B) {
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
			sinkM4, _ = Mat4Translation(a3)
		}
	})
}

func BenchmarkMat4Scale(b *testing.B) {
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
			sinkM4, _ = Mat4Scale(a3)
		}
	})
}

func BenchmarkMat4FromMat3(b *testing.B) {
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
			sinkM4, _ = Mat4FromMat3(m3a)
		}
	})
}

func BenchmarkMat4TransformPoint(b *testing.B) {
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
			sinkV3, _ = Mat4TransformPoint(m4a, a3)
		}
	})
}

func BenchmarkMat4TransformDirection(b *testing.B) {
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
			sinkV3, _ = Mat4TransformDirection(m4a, a3)
		}
	})
}
