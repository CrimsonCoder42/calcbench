package linalg

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

func nearQ(a, b Quat, tol float64) bool {
	return near(a.W, b.W, tol) && near(a.X, b.X, tol) && near(a.Y, b.Y, tol) && near(a.Z, b.Z, tol)
}

// sameRotation reports whether a and b are equal up to sign (q and −q are one rotation).
func sameRotation(a, b Quat, tol float64) bool {
	return nearQ(a, b, tol) || nearQ(a, Quat{-b.W, -b.X, -b.Y, -b.Z}, tol)
}

func randQ(r *rand.Rand) Quat { return Quat{sym(r), sym(r), sym(r), sym(r)} }

func TestQuat_workedExamples(t *testing.T) {
	h := math.Sqrt2 / 2
	if q, err := QuatAdd(Quat{1, 2, 3, 4}, Quat{5, 6, 7, 8}); err != nil || q != (Quat{6, 8, 10, 12}) {
		t.Errorf("QuatAdd = %v, %v", q, err)
	}
	if q, err := QuatSub(Quat{5, 6, 7, 8}, Quat{1, 2, 3, 4}); err != nil || q != (Quat{4, 4, 4, 4}) {
		t.Errorf("QuatSub = %v, %v", q, err)
	}
	// Hamilton: ij = k, ji = −k, i² = −1 (Shoemake 1985 §2.3).
	i, j, k := Quat{0, 1, 0, 0}, Quat{0, 0, 1, 0}, Quat{0, 0, 0, 1}
	if q, _ := QuatMul(i, j); q != k {
		t.Errorf("ij = %v", q)
	}
	if q, _ := QuatMul(j, i); q != (Quat{0, 0, 0, -1}) {
		t.Errorf("ji = %v", q)
	}
	if q, _ := QuatMul(i, i); q != (Quat{-1, 0, 0, 0}) {
		t.Errorf("i² = %v", q)
	}
	if q, _ := QuatMul(Quat{0, 1, 0, 0}, Quat{0, 0, 1, 0}); q != (Quat{0, 0, 0, 1}) {
		t.Errorf("Shoemake's book flip [0,(0,1,0)]·[0,(0,0,1)] in Hamilton order = %v", q)
	}
	if q, err := QuatConjugate(Quat{1, 2, 3, 4}); err != nil || q != (Quat{1, -2, -3, -4}) {
		t.Errorf("QuatConjugate = %v, %v", q, err)
	}
	if n, err := QuatNorm(Quat{1, 2, 2, 4}); err != nil || n != 5 {
		t.Errorf("QuatNorm = %v, %v", n, err)
	}
	if q, err := QuatNormalize(Quat{1, 2, 2, 4}); err != nil || !nearQ(q, Quat{0.2, 0.4, 0.4, 0.8}, ulp1) {
		t.Errorf("QuatNormalize = %v, %v", q, err)
	}
	if q, err := QuatInverse(Quat{1, 1, 1, 1}); err != nil || !nearQ(q, Quat{0.25, -0.25, -0.25, -0.25}, ulp1) {
		t.Errorf("QuatInverse = %v, %v", q, err)
	}
	qz := Quat{h, 0, 0, h}
	if q, err := QuatFromAxisAngle(Vec3{0, 0, 1}, math.Pi/2); err != nil || !nearQ(q, qz, ulp1) {
		t.Errorf("QuatFromAxisAngle = %v, %v", q, err)
	}
	if ax, ang, err := QuatToAxisAngle(qz); err != nil || !nearV3(ax, Vec3{0, 0, 1}, ulp1) || !near(ang, math.Pi/2, 2*ulp1) {
		t.Errorf("QuatToAxisAngle = %v, %v, %v", ax, ang, err)
	}
	rz, _ := Mat3RotationZ(math.Pi / 2)
	if q, err := QuatFromMat3(rz); err != nil || !nearQ(q, qz, 2*ulp1) {
		t.Errorf("QuatFromMat3 = %v, %v", q, err)
	}
	if m, err := QuatToMat3(qz); err != nil || !nearM(m[:], rz[:], 4*ulp1) {
		t.Errorf("QuatToMat3 = %v, %v", m, err)
	}
	if v, err := QuatRotateVec3(qz, Vec3{1, 0, 0}); err != nil || !nearV3(v, Vec3{0, 1, 0}, 4*ulp1) {
		t.Errorf("QuatRotateVec3 = %v, %v", v, err)
	}
	q45, _ := QuatFromAxisAngle(Vec3{0, 0, 1}, math.Pi/4)
	if q, err := QuatSlerp(Quat{1, 0, 0, 0}, qz, 0.5); err != nil || !nearQ(q, q45, 2*ulp1) {
		t.Errorf("QuatSlerp = %v, %v", q, err)
	}
	if q, err := QuatFromEulerZYX(math.Pi/2, 0, 0); err != nil || !nearQ(q, qz, ulp1) {
		t.Errorf("QuatFromEulerZYX = %v, %v", q, err)
	}
	if y, p, r, err := QuatToEulerZYX(qz, -1); err != nil || !near(y, math.Pi/2, 2*ulp1) || p != 0 || r != 0 {
		t.Errorf("QuatToEulerZYX = %v, %v, %v, %v", y, p, r, err)
	}
}

