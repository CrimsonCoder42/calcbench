package linalg

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Quat is a quaternion q = W + Xi + Yj + Zk: W the scalar part, (X, Y, Z) the vector part. The fields
// are named so the scalar-first / scalar-last ambiguity of array layouts cannot arise. Products are
// Hamilton's (i² = j² = k² = ijk = −1). A unit quaternion represents the rotation v ↦ q v q⁻¹ (active,
// right-hand rule); q and −q represent the same rotation.
type Quat struct {
	W, X, Y, Z float64
}

func finiteQ(q Quat) bool { return reduce.AllFinite(q.W, q.X, q.Y, q.Z) }

// norm4 returns |q| without overflow or underflow in the squares (exact power-of-two scaling).
func norm4(q Quat) float64 {
	m := max(math.Abs(q.W), math.Abs(q.X), math.Abs(q.Y), math.Abs(q.Z))
	if m == 0 {
		return 0
	}
	s, e := pow2Scale(m)
	w, x, y, z := q.W*s, q.X*s, q.Y*s, q.Z*s
	return math.Ldexp(math.Sqrt(math.FMA(w, w, math.FMA(x, x, math.FMA(y, y, z*z)))), e)
}

// unitQ returns q/|q| and false for the zero quaternion.
func unitQ(q Quat) (Quat, bool) {
	m := max(math.Abs(q.W), math.Abs(q.X), math.Abs(q.Y), math.Abs(q.Z))
	if m == 0 {
		return Quat{}, false
	}
	s, _ := pow2Scale(m)
	w, x, y, z := q.W*s, q.X*s, q.Y*s, q.Z*s
	n := math.Sqrt(math.FMA(w, w, math.FMA(x, x, math.FMA(y, y, z*z))))
	return Quat{w / n, x / n, y / n, z / n}, true
}

// mulQ is the Hamilton product [p₀, p] [q₀, q] = [p₀q₀ − p·q, p₀q + q₀p + p × q].
func mulQ(p, q Quat) Quat {
	return Quat{
		W: math.FMA(p.W, q.W, -math.FMA(p.X, q.X, math.FMA(p.Y, q.Y, p.Z*q.Z))),
		X: math.FMA(p.W, q.X, math.FMA(p.X, q.W, det2(p.Y, p.Z, q.Y, q.Z))),
		Y: math.FMA(p.W, q.Y, math.FMA(p.Y, q.W, det2(p.Z, p.X, q.Z, q.X))),
		Z: math.FMA(p.W, q.Z, math.FMA(p.Z, q.W, det2(p.X, p.Y, q.X, q.Y))),
	}
}

