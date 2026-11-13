package physics

// AverageSpeed returns the average speed over a trip, v = d / t (total distance over elapsed time), in m/s.
// Errors: ErrNotFinite; ErrNegative for a negative distance or duration; ErrZero for durationSeconds = 0;
// ErrOverflow.
//
// ID:      SCI-PHYS-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ling SJ, Sanny J, Moebs W. University Physics Volume 1. OpenStax, 2016, https://openstax.org/details/books/university-physics-volume-1, §3.1 "Position, Displacement, and Average Velocity" (average speed = total distance / elapsed time)
// Tags:    speed, average-speed, distance, time, kinematics, physics
// Units:   distanceMeters m; durationSeconds s; result m/s
// Example: 100 000 m in 3600 s → 27.777… m/s (100 km/h)
func AverageSpeed(distanceMeters, durationSeconds float64) (float64, error) {
	return quotient(distanceMeters, "distanceMeters", durationSeconds, "durationSeconds")
}

// DistanceFromSpeed returns the distance covered at an average speed over a duration, d = v · t, in meters.
// Errors: ErrNotFinite; ErrNegative for a negative speed or duration; ErrOverflow.
//
// ID:      SCI-PHYS-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ling SJ, Sanny J, Moebs W. University Physics Volume 1. OpenStax, 2016, https://openstax.org/details/books/university-physics-volume-1, §3.1 "Position, Displacement, and Average Velocity" (average speed = total distance / elapsed time, solved for distance)
// Tags:    speed, distance, time, kinematics, physics
// Units:   speedMetersPerSecond m/s; durationSeconds s; result m
// Example: 25 m/s for 7200 s → 180 000 m
func DistanceFromSpeed(speedMetersPerSecond, durationSeconds float64) (float64, error) {
	return product(speedMetersPerSecond, "speedMetersPerSecond", durationSeconds, "durationSeconds")
}

// DurationFromSpeed returns the time to cover a distance at an average speed, t = d / v, in seconds.
// Errors: ErrNotFinite; ErrNegative for a negative distance or speed; ErrZero for speed 0; ErrOverflow.
//
// ID:      SCI-PHYS-006
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ling SJ, Sanny J, Moebs W. University Physics Volume 1. OpenStax, 2016, https://openstax.org/details/books/university-physics-volume-1, §3.1 "Position, Displacement, and Average Velocity" (average speed = total distance / elapsed time, solved for time)
// Tags:    speed, duration, distance, time, kinematics, physics
// Units:   distanceMeters m; speedMetersPerSecond m/s; result s
// Example: 180 000 m at 25 m/s → 7200 s
func DurationFromSpeed(distanceMeters, speedMetersPerSecond float64) (float64, error) {
	return quotient(distanceMeters, "distanceMeters", speedMetersPerSecond, "speedMetersPerSecond")
}
