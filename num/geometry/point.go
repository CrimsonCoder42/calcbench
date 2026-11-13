package geometry

// Point is a point of the Cartesian plane. X and Y are in one length unit, chosen by the caller; every
// length a function derives from points is in that unit.
type Point struct {
	X, Y float64
}