func TestQuat_errors(t *testing.T) {
	nan, big := math.NaN(), math.MaxFloat64
	bad, zero, one := Quat{nan, 0, 0, 0}, Quat{}, Quat{1, 0, 0, 0}
	cases := []struct {
		name string
		err  error
		want error
	}{
		{"Add NaN", second(QuatAdd(bad, one)), ErrNotFinite},
		{"Add overflow", second(QuatAdd(Quat{W: big}, Quat{W: big})), ErrOverflow},
		{"Sub NaN", second(QuatSub(one, bad)), ErrNotFinite},
		{"Sub overflow", second(QuatSub(Quat{W: big}, Quat{W: -big})), ErrOverflow},
		{"Mul NaN", second(QuatMul(bad, one)), ErrNotFinite},
		{"Mul overflow", second(QuatMul(Quat{W: big}, Quat{W: 2})), ErrOverflow},
		{"Conjugate NaN", second(QuatConjugate(bad)), ErrNotFinite},
		{"Norm NaN", second(QuatNorm(bad)), ErrNotFinite},
		{"Norm overflow", second(QuatNorm(Quat{big, big, big, big})), ErrOverflow},
		{"Normalize zero", second(QuatNormalize(zero)), ErrZeroVector},
		{"Normalize NaN", second(QuatNormalize(bad)), ErrNotFinite},
		{"Inverse zero", second(QuatInverse(zero)), ErrZeroVector},
		{"Inverse NaN", second(QuatInverse(bad)), ErrNotFinite},
		{"Inverse overflow", second(QuatInverse(Quat{1e-310, 0, 0, 0})), ErrOverflow},
		{"FromAxisAngle zero", second(QuatFromAxisAngle(Vec3{}, 1)), ErrZeroVector},
		{"FromAxisAngle NaN", second(QuatFromAxisAngle(Vec3{1}, nan)), ErrNotFinite},
		{"FromAxisAngle NaN axis", second(QuatFromAxisAngle(Vec3{nan}, 1)), ErrNotFinite},
		{"ToAxisAngle zero", second3(QuatToAxisAngle(zero)), ErrZeroVector},
		{"ToAxisAngle NaN", second3(QuatToAxisAngle(bad)), ErrNotFinite},
		{"FromMat3 NaN", second(QuatFromMat3(Mat3{nan})), ErrNotFinite},
		{"FromMat3 overflow", second(QuatFromMat3(Mat3{big, 0, 0, 0, big, 0, 0, 0, big})), ErrOverflow},
		{"ToMat3 zero", second(QuatToMat3(zero)), ErrZeroVector},
		{"ToMat3 NaN", second(QuatToMat3(bad)), ErrNotFinite},
		{"Rotate zero", second(QuatRotateVec3(zero, Vec3{1})), ErrZeroVector},
		{"Rotate NaN v", second(QuatRotateVec3(one, Vec3{nan})), ErrNotFinite},
		{"Rotate NaN q", second(QuatRotateVec3(bad, Vec3{1})), ErrNotFinite},
		{"Rotate overflow", second(QuatRotateVec3(Quat{1, 1, 0, 0}, Vec3{0, big, big})), ErrOverflow},
		{"Slerp t<0", second(QuatSlerp(one, one, -0.1)), ErrOutOfRange},
		{"Slerp t>1", second(QuatSlerp(one, one, 1.1)), ErrOutOfRange},
		{"Slerp zero", second(QuatSlerp(zero, one, 0.5)), ErrZeroVector},
		{"Slerp zero q1", second(QuatSlerp(one, zero, 0.5)), ErrZeroVector},
		{"Slerp NaN t", second(QuatSlerp(one, one, nan)), ErrNotFinite},
		{"Slerp NaN q0", second(QuatSlerp(bad, one, 0.5)), ErrNotFinite},
		{"Slerp NaN q1", second(QuatSlerp(one, bad, 0.5)), ErrNotFinite},
		{"FromEuler NaN yaw", second(QuatFromEulerZYX(nan, 0, 0)), ErrNotFinite},
		{"FromEuler NaN pitch", second(QuatFromEulerZYX(0, nan, 0)), ErrNotFinite},
		{"FromEuler NaN roll", second(QuatFromEulerZYX(0, 0, nan)), ErrNotFinite},
		{"ToEuler zero", euler4(QuatToEulerZYX(zero, -1)), ErrZeroVector},
		{"ToEuler NaN", euler4(QuatToEulerZYX(bad, -1)), ErrNotFinite},
		{"ToEuler NaN tol", euler4(QuatToEulerZYX(one, nan)), ErrNotFinite},
	}
	for _, c := range cases {
		if !errors.Is(c.err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, c.err, c.want)
		}
	}
}

