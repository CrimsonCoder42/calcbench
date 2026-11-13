# S7: Geometry (`num/geometry`)

Status: SHIPPED 2026-11-12
Work folder: `doc/math_build_2026-11-12/` (hand-off: `notes_geometry.md`)
Decisions: D-011 and Devin's request of 2026-11-12 (radians only; package-local `Point`; plain `float64` lengths;
"Interface Angle" = angle between two lines; golden ratio is NUM-LOGS-001; the law of cosines is num/trig, so
NUM-GEOM-018 and -019 are not built; "Bell-shaped" gets no ID) · Assumptions: none

Detail for every entry (formula, inputs, variants, notes) is in [../catalog/geometry.md](../catalog/geometry.md),
STEP 2. This spec records what was built, the sources as confirmed during the build, and the choices made.

Conventions: every angle is in radians with `Rad` in its name; lengths are `float64` in one caller-chosen unit,
areas in that unit squared, volumes cubed (stated under `Units:`). Every function returns `(value, error)` and only
reads its inputs. Sentinels, each wrapped with the parameter name: `ErrNotFinite` (NaN/±Inf in), `ErrOverflow`
(finite inputs, result beyond float64), `ErrNonPositive`, `ErrNegative`, `ErrOutOfRange`, `ErrNotTriangle`,
`ErrDegenerate`, `ErrTooFewSides`, `ErrVertical`, `ErrNoConvergence` (an iteration cap; not reachable from valid
input). Imports: standard library only (`errors`, `fmt`, `math`); no other calcbench package is needed.

## Choices that apply to many functions

- **No spurious overflow.** Products of three or four factors go through an unexported `product` that multiplies
  the mantissas and adds the exponents (`math.Frexp`/`math.Ldexp`), so a result is ±Inf only when the true value
  is beyond float64 and no intermediate underflows first. Lost: plain left-to-right products (overflow or
  underflow in an intermediate for e.g. `πr²h` with r = 1e200, h = 1e-300) and log-space evaluation (an extra
  `exp`/`log` and lost accuracy).
- **Differences of squares are factored**: √(c² − a²) = √(c − a)·√(c + a) (right-triangle leg, tangent length,
  ellipse chord, helix height), π(R² − r²) = 2π(R − r)(R/2 + r/2) (annulus, tube). Lost: the textbook forms,
  which cancel for nearly equal arguments and overflow for large ones.
- **Sums that can overflow are halved first** (`a/2 + b/2`): midpoint, trapezoid, Ramanujan's a + b.
- **Polygons are translated by vertex 0 and scaled by a power of two** before the shoelace sum, and the sum is
  compensated (Neumaier 1974). Lost: the plain shoelace sum (loses every digit for a polygon far from the origin:
  a unit square at 1e15 gives garbage; here it gives exactly 1) and pairwise summation (no better than Neumaier
  here and needs a buffer or recursion). Best known: Θ(n) time (reading the input), O(1) extra space, met.
- **Kahan's triangle**: TriangleAreaHeron sorts the sides and uses ¼√((a+(b+c))(c−(a−b))(c+(a−b))(a+(b−c))),
  accurate to a few ulps for needles. Lost: √(s(s−a)(s−b)(s−c)) (Kahan 2014 Table 1: 17.6 for 9.999999990).
- **Elliptic integrals stay unexported** (catalog; PENDING_SHARED 3b). EllipsePerimeter: E(k) by Carlson's AGM
  (Carlson 1995 eqs. (42)–(46)), 4a for b/a < 2⁻³⁰ where E(k) rounds to 1 (DLMF 19.12.2). Lost: Adlaj's modified
  AGM (its citation was not confirmed) and the R_G route (slower, no more accurate for b/a ≥ 2⁻³⁰).
  EllipsoidSurfaceArea: S = 4πab·R_G(c²/a², c²/b², 1) with R_G from three non-negative R_D terms (DLMF 19.21.11),
  R_D by Carlson's duplication. Lost: DLMF 19.21.10 (z·R_F − ⅓(x−z)(y−z)R_D + √(xy/z)), whose terms cancel for a
  flat ellipsoid; the Legendre form 19.33.2 (needs incomplete F and E); Thomsen's approximation (no source).
- **SegmentArea** uses the sine series for θ < 0.5 (θ − sin θ cancels there). The switch point 0.5 is the
  catalog's; the verify phase may move it (at 0.5 the direct form loses about 5 bits; at 1.0 about 2.5).

## Source confirmation (summary; details in notes_geometry.md)

The catalog cited the CRC **32nd** edition with inferred section numbers (†) and Bronshtein et al. with snippet
locators (‡). Neither could be seen. The CRC **31st** edition (2003) was read page by page (formulas seen on the
rendered pages; book page = PDF page − 8, checked against three in-text page references), so every CRC citation is
now to the 31st edition with section, equation and page as printed. All ‡ (Bronshtein) lines were dropped: every
function they covered has a confirmed source. Added sources: OpenStax Prealgebra 2e §9.3 (complementary and
supplementary angles; not in the CRC), Euclid III.36 (tangent length), Encyclopædia Britannica 1911
"Mensuration" §58 (prismoidal formula; not in the CRC), DLMF 19.12.2 and 19.21.11, Carlson 1995 (read in the
arXiv preprint). Held back: none.

## Functions

