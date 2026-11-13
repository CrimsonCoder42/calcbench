# Catalog: geometry (`num/geometry`, SHORT `GEOM`)

Status: CATALOG, awaiting Devin's approval. Nothing is built until it is approved.
Work folder: [../geometry_2026-11-12/](../geometry_2026-11-12/README.md). Rows for shared files:
[PENDING_SHARED.md](../geometry_2026-11-12/PENDING_SHARED.md).

The input list is the names on vCalc's Geometry top-level page, taken as NAMES only (`.claude/rules/sources.md`).
No vCalc page was fetched, so where a name is ambiguous ("Interface Angle", "Bell-shaped") the reading below is
ours and is marked as such. Formulas come from the primary sources cited in STEP 2.

Conventions used throughout (each is a question in the report until Devin decides):
- Angles are radians, and every angle parameter and result says `Rad` in its name.
- Lengths are `float64` in any one consistent unit; areas come out in that unit squared and volumes in that
  unit cubed. `num/geometry` cannot use `unit/length` types: `unit` may import `num`, so the reverse would be
  a cycle (CLAUDE.md import rule).
- Points are a package-local `Point{X, Y float64}`.
- Every function returns `(value, error)`. Shared sentinels: `ErrNonPositive` (a length, count or area that
  must be > 0), `ErrNotFinite` (NaN or ±Inf in), `ErrOverflow` (a finite input whose result leaves the float64
  range), `ErrOutOfRange` (an angle or ratio outside its domain), `ErrNotTriangle` (triangle inequality fails),
  `ErrDegenerate` (coincident points, zero-length line, fewer than 3 vertices, zero-area polygon), `ErrTooFewSides`
  (n < 3), `ErrNegative` (a value that may be 0 but not < 0), `ErrVertical` (slope of a vertical line),
  `ErrNoConvergence` (an iteration cap reached; never expected, but a library must not loop).

## STEP 1a: the input list, cleaned