func euler4(_, _, _ float64, err error) error { return err }

// diebelRq is Diebel 2006 §6.4 eq. (125) written out literally: the frame rotation of a unit quaternion.
func diebelRq(q Quat) Mat3 {
	q0, q1, q2, q3 := q.W, q.X, q.Y, q.Z
	return Mat3{
		q0*q0 + q1*q1 - q2*q2 - q3*q3, 2*q1*q2 + 2*q0*q3, 2*q1*q3 - 2*q0*q2,
		2*q1*q2 - 2*q0*q3, q0*q0 - q1*q1 + q2*q2 - q3*q3, 2*q2*q3 + 2*q0*q1,
		2*q1*q3 + 2*q0*q2, 2*q2*q3 - 2*q0*q1, q0*q0 - q1*q1 - q2*q2 + q3*q3,
	}
}

// diebelR123 is Diebel 2006 §5.6.2 eq. (67) written out literally: R₁(φ)R₂(θ)R₃(ψ), φ roll, θ pitch, ψ yaw.
func diebelR123(phi, theta, psi float64) Mat3 {
	sf, cf := math.Sincos(phi)
	st, ct := math.Sincos(theta)
	ss, cs := math.Sincos(psi)
	return Mat3{
		ct * cs, ct * ss, -st,
		sf*st*cs - cf*ss, sf*st*ss + cf*cs, ct * sf,
		cf*st*cs + sf*ss, cf*st*ss - sf*cs, ct * cf,
	}
}