| ID | Name | Signature | Units |
|---|---|---|---|
| `NUM-GEOM-001` | `PointDistance` | `func PointDistance(p, q Point) (float64, error)` | coordinates of p and q in one length unit; result in that unit |
| `NUM-GEOM-002` | `Midpoint` | `func Midpoint(p, q Point) (Point, error)` | coordinates of p and q in one length unit; result in that unit |
| `NUM-GEOM-003` | `Slope` | `func Slope(p, q Point) (float64, error)` | coordinates of p and q in one length unit; result dimensionless |
| `NUM-GEOM-004` | `SlopeAngleRad` | `func SlopeAngleRad(rise, run float64) (float64, error)` | rise and run in one length unit (or a ratio with run = 1); result in radians |
| `NUM-GEOM-005` | `PointLineDistance` | `func PointLineDistance(p, a, b Point) (float64, error)` | coordinates of p, a, b in one length unit; result in that unit |
| `NUM-GEOM-006` | `AngleBetweenLinesRad` | `func AngleBetweenLinesRad(a, b, c, d Point) (float64, error)` | coordinates in one length unit; result in radians |
| `NUM-GEOM-007` | `PolarToCartesian` | `func PolarToCartesian(r, thetaRad float64) (Point, error)` | r in a length unit; thetaRad in radians; result coordinates in r's unit |
| `NUM-GEOM-008` | `CartesianToPolar` | `func CartesianToPolar(p Point) (r, thetaRad float64, err error)` | coordinates of p in one length unit; r in that unit; thetaRad in radians |
| `NUM-GEOM-009` | `ComplementaryAngleRad` | `func ComplementaryAngleRad(angleRad float64) (float64, error)` | angleRad and the result in radians |
| `NUM-GEOM-010` | `SupplementaryAngleRad` | `func SupplementaryAngleRad(angleRad float64) (float64, error)` | angleRad and the result in radians |
| `NUM-GEOM-011` | `PolygonInteriorAngleSumRad` | `func PolygonInteriorAngleSumRad(sides int) (float64, error)` | sides is a count; result in radians |
| `NUM-GEOM-012` | `RegularPolygonInteriorAngleRad` | `func RegularPolygonInteriorAngleRad(sides int) (float64, error)` | sides is a count; result in radians |
| `NUM-GEOM-013` | `RegularPolygonExteriorAngleRad` | `func RegularPolygonExteriorAngleRad(sides int) (float64, error)` | sides is a count; result in radians |
| `NUM-GEOM-014` | `TrianglePerimeter` | `func TrianglePerimeter(sideA, sideB, sideC float64) (float64, error)` | sides in one length unit; result in that unit |
| `NUM-GEOM-015` | `TriangleAreaBaseHeight` | `func TriangleAreaBaseHeight(base, height float64) (float64, error)` | base and height in one length unit; result in that unit squared |
| `NUM-GEOM-016` | `TriangleAreaHeron` | `func TriangleAreaHeron(sideA, sideB, sideC float64) (float64, error)` | sides in one length unit; result in that unit squared |
| `NUM-GEOM-017` | `TriangleAreaSAS` | `func TriangleAreaSAS(sideA, sideB, includedAngleRad float64) (float64, error)` | sides in one length unit; includedAngleRad in radians; result in that unit squared |
| `NUM-GEOM-020` | `RightTriangleHypotenuse` | `func RightTriangleHypotenuse(legA, legB float64) (float64, error)` | legs in one length unit; result in that unit |
| `NUM-GEOM-021` | `RightTriangleLeg` | `func RightTriangleLeg(hypotenuse, leg float64) (float64, error)` | hypotenuse and leg in one length unit; result in that unit |
| `NUM-GEOM-022` | `RectangleArea` | `func RectangleArea(width, height float64) (float64, error)` | width and height in one length unit; result in that unit squared |
| `NUM-GEOM-023` | `RectanglePerimeter` | `func RectanglePerimeter(sideA, sideB float64) (float64, error)` | sides in one length unit; result in that unit |
| `NUM-GEOM-024` | `ParallelogramAreaSAS` | `func ParallelogramAreaSAS(sideA, sideB, angleRad float64) (float64, error)` | sides in one length unit; angleRad in radians; result in that unit squared |
| `NUM-GEOM-025` | `TrapezoidArea` | `func TrapezoidArea(baseA, baseB, height float64) (float64, error)` | baseA, baseB and height in one length unit; result in that unit squared |
| `NUM-GEOM-026` | `KiteArea` | `func KiteArea(diagonal1, diagonal2 float64) (float64, error)` | diagonals in one length unit; result in that unit squared |
| `NUM-GEOM-027` | `RegularPolygonAreaFromSide` | `func RegularPolygonAreaFromSide(sides int, side float64) (float64, error)` | sides is a count; side in a length unit; result in that unit squared |
| `NUM-GEOM-028` | `RegularPolygonAreaFromApothem` | `func RegularPolygonAreaFromApothem(sides int, apothem float64) (float64, error)` | sides is a count; apothem in a length unit; result in that unit squared |
| `NUM-GEOM-029` | `RegularPolygonAreaFromCircumradius` | `func RegularPolygonAreaFromCircumradius(sides int, circumradius float64) (float64, error)` | sides is a count; circumradius in a length unit; result in that unit squared |
| `NUM-GEOM-030` | `RegularPolygonSideInscribed` | `func RegularPolygonSideInscribed(sides int, circumradius float64) (float64, error)` | sides is a count; circumradius in a length unit; result in that unit |
| `NUM-GEOM-031` | `RegularPolygonSideCircumscribed` | `func RegularPolygonSideCircumscribed(sides int, inradius float64) (float64, error)` | sides is a count; inradius in a length unit; result in that unit |
| `NUM-GEOM-032` | `PolygonArea` | `func PolygonArea(vertices []Point) (float64, error)` | coordinates in one length unit; result in that unit squared |
| `NUM-GEOM-033` | `PolygonPerimeter` | `func PolygonPerimeter(vertices []Point) (float64, error)` | coordinates in one length unit; result in that unit |
| `NUM-GEOM-034` | `CircleCircumference` | `func CircleCircumference(radius float64) (float64, error)` | radius in a length unit; result in that unit |
| `NUM-GEOM-035` | `CircleArea` | `func CircleArea(radius float64) (float64, error)` | radius in a length unit; result in that unit squared |
| `NUM-GEOM-036` | `AnnulusArea` | `func AnnulusArea(outerRadius, innerRadius float64) (float64, error)` | radii in one length unit; result in that unit squared |
| `NUM-GEOM-037` | `ArcLength` | `func ArcLength(radius, angleRad float64) (float64, error)` | radius in a length unit; angleRad in radians; result in radius's unit |
| `NUM-GEOM-038` | `ChordLength` | `func ChordLength(radius, centralAngleRad float64) (float64, error)` | radius in a length unit; centralAngleRad in radians; result in radius's unit |
| `NUM-GEOM-039` | `SectorArea` | `func SectorArea(radius, angleRad float64) (float64, error)` | radius in a length unit; angleRad in radians; result in that unit squared |
| `NUM-GEOM-040` | `SegmentArea` | `func SegmentArea(radius, centralAngleRad float64) (float64, error)` | radius in a length unit; centralAngleRad in radians; result in that unit squared |
| `NUM-GEOM-041` | `ArcRadiusFromChordSagitta` | `func ArcRadiusFromChordSagitta(chord, sagitta float64) (float64, error)` | chord and sagitta in one length unit; result in that unit |
| `NUM-GEOM-042` | `CircleTangentLength` | `func CircleTangentLength(radius, distanceFromCenter float64) (float64, error)` | radius and distanceFromCenter in one length unit; result in that unit |
| `NUM-GEOM-043` | `EllipseArea` | `func EllipseArea(semiAxisA, semiAxisB float64) (float64, error)` | semi-axes in one length unit; result in that unit squared |
| `NUM-GEOM-044` | `EllipseChordFromCenter` | `func EllipseChordFromCenter(semiAxisAlong, semiAxisAcross, offsetFromCenter float64) (float64, error)` | semi-axes and offsetFromCenter in one length unit; result in that unit |
| `NUM-GEOM-045` | `EllipseChordFromEdge` | `func EllipseChordFromEdge(semiAxisAlong, semiAxisAcross, depthFromEdge float64) (float64, error)` | semi-axes and depthFromEdge in one length unit; result in that unit |
| `NUM-GEOM-046` | `EllipsePerimeter` | `func EllipsePerimeter(semiAxisA, semiAxisB float64) (float64, error)` | semi-axes in one length unit; result in that unit |
| `NUM-GEOM-047` | `EllipsePerimeterRamanujan` | `func EllipsePerimeterRamanujan(semiAxisA, semiAxisB float64) (float64, error)` | semi-axes in one length unit; result in that unit |
| `NUM-GEOM-048` | `PrismVolume` | `func PrismVolume(baseArea, height float64) (float64, error)` | baseArea in a length unit squared; height in that unit; result in that unit cubed |
| `NUM-GEOM-049` | `CylinderVolume` | `func CylinderVolume(radius, height float64) (float64, error)` | radius and height in one length unit; result in that unit cubed |
| `NUM-GEOM-050` | `CylinderLateralArea` | `func CylinderLateralArea(radius, height float64) (float64, error)` | radius and height in one length unit; result in that unit squared |
| `NUM-GEOM-051` | `CylinderSurfaceArea` | `func CylinderSurfaceArea(radius, height float64) (float64, error)` | radius and height in one length unit; result in that unit squared |
| `NUM-GEOM-052` | `HollowCylinderVolume` | `func HollowCylinderVolume(outerRadius, innerRadius, height float64) (float64, error)` | radii and height in one length unit; result in that unit cubed |
| `NUM-GEOM-053` | `HollowCylinderSurfaceArea` | `func HollowCylinderSurfaceArea(outerRadius, innerRadius, height float64) (float64, error)` | radii and height in one length unit; result in that unit squared |
| `NUM-GEOM-054` | `PyramidVolume` | `func PyramidVolume(baseArea, height float64) (float64, error)` | baseArea in a length unit squared; height in that unit; result in that unit cubed |
| `NUM-GEOM-055` | `RegularPyramidSurfaceArea` | `func RegularPyramidSurfaceArea(sides int, side, height float64) (float64, error)` | sides is a count; side and height in one length unit; result in that unit squared |
| `NUM-GEOM-056` | `PyramidFrustumVolume` | `func PyramidFrustumVolume(baseArea1, baseArea2, height float64) (float64, error)` | base areas in a length unit squared; height in that unit; result in that unit cubed |
| `NUM-GEOM-057` | `PrismoidVolume` | `func PrismoidVolume(area1, areaMid, area2, height float64) (float64, error)` | areas in a length unit squared; height in that unit; result in that unit cubed |
| `NUM-GEOM-058` | `ConeVolume` | `func ConeVolume(radius, height float64) (float64, error)` | radius and height in one length unit; result in that unit cubed |
| `NUM-GEOM-059` | `ConeLateralArea` | `func ConeLateralArea(radius, height float64) (float64, error)` | radius and height in one length unit; result in that unit squared |
| `NUM-GEOM-060` | `ConeSurfaceArea` | `func ConeSurfaceArea(radius, height float64) (float64, error)` | radius and height in one length unit; result in that unit squared |
| `NUM-GEOM-061` | `ConeFrustumVolume` | `func ConeFrustumVolume(radius1, radius2, height float64) (float64, error)` | radii and height in one length unit; result in that unit cubed |
| `NUM-GEOM-062` | `ConeFrustumSurfaceArea` | `func ConeFrustumSurfaceArea(radius1, radius2, height float64) (float64, error)` | radii and height in one length unit; result in that unit squared |
| `NUM-GEOM-063` | `SphereVolume` | `func SphereVolume(radius float64) (float64, error)` | radius in a length unit; result in that unit cubed |
| `NUM-GEOM-064` | `SphereSurfaceArea` | `func SphereSurfaceArea(radius float64) (float64, error)` | radius in a length unit; result in that unit squared |
| `NUM-GEOM-065` | `SphericalCapVolume` | `func SphericalCapVolume(sphereRadius, capHeight float64) (float64, error)` | sphereRadius and capHeight in one length unit; result in that unit cubed |
| `NUM-GEOM-066` | `SphericalCapArea` | `func SphericalCapArea(sphereRadius, capHeight float64) (float64, error)` | sphereRadius and capHeight in one length unit; result in that unit squared |
| `NUM-GEOM-067` | `TorusVolume` | `func TorusVolume(majorRadius, minorRadius float64) (float64, error)` | majorRadius and minorRadius in one length unit; result in that unit cubed |
| `NUM-GEOM-068` | `TorusSurfaceArea` | `func TorusSurfaceArea(majorRadius, minorRadius float64) (float64, error)` | majorRadius and minorRadius in one length unit; result in that unit squared |
| `NUM-GEOM-069` | `EllipsoidVolume` | `func EllipsoidVolume(semiAxisA, semiAxisB, semiAxisC float64) (float64, error)` | semi-axes in one length unit; result in that unit cubed |
| `NUM-GEOM-070` | `EllipsoidSurfaceArea` | `func EllipsoidSurfaceArea(semiAxisA, semiAxisB, semiAxisC float64) (float64, error)` | semi-axes in one length unit; result in that unit squared |
| `NUM-GEOM-071` | `TetrahedronVolume` | `func TetrahedronVolume(edge float64) (float64, error)` | edge in a length unit; result in that unit cubed |
| `NUM-GEOM-072` | `TetrahedronSurfaceArea` | `func TetrahedronSurfaceArea(edge float64) (float64, error)` | edge in a length unit; result in that unit squared |
| `NUM-GEOM-073` | `CubeVolume` | `func CubeVolume(edge float64) (float64, error)` | edge in a length unit; result in that unit cubed |
| `NUM-GEOM-074` | `CubeSurfaceArea` | `func CubeSurfaceArea(edge float64) (float64, error)` | edge in a length unit; result in that unit squared |
| `NUM-GEOM-075` | `OctahedronVolume` | `func OctahedronVolume(edge float64) (float64, error)` | edge in a length unit; result in that unit cubed |
| `NUM-GEOM-076` | `OctahedronSurfaceArea` | `func OctahedronSurfaceArea(edge float64) (float64, error)` | edge in a length unit; result in that unit squared |
| `NUM-GEOM-077` | `DodecahedronVolume` | `func DodecahedronVolume(edge float64) (float64, error)` | edge in a length unit; result in that unit cubed |
| `NUM-GEOM-078` | `DodecahedronSurfaceArea` | `func DodecahedronSurfaceArea(edge float64) (float64, error)` | edge in a length unit; result in that unit squared |
| `NUM-GEOM-079` | `IcosahedronVolume` | `func IcosahedronVolume(edge float64) (float64, error)` | edge in a length unit; result in that unit cubed |
| `NUM-GEOM-080` | `IcosahedronSurfaceArea` | `func IcosahedronSurfaceArea(edge float64) (float64, error)` | edge in a length unit; result in that unit squared |
| `NUM-GEOM-081` | `HelixArcLength` | `func HelixArcLength(radius, pitch, turns float64) (float64, error)` | radius and pitch in one length unit; turns is a count (fractional allowed); result in that unit |
| `NUM-GEOM-082` | `HelixHeight` | `func HelixHeight(radius, arcLength, turns float64) (float64, error)` | radius and arcLength in one length unit; turns is a count (fractional allowed); result in that unit |