// QuatAdd returns p + q, component by component.
// Errors: ErrNotFinite for a NaN or ±Inf component; ErrOverflow when a component of the sum is outside the
// float64 range.
//
// ID:      NUM-LINALG-101
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Shoemake, "Animating Rotation with Quaternion Curves", Computer Graphics (SIGGRAPH '85), 1985, 19(3):245–254, doi:10.1145/325165.325242, §2.3 "Quaternions", p. 246 (addition: scalar to scalar, vector to vector)
// Tags:    quaternion, add, linear-algebra
// Units:   dimensionless
// Example: (1, 2, 3, 4) + (5, 6, 7, 8) → (6, 8, 10, 12)
func QuatAdd(p, q Quat) (Quat, error) {
	if !finiteQ(p) {
		return Quat{}, errNotFinite("p")
	}
	if !finiteQ(q) {
		return Quat{}, errNotFinite("q")
	}
	r := Quat{p.W + q.W, p.X + q.X, p.Y + q.Y, p.Z + q.Z}
	if !finiteQ(r) {
		return Quat{}, errOverflow("p + q")
	}
	return r, nil
}

// QuatSub returns p − q, component by component.
// Errors: as QuatAdd.
//
// ID:      NUM-LINALG-102
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Shoemake, "Animating Rotation with Quaternion Curves", Computer Graphics (SIGGRAPH '85), 1985, 19(3):245–254, doi:10.1145/325165.325242, §2.3 "Quaternions", p. 246 (componentwise addition, with q negated)
// Tags:    quaternion, subtract, linear-algebra
// Units:   dimensionless
// Example: (5, 6, 7, 8) − (1, 2, 3, 4) → (4, 4, 4, 4)
func QuatSub(p, q Quat) (Quat, error) {
	if !finiteQ(p) {
		return Quat{}, errNotFinite("p")
	}
	if !finiteQ(q) {
		return Quat{}, errNotFinite("q")
	}
	r := Quat{p.W - q.W, p.X - q.X, p.Y - q.Y, p.Z - q.Z}
	if !finiteQ(r) {
		return Quat{}, errOverflow("p − q")
	}
	return r, nil
}

// QuatMul returns the Hamilton product p·q = [p₀q₀ − p⃗·q⃗, p₀q⃗ + q₀p⃗ + p⃗ × q⃗]. It is not commutative.
// For unit quaternions, the rotation of p·q applies q first, then p.
// Errors: ErrNotFinite for a NaN or ±Inf component; ErrOverflow when a component of the product is
// outside the float64 range.
//
// ID:      NUM-LINALG-103
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Shoemake, "Animating Rotation with Quaternion Curves", Computer Graphics (SIGGRAPH '85), 1985, 19(3):245–254, doi:10.1145/325165.325242, §2.3 "Quaternions", p. 246 (the general product rule)
// Tags:    quaternion, multiply, hamilton-product, composition, linear-algebra
// Units:   dimensionless
// Example: i·j → k; j·i → −k (hand-checked)
func QuatMul(p, q Quat) (Quat, error) {
	if !finiteQ(p) {
		return Quat{}, errNotFinite("p")
	}
	if !finiteQ(q) {
		return Quat{}, errNotFinite("q")
	}
	r := mulQ(p, q)
	if !finiteQ(r) {
		return Quat{}, errOverflow("p·q")
	}
	return r, nil
}

// QuatConjugate returns q* = [q₀, −q⃗]. For a unit quaternion it is the inverse rotation.
// Errors: ErrNotFinite for a NaN or ±Inf component.
//
// ID:      NUM-LINALG-104
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §6.1, eq. (98)
// Tags:    quaternion, conjugate, linear-algebra
// Units:   dimensionless
// Example: (1, 2, 3, 4) → (1, −2, −3, −4)
func QuatConjugate(q Quat) (Quat, error) {
	if !finiteQ(q) {
		return Quat{}, errNotFinite("q")
	}
	return Quat{q.W, -q.X, -q.Y, -q.Z}, nil
}

// QuatNorm returns the magnitude |q| = √(W² + X² + Y² + Z²), without overflow or underflow in the
// squares. (Diebel's eq. 99 calls the squared value the norm; this returns the square root.)
// Errors: ErrNotFinite for a NaN or ±Inf component; ErrOverflow when |q| exceeds the float64 range.
//
// ID:      NUM-LINALG-105
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §6.1, eq. (99)
// Source:  Blue, "A Portable Fortran Program to Find the Euclidean Norm of a Vector", ACM Transactions on Mathematical Software, 1978, 4(1):15–23, doi:10.1145/355769.355771
// Tags:    quaternion, norm, magnitude, linear-algebra
// Units:   dimensionless
// Example: (1, 2, 2, 4) → 5 (hand-checked)
func QuatNorm(q Quat) (float64, error) {
	if !finiteQ(q) {
		return 0, errNotFinite("q")
	}
	n := norm4(q)
	if !reduce.AllFinite(n) {
		return 0, errOverflow("|q|")
	}
	return n, nil
}

// QuatNormalize returns the versor (unit quaternion) q/|q|, which represents the same rotation as q.
// Errors: ErrZeroVector when q = 0; ErrNotFinite for a NaN or ±Inf component.
//
// ID:      NUM-LINALG-106
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §6.4, eq. (120) (unit quaternions)
// Tags:    quaternion, normalize, versor, unit-quaternion, linear-algebra
// Units:   dimensionless
// Example: (1, 2, 2, 4) → (0.2, 0.4, 0.4, 0.8) (hand-checked)
func QuatNormalize(q Quat) (Quat, error) {
	if !finiteQ(q) {
		return Quat{}, errNotFinite("q")
	}
	u, ok := unitQ(q)
	if !ok {
		return Quat{}, fmt.Errorf("%w: q", ErrZeroVector)
	}
	return u, nil
}

// QuatInverse returns q⁻¹ = q*/|q|², so that q·q⁻¹ = q⁻¹·q = 1.
// Errors: ErrZeroVector when q = 0; ErrNotFinite for a NaN or ±Inf component; ErrOverflow when a
// component of the inverse is outside the float64 range (q tiny).
//
// ID:      NUM-LINALG-107
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Shoemake, "Animating Rotation with Quaternion Curves", Computer Graphics (SIGGRAPH '85), 1985, 19(3):245–254, doi:10.1145/325165.325242, §3.2 "How quaternions rotate", pp. 247–248 (negate the vector part and divide by the magnitude squared)
// Source:  Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §6.1, eq. (100)
// Tags:    quaternion, inverse, reciprocal, linear-algebra
// Units:   dimensionless
// Example: (1, 1, 1, 1) → (0.25, −0.25, −0.25, −0.25) (hand-checked)
func QuatInverse(q Quat) (Quat, error) {
	if !finiteQ(q) {
		return Quat{}, errNotFinite("q")
	}
	u, ok := unitQ(q)
	if !ok {
		return Quat{}, fmt.Errorf("%w: q", ErrZeroVector)
	}
	// |q| = ns·2^e with ns in [1, 2]; divide by ns, then scale by 2^(−e), so |q| itself never overflows.
	m := max(math.Abs(q.W), math.Abs(q.X), math.Abs(q.Y), math.Abs(q.Z))
	_, e := pow2Scale(m)
	ns := norm4(Quat{math.Ldexp(q.W, -e), math.Ldexp(q.X, -e), math.Ldexp(q.Y, -e), math.Ldexp(q.Z, -e)})
	r := Quat{math.Ldexp(u.W/ns, -e), math.Ldexp(-u.X/ns, -e), math.Ldexp(-u.Y/ns, -e), math.Ldexp(-u.Z/ns, -e)}
	if !finiteQ(r) {
		return Quat{}, errOverflow("q⁻¹")
	}
	return r, nil
}
