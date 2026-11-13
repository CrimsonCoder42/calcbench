// Package geometry computes lengths, angles, areas and volumes of plane and solid figures: points and
// lines, triangles, quadrilaterals and polygons, the circle and the ellipse, prisms, pyramids, cylinders,
// cones, the sphere, the torus, the ellipsoid, the five regular polyhedra, and the circular helix.
//
// ID short form: GEOM (IDs are NUM-GEOM-NNN).
//
// Conventions (doc/DECISIONS.md D-011): every angle is in radians and says Rad in its name. Lengths are
// plain float64 in any one unit the caller chooses; areas come out in that unit squared and volumes in
// that unit cubed (no unit types: unit/ may import num/, so the reverse would be a cycle). Points are this
// package's Point{X, Y}. Every function returns (value, error); inputs are only read. NaN or ±Inf in
// returns ErrNotFinite; finite inputs whose result is beyond float64 return ErrOverflow, and products are
// formed so that no intermediate overflows before the result does. Errors are wrapped with the parameter
// name; test with errors.Is. The two elliptic-integral helpers (the arithmetic-geometric mean for E(k),
// Carlson's R_D and R_G) are unexported here until a second package needs them.
//
// The law of cosines is in num/trig: the third side from two sides and the included angle is
// NUM-TRIG-020 LawOfCosinesSide, and an angle from three sides is NUM-TRIG-021 LawOfCosinesAngle. The IDs
// NUM-GEOM-018 and NUM-GEOM-019, reserved for them here, are not used. The golden ratio is NUM-LOGS-001.
//
// Functions:
//
//	NUM-GEOM-001  PointDistance                       distance between two points
//	NUM-GEOM-002  Midpoint                            midpoint of a segment
//	NUM-GEOM-003  Slope                               rise over run of the line through two points
//	NUM-GEOM-004  SlopeAngleRad                       angle of a direction from rise and run
//	NUM-GEOM-005  PointLineDistance                   perpendicular distance from a point to a line
//	NUM-GEOM-006  AngleBetweenLinesRad                acute angle between two lines
//	NUM-GEOM-007  PolarToCartesian                    (r, θ) to (x, y)
//	NUM-GEOM-008  CartesianToPolar                    (x, y) to (r, θ)
//	NUM-GEOM-009  ComplementaryAngleRad               π/2 − θ
//	NUM-GEOM-010  SupplementaryAngleRad               π − θ
//	NUM-GEOM-011  PolygonInteriorAngleSumRad          sum of the interior angles of an n-gon
//	NUM-GEOM-012  RegularPolygonInteriorAngleRad      one interior angle of a regular n-gon
//	NUM-GEOM-013  RegularPolygonExteriorAngleRad      one exterior angle of a regular n-gon
//	NUM-GEOM-014  TrianglePerimeter                   perimeter of a triangle from three sides
//	NUM-GEOM-015  TriangleAreaBaseHeight              area of a triangle from base and height
//	NUM-GEOM-016  TriangleAreaHeron                   area of a triangle from three sides (Kahan's Heron)
//	NUM-GEOM-017  TriangleAreaSAS                     area of a triangle from two sides and the included angle
//	NUM-GEOM-020  RightTriangleHypotenuse             hypotenuse of a right triangle from two legs
//	NUM-GEOM-021  RightTriangleLeg                    other leg of a right triangle from hypotenuse and leg
//	NUM-GEOM-022  RectangleArea                       area of a rectangle (square; parallelogram by base × height)
//	NUM-GEOM-023  RectanglePerimeter                  perimeter of a quadrilateral with two pairs of equal sides
//	NUM-GEOM-024  ParallelogramAreaSAS                area of a parallelogram from two sides and the angle
//	NUM-GEOM-025  TrapezoidArea                       area of a trapezoid
//	NUM-GEOM-026  KiteArea                            area of a quadrilateral with perpendicular diagonals
//	NUM-GEOM-027  RegularPolygonAreaFromSide          area of a regular n-gon from its side
//	NUM-GEOM-028  RegularPolygonAreaFromApothem       area of a regular n-gon from its apothem
//	NUM-GEOM-029  RegularPolygonAreaFromCircumradius  area of a regular n-gon from its circumradius
//	NUM-GEOM-030  RegularPolygonSideInscribed         side of a regular n-gon inscribed in a circle
//	NUM-GEOM-031  RegularPolygonSideCircumscribed     side of a regular n-gon circumscribed about a circle
//	NUM-GEOM-032  PolygonArea                         area of a simple polygon from its vertices (shoelace)
//	NUM-GEOM-033  PolygonPerimeter                    perimeter of a polygon from its vertices
//	NUM-GEOM-034  CircleCircumference                 circumference of a circle
//	NUM-GEOM-035  CircleArea                          area of a circle
//	NUM-GEOM-036  AnnulusArea                         area of an annulus (ring)
//	NUM-GEOM-037  ArcLength                           length of a circular arc
//	NUM-GEOM-038  ChordLength                         length of a circular chord
//	NUM-GEOM-039  SectorArea                          area of a circular sector
//	NUM-GEOM-040  SegmentArea                         area of a circular segment
//	NUM-GEOM-041  ArcRadiusFromChordSagitta           radius of an arc from its chord and sagitta
//	NUM-GEOM-042  CircleTangentLength                 tangent length from an external point to a circle
//	NUM-GEOM-043  EllipseArea                         area of an ellipse
//	NUM-GEOM-044  EllipseChordFromCenter              chord of an ellipse at an offset from the centre
//	NUM-GEOM-045  EllipseChordFromEdge                chord of an ellipse at a depth from the vertex
//	NUM-GEOM-046  EllipsePerimeter                    perimeter of an ellipse, exact (AGM)
//	NUM-GEOM-047  EllipsePerimeterRamanujan           perimeter of an ellipse, Ramanujan's approximation
//	NUM-GEOM-048  PrismVolume                         volume of a prism from base area and height
//	NUM-GEOM-049  CylinderVolume                      volume of a cylinder
//	NUM-GEOM-050  CylinderLateralArea                 lateral area of a cylinder
//	NUM-GEOM-051  CylinderSurfaceArea                 total surface area of a cylinder
//	NUM-GEOM-052  HollowCylinderVolume                volume of a hollow cylinder (tube)
//	NUM-GEOM-053  HollowCylinderSurfaceArea           total surface area of a hollow cylinder (tube)
//	NUM-GEOM-054  PyramidVolume                       volume of a pyramid, any base
//	NUM-GEOM-055  RegularPyramidSurfaceArea           total surface area of a right pyramid with a regular base
//	NUM-GEOM-056  PyramidFrustumVolume                volume of a frustum of a pyramid
//	NUM-GEOM-057  PrismoidVolume                      volume by the prismoidal formula
//	NUM-GEOM-058  ConeVolume                          volume of a cone
//	NUM-GEOM-059  ConeLateralArea                     lateral area of a right circular cone
//	NUM-GEOM-060  ConeSurfaceArea                     total surface area of a right circular cone
//	NUM-GEOM-061  ConeFrustumVolume                   volume of a frustum of a cone
//	NUM-GEOM-062  ConeFrustumSurfaceArea              total surface area of a frustum of a cone
//	NUM-GEOM-063  SphereVolume                        volume of a sphere
//	NUM-GEOM-064  SphereSurfaceArea                   surface area of a sphere
//	NUM-GEOM-065  SphericalCapVolume                  volume of a spherical cap
//	NUM-GEOM-066  SphericalCapArea                    curved area of a spherical cap
//	NUM-GEOM-067  TorusVolume                         volume of a ring torus
//	NUM-GEOM-068  TorusSurfaceArea                    surface area of a ring torus
//	NUM-GEOM-069  EllipsoidVolume                     volume of an ellipsoid
//	NUM-GEOM-070  EllipsoidSurfaceArea                surface area of an ellipsoid, exact (Carlson)
//	NUM-GEOM-071  TetrahedronVolume                   volume of a regular tetrahedron
//	NUM-GEOM-072  TetrahedronSurfaceArea              surface area of a regular tetrahedron
//	NUM-GEOM-073  CubeVolume                          volume of a cube
//	NUM-GEOM-074  CubeSurfaceArea                     surface area of a cube
//	NUM-GEOM-075  OctahedronVolume                    volume of a regular octahedron
//	NUM-GEOM-076  OctahedronSurfaceArea               surface area of a regular octahedron
//	NUM-GEOM-077  DodecahedronVolume                  volume of a regular dodecahedron
//	NUM-GEOM-078  DodecahedronSurfaceArea             surface area of a regular dodecahedron
//	NUM-GEOM-079  IcosahedronVolume                   volume of a regular icosahedron
//	NUM-GEOM-080  IcosahedronSurfaceArea              surface area of a regular icosahedron
//	NUM-GEOM-081  HelixArcLength                      arc length of a circular helix
//	NUM-GEOM-082  HelixHeight                         height of a circular helix from its arc length
package geometry