Decisions: KEEP (goes in `num/geometry`) · MOVE (another home) · MERGE (folded into another row's function) ·
DROP (not a formula) · LATER (real geometry, not in this pass).

| # | Name as given | Decision | Home | Reason |
|---|---|---|---|---|
| 1 | Cylinder Calculator (collection) | DROP | — | A calculator, not a formula. Bundles cylinder volume, lateral and total area: NUM-GEOM-049…051 |
| 2 | Ellipsoid Calc (calculator) | DROP | — | A calculator. Bundles ellipsoid volume and surface area: NUM-GEOM-069, -070 |
| 3 | Geometric Area (calculator) | DROP | — | A calculator. Bundles the plane areas of the standard set (triangle, quadrilaterals, circle, ellipse, polygons): NUM-GEOM-015…047 |
| 4 | LM 15.2 Angular momentum in planetary motion (collection) | MOVE | `sci/physics` | Mechanics (L = m·r×v, Kepler's second law), not geometry. A textbook-lesson collection: the formulas move, the bundle does not |
| 5 | Moments of Inertia (calculator) | DROP | — | A calculator. Bundles section and body moments of inertia; the formula is row 35 (MOVE to `eng/mechanical`) |
| 6 | Quilting (calculator) | MOVE | `life/home` | taxonomy.md: "Home Improvement, Quilting" → `life/home`. Uses rectangle and triangle areas from here |
| 7 | Sphere Calc (calculator) | DROP | — | A calculator. Bundles sphere volume and surface area: NUM-GEOM-063, -064 (and cap: -065, -066) |
| 8 | Trapezoidal Volume (collection) | MERGE | `num/geometry` | Into the prismoid, NUM-GEOM-057 (with row 43) |
| 9 | Area of Circular Annulus | KEEP | `num/geometry` | NUM-GEOM-036 AnnulusArea |
| 10 | Area of Regular Polygon [apothem, number of sides] | KEEP | `num/geometry` | NUM-GEOM-028 RegularPolygonAreaFromApothem |
| 11 | Bell-shaped Arc Length | MERGE | `num/geometry` | Rows 11–13 are one shape. SHAPE UNDEFINED: see the STEP 2 note "Bell-shaped". No ID reserved |
| 12 | Bell-shaped Area | MERGE | `num/geometry` | Into row 11 |
| 13 | Bell-shaped Height | MERGE | `num/geometry` | Into row 11 |
| 14 | Bottle - mass | DROP | — | Not a geometry formula (a mass of a specific object); no general formula behind the name |
| 15 | Central angle of horizon | MOVE | `sci/astronomy` | Needs a planet's radius (a physical constant). The geometric core, arccos(R/(R+h)), is one `math.Acos` call |
| 16 | Centroid of a Beam Section | MOVE | `eng/mechanical` | Section properties of a structural member (mechanics of materials). See STEP 1b note on polygon centroid |
| 17 | Complementary Angle | KEEP | `num/geometry` | NUM-GEOM-009 ComplementaryAngleRad |
| 18 | Distance to Sea Level Horizon | MOVE | `sci/astronomy` | Needs Earth's radius and refraction policy. Its geometric core is NUM-GEOM-042 CircleTangentLength, which it will call |
| 19 | Ellipse Horizontal Chord from Center | MERGE | `num/geometry` | Rows 19–22 are two functions: horizontal vs vertical is only which semi-axis is "along". NUM-GEOM-044 EllipseChordFromCenter |
| 20 | Ellipse Horizontal Chord from Edge | MERGE | `num/geometry` | NUM-GEOM-045 EllipseChordFromEdge |
| 21 | Ellipse Vertical Chord from Center | MERGE | `num/geometry` | NUM-GEOM-044 (axes swapped) |
| 22 | Ellipse Vertical Chord from Edge | MERGE | `num/geometry` | NUM-GEOM-045 (axes swapped) |
| 23 | Five Sided Polygon Area | MERGE | `num/geometry` | One shoelace polygon area with rows 30 and 37: NUM-GEOM-032 PolygonArea |
| 24 | Golden Ratio Sequence | MOVE | `num/series` | A sequence (powers of φ / Fibonacci ratios), not a shape |
| 25 | Helix Height | KEEP | `num/geometry` | NUM-GEOM-082 HelixHeight |
| 26 | Helix Length | KEEP | `num/geometry` | NUM-GEOM-081 HelixArcLength |
| 27 | Helix Volume | LATER | `num/geometry` | Volume of a tube swept along a helix (Pappus); needs a decision on self-overlap of tight coils |
| 28 | Interface Angle | KEEP | `num/geometry` | Read here as the angle between two lines: NUM-GEOM-006 AngleBetweenLinesRad. OUR READING of the name; vCalc's page was not read. Devin to confirm |
| 29 | Interface Angle from Rise over Run | KEEP | `num/geometry` | NUM-GEOM-004 SlopeAngleRad |
| 30 | Irregular Pentagon Area | MERGE | `num/geometry` | NUM-GEOM-032 PolygonArea |
| 31 | Length of Spiral Christmas Garland | MOVE | `life/home` | An everyday use of NUM-GEOM-081 HelixArcLength (garland around a tree = helix on a cone/cylinder) |
| 32 | Mass of the Earth As Told By A Pig | DROP | — | Not a formula (a story page). Earth's mass, if ever needed, is a cited constant in `sci/astronomy` |
| 33 | Mass Unit Conversion | MOVE | `unit/mass` | A unit conversion |
| 34 | Möbius strip coordinates | LATER | `num/geometry` | A parametric surface; needs a parametric-surface API shape first |
| 35 | Moment of Inertia of a Beam Section | MOVE | `eng/mechanical` | Second moment of area of a structural section |
| 36 | Morley's Trisector Theorem | LATER | `num/geometry` | A theorem; as a function it is "Morley triangle side from the circumradius and angles". Not in this pass |
| 37 | Polygon - area4 | MERGE | `num/geometry` | NUM-GEOM-032 PolygonArea |
| 38 | Polygon Side around a Circle | KEEP | `num/geometry` | NUM-GEOM-031 RegularPolygonSideCircumscribed |
| 39 | Polygon Side within a Circle | KEEP | `num/geometry` | NUM-GEOM-030 RegularPolygonSideInscribed |
| 40 | Pyramid Geometries | KEEP | `num/geometry` | NUM-GEOM-054 PyramidVolume, -055 RegularPyramidSurfaceArea |
| 41 | Radius of an Arch | KEEP | `num/geometry` | NUM-GEOM-041 ArcRadiusFromChordSagitta |
| 42 | Supplementary Angle | KEEP | `num/geometry` | NUM-GEOM-010 SupplementaryAngleRad |
| 43 | Trapezoid Volume | MERGE | `num/geometry` | With row 8: one prismoid, NUM-GEOM-057 PrismoidVolume |
| 44 | Velocity Unit Conversion | MOVE | `unit/speed` | A unit conversion |
| 45 | Golden Ratio (data item) | MOVE | `num/logs` | A constant, φ = (1+√5)/2. `num/logs` already holds "transcendental constants"; proposed wording change in PENDING_SHARED.md |
| 46 | WikiClips: Angle Function Information | DROP | — | Not a formula (an explanatory clip). Its subject is `num/trig` |
| 47 | WikiClips: Dodecahedron Information | DROP | — | Not a formula. The dodecahedron is NUM-GEOM-077, -078 |
| 48 | WikiClips: Ellipse Calculator Information | DROP | — | Not a formula. Ellipse: NUM-GEOM-043…047 |
| 49 | WikiClips: Ellipsoid Calc Info | DROP | — | Not a formula. Ellipsoid: NUM-GEOM-069, -070 |
| 50 | WikiClips: Polar Coordinate Functions | DROP | — | Not a formula. Polar conversion: NUM-GEOM-007, -008 |
| 51 | WikiClips: Pyramid Calculator Info | DROP | — | Not a formula. Pyramid: NUM-GEOM-054…056 |

Counts: KEEP 12 · MERGE 12 (into 6 functions and 1 undefined shape) · MOVE 11 · DROP 13 · LATER 3 · total 51.

## STEP 1b: the standard set

What any geometry library must have, catalogued independently of the input list. Each line is one function;
"Covers" names the input rows (STEP 1a numbers) it answers. Shapes are grouped in build order (STEP 3).

| ID | Function | Shape / relation | Quantity | Covers |
|---|---|---|---|---|
| NUM-GEOM-001 | PointDistance | two points | distance | — |
| NUM-GEOM-002 | Midpoint | two points | midpoint | — |
| NUM-GEOM-003 | Slope | line through two points | rise over run | — |
| NUM-GEOM-004 | SlopeAngleRad | rise and run | angle of the slope | 29 |
| NUM-GEOM-005 | PointLineDistance | point and line | perpendicular distance | — |
| NUM-GEOM-006 | AngleBetweenLinesRad | two lines | acute angle between | 28 (our reading) |
| NUM-GEOM-007 | PolarToCartesian | point | (r, θ) → (x, y) | 51 (clip) |
| NUM-GEOM-008 | CartesianToPolar | point | (x, y) → (r, θ) | 51 (clip) |
| NUM-GEOM-009 | ComplementaryAngleRad | angle | π/2 − θ | 17 |
| NUM-GEOM-010 | SupplementaryAngleRad | angle | π − θ | 42 |
| NUM-GEOM-011 | PolygonInteriorAngleSumRad | n-gon | sum of interior angles | — |
| NUM-GEOM-012 | RegularPolygonInteriorAngleRad | regular n-gon | one interior angle | — |
| NUM-GEOM-013 | RegularPolygonExteriorAngleRad | regular n-gon | one exterior angle | — |
| NUM-GEOM-014 | TrianglePerimeter | triangle, three sides | perimeter (checks the triangle inequality) | — |
| NUM-GEOM-015 | TriangleAreaBaseHeight | triangle, base and height | area | 3 |
| NUM-GEOM-016 | TriangleAreaHeron | triangle, three sides | area (Kahan's stable Heron) | 3 |
| NUM-GEOM-017 | TriangleAreaSAS | two sides and included angle | area | 3 |
| NUM-GEOM-018 | TriangleSideSAS | two sides and included angle | third side (law of cosines) | — |
| NUM-GEOM-019 | TriangleAngleSSSRad | three sides | angle opposite one side (Kahan) | — |
| NUM-GEOM-020 | RightTriangleHypotenuse | right triangle, two legs | hypotenuse | — |
| NUM-GEOM-021 | RightTriangleLeg | right triangle, hypotenuse and leg | other leg | — |
| NUM-GEOM-022 | RectangleArea | rectangle, square; parallelogram and rhombus by base × height | area | 3, 6 (via life/home) |
| NUM-GEOM-023 | RectanglePerimeter | rectangle, square, parallelogram, rhombus, kite (two pairs of equal sides) | perimeter | — |
| NUM-GEOM-024 | ParallelogramAreaSAS | parallelogram, rhombus: two sides and included angle | area | 3 |
| NUM-GEOM-025 | TrapezoidArea | trapezoid | area | 3 |
| NUM-GEOM-026 | KiteArea | kite, rhombus, square: perpendicular diagonals | area | 3 |
| NUM-GEOM-027 | RegularPolygonAreaFromSide | regular n-gon | area | 3 |
| NUM-GEOM-028 | RegularPolygonAreaFromApothem | regular n-gon | area | 10 |
| NUM-GEOM-029 | RegularPolygonAreaFromCircumradius | regular n-gon | area | — |
| NUM-GEOM-030 | RegularPolygonSideInscribed | n-gon inscribed in a circle | side | 39 |
| NUM-GEOM-031 | RegularPolygonSideCircumscribed | n-gon circumscribed about a circle | side | 38 |
| NUM-GEOM-032 | PolygonArea | simple polygon by vertices | area (shoelace) | 23, 30, 37, 3 |
| NUM-GEOM-033 | PolygonPerimeter | polygon by vertices | perimeter | — |
| NUM-GEOM-034 | CircleCircumference | circle | perimeter | — |
| NUM-GEOM-035 | CircleArea | circle | area | 3 |
| NUM-GEOM-036 | AnnulusArea | annulus | area | 9 |
| NUM-GEOM-037 | ArcLength | circular arc | length | — |
| NUM-GEOM-038 | ChordLength | circular chord | length | — |
| NUM-GEOM-039 | SectorArea | circular sector | area | — |
| NUM-GEOM-040 | SegmentArea | circular segment | area | — |
| NUM-GEOM-041 | ArcRadiusFromChordSagitta | arc from chord and sagitta (rise) | radius | 41 |
| NUM-GEOM-042 | CircleTangentLength | circle and external point | tangent length | 18 (core, via sci/astronomy) |
| NUM-GEOM-043 | EllipseArea | ellipse | area | 3 |
| NUM-GEOM-044 | EllipseChordFromCenter | ellipse | chord at an offset from the centre | 19, 21 |
| NUM-GEOM-045 | EllipseChordFromEdge | ellipse | chord at a depth from the vertex | 20, 22 |
| NUM-GEOM-046 | EllipsePerimeter | ellipse | perimeter, exact (AGM) | — |
| NUM-GEOM-047 | EllipsePerimeterRamanujan | ellipse | perimeter, Ramanujan's approximation | — |
| NUM-GEOM-048 | PrismVolume | right or oblique prism | volume | — |
| NUM-GEOM-049 | CylinderVolume | cylinder | volume | 1 |
| NUM-GEOM-050 | CylinderLateralArea | cylinder | lateral area | 1 |
| NUM-GEOM-051 | CylinderSurfaceArea | cylinder | total surface area | 1 |
| NUM-GEOM-052 | HollowCylinderVolume | hollow cylinder (tube) | volume | — |
| NUM-GEOM-053 | HollowCylinderSurfaceArea | hollow cylinder (tube) | total surface area | — |
| NUM-GEOM-054 | PyramidVolume | pyramid, any base | volume | 40 |
| NUM-GEOM-055 | RegularPyramidSurfaceArea | right pyramid, regular base | total surface area | 40 |
| NUM-GEOM-056 | PyramidFrustumVolume | frustum of a pyramid | volume | 40 |
| NUM-GEOM-057 | PrismoidVolume | prismoid | volume (prismoidal formula) | 8, 43 |
| NUM-GEOM-058 | ConeVolume | right circular cone | volume | — |
| NUM-GEOM-059 | ConeLateralArea | right circular cone | lateral area | — |
| NUM-GEOM-060 | ConeSurfaceArea | right circular cone | total surface area | — |
| NUM-GEOM-061 | ConeFrustumVolume | frustum of a cone | volume | — |
| NUM-GEOM-062 | ConeFrustumSurfaceArea | frustum of a cone | total surface area | — |
| NUM-GEOM-063 | SphereVolume | sphere | volume | 7 |
| NUM-GEOM-064 | SphereSurfaceArea | sphere | surface area | 7 |
| NUM-GEOM-065 | SphericalCapVolume | spherical cap | volume | 7 |
| NUM-GEOM-066 | SphericalCapArea | spherical cap | curved area | 7 |
| NUM-GEOM-067 | TorusVolume | ring torus | volume | — |
| NUM-GEOM-068 | TorusSurfaceArea | ring torus | surface area | — |
| NUM-GEOM-069 | EllipsoidVolume | ellipsoid | volume | 2 |
| NUM-GEOM-070 | EllipsoidSurfaceArea | ellipsoid | surface area, exact (Carlson) | 2 |
| NUM-GEOM-071 | TetrahedronVolume | regular tetrahedron | volume | — |
| NUM-GEOM-072 | TetrahedronSurfaceArea | regular tetrahedron | surface area | — |
| NUM-GEOM-073 | CubeVolume | cube | volume | — |
| NUM-GEOM-074 | CubeSurfaceArea | cube | surface area | — |
| NUM-GEOM-075 | OctahedronVolume | regular octahedron | volume | — |
| NUM-GEOM-076 | OctahedronSurfaceArea | regular octahedron | surface area | — |
| NUM-GEOM-077 | DodecahedronVolume | regular dodecahedron | volume | 47 (clip) |
| NUM-GEOM-078 | DodecahedronSurfaceArea | regular dodecahedron | surface area | 47 (clip) |
| NUM-GEOM-079 | IcosahedronVolume | regular icosahedron | volume | — |
| NUM-GEOM-080 | IcosahedronSurfaceArea | regular icosahedron | surface area | — |
| NUM-GEOM-081 | HelixArcLength | circular helix | arc length | 26, 31 (via life/home) |
| NUM-GEOM-082 | HelixHeight | circular helix | height from arc length | 25 |

Deliberately NOT separate functions (one function, one home):
- Square area/perimeter: RectangleArea / RectanglePerimeter with equal sides, or the regular 4-gon.
- Parallelogram and rhombus area by base × height: identical formula to RectangleArea (tagged `parallelogram`, `rhombus`).
- Rhombus and square area by diagonals: KiteArea (any quadrilateral with perpendicular diagonals).
- Trapezoid, kite, irregular quadrilateral perimeter: a sum of sides, or PolygonPerimeter over vertices.
- Right-triangle acute angle from two legs: SlopeAngleRad (opposite = rise, adjacent = run).
- Cone slant height: RightTriangleHypotenuse(radius, height).
- Regular polygon perimeter: n × side.
- Ellipsoid surface by Thomsen's approximation: SOURCE NOT FOUND (see NUM-GEOM-070 notes); no ID.

Candidates noticed, not in this pass (LATER, each needs Devin's yes): polygon centroid and second moment of a
polygon (pure geometry; the beam-section rows in `eng/mechanical` would then call them), polygon orientation
(signed area), point-in-polygon, circle through three points (circumradius), inradius, 3D point distance.

## STEP 2: one entry per function

Source short forms used below (full citations, each read or confirmed 2026-11-12, are in PENDING_SHARED.md
under "doc/sources/ files to create"):
- **CRC32**: Zwillinger (ed.), *CRC Standard Mathematical Tables and Formulae*, 32nd ed., CRC Press, 2011.
- **BSMM**: Bronshtein, Semendyayev, Musiol, Mühlig, *Handbook of Mathematics*, 6th ed., Springer, 2015.
- **DLMF**: NIST Digital Library of Mathematical Functions, Release 1.2.x, https://dlmf.nist.gov (read 2026-11-12).
- **Kahan2014**: W. Kahan, "Miscalculating Area and Angles of a Needle-like Triangle", UC Berkeley lecture
  notes, work in progress dated September 4, 2014, https://people.eecs.berkeley.edu/~wkahan/Triangle.pdf (read 2026-11-12).
- **MW**: Weisstein, MathWorld, https://mathworld.wolfram.com (secondary pointer only, never the primary).

Locator confidence (no locator is presented as more certain than it is):
- CRC32 §4.6.2 (p. 210), §4.10.3 (p. 227) are CONFIRMED by the author's 32nd-edition errata
  (mathtable.com/errata/smtf32_errata.pdf, read 2026-11-12). **†** = INFERRED: the 32nd edition's chapter 4 section
  order (publisher contents, read 2026-11-12) counted from 4.1 and anchored by the errata. To be checked against a
  copy before build. The 31st and 33rd editions number differently; only the 32nd is cited.
- **‡** = BSMM subsection from a contents snippet whose page numbers match Crossref for chapter 3 (doi:10.1007/978-3-662-46221-8_3,
  pp. 129–268); §3.5 and §3.6 numbers are from memory. To be checked against a copy before build.
- DLMF equation numbers 19.8.x, 19.9.9, 19.9.10, 19.33.1–3 were read on the pages; 19.21.10 and 4.19.1 were not.
- Kahan2014 §2, §3 and Table 1 rows were read in the PDF itself.

Every function: O(1) time and O(1) space unless stated; no allocation; NaN or ±Inf in → `ErrNotFinite`; a
finite input whose result is ±Inf → `ErrOverflow`. These two lines are not repeated per entry. "HAND-CHECKED"
means the value was computed from the cited formula here (double precision, `doc/geometry_2026-11-12/`) and
agrees with an independent route named in the entry.

### PointDistance      (ID reserved: NUM-GEOM-001)
Formula: d = √((x₂−x₁)² + (y₂−y₁)²).
Inputs and units: `p, q Point`, coordinates in one length unit.
Output: `float64` distance, same unit.
Errors: none beyond the shared two.
Primary source: CRC32 §4.5 (lines)† (distance between two points); BSMM §3.5 (analytical geometry)‡.
Worked example: (0,0), (3,4) → 5. HAND-CHECKED (3-4-5 triangle).
Variants: none.
Best approach: `math.Hypot(dx, dy)`, which avoids overflow and underflow of the squares. O(1).
Notes: must NOT square and add directly (overflows at |dx| > 1.3e154).

### Midpoint      (ID reserved: NUM-GEOM-002)
Formula: M = ((x₁+x₂)/2, (y₁+y₂)/2).
Inputs and units: `p, q Point`, one length unit.
Output: `Point`, same unit.
Errors: shared only.
Primary source: CRC32 §4.5 (lines)† (midpoint of a segment); BSMM §3.5 (analytical geometry)‡.
Worked example: (0,0), (4,−2) → (2,−1). HAND-CHECKED.
Variants: none.
Best approach: x₁ + (x₂−x₁)/2 overflows when signs differ; x₁/2 + x₂/2 never overflows and loses at most
one ulp to underflow at subnormals. Use x₁/2 + x₂/2. O(1).
Notes: none.

### Slope      (ID reserved: NUM-GEOM-003)
Formula: m = (y₂−y₁)/(x₂−x₁).
Inputs and units: `p, q Point`, one length unit.
Output: `float64`, dimensionless (rise over run).
Errors: `ErrDegenerate` if p == q; `ErrVertical` (new sentinel) if x₂ == x₁ and y₂ ≠ y₁.
Primary source: CRC32 §4.5 (lines)† (slope of a line through two points); BSMM §3.5 (analytical geometry)‡.
Worked example: (1,1), (3,5) → 2. HAND-CHECKED.
Variants: none.
Best approach: one subtraction each, one division. O(1).
Notes: must NOT return ±Inf for a vertical line; SlopeAngleRad handles vertical lines.

### SlopeAngleRad      (ID reserved: NUM-GEOM-004)
Formula: θ = atan2(rise, run), the angle of the line from the +x axis.
Inputs and units: `rise, run float64`, same length unit (or a ratio with run = 1).
Output: `float64` radians in (−π, π].
Errors: `ErrDegenerate` if rise == run == 0.
Primary source: CRC32 §4.5 (lines)† (angle of inclination, tan θ = m); CRC32 §4.2.4 (polar coordinates)† (the quadrant-correct angle).
Worked example: rise 1, run 1 → π/4. HAND-CHECKED.
Variants: grade in percent = 100·rise/run is NOT this function (a ratio, `life/home` or `eng/civil`).
Best approach: `math.Atan2`, exact in the quadrant and finite for run = 0. O(1).
Notes: covers input row 29 and the right-triangle acute angle (opposite = rise, adjacent = run).

### PointLineDistance      (ID reserved: NUM-GEOM-005)
Formula: for the line through A and B, d = |(Bx−Ax)(Py−Ay) − (By−Ay)(Px−Ax)| / |B−A|.
Inputs and units: `p, a, b Point`, one length unit.
Output: `float64` perpendicular distance from p to the infinite line AB, same unit.
Errors: `ErrDegenerate` if a == b.
Primary source: CRC32 §4.5 (lines)† (distance from a point to a line); MW "Point-Line Distance--2-Dimensional" (pointer).
Worked example: P (0,0), A (0,1), B (1,0) → √2/2 ≈ 0.707106781187. HAND-CHECKED (equal to half the diagonal of the unit square).
Variants: distance to the SEGMENT (clamped) is a different function, not in this pass.
Best approach: the cross-product form above with `math.Hypot` in the denominator; the a·x+b·y+c form loses
digits when the line is far from the origin. O(1).
Notes: none.

### AngleBetweenLinesRad      (ID reserved: NUM-GEOM-006)
Formula: with direction vectors u = B−A and v = D−C, θ = atan2(|u×v|, |u·v|), the acute angle in [0, π/2].
Inputs and units: `a, b, c, d Point` (line 1 through a, b; line 2 through c, d), one length unit.
Output: `float64` radians in [0, π/2].
Errors: `ErrDegenerate` if a == b or c == d.
Primary source: CRC32 §4.5 (lines)† (angle between two lines, tan θ = abs((m₂−m₁)/(1+m₁m₂))).
Worked example: line y = 0 and line through (0,0), (√3, 1) → π/6. HAND-CHECKED.
Variants: from two slopes is the source's form; it fails for vertical lines and perpendicular lines
(1 + m₁m₂ = 0), so the vector form is used.
Best approach: atan2 of cross and dot: accurate at both 0 and π/2, where acos(dot/(|u||v|)) loses half its digits. O(1).
Notes: input row 28 "Interface Angle" is mapped here BY OUR READING of the name only; Devin to confirm.

### PolarToCartesian      (ID reserved: NUM-GEOM-007)
Formula: x = r cos θ, y = r sin θ.
Inputs and units: `r float64` (length, ≥ 0), `thetaRad float64`.
Output: `Point`, same length unit.
Errors: `ErrNegative` if r < 0.
Primary source: CRC32 §4.2.4 (polar coordinates)†.
Worked example: r 2, θ π/6 → (√3, 1) ≈ (1.732050807569, 1). HAND-CHECKED.
Variants: none.
Best approach: `math.Sincos`, one call. O(1).
Notes: θ is not reduced; very large θ loses accuracy in `math.Sincos` (documented, not an error).

### CartesianToPolar      (ID reserved: NUM-GEOM-008)
Formula: r = √(x²+y²), θ = atan2(y, x).
Inputs and units: `p Point`, one length unit.
Output: `r float64` (same unit), `thetaRad float64` in (−π, π].
Errors: none beyond the shared two; the origin returns (0, 0) by the atan2 convention, documented.
Primary source: CRC32 §4.2.4 (polar coordinates)†.
Worked example: (−1, 0) → (1, π). HAND-CHECKED.
Variants: θ in [0, 2π) is a caller's adjustment, not a second function.
Best approach: `math.Hypot`, `math.Atan2`. O(1).
Notes: −0 for y gives θ = −π from atan2; the function maps it to +π so the range is (−π, π].

### ComplementaryAngleRad      (ID reserved: NUM-GEOM-009)
Formula: θc = π/2 − θ.
Inputs and units: `angleRad float64` in [0, π/2].
Output: `float64` radians in [0, π/2].
Errors: `ErrOutOfRange` outside [0, π/2].
Primary source: BSMM §3.1 (plane geometry, p. 129)‡ (complementary and supplementary angles); a CRC32 line was not located (to look for at build;
BSMM alone suffices for a definition).
Worked example: π/6 → π/3. HAND-CHECKED.
Variants: none.
Best approach: one subtraction. O(1).
Notes: must NOT accept negative or obtuse angles silently.

### SupplementaryAngleRad      (ID reserved: NUM-GEOM-010)
Formula: θs = π − θ.
Inputs and units: `angleRad float64` in [0, π].
Output: `float64` radians in [0, π].
Errors: `ErrOutOfRange` outside [0, π].
Primary source: BSMM §3.1 (plane geometry, p. 129)‡.
Worked example: π/3 → 2π/3. HAND-CHECKED.
Variants: none.
Best approach: one subtraction. O(1).
Notes: none.

### PolygonInteriorAngleSumRad      (ID reserved: NUM-GEOM-011)
Formula: S = (n − 2)π.
Inputs and units: `sides int`.
Output: `float64` radians.
Errors: `ErrTooFewSides` if n < 3.
Primary source: CRC32 §4.6.3 (regular polygons)†; BSMM §3.1 (plane geometry, p. 129)‡ (polygons).
Worked example: n 5 → 3π. HAND-CHECKED.
Variants: none.
Best approach: one multiply; n converted to float64 exactly for n < 2⁵³. O(1).
Notes: none.

### RegularPolygonInteriorAngleRad      (ID reserved: NUM-GEOM-012)
Formula: α = (n − 2)π/n.
Inputs and units: `sides int`.
Output: `float64` radians.
Errors: `ErrTooFewSides` if n < 3.
Primary source: CRC32 §4.6.3 (regular polygons)†.
Worked example: n 6 → 2π/3. HAND-CHECKED.
Variants: none.
Best approach: compute π − 2π/n (one rounding fewer). O(1).
Notes: none.

### RegularPolygonExteriorAngleRad      (ID reserved: NUM-GEOM-013)
Formula: β = 2π/n.
Inputs and units: `sides int`.
Output: `float64` radians.
Errors: `ErrTooFewSides` if n < 3.
Primary source: CRC32 §4.6.3 (regular polygons)†.
Worked example: n 8 → π/4. HAND-CHECKED.
Variants: none.
Best approach: one division. O(1).
Notes: none.

### TrianglePerimeter      (ID reserved: NUM-GEOM-014)
Formula: P = a + b + c.
Inputs and units: `a, b, c float64`, one length unit.
Output: `float64`, same unit.
Errors: `ErrNonPositive` if any side ≤ 0; `ErrNotTriangle` if Kahan's test c − (a − b) < 0 fails after sorting.
Primary source: CRC32 §4.6.1 (triangles)†; triangle test: Kahan2014 §2.
Worked example: 3, 4, 5 → 12. HAND-CHECKED.
Variants: none.
Best approach: sort three values, test, add. O(1).
Notes: the same validity test as TriangleAreaHeron, so the two agree on what a triangle is.

### TriangleAreaBaseHeight      (ID reserved: NUM-GEOM-015)
Formula: A = ½·b·h.
Inputs and units: `base, height float64`, one length unit.
Output: `float64` area, unit squared.
Errors: `ErrNonPositive` if either ≤ 0.
Primary source: CRC32 §4.6.1 (triangles)†.
Worked example: base 6, height 4 → 12. HAND-CHECKED.
Variants: none.
Best approach: one multiply, one halving. O(1).
Notes: none.

### TriangleAreaHeron      (ID reserved: NUM-GEOM-016)
Formula: sort so a ≥ b ≥ c; A = ¼√((a+(b+c))(c−(a−b))(c+(a−b))(a+(b−c))). The parentheses are required.
Inputs and units: `a, b, c float64`, one length unit.
Output: `float64` area, unit squared.
Errors: `ErrNonPositive` if any side ≤ 0; `ErrNotTriangle` if c − (a − b) < 0.
Primary source: Kahan2014 §2 "How to compute ∆"; Heron's formula itself: CRC32 §4.6.1 (triangles)†; MW "Heron's Formula" (pointer).
Worked example: Kahan2014 Table 1 (§5, p. 4): (10, 10, 10) → 43.30127020; (100000, 99999.99979, 0.00029) →
9.999999990; (10000, 5000.000001, 15000) → 612.3724358. Plus (3,4,5) → 6, (13,14,15) → 84, HAND-CHECKED.
Variants: classical √(s(s−a)(s−b)(s−c)) is the same quantity but loses all digits for needle triangles (Kahan2014 §1).
Best approach: Kahan's form above: error a few ulps for every triangle. O(1).
Notes: must NOT use s = (a+b+c)/2. A degenerate triangle (c − (a−b) == 0) returns 0, no error.

### TriangleAreaSAS      (ID reserved: NUM-GEOM-017)
Formula: A = ½·a·b·sin C.
Inputs and units: `a, b float64` (one length unit), `includedAngleRad float64` in (0, π).
Output: `float64` area, unit squared.
Errors: `ErrNonPositive` for a side ≤ 0; `ErrOutOfRange` if the angle ∉ (0, π).
Primary source: CRC32 §4.6.1 (triangles)†.
Worked example: 3, 4, π/6 → 3. HAND-CHECKED.
Variants: none.
Best approach: `math.Sin`, two multiplies. O(1).
Notes: none.

### TriangleSideSAS      (ID reserved: NUM-GEOM-018)
Formula: law of cosines c² = a² + b² − 2ab cos C, computed as c = √((a−b)² + 4ab·sin²(C/2)).
Inputs and units: `a, b float64` (one length unit), `includedAngleRad` in (0, π).
Output: `float64` side, same unit.
Errors: `ErrNonPositive`, `ErrOutOfRange` as above.
Primary source: law of cosines CRC32 §4.6.1 (triangles)†; half-angle identity cos C = 1 − 2sin²(C/2): DLMF §4.21(iii) (multiples of the argument).
Worked example: 3, 4, π/3 → √13 ≈ 3.605551275464. HAND-CHECKED (cos π/3 = ½: 9 + 16 − 12 = 13).
Variants: none.
Best approach: the half-angle form: both terms are non-negative, so no cancellation when C is small and
a ≈ b (the textbook form loses all digits there). O(1).
Notes: the rewrite is an identity, not an approximation; the test includes a needle case checked against Kahan2014 Table 1 row 2 (C from §3).

### TriangleAngleSSSRad      (ID reserved: NUM-GEOM-019)
Formula: angle C opposite c. With a ≥ b: μ = c − (a − b) if b ≥ c ≥ 0, else μ = b − (a − c) if c > b ≥ 0;
C = 2·atan(√(((a−b)+c)·μ / ((a+(b+c))·((a−c)+b)))).
Inputs and units: `a, b, c float64` (one length unit); the result is the angle opposite `c`.
Output: `float64` radians in [0, π].
Errors: `ErrNonPositive`; `ErrNotTriangle`.
Primary source: Kahan2014 §3 "How to compute C"; law of cosines CRC32 §4.6.1 (triangles)†.
Worked example: Kahan2014 Table 1: (10,10,10) → 60°; (100000, 99999.99979, 0.00029) → 1.145915591E-7°;
(10000, 5000.000001, 15000) → 179.9985965°. Plus (4, 3, 5) → π/2, HAND-CHECKED.
Variants: acos((a²+b²−c²)/2ab) loses digits near 0 and π; not offered.
Best approach: Kahan's form. O(1).
Notes: the caller passes a, b in either order; the function swaps so a ≥ b, as §3 requires.

### RightTriangleHypotenuse      (ID reserved: NUM-GEOM-020)
Formula: c = √(a² + b²).
Inputs and units: `legA, legB float64`, one length unit.
Output: `float64`, same unit.
Errors: `ErrNonPositive`.
Primary source: CRC32 §4.6.1 (triangles)† (Pythagorean theorem).
Worked example: 3, 4 → 5. HAND-CHECKED.
Variants: none.
Best approach: `math.Hypot`. O(1).
Notes: none.

### RightTriangleLeg      (ID reserved: NUM-GEOM-021)
Formula: b = √(c² − a²), computed as √((c−a)(c+a)).
Inputs and units: `hypotenuse, leg float64`, one length unit.
Output: `float64`, same unit.
Errors: `ErrNonPositive`; `ErrNotTriangle` if leg ≥ hypotenuse.
Primary source: CRC32 §4.6.1 (triangles)† (Pythagorean theorem).
Worked example: 5, 3 → 4. HAND-CHECKED.
Variants: none.
Best approach: the factored form: no cancellation when leg ≈ hypotenuse, no overflow of c². O(1).
Notes: none.

### RectangleArea      (ID reserved: NUM-GEOM-022)
Formula: A = w·h.
Inputs and units: `width, height float64`, one length unit.
Output: `float64` area, unit squared.
Errors: `ErrNonPositive`.
Primary source: CRC32 §4.6.2 (quadrilaterals, p. 210).
Worked example: 3, 4 → 12. HAND-CHECKED.
Variants: square (w = h); parallelogram and rhombus by base × perpendicular height (same formula, CRC32 §4.6.2 (quadrilaterals, p. 210)).
Best approach: one multiply. O(1).
Notes: tags carry `square`, `parallelogram`, `rhombus` so INDEX searches find it.

### RectanglePerimeter      (ID reserved: NUM-GEOM-023)
Formula: P = 2(a + b).
Inputs and units: `sideA, sideB float64`, one length unit.
Output: `float64`, same unit.
Errors: `ErrNonPositive`.
Primary source: CRC32 §4.6.2 (quadrilaterals, p. 210).
Worked example: 3, 4 → 14. HAND-CHECKED.
Variants: any quadrilateral with two pairs of equal sides: square, parallelogram, rhombus, kite (CRC32 §4.6.2 (quadrilaterals, p. 210)).
Best approach: one add, one doubling. O(1).
Notes: none.

### ParallelogramAreaSAS      (ID reserved: NUM-GEOM-024)
Formula: A = a·b·sin θ.
Inputs and units: `sideA, sideB float64` (one length unit), `angleRad` in (0, π).
Output: `float64` area, unit squared.
Errors: `ErrNonPositive`; `ErrOutOfRange`.
Primary source: CRC32 §4.6.2 (quadrilaterals, p. 210).
Worked example: 3, 4, π/6 → 6. HAND-CHECKED.
Variants: rhombus (a = b).
Best approach: `math.Sin`, two multiplies. O(1).
Notes: none.

### TrapezoidArea      (ID reserved: NUM-GEOM-025)
Formula: A = ½(a + b)·h, a and b the parallel sides.
Inputs and units: `baseA, baseB, height float64`, one length unit.
Output: `float64` area, unit squared.
Errors: `ErrNonPositive` (a base of 0 is a triangle: use TriangleAreaBaseHeight).
Primary source: CRC32 §4.6.2 (quadrilaterals, p. 210).
Worked example: 3, 5, 2 → 8. HAND-CHECKED.
Variants: none.
Best approach: O(1).
Notes: none.

### KiteArea      (ID reserved: NUM-GEOM-026)
Formula: A = ½·d₁·d₂.
Inputs and units: `diagonal1, diagonal2 float64`, one length unit.
Output: `float64` area, unit squared.
Errors: `ErrNonPositive`.
Primary source: CRC32 §4.6.2 (quadrilaterals, p. 210); MW "Kite" (pointer).
Worked example: 6, 4 → 12. HAND-CHECKED.
Variants: any quadrilateral whose diagonals are perpendicular: rhombus, square.
Best approach: O(1).
Notes: none.

### RegularPolygonAreaFromSide      (ID reserved: NUM-GEOM-027)
Formula: A = n·s² / (4·tan(π/n)).
Inputs and units: `sides int`, `side float64` (length unit).
Output: `float64` area, unit squared.
Errors: `ErrTooFewSides`; `ErrNonPositive`.
Primary source: CRC32 §4.6.3 (regular polygons)†; MW "Regular Polygon" (pointer).
Worked example: n 6, s 1 → 3√3/2 ≈ 2.598076211353. HAND-CHECKED (six equilateral triangles of area √3/4).
Variants: by apothem (-028), by circumradius (-029): same shape, other defining length.
Best approach: `math.Tan(math.Pi/n)`. O(1).
Notes: none.

### RegularPolygonAreaFromApothem      (ID reserved: NUM-GEOM-028)
Formula: A = n·a²·tan(π/n), a = apothem (inradius).
Inputs and units: `sides int`, `apothem float64`.
Output: `float64` area, unit squared.
Errors: `ErrTooFewSides`; `ErrNonPositive`.
Primary source: CRC32 §4.6.3 (regular polygons)†.
Worked example: n 6, a √3/2 → 2.598076211353. HAND-CHECKED (agrees with -027).
Variants: see -027.
Best approach: O(1).
Notes: covers input row 10.

### RegularPolygonAreaFromCircumradius      (ID reserved: NUM-GEOM-029)
Formula: A = ½·n·R²·sin(2π/n).
Inputs and units: `sides int`, `circumradius float64`.
Output: `float64` area, unit squared.
Errors: `ErrTooFewSides`; `ErrNonPositive`.
Primary source: CRC32 §4.6.3 (regular polygons)†.
Worked example: n 6, R 1 → 2.598076211353. HAND-CHECKED (agrees with -027).
Variants: see -027.
Best approach: O(1).
Notes: none.

### RegularPolygonSideInscribed      (ID reserved: NUM-GEOM-030)
Formula: s = 2R·sin(π/n), the side of a regular n-gon inscribed in a circle of radius R.
Inputs and units: `sides int`, `circumradius float64`.
Output: `float64` side, same unit.
Errors: `ErrTooFewSides`; `ErrNonPositive`.
Primary source: CRC32 §4.6.3 (regular polygons)†.
Worked example: n 4, R 1 → √2. HAND-CHECKED. n 6, R 1 → 1.
Variants: none.
Best approach: O(1).
Notes: covers input row 39 ("within a circle").

### RegularPolygonSideCircumscribed      (ID reserved: NUM-GEOM-031)
Formula: s = 2r·tan(π/n), the side of a regular n-gon circumscribed about a circle of radius r.
Inputs and units: `sides int`, `inradius float64`.
Output: `float64` side, same unit.
Errors: `ErrTooFewSides`; `ErrNonPositive`.
Primary source: CRC32 §4.6.3 (regular polygons)†.
Worked example: n 4, r 1 → 2. HAND-CHECKED.
Variants: none.
Best approach: O(1).
Notes: covers input row 38 ("around a circle").

### PolygonArea      (ID reserved: NUM-GEOM-032)
Formula: shoelace (surveyor's) formula A = ½|Σᵢ (xᵢ·yᵢ₊₁ − xᵢ₊₁·yᵢ)|, indices mod n.
Inputs and units: `vertices []Point` in order (either orientation), one length unit; the polygon closes itself
(the last vertex is NOT repeated).
Output: `float64` area, unit squared.
Errors: `ErrDegenerate` for fewer than 3 vertices or zero area (all collinear); shared two.
Primary source: CRC32 §4.6 (polygons)†; B. Braden, "The Surveyor's Area Formula", The College Mathematics Journal
17(4), 1986, 326–337, doi:10.2307/2686282 (DOI to confirm at build); MW "Polygon Area" (pointer).
Worked example: (0,0), (4,0), (4,3), (2,5), (0,3) → 16. HAND-CHECKED (rectangle 12 + triangle 4).
Variants: signed area (orientation) is a LATER candidate; this returns |A|.
Best approach: translate every vertex by vertex 0 first (removes the cancellation when coordinates are far
from the origin), then sum the cross terms with Neumaier compensated summation (Neumaier 1974, already in
`doc/sources/neumaier_1974.md`), fused into the one loop. O(n) time, O(1) space, no allocation.
Notes: self-intersection is NOT detected: detecting it costs O(n log n) (Shamos–Hoey) and more than the area
itself; a self-intersecting polygon is a documented precondition, and the function returns the signed-sum
magnitude. The input slice is not modified.

### PolygonPerimeter      (ID reserved: NUM-GEOM-033)
Formula: P = Σᵢ |Vᵢ₊₁ − Vᵢ|, indices mod n.
Inputs and units: `vertices []Point`, one length unit, not repeating the first vertex.
Output: `float64`, same unit.
Errors: `ErrDegenerate` for fewer than 2 distinct vertices; shared two.
Primary source: CRC32 §4.6 (polygons)† (perimeter as the sum of side lengths).
Worked example: the pentagon of -032 → 12 + 2√8 ≈ 15.656854249492. HAND-CHECKED (4 + 3 + 2√2 + 2√2 + 3).
Variants: open polyline length is a different function, not in this pass.
Best approach: `math.Hypot` per edge, Neumaier summation. O(n) time, O(1) space.
Notes: none.

### CircleCircumference      (ID reserved: NUM-GEOM-034)
Formula: C = 2πr.
Inputs and units: `radius float64`.
Output: `float64`, same unit.
Errors: `ErrNonPositive`.
Primary source: CRC32 §4.10.6 (circles)†; BSMM §3.1.6 (the circle and related shapes, p. 139)‡.
Worked example: r 1 → 2π ≈ 6.283185307180. HAND-CHECKED.
Variants: from diameter is the caller's halving, not a function.
Best approach: O(1).
Notes: none.

### CircleArea      (ID reserved: NUM-GEOM-035)
Formula: A = πr².
Inputs and units: `radius float64`.
Output: `float64` area, unit squared.
Errors: `ErrNonPositive`.
Primary source: CRC32 §4.10.6 (circles)†; BSMM §3.1.6 (the circle and related shapes, p. 139)‡.
Worked example: r 2 → 4π ≈ 12.566370614359. HAND-CHECKED.
Variants: none.
Best approach: O(1).
Notes: none.

### AnnulusArea      (ID reserved: NUM-GEOM-036)
Formula: A = π(R² − r²), computed as π(R − r)(R + r).
Inputs and units: `outerRadius, innerRadius float64`, one length unit.
Output: `float64` area, unit squared.
Errors: `ErrNonPositive`; `ErrOutOfRange` if inner ≥ outer.
Primary source: CRC32 §4.10.6 (circles)†; BSMM §3.1.6 (the circle and related shapes, p. 139)‡; MW "Annulus" (pointer).
Worked example: 5, 3 → 16π ≈ 50.265482457437. HAND-CHECKED.
Variants: none.
Best approach: factored form: no cancellation for thin rings, no overflow of R². O(1).
Notes: covers input row 9.

### ArcLength      (ID reserved: NUM-GEOM-037)
Formula: s = r·θ.
Inputs and units: `radius float64`, `angleRad float64` in (0, 2π].
Output: `float64`, same unit as radius.
Errors: `ErrNonPositive`; `ErrOutOfRange`.
Primary source: CRC32 §4.10.6 (circles)†; BSMM §3.1.6 (the circle and related shapes, p. 139)‡.
Worked example: r 2, π/3 → 2π/3 ≈ 2.094395102393. HAND-CHECKED.
Variants: none.
Best approach: O(1).
Notes: θ > 2π is rejected, not wrapped.

### ChordLength      (ID reserved: NUM-GEOM-038)
Formula: c = 2r·sin(θ/2).
Inputs and units: `radius float64`, `centralAngleRad` in (0, 2π].
Output: `float64`, same unit.
Errors: `ErrNonPositive`; `ErrOutOfRange`.
Primary source: CRC32 §4.10.6 (circles)†; BSMM §3.1.6 (the circle and related shapes, p. 139)‡.
Worked example: r 2, π/3 → 2. HAND-CHECKED (equilateral triangle).
Variants: none.
Best approach: the half-angle sine form (the law-of-cosines form loses digits for small θ). O(1).
Notes: none.

### SectorArea      (ID reserved: NUM-GEOM-039)
Formula: A = ½·r²·θ.
Inputs and units: `radius float64`, `angleRad` in (0, 2π].
Output: `float64` area, unit squared.
Errors: `ErrNonPositive`; `ErrOutOfRange`.
Primary source: CRC32 §4.10.6 (circles)†; BSMM §3.1.6 (the circle and related shapes, p. 139)‡.
Worked example: r 2, π/3 → 2π/3 ≈ 2.094395102393. HAND-CHECKED.
Variants: none.
Best approach: O(1).
Notes: none.

### SegmentArea      (ID reserved: NUM-GEOM-040)
Formula: A = ½·r²·(θ − sin θ).
Inputs and units: `radius float64`, `centralAngleRad` in (0, 2π].
Output: `float64` area, unit squared.
Errors: `ErrNonPositive`; `ErrOutOfRange`.
Primary source: CRC32 §4.10.6 (circles)†; BSMM §3.1.6 (the circle and related shapes, p. 139)‡; MW "Circular Segment" (pointer).
Worked example: r 2, π/3 → 2(π/3 − √3/2) ≈ 0.362344294824; r 1, π → π/2. HAND-CHECKED.
Variants: from chord and sagitta is a caller's composition with -041.
Best approach: θ − sin θ cancels catastrophically for small θ (θ = 1e-4 keeps ~8 digits). For θ < 0.5 use the
Maclaurin series θ³/6 − θ⁵/120 + θ⁷/5040 − … (DLMF §4.19, eq. 4.19.1, the sine series), summed until the term
is below 2⁻⁵³ of the sum (≤ 8 terms at θ = 0.5); otherwise the direct form. O(1).
Notes: the series switch point is measured in verify, not guessed; the test includes θ = 1e-6.

### ArcRadiusFromChordSagitta      (ID reserved: NUM-GEOM-041)
Formula: R = c²/(8s) + s/2, c = chord, s = sagitta (the rise of the arc at the chord's midpoint).
Inputs and units: `chord, sagitta float64`, one length unit.
Output: `float64` radius, same unit.
Errors: `ErrNonPositive`.
Primary source: CRC32 §4.10.6 (circles)† (relations between chord, sagitta and radius); MW "Sagitta" (pointer).
Worked example: chord 8, sagitta 2 → 5. HAND-CHECKED (the 3-4-5 triangle: 5 − 2 = 3, half-chord 4).
Variants: sagitta from radius and chord is LATER.
Best approach: O(1). Both terms positive: no cancellation.
Notes: covers input row 41 "Radius of an Arch". sagitta > chord/2 is a valid major arc, not an error.

### CircleTangentLength      (ID reserved: NUM-GEOM-042)
Formula: t = √(d² − r²), computed as √((d − r)(d + r)), d = distance from the external point to the centre.
Inputs and units: `radius, distanceFromCenter float64`, one length unit.
Output: `float64` tangent length, same unit.
Errors: `ErrNonPositive`; `ErrOutOfRange` if d < r (point inside the circle); d == r → 0.
Primary source: CRC32 §4.10.6 (circles)† (power of a point / tangent length); second source to find at build (power of a point).
Worked example: r 6 371 000, d r + 2 → √(2·6371000·2 + 4) ≈ 5048.168380706. HAND-CHECKED.
Variants: none.
Best approach: factored form: for a point just above a planet (d − r ≪ r) the textbook form loses
all digits of d² − r². O(1).
Notes: the horizon distance in `sci/astronomy` calls this with the cited Earth radius; this function knows no planet.

### EllipseArea      (ID reserved: NUM-GEOM-043)
Formula: A = π·a·b, a and b the semi-axes.
Inputs and units: `semiAxisA, semiAxisB float64`, one length unit, either order.
Output: `float64` area, unit squared.
Errors: `ErrNonPositive`.
Primary source: CRC32 §4.10.3 (ellipses, p. 227); BSMM §3.5 (curves of second order)‡.
Worked example: 5, 3 → 15π ≈ 47.123889803847. HAND-CHECKED.
Variants: none.
Best approach: O(1).
Notes: none.

### EllipseChordFromCenter      (ID reserved: NUM-GEOM-044)
Formula: from x²/a² + y²/b² = 1, the chord parallel to the a-axis at offset y from the centre is
L = 2a·√(1 − y²/b²), computed as (2a/b)·√((b − y)(b + y)).
Inputs and units: `semiAxisAlong` (a, the semi-axis parallel to the chord), `semiAxisAcross` (b),
`offsetFromCenter` (|y| ≤ b), one length unit.
Output: `float64` chord length, same unit.
Errors: `ErrNonPositive` for a semi-axis ≤ 0; `ErrNegative` for offset < 0; `ErrOutOfRange` for offset > b.
Primary source: equation of the ellipse in standard position: CRC32 §4.10.3 (ellipses, p. 227); BSMM §3.5 (curves of second order)‡.
Worked example: a 5, b 3, offset 1.8 → 8. HAND-CHECKED (√(1 − 0.36) = 0.8).
Variants: horizontal vs vertical chord is only which semi-axis is passed as "along" (input rows 19 and 21).
Best approach: factored form (no cancellation near the vertex). O(1).
Notes: none.

### EllipseChordFromEdge      (ID reserved: NUM-GEOM-045)
Formula: depth h measured from the vertex on the b-axis: y = b − h, L = (2a/b)·√(h·(2b − h)).
Inputs and units: `semiAxisAlong`, `semiAxisAcross`, `depthFromEdge` (0 ≤ h ≤ 2b), one length unit.
Output: `float64` chord length, same unit.
Errors: as -044, with `ErrOutOfRange` for h > 2b.
Primary source: as -044.
Worked example: a 5, b 3, depth 1.2 → 8 (same chord as -044's example). HAND-CHECKED.
Variants: input rows 20 and 22.
Best approach: h·(2b − h) directly: exact near the vertex where b − y would cancel. O(1).
Notes: none.

### EllipsePerimeter      (ID reserved: NUM-GEOM-046)
Formula: L = 4a·E(k), k² = 1 − b²/a², a ≥ b (DLMF 19.9.9); E(k) from the arithmetic-geometric mean:
K(k) = π/(2·M(1, k′)) (DLMF 19.8.5) and E(k) = K(k)·(1 − Σₙ 2ⁿ⁻¹cₙ²) with a₀ = 1, g₀ = k′ (DLMF 19.8.6, cₙ per
19.8.2). Equivalent single-loop form: L = 2π·N(a², b²)/M(a, b), N the modified AGM (S. Adlaj, "An Eloquent
Formula for the Perimeter of an Ellipse", Notices of the AMS 59(8), 2012, 1094–1099, doi:10.1090/noti879,
DOI to confirm at build).
Inputs and units: `semiAxisA, semiAxisB float64`, one length unit, either order.
Output: `float64` perimeter, same unit.
Errors: `ErrNonPositive`.
Primary source: DLMF §19.9 eq. 19.9.9; DLMF §19.8 eqs. 19.8.1, 19.8.2, 19.8.5, 19.8.6 (formulas to be re-read
on the page at build: a research summary garbled 19.8.5).
Worked example: a 1, b 1 → 2π; a 2, b 1 → 9.688448220547636 (±2e-15); a 5, b 3 → 25.526998863398. HAND-CHECKED
three ways: AGM/MAGM, 200 000-point midpoint quadrature of the arc-length integral, and Ramanujan (-047)
agreeing to its stated error.
Variants: exact (this) vs approximate (-047).
Best approach: scale to a = 1, b = t ≤ 1 (homogeneous of degree 1), run AGM and MAGM until |x − y| ≤ 2⁻⁵⁰·x
(quadratic convergence: ≤ 8 iterations for t ≥ 1e-9 in the probe), multiply back by a. For t < 2⁻²⁶ return 4a
(the next term is O(a·t²·ln(1/t)), below one ulp). O(log log(1/ε)) iterations, O(1) space.
Notes: must NOT run unscaled (the probe at a = 1e6, b = 1 failed to converge without scaling); must cap the
iterations and return `ErrNoConvergence` (new sentinel) rather than loop. The circle (a == b) returns 2πa exactly.

### EllipsePerimeterRamanujan      (ID reserved: NUM-GEOM-047)
Formula: L ≈ π·(a + b)·(1 + 3λ²/(10 + √(4 − 3λ²))), λ = (a − b)/(a + b) (Ramanujan's second approximation).
Inputs and units: as -046.
Output: `float64` perimeter, same unit.
Errors: `ErrNonPositive`.
Primary source: S. Ramanujan, "Modular equations and approximations to π", Quarterly Journal of Pure and
Applied Mathematics 45, 1914, 350–372, §15 p. 372 ("obtained empirically"; error ε "about 3a·e²⁰/2³⁶", e the
eccentricity); DLMF §19.9 eq. 19.9.10 (the same formula plus its leading error term 3λ¹⁰/2¹⁷, citing
G. Almkvist and B. Berndt, Amer. Math. Monthly 95(7), 1988, 585–608, doi:10.2307/2323302).
Worked example: a 2, b 1 → 9.688448216130 (exact 9.688448220548, relative error −4.6e-10); a 1, b 0⁺ limit →
3.998390650 vs exact 4 (relative error −4.0e-4, the worst case). HAND-CHECKED against -046.
Variants: Ramanujan's first approximation π[3(a+b) − √((a+3b)(3a+b))] (same §15, error ≈ a·e¹²/2²⁰) is not
offered: the second is better everywhere at the same cost.
Best approach: λ² form. O(1).
Notes: Error bound (documented in the doc comment): relative error ≤ 4.1e-4 for all a, b > 0 (worst as b/a → 0,
HAND-CHECKED), and ≈ 3λ¹⁰/2¹⁷ for small λ (DLMF 19.9.10). Exists for callers that want the closed form; -046 is
the default.

### PrismVolume      (ID reserved: NUM-GEOM-048)
Formula: V = A·h, A the base area, h the perpendicular height (Cavalieri: holds for oblique prisms).
Inputs and units: `baseArea` (length unit squared), `height` (length unit).
Output: `float64` volume, unit cubed.
Errors: `ErrNonPositive`.
Primary source: CRC32 §4.18 (polyhedra)†; BSMM §3.3.3 (polyhedra, p. 153)‡.
Worked example: base 12, height 5 → 60. HAND-CHECKED.
Variants: any base: the caller passes an area from -022…-032.
Best approach: O(1).
Notes: none.

### CylinderVolume      (ID reserved: NUM-GEOM-049)
Formula: V = πr²h.
Inputs and units: `radius, height float64`, one length unit.
Output: `float64` volume, unit cubed.
Errors: `ErrNonPositive`.
Primary source: CRC32 §4.19 (cylinders)†; BSMM §3.3.4 (solids bounded by curved surfaces, p. 156)‡.
Worked example: r 2, h 5 → 20π ≈ 62.831853071796. HAND-CHECKED.
Variants: oblique cylinder: same formula with perpendicular height.
Best approach: O(1).
Notes: none.

### CylinderLateralArea      (ID reserved: NUM-GEOM-050)
Formula: S = 2πrh.
Inputs and units: as -049.
Output: `float64` area, unit squared.
Errors: `ErrNonPositive`.
Primary source: CRC32 §4.19 (cylinders)†; BSMM §3.3.4 (solids bounded by curved surfaces, p. 156)‡.
Worked example: r 2, h 5 → 20π ≈ 62.831853071796. HAND-CHECKED.
Variants: right circular cylinder only.
Best approach: O(1).
Notes: none.

### CylinderSurfaceArea      (ID reserved: NUM-GEOM-051)
Formula: S = 2πr(h + r).
Inputs and units: as -049.
Output: `float64` area, unit squared.
Errors: `ErrNonPositive`.
Primary source: CRC32 §4.19 (cylinders)†; BSMM §3.3.4 (solids bounded by curved surfaces, p. 156)‡.
Worked example: r 2, h 5 → 28π ≈ 87.964594300514. HAND-CHECKED.
Variants: none.
Best approach: factored form, one multiply chain. O(1).
Notes: none.

### HollowCylinderVolume      (ID reserved: NUM-GEOM-052)
Formula: V = π(R² − r²)h = π(R − r)(R + r)h.
Inputs and units: `outerRadius, innerRadius, height float64`, one length unit.
Output: `float64` volume, unit cubed.
Errors: `ErrNonPositive`; `ErrOutOfRange` if inner ≥ outer.
Primary source: CRC32 §4.19 (cylinders)† with the annulus of CRC32 §4.10.6 (circles)† (prism of an annular base); BSMM §3.3.4 (solids bounded by curved surfaces, p. 156)‡.
Worked example: 3, 2, 5 → 25π ≈ 78.539816339745. HAND-CHECKED.
Variants: none.
Best approach: factored form. O(1).
Notes: none.

### HollowCylinderSurfaceArea      (ID reserved: NUM-GEOM-053)
Formula: S = 2π(R + r)h + 2π(R − r)(R + r) (outer and inner walls plus the two annular ends).
Inputs and units: as -052.
Output: `float64` area, unit squared.
Errors: as -052.
Primary source: CRC32 §4.19 (cylinders)† (lateral area) + CRC32 §4.10.6 (circles)† (annulus).
Worked example: 3, 2, 5 → 60π ≈ 188.495559215388. HAND-CHECKED.
Variants: none.
Best approach: 2π(R + r)(h + R − r), one product. O(1).
Notes: none.

### PyramidVolume      (ID reserved: NUM-GEOM-054)
Formula: V = ⅓·A·h.
Inputs and units: `baseArea` (unit squared), `height` (unit).
Output: `float64` volume, unit cubed.
Errors: `ErrNonPositive`.
Primary source: CRC32 §4.18 (polyhedra)†; BSMM §3.3.3 (polyhedra, p. 153)‡.
Worked example: base 36 (square, side 6), height 4 → 48. HAND-CHECKED.
Variants: any base; the cone is the circular case (-058).
Best approach: O(1).
Notes: covers input row 40 with -055.

### RegularPyramidSurfaceArea      (ID reserved: NUM-GEOM-055)
Formula: S = base + ½·n·s·ℓ, base = n s²/(4 tan(π/n)), slant height ℓ = √(h² + apothem²), apothem = s/(2 tan(π/n)).
Inputs and units: `sides int`, `side, height float64` (one length unit).
Output: `float64` total area, unit squared.
Errors: `ErrTooFewSides`; `ErrNonPositive`.
Primary source: CRC32 §4.18 (polyhedra)†; BSMM §3.3.3 (polyhedra, p. 153)‡; base area from CRC32 §4.6.3 (regular polygons)†.
Worked example: n 4, s 6, h 4 → 36 + 4·(6·5/2) = 96. HAND-CHECKED (slant √(16 + 9) = 5).
Variants: right pyramid with a regular base only.
Best approach: `math.Hypot` for ℓ; reuses -027. O(1).
Notes: none.

### PyramidFrustumVolume      (ID reserved: NUM-GEOM-056)
Formula: V = (h/3)(A₁ + A₂ + √(A₁A₂)).
Inputs and units: `baseArea1, baseArea2` (unit squared, parallel and similar), `height` (unit).
Output: `float64` volume, unit cubed.
Errors: `ErrNonPositive` (a top area of 0 is a pyramid: use -054).
Primary source: CRC32 §4.18 (polyhedra)†; BSMM §3.3.3 (polyhedra, p. 153)‡; MW "Pyramidal Frustum" (pointer).
Worked example: 16, 4, 3 → 28. HAND-CHECKED (agrees with -057 for the same solid).
Variants: the cone frustum is the circular case (-061).
Best approach: √(A₁)·√(A₂) to avoid overflow of the product. O(1).
Notes: none.

### PrismoidVolume      (ID reserved: NUM-GEOM-057)
Formula: prismoidal formula V = (h/6)(A₁ + 4·Aₘ + A₂), Aₘ the area of the mid-height section.
Inputs and units: `area1, areaMid, area2` (unit squared), `height` (unit).
Output: `float64` volume, unit cubed.
Errors: `ErrNonPositive` for height; `ErrNegative` for an area < 0 (an end area may be 0: a wedge or pyramid).
Primary source: CRC32 §4.18 (polyhedra)† (prismatoid); BSMM §3.3.3 (polyhedra, p. 153)‡; MW "Prismoid" (pointer).
Worked example: square frustum, end sides 4 and 2, h 3: A₁ 16, Aₘ 9, A₂ 4 → 28. HAND-CHECKED (equals -056).
Variants: exact for any solid whose cross-section area is a polynomial of degree ≤ 3 in height (prismatoids);
it is Simpson's rule, so for other solids it is an approximation, documented.
Best approach: O(1).
Notes: covers input rows 8 and 43.

### ConeVolume      (ID reserved: NUM-GEOM-058)
Formula: V = ⅓πr²h.
Inputs and units: `radius, height float64`, one length unit.
Output: `float64` volume, unit cubed.
Errors: `ErrNonPositive`.
Primary source: CRC32 §4.20 (cones)†; BSMM §3.3.4 (solids bounded by curved surfaces, p. 156)‡.
Worked example: r 3, h 4 → 12π ≈ 37.699111843078. HAND-CHECKED.
Variants: oblique cone: same formula.
Best approach: O(1).
Notes: none.

### ConeLateralArea      (ID reserved: NUM-GEOM-059)
Formula: S = πr·ℓ, ℓ = √(r² + h²).
Inputs and units: as -058.
Output: `float64` area, unit squared.
Errors: `ErrNonPositive`.
Primary source: CRC32 §4.20 (cones)†; BSMM §3.3.4 (solids bounded by curved surfaces, p. 156)‡.
Worked example: r 3, h 4 → 15π ≈ 47.123889803847. HAND-CHECKED.
Variants: right circular cone only.
Best approach: `math.Hypot`. O(1).
Notes: none.

### ConeSurfaceArea      (ID reserved: NUM-GEOM-060)
Formula: S = πr(r + ℓ).
Inputs and units: as -058.
Output: `float64` area, unit squared.
Errors: `ErrNonPositive`.
Primary source: CRC32 §4.20 (cones)†; BSMM §3.3.4 (solids bounded by curved surfaces, p. 156)‡.
Worked example: r 3, h 4 → 24π ≈ 75.398223686155. HAND-CHECKED.
Variants: none.
Best approach: O(1).
Notes: none.

### ConeFrustumVolume      (ID reserved: NUM-GEOM-061)
Formula: V = (πh/3)(R² + Rr + r²).
Inputs and units: `radius1, radius2, height float64`, one length unit, radii in either order.
Output: `float64` volume, unit cubed.
Errors: `ErrNonPositive` (radius 0 is a cone: -058).
Primary source: CRC32 §4.20 (cones)†; BSMM §3.3.4 (solids bounded by curved surfaces, p. 156)‡; MW "Conical Frustum" (pointer).
Worked example: R 4, r 1, h 4 → 28π ≈ 87.964594300514. HAND-CHECKED.
Variants: none.
Best approach: O(1).
Notes: none.

### ConeFrustumSurfaceArea      (ID reserved: NUM-GEOM-062)
Formula: S = π(R + r)ℓ + πR² + πr², ℓ = √((R − r)² + h²).
Inputs and units: as -061.
Output: `float64` total area, unit squared.
Errors: as -061.
Primary source: CRC32 §4.20 (cones)†; BSMM §3.3.4 (solids bounded by curved surfaces, p. 156)‡.
Worked example: R 4, r 1, h 4 → 42π ≈ 131.946891450771. HAND-CHECKED (ℓ = 5).
Variants: lateral only = first term; not a separate function in this pass.
Best approach: `math.Hypot`. O(1).
Notes: none.

### SphereVolume      (ID reserved: NUM-GEOM-063)
Formula: V = (4/3)πr³.
Inputs and units: `radius float64`.
Output: `float64` volume, unit cubed.
Errors: `ErrNonPositive`.
Primary source: CRC32 §4.8.1 (spheres; §4.8 quadrics)†; BSMM §3.3.4 (solids bounded by curved surfaces, p. 156)‡.
Worked example: r 3 → 36π ≈ 113.097335529233. HAND-CHECKED.
Variants: none.
Best approach: O(1).
Notes: none.

### SphereSurfaceArea      (ID reserved: NUM-GEOM-064)
Formula: S = 4πr².
Inputs and units: `radius float64`.
Output: `float64` area, unit squared.
Errors: `ErrNonPositive`.
Primary source: CRC32 §4.8.1 (spheres; §4.8 quadrics)†; BSMM §3.3.4 (solids bounded by curved surfaces, p. 156)‡.
Worked example: r 3 → 36π ≈ 113.097335529233. HAND-CHECKED.
Variants: none.
Best approach: O(1).
Notes: none.

### SphericalCapVolume      (ID reserved: NUM-GEOM-065)
Formula: V = (πh²/3)(3r − h), r the sphere's radius, h the cap height.
Inputs and units: `sphereRadius, capHeight float64`, one length unit, 0 < h ≤ 2r.
Output: `float64` volume, unit cubed.
Errors: `ErrNonPositive`; `ErrOutOfRange` for h > 2r.
Primary source: CRC32 §4.8.1 (spheres; §4.8 quadrics)†; BSMM §3.3.4 (solids bounded by curved surfaces, p. 156)‡; MW "Spherical Cap" (pointer).
Worked example: r 5, h 2 → 52π/3 ≈ 54.454272662223. HAND-CHECKED; h = 2r gives the sphere (−063).
Variants: from the cap's base radius is LATER.
Best approach: O(1).
Notes: none.

### SphericalCapArea      (ID reserved: NUM-GEOM-066)
Formula: S = 2πrh (curved surface only).
Inputs and units: as -065.
Output: `float64` area, unit squared.
Errors: as -065.
Primary source: CRC32 §4.8.1 (spheres; §4.8 quadrics)†; BSMM §3.3.4 (solids bounded by curved surfaces, p. 156)‡.
Worked example: r 5, h 2 → 20π ≈ 62.831853071796. HAND-CHECKED.
Variants: + the base disc is the caller's sum.
Best approach: O(1).
Notes: none.

### TorusVolume      (ID reserved: NUM-GEOM-067)
Formula: V = 2π²Rr², R the distance from the axis to the tube's centre, r the tube radius.
Inputs and units: `majorRadius, minorRadius float64`, one length unit.
Output: `float64` volume, unit cubed.
Errors: `ErrNonPositive`; `ErrOutOfRange` if r > R (a self-intersecting spindle torus; the formula no longer
gives the enclosed volume).
Primary source: CRC32 §4.7 (the torus)†; BSMM §3.3.4 (solids bounded by curved surfaces, p. 156)‡; MW "Torus" (pointer).
Worked example: R 3, r 1 → 6π² ≈ 59.217626406536. HAND-CHECKED (Pappus: 2πR·πr²).
Variants: horn torus r = R allowed.
Best approach: O(1).
Notes: none.

### TorusSurfaceArea      (ID reserved: NUM-GEOM-068)
Formula: S = 4π²Rr.
Inputs and units: as -067.
Output: `float64` area, unit squared.
Errors: as -067.
Primary source: CRC32 §4.7 (the torus)†; BSMM §3.3.4 (solids bounded by curved surfaces, p. 156)‡.
Worked example: R 3, r 1 → 12π² ≈ 118.435252813072. HAND-CHECKED (Pappus: 2πR·2πr).
Variants: none.
Best approach: O(1).
Notes: none.

### EllipsoidVolume      (ID reserved: NUM-GEOM-069)
Formula: V = (4/3)πabc.
Inputs and units: `semiAxisA, semiAxisB, semiAxisC float64`, one length unit, any order.
Output: `float64` volume, unit cubed.
Errors: `ErrNonPositive`.
Primary source: CRC32 §4.8.1 (spheres; §4.8 quadrics)† (quadrics); DLMF §19.33 (V as used in 19.33.1); MW "Ellipsoid" (pointer).
Worked example: 3, 2, 1 → 8π ≈ 25.132741228718. HAND-CHECKED.
Variants: none.
Best approach: O(1).
Notes: none.

### EllipsoidSurfaceArea      (ID reserved: NUM-GEOM-070)
Formula: S = 3V·R_G(a⁻², b⁻², c⁻²), V = 4πabc/3 (DLMF 19.33.1); R_G from Carlson's symmetric integrals:
2R_G(x,y,z) = z·R_F(x,y,z) − ⅓(x − z)(y − z)·R_D(x,y,z) + √(xy/z) (DLMF 19.21.10, to re-read at build), with
R_F and R_D by Carlson's duplication algorithm.
Inputs and units: `semiAxisA, semiAxisB, semiAxisC float64`, one length unit, any order.
Output: `float64` area, unit squared.
Errors: `ErrNonPositive`; `ErrNoConvergence` (cap on duplication steps; never expected).
Primary source: DLMF §19.33 eq. 19.33.1 (and the Legendre form 19.33.2–19.33.3); DLMF §19.21 eq. 19.21.10;
B. C. Carlson, "Numerical computation of real or complex elliptic integrals", Numerical Algorithms 10, 1995,
13–26, doi:10.1007/BF02198293, §2 (duplication algorithms for R_F, R_D) and §3 (check values).
Worked example: Carlson 1995 §3 check values for R_F and R_D (test of the internal helpers); sphere 1,1,1 → 4π;
prolate 2,1,1 → 21.478435327884 (equals the closed-form spheroid area 2πb²(1 + (a/(b·e))·arcsin e), MW "Prolate
Spheroid", pointer); 3, 2, 1 → 48.882146302582. HAND-CHECKED: Carlson route vs 640 000-point surface
quadrature (agree to 1e-6, the quadrature's own error).
Variants: the closed forms for spheroids are special cases, not separate functions. Thomsen's approximation
(p ≈ 1.6075, ≈ ±1.06 %) is NOT offered: SOURCE NOT FOUND (its only origin is a 2004 personal e-mail reported on
numericana.com; no peer-reviewed or textbook source).
Best approach: Carlson duplication: each step cuts the error by 4⁶; stop when the spread ≤ (3ε)^(1/6)·mean per
Carlson §2; ≈ 10 steps in double precision. O(1) (bounded iterations), O(1) space. R_F, R_D are unexported
helpers in `num/geometry` until a second consumer needs them (see PENDING_SHARED.md, taxonomy question).
Notes: sort a ≥ b ≥ c and scale by a (homogeneous) before evaluating.

### TetrahedronVolume      (ID reserved: NUM-GEOM-071)
Formula: V = a³/(6√2) = a³√2/12, a the edge.
Inputs and units: `edge float64`.
Output: `float64` volume, unit cubed.
Errors: `ErrNonPositive`.
Primary source: CRC32 §4.18.1 (convex regular polyhedra)† (table of convex regular polyhedra); BSMM §3.3.3 (polyhedra, p. 153)‡; MW "Platonic Solid" (pointer).
Worked example: a 1 → 0.117851130198. HAND-CHECKED.
Variants: none.
Best approach: precomputed coefficient constant (written as the exact expression, evaluated at compile time). O(1).
Notes: same table and approach for -072…-080.

### TetrahedronSurfaceArea      (ID reserved: NUM-GEOM-072)
Formula: S = √3·a².
Inputs and units: `edge float64`.
Output: `float64` area, unit squared.
Errors: `ErrNonPositive`.
Primary source: CRC32 §4.18.1 (convex regular polyhedra)†; BSMM §3.3.3 (polyhedra, p. 153)‡.
Worked example: a 1 → 1.732050807569. HAND-CHECKED (4 equilateral triangles of √3/4).
Variants: none. Best approach: O(1). Notes: none.

### CubeVolume      (ID reserved: NUM-GEOM-073)
Formula: V = a³.
Inputs and units: `edge float64`.
Output: `float64` volume, unit cubed.
Errors: `ErrNonPositive`.
Primary source: CRC32 §4.18.1 (convex regular polyhedra)†; BSMM §3.3.3 (polyhedra, p. 153)‡.
Worked example: a 2 → 8. HAND-CHECKED.
Variants: none. Best approach: O(1). Notes: none.

### CubeSurfaceArea      (ID reserved: NUM-GEOM-074)
Formula: S = 6a².
Inputs and units: `edge float64`.
Output: `float64` area, unit squared.
Errors: `ErrNonPositive`.
Primary source: CRC32 §4.18.1 (convex regular polyhedra)†; BSMM §3.3.3 (polyhedra, p. 153)‡.
Worked example: a 2 → 24. HAND-CHECKED.
Variants: none. Best approach: O(1). Notes: none.

### OctahedronVolume      (ID reserved: NUM-GEOM-075)
Formula: V = (√2/3)·a³.
Inputs and units: `edge float64`.
Output: `float64` volume, unit cubed.
Errors: `ErrNonPositive`.
Primary source: CRC32 §4.18.1 (convex regular polyhedra)†; BSMM §3.3.3 (polyhedra, p. 153)‡.
Worked example: a 1 → 0.471404520791. HAND-CHECKED (two square pyramids, base 1, height √2/2).
Variants: none. Best approach: O(1). Notes: none.

### OctahedronSurfaceArea      (ID reserved: NUM-GEOM-076)
Formula: S = 2√3·a².
Inputs and units: `edge float64`.
Output: `float64` area, unit squared.
Errors: `ErrNonPositive`.
Primary source: CRC32 §4.18.1 (convex regular polyhedra)†; BSMM §3.3.3 (polyhedra, p. 153)‡.
Worked example: a 1 → 3.464101615138. HAND-CHECKED (8 equilateral triangles).
Variants: none. Best approach: O(1). Notes: none.

### DodecahedronVolume      (ID reserved: NUM-GEOM-077)
Formula: V = ((15 + 7√5)/4)·a³.
Inputs and units: `edge float64`.
Output: `float64` volume, unit cubed.
Errors: `ErrNonPositive`.
Primary source: CRC32 §4.18.1 (convex regular polyhedra)†; BSMM §3.3.3 (polyhedra, p. 153)‡.
Worked example: a 1 → 7.663118960625. HAND-CHECKED against MW "Platonic Solid" (pointer) table.
Variants: none. Best approach: O(1). Notes: none.

### DodecahedronSurfaceArea      (ID reserved: NUM-GEOM-078)
Formula: S = 3√(25 + 10√5)·a².
Inputs and units: `edge float64`.
Output: `float64` area, unit squared.
Errors: `ErrNonPositive`.
Primary source: CRC32 §4.18.1 (convex regular polyhedra)†; BSMM §3.3.3 (polyhedra, p. 153)‡.
Worked example: a 1 → 20.645728807068. HAND-CHECKED (12 × RegularPolygonAreaFromSide(5, 1) = 12 × 1.720477400589).
Variants: none. Best approach: O(1). Notes: none.

### IcosahedronVolume      (ID reserved: NUM-GEOM-079)
Formula: V = (5(3 + √5)/12)·a³.
Inputs and units: `edge float64`.
Output: `float64` volume, unit cubed.
Errors: `ErrNonPositive`.
Primary source: CRC32 §4.18.1 (convex regular polyhedra)†; BSMM §3.3.3 (polyhedra, p. 153)‡.
Worked example: a 1 → 2.181694990625. HAND-CHECKED against MW "Platonic Solid" (pointer) table.
Variants: none. Best approach: O(1). Notes: none.

### IcosahedronSurfaceArea      (ID reserved: NUM-GEOM-080)
Formula: S = 5√3·a².
Inputs and units: `edge float64`.
Output: `float64` area, unit squared.
Errors: `ErrNonPositive`.
Primary source: CRC32 §4.18.1 (convex regular polyhedra)†; BSMM §3.3.3 (polyhedra, p. 153)‡.
Worked example: a 1 → 8.660254037844. HAND-CHECKED (20 equilateral triangles).
Variants: none. Best approach: O(1). Notes: none.

### HelixArcLength      (ID reserved: NUM-GEOM-081)
Formula: circular helix x = r cos t, y = r sin t, z = c·t with pitch p = 2πc; one turn has length
√((2πr)² + p²), so L = n·√((2πr)² + p²) for n turns.
Inputs and units: `radius, pitch float64` (one length unit), `turns float64` (> 0, fractional allowed).
Output: `float64` length, same unit.
Errors: `ErrNonPositive` (pitch 0 is a circle: n·CircleCircumference).
Primary source: CRC32 §4.21.1 (curves)† (arc length of a space curve, s = ∫‖r′(t)‖dt, applied to the helix); BSMM §3.6 (differential geometry)‡;
MW "Helix" (pointer). The helix as a worked case in a primary source is NOT CONFIRMED: see note.
Worked example: r 1, p 2, n 3 → 3·√(4π² + 4) ≈ 19.781449856854. HAND-CHECKED (the helix unrolls to the
hypotenuse of a right triangle with legs 2πr·n and p·n).
Variants: conical helix (garland on a cone) has no elementary closed form; LATER, needs `num/calculus` quadrature.
Best approach: `math.Hypot(2πr, p)·n`. O(1).
Notes: the formula follows from the arc-length integral in one line, and the unrolled-cylinder check is
exact; if Devin requires the helix to appear verbatim in a source, it is SOURCE TO CONFIRM until a copy is checked.

### HelixHeight      (ID reserved: NUM-GEOM-082)
Formula: h = n·p = √(L² − (2πr·n)²), computed as √((L − 2πrn)(L + 2πrn)).
Inputs and units: `radius, arcLength float64` (one length unit), `turns float64`.
Output: `float64` total height, same unit.
Errors: `ErrNonPositive`; `ErrOutOfRange` if L < 2πrn (too short to make n turns).
Primary source: as -081 (the same right-triangle relation solved for the other leg).
Worked example: r 1, L 19.781449856854, n 3 → 6. HAND-CHECKED (inverse of -081's example).
Variants: height as pitch × turns is a multiplication, not a function.
Best approach: factored form. O(1).
Notes: covers input row 25. Reading of vCalc's "Helix Height" (height from length) is ours.

### Bell-shaped (input rows 11–13)      (no ID reserved)
SOURCE NOT FOUND / SHAPE UNDEFINED. The three names ("Arc Length", "Area", "Height") describe one curve, but
which curve (a Gaussian y = H·exp(−x²/2σ²), a raised cosine, a "bell" in architecture or bell-founding) is
vCalc's choice and vCalc is not read. The Gaussian's area (H·σ·√(2π)) is `num/stats`/`num/calculus`
territory and its arc length has no closed form. Devin to name the shape; until then nothing is built.

## STEP 3: order and dependencies

IDs are already in build order (catalog order = build order): primitives first, 2D before 3D, exact before
approximate. Suggested build batches (each one file per concept, `.claude/rules/go-functions.md`):

| Batch | IDs | Files (proposed) | Calls inside calcbench |
|---|---|---|---|
| 0 | — | `doc.go`, `point.go` (the Point type). No `errors.go`: each sentinel lives in the first file that defines it (A-008) | — |
| 1 Points, lines, angles | 001–013 | `point_distance.go`, `midpoint.go`, `slope.go`, `slope_angle.go`, `point_line_distance.go`, `angle_between_lines.go`, `polar.go` (007, 008), `complementary_angle.go`, `supplementary_angle.go`, `polygon_angles.go` (011–013) | none |
| 2 Triangles | 014–021 | `triangle_perimeter.go`, `triangle_area.go` (015–017), `triangle_side_sas.go`, `triangle_angle_sss.go`, `right_triangle.go` (020, 021) | 016, 019 share the Kahan sort-and-test helper (unexported) |
| 3 Quadrilaterals, polygons | 022–033 | `rectangle.go`, `parallelogram_area.go`, `trapezoid_area.go`, `kite_area.go`, `regular_polygon_area.go`, `regular_polygon_side.go`, `polygon_area.go`, `polygon_perimeter.go` | 032, 033 inline Neumaier summation (no import) |
| 4 Circle | 034–042 | `circle.go` (034, 035), `annulus_area.go`, `arc.go` (037–039), `segment_area.go`, `arc_radius.go`, `circle_tangent_length.go` | none |
| 5 Ellipse | 043–047 | `ellipse_area.go`, `ellipse_chord.go`, `ellipse_perimeter.go` (046 + unexported AGM/MAGM), `ellipse_perimeter_ramanujan.go` | 047's tests call 046 |
| 6 Prisms, pyramids, cylinders, cones | 048–062 | `prism_volume.go`, `cylinder.go`, `hollow_cylinder.go`, `pyramid.go`, `prismoid_volume.go`, `cone.go`, `cone_frustum.go` | 052, 053 reuse the annulus form; 055 calls 027 |
| 7 Sphere, torus, ellipsoid | 063–070 | `sphere.go`, `spherical_cap.go`, `torus.go`, `ellipsoid_volume.go`, `ellipsoid_surface_area.go` (+ unexported Carlson R_F, R_D, R_G) | none outside the package |
| 8 Regular polyhedra | 071–080 | `tetrahedron.go`, `cube.go`, `octahedron.go`, `dodecahedron.go`, `icosahedron.go` | none |
| 9 Helix | 081–082 | `helix.go` | none |

Dependencies on other subjects: NONE BLOCKING.
- `num/trig`: not needed. Every angle operation is in Go's `math` (Sin, Cos, Tan, Sincos, Atan, Atan2, Acos,
  Hypot). With a radians-only API there is no degree conversion here. If Devin chooses "both radians and
  degrees", the degree variants call a `num/trig` conversion and become BLOCKED ON trig (request line ready in
  PENDING_SHARED.md).
- S1 algorithms: none called. Compensated summation is inlined in 032/033 (fused with the cross-product loop;
  `algo/window` exports only rolling sums).
- `num/linalg`: none, unless Devin chooses a shared Point/Vec2 type (question 2 in the report).
- Elliptic integrals (R_F, R_D) and AGM: kept unexported here; if `num/calculus` or a special-functions home
  later exports them, 046 and 070 switch to the shared ones with no API change.
- Callers waiting on this catalog (MOVE rows): `sci/astronomy` horizon distance → 042; `life/home` garland →
  081, quilting → 015/022; `eng/mechanical` beam-section rows → a future polygon centroid / second moment (LATER).