// Oracle (fixed seed): the quaternion and rotation-matrix paths agree on random rotations and vectors, and
// both agree with the sources written out literally (Diebel eq. 125 and eq. 67, transposed to the active
// convention); identities |pq| = |p||q|, q·q⁻¹ = 1, |normalize(q)| = 1, and round trips hold.
func TestQuat_oracle(t *testing.T) {
	r := newRNG()
	for trial := range 500 {
		q, p := randQ(r), randQ(r)
		u, _ := QuatNormalize(q)
		if n, _ := QuatNorm(u); !near(n, 1, 2*ulp1) {
			t.Fatalf("trial %d: |normalize(q)| = %v", trial, n)
		}
		m, _ := QuatToMat3(q)
		if lit := transpose3(diebelRq(u)); !nearM(m[:], lit[:], 8*ulp1) {
			t.Fatalf("trial %d: QuatToMat3 ≠ Diebel eq. 125 transposed", trial)
		}
		v := Vec3{sym(r) * 5, sym(r) * 5, sym(r) * 5}
		qv, _ := QuatRotateVec3(q, v)
		mv, _ := Mat3MulVec(m, v)
		if !nearV3(qv, mv, 32*ulp1) {
			t.Fatalf("trial %d: quaternion rotation %v ≠ matrix rotation %v", trial, qv, mv)
		}
		back, _ := QuatFromMat3(m)
		if !sameRotation(back, u, 16*ulp1) {
			t.Fatalf("trial %d: FromMat3(ToMat3(q)) = %v, want ±%v", trial, back, u)
		}
		pq, _ := QuatMul(p, q)
		np, _ := QuatNorm(p)
		nq, _ := QuatNorm(q)
		npq, _ := QuatNorm(pq)
		if !near(npq, np*nq, 8*ulp1) {
			t.Fatalf("trial %d: |pq| ≠ |p||q|", trial)
		}
		// Composition: rotating by p·q applies q, then p.
		v1, _ := QuatRotateVec3(q, v)
		v2, _ := QuatRotateVec3(p, v1)
		v3, _ := QuatRotateVec3(pq, v)
		if !nearV3(v2, v3, 64*ulp1) {
			t.Fatalf("trial %d: composition order", trial)
		}
		inv, _ := QuatInverse(q)
		one, _ := QuatMul(q, inv)
		if !nearQ(one, Quat{1, 0, 0, 0}, 8*ulp1) {
			t.Fatalf("trial %d: q·q⁻¹ = %v", trial, one)
		}
		ax, ang, _ := QuatToAxisAngle(q)
		q2, _ := QuatFromAxisAngle(ax, ang)
		if !sameRotation(q2, u, 16*ulp1) || ang < 0 || ang > math.Pi {
			t.Fatalf("trial %d: axis-angle round trip %v → %v", trial, u, q2)
		}
		// Euler: the quaternion of (yaw, pitch, roll) is Rz·Ry·Rx, which is Diebel's R₁₂₃ transposed.
		yaw, pitch, roll := sym(r)*math.Pi, sym(r)*1.5, sym(r)*math.Pi
		qe, _ := QuatFromEulerZYX(yaw, pitch, roll)
		me, _ := QuatToMat3(qe)
		if lit := transpose3(diebelR123(roll, pitch, yaw)); !nearM(me[:], lit[:], 8*ulp1) {
			t.Fatalf("trial %d: Euler quaternion ≠ Diebel eq. 67 transposed", trial)
		}
		rz, _ := Mat3RotationZ(yaw)
		ry, _ := Mat3RotationY(pitch)
		rx, _ := Mat3RotationX(roll)
		zy, _ := Mat3Mul(rz, ry)
		zyx, _ := Mat3Mul(zy, rx)
		if !nearM(me[:], zyx[:], 8*ulp1) {
			t.Fatalf("trial %d: Euler quaternion ≠ Rz·Ry·Rx", trial)
		}
		y2, p2, r2, _ := QuatToEulerZYX(qe, -1)
		// |pitch| ≤ 1.5 keeps cos θ ≥ 0.07, so atan2/asin conditioning costs at most ~2 digits: 1e−12.
		if !near(y2, yaw, 1e-12) || !near(p2, pitch, 1e-12) || !near(r2, roll, 1e-12) {
			t.Fatalf("trial %d: Euler round trip (%v, %v, %v) → (%v, %v, %v)", trial, yaw, pitch, roll, y2, p2, r2)
		}
	}
}