## Per function

### `NUM-GEOM-001` `PointDistance`  (point_distance.go)
- **Source**: CRC31 §4.4.2 "Distances", eq. (4.4.5), p. 316.
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: (0, 0), (3, 4) → 5 (hand-derived from eq. (4.4.5)).
- **Errors**: ErrNotFinite when a coordinate is NaN or ±Inf; ErrOverflow when the distance is beyond float64.

### `NUM-GEOM-002` `Midpoint`  (midpoint.go)
- **Source**: CRC31 §4.4.2 "Distances", the midpoint of P₀P₁ (after eq. (4.4.8)), p. 316.
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: (0, 0), (4, −2) → (2, −1) (hand-derived).
- **Errors**: ErrNotFinite when a coordinate is NaN or ±Inf.

### `NUM-GEOM-003` `Slope`  (slope.go)
- **Source**: CRC31 §4.4.1 "Lines with prescribed properties", item 6 (slope of the line through two points), p. 315.
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: (1, 1), (3, 5) → 2 (hand-derived).
- **Errors**: ErrNotFinite when a coordinate is NaN or ±Inf; ErrDegenerate when p == q; ErrVertical when x₁ == x₂ and y₁ ≠ y₂; ErrOverflow when the slope is beyond float64.

### `NUM-GEOM-004` `SlopeAngleRad`  (slope_angle.go)
- **Source**: CRC31 §4.4.3 "Angles", p. 316 (the angle of a line of slope m from the x-axis is tan⁻¹ m); CRC31 §4.1.4.1 "Relations between Cartesian and polar coordinates", p. 303 (θ from x and y, sin θ = y/√(x²+y²), cos θ = x/√(x²+y²)).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: rise 1, run 1 → π/4 (hand-derived).
- **Errors**: ErrNotFinite when rise or run is NaN or ±Inf; ErrDegenerate when rise == run == 0.

