package linalg

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Euler angles here are the aerospace yaw-pitch-roll set (Tait-Bryan, sequence Z-Y-X): the body turns by
// yaw ψ about z, then by pitch θ about the new y, then by roll φ about the newest x. The active rotation
// matrix is Rz(ψ)·Ry(θ)·Rx(φ). Diebel 2006 calls this the (1, 2, 3) sequence because he writes the
// frame rotation R₁(φ)R₂(θ)R₃(ψ), the transpose; the quaternion components are the same.

// QuatFromEulerZYX returns the unit quaternion of the yaw-pitch-roll rotation:
// W = cφcθcψ + sφsθsψ, X = sφcθcψ − cφsθsψ, Y = cφsθcψ + sφcθsψ, Z = cφcθsψ − sφsθcψ, where c and s are
// the cosine and sine of half of each angle.
// Errors: ErrNotFinite for a NaN or ±Inf angle.
//
// ID:      NUM-LINALG-114
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §5.6.5, eq. (84) (sequence (1, 2, 3): φ roll, θ pitch, ψ yaw)
// Tags:    quaternion, euler-angles, yaw, pitch, roll, tait-bryan, aerospace, conversion, linear-algebra
// Units:   angles in radians; the result is dimensionless
// Example: yaw π/2, pitch 0, roll 0 → (√2/2, 0, 0, √2/2) (hand-checked)
func QuatFromEulerZYX(yawRad, pitchRad, rollRad float64) (Quat, error) {
	if !reduce.AllFinite(yawRad) {
		return Quat{}, errNotFinite("yawRad")
	}
	if !reduce.AllFinite(pitchRad) {
		return Quat{}, errNotFinite("pitchRad")
	}
	if !reduce.AllFinite(rollRad) {
		return Quat{}, errNotFinite("rollRad")
	}
	sp, cp := math.Sincos(rollRad / 2)
	st, ct := math.Sincos(pitchRad / 2)
	ss, cs := math.Sincos(yawRad / 2)
	return Quat{
		W: cp*ct*cs + sp*st*ss,
		X: sp*ct*cs - cp*st*ss,
		Y: cp*st*cs + sp*ct*ss,
		Z: cp*ct*ss - sp*st*cs,
	}, nil
}

// QuatToEulerZYX returns the yaw ψ ∈ (−π, π], pitch θ ∈ [−π/2, π/2] and roll φ ∈ (−π, π] of the rotation
// q (normalised first), by Diebel's formulas: φ = atan2(2(WX + YZ), W² − X² − Y² + Z²),
// θ = asin(2(WY − XZ)), ψ = atan2(2(XY + WZ), W² + X² − Y² − Z²).
// Gimbal lock: at pitch ±π/2 roll and yaw turn about the same axis and only their difference (θ = +π/2)
// or sum (θ = −π/2) is defined. When 1 − |2(WY − XZ)| ≤ lockTol the pitch is set to ±π/2, roll to 0, and
// yaw to 2·atan2(Z, W), which reproduces the rotation exactly. lockTol < 0 selects the default 2⁻⁵⁰
// (four rounding units of the sine of the pitch). Near the lock, roll and yaw separately lose accuracy as
// the lock is approached (Diebel §5.6.6); the rotation they describe does not.
// Errors: ErrZeroVector when q = 0; ErrNotFinite for a NaN or ±Inf input.
//
// ID:      NUM-LINALG-115
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §5.6.3, eq. (73) and §5.6.6 "Singularities" (sequence (1, 2, 3): φ roll, θ pitch, ψ yaw)
// Tags:    quaternion, euler-angles, yaw, pitch, roll, tait-bryan, aerospace, gimbal-lock, conversion, linear-algebra
// Units:   q dimensionless, lockTol dimensionless; the angles in radians
// Example: (√2/2, 0, 0, √2/2) → yaw π/2, pitch 0, roll 0 (hand-checked)
func QuatToEulerZYX(q Quat, lockTol float64) (yawRad, pitchRad, rollRad float64, err error) {
	if !finiteQ(q) {
		return 0, 0, 0, errNotFinite("q")
	}
	if !reduce.AllFinite(lockTol) {
		return 0, 0, 0, errNotFinite("lockTol")
	}
	if lockTol < 0 {
		lockTol = 0x1p-50
	}
	u, ok := unitQ(q)
	if !ok {
		return 0, 0, 0, fmt.Errorf("%w: q", ErrZeroVector)
	}
	w, x, y, z := u.W, u.X, u.Y, u.Z
	s := 2 * (w*y - x*z)
	if 1-math.Abs(s) <= lockTol {
		pitch := math.Copysign(math.Pi/2, s)
		yaw := 2 * math.Atan2(z, w)
		if yaw > math.Pi {
			yaw -= 2 * math.Pi
		} else if yaw <= -math.Pi {
			yaw += 2 * math.Pi
		}
		return yaw, pitch, 0, nil
	}
	roll := math.Atan2(2*(w*x+y*z), w*w-x*x-y*y+z*z)
	pitch := math.Asin(max(-1, min(1, s)))
	yaw := math.Atan2(2*(x*y+w*z), w*w+x*x-y*y-z*z)
	return yaw, pitch, roll, nil
}