// Gimbal lock (Diebel §5.6.6): at pitch ±π/2 only yaw − roll (or yaw + roll) is defined. The returned
// angles have roll 0 and reproduce the same rotation.
func TestQuatToEulerZYX_gimbalLock(t *testing.T) {
	for _, pitch := range []float64{math.Pi / 2, -math.Pi / 2} {
		q, _ := QuatFromEulerZYX(0.3, pitch, 0.1)
		y, p, r, err := QuatToEulerZYX(q, -1)
		if err != nil || r != 0 || p != pitch {
			t.Fatalf("pitch %v: got (%v, %v, %v), %v", pitch, y, p, r, err)
		}
		want, _ := QuatToMat3(q)
		q2, _ := QuatFromEulerZYX(y, p, r)
		got, _ := QuatToMat3(q2)
		if !nearM(got[:], want[:], 1e-15) {
			t.Fatalf("pitch %v: the returned angles describe another rotation", pitch)
		}
	}
	// Just off the lock the generic formulas run; with a looser lockTol the same input is treated as locked.
	q, _ := QuatFromEulerZYX(0.3, math.Pi/2-1e-5, 0.1)
	if _, _, r, _ := QuatToEulerZYX(q, -1); !near(r, 0.1, 1e-6) {
		t.Errorf("near lock, default tol: roll = %v, want ≈ 0.1", r)
	}
	if _, _, r, _ := QuatToEulerZYX(q, 1e-9); r != 0 {
		t.Errorf("near lock, lockTol 1e−9: roll = %v, want 0", r)
	}
}

func TestQuatSlerp_properties(t *testing.T) {
	r := newRNG()
	for trial := range 300 {
		a, _ := QuatNormalize(randQ(r))
		b, _ := QuatNormalize(randQ(r))
		s0, _ := QuatSlerp(a, b, 0)
		if na, _ := QuatNormalize(a); s0 != na {
			t.Fatalf("trial %d: slerp(0) = %v, want exactly the normalised q0 %v", trial, s0, na)
		}
		s1, _ := QuatSlerp(a, b, 1)
		if !sameRotation(s1, b, 4*ulp1) {
			t.Fatalf("trial %d: slerp(1) = %v, want ±%v", trial, s1, b)
		}
		// Constant angular rate: the angle from a grows linearly in t, and |slerp| = 1.
		_, full, _ := QuatToAxisAngle(mulConjQ(a, s1))
		tt := r.Float64()
		st, _ := QuatSlerp(a, b, tt)
		_, part, _ := QuatToAxisAngle(mulConjQ(a, st))
		n, _ := QuatNorm(st)
		// The angle is recovered through two quaternion products and atan2: ~1e−15; 1e−12 leaves margin.
		if !near(part, tt*full, 1e-12) || !near(n, 1, 4*ulp1) {
			t.Fatalf("trial %d: angle at t = %v is %v, want %v (|q| = %v)", trial, tt, part, tt*full, n)
		}
	}
	// Shorter arc: q and −q give the same path.
	a, _ := QuatFromAxisAngle(Vec3{0, 0, 1}, 0.2)
	b, _ := QuatFromAxisAngle(Vec3{0, 0, 1}, 1.0)
	s, _ := QuatSlerp(a, Quat{-b.W, -b.X, -b.Y, -b.Z}, 0.5)
	want, _ := QuatFromAxisAngle(Vec3{0, 0, 1}, 0.6)
	if !sameRotation(s, want, 4*ulp1) {
		t.Errorf("shorter arc: %v, want %v", s, want)
	}
	// Identical inputs: θ = 0 uses the limit weights (1 − t, t) without dividing by zero.
	if s, err := QuatSlerp(a, a, 0.3); err != nil || !nearQ(s, a, 2*ulp1) {
		t.Errorf("slerp(a, a) = %v, %v", s, err)
	}
}

// mulConjQ returns conj(a)·b, the relative rotation from a to b.
func mulConjQ(a, b Quat) Quat { return mulQ(Quat{a.W, -a.X, -a.Y, -a.Z}, b) }

func TestQuatToAxisAngle_smallAngle(t *testing.T) {
	q, _ := QuatFromAxisAngle(Vec3{1, 0, 0}, 1e-12)
	_, ang, _ := QuatToAxisAngle(q)
	if math.Abs(ang-1e-12) > 1e-27 {
		t.Errorf("small angle = %v (2·acos(W) would give 0)", ang)
	}
	if ax, ang, _ := QuatToAxisAngle(Quat{1, 0, 0, 0}); ax != (Vec3{1, 0, 0}) || ang != 0 {
		t.Errorf("identity = %v, %v", ax, ang)
	}
}