### `NUM-GEOM-005` `PointLineDistance`  (point_line_distance.go)
- **Source**: CRC31 §4.4.2 "Distances", eq. (4.4.9), p. 316 (distance from a point to a line), with the line through two points of §4.4.1, eq. (4.4.3), p. 315.
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: p (0, 0), a (0, 1), b (1, 0) → √2/2 (hand-derived).
- **Errors**: ErrNotFinite when a coordinate is NaN or ±Inf; ErrDegenerate when a == b; ErrOverflow when the distance is beyond float64.

### `NUM-GEOM-006` `AngleBetweenLinesRad`  (angle_between_lines.go)
- **Source**: CRC31 §4.4.3 "Angles", eq. (4.4.10), p. 316 (angle between two lines; the form in slopes, tan⁻¹((m₁−m₀)/(1+m₀m₁)), on the same page).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: the line y = 0 and the line through (0, 0), (√3, 1) → π/6 (hand-derived).
- **Errors**: ErrNotFinite when a coordinate is NaN or ±Inf; ErrDegenerate when a == b or c == d.

### `NUM-GEOM-007` `PolarToCartesian`  (polar.go)
- **Source**: CRC31 §4.1.4.1 "Relations between Cartesian and polar coordinates", p. 303 (x = r cos θ, y = r sin θ).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: r 2, θ π/6 → (√3, 1) (hand-derived).
- **Errors**: ErrNotFinite when r or thetaRad is NaN or ±Inf; ErrNegative when r < 0.

### `NUM-GEOM-008` `CartesianToPolar`  (polar.go)
- **Source**: CRC31 §4.1.4.1 "Relations between Cartesian and polar coordinates", p. 303 (r = √(x²+y²), θ = tan⁻¹(y/x)).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: (−1, 0) → (1, π) (hand-derived).
- **Errors**: ErrNotFinite when a coordinate is NaN or ±Inf; ErrOverflow when r is beyond float64.

### `NUM-GEOM-009` `ComplementaryAngleRad`  (complementary_angle.go)
- **Source**: Marecek, Anthony-Smith, Honeycutt Mathis, Prealgebra, 2nd ed., OpenStax (Rice University), 2020, §9.3 "Use Properties of Angles, Triangles, and the Pythagorean Theorem" (complementary angles: measures that add to 90°).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: π/6 → π/3 (hand-derived).
- **Errors**: ErrNotFinite when angleRad is NaN or ±Inf; ErrOutOfRange when angleRad is outside [0, π/2].

### `NUM-GEOM-010` `SupplementaryAngleRad`  (supplementary_angle.go)
- **Source**: Marecek, Anthony-Smith, Honeycutt Mathis, Prealgebra, 2nd ed., OpenStax (Rice University), 2020, §9.3 "Use Properties of Angles, Triangles, and the Pythagorean Theorem" (supplementary angles: measures that add to 180°).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: π/3 → 2π/3 (hand-derived).
- **Errors**: ErrNotFinite when angleRad is NaN or ±Inf; ErrOutOfRange when angleRad is outside [0, π].

### `NUM-GEOM-011` `PolygonInteriorAngleSumRad`  (polygon_angles.go)
- **Source**: CRC31 §4.5 "Polygons", p. 318 ("In any k-gon, the sum of the angles equals 2(k − 2) right angles").
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 5 → 3π (hand-derived).
- **Errors**: ErrTooFewSides when sides < 3.

### `NUM-GEOM-012` `RegularPolygonInteriorAngleRad`  (polygon_angles.go)
- **Source**: CRC31 §4.5.3 "Regular polygons", p. 324 (θ = ((k − 2)/k)·180°).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 6 → 2π/3 (hand-derived).
- **Errors**: ErrTooFewSides when sides < 3.

### `NUM-GEOM-013` `RegularPolygonExteriorAngleRad`  (polygon_angles.go)
- **Source**: CRC31 §4.5 "Polygons", p. 318 (the exterior angle: between one side and the extension of the adjacent side), with θ of §4.5.3 "Regular polygons", p. 324.
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 8 → π/4 (hand-derived).
- **Errors**: ErrTooFewSides when sides < 3.

### `NUM-GEOM-014` `TrianglePerimeter`  (triangle_perimeter.go)
- **Source**: CRC31 §4.5.1 "Triangles", p. 319 (semi-perimeter s = ½(a + b + c)) and p. 321 (additional fact 4: each side is less than the sum of the other two); Kahan 2014 §2 "How to compute ∆" (sort a ≥ b ≥ c; the data are not a triangle if c − (a − b) < 0), p. 3.
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 3, 4, 5 → 12 (hand-derived).
- **Errors**: ErrNotFinite when a side is NaN or ±Inf; ErrNonPositive when a side ≤ 0; ErrNotTriangle when the triangle inequality fails; ErrOverflow when the perimeter is beyond float64.

### `NUM-GEOM-015` `TriangleAreaBaseHeight`  (triangle_area.go)
- **Source**: CRC31 §4.5.1 "Triangles", p. 319 (area = ½h_c·c).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: base 6, height 4 → 12 (hand-derived).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when the area is beyond float64.

### `NUM-GEOM-016` `TriangleAreaHeron`  (triangle_area.go)
- **Source**: Kahan 2014 §2 "How to compute ∆", p. 3, and §5 Table 1 "Area ∆ and Angle C", p. 4; CRC31 §4.5.1 "Triangles", p. 319 (area = √(s(s−a)(s−b)(s−c)), Heron).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: (10, 10, 10) → 43.30127020; (100000, 99999.99979, 0.00029) → 9.999999990; (10000, 5000.000001, 15000) → 612.3724358 (Kahan 2014, Table 1, "Accurate ∆").
- **Errors**: ErrNotFinite when a side is NaN or ±Inf; ErrNonPositive when a side ≤ 0; ErrNotTriangle when c − (a − b) < 0; ErrOverflow when the area is beyond float64.

### `NUM-GEOM-017` `TriangleAreaSAS`  (triangle_area.go)
- **Source**: CRC31 §4.5.1 "Triangles", p. 319 (area = ½ab sin C).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 3, 4, π/6 → 3 (hand-derived).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when a side ≤ 0; ErrOutOfRange when includedAngleRad is outside (0, π); ErrOverflow when the area is beyond float64.

### `NUM-GEOM-020` `RightTriangleHypotenuse`  (right_triangle.go)
- **Source**: CRC31 §4.5.1 "Triangles", p. 320 (right triangle: c² = a² + b², Pythagoras).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 3, 4 → 5 (hand-derived).
- **Errors**: ErrNotFinite when a leg is NaN or ±Inf; ErrNonPositive when a leg ≤ 0; ErrOverflow when the hypotenuse is beyond float64.

### `NUM-GEOM-021` `RightTriangleLeg`  (right_triangle.go)
- **Source**: CRC31 §4.5.1 "Triangles", p. 320 (right triangle: c² = a² + b², Pythagoras).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 5, 3 → 4 (hand-derived).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrNotTriangle when leg ≥ hypotenuse.

### `NUM-GEOM-022` `RectangleArea`  (rectangle.go)
- **Source**: CRC31 §4.5.2 "Quadrilaterals", p. 322 (parallelogram: area = bh; rectangle: area = ab, eq. (4.5.10)).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 3, 4 → 12 (hand-derived).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when the area is beyond float64.