func TestQuat_noAllocation(t *testing.T) {
	p, q := Quat{1, 2, 3, 4}, Quat{0.5, -1, 0.25, 2}
	m, _ := QuatToMat3(q)
	if n := testing.AllocsPerRun(100, func() {
		sinkQ, _ = QuatAdd(p, q)
		sinkQ, _ = QuatSub(p, q)
		sinkQ, _ = QuatMul(p, q)
		sinkQ, _ = QuatConjugate(q)
		sinkF, _ = QuatNorm(q)
		sinkQ, _ = QuatNormalize(q)
		sinkQ, _ = QuatInverse(q)
		sinkQ, _ = QuatFromAxisAngle(Vec3{1, 2, 3}, 1)
		sinkV3, sinkF, _ = QuatToAxisAngle(q)
		sinkQ, _ = QuatFromMat3(m)
		sinkM3, _ = QuatToMat3(q)
		sinkV3, _ = QuatRotateVec3(q, Vec3{1, 2, 3})
		sinkQ, _ = QuatSlerp(p, q, 0.3)
		sinkQ, _ = QuatFromEulerZYX(0.1, 0.2, 0.3)
		sinkF, _, _, _ = QuatToEulerZYX(q, -1)
	}); n != 0 {
		t.Errorf("allocations = %v, want 0", n)
	}
}

func BenchmarkQuatAdd(b *testing.B) {
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
			sinkQ, _ = QuatAdd(qa, qb)
		}
	})
}

func BenchmarkQuatSub(b *testing.B) {
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
			sinkQ, _ = QuatSub(qa, qb)
		}
	})
}

func BenchmarkQuatMul(b *testing.B) {
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
			sinkQ, _ = QuatMul(qa, qb)
		}
	})
}

func BenchmarkQuatConjugate(b *testing.B) {
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
			sinkQ, _ = QuatConjugate(qa)
		}
	})
}

func BenchmarkQuatNorm(b *testing.B) {
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
			sinkF, _ = QuatNorm(qa)
		}
	})
}

func BenchmarkQuatNormalize(b *testing.B) {
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
			sinkQ, _ = QuatNormalize(qa)
		}
	})
}

func BenchmarkQuatInverse(b *testing.B) {
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
			sinkQ, _ = QuatInverse(qa)
		}
	})
}

func BenchmarkQuatFromAxisAngle(b *testing.B) {
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
			sinkQ, _ = QuatFromAxisAngle(a3, 0.7)
		}
	})
}

func BenchmarkQuatToAxisAngle(b *testing.B) {
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
			sinkV3, sinkF, _ = QuatToAxisAngle(qa)
		}
	})
}

func BenchmarkQuatFromMat3(b *testing.B) {
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
			sinkQ, _ = QuatFromMat3(m3r)
		}
	})
}

func BenchmarkQuatToMat3(b *testing.B) {
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
			sinkM3, _ = QuatToMat3(qa)
		}
	})
}

func BenchmarkQuatRotateVec3(b *testing.B) {
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
			sinkV3, _ = QuatRotateVec3(qa, a3)
		}
	})
}

func BenchmarkQuatSlerp(b *testing.B) {
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
			sinkQ, _ = QuatSlerp(qa, qb, 0.3)
		}
	})
}

func BenchmarkQuatFromEulerZYX(b *testing.B) {
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
			sinkQ, _ = QuatFromEulerZYX(0.1, 0.2, 0.3)
		}
	})
}

func BenchmarkQuatToEulerZYX(b *testing.B) {
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
			sinkF, _, _, _ = QuatToEulerZYX(qa, -1)
		}
	})
}

// Review regressions (finding 7): |q| beyond the float64 range still inverts.
func TestQuatInverse_large(t *testing.T) {
	q, err := QuatInverse(Quat{1.5e308, 1.5e308, 0, 0})
	if want := 1 / 3e308; err != nil || !near(q.W/want, 1, 1e-12) || !near(q.X/-want, 1, 1e-12) {
		t.Errorf("QuatInverse = %v, %v; want ≈ (%v, %v, 0, 0)", q, err, want, -want) // subnormal result: ~12 digits
	}
}