### `NUM-GEOM-023` `RectanglePerimeter`  (rectangle.go)
- **Source**: CRC31 §4.5.2 "Quadrilaterals", p. 322 (Figure 4.13: half-perimeter s = ½(a + b + c + d); in a parallelogram opposite sides have the same length).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 3, 4 → 14 (hand-derived).
- **Errors**: ErrNotFinite when a side is NaN or ±Inf; ErrNonPositive when a side ≤ 0; ErrOverflow when the perimeter is beyond float64.

### `NUM-GEOM-024` `ParallelogramAreaSAS`  (parallelogram_area.go)
- **Source**: CRC31 §4.5.2 "Quadrilaterals", p. 322 (parallelogram: h = a sin A, area = bh).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 3, 4, π/6 → 6 (hand-derived).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when a side ≤ 0; ErrOutOfRange when angleRad is outside (0, π); ErrOverflow when the area is beyond float64.

### `NUM-GEOM-025` `TrapezoidArea`  (trapezoid_area.go)
- **Source**: CRC31 §4.5.2 "Quadrilaterals", p. 323 (trapezoid: area = ½(AB + CD)h).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 3, 5, 2 → 8 (hand-derived).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when the area is beyond float64.

### `NUM-GEOM-026` `KiteArea`  (kite_area.go)
- **Source**: CRC31 §4.5.2 "Quadrilaterals", eq. (4.5.9), p. 322 (area = ½pq sin θ) and p. 323 (rhombus: area = ½pq).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 6, 4 → 12 (hand-derived).
- **Errors**: ErrNotFinite when a diagonal is NaN or ±Inf; ErrNonPositive when a diagonal ≤ 0; ErrOverflow when the area is beyond float64.

### `NUM-GEOM-027` `RegularPolygonAreaFromSide`  (regular_polygon_area.go)
- **Source**: CRC31 §4.5.3 "Regular polygons", p. 324 (area = ¼ka² cot(180°/k), and the table of areas for k = 3…12).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 6, 1 → 2.59808 (CRC, §4.5.3 table, regular hexagon: 2.59808a²).
- **Errors**: ErrTooFewSides when sides < 3; ErrNotFinite when side is NaN or ±Inf; ErrNonPositive when side ≤ 0; ErrOverflow when the area is beyond float64.

### `NUM-GEOM-028` `RegularPolygonAreaFromApothem`  (regular_polygon_area.go)
- **Source**: CRC31 §4.5.3 "Regular polygons", p. 324 (area = kr² tan(180°/k), r the apothem).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 6, √3/2 → 3√3/2 ≈ 2.598076211353 (hand-derived; equals the hexagon of side 1).
- **Errors**: ErrTooFewSides when sides < 3; ErrNotFinite when apothem is NaN or ±Inf; ErrNonPositive when apothem ≤ 0; ErrOverflow when the area is beyond float64.

### `NUM-GEOM-029` `RegularPolygonAreaFromCircumradius`  (regular_polygon_area.go)
- **Source**: CRC31 §4.5.3 "Regular polygons", p. 324 (area = ½kR² sin(360°/k)).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 6, 1 → 3√3/2 ≈ 2.598076211353 (hand-derived).
- **Errors**: ErrTooFewSides when sides < 3; ErrNotFinite when circumradius is NaN or ±Inf; ErrNonPositive when circumradius ≤ 0; ErrOverflow when the area is beyond float64.

### `NUM-GEOM-030` `RegularPolygonSideInscribed`  (regular_polygon_side.go)
- **Source**: CRC31 §4.5.3 "Regular polygons", p. 324 (a = 2R sin(180°/k)).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 4, 1 → √2; 6, 1 → 1 (hand-derived).
- **Errors**: ErrTooFewSides when sides < 3; ErrNotFinite when circumradius is NaN or ±Inf; ErrNonPositive when circumradius ≤ 0; ErrOverflow when the side is beyond float64.

### `NUM-GEOM-031` `RegularPolygonSideCircumscribed`  (regular_polygon_side.go)
- **Source**: CRC31 §4.5.3 "Regular polygons", p. 324 (a = 2r tan(180°/k)).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 4, 1 → 2 (hand-derived).
- **Errors**: ErrTooFewSides when sides < 3; ErrNotFinite when inradius is NaN or ±Inf; ErrNonPositive when inradius ≤ 0; ErrOverflow when the side is beyond float64.

### `NUM-GEOM-032` `PolygonArea`  (polygon_area.go)
- **Source**: CRC31 §4.5 "Polygons", eq. (4.5.1), p. 318; Braden, "The Surveyor's Area Formula", The College Mathematics Journal, 1986, 17(4):326–337, doi:10.2307/2686282; Neumaier 1974 (the compensated sum).
- **Time / space**: O(n)   n = len(vertices); O(1) extra, no allocation. **Input modified**: no.
- **Worked example**: (0,0), (4,0), (4,3), (2,5), (0,3) → 16 (hand-derived: rectangle 12 + triangle 4).
- **Errors**: ErrDegenerate when there are fewer than 3 vertices or the computed area is exactly 0 (all vertices collinear); ErrNotFinite when a coordinate is NaN or ±Inf; ErrOverflow when the area is beyond float64.

### `NUM-GEOM-033` `PolygonPerimeter`  (polygon_perimeter.go)
- **Source**: CRC31 §4.5 "Polygons", p. 318 (the sides A_iA_{i+1}, the last vertex joined to the first), with the side lengths of §4.4.2, eq. (4.4.5), p. 316; Neumaier 1974 (the compensated sum).
- **Time / space**: O(n)   n = len(vertices); O(1) extra, no allocation. **Input modified**: no.
- **Worked example**: (0,0), (4,0), (4,3), (2,5), (0,3) → 10 + 4√2 ≈ 15.656854249492 (hand-derived: 4 + 3 + 2√2 + 2√2 + 3).
- **Errors**: ErrDegenerate when there are fewer than 2 vertices or all vertices coincide; ErrNotFinite when a coordinate is NaN or ±Inf; ErrOverflow when the perimeter is beyond float64.

### `NUM-GEOM-034` `CircleCircumference`  (circle.go)
- **Source**: CRC31 §4.6.6 "Circles", p. 334 ("The length or circumference of a circle of radius r is 2πr").
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 1 → 2π ≈ 6.283185307180 (hand-derived).
- **Errors**: ErrNotFinite when radius is NaN or ±Inf; ErrNonPositive when radius ≤ 0; ErrOverflow when the circumference is beyond float64.

### `NUM-GEOM-035` `CircleArea`  (circle.go)
- **Source**: CRC31 §4.6.6 "Circles", p. 334 ("the area is πr²").
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 2 → 4π ≈ 12.566370614359 (hand-derived).
- **Errors**: ErrNotFinite when radius is NaN or ±Inf; ErrNonPositive when radius ≤ 0; ErrOverflow when the area is beyond float64.

### `NUM-GEOM-036` `AnnulusArea`  (annulus_area.go)
- **Source**: CRC31 §4.6.6 "Circles", p. 334 (area πr², of the outer circle less the inner).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 5, 3 → 16π ≈ 50.265482457437 (hand-derived).
- **Errors**: ErrNotFinite when a radius is NaN or ±Inf; ErrNonPositive when a radius ≤ 0; ErrOutOfRange when innerRadius ≥ outerRadius; ErrOverflow when the area is beyond float64.

### `NUM-GEOM-037` `ArcLength`  (arc.go)
- **Source**: CRC31 §4.6.6 "Circles", p. 334 ("The length of the arc of circle subtended by an angle θ … is rθ").
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 2, π/3 → 2π/3 ≈ 2.094395102393 (hand-derived).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when radius ≤ 0; ErrOutOfRange when angleRad is outside (0, 2π]; ErrOverflow when the length is beyond float64.

### `NUM-GEOM-038` `ChordLength`  (arc.go)
- **Source**: CRC31 §4.6.6 "Circles", p. 334 (c = 2R sin ½θ).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 2, π/3 → 2 (hand-derived: an equilateral triangle).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when radius ≤ 0; ErrOutOfRange when centralAngleRad is outside (0, 2π]; ErrOverflow when the chord is beyond float64.

### `NUM-GEOM-039` `SectorArea`  (arc.go)
- **Source**: CRC31 §4.6.6 "Circles", p. 334 (area of sector = ½Rs = ½R²θ).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 2, π/3 → 2π/3 ≈ 2.094395102393 (hand-derived).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when radius ≤ 0; ErrOutOfRange when angleRad is outside (0, 2π]; ErrOverflow when the area is beyond float64.

### `NUM-GEOM-040` `SegmentArea`  (segment_area.go)
- **Source**: CRC31 §4.6.6 "Circles", p. 334 (area of segment = ½R²(θ − sin θ)); DLMF eq. 4.19.1 (sin z = z − z³/3! + z⁵/5! − …).
- **Time / space**: O(1); at most 30 series terms; O(1), no allocation. **Input modified**: no.
- **Worked example**: 2, π/3 → 2(π/3 − √3/2) ≈ 0.362344294824; 1, π → π/2 (hand-derived).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when radius ≤ 0; ErrOutOfRange when centralAngleRad is outside (0, 2π]; ErrOverflow when the area is beyond float64.

### `NUM-GEOM-041` `ArcRadiusFromChordSagitta`  (arc_radius.go)
- **Source**: CRC31 §4.6.6 "Circles", p. 334 (c = √(4h(2R − h)), h the height of the segment; Figure 4.18, p. 335).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: chord 8, sagitta 2 → 5 (hand-derived: half-chord 4 and R − s = 3 make a 3-4-5 triangle).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when the radius is beyond float64.

### `NUM-GEOM-042` `CircleTangentLength`  (circle_tangent_length.go)
- **Source**: Euclid, Elements, Book III, Proposition 36 (the square on the tangent equals the rectangle of the whole secant and its part outside the circle); D. E. Joyce (ed.), Euclid's Elements, Clark University.
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 6371000, 6371002 → √(2·6371000·2 + 4) ≈ 5048.168380706 (hand-derived).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOutOfRange when distanceFromCenter < radius (the point is inside the circle).

### `NUM-GEOM-043` `EllipseArea`  (ellipse_area.go)
- **Source**: CRC31 §4.6.3 "Additional properties of ellipses", eq. (4.6.11), p. 330 (area of C = πab).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 5, 3 → 15π ≈ 47.123889803847 (hand-derived).
- **Errors**: ErrNotFinite when a semi-axis is NaN or ±Inf; ErrNonPositive when a semi-axis ≤ 0; ErrOverflow when the area is beyond float64.

### `NUM-GEOM-044` `EllipseChordFromCenter`  (ellipse_chord.go)
- **Source**: CRC31 §4.6.3 "Additional properties of ellipses", p. 330 (the ellipse x²/a² + y²/b² = 1).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: a 5, b 3, offset 1.8 → 8 (hand-derived: √(1 − 0.36) = 0.8).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when a semi-axis ≤ 0; ErrNegative when offsetFromCenter < 0; ErrOutOfRange when offsetFromCenter > semiAxisAcross; ErrOverflow when the chord is beyond float64.

### `NUM-GEOM-045` `EllipseChordFromEdge`  (ellipse_chord.go)
- **Source**: CRC31 §4.6.3 "Additional properties of ellipses", p. 330 (the ellipse x²/a² + y²/b² = 1).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: a 5, b 3, depth 1.2 → 8 (hand-derived; the same chord as NUM-GEOM-044's example).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when a semi-axis ≤ 0; ErrNegative when depthFromEdge < 0; ErrOutOfRange when depthFromEdge > 2·semiAxisAcross; ErrOverflow when the chord is beyond float64.

### `NUM-GEOM-046` `EllipsePerimeter`  (ellipse_perimeter.go)
- **Source**: DLMF eq. 19.9.9 (L(a,b) = 4aE(k), k² = 1 − b²/a²); Carlson 1995 §2, eqs. (42)–(46) (E(k) = 2R_G(1 − k², 1, 0) by arithmetic and geometric means); preprint arXiv:math/9409227 (read 2026-11-12); DLMF eq. 19.12.2 (E(k) as k′ → 0, used for b/a < 2⁻³⁰); CRC31 §4.6.3 "Additional properties of ellipses", eq. (4.6.11), p. 330 (perimeter of C = 4aE(π/2, e)).
- **Time / space**: O(1); at most 64 AGM steps, about 7 in practice; O(1), no allocation. **Input modified**: no.
- **Worked example**: 2, 1 → 9.688448220547676; 5, 3 → 25.526998863398128 (hand-derived, 60-digit AGM).
- **Errors**: ErrNotFinite when a semi-axis is NaN or ±Inf; ErrNonPositive when a semi-axis ≤ 0; ErrOverflow when the perimeter is beyond float64; ErrNoConvergence if the AGM reaches its cap (not expected).

### `NUM-GEOM-047` `EllipsePerimeterRamanujan`  (ellipse_perimeter_ramanujan.go)
- **Source**: Ramanujan, "Modular equations and approximations to π", Quarterly Journal of Mathematics, 1914, 45:350–372, §15 (p = π{(a+b) + 3(a−b)²/[10(a+b) + √(a²+14ab+b²)]}, "obtained empirically", error about 3a·e²⁰/2³⁶); DLMF eq. 19.9.10 (the same approximation in λ, with its leading error term 3λ¹⁰/2¹⁷).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 2, 1 → 9.688448216130 (exact 9.688448220548; hand-derived, 60 digits).
- **Errors**: ErrNotFinite when a semi-axis is NaN or ±Inf; ErrNonPositive when a semi-axis ≤ 0; ErrOverflow when the perimeter is beyond float64.

### `NUM-GEOM-048` `PrismVolume`  (prism_volume.go)
- **Source**: CRC31 §4.15 "Cylinders", eq. (4.15.1), p. 361 (V = Ah, h the perpendicular distance between the planes), and §4.14 "Polyhedra", p. 357 (a cylinder with a polygonal base is a prism).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: base 12, height 5 → 60 (hand-derived).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when the volume is beyond float64.

### `NUM-GEOM-049` `CylinderVolume`  (cylinder.go)
- **Source**: CRC31 §4.15 "Cylinders", p. 361 (right circular cylinder: volume πr²h; eq. (4.15.1) for the oblique case).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 2, 5 → 20π ≈ 62.831853071796 (hand-derived).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when the volume is beyond float64.

### `NUM-GEOM-050` `CylinderLateralArea`  (cylinder.go)
- **Source**: CRC31 §4.15 "Cylinders", p. 361 (right circular cylinder: lateral area 2πrh).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 2, 5 → 20π ≈ 62.831853071796 (hand-derived).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when the area is beyond float64.

### `NUM-GEOM-051` `CylinderSurfaceArea`  (cylinder.go)
- **Source**: CRC31 §4.15 "Cylinders", p. 361 (right circular cylinder: total area 2πr(r + h)).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 2, 5 → 28π ≈ 87.964594300514 (hand-derived).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when the area is beyond float64.

### `NUM-GEOM-052` `HollowCylinderVolume`  (hollow_cylinder.go)
- **Source**: CRC31 §4.15 "Cylinders", eq. (4.15.1), p. 361 (V = Ah), with the area of the annular base from §4.6.6 "Circles", p. 334 (area πr²).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 3, 2, 5 → 25π ≈ 78.539816339745 (hand-derived).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOutOfRange when innerRadius ≥ outerRadius; ErrOverflow when the volume is beyond float64.

### `NUM-GEOM-053` `HollowCylinderSurfaceArea`  (hollow_cylinder.go)
- **Source**: CRC31 §4.15 "Cylinders", p. 361 (lateral area 2πrh, for each wall), with the annular ends from §4.6.6 "Circles", p. 334 (area πr²).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 3, 2, 5 → 60π ≈ 188.495559215388 (hand-derived).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOutOfRange when innerRadius ≥ outerRadius; ErrOverflow when the area is beyond float64.

### `NUM-GEOM-054` `PyramidVolume`  (pyramid.go)
- **Source**: CRC31 §4.16 "Cones", eq. (4.16.1), p. 362 (V = ⅓Ah), and §4.14 "Polyhedra", p. 357 (a cone with a polygonal base is a pyramid).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: base 36 (a square of side 6), height 4 → 48 (hand-derived).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when the volume is beyond float64.

### `NUM-GEOM-055` `RegularPyramidSurfaceArea`  (pyramid.go)
- **Source**: CRC31 §4.5.3 "Regular polygons", p. 324 (base area ¼ka² cot(180°/k) and apothem r = ½a cot(180°/k)), §4.5.1 "Triangles", p. 319 (face area ½h_c c) and p. 320 (Pythagoras, for the slant height).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 4, 6, 4 → 36 + 4·(½·6·5) = 96 (hand-derived: slant height √(16 + 9) = 5).
- **Errors**: ErrTooFewSides when sides < 3; ErrNotFinite when side or height is NaN or ±Inf; ErrNonPositive when side or height ≤ 0; ErrOverflow when the area is beyond float64.

### `NUM-GEOM-056` `PyramidFrustumVolume`  (pyramid.go)
- **Source**: CRC31 §4.16 "Cones", eq. (4.16.2), p. 362 (V = ⅓h(A + A′ + √(AA′))).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 16, 4, 3 → 28 (hand-derived; PrismoidVolume gives the same for this solid).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when the volume is beyond float64.

### `NUM-GEOM-057` `PrismoidVolume`  (prismoid_volume.go)
- **Source**: Encyclopædia Britannica, 11th ed., 1911, vol. 18, "Mensuration", §58 "The Prismoidal Formula", p. 160 (volume ⅙H(S₀ + 4S₁ + S₂) when the cross-section area is a quadratic function of the distance).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: square frustum with end sides 4 and 2, height 3: 16, 9, 4, 3 → 28 (hand-derived; equals PyramidFrustumVolume).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when height ≤ 0; ErrNegative when an area < 0; ErrOverflow when the volume is beyond float64.

### `NUM-GEOM-058` `ConeVolume`  (cone.go)
- **Source**: CRC31 §4.16 "Cones", p. 362 (right circular cone: volume ⅓πr²h; eq. (4.16.1) for any base).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 3, 4 → 12π ≈ 37.699111843078 (hand-derived).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when the volume is beyond float64.

### `NUM-GEOM-059` `ConeLateralArea`  (cone.go)
- **Source**: CRC31 §4.16 "Cones", p. 362 (ℓ = √(r² + h²), lateral area = πrℓ).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 3, 4 → 15π ≈ 47.123889803847 (hand-derived: ℓ = 5).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when the area is beyond float64.

### `NUM-GEOM-060` `ConeSurfaceArea`  (cone.go)
- **Source**: CRC31 §4.16 "Cones", p. 362 (total area = πr(ℓ + r)).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 3, 4 → 24π ≈ 75.398223686155 (hand-derived).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when the area is beyond float64.

### `NUM-GEOM-061` `ConeFrustumVolume`  (cone_frustum.go)
- **Source**: CRC31 §4.16 "Cones", p. 363 (frustum of a right circular cone: volume = ⅓πh(r₁² + r₂² + r₁r₂)).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 4, 1, 4 → 28π ≈ 87.964594300514 (hand-derived).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when the volume is beyond float64.

### `NUM-GEOM-062` `ConeFrustumSurfaceArea`  (cone_frustum.go)
- **Source**: CRC31 §4.16 "Cones", p. 363 (frustum: ℓ = √((r₁ − r₂)² + h²), total area = π(r₁² + r₂² + (r₁ + r₂)ℓ)).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 4, 1, 4 → 42π ≈ 131.946891450771 (hand-derived: ℓ = 5).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when the area is beyond float64.

### `NUM-GEOM-063` `SphereVolume`  (sphere.go)
- **Source**: CRC31 §4.18.1 "Spheres", item 3, p. 367 (the volume is (4/3)πr³).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 3 → 36π ≈ 113.097335529233 (hand-derived).
- **Errors**: ErrNotFinite when radius is NaN or ±Inf; ErrNonPositive when radius ≤ 0; ErrOverflow when the volume is beyond float64.

### `NUM-GEOM-064` `SphereSurfaceArea`  (sphere.go)
- **Source**: CRC31 §4.18.1 "Spheres", item 3, p. 367 (the area of a sphere of radius r is 4πr²).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 3 → 36π ≈ 113.097335529233 (hand-derived).
- **Errors**: ErrNotFinite when radius is NaN or ±Inf; ErrNonPositive when radius ≤ 0; ErrOverflow when the area is beyond float64.

### `NUM-GEOM-065` `SphericalCapVolume`  (spherical_cap.go)
- **Source**: CRC31 §4.18.1.1 "Spherical cap", p. 367 (the volume of the cap is ⅓πh²(3r − h)).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 5, 2 → 52π/3 ≈ 54.454272662223 (hand-derived).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOutOfRange when capHeight > 2·sphereRadius; ErrOverflow when the volume is beyond float64.

### `NUM-GEOM-066` `SphericalCapArea`  (spherical_cap.go)
- **Source**: CRC31 §4.18.1.1 "Spherical cap", p. 367 (the area of the curved region is 2πrh).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 5, 2 → 20π ≈ 62.831853071796 (hand-derived).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOutOfRange when capHeight > 2·sphereRadius; ErrOverflow when the area is beyond float64.

### `NUM-GEOM-067` `TorusVolume`  (torus.go)
- **Source**: CRC31 §4.17 "Surfaces of revolution: the torus", p. 363 (its volume is 2π²Rr²; the Pappus–Guldinus theorem).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 3, 1 → 6π² ≈ 59.217626406536 (hand-derived).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOutOfRange when minorRadius > majorRadius; ErrOverflow when the volume is beyond float64.

### `NUM-GEOM-068` `TorusSurfaceArea`  (torus.go)
- **Source**: CRC31 §4.17 "Surfaces of revolution: the torus", p. 363 (the area of the torus is 4π²Rr).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 3, 1 → 12π² ≈ 118.435252813072 (hand-derived).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOutOfRange when minorRadius > majorRadius; ErrOverflow when the area is beyond float64.

### `NUM-GEOM-069` `EllipsoidVolume`  (ellipsoid_volume.go)
- **Source**: CRC31 §4.18 "Quadrics", p. 364 (the volume of the ellipsoid with semi-axes a, b, c is (4/3)πabc); DLMF §19.33 (V = 4πabc/3, beside eq. 19.33.1).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 3, 2, 1 → 8π ≈ 25.132741228718 (hand-derived).
- **Errors**: ErrNotFinite when a semi-axis is NaN or ±Inf; ErrNonPositive when a semi-axis ≤ 0; ErrOverflow when the volume is beyond float64.

### `NUM-GEOM-070` `EllipsoidSurfaceArea`  (ellipsoid_surface_area.go)
- **Source**: DLMF eq. 19.33.1 (S = 3V R_G(a⁻², b⁻², c⁻²)) and eq. 19.21.11 (R_G from R_D); Carlson 1995 §2 "Algorithm for R_D", eqs. (35)–(41), and §3 check values; preprint arXiv:math/9409227 (read 2026-11-12).
- **Time / space**: O(1); three R_D evaluations of at most 100 duplication steps, under 20 in practice; O(1), no allocation. **Input modified**: no.
- **Worked example**: 3, 2, 1 → 48.882146302582; 2, 1, 1 → 21.478435327884 (hand-derived, 60 digits; the second equals the prolate-spheroid closed form, CRC 31st ed. eq. (4.18.1)).
- **Errors**: ErrNotFinite when a semi-axis is NaN or ±Inf; ErrNonPositive when a semi-axis ≤ 0; ErrOverflow when the area is beyond float64; ErrNoConvergence if a duplication loop reaches its cap (not expected).

### `NUM-GEOM-071` `TetrahedronVolume`  (tetrahedron.go)
- **Source**: CRC31 §4.14.1 "Convex regular polyhedra", p. 360 (table of V/a³ for the five Platonic solids).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 1 → 0.117851 (CRC §4.14.1 table, V/a³).
- **Errors**: ErrNotFinite when edge is NaN or ±Inf; ErrNonPositive when edge ≤ 0; ErrOverflow when the volume is beyond float64.

### `NUM-GEOM-072` `TetrahedronSurfaceArea`  (tetrahedron.go)
- **Source**: CRC31 §4.14.1 "Convex regular polyhedra", p. 360 (table of S/a² for the five Platonic solids).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 1 → 1.73205 (CRC §4.14.1 table, S/a²).
- **Errors**: ErrNotFinite when edge is NaN or ±Inf; ErrNonPositive when edge ≤ 0; ErrOverflow when the area is beyond float64.

### `NUM-GEOM-073` `CubeVolume`  (cube.go)
- **Source**: CRC31 §4.14.1 "Convex regular polyhedra", p. 360 (table of V/a³ for the five Platonic solids).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 1 → 1 (CRC §4.14.1 table, V/a³).
- **Errors**: ErrNotFinite when edge is NaN or ±Inf; ErrNonPositive when edge ≤ 0; ErrOverflow when the volume is beyond float64.

### `NUM-GEOM-074` `CubeSurfaceArea`  (cube.go)
- **Source**: CRC31 §4.14.1 "Convex regular polyhedra", p. 360 (table of S/a² for the five Platonic solids).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 1 → 6 (CRC §4.14.1 table, S/a²).
- **Errors**: ErrNotFinite when edge is NaN or ±Inf; ErrNonPositive when edge ≤ 0; ErrOverflow when the area is beyond float64.

### `NUM-GEOM-075` `OctahedronVolume`  (octahedron.go)
- **Source**: CRC31 §4.14.1 "Convex regular polyhedra", p. 360 (table of V/a³ for the five Platonic solids).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 1 → 0.471405 (CRC §4.14.1 table, V/a³).
- **Errors**: ErrNotFinite when edge is NaN or ±Inf; ErrNonPositive when edge ≤ 0; ErrOverflow when the volume is beyond float64.

### `NUM-GEOM-076` `OctahedronSurfaceArea`  (octahedron.go)
- **Source**: CRC31 §4.14.1 "Convex regular polyhedra", p. 360 (table of S/a² for the five Platonic solids).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 1 → 3.46410 (CRC §4.14.1 table, S/a²).
- **Errors**: ErrNotFinite when edge is NaN or ±Inf; ErrNonPositive when edge ≤ 0; ErrOverflow when the area is beyond float64.

### `NUM-GEOM-077` `DodecahedronVolume`  (dodecahedron.go)
- **Source**: CRC31 §4.14.1 "Convex regular polyhedra", p. 360 (table of V/a³ for the five Platonic solids).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 1 → 7.663119 (CRC §4.14.1 table, V/a³).
- **Errors**: ErrNotFinite when edge is NaN or ±Inf; ErrNonPositive when edge ≤ 0; ErrOverflow when the volume is beyond float64.

### `NUM-GEOM-078` `DodecahedronSurfaceArea`  (dodecahedron.go)
- **Source**: CRC31 §4.14.1 "Convex regular polyhedra", p. 360 (table of S/a² for the five Platonic solids).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 1 → 20.64573 (CRC §4.14.1 table, S/a²).
- **Errors**: ErrNotFinite when edge is NaN or ±Inf; ErrNonPositive when edge ≤ 0; ErrOverflow when the area is beyond float64.

### `NUM-GEOM-079` `IcosahedronVolume`  (icosahedron.go)
- **Source**: CRC31 §4.14.1 "Convex regular polyhedra", p. 360 (table of V/a³ for the five Platonic solids).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 1 → 2.181695 (CRC §4.14.1 table, V/a³).
- **Errors**: ErrNotFinite when edge is NaN or ±Inf; ErrNonPositive when edge ≤ 0; ErrOverflow when the volume is beyond float64.

### `NUM-GEOM-080` `IcosahedronSurfaceArea`  (icosahedron.go)
- **Source**: CRC31 §4.14.1 "Convex regular polyhedra", p. 360 (table of S/a² for the five Platonic solids).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 1 → 8.66025 (CRC §4.14.1 table, S/a²).
- **Errors**: ErrNotFinite when edge is NaN or ±Inf; ErrNonPositive when edge ≤ 0; ErrOverflow when the area is beyond float64.

### `NUM-GEOM-081` `HelixArcLength`  (helix.go)
- **Source**: CRC31 §4.20.1.3 "Example", p. 375 (the circular helix (a cos t, a sin t, bt): ds/dt = √(a² + b²), arc length s = t√(a² + b²)).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 1, 2, 3 → 3·√(4π² + 4) ≈ 19.781449856854 (hand-derived).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when the length is beyond float64.

### `NUM-GEOM-082` `HelixHeight`  (helix.go)
- **Source**: CRC31 §4.20.1.3 "Example", p. 375 (arc length s = t√(a² + b²) of the circular helix, solved for the height bt).
- **Time / space**: O(1); O(1), no allocation. **Input modified**: no.
- **Worked example**: 1, 19.781449856854, 3 → 6 (hand-derived; the inverse of NUM-GEOM-081's example).
- **Errors**: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOutOfRange when arcLength < 2π·radius·turns (too short to make that many turns).

## Out of scope

- NUM-GEOM-018 TriangleSideSAS and NUM-GEOM-019 TriangleAngleSSSRad: the law of cosines is num/trig
  (NUM-TRIG-020 LawOfCosinesSide, NUM-TRIG-021 LawOfCosinesAngle). The IDs are REMOVED rows, written by the lead.
- "Bell-shaped" arc length, area, height (catalog input rows 11–13): shape undefined; no ID.
- LATER (catalog STEP 1a/1b): helix volume, Möbius strip, Morley's theorem, conical helix, polygon centroid and
  second moment, signed polygon area, point in polygon, circumradius and inradius, 3D point distance, distance to a
  segment, sagitta from radius and chord, cap from its base radius.
- Degree variants: none (radians only).
