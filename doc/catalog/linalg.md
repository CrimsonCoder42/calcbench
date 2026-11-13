# Catalog: linear algebra (`num/linalg`, SHORT `LINALG`)

Catalog and build pass, 2026-11-12. Input: the vCalc "Vectors" list, NAMES ONLY (no vCalc page was fetched).
Work folder: [../linalg_2026-11-12/](../linalg_2026-11-12/README.md) (decisions, citations, benchmarks, shared rows).
Read, not edited: `doc/catalog/geometry.md` (owns Point, 2D distance, polar ↔ cartesian, arc segments) and
`doc/catalog/trig.md` (owns angle normalisation).

## Part 1a: clean (136 rows; "(x2)" rows and grouped collections are expanded)

| # | Name as given | Decision | Home | Reason |
|---|---|---|---|---|
| 1 | Common: scalar multiplication | KEEP | `num/linalg` → NUM-LINALG-003, NUM-LINALG-013, NUM-LINALG-037 | k·v |
| 2 | Common: scalar division | KEEP | `num/linalg` → NUM-LINALG-004, NUM-LINALG-014, NUM-LINALG-038 | v/k |
| 3 | Common: unit vector | KEEP | `num/linalg` → NUM-LINALG-008, NUM-LINALG-018, NUM-LINALG-045 | v/|v|, zero-vector error |
| 4 | Common: magnitude | KEEP | `num/linalg` → NUM-LINALG-007, NUM-LINALG-017, NUM-LINALG-041 | Euclidean length |
| 5 | Common: addition | KEEP | `num/linalg` → NUM-LINALG-001, NUM-LINALG-011, NUM-LINALG-035 | u + v |
| 6 | Common: subtraction | KEEP | `num/linalg` → NUM-LINALG-002, NUM-LINALG-012, NUM-LINALG-036 | u − v |
| 7 | Common: distance between endpoints | KEEP | `num/linalg` → NUM-LINALG-019, NUM-LINALG-046 | |u − v| (2D point distance is num/geometry's PointDistance) |
| 8 | Common: magnitude of sum | KEEP | `num/linalg` → NUM-LINALG-011, NUM-LINALG-017 | |u + v| is Norm(Add(u, v)); no separate function |
| 9 | Common: dot product | KEEP | `num/linalg` → NUM-LINALG-005, NUM-LINALG-015, NUM-LINALG-039, NUM-LINALG-040 | u·v |
| 10 | Common: cross product | KEEP | `num/linalg` → NUM-LINALG-016, NUM-LINALG-006 | u × v |
| 11 | Common: vector angle | KEEP | `num/linalg` → NUM-LINALG-009, NUM-LINALG-020, NUM-LINALG-047 | Kahan's angle formula |
| 12 | Common: vector area | KEEP | `num/linalg` → NUM-LINALG-026, NUM-LINALG-027 | |u × v| and half of it |
| 13 | Common: vector projection | KEEP | `num/linalg` → NUM-LINALG-023, NUM-LINALG-048 | projection onto a line |
| 14 | Common: vector rotation | KEEP | `num/linalg` → NUM-LINALG-031 | Rodrigues about an axis |
| 15 | Common: vector components 3D (magnitude, unit vector, spherical, cylindrical, angle from each axis) | KEEP | `num/linalg` → NUM-LINALG-017, NUM-LINALG-018, NUM-LINALG-054, NUM-LINALG-052, NUM-LINALG-021 | each part is its own function |
| 16 | 3D Vector Calc | DROP | — | a calculator, not a formula; its operations are the KEEP rows |
| 17 | Collection: LM 10.2 | MOVE | `sci/physics` | an LM (Light and Matter) physics textbook collection, routed as the request says |
| 18 | Collection: LM 10.4 | MOVE | `sci/physics` | an LM (Light and Matter) physics textbook collection, routed as the request says |
| 19 | Collection: LM 10.8 | MOVE | `sci/physics` | an LM (Light and Matter) physics textbook collection, routed as the request says |
| 20 | Collection: LM 11.2 | MOVE | `sci/physics` | an LM (Light and Matter) physics textbook collection, routed as the request says |
| 21 | Collection: LM 13.2 | MOVE | `sci/physics` | an LM (Light and Matter) physics textbook collection, routed as the request says |
| 22 | Collection: LM 13.7 | MOVE | `sci/physics` | an LM (Light and Matter) physics textbook collection, routed as the request says |
| 23 | Collection: LM 14.1 | MOVE | `sci/physics` | an LM (Light and Matter) physics textbook collection, routed as the request says |
| 24 | Collection: LM 14.3 | MOVE | `sci/physics` | an LM (Light and Matter) physics textbook collection, routed as the request says |
| 25 | Collection: LM 14.4 | MOVE | `sci/physics` | an LM (Light and Matter) physics textbook collection, routed as the request says |
| 26 | Collection: LM 14.7 | MOVE | `sci/physics` | an LM (Light and Matter) physics textbook collection, routed as the request says |
| 27 | Collection: LM 15.5 | MOVE | `sci/physics` | an LM (Light and Matter) physics textbook collection, routed as the request says |
| 28 | Collection: LM 15.6 | MOVE | `sci/physics` | an LM (Light and Matter) physics textbook collection, routed as the request says |
| 29 | Collection: LM 15.7 | MOVE | `sci/physics` | an LM (Light and Matter) physics textbook collection, routed as the request says |
| 30 | Collection: LM 15.8 | MOVE | `sci/physics` | an LM (Light and Matter) physics textbook collection, routed as the request says |
| 31 | Collection: LM 19.1 | MOVE | `sci/physics` | an LM (Light and Matter) physics textbook collection, routed as the request says |
| 32 | Collection: LM 7_0 | MOVE | `sci/physics` | an LM (Light and Matter) physics textbook collection, routed as the request says |
| 33 | Collection: LM 7_1 | MOVE | `sci/physics` | an LM (Light and Matter) physics textbook collection, routed as the request says |
| 34 | Collection: LM 7_2 | MOVE | `sci/physics` | an LM (Light and Matter) physics textbook collection, routed as the request says |
| 35 | Collection: LM 7_3 | MOVE | `sci/physics` | an LM (Light and Matter) physics textbook collection, routed as the request says |
| 36 | Collection: LM 7_4 | MOVE | `sci/physics` | an LM (Light and Matter) physics textbook collection, routed as the request says |
| 37 | Collection: LM 7_5 | MOVE | `sci/physics` | an LM (Light and Matter) physics textbook collection, routed as the request says |
| 38 | Collection: LM 7_6 | MOVE | `sci/physics` | an LM (Light and Matter) physics textbook collection, routed as the request says |
| 39 | Collection: LM 8_1 | MOVE | `sci/physics` | an LM (Light and Matter) physics textbook collection, routed as the request says |
| 40 | Collection: LM 8_2 | MOVE | `sci/physics` | an LM (Light and Matter) physics textbook collection, routed as the request says |
| 41 | Collection: LM 8_3 | MOVE | `sci/physics` | an LM (Light and Matter) physics textbook collection, routed as the request says |
| 42 | Collection: LM 8_4 | MOVE | `sci/physics` | an LM (Light and Matter) physics textbook collection, routed as the request says |
| 43 | Collection: LM 8_5 | MOVE | `sci/physics` | an LM (Light and Matter) physics textbook collection, routed as the request says |
| 44 | Collection: LM 9_3 | MOVE | `sci/physics` | an LM (Light and Matter) physics textbook collection, routed as the request says |
| 45 | Collection: LM 9_4 | MOVE | `sci/physics` | an LM (Light and Matter) physics textbook collection, routed as the request says |
| 46 | Physics: Motion and Acceleration | MOVE | `sci/physics` | a physics collection |
| 47 | Quaternion Calculator | DROP | — | a calculator; its operations are the quaternion KEEP rows |
| 48 | Vector Calculator (3D) | DROP | — | a calculator, not a formula |
| 49 | Vector Calculator (3D) | MERGE | → #48 | duplicate |
| 50 | Vector Spaces Collection | DROP | — | a collection, not a formula |
| 51 | 4-Vector Norm | KEEP | `num/linalg` → NUM-LINALG-041 | the n-dimensional norm |
| 52 | 5-Vector Norm | MERGE | `num/linalg` → NUM-LINALG-041 | same function as the 4-vector norm |
| 53 | Angle between Vectors | KEEP | `num/linalg` → NUM-LINALG-020, NUM-LINALG-047 | Kahan's formula |
| 54 | Angle of a Circle Arc from Radius and Chord Depth | MOVE | `num/geometry` | circle geometry, routed as the request says |
| 55 | Angle of Vector (vy, vx) | MOVE | `num/geometry` (NUM-GEOM-008 CartesianToPolar) | the polar angle of a 2D vector; geometry owns polar coordinates |
| 56 | Angle of Vector (|v|, vx) | MOVE | `num/geometry` (NUM-GEOM-008) | the same polar angle from other components |
| 57 | Angle of Vector (|v|, vy) | MOVE | `num/geometry` (NUM-GEOM-008) | the same polar angle from other components |
| 58 | Area between Vectors | KEEP | `num/linalg` → NUM-LINALG-026, NUM-LINALG-027 | |a × b| |
| 59 | Area of a Circle Arc Segment | MOVE | `num/geometry` | circle geometry |
| 60 | Cartesian to Cylindrical | KEEP | `num/linalg` → NUM-LINALG-052 | DLMF eq. 1.5.14 |
| 61 | Cartesian to Cylindrical | MERGE | → #60 | duplicate |
| 62 | Cartesian to Polar | MOVE | `num/geometry` (NUM-GEOM-008) | geometry owns 2D polar |
| 63 | Cartesian to Spherical | KEEP | `num/linalg` → NUM-LINALG-054 | DLMF eq. 1.5.16 |
| 64 | Cartesian to Spherical | MERGE | → #63 | duplicate |
| 65 | Circle - area of an arc segment | MERGE | `num/geometry` (with the arc-segment row above) | duplicate circle-segment area |
| 66 | Component Force 2-D | MOVE | `sci/physics` | force |
| 67 | Component Force 3D | MOVE | `sci/physics` | force |
| 68 | Cylindrical to Cartesian | KEEP | `num/linalg` → NUM-LINALG-053 | DLMF eq. 1.5.14 |
| 69 | Distance (Velocity Vector, Time) | MOVE | `sci/physics` | kinematics, routed as the request says |
| 70 | Ephemeris Generation | MOVE | `sci/astronomy` | orbital mechanics |
| 71 | Force Vector Addition | MOVE | `sci/physics` | force (the maths is Vec2Add/Vec3Add) |
| 72 | Force Vector Addition (angle) | MOVE | `sci/physics` | force |
| 73 | Force Vector Components 2D | MOVE | `sci/physics` | force |
| 74 | Horizontal and Vertical Velocity Components | MOVE | `sci/physics` | velocity |
| 75 | LLA to XYZ | MOVE | `sci/geoscience` | WGS 84 geodetic to ECEF; the standard is NIMA TR8350.2 (EPSG:4979 → 4978) |
| 76 | Magnitude of Acceleration | MOVE | `sci/physics` | acceleration (the maths is Vec3Norm) |
| 77 | Magnitude of Force | MOVE | `sci/physics` | force (the maths is Vec3Norm) |
| 78 | Magnitude of Position | MERGE | `num/linalg` → NUM-LINALG-017 | a position vector's length is the vector norm |
| 79 | Magnitude of Position Vector | MERGE | `num/linalg` → NUM-LINALG-017 | duplicate of the previous row |
| 80 | Magnitude of Velocity Vector | MOVE | `sci/physics` | velocity (the maths is Vec3Norm) |
| 81 | Multiplication of a 2x2 Matrix by a 2x1 Matrix | KEEP | `num/linalg` → NUM-LINALG-058 | matrix times vector |
| 82 | Orbit Position and Velocity | MOVE | `sci/astronomy` | orbital mechanics |
| 83 | Polar to Cartesian | MOVE | `num/geometry` (NUM-GEOM-007) | geometry owns 2D polar |
| 84 | Quaternion Addition | KEEP | `num/linalg` → NUM-LINALG-101 | — |
| 85 | Quaternion Conjugate | KEEP | `num/linalg` → NUM-LINALG-104 | — |
| 86 | Quaternion Inverse | KEEP | `num/linalg` → NUM-LINALG-107 | — |
| 87 | Quaternion Magnitude | KEEP | `num/linalg` → NUM-LINALG-105 | — |
| 88 | Quaternion Multiplication | KEEP | `num/linalg` → NUM-LINALG-103 | Hamilton product |
| 89 | Quaternion of Rotation | KEEP | `num/linalg` → NUM-LINALG-108 | axis-angle |
| 90 | Quaternion Subtraction | KEEP | `num/linalg` → NUM-LINALG-102 | — |
| 91 | Quaternion to Roll Pitch Yaw | KEEP | `num/linalg` → NUM-LINALG-115 | aerospace Z-Y-X |
| 92 | Quaternion Versor | KEEP | `num/linalg` → NUM-LINALG-106 | a versor is a unit quaternion |
| 93 | River Flow Compensation | MOVE | `sci/geoscience` | hydrology/navigation, routed as the request says |
| 94 | Roll Pitch Yaw to Quaternion | KEEP | `num/linalg` → NUM-LINALG-114 | aerospace Z-Y-X |
| 95 | Slant Range (Vector) | MOVE | `sci/astronomy` | observer-to-satellite range; the maths is Vec3Distance |
| 96 | Spherical to Cartesian | KEEP | `num/linalg` → NUM-LINALG-055 | DLMF eq. 1.5.16 |
| 97 | Two Vector Functions | DROP | — | a bundle of the KEEP rows, not a formula |
| 98 | V2 - Vector (x component) | MERGE | `num/geometry` (NUM-GEOM-007 PolarToCartesian) | |v| cos θ |
| 99 | V2 - Vector (y component) | MERGE | `num/geometry` (NUM-GEOM-007) | |v| sin θ |
| 100 | V2 - Vector addition orthogonal vectors (magnitude) | MERGE | `num/linalg` → NUM-LINALG-007 | √(x² + y²) |
| 101 | V2 - Vector addition orthogonal vectors (x,y) | MERGE | `num/linalg` → NUM-LINALG-001 | component addition |
| 102 | V2 - Determinant as Related to the Right Hand Rule | KEEP | `num/linalg` → NUM-LINALG-006 | 2D scalar cross product |
| 103 | V3 - Distance between Vectors | KEEP | `num/linalg` → NUM-LINALG-019 | — |
| 104 | V3 - Mixed Product (cross and dot) | KEEP | `num/linalg` → NUM-LINALG-028 | — |
| 105 | V3 - scalar division | KEEP | `num/linalg` → NUM-LINALG-014 | — |
| 106 | V3 - scalar multiplication | KEEP | `num/linalg` → NUM-LINALG-013 | — |
| 107 | V3 - Unit Vector | KEEP | `num/linalg` → NUM-LINALG-018 | — |
| 108 | V3 - vector addition | KEEP | `num/linalg` → NUM-LINALG-011 | — |
| 109 | V3 - Vector Cross Product | KEEP | `num/linalg` → NUM-LINALG-016 | — |
| 110 | V3 - Vector Dot Product | KEEP | `num/linalg` → NUM-LINALG-015 | — |
| 111 | V3 - Vector Magnitude | KEEP | `num/linalg` → NUM-LINALG-017 | — |
| 112 | V3 - Vector Norm | MERGE | `num/linalg` → NUM-LINALG-017 | same as magnitude |
| 113 | V3 - vector subtraction | KEEP | `num/linalg` → NUM-LINALG-012 | — |
| 114 | V3 - XYZ to LLA | MOVE | `sci/geoscience` | ECEF to WGS 84 geodetic (NIMA TR8350.2) |
| 115 | V3 - |V+U| | MERGE | `num/linalg` → NUM-LINALG-011 + NUM-LINALG-017 | composition |
| 116 | V3 - |V| | MERGE | `num/linalg` → NUM-LINALG-017 | same as magnitude |
| 117 | V3 - Determinant as Related to the Right Hand Rule | KEEP | `num/linalg` → NUM-LINALG-016, NUM-LINALG-028 | cross product as a determinant |
| 118 | vector angle (2D) | KEEP | `num/linalg` → NUM-LINALG-009 | — |
| 119 | Vector Components 3D | MERGE | (the common-functions row above) | duplicate of the page's own list |
| 120 | Vector Norm | MERGE | `num/linalg` → NUM-LINALG-041 | same as the n-vector norm |
| 121 | Vector Normal to a Plane | KEEP | `num/linalg` → NUM-LINALG-030 | — |
| 122 | Vector Projection | KEEP | `num/linalg` → NUM-LINALG-023 | — |
| 123 | Vector Rotation | KEEP | `num/linalg` → NUM-LINALG-031 | — |
| 124 | Vector Sum of Forces (2D) | MOVE | `sci/physics` | force |
| 125 | Unlabeled: cos θ = Δx/|Δr| | MERGE | `num/geometry` (NUM-GEOM-008) | 2D basics already covered |
| 126 | Unlabeled: Δx = |Δr| cos θ | MERGE | `num/geometry` (NUM-GEOM-007) | 2D basics already covered |
| 127 | Unlabeled: Δy = |Δr| sin θ | MERGE | `num/geometry` (NUM-GEOM-007) | 2D basics already covered |
| 128 | Unlabeled: Δx_t = Δx1 + Δx2 | MERGE | `num/linalg` → NUM-LINALG-001 | 2D basics already covered |
| 129 | Unlabeled: Δy_t = Δy1 + Δy2 | MERGE | `num/linalg` → NUM-LINALG-001 | 2D basics already covered |
| 130 | Unlabeled: sin θ = Δy/|Δr| | MERGE | `num/geometry` (NUM-GEOM-008) | 2D basics already covered |
| 131 | Unlabeled: θ = acos(Â·B̂) | MERGE | `num/linalg` → NUM-LINALG-009 | covered, by Kahan's formula instead of arccos |
| 132 | Unlabeled: A·B | MERGE | `num/linalg` → NUM-LINALG-005 | 2D basics already covered |
| 133 | Unlabeled: A+B | MERGE | `num/linalg` → NUM-LINALG-001 | 2D basics already covered |
| 134 | Unlabeled: A−B | MERGE | `num/linalg` → NUM-LINALG-002 | 2D basics already covered |
| 135 | Unlabeled: |r| = √(Δx²+Δy²) | MERGE | `num/linalg` → NUM-LINALG-007 | 2D basics already covered |
| 136 | Unlabeled: |A+B| | MERGE | `num/linalg` → NUM-LINALG-001 + NUM-LINALG-007 | composition |

Counts: KEEP 49 · MERGE 28 · MOVE 54 · DROP 5 · LATER 0 (LATER items of the standard set are listed below).

## Part 1b: the standard set (115 functions, all built)

Vectors 001–051 (Vec2, Vec3, n-D) · coordinates 052–055 · matrices 056–100 (Mat2, Mat3, Mat4, n×n with LU) ·
quaternions 101–115. Nothing in the standard set was held back: every citation was confirmed during the build.

| ID | Function | Does |
|---|---|---|
| NUM-LINALG-001 | Vec2Add | a + b in the plane |
| NUM-LINALG-002 | Vec2Sub | a − b in the plane |
| NUM-LINALG-003 | Vec2Scale | k·v in the plane |
| NUM-LINALG-004 | Vec2Div | v/k in the plane, error at k = 0 |
| NUM-LINALG-005 | Vec2Dot | dot product of two 2D vectors, within 2 ulps |
| NUM-LINALG-006 | Vec2Cross | 2D scalar cross product (determinant, right-hand-rule sign) |
| NUM-LINALG-007 | Vec2Norm | length of a 2D vector without overflow |
| NUM-LINALG-008 | Vec2Normalize | unit vector of a 2D vector, error at zero |
| NUM-LINALG-009 | Vec2AngleBetween | angle between two 2D vectors (Kahan's formula) |
| NUM-LINALG-010 | Vec2Rotate | 2D vector turned counterclockwise by an angle |
| NUM-LINALG-011 | Vec3Add | a + b in space |
| NUM-LINALG-012 | Vec3Sub | a − b in space |
| NUM-LINALG-013 | Vec3Scale | k·v in space |
| NUM-LINALG-014 | Vec3Div | v/k in space, error at k = 0 |
| NUM-LINALG-015 | Vec3Dot | dot product of two 3D vectors |
| NUM-LINALG-016 | Vec3Cross | cross product, each component within 2 ulps |
| NUM-LINALG-017 | Vec3Norm | length of a 3D vector without overflow |
| NUM-LINALG-018 | Vec3Normalize | unit vector of a 3D vector, error at zero |
| NUM-LINALG-019 | Vec3Distance | distance between two points in space |
| NUM-LINALG-020 | Vec3AngleBetween | angle between two 3D vectors (Kahan's formula) |
| NUM-LINALG-021 | Vec3DirectionAngles | angles between a vector and the x, y, z axes |
| NUM-LINALG-022 | Vec3ComponentAlong | scalar component of v along an axis |
| NUM-LINALG-023 | Vec3Project | projection of v onto a line |
| NUM-LINALG-024 | Vec3Reject | component of v perpendicular to a line |
| NUM-LINALG-025 | Vec3Reflect | mirror image of v in a plane through the origin |
| NUM-LINALG-026 | Vec3ParallelogramArea | area of the parallelogram spanned by two vectors |
| NUM-LINALG-027 | Vec3TriangleArea | area of the triangle spanned by two vectors |
| NUM-LINALG-028 | Vec3ScalarTripleProduct | a·(b × c), the signed volume (mixed product) |
| NUM-LINALG-029 | Vec3VectorTripleProduct | a × (b × c) |
| NUM-LINALG-030 | Vec3PlaneNormal | unit normal of the plane through three points |
| NUM-LINALG-031 | Vec3Rotate | v turned about an axis by an angle (Rodrigues) |
| NUM-LINALG-032 | Vec3Lerp | linear interpolation, exact at both ends |
| NUM-LINALG-033 | Vec3Min | componentwise minimum |
| NUM-LINALG-034 | Vec3Max | componentwise maximum |
| NUM-LINALG-035 | VecAdd | a + b in n dimensions |
| NUM-LINALG-036 | VecSub | a − b in n dimensions |
| NUM-LINALG-037 | VecScale | k·v in n dimensions |
| NUM-LINALG-038 | VecDiv | v/k in n dimensions, error at k = 0 |
| NUM-LINALG-039 | VecDot | dot product in n dimensions |
| NUM-LINALG-040 | VecDotCompensated | dot product as if in twice the working precision (Dot2) |
| NUM-LINALG-041 | VecNorm | Euclidean norm, scaled and compensated |
| NUM-LINALG-042 | VecNorm1 | 1-norm, compensated sum |
| NUM-LINALG-043 | VecNormInf | infinity norm (largest |element|) |
| NUM-LINALG-044 | VecNormP | p-norm for p ≥ 1, scaled |
| NUM-LINALG-045 | VecNormalize | unit vector in n dimensions, error at zero |
| NUM-LINALG-046 | VecDistance | Euclidean distance between two points in n dimensions |
| NUM-LINALG-047 | VecAngleBetween | angle between two n-D vectors (Kahan's formula) |
| NUM-LINALG-048 | VecProject | projection of v onto a line in n dimensions |
| NUM-LINALG-049 | VecLerp | linear interpolation in n dimensions, exact at both ends |
| NUM-LINALG-050 | VecMin | element-by-element minimum |
| NUM-LINALG-051 | VecMax | element-by-element maximum |
| NUM-LINALG-052 | CartesianToCylindrical | (x, y, z) to cylindrical (R, φ, z) |
| NUM-LINALG-053 | CylindricalToCartesian | cylindrical (R, φ, z) to (x, y, z) |
| NUM-LINALG-054 | CartesianToSpherical | (x, y, z) to spherical (R, θ polar, φ azimuth) |
| NUM-LINALG-055 | SphericalToCartesian | spherical (R, θ, φ) to (x, y, z) |
| NUM-LINALG-056 | Mat2Identity | 2×2 identity |
| NUM-LINALG-057 | Mat2Mul | 2×2 matrix product, elements within 2 ulps |
| NUM-LINALG-058 | Mat2MulVec | 2×2 matrix times a 2D vector |
| NUM-LINALG-059 | Mat2Transpose | 2×2 transpose |
| NUM-LINALG-060 | Mat2Det | 2×2 determinant within 2 ulps |
| NUM-LINALG-061 | Mat2Inverse | 2×2 inverse, error when singular |
| NUM-LINALG-062 | Mat2Trace | 2×2 trace |
| NUM-LINALG-063 | Mat2Rotation | 2×2 rotation matrix |
| NUM-LINALG-064 | Mat2SymmetricEigenvalues | eigenvalues of a symmetric 2×2 matrix without cancellation |
| NUM-LINALG-065 | Mat3Identity | 3×3 identity |
| NUM-LINALG-066 | Mat3Mul | 3×3 matrix product |
| NUM-LINALG-067 | Mat3MulVec | 3×3 matrix times a 3D vector |
| NUM-LINALG-068 | Mat3Transpose | 3×3 transpose |
| NUM-LINALG-069 | Mat3Det | 3×3 determinant as a triple product |
| NUM-LINALG-070 | Mat3Inverse | 3×3 inverse by cofactors, error when singular |
| NUM-LINALG-071 | Mat3Trace | 3×3 trace |
| NUM-LINALG-072 | Mat3RotationX | rotation matrix about the x axis |
| NUM-LINALG-073 | Mat3RotationY | rotation matrix about the y axis |
| NUM-LINALG-074 | Mat3RotationZ | rotation matrix about the z axis |
| NUM-LINALG-075 | Mat3RotationAxisAngle | rotation matrix about any axis (Rodrigues) |
| NUM-LINALG-076 | Mat3SymmetricEigenvalues | eigenvalues of a symmetric 3×3 matrix (Smith's closed form) |
| NUM-LINALG-077 | Mat4Identity | 4×4 identity |
| NUM-LINALG-078 | Mat4Mul | 4×4 matrix product |
| NUM-LINALG-079 | Mat4MulVec | 4×4 matrix times a 4-vector |
| NUM-LINALG-080 | Mat4Transpose | 4×4 transpose |
| NUM-LINALG-081 | Mat4Det | 4×4 determinant by Laplace expansion |
| NUM-LINALG-082 | Mat4Inverse | 4×4 inverse by cofactors, error when singular |
| NUM-LINALG-083 | Mat4Trace | 4×4 trace |
| NUM-LINALG-084 | Mat4Translation | homogeneous translation matrix |
| NUM-LINALG-085 | Mat4Scale | homogeneous scaling matrix |
| NUM-LINALG-086 | Mat4FromMat3 | 3×3 linear map embedded in a 4×4 homogeneous matrix |
| NUM-LINALG-087 | Mat4TransformPoint | point transformed by a homogeneous matrix (with the w divide) |
| NUM-LINALG-088 | Mat4TransformDirection | direction transformed by a homogeneous matrix (no translation) |
| NUM-LINALG-089 | MatIdentity | n×n identity |
| NUM-LINALG-090 | MatMul | product of two general matrices |
| NUM-LINALG-091 | MatVecMul | general matrix times a vector |
| NUM-LINALG-092 | MatTranspose | transpose of a general matrix |
| NUM-LINALG-093 | MatTrace | trace of an n×n matrix, compensated |
| NUM-LINALG-094 | LUDecompose | PA = LU with partial pivoting |
| NUM-LINALG-095 | LUSolve | solve Ax = b from an LU factorisation |
| NUM-LINALG-096 | MatSolve | solve Ax = b by LU with partial pivoting |
| NUM-LINALG-097 | MatDet | n×n determinant from LU, without intermediate overflow |
| NUM-LINALG-098 | MatInverse | n×n inverse from LU, error when singular |
| NUM-LINALG-099 | MatRank | numerical rank by complete pivoting with a tolerance |
| NUM-LINALG-100 | MatConditionEstimate1 | estimate of the 1-norm condition number (Hager–Higham) |
| NUM-LINALG-101 | QuatAdd | p + q for quaternions |
| NUM-LINALG-102 | QuatSub | p − q for quaternions |
| NUM-LINALG-103 | QuatMul | Hamilton product p·q |
| NUM-LINALG-104 | QuatConjugate | quaternion conjugate |
| NUM-LINALG-105 | QuatNorm | quaternion magnitude without overflow |
| NUM-LINALG-106 | QuatNormalize | unit quaternion (versor), error at zero |
| NUM-LINALG-107 | QuatInverse | quaternion inverse, error at zero |
| NUM-LINALG-108 | QuatFromAxisAngle | unit quaternion of a rotation about an axis |
| NUM-LINALG-109 | QuatToAxisAngle | axis and angle of a quaternion's rotation |
| NUM-LINALG-110 | QuatFromMat3 | unit quaternion of a rotation matrix (Shepperd) |
| NUM-LINALG-111 | QuatToMat3 | rotation matrix of a quaternion |
| NUM-LINALG-112 | QuatRotateVec3 | 3D vector rotated by a quaternion |
| NUM-LINALG-113 | QuatSlerp | spherical linear interpolation between rotations |
| NUM-LINALG-114 | QuatFromEulerZYX | quaternion from yaw, pitch, roll (aerospace Z-Y-X) |
| NUM-LINALG-115 | QuatToEulerZYX | yaw, pitch, roll from a quaternion, gimbal lock handled |

LATER (not built, with the reason):
- General (non-symmetric) eigenvalues, symmetric eigenvectors, SVD: iterative methods (Golub, Van Loan ch. 7–8), a later slice.
- QR / Gram–Schmidt (Strang §4.4 "Orthonormal Bases and Gram-Schmidt"): no consumer yet; classical Gram–Schmidt
  is unstable and Householder QR (Golub, Van Loan ch. 5) is the right first version when least squares is needed.
- Complex vectors and matrices; Cholesky; sparse matrices: not requested.
- A 2D `Vec2Distance`: not built, num/geometry's PointDistance owns it.

## Part 1c: research entries (one per KEEP function, generated from the doc comments, which are the source of truth)

### Vec2Add      (ID reserved: NUM-LINALG-001)
Formula: Vec2Add returns a + b, component by component.
Inputs and units: `(a, b Vec2)`; a and b in one unit; the result in that unit.
Output: `(Vec2, error)`.
Errors: ErrNotFinite for a NaN or ±Inf component; ErrOverflow when a component of the sum is outside the float64 range. The result is the zero vector with any error.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations" [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, 2) + (3, −5) → (4, −3) (hand-checked)
Variants: Vec3 and n-D (Vec) forms exist; 2D points and polar coordinates are num/geometry's.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec2.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec2Sub      (ID reserved: NUM-LINALG-002)
Formula: Vec2Sub returns a − b, component by component.
Inputs and units: `(a, b Vec2)`; a and b in one unit; the result in that unit.
Output: `(Vec2, error)`.
Errors: as Vec2Add.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations" [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, 2) − (3, −5) → (−2, 7) (hand-checked)
Variants: Vec3 and n-D (Vec) forms exist; 2D points and polar coordinates are num/geometry's.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec2.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec2Scale      (ID reserved: NUM-LINALG-003)
Formula: Vec2Scale returns k·v.
Inputs and units: `(v Vec2, k float64)`; v in any unit, k dimensionless; the result in v's unit.
Output: `(Vec2, error)`.
Errors: ErrNotFinite for a NaN or ±Inf input; ErrOverflow when a component of the product is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations" [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: 3·(1, −2) → (3, −6) (hand-checked)
Variants: Vec3 and n-D (Vec) forms exist; 2D points and polar coordinates are num/geometry's.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec2.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec2Div      (ID reserved: NUM-LINALG-004)
Formula: Vec2Div returns v/k.
Inputs and units: `(v Vec2, k float64)`; v in any unit, k dimensionless; the result in v's unit.
Output: `(Vec2, error)`.
Errors: ErrZeroDivisor when k = 0; ErrNotFinite for a NaN or ±Inf input; ErrOverflow when a component of the quotient is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations" [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (3, −6)/3 → (1, −2) (hand-checked)
Variants: Vec3 and n-D (Vec) forms exist; 2D points and polar coordinates are num/geometry's.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec2.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec2Dot      (ID reserved: NUM-LINALG-005)
Formula: Vec2Dot returns the dot product a·b = a₀b₀ + a₁b₁, computed as a 2×2 determinant by Kahan's fused multiply-add algorithm so that it is within 2 ulps even when the two products cancel.
Inputs and units: `(a, b Vec2)`; a in unit A, b in unit B; the result in A·B.
Output: `(float64, error)`.
Errors: ErrNotFinite for a NaN or ±Inf component; ErrOverflow when the product is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.2 "Lengths and Dot Products" — and — Jeannerod, Louvet, Muller, "Further analysis of Kahan's algorithm for the accurate computation of 2×2 determinants", Mathematics of Computation, 2013, 82(284):2245–2264, doi:10.1090/S0025-5718-2013-02679-8 [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, 2)·(3, 4) → 11 (hand-checked)
Variants: Vec3 and n-D (Vec) forms exist; 2D points and polar coordinates are num/geometry's.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec2.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec2Cross      (ID reserved: NUM-LINALG-006)
Formula: Vec2Cross returns the scalar cross product a × b = a₀b₁ − a₁b₀, the determinant of the matrix with rows a and b. It is positive when b is counterclockwise from a (right-hand rule: it is the z component of the 3D cross product), negative when clockwise, zero when parallel; its absolute value is the area of the parallelogram the two vectors span. Computed by Kahan's fused multiply-add algorithm, within 2 ulps.
Inputs and units: `(a, b Vec2)`; a in unit A, b in unit B; the result in A·B.
Output: `(float64, error)`.
Errors: as Vec2Dot.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.3 "Cramer's Rule, Inverses, and Volumes" — and — Jeannerod, Louvet, Muller, "Further analysis of Kahan's algorithm for the accurate computation of 2×2 determinants", Mathematics of Computation, 2013, 82(284):2245–2264, doi:10.1090/S0025-5718-2013-02679-8 [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, 0) × (0, 1) → 1; (3, 1) × (2, 4) → 10 (hand-checked)
Variants: Vec3 and n-D (Vec) forms exist; 2D points and polar coordinates are num/geometry's.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec2.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec2Norm      (ID reserved: NUM-LINALG-007)
Formula: Vec2Norm returns the Euclidean length |v| = √(x² + y²), without overflow or underflow in the squares.
Inputs and units: `(v Vec2)`; v in any unit; the result in that unit.
Output: `(float64, error)`.
Errors: ErrNotFinite for a NaN or ±Inf component; ErrOverflow when the length exceeds the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.2 "Lengths and Dot Products" — and — Blue, "A Portable Fortran Program to Find the Euclidean Norm of a Vector", ACM Transactions on Mathematical Software, 1978, 4(1):15–23, doi:10.1145/355769.355771 [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (3, 4) → 5 (hand-checked)
Variants: Vec3 and n-D (Vec) forms exist; 2D points and polar coordinates are num/geometry's.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec2.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec2Normalize      (ID reserved: NUM-LINALG-008)
Formula: Vec2Normalize returns the unit vector v/|v|.
Inputs and units: `(v Vec2)`; v in any unit; the result is dimensionless.
Output: `(Vec2, error)`.
Errors: ErrZeroVector when v = 0; ErrNotFinite for a NaN or ±Inf component.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.2 "Lengths and Dot Products" [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (3, 4) → (0.6, 0.8) (hand-checked)
Variants: Vec3 and n-D (Vec) forms exist; 2D points and polar coordinates are num/geometry's.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec2.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec2AngleBetween      (ID reserved: NUM-LINALG-009)
Formula: Vec2AngleBetween returns the unoriented angle between a and b, in [0, π] radians, by Kahan's formula θ = 2·atan2(|â − b̂|, |â + b̂|) with â = a/|a|, b̂ = b/|b|. Unlike arccos of the normalised dot product, it keeps full accuracy when the vectors are nearly parallel or antiparallel, and unlike arcsin of the cross product it keeps it near π/2.
Inputs and units: `(a, b Vec2)`; a and b in any units; the result in radians.
Output: `(float64, error)`.
Errors: ErrZeroVector when a or b is zero; ErrNotFinite for a NaN or ±Inf component.
Primary source: Kahan, "How Futile are Mindless Assessments of Roundoff in Floating-Point Computation?", University of California, Berkeley, 11 January 2006, §12 "Mangled Angles", pp. 46–47, https://people.eecs.berkeley.edu/~wkahan/Mindless.pdf (read 2026-11-12) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, 0) and (1, 1) → π/4 (hand-checked)
Variants: Vec3 and n-D (Vec) forms exist; 2D points and polar coordinates are num/geometry's.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec2.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec2Rotate      (ID reserved: NUM-LINALG-010)
Formula: Vec2Rotate returns v turned counterclockwise by angleRad about the origin: the active rotation [cos −sin; sin cos]·v. Each component is a 2×2 determinant computed by Kahan's algorithm.
Inputs and units: `(v Vec2, angleRad float64)`; v in any unit, angleRad in radians; the result in v's unit.
Output: `(Vec2, error)`.
Errors: ErrNotFinite for a NaN or ±Inf input; ErrOverflow when a component of the result is outside the float64 range (possible only for |v| near the float64 maximum).
Primary source: Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §2.4, eq. (16) (the upper-left 2×2 block of R₃, transposed: Diebel's matrices are frame rotations) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, 0) by π/2 → (0, 1) (hand-checked; the computed x is cos(π/2) ≈ 6.1e−17)
Variants: Vec3 and n-D (Vec) forms exist; 2D points and polar coordinates are num/geometry's.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec2.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec3Add      (ID reserved: NUM-LINALG-011)
Formula: Vec3Add returns a + b, component by component.
Inputs and units: `(a, b Vec3)`; a and b in one unit; the result in that unit.
Output: `(Vec3, error)`.
Errors: ErrNotFinite for a NaN or ±Inf component; ErrOverflow when a component of the sum is outside the float64 range. The result is the zero vector with any error.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations" [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, 2, 3) + (4, −5, 6) → (5, −3, 9) (hand-checked)
Variants: Vec2 and n-D (Vec) forms where meaningful.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec3.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec3Sub      (ID reserved: NUM-LINALG-012)
Formula: Vec3Sub returns a − b, component by component.
Inputs and units: `(a, b Vec3)`; a and b in one unit; the result in that unit.
Output: `(Vec3, error)`.
Errors: as Vec3Add.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations" [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, 2, 3) − (4, −5, 6) → (−3, 7, −3) (hand-checked)
Variants: Vec2 and n-D (Vec) forms where meaningful.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec3.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec3Scale      (ID reserved: NUM-LINALG-013)
Formula: Vec3Scale returns k·v.
Inputs and units: `(v Vec3, k float64)`; v in any unit, k dimensionless; the result in v's unit.
Output: `(Vec3, error)`.
Errors: ErrNotFinite for a NaN or ±Inf input; ErrOverflow when a component of the product is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations" [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: 2·(1, −2, 3) → (2, −4, 6) (hand-checked)
Variants: Vec2 and n-D (Vec) forms where meaningful.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec3.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec3Div      (ID reserved: NUM-LINALG-014)
Formula: Vec3Div returns v/k.
Inputs and units: `(v Vec3, k float64)`; v in any unit, k dimensionless; the result in v's unit.
Output: `(Vec3, error)`.
Errors: ErrZeroDivisor when k = 0; ErrNotFinite for a NaN or ±Inf input; ErrOverflow when a component of the quotient is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations" [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (2, −4, 6)/2 → (1, −2, 3) (hand-checked)
Variants: Vec2 and n-D (Vec) forms where meaningful.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec3.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec3Dot      (ID reserved: NUM-LINALG-015)
Formula: Vec3Dot returns the dot product a·b = a₀b₀ + a₁b₁ + a₂b₂, accumulated with fused multiply-adds. For an ill-conditioned dot product (heavy cancellation) use VecDotCompensated.
Inputs and units: `(a, b Vec3)`; a in unit A, b in unit B; the result in A·B.
Output: `(float64, error)`.
Errors: ErrNotFinite for a NaN or ±Inf component; ErrOverflow when the result is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.2 "Lengths and Dot Products" [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, 2, 3)·(4, −5, 6) → 12 (hand-checked)
Variants: Vec2 and n-D (Vec) forms where meaningful.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec3_product.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec3Cross      (ID reserved: NUM-LINALG-016)
Formula: Vec3Cross returns the cross product a × b, perpendicular to both by the right-hand rule, with length |a||b| sin θ. Each component is a 2×2 determinant computed by Kahan's fused multiply-add algorithm, within 2 ulps.
Inputs and units: `(a, b Vec3)`; a in unit A, b in unit B; the result in A·B.
Output: `(Vec3, error)`.
Errors: ErrNotFinite for a NaN or ±Inf component; ErrOverflow when a component is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.3 "Cramer's Rule, Inverses, and Volumes" (cross product) — and — Jeannerod, Louvet, Muller, "Further analysis of Kahan's algorithm for the accurate computation of 2×2 determinants", Mathematics of Computation, 2013, 82(284):2245–2264, doi:10.1090/S0025-5718-2013-02679-8 [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: x̂ × ŷ → ẑ; (1, 2, 3) × (4, 5, 6) → (−3, 6, −3) (hand-checked)
Variants: Vec2 and n-D (Vec) forms where meaningful.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec3_product.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec3Norm      (ID reserved: NUM-LINALG-017)
Formula: Vec3Norm returns the Euclidean length |v| = √(x² + y² + z²), without overflow or underflow in the squares.
Inputs and units: `(v Vec3)`; v in any unit; the result in that unit.
Output: `(float64, error)`.
Errors: ErrNotFinite for a NaN or ±Inf component; ErrOverflow when the length exceeds the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.2 "Lengths and Dot Products" — and — Blue, "A Portable Fortran Program to Find the Euclidean Norm of a Vector", ACM Transactions on Mathematical Software, 1978, 4(1):15–23, doi:10.1145/355769.355771 [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (2, 3, 6) → 7 (hand-checked)
Variants: Vec2 and n-D (Vec) forms where meaningful.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec3_norm.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec3Normalize      (ID reserved: NUM-LINALG-018)
Formula: Vec3Normalize returns the unit vector v/|v|.
Inputs and units: `(v Vec3)`; v in any unit; the result is dimensionless.
Output: `(Vec3, error)`.
Errors: ErrZeroVector when v = 0; ErrNotFinite for a NaN or ±Inf component.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.2 "Lengths and Dot Products" [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (2, 3, 6) → (2/7, 3/7, 6/7) (hand-checked)
Variants: Vec2 and n-D (Vec) forms where meaningful.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec3_norm.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec3Distance      (ID reserved: NUM-LINALG-019)
Formula: Vec3Distance returns |a − b|, the distance between the endpoints of a and b (two points in space). The 2D distance between points is num/geometry's PointDistance.
Inputs and units: `(a, b Vec3)`; a and b in one unit; the result in that unit.
Output: `(float64, error)`.
Errors: ErrNotFinite for a NaN or ±Inf component; ErrOverflow when the distance exceeds the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.2 "Lengths and Dot Products" (length of v − w) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, 2, 3) and (4, 6, 15) → 13 (hand-checked: 3² + 4² + 12² = 169)
Variants: Vec2 and n-D (Vec) forms where meaningful.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec3_norm.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec3AngleBetween      (ID reserved: NUM-LINALG-020)
Formula: Vec3AngleBetween returns the unoriented angle between a and b, in [0, π] radians, by Kahan's formula θ = 2·atan2(|â − b̂|, |â + b̂|) with â = a/|a|, b̂ = b/|b|. Unlike arccos of the normalised dot product it keeps full accuracy for nearly parallel or antiparallel vectors, and unlike arcsin of the cross product it keeps it near π/2.
Inputs and units: `(a, b Vec3)`; a and b in any units; the result in radians.
Output: `(float64, error)`.
Errors: ErrZeroVector when a or b is zero; ErrNotFinite for a NaN or ±Inf component.
Primary source: Kahan, "How Futile are Mindless Assessments of Roundoff in Floating-Point Computation?", University of California, Berkeley, 11 January 2006, §12 "Mangled Angles", pp. 46–47, https://people.eecs.berkeley.edu/~wkahan/Mindless.pdf (read 2026-11-12) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, 0, 0) and (0, 1, 1) → π/2; (1, 0, 0) and (1, 1, 0) → π/4 (hand-checked)
Variants: Vec2 and n-D (Vec) forms where meaningful.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec3_angle.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec3DirectionAngles      (ID reserved: NUM-LINALG-021)
Formula: Vec3DirectionAngles returns the angles (α, β, γ) between v and the +x, +y and +z axes, each in [0, π] radians; their cosines are the direction cosines v/|v|. Each angle is computed as atan2(|v × eᵢ|, v·eᵢ), e.g. α = atan2(√(y² + z²), x), which is accurate at every angle (the same value Kahan's formula gives for v and the axis).
Inputs and units: `(v Vec3)`; v in any unit; the result in radians.
Output: `(Vec3, error)`.
Errors: ErrZeroVector when v = 0; ErrNotFinite for a NaN or ±Inf component.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.2 "Lengths and Dot Products" (cosine formula with w a coordinate axis) — and — Kahan, "How Futile are Mindless Assessments of Roundoff in Floating-Point Computation?", University of California, Berkeley, 11 January 2006, §12 "Mangled Angles", pp. 46–47, https://people.eecs.berkeley.edu/~wkahan/Mindless.pdf (read 2026-11-12) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, 1, 0) → (π/4, π/4, π/2) (hand-checked)
Variants: Vec2 and n-D (Vec) forms where meaningful.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec3_angle.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec3ComponentAlong      (ID reserved: NUM-LINALG-022)
Formula: Vec3ComponentAlong returns the scalar component of v along axis, v·â with â = axis/|axis|: the signed length of v's projection onto the axis line.
Inputs and units: `(v, axis Vec3)`; v in any unit, axis in any unit; the result in v's unit.
Output: `(float64, error)`.
Errors: ErrZeroVector when axis = 0; ErrNotFinite for a NaN or ±Inf component; ErrOverflow when the result is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.2 "Lengths and Dot Products" (dot product with a unit vector) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: v = (3, 4, 5), axis = (0, 2, 0) → 4 (hand-checked)
Variants: Vec2 and n-D (Vec) forms where meaningful.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec3_projection.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec3Project      (ID reserved: NUM-LINALG-023)
Formula: Vec3Project returns the projection of v onto the line of onto: (v·n)n with n = onto/|onto|, the matrix nnᵀ applied to v.
Inputs and units: `(v, onto Vec3)`; v in any unit, onto in any unit; the result in v's unit.
Output: `(Vec3, error)`.
Errors: ErrZeroVector when onto = 0; ErrNotFinite for a NaN or ±Inf component; ErrOverflow when a component is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §4.2 "Projections" (projection onto a line: p = a(aᵀb)/(aᵀa)) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: v = (3, 4, 5), onto = (0, 2, 0) → (0, 4, 0) (hand-checked)
Variants: Vec2 and n-D (Vec) forms where meaningful.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec3_projection.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec3Reject      (ID reserved: NUM-LINALG-024)
Formula: Vec3Reject returns the rejection of v from the line of from: v − (v·n)n with n = from/|from|, the component of v perpendicular to from (the projection I − nnᵀ onto the plane with normal n).
Inputs and units: `(v, from Vec3)`; v in any unit, from in any unit; the result in v's unit.
Output: `(Vec3, error)`.
Errors: as Vec3Project, with ErrZeroVector when from = 0.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §10.6 "Computer Graphics", p. 497 (projection onto the plane nᵀv = 0 is I − nnᵀ) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: v = (3, 4, 5), from = (0, 2, 0) → (3, 0, 5) (hand-checked)
Variants: Vec2 and n-D (Vec) forms where meaningful.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec3_projection.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec3Reflect      (ID reserved: NUM-LINALG-025)
Formula: Vec3Reflect returns the mirror image of v in the plane through the origin with the given normal: v − 2(v·n)n with n = normal/|normal|, the matrix I − 2nnᵀ.
Inputs and units: `(v, normal Vec3)`; v in any unit, normal in any unit; the result in v's unit.
Output: `(Vec3, error)`.
Errors: as Vec3Project, with ErrZeroVector when normal = 0.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §10.6 "Computer Graphics", p. 498 (mirror matrix I − 2nnᵀ) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: v = (1, −1, 2), normal = (0, 3, 0) → (1, 1, 2) (hand-checked)
Variants: Vec2 and n-D (Vec) forms where meaningful.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec3_projection.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec3ParallelogramArea      (ID reserved: NUM-LINALG-026)
Formula: Vec3ParallelogramArea returns |a × b|, the area of the parallelogram spanned by a and b. In the plane, the area is |Vec2Cross(a, b)|.
Inputs and units: `(a, b Vec3)`; a and b in one length unit; the result in that unit squared.
Output: `(float64, error)`.
Errors: ErrNotFinite for a NaN or ±Inf component; ErrOverflow when the area exceeds the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.3 "Cramer's Rule, Inverses, and Volumes" (|u × v| is the area of the parallelogram) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (3, 0, 0) and (1, 2, 0) → 6 (hand-checked)
Variants: Vec2 and n-D (Vec) forms where meaningful.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec3_area.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec3TriangleArea      (ID reserved: NUM-LINALG-027)
Formula: Vec3TriangleArea returns ½|a × b|, the area of the triangle with two sides a and b from a common vertex.
Inputs and units: `(a, b Vec3)`; a and b in one length unit; the result in that unit squared.
Output: `(float64, error)`.
Errors: as Vec3ParallelogramArea.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.3 "Cramer's Rule, Inverses, and Volumes" (a triangle is half of the parallelogram) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (3, 0, 0) and (1, 2, 0) → 3 (hand-checked)
Variants: Vec2 and n-D (Vec) forms where meaningful.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec3_area.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec3ScalarTripleProduct      (ID reserved: NUM-LINALG-028)
Formula: Vec3ScalarTripleProduct returns a·(b × c), the determinant of the matrix with rows a, b, c: the signed volume of the parallelepiped they span (positive when a, b, c are right-handed). It is the "mixed product".
Inputs and units: `(a, b, c Vec3)`; a, b, c in units A, B, C; the result in A·B·C.
Output: `(float64, error)`.
Errors: ErrNotFinite for a NaN or ±Inf component; ErrOverflow when the result is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.3 "Cramer's Rule, Inverses, and Volumes" (triple product = determinant = volume) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: x̂, ŷ, ẑ → 1; (1, 2, 3), (4, 5, 6), (7, 8, 10) → −3 (hand-checked)
Variants: Vec2 and n-D (Vec) forms where meaningful.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec3_product.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec3VectorTripleProduct      (ID reserved: NUM-LINALG-029)
Formula: Vec3VectorTripleProduct returns a × (b × c). It lies in the plane of b and c and equals b(a·c) − c(a·b); it is computed here as two cross products, each component within 2 ulps of the rounded inner product.
Inputs and units: `(a, b, c Vec3)`; a, b, c in units A, B, C; the result in A·B·C.
Output: `(Vec3, error)`.
Errors: as Vec3ScalarTripleProduct.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.3 "Cramer's Rule, Inverses, and Volumes" (cross product, applied twice) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: a = x̂, b = x̂, c = ŷ → x̂ × ẑ = (0, −1, 0) (hand-checked)
Variants: Vec2 and n-D (Vec) forms where meaningful.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec3_product.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec3PlaneNormal      (ID reserved: NUM-LINALG-030)
Formula: Vec3PlaneNormal returns the unit normal of the plane through p0, p1 and p2: (p1 − p0) × (p2 − p0) normalised. Its direction follows the right-hand rule: it points toward an observer who sees p0 → p1 → p2 counterclockwise.
Inputs and units: `(p0, p1, p2 Vec3)`; points in one length unit; the result is dimensionless.
Output: `(Vec3, error)`.
Errors: ErrZeroVector when the three points are collinear or coincide (no plane); ErrNotFinite for a NaN or ±Inf component; ErrOverflow when an edge vector is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.3 "Cramer's Rule, Inverses, and Volumes" (u × v is perpendicular to u and v) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (0, 0, 0), (1, 0, 0), (0, 1, 0) → (0, 0, 1) (hand-checked)
Variants: Vec2 and n-D (Vec) forms where meaningful.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec3_area.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec3Rotate      (ID reserved: NUM-LINALG-031)
Formula: Vec3Rotate returns v turned by angleRad about axis through the origin, counterclockwise when the axis points at the viewer (right-hand rule). With k = axis/|axis| this is Rodrigues' formula v cos θ + (k × v) sin θ + k (k·v)(1 − cos θ), the active form of the axis-angle rotation matrix.
Inputs and units: `(v, axis Vec3, angleRad float64)`; v in any unit, axis in any unit, angleRad in radians; the result in v's unit.
Output: `(Vec3, error)`.
Errors: ErrZeroVector when axis = 0; ErrNotFinite for a NaN or ±Inf input; ErrOverflow when a component of the result is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §10.6 "Computer Graphics", p. 497, eq. (1) (rotation about the unit vector a; Strang multiplies row vectors, so the column-vector matrix used here is its transpose) — and — Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §7.7, eqs. (243)–(244) (transposed: Diebel's matrix is the frame rotation) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: v = (1, 0, 0), axis = (0, 0, 1), angle π/2 → (0, 1, 0) (hand-checked)
Variants: Vec2 and n-D (Vec) forms where meaningful.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec3_rotate.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec3Lerp      (ID reserved: NUM-LINALG-032)
Formula: Vec3Lerp returns the linear interpolation (1−t)·a + t·b: a at t = 0 and b at t = 1, exactly. t outside [0, 1] extrapolates along the line through a and b.
Inputs and units: `(a, b Vec3, t float64)`; a and b in one unit, t dimensionless; the result in that unit.
Output: `(Vec3, error)`.
Errors: ErrNotFinite for a NaN or ±Inf input; ErrOverflow when a component of the result is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations" (the combination (1−t)a + tb) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: a = (0, 0, 0), b = (2, 4, 6), t = 0.25 → (0.5, 1, 1.5) (hand-checked)
Variants: Vec2 and n-D (Vec) forms where meaningful.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec3.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec3Min      (ID reserved: NUM-LINALG-033)
Formula: Vec3Min returns the componentwise minimum of a and b.
Inputs and units: `(a, b Vec3)`; a and b in one unit; the result in that unit.
Output: `(Vec3, error)`.
Errors: ErrNotFinite for a NaN or ±Inf component.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations" (a vector as its list of components; the minimum is taken component by component) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, 5, −2) and (3, 0, −4) → (1, 0, −4) (hand-checked)
Variants: Vec2 and n-D (Vec) forms where meaningful.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec3.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### Vec3Max      (ID reserved: NUM-LINALG-034)
Formula: Vec3Max returns the componentwise maximum of a and b.
Inputs and units: `(a, b Vec3)`; a and b in one unit; the result in that unit.
Output: `(Vec3, error)`.
Errors: ErrNotFinite for a NaN or ±Inf component.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations" (a vector as its list of components; the maximum is taken component by component) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, 5, −2) and (3, 0, −4) → (3, 5, −2) (hand-checked)
Variants: Vec2 and n-D (Vec) forms where meaningful.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/vec3.go`.
Notes: value types never allocate; every function returns an error rather than a silent NaN.

### VecAdd      (ID reserved: NUM-LINALG-035)
Formula: VecAdd returns a + b, element by element, written into dst (reused when its capacity suffices, else one allocation; dst may be a or b).
Inputs and units: `(dst, a, b []float64)`; a and b in one unit; the result in that unit.
Output: `([]float64, error)`.
Errors: ErrDimension when len(a) ≠ len(b); ErrNotFinite for a NaN or ±Inf element; ErrOverflow when an element of the sum is outside the float64 range. On error the result is nil and dst's contents are unspecified.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations" [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, 2, 3, 4) + (4, 3, 2, 1) → (5, 5, 5, 5) (hand-checked)
Variants: fixed-size Vec2/Vec3 forms.
Best approach: as in Formula. Time O(n)   n = len(a). Space O(1) extra; O(n) for the result, one allocation only if dst lacks capacity. Code: `num/linalg/vec.go`.
Notes: dst may be reused; one allocation only when dst lacks capacity; inputs are never modified.

### VecSub      (ID reserved: NUM-LINALG-036)
Formula: VecSub returns a − b, element by element, written into dst (as VecAdd).
Inputs and units: `(dst, a, b []float64)`; a and b in one unit; the result in that unit.
Output: `([]float64, error)`.
Errors: as VecAdd.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations" [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, 2, 3, 4) − (4, 3, 2, 1) → (−3, −1, 1, 3) (hand-checked)
Variants: fixed-size Vec2/Vec3 forms.
Best approach: as in Formula. Time O(n)   n = len(a). Space O(1) extra; O(n) for the result, one allocation only if dst lacks capacity. Code: `num/linalg/vec.go`.
Notes: dst may be reused; one allocation only when dst lacks capacity; inputs are never modified.

### VecScale      (ID reserved: NUM-LINALG-037)
Formula: VecScale returns k·v written into dst (reused when its capacity suffices; dst may be v).
Inputs and units: `(dst, v []float64, k float64)`; v in any unit, k dimensionless; the result in v's unit.
Output: `([]float64, error)`.
Errors: ErrNotFinite for a NaN or ±Inf input; ErrOverflow when an element of the product is outside the float64 range. On error the result is nil.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations" [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: 2·(1, −2, 3, 0) → (2, −4, 6, 0) (hand-checked)
Variants: fixed-size Vec2/Vec3 forms.
Best approach: as in Formula. Time O(n)   n = len(v). Space O(1) extra; O(n) for the result, one allocation only if dst lacks capacity. Code: `num/linalg/vec.go`.
Notes: dst may be reused; one allocation only when dst lacks capacity; inputs are never modified.

### VecDiv      (ID reserved: NUM-LINALG-038)
Formula: VecDiv returns v/k written into dst (reused when its capacity suffices; dst may be v).
Inputs and units: `(dst, v []float64, k float64)`; v in any unit, k dimensionless; the result in v's unit.
Output: `([]float64, error)`.
Errors: ErrZeroDivisor when k = 0; otherwise as VecScale.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations" [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (2, −4, 6, 0)/2 → (1, −2, 3, 0) (hand-checked)
Variants: fixed-size Vec2/Vec3 forms.
Best approach: as in Formula. Time O(n)   n = len(v). Space O(1) extra; O(n) for the result, one allocation only if dst lacks capacity. Code: `num/linalg/vec.go`.
Notes: dst may be reused; one allocation only when dst lacks capacity; inputs are never modified.

### VecDot      (ID reserved: NUM-LINALG-039)
Formula: VecDot returns the dot product Σ aᵢbᵢ, accumulated with fused multiply-adds. Its error grows with n and with cancellation; VecDotCompensated is as accurate as twice the working precision.
Inputs and units: `(a, b []float64)`; a in unit A, b in unit B; the result in A·B.
Output: `(float64, error)`.
Errors: ErrDimension when len(a) ≠ len(b); ErrNotFinite for a NaN or ±Inf element; ErrOverflow when the sum is outside the float64 range. The dot product of two empty vectors is 0.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.2 "Lengths and Dot Products" [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, 2, 3, 4)·(4, 3, 2, 1) → 20 (hand-checked)
Variants: fixed-size Vec2/Vec3 forms.
Best approach: as in Formula. Time O(n)   n = len(a). Space O(1), no allocation. Code: `num/linalg/vec_dot.go`.
Notes: dst may be reused; one allocation only when dst lacks capacity; inputs are never modified.

### VecDotCompensated      (ID reserved: NUM-LINALG-040)
Formula: VecDotCompensated returns the dot product Σ aᵢbᵢ as accurately as if it were computed in twice the working precision and then rounded: |res − aᵀb| ≤ u|aᵀb| + γ²ₙ|a|ᵀ|b| (u = 2⁻⁵³, γₙ = nu/(1 − nu)), by Ogita, Rump and Oishi's Dot2 with error-free products from fused multiply-adds. Use it for long vectors and for dot products with heavy cancellation.
Inputs and units: `(a, b []float64)`; a in unit A, b in unit B; the result in A·B.
Output: `(float64, error)`.
Errors: as VecDot.
Primary source: Ogita, Rump, Oishi, "Accurate Sum and Dot Product", SIAM Journal on Scientific Computing, 2005, 26(6):1955–1988, doi:10.1137/030601818, Algorithm 5.3 (Dot2), Algorithms 3.1 (TwoSum) and 3.5 (TwoProductFMA), Proposition 5.5 (error bound) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1e16, 1, −1e16)·(1, 1, 1) → 1 (the plain sum gives 0; hand-checked)
Variants: fixed-size Vec2/Vec3 forms.
Best approach: as in Formula. Time O(n)   n = len(a); about 6n floating-point operations. Space O(1), no allocation. Code: `num/linalg/vec_dot.go`.
Notes: dst may be reused; one allocation only when dst lacks capacity; inputs are never modified.

### VecNorm      (ID reserved: NUM-LINALG-041)
Formula: VecNorm returns the Euclidean length ‖v‖₂ = √(Σ vᵢ²). The vector is scaled by an exact power of two so that no square overflows or underflows, and the squares are summed as if in twice the working precision (Dot2), so the result is within about one ulp for any length of vector.
Inputs and units: `(v []float64)`; v in any unit; the result in that unit.
Output: `(float64, error)`.
Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when the length exceeds the float64 range. The norm of an empty vector is 0.
Primary source: Blue, "A Portable Fortran Program to Find the Euclidean Norm of a Vector", ACM Transactions on Mathematical Software, 1978, 4(1):15–23, doi:10.1145/355769.355771 — and — Ogita, Rump, Oishi, "Accurate Sum and Dot Product", SIAM Journal on Scientific Computing, 2005, 26(6):1955–1988, doi:10.1137/030601818, Algorithm 5.3 (Dot2) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, 2, 2, 4) → 5 (hand-checked)
Variants: fixed-size Vec2/Vec3 forms.
Best approach: as in Formula. Time O(n)   n = len(v); two passes. Space O(1), no allocation. Code: `num/linalg/vec_norm.go`.
Notes: dst may be reused; one allocation only when dst lacks capacity; inputs are never modified.

### VecNorm1      (ID reserved: NUM-LINALG-042)
Formula: VecNorm1 returns the 1-norm ‖v‖₁ = Σ|vᵢ|, summed by Ogita, Rump and Oishi's cascaded summation (Sum2), as accurate as in twice the working precision.
Inputs and units: `(v []float64)`; v in any unit; the result in that unit.
Output: `(float64, error)`.
Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when the sum exceeds the float64 range.
Primary source: Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002, §6.1 "Vector Norms" — and — Ogita, Rump, Oishi, "Accurate Sum and Dot Product", SIAM Journal on Scientific Computing, 2005, 26(6):1955–1988, doi:10.1137/030601818, Algorithm 4.4 (Sum2) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, −2, 3, −4) → 10 (hand-checked)
Variants: fixed-size Vec2/Vec3 forms.
Best approach: as in Formula. Time O(n)   n = len(v). Space O(1), no allocation. Code: `num/linalg/vec_norm.go`.
Notes: dst may be reused; one allocation only when dst lacks capacity; inputs are never modified.

### VecNormInf      (ID reserved: NUM-LINALG-043)
Formula: VecNormInf returns the infinity norm ‖v‖∞ = max|vᵢ| (0 for an empty vector).
Inputs and units: `(v []float64)`; v in any unit; the result in that unit.
Output: `(float64, error)`.
Errors: ErrNotFinite for a NaN or ±Inf element.
Primary source: Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002, §6.1 "Vector Norms" [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, −7, 3) → 7 (hand-checked)
Variants: fixed-size Vec2/Vec3 forms.
Best approach: as in Formula. Time O(n)   n = len(v). Space O(1), no allocation. Code: `num/linalg/vec_norm.go`.
Notes: dst may be reused; one allocation only when dst lacks capacity; inputs are never modified.

### VecNormP      (ID reserved: NUM-LINALG-044)
Formula: VecNormP returns the p-norm ‖v‖ₚ = (Σ|vᵢ|ᵖ)^(1/p) for p ≥ 1, computed as m·(Σ(|vᵢ|/m)ᵖ)^(1/p) with m = max|vᵢ| so that no power overflows, and the sum taken by Sum2. p = 1 and p = 2 agree with VecNorm1 and VecNorm to within rounding; use VecNormInf for p = ∞.
Inputs and units: `(v []float64, p float64)`; v in any unit, p dimensionless; the result in v's unit.
Output: `(float64, error)`.
Errors: ErrOutOfRange when p < 1; ErrNotFinite for a NaN or ±Inf element or p; ErrOverflow when the norm exceeds the float64 range.
Primary source: Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002, §6.1 "Vector Norms" (Hölder p-norm) — and — Ogita, Rump, Oishi, "Accurate Sum and Dot Product", SIAM Journal on Scientific Computing, 2005, 26(6):1955–1988, doi:10.1137/030601818, Algorithm 4.4 (Sum2) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (3, 4), p = 2 → 5; (1, 1, 1, 1), p = 3 → 4^(1/3) ≈ 1.5874010519681994 (hand-checked)
Variants: fixed-size Vec2/Vec3 forms.
Best approach: as in Formula. Time O(n)   n = len(v); two passes and one math.Pow per element. Space O(1), no allocation. Code: `num/linalg/vec_norm.go`.
Notes: dst may be reused; one allocation only when dst lacks capacity; inputs are never modified.

### VecNormalize      (ID reserved: NUM-LINALG-045)
Formula: VecNormalize returns v/‖v‖₂ written into dst (reused when its capacity suffices; dst may be v). The division uses the scaled vector, so it is exact in range even when ‖v‖₂ itself would overflow.
Inputs and units: `(dst, v []float64)`; v in any unit; the result is dimensionless.
Output: `([]float64, error)`.
Errors: ErrZeroVector when v is all zeros (or empty); ErrNotFinite for a NaN or ±Inf element.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.2 "Lengths and Dot Products" — and — Blue, "A Portable Fortran Program to Find the Euclidean Norm of a Vector", ACM Transactions on Mathematical Software, 1978, 4(1):15–23, doi:10.1145/355769.355771 [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, 2, 2, 4) → (0.2, 0.4, 0.4, 0.8) (hand-checked)
Variants: fixed-size Vec2/Vec3 forms.
Best approach: as in Formula. Time O(n)   n = len(v); three passes. Space O(1) extra; O(n) for the result, one allocation only if dst lacks capacity. Code: `num/linalg/vec_norm.go`.
Notes: dst may be reused; one allocation only when dst lacks capacity; inputs are never modified.

### VecDistance      (ID reserved: NUM-LINALG-046)
Formula: VecDistance returns ‖a − b‖₂, the Euclidean distance between the points a and b, without forming a − b in memory, scaled against overflow and summed by Dot2.
Inputs and units: `(a, b []float64)`; a and b in one unit; the result in that unit.
Output: `(float64, error)`.
Errors: ErrDimension when len(a) ≠ len(b); ErrNotFinite for a NaN or ±Inf element; ErrOverflow when an element difference or the distance exceeds the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.2 "Lengths and Dot Products" (length of v − w) — and — Blue, "A Portable Fortran Program to Find the Euclidean Norm of a Vector", ACM Transactions on Mathematical Software, 1978, 4(1):15–23, doi:10.1145/355769.355771 [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, 1, 1, 1) and (2, 3, 3, 5) → 5 (hand-checked)
Variants: fixed-size Vec2/Vec3 forms.
Best approach: as in Formula. Time O(n)   n = len(a); two passes. Space O(1), no allocation. Code: `num/linalg/vec_norm.go`.
Notes: dst may be reused; one allocation only when dst lacks capacity; inputs are never modified.

### VecAngleBetween      (ID reserved: NUM-LINALG-047)
Formula: VecAngleBetween returns the unoriented angle between a and b, in [0, π] radians, in any dimension, by Kahan's formula θ = 2·atan2(‖â − b̂‖, ‖â + b̂‖) with â = a/‖a‖, b̂ = b/‖b‖. Its absolute error stays near one ulp of π for nearly parallel, perpendicular and antiparallel vectors alike.
Inputs and units: `(a, b []float64)`; a and b in any units; the result in radians.
Output: `(float64, error)`.
Errors: ErrDimension when len(a) ≠ len(b); ErrZeroVector when a or b is all zeros (or empty); ErrNotFinite for a NaN or ±Inf element.
Primary source: Kahan, "How Futile are Mindless Assessments of Roundoff in Floating-Point Computation?", University of California, Berkeley, 11 January 2006, §12 "Mangled Angles", pp. 46–47, https://people.eecs.berkeley.edu/~wkahan/Mindless.pdf (read 2026-11-12) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, 0, 0, 0) and (1, 1, 1, 1) → π/3 (hand-checked: cos θ = 1/2)
Variants: fixed-size Vec2/Vec3 forms.
Best approach: as in Formula. Time O(n)   n = len(a); five passes. Space O(1), no allocation. Code: `num/linalg/vec_angle.go`.
Notes: dst may be reused; one allocation only when dst lacks capacity; inputs are never modified.

### VecProject      (ID reserved: NUM-LINALG-048)
Formula: VecProject returns the projection of v onto the line of onto, (v·n)n with n = onto/‖onto‖₂, written into dst (reused when its capacity suffices; dst may be v or onto). The dot product is taken by Dot2.
Inputs and units: `(dst, v, onto []float64)`; v in any unit, onto in any unit; the result in v's unit.
Output: `([]float64, error)`.
Errors: ErrDimension when len(v) ≠ len(onto); ErrZeroVector when onto is all zeros (or empty); ErrNotFinite for a NaN or ±Inf element; ErrOverflow when an element of the result is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §4.2 "Projections" (projection onto a line: p = a(aᵀb)/(aᵀa)) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: v = (1, 2, 3, 4), onto = (0, 0, 0, 2) → (0, 0, 0, 4) (hand-checked)
Variants: fixed-size Vec2/Vec3 forms.
Best approach: as in Formula. Time O(n)   n = len(v); four passes. Space O(1) extra; O(n) for the result, one allocation only if dst lacks capacity. Code: `num/linalg/vec_angle.go`.
Notes: dst may be reused; one allocation only when dst lacks capacity; inputs are never modified.

### VecLerp      (ID reserved: NUM-LINALG-049)
Formula: VecLerp returns (1−t)·a + t·b written into dst: a at t = 0 and b at t = 1, exactly; t outside [0, 1] extrapolates. dst may be a or b.
Inputs and units: `(dst, a, b []float64, t float64)`; a and b in one unit, t dimensionless; the result in that unit.
Output: `([]float64, error)`.
Errors: ErrDimension when len(a) ≠ len(b); ErrNotFinite for a NaN or ±Inf input; ErrOverflow when an element is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations" (the combination (1−t)a + tb) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: a = (0, 10), b = (4, 20), t = 0.5 → (2, 15) (hand-checked)
Variants: fixed-size Vec2/Vec3 forms.
Best approach: as in Formula. Time O(n)   n = len(a). Space O(1) extra; O(n) for the result, one allocation only if dst lacks capacity. Code: `num/linalg/vec.go`.
Notes: dst may be reused; one allocation only when dst lacks capacity; inputs are never modified.

### VecMin      (ID reserved: NUM-LINALG-050)
Formula: VecMin returns the element-by-element minimum of a and b written into dst (dst may be a or b).
Inputs and units: `(dst, a, b []float64)`; a and b in one unit; the result in that unit.
Output: `([]float64, error)`.
Errors: ErrDimension when len(a) ≠ len(b); ErrNotFinite for a NaN or ±Inf element.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations" (a vector as its list of components; the minimum is taken component by component) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, 5, −2) and (3, 0, −4) → (1, 0, −4) (hand-checked)
Variants: fixed-size Vec2/Vec3 forms.
Best approach: as in Formula. Time O(n)   n = len(a). Space O(1) extra; O(n) for the result, one allocation only if dst lacks capacity. Code: `num/linalg/vec.go`.
Notes: dst may be reused; one allocation only when dst lacks capacity; inputs are never modified.

### VecMax      (ID reserved: NUM-LINALG-051)
Formula: VecMax returns the element-by-element maximum of a and b written into dst (dst may be a or b).
Inputs and units: `(dst, a, b []float64)`; a and b in one unit; the result in that unit.
Output: `([]float64, error)`.
Errors: as VecMin.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.1 "Vectors and Linear Combinations" (a vector as its list of components; the maximum is taken component by component) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, 5, −2) and (3, 0, −4) → (3, 5, −2) (hand-checked)
Variants: fixed-size Vec2/Vec3 forms.
Best approach: as in Formula. Time O(n)   n = len(a). Space O(1) extra; O(n) for the result, one allocation only if dst lacks capacity. Code: `num/linalg/vec.go`.
Notes: dst may be reused; one allocation only when dst lacks capacity; inputs are never modified.

### CartesianToCylindrical      (ID reserved: NUM-LINALG-052)
Formula: CartesianToCylindrical returns the cylindrical coordinates (R, φ, z) of the point v = (x, y, z): R = √(x² + y²) (without overflow in the squares), φ = atan2(y, x) in [0, 2π), z unchanged. On the z axis (x = y = 0) the azimuth is undefined and 0 is returned.
Inputs and units: `(v Vec3)`; x, y, z in one length unit; R and Z in that unit, PhiRad in radians.
Output: `(Cylindrical, error)`.
Errors: ErrNotFinite for a NaN or ±Inf component; ErrOverflow when R exceeds the float64 range.
Primary source: NIST Digital Library of Mathematical Functions, §1.5(ii) "Coordinate Systems", eq. 1.5.14 (inverted), https://dlmf.nist.gov/1.5 (read 2026-11-12) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, 1, 5) → (√2, π/4, 5) (hand-checked)
Variants: the mathematicians' convention swaps θ and φ (DLMF §1.5(ii) notes it); geographic latitude is π/2 − θ.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/cylindrical.go`.
Notes: angles in radians; φ is returned in [0, 2π).

### CylindricalToCartesian      (ID reserved: NUM-LINALG-053)
Formula: CylindricalToCartesian returns the point (R cos φ, R sin φ, z). Any finite φ is accepted (it need not be reduced to [0, 2π)).
Inputs and units: `(c Cylindrical)`; R and Z in one length unit, PhiRad in radians; the result in that length unit.
Output: `(Vec3, error)`.
Errors: ErrOutOfRange when R < 0; ErrNotFinite for a NaN or ±Inf coordinate.
Primary source: NIST Digital Library of Mathematical Functions, §1.5(ii) "Coordinate Systems", eq. 1.5.14, https://dlmf.nist.gov/1.5 (read 2026-11-12) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (2, π/2, 3) → (0, 2, 3) (hand-checked; the computed x is 2·cos(π/2) ≈ 1.2e−16)
Variants: as above.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/cylindrical.go`.
Notes: angles in radians.

### CartesianToSpherical      (ID reserved: NUM-LINALG-054)
Formula: CartesianToSpherical returns the spherical coordinates (R, θ, φ) of v = (x, y, z): R = |v| (without overflow in the squares), θ = atan2(√(x² + y²), z) in [0, π] (accurate at the poles, where arccos(z/R) loses half its digits), φ = atan2(y, x) in [0, 2π). At the origin both angles are undefined and 0 is returned; on the z axis φ is undefined and 0 is returned.
Inputs and units: `(v Vec3)`; x, y, z in one length unit; R in that unit, ThetaRad and PhiRad in radians.
Output: `(Spherical, error)`.
Errors: ErrNotFinite for a NaN or ±Inf component; ErrOverflow when R exceeds the float64 range.
Primary source: NIST Digital Library of Mathematical Functions, §1.5(ii) "Coordinate Systems", eq. 1.5.16 (inverted), https://dlmf.nist.gov/1.5 (read 2026-11-12) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (0, 1, 1) → (√2, π/4, π/2) (hand-checked)
Variants: the mathematicians' convention swaps θ and φ (DLMF §1.5(ii) notes it); geographic latitude is π/2 − θ.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/spherical.go`.
Notes: angles in radians; φ is returned in [0, 2π).

### SphericalToCartesian      (ID reserved: NUM-LINALG-055)
Formula: SphericalToCartesian returns the point (R sin θ cos φ, R sin θ sin φ, R cos θ). Any finite angles are accepted.
Inputs and units: `(s Spherical)`; R in a length unit, ThetaRad and PhiRad in radians; the result in that length unit.
Output: `(Vec3, error)`.
Errors: ErrOutOfRange when R < 0; ErrNotFinite for a NaN or ±Inf coordinate.
Primary source: NIST Digital Library of Mathematical Functions, §1.5(ii) "Coordinate Systems", eq. 1.5.16, https://dlmf.nist.gov/1.5 (read 2026-11-12) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (2, π/2, 0) → (2, 0, 0) (hand-checked; the computed z is 2·cos(π/2) ≈ 1.2e−16)
Variants: as above.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/spherical.go`.
Notes: angles in radians.

### Mat2Identity      (ID reserved: NUM-LINALG-056)
Formula: Mat2Identity returns the 2×2 identity matrix.
Inputs and units: `()`; dimensionless.
Output: `Mat2`.
Errors: none (no input can fail)
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.5 "Inverse Matrices" (the identity I) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: → [1 0; 0 1]
Variants: Mat3, Mat4 and n×n forms.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/mat2.go`.
Notes: column vectors; active rotations.

### Mat2Mul      (ID reserved: NUM-LINALG-057)
Formula: Mat2Mul returns the matrix product a·b. Each element is a two-term dot product computed within 2 ulps (Kahan's algorithm).
Inputs and units: `(a, b Mat2)`; a in unit A, b in unit B; the result in A·B.
Output: `(Mat2, error)`.
Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when an element of the product is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.4 "Rules for Matrix Operations" (row times column) — and — Jeannerod, Louvet, Muller, "Further analysis of Kahan's algorithm for the accurate computation of 2×2 determinants", Mathematics of Computation, 2013, 82(284):2245–2264, doi:10.1090/S0025-5718-2013-02679-8 [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: [1 2; 3 4]·[5 6; 7 8] → [19 22; 43 50] (hand-checked)
Variants: Mat3, Mat4 and n×n forms.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/mat2.go`.
Notes: column vectors; active rotations.

### Mat2MulVec      (ID reserved: NUM-LINALG-058)
Formula: Mat2MulVec returns the matrix-vector product m·v (v a column), each element within 2 ulps. This is the "2×2 matrix times 2×1 matrix" product.
Inputs and units: `(m Mat2, v Vec2)`; m in unit A, v in unit B; the result in A·B.
Output: `(Vec2, error)`.
Errors: ErrNotFinite for a NaN or ±Inf input; ErrOverflow when an element is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.4 "Rules for Matrix Operations" (matrix times vector) — and — Jeannerod, Louvet, Muller, "Further analysis of Kahan's algorithm for the accurate computation of 2×2 determinants", Mathematics of Computation, 2013, 82(284):2245–2264, doi:10.1090/S0025-5718-2013-02679-8 [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: [1 2; 3 4]·(5, 6) → (17, 39) (hand-checked)
Variants: Mat3, Mat4 and n×n forms.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/mat2.go`.
Notes: column vectors; active rotations.

### Mat2Transpose      (ID reserved: NUM-LINALG-059)
Formula: Mat2Transpose returns mᵀ.
Inputs and units: `(m Mat2)`; the unit of m.
Output: `(Mat2, error)`.
Errors: ErrNotFinite for a NaN or ±Inf element.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.7 "Transposes and Permutations" [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: [1 2; 3 4] → [1 3; 2 4]
Variants: Mat3, Mat4 and n×n forms.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/mat2.go`.
Notes: column vectors; active rotations.

### Mat2Det      (ID reserved: NUM-LINALG-060)
Formula: Mat2Det returns det m = m₀₀m₁₁ − m₀₁m₁₀, within 2 ulps by Kahan's fused multiply-add algorithm, even when the two products nearly cancel.
Inputs and units: `(m Mat2)`; m in unit A; the result in A².
Output: `(float64, error)`.
Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when the determinant is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.1 "The Properties of Determinants" — and — Jeannerod, Louvet, Muller, "Further analysis of Kahan's algorithm for the accurate computation of 2×2 determinants", Mathematics of Computation, 2013, 82(284):2245–2264, doi:10.1090/S0025-5718-2013-02679-8 [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: [1 2; 3 4] → −2 (hand-checked)
Variants: Mat3, Mat4 and n×n forms.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/mat2.go`.
Notes: column vectors; active rotations.

### Mat2Inverse      (ID reserved: NUM-LINALG-061)
Formula: Mat2Inverse returns m⁻¹ = [m₁₁ −m₀₁; −m₁₀ m₀₀]/det m. Each row and then each column is first scaled by an exact power of two (equilibration), so the determinant neither overflows nor underflows for a badly scaled but well-conditioned matrix. An ill-conditioned matrix has an inaccurate inverse; check with MatConditionEstimate1.
Inputs and units: `(m Mat2)`; m in unit A; the result in 1/A.
Output: `(Mat2, error)`.
Errors: ErrSingular when det m = 0; ErrNotFinite for a NaN or ±Inf element; ErrOverflow when an element of the inverse is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.5 "Inverse Matrices" (the 2 by 2 inverse) and §5.3 "Cramer's Rule, Inverses, and Volumes" [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: [4 7; 2 6] → [0.6 −0.7; −0.2 0.4] (hand-checked: det = 10)
Variants: Mat3, Mat4 and n×n forms.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/mat2.go`.
Notes: column vectors; active rotations.

### Mat2Trace      (ID reserved: NUM-LINALG-062)
Formula: Mat2Trace returns m₀₀ + m₁₁, the sum of the eigenvalues.
Inputs and units: `(m Mat2)`; the unit of m.
Output: `(float64, error)`.
Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when the sum is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.1 "Introduction to Eigenvalues" (the trace equals the sum of the eigenvalues) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: [1 2; 3 4] → 5
Variants: Mat3, Mat4 and n×n forms.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/mat2.go`.
Notes: column vectors; active rotations.

### Mat2Rotation      (ID reserved: NUM-LINALG-063)
Formula: Mat2Rotation returns the matrix that turns a column vector counterclockwise by angleRad: [cos −sin; sin cos].
Inputs and units: `(angleRad float64)`; angleRad in radians; the result is dimensionless.
Output: `(Mat2, error)`.
Errors: ErrNotFinite when angleRad is NaN or ±Inf.
Primary source: Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §2.4, eq. (16) (the upper-left 2×2 block of R₃, transposed: Diebel's matrices are frame rotations) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: π/2 → [0 −1; 1 0] (to within cos(π/2) ≈ 6.1e−17)
Variants: Mat3, Mat4 and n×n forms.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/mat2.go`.
Notes: column vectors; active rotations.

### Mat2SymmetricEigenvalues      (ID reserved: NUM-LINALG-064)
Formula: Mat2SymmetricEigenvalues returns the two eigenvalues lo ≤ hi of the symmetric matrix [a b; b c]. They are the roots of λ² − (a + c)λ + (ac − b²) = 0; the one of larger magnitude is (a + c)/2 ± hypot((a − c)/2, b) with the sign of (a + c), and the other is det/that, so neither suffers cancellation (the determinant itself is computed within 2 ulps by Kahan's algorithm).
Inputs and units: `(m Mat2)`; m in any unit; the eigenvalues in that unit.
Output: `(lo, hi float64, err error)`.
Errors: ErrNotSymmetric when m₀₁ ≠ m₁₀ (compared exactly); ErrNotFinite for a NaN or ±Inf element; ErrOverflow when an eigenvalue is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.1 "Introduction to Eigenvalues" (λ₁ + λ₂ = trace, λ₁λ₂ = det) and §6.4 "Symmetric Matrices" (real eigenvalues) — and — Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §8.5 "Jacobi Methods" (the 2-by-2 symmetric Schur decomposition) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: [2 1; 1 2] → (1, 3) (hand-checked)
Variants: Mat3, Mat4 and n×n forms.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/symmetric_eigen.go`.
Notes: column vectors; active rotations.

### Mat3Identity      (ID reserved: NUM-LINALG-065)
Formula: Mat3Identity returns the 3×3 identity matrix.
Inputs and units: `()`; dimensionless.
Output: `Mat3`.
Errors: none (no input can fail)
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.5 "Inverse Matrices" (the identity I) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: → [1 0 0; 0 1 0; 0 0 1]
Variants: Mat2, Mat4 and n×n forms.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/mat3.go`.
Notes: column vectors; active rotations (Diebel's are the transposes).

### Mat3Mul      (ID reserved: NUM-LINALG-066)
Formula: Mat3Mul returns the matrix product a·b, each element a three-term dot product accumulated with fused multiply-adds.
Inputs and units: `(a, b Mat3)`; a in unit A, b in unit B; the result in A·B.
Output: `(Mat3, error)`.
Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when an element of the product is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.4 "Rules for Matrix Operations" (row times column) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: [1 2 3; 4 5 6; 7 8 9]·I → the same matrix; [1 2 0; 0 1 0; 0 0 2]·[1 0 0; 3 1 0; 0 0 1] → [7 2 0; 3 1 0; 0 0 2] (hand-checked)
Variants: Mat2, Mat4 and n×n forms.
Best approach: as in Formula. Time O(1)   27 multiply-adds. Space O(1), no allocation. Code: `num/linalg/mat3.go`.
Notes: column vectors; active rotations (Diebel's are the transposes).

### Mat3MulVec      (ID reserved: NUM-LINALG-067)
Formula: Mat3MulVec returns the matrix-vector product m·v (v a column).
Inputs and units: `(m Mat3, v Vec3)`; m in unit A, v in unit B; the result in A·B.
Output: `(Vec3, error)`.
Errors: ErrNotFinite for a NaN or ±Inf input; ErrOverflow when an element is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.4 "Rules for Matrix Operations" (matrix times vector) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: [1 2 3; 4 5 6; 7 8 9]·(1, 0, −1) → (−2, −2, −2) (hand-checked)
Variants: Mat2, Mat4 and n×n forms.
Best approach: as in Formula. Time O(1)   9 multiply-adds. Space O(1), no allocation. Code: `num/linalg/mat3.go`.
Notes: column vectors; active rotations (Diebel's are the transposes).

### Mat3Transpose      (ID reserved: NUM-LINALG-068)
Formula: Mat3Transpose returns mᵀ.
Inputs and units: `(m Mat3)`; the unit of m.
Output: `(Mat3, error)`.
Errors: ErrNotFinite for a NaN or ±Inf element.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.7 "Transposes and Permutations" [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: [1 2 3; 4 5 6; 7 8 9] → [1 4 7; 2 5 8; 3 6 9]
Variants: Mat2, Mat4 and n×n forms.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/mat3.go`.
Notes: column vectors; active rotations (Diebel's are the transposes).

### Mat3Det      (ID reserved: NUM-LINALG-069)
Formula: Mat3Det returns det m as the triple product r₀·(r₁ × r₂) of its rows: the cross product's components each within 2 ulps (Kahan's algorithm), the final dot product as if in twice the working precision (Dot2), so an exactly singular matrix gives 0. When an intermediate overflows, the rows are rescaled by exact powers of two and the determinant recomputed.
Inputs and units: `(m Mat3)`; m in unit A; the result in A³.
Output: `(float64, error)`.
Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when the determinant is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.2 "Permutations and Cofactors" and §5.3 "Cramer's Rule, Inverses, and Volumes" (the triple product is the determinant) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: [1 2 3; 0 1 4; 5 6 0] → 1 (hand-checked: 1·(−24) − 2·(−20) + 3·(−5))
Variants: Mat2, Mat4 and n×n forms.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/mat3.go`.
Notes: column vectors; active rotations (Diebel's are the transposes).

### Mat3Inverse      (ID reserved: NUM-LINALG-070)
Formula: Mat3Inverse returns m⁻¹ = Cᵀ/det m (the cofactor formula): column j of the inverse is the cross product of the two rows other than j, in cyclic order, divided by det m. Each row and then each column is first scaled by an exact power of two (equilibration), so the cofactors and the determinant neither overflow nor underflow for a badly scaled but well-conditioned matrix. An ill-conditioned matrix has an inaccurate inverse; check with MatConditionEstimate1.
Inputs and units: `(m Mat3)`; m in unit A; the result in 1/A.
Output: `(Mat3, error)`.
Errors: ErrSingular when det m = 0; ErrNotFinite for a NaN or ±Inf element; ErrOverflow when an element of the inverse is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.3 "Cramer's Rule, Inverses, and Volumes" (A⁻¹ = Cᵀ/det A) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: [1 2 3; 0 1 4; 5 6 0] → [−24 18 5; 20 −15 −4; −5 4 1] (hand-checked: det = 1, so the inverse is Cᵀ)
Variants: Mat2, Mat4 and n×n forms.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/mat3.go`.
Notes: column vectors; active rotations (Diebel's are the transposes).

### Mat3Trace      (ID reserved: NUM-LINALG-071)
Formula: Mat3Trace returns m₀₀ + m₁₁ + m₂₂, the sum of the eigenvalues.
Inputs and units: `(m Mat3)`; the unit of m.
Output: `(float64, error)`.
Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when the sum is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.1 "Introduction to Eigenvalues" (the trace equals the sum of the eigenvalues) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: [1 2 3; 4 5 6; 7 8 9] → 15
Variants: Mat2, Mat4 and n×n forms.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/mat3.go`.
Notes: column vectors; active rotations (Diebel's are the transposes).

### Mat3RotationX      (ID reserved: NUM-LINALG-072)
Formula: Mat3RotationX returns the matrix that turns a vector by angleRad about the +x axis: [1 0 0; 0 cos −sin; 0 sin cos].
Inputs and units: `(angleRad float64)`; angleRad in radians; the result is dimensionless.
Output: `(Mat3, error)`.
Errors: ErrNotFinite when angleRad is NaN or ±Inf.
Primary source: Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §2.4, eq. (14) (R₁, transposed: Diebel's matrices are frame rotations) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: π/2 turns ŷ into ẑ
Variants: Mat2, Mat4 and n×n forms.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/mat3_rotation.go`.
Notes: column vectors; active rotations (Diebel's are the transposes).

### Mat3RotationY      (ID reserved: NUM-LINALG-073)
Formula: Mat3RotationY returns the matrix that turns a vector by angleRad about the +y axis: [cos 0 sin; 0 1 0; −sin 0 cos].
Inputs and units: `(angleRad float64)`; angleRad in radians; the result is dimensionless.
Output: `(Mat3, error)`.
Errors: ErrNotFinite when angleRad is NaN or ±Inf.
Primary source: Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §2.4, eq. (15) (R₂, transposed: Diebel's matrices are frame rotations) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: π/2 turns ẑ into x̂
Variants: Mat2, Mat4 and n×n forms.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/mat3_rotation.go`.
Notes: column vectors; active rotations (Diebel's are the transposes).

### Mat3RotationZ      (ID reserved: NUM-LINALG-074)
Formula: Mat3RotationZ returns the matrix that turns a vector by angleRad about the +z axis: [cos −sin 0; sin cos 0; 0 0 1].
Inputs and units: `(angleRad float64)`; angleRad in radians; the result is dimensionless.
Output: `(Mat3, error)`.
Errors: ErrNotFinite when angleRad is NaN or ±Inf.
Primary source: Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §2.4, eq. (16) (R₃, transposed: Diebel's matrices are frame rotations) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: π/2 turns x̂ into ŷ
Variants: Mat2, Mat4 and n×n forms.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/mat3_rotation.go`.
Notes: column vectors; active rotations (Diebel's are the transposes).

### Mat3RotationAxisAngle      (ID reserved: NUM-LINALG-075)
Formula: Mat3RotationAxisAngle returns the matrix that turns a vector by angleRad about axis (right-hand rule): with k = axis/|axis|, R = cos θ·I + (1 − cos θ)·kkᵀ + sin θ·[k]×, where [k]× is the cross-product matrix (Rodrigues' formula in matrix form). 1 − cos θ is computed as 2 sin²(θ/2) for accuracy at small angles.
Inputs and units: `(axis Vec3, angleRad float64)`; axis in any unit, angleRad in radians; the result is dimensionless.
Output: `(Mat3, error)`.
Errors: ErrZeroVector when axis = 0; ErrNotFinite for a NaN or ±Inf input.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §10.6 "Computer Graphics", p. 497, eq. (1) (rotation about the unit vector a; Strang multiplies row vectors, so the matrix here is its transpose) — and — Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §7.7, eqs. (243)–(244) (transposed: Diebel's matrix is the frame rotation) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: axis ẑ, π/2 → [0 −1 0; 1 0 0; 0 0 1] (to within cos(π/2) ≈ 6.1e−17)
Variants: Mat2, Mat4 and n×n forms.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/mat3_rotation.go`.
Notes: column vectors; active rotations (Diebel's are the transposes).

### Mat3SymmetricEigenvalues      (ID reserved: NUM-LINALG-076)
Formula: Mat3SymmetricEigenvalues returns the three eigenvalues of a symmetric 3×3 matrix in ascending order, by Smith's closed form: with m = trace/3, K = A − mI, q = det(K)/2 and p = (Σ Kᵢⱼ²)/6, the angle φ = ⅓·atan2(√(p³ − q²), q) gives the eigenvalues m + 2√p cos φ and m − √p(cos φ ± √3 sin φ). The matrix is first scaled by an exact power of two (eigenvalues scale with it) so no square overflows. For distinct eigenvalues the absolute error is a few ulps of the largest |eigenvalue|. For two equal or nearly equal eigenvalues p³ − q² nearly vanishes and its rounding error passes through a square root, so those two carry an absolute error up to about √ε·max|λ| (≈ 1.5e−8 relative); an eigenvalue much smaller than the largest also loses relative accuracy. An iterative symmetric eigensolver (LATER) does not have these limits.
Inputs and units: `(m Mat3)`; m in any unit; the eigenvalues in that unit.
Output: `([3]float64, error)`.
Errors: ErrNotSymmetric unless m₀₁ = m₁₀, m₀₂ = m₂₀ and m₁₂ = m₂₁ exactly; ErrNotFinite for a NaN or ±Inf element; ErrOverflow when an eigenvalue is outside the float64 range.
Primary source: Smith, "Eigenvalues of a symmetric 3 × 3 matrix", Communications of the ACM, 1961, 4(4):168, doi:10.1145/355578.366316 [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: [2 0 0; 0 3 4; 0 4 9] → (1, 2, 11) (hand-checked: the lower block has trace 12 and det 11)
Variants: Mat2, Mat4 and n×n forms.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/symmetric_eigen.go`.
Notes: column vectors; active rotations (Diebel's are the transposes).

### Mat4Identity      (ID reserved: NUM-LINALG-077)
Formula: Mat4Identity returns the 4×4 identity matrix.
Inputs and units: `()`; dimensionless.
Output: `Mat4`.
Errors: none (no input can fail)
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.5 "Inverse Matrices" (the identity I) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: → the 4×4 identity
Variants: Mat2, Mat3 and n×n forms.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/mat4.go`.
Notes: homogeneous transforms act on column vectors (Strang §10.6 uses rows: transpose).

### Mat4Mul      (ID reserved: NUM-LINALG-078)
Formula: Mat4Mul returns the matrix product a·b, each element a four-term dot product accumulated with fused multiply-adds. Composing transforms: (a·b)·v applies b first, then a.
Inputs and units: `(a, b Mat4)`; a in unit A, b in unit B; the result in A·B.
Output: `(Mat4, error)`.
Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when an element of the product is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.4 "Rules for Matrix Operations" (row times column) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: translation by (1, 2, 3) times translation by (4, 5, 6) → translation by (5, 7, 9) (hand-checked)
Variants: Mat2, Mat3 and n×n forms.
Best approach: as in Formula. Time O(1)   64 multiply-adds. Space O(1), no allocation. Code: `num/linalg/mat4.go`.
Notes: homogeneous transforms act on column vectors (Strang §10.6 uses rows: transpose).

### Mat4MulVec      (ID reserved: NUM-LINALG-079)
Formula: Mat4MulVec returns the matrix-vector product m·v (v a column).
Inputs and units: `(m Mat4, v Vec4)`; m in unit A, v in unit B; the result in A·B.
Output: `(Vec4, error)`.
Errors: ErrNotFinite for a NaN or ±Inf input; ErrOverflow when an element is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.4 "Rules for Matrix Operations" (matrix times vector) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: translation by (1, 2, 3) times (0, 0, 0, 1) → (1, 2, 3, 1) (hand-checked)
Variants: Mat2, Mat3 and n×n forms.
Best approach: as in Formula. Time O(1)   16 multiply-adds. Space O(1), no allocation. Code: `num/linalg/mat4.go`.
Notes: homogeneous transforms act on column vectors (Strang §10.6 uses rows: transpose).

### Mat4Transpose      (ID reserved: NUM-LINALG-080)
Formula: Mat4Transpose returns mᵀ.
Inputs and units: `(m Mat4)`; the unit of m.
Output: `(Mat4, error)`.
Errors: ErrNotFinite for a NaN or ±Inf element.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.7 "Transposes and Permutations" [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: the transpose of a translation matrix moves the offset from the last column to the last row
Variants: Mat2, Mat3 and n×n forms.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/mat4.go`.
Notes: homogeneous transforms act on column vectors (Strang §10.6 uses rows: transpose).

### Mat4Det      (ID reserved: NUM-LINALG-081)
Formula: Mat4Det returns det m by the Laplace expansion along the first two rows: the sum of the six products of complementary 2×2 minors, with signs. Each minor is within 2 ulps (Kahan's algorithm) and the sum is taken by Dot2. When an intermediate overflows, the rows are rescaled by exact powers of two and the determinant recomputed.
Inputs and units: `(m Mat4)`; m in unit A; the result in A⁴.
Output: `(float64, error)`.
Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when the determinant is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.2 "Permutations and Cofactors" [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: diag(1, 2, 3, 4) → 24; a homogeneous translation → 1 (hand-checked)
Variants: Mat2, Mat3 and n×n forms.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/mat4.go`.
Notes: homogeneous transforms act on column vectors (Strang §10.6 uses rows: transpose).

### Mat4Inverse      (ID reserved: NUM-LINALG-082)
Formula: Mat4Inverse returns m⁻¹ = Cᵀ/det m (the cofactor formula), building every cofactor from the twelve 2×2 minors of the Laplace expansion. Each row and then each column is first scaled by an exact power of two (equilibration), so the minors and the determinant neither overflow nor underflow for a badly scaled but well-conditioned matrix. A tiny determinant relative to the elements means an ill-conditioned matrix; check with MatConditionEstimate1, or use MatInverse (LU with partial pivoting).
Inputs and units: `(m Mat4)`; m in unit A; the result in 1/A.
Output: `(Mat4, error)`.
Errors: ErrSingular when det m = 0; ErrNotFinite for a NaN or ±Inf element; ErrOverflow when an element of the inverse is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.2 "Permutations and Cofactors" and §5.3 "Cramer's Rule, Inverses, and Volumes" (A⁻¹ = Cᵀ/det A) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: translation by (1, 2, 3) → translation by (−1, −2, −3) (hand-checked)
Variants: Mat2, Mat3 and n×n forms.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/mat4.go`.
Notes: homogeneous transforms act on column vectors (Strang §10.6 uses rows: transpose).

### Mat4Trace      (ID reserved: NUM-LINALG-083)
Formula: Mat4Trace returns the sum of the diagonal elements.
Inputs and units: `(m Mat4)`; the unit of m.
Output: `(float64, error)`.
Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when the sum is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.1 "Introduction to Eigenvalues" (the trace equals the sum of the eigenvalues) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: diag(1, 2, 3, 4) → 10
Variants: Mat2, Mat3 and n×n forms.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/mat4.go`.
Notes: homogeneous transforms act on column vectors (Strang §10.6 uses rows: transpose).

### Mat4Translation      (ID reserved: NUM-LINALG-084)
Formula: Mat4Translation returns the homogeneous matrix that moves every point by offset: [I t; 0 1] with t = offset in the last column.
Inputs and units: `(offset Vec3)`; offset in a length unit; the matrix carries that unit in its last column.
Output: `(Mat4, error)`.
Errors: ErrNotFinite for a NaN or ±Inf component.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §10.6 "Computer Graphics", p. 495 (translation matrix T in homogeneous coordinates; transposed here for column vectors) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: offset (1, 2, 3) applied to the point (0, 0, 0) → (1, 2, 3)
Variants: Mat2, Mat3 and n×n forms.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/mat4_transform.go`.
Notes: homogeneous transforms act on column vectors (Strang §10.6 uses rows: transpose).

### Mat4Scale      (ID reserved: NUM-LINALG-085)
Formula: Mat4Scale returns the homogeneous matrix diag(sx, sy, sz, 1) that rescales the x, y and z directions.
Inputs and units: `(factors Vec3)`; factors dimensionless.
Output: `(Mat4, error)`.
Errors: ErrNotFinite for a NaN or ±Inf factor.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §10.6 "Computer Graphics", pp. 495–496 (scaling matrix S in homogeneous coordinates) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: factors (2, 3, 4) applied to (1, 1, 1) → (2, 3, 4)
Variants: Mat2, Mat3 and n×n forms.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/mat4_transform.go`.
Notes: homogeneous transforms act on column vectors (Strang §10.6 uses rows: transpose).

### Mat4FromMat3      (ID reserved: NUM-LINALG-086)
Formula: Mat4FromMat3 returns the homogeneous matrix [m 0; 0 1]: the linear map m (a rotation, say) in the upper-left block, no translation.
Inputs and units: `(m Mat3)`; the unit of m.
Output: `(Mat4, error)`.
Errors: ErrNotFinite for a NaN or ±Inf element.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §10.6 "Computer Graphics", p. 497 (the rotation Q in the upper-left block of the 4 by 4 matrix R) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: Mat3RotationZ(π/2) embedded, applied to (1, 0, 0, 1) → (0, 1, 0, 1)
Variants: Mat2, Mat3 and n×n forms.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/mat4_transform.go`.
Notes: homogeneous transforms act on column vectors (Strang §10.6 uses rows: transpose).

### Mat4TransformPoint      (ID reserved: NUM-LINALG-087)
Formula: Mat4TransformPoint returns the point p transformed by m: (x′, y′, z′, w′) = m·(p, 1), then divided by w′ (exactly 1 for an affine m; the perspective divide otherwise).
Inputs and units: `(m Mat4, p Vec3)`; p in the length unit of m's translation; the result in that unit.
Output: `(Vec3, error)`.
Errors: ErrOutOfRange when w′ = 0 (the point maps to infinity); ErrNotFinite for a NaN or ±Inf input; ErrOverflow when a component of the result is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §10.6 "Computer Graphics", p. 495 (points as (x, y, z, 1) in homogeneous coordinates) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: translation by (1, 2, 3) applied to (1, 1, 1) → (2, 3, 4)
Variants: Mat2, Mat3 and n×n forms.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/mat4_transform.go`.
Notes: homogeneous transforms act on column vectors (Strang §10.6 uses rows: transpose).

### Mat4TransformDirection      (ID reserved: NUM-LINALG-088)
Formula: Mat4TransformDirection returns the direction d transformed by m: the upper-left 3×3 block of m applied to d, ignoring translation (a direction is (x, y, z, 0) in homogeneous coordinates). For normals under a non-uniform scale, transform by the inverse transpose instead.
Inputs and units: `(m Mat4, d Vec3)`; d in any unit; the result in that unit.
Output: `(Vec3, error)`.
Errors: ErrNotFinite for a NaN or ±Inf input; ErrOverflow when a component of the result is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §10.6 "Computer Graphics", p. 497 (the linear part in the upper-left block; translation does not act on vectors) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: translation by (1, 2, 3) applied to the direction (1, 1, 1) → (1, 1, 1)
Variants: Mat2, Mat3 and n×n forms.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/mat4_transform.go`.
Notes: homogeneous transforms act on column vectors (Strang §10.6 uses rows: transpose).

### MatIdentity      (ID reserved: NUM-LINALG-089)
Formula: MatIdentity returns the n×n identity matrix in dst.
Inputs and units: `(dst []float64, n int)`; dimensionless.
Output: `([]float64, error)`.
Errors: ErrDimension when n < 0 or n·n overflows.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.5 "Inverse Matrices" (the identity I) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: n = 2 → [1 0; 0 1]
Variants: fixed-size Mat2/Mat3/Mat4 forms.
Best approach: as in Formula. Time O(n²)   n = order of the matrix. Space O(1) extra; O(n²) for the result, one allocation only if dst lacks capacity. Code: `num/linalg/mat.go`.
Notes: row-major []float64 with dimensions passed alongside; dst must not overlap the inputs.

### MatMul      (ID reserved: NUM-LINALG-090)
Formula: MatMul returns the product of the rows×inner matrix a and the inner×cols matrix b, a rows×cols matrix, in dst (which must not overlap a or b). Each element is accumulated with fused multiply-adds in i-k-j order (stride-1 access to b and dst).
Inputs and units: `(dst, a []float64, rows, inner int, b []float64, cols int)`; a in unit A, b in unit B; the result in A·B.
Output: `([]float64, error)`.
Errors: ErrDimension when a slice's length disagrees with its dimensions; ErrNotFinite for a NaN or ±Inf element; ErrOverflow when an element of the product is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.4 "Rules for Matrix Operations" (row times column) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: [1 2 3; 4 5 6] (2×3) · [7 8; 9 10; 11 12] (3×2) → [58 64; 139 154] (hand-checked)
Variants: fixed-size Mat2/Mat3/Mat4 forms.
Best approach: as in Formula. Time O(rows·inner·cols). Space O(1) extra; O(rows·cols) for the result, one allocation only if dst lacks capacity. Code: `num/linalg/mat.go`.
Notes: row-major []float64 with dimensions passed alongside; dst must not overlap the inputs.

### MatVecMul      (ID reserved: NUM-LINALG-091)
Formula: MatVecMul returns the product of the rows×cols matrix a and the column vector x (length cols), a vector of length rows, in dst (which must not overlap a or x).
Inputs and units: `(dst, a []float64, rows, cols int, x []float64)`; a in unit A, x in unit B; the result in A·B.
Output: `([]float64, error)`.
Errors: ErrDimension when a's length disagrees with its dimensions or len(x) ≠ cols; ErrNotFinite for a NaN or ±Inf element; ErrOverflow when an element of the result is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.4 "Rules for Matrix Operations" (matrix times vector) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: [1 2 3; 4 5 6]·(1, 0, −1) → (−2, −2) (hand-checked)
Variants: fixed-size Mat2/Mat3/Mat4 forms.
Best approach: as in Formula. Time O(rows·cols). Space O(1) extra; O(rows) for the result, one allocation only if dst lacks capacity. Code: `num/linalg/mat.go`.
Notes: row-major []float64 with dimensions passed alongside; dst must not overlap the inputs.

### MatTranspose      (ID reserved: NUM-LINALG-092)
Formula: MatTranspose returns the cols×rows transpose of the rows×cols matrix a in dst (which must not overlap a).
Inputs and units: `(dst, a []float64, rows, cols int)`; the unit of a.
Output: `([]float64, error)`.
Errors: ErrDimension when a's length disagrees with its dimensions; ErrNotFinite for a NaN or ±Inf element.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.7 "Transposes and Permutations" [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: [1 2 3; 4 5 6] → [1 4; 2 5; 3 6]
Variants: fixed-size Mat2/Mat3/Mat4 forms.
Best approach: as in Formula. Time O(rows·cols). Space O(1) extra; O(rows·cols) for the result, one allocation only if dst lacks capacity. Code: `num/linalg/mat.go`.
Notes: row-major []float64 with dimensions passed alongside; dst must not overlap the inputs.

### MatTrace      (ID reserved: NUM-LINALG-093)
Formula: MatTrace returns the sum of the diagonal of the n×n matrix a, summed by Sum2 (as accurate as in twice the working precision).
Inputs and units: `(a []float64, n int)`; the unit of a.
Output: `(float64, error)`.
Errors: ErrDimension when len(a) ≠ n·n; ErrNotFinite for a NaN or ±Inf diagonal element; ErrOverflow when the sum is outside the float64 range.
Primary source: Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.1 "Introduction to Eigenvalues" (the trace equals the sum of the eigenvalues) — and — Ogita, Rump, Oishi, "Accurate Sum and Dot Product", SIAM Journal on Scientific Computing, 2005, 26(6):1955–1988, doi:10.1137/030601818, Algorithm 4.4 (Sum2) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: [1 2 3; 4 5 6; 7 8 9] → 15
Variants: fixed-size Mat2/Mat3/Mat4 forms.
Best approach: as in Formula. Time O(n)   n = order of the matrix. Space O(1), no allocation. Code: `num/linalg/mat.go`.
Notes: row-major []float64 with dimensions passed alongside; dst must not overlap the inputs.

### LUDecompose      (ID reserved: NUM-LINALG-094)
Formula: LUDecompose returns the factorisation PA = LU of the n×n matrix a by Gaussian elimination with partial pivoting: at step k the row with the largest |entry| in column k is swapped up, so every multiplier has |l| ≤ 1. The input is not modified. A singular matrix still factors (some pivot is exactly 0); LUSolve then reports ErrSingular and MatDet returns 0.
Inputs and units: `(a []float64, n int)`; the unit of a.
Output: `(LU, error)`.
Errors: ErrDimension when n < 1 or len(a) ≠ n·n; ErrNotFinite for a NaN or ±Inf element; ErrOverflow when an element grows beyond the float64 range during elimination.
Primary source: Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §3.2 "The LU Factorization" and §3.4 "Pivoting" (partial pivoting) — and — Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.6 "Elimination = Factorization: A = LU" [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: [0 1; 2 3] → rows swapped (det P = −1), L = [1 0; 0 1], U = [2 3; 0 1] (hand-checked)
Variants: MatSolve, MatDet, MatInverse are built on it.
Best approach: as in Formula. Time O(n³)   n = order of the matrix; about n³/3 multiply-adds (2n³/3 flops). Space O(n²) for the factor (one allocation of n² floats and one of n ints). Code: `num/linalg/lu.go`.
Notes: partial pivoting keeps every multiplier ≤ 1 in magnitude.

### LUSolve      (ID reserved: NUM-LINALG-095)
Formula: LUSolve returns the solution x of Ax = b from the factorisation f = LUDecompose(A): x = U⁻¹L⁻¹Pb by forward and back substitution, written into dst (which must not overlap b).
Inputs and units: `(dst []float64, f LU, b []float64)`; b in unit B with A in unit A; x in B/A.
Output: `([]float64, error)`.
Errors: ErrSingular when A is singular; ErrDimension when f is empty or len(b) ≠ n; ErrNotFinite for a NaN or ±Inf element of b; ErrOverflow when an element of x is outside the float64 range.
Primary source: Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §3.1 "Triangular Systems" and §3.4 "Pivoting" — and — Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.6 "Elimination = Factorization: A = LU" [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: A = [2 1; 1 3], b = (3, 5) → x = (0.8, 1.4) (hand-checked)
Variants: MatSolve, MatDet, MatInverse are built on it.
Best approach: as in Formula. Time O(n²)   n = order of the matrix. Space O(1) extra; O(n) for the result, one allocation only if dst lacks capacity. Code: `num/linalg/lu.go`.
Notes: partial pivoting keeps every multiplier ≤ 1 in magnitude.

### MatSolve      (ID reserved: NUM-LINALG-096)
Formula: MatSolve returns the solution x of Ax = b for the n×n matrix a, by LU with partial pivoting, written into dst (which must not overlap b). To solve for several right-hand sides, factor once with LUDecompose and call LUSolve. For an ill-conditioned A the answer is inaccurate; check with MatConditionEstimate1.
Inputs and units: `(dst, a []float64, n int, b []float64)`; b in unit B with A in unit A; x in B/A.
Output: `([]float64, error)`.
Errors: as LUDecompose and LUSolve.
Primary source: Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §3.2 "The LU Factorization" and §3.4 "Pivoting" [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: A = [2 1; 1 3], b = (3, 5) → x = (0.8, 1.4) (hand-checked)
Variants: fixed-size Mat2/Mat3/Mat4 forms.
Best approach: as in Formula. Time O(n³)   n = order of the matrix. Space O(n²) for the factor; O(n) for the result, one allocation only if dst lacks capacity. Code: `num/linalg/lu.go`.
Notes: row-major []float64 with dimensions passed alongside; dst must not overlap the inputs.

### MatDet      (ID reserved: NUM-LINALG-097)
Formula: MatDet returns the determinant of the n×n matrix a as det P · Π uᵢᵢ from LU with partial pivoting. The product is kept as a mantissa and a separate binary exponent, so it does not overflow or underflow part-way. A singular matrix has determinant 0 (not an error).
Inputs and units: `(a []float64, n int)`; a in unit A; the result in Aⁿ.
Output: `(float64, error)`.
Errors: ErrDimension when n < 1 or len(a) ≠ n·n; ErrNotFinite for a NaN or ±Inf element; ErrOverflow when |det| exceeds the float64 range.
Primary source: Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §3.2 "The LU Factorization" (det A = det P · u₁₁⋯uₙₙ) — and — Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.1 "The Properties of Determinants" (the determinant is ± the product of the pivots) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: [1 2 3; 0 1 4; 5 6 0] → 1 (hand-checked)
Variants: fixed-size Mat2/Mat3/Mat4 forms.
Best approach: as in Formula. Time O(n³)   n = order of the matrix. Space O(n²) for the factor. Code: `num/linalg/lu.go`.
Notes: row-major []float64 with dimensions passed alongside; dst must not overlap the inputs.

### MatInverse      (ID reserved: NUM-LINALG-098)
Formula: MatInverse returns the inverse of the n×n matrix a, column by column from LU with partial pivoting, in dst (which must not overlap a). Forming an inverse is rarely needed to solve Ax = b (use MatSolve); for an ill-conditioned A the inverse is inaccurate, so check with MatConditionEstimate1.
Inputs and units: `(dst, a []float64, n int)`; a in unit A; the result in 1/A.
Output: `([]float64, error)`.
Errors: ErrSingular when a is singular; otherwise as LUDecompose; ErrOverflow when an element of the inverse is outside the float64 range.
Primary source: Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §3.2 "The LU Factorization" and §3.4 "Pivoting" — and — Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.5 "Inverse Matrices" (A⁻¹ column by column: Axⱼ = eⱼ) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: [1 2 3; 0 1 4; 5 6 0] → [−24 18 5; 20 −15 −4; −5 4 1] (hand-checked)
Variants: fixed-size Mat2/Mat3/Mat4 forms.
Best approach: as in Formula. Time O(n³)   n = order of the matrix. Space O(n²) for the factor and O(n) work; O(n²) for the result, one allocation only if dst lacks capacity. Code: `num/linalg/lu.go`.
Notes: row-major []float64 with dimensions passed alongside; dst must not overlap the inputs.

### MatRank      (ID reserved: NUM-LINALG-099)
Formula: MatRank returns the numerical rank of the rows×cols matrix a: the number of pivots larger than tol in Gaussian elimination with complete pivoting (at each step the largest |entry| of the remaining block is brought to the pivot position). Elimination stops when no remaining entry exceeds tol. tol < 0 selects the default max(rows, cols)·ε·max|aᵢⱼ| (ε = 2⁻⁵², machine epsilon): elimination's backward error is of that order, so pivots below it cannot be told from zero. tol = 0 counts every nonzero pivot. The input is not modified. The work copy is scaled so its largest |element| is in [1, 2), so elimination cannot overflow.
Inputs and units: `(a []float64, rows, cols int, tol float64)`; tol in the unit of a.
Output: `(int, error)`.
Errors: ErrDimension when len(a) ≠ rows·cols; ErrNotFinite for a NaN or ±Inf element or tol.
Primary source: Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §3.4 "Pivoting" (complete pivoting) — and — Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §3.2 "The Nullspace of A: Solving Ax = 0 and Rx = 0" (rank = number of pivots) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: [1 2 3; 2 4 6; 1 0 1] → 2 (the second row is twice the first; hand-checked)
Variants: fixed-size Mat2/Mat3/Mat4 forms.
Best approach: as in Formula. Time O(rows·cols·min(rows, cols)). Space O(rows·cols) work (one allocation). Code: `num/linalg/mat_rank.go`.
Notes: row-major []float64 with dimensions passed alongside; dst must not overlap the inputs.

### MatConditionEstimate1      (ID reserved: NUM-LINALG-100)
Formula: MatConditionEstimate1 returns an estimate of the 1-norm condition number κ₁(A) = ‖A‖₁·‖A⁻¹‖₁ of the n×n matrix a, without forming A⁻¹. ‖A⁻¹‖₁ is estimated by Hager's method with Higham's refinements: at most five iterations, each one solve with A and one with Aᵀ from a single LU factorisation, stopping when the estimate stops growing or the sign vector repeats, then Higham's alternating test vector bᵢ = (−1)ⁱ(1 + i/(n−1)). The estimate never exceeds the true κ₁ and is usually within a factor of 3 of it. Roughly log₁₀ κ₁ decimal digits are lost when solving Ax = b or inverting A; κ₁ near 1/ε (about 4.5e15) means the matrix is numerically singular.
Inputs and units: `(a []float64, n int)`; dimensionless.
Output: `(float64, error)`.
Errors: ErrSingular when a is exactly singular (κ₁ = ∞); ErrOverflow when the estimate exceeds the float64 range; otherwise as LUDecompose.
Primary source: Hager, "Condition Estimates", SIAM Journal on Scientific and Statistical Computing, 1984, 5(2):311–316, doi:10.1137/0905023 — and — Higham, "FORTRAN codes for estimating the one-norm of a real or complex matrix, with applications to condition estimation", ACM Transactions on Mathematical Software, 1988, 14(4):381–396, doi:10.1145/50063.214386 [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: diag(1, 1e−8) → 1e8 (exact for a diagonal matrix; hand-checked)
Variants: fixed-size Mat2/Mat3/Mat4 forms.
Best approach: as in Formula. Time O(n³)   n = order of the matrix (the LU factorisation; each iteration is O(n²)). Space O(n²) for the factor and O(n) work. Code: `num/linalg/mat_condition.go`.
Notes: row-major []float64 with dimensions passed alongside; dst must not overlap the inputs.

### QuatAdd      (ID reserved: NUM-LINALG-101)
Formula: QuatAdd returns p + q, component by component.
Inputs and units: `(p, q Quat)`; dimensionless.
Output: `(Quat, error)`.
Errors: ErrNotFinite for a NaN or ±Inf component; ErrOverflow when a component of the sum is outside the float64 range.
Primary source: Shoemake, "Animating Rotation with Quaternion Curves", Computer Graphics (SIGGRAPH '85), 1985, 19(3):245–254, doi:10.1145/325165.325242, §2.3 "Quaternions", p. 246 (addition: scalar to scalar, vector to vector) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, 2, 3, 4) + (5, 6, 7, 8) → (6, 8, 10, 12)
Variants: rotation matrices (Mat3) and axis-angle are equivalent representations.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/quat.go`.
Notes: Hamilton product; active rotation v ↦ qvq⁻¹; q and −q are the same rotation.

### QuatSub      (ID reserved: NUM-LINALG-102)
Formula: QuatSub returns p − q, component by component.
Inputs and units: `(p, q Quat)`; dimensionless.
Output: `(Quat, error)`.
Errors: as QuatAdd.
Primary source: Shoemake, "Animating Rotation with Quaternion Curves", Computer Graphics (SIGGRAPH '85), 1985, 19(3):245–254, doi:10.1145/325165.325242, §2.3 "Quaternions", p. 246 (componentwise addition, with q negated) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (5, 6, 7, 8) − (1, 2, 3, 4) → (4, 4, 4, 4)
Variants: rotation matrices (Mat3) and axis-angle are equivalent representations.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/quat.go`.
Notes: Hamilton product; active rotation v ↦ qvq⁻¹; q and −q are the same rotation.

### QuatMul      (ID reserved: NUM-LINALG-103)
Formula: QuatMul returns the Hamilton product p·q = [p₀q₀ − p⃗·q⃗, p₀q⃗ + q₀p⃗ + p⃗ × q⃗]. It is not commutative. For unit quaternions, the rotation of p·q applies q first, then p.
Inputs and units: `(p, q Quat)`; dimensionless.
Output: `(Quat, error)`.
Errors: ErrNotFinite for a NaN or ±Inf component; ErrOverflow when a component of the product is outside the float64 range.
Primary source: Shoemake, "Animating Rotation with Quaternion Curves", Computer Graphics (SIGGRAPH '85), 1985, 19(3):245–254, doi:10.1145/325165.325242, §2.3 "Quaternions", p. 246 (the general product rule) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: i·j → k; j·i → −k (hand-checked)
Variants: rotation matrices (Mat3) and axis-angle are equivalent representations.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/quat.go`.
Notes: Hamilton product; active rotation v ↦ qvq⁻¹; q and −q are the same rotation.

### QuatConjugate      (ID reserved: NUM-LINALG-104)
Formula: QuatConjugate returns q* = [q₀, −q⃗]. For a unit quaternion it is the inverse rotation.
Inputs and units: `(q Quat)`; dimensionless.
Output: `(Quat, error)`.
Errors: ErrNotFinite for a NaN or ±Inf component.
Primary source: Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §6.1, eq. (98) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, 2, 3, 4) → (1, −2, −3, −4)
Variants: rotation matrices (Mat3) and axis-angle are equivalent representations.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/quat.go`.
Notes: Hamilton product; active rotation v ↦ qvq⁻¹; q and −q are the same rotation.

### QuatNorm      (ID reserved: NUM-LINALG-105)
Formula: QuatNorm returns the magnitude |q| = √(W² + X² + Y² + Z²), without overflow or underflow in the squares. (Diebel's eq. 99 calls the squared value the norm; this returns the square root.)
Inputs and units: `(q Quat)`; dimensionless.
Output: `(float64, error)`.
Errors: ErrNotFinite for a NaN or ±Inf component; ErrOverflow when |q| exceeds the float64 range.
Primary source: Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §6.1, eq. (99) — and — Blue, "A Portable Fortran Program to Find the Euclidean Norm of a Vector", ACM Transactions on Mathematical Software, 1978, 4(1):15–23, doi:10.1145/355769.355771 [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, 2, 2, 4) → 5 (hand-checked)
Variants: rotation matrices (Mat3) and axis-angle are equivalent representations.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/quat.go`.
Notes: Hamilton product; active rotation v ↦ qvq⁻¹; q and −q are the same rotation.

### QuatNormalize      (ID reserved: NUM-LINALG-106)
Formula: QuatNormalize returns the versor (unit quaternion) q/|q|, which represents the same rotation as q.
Inputs and units: `(q Quat)`; dimensionless.
Output: `(Quat, error)`.
Errors: ErrZeroVector when q = 0; ErrNotFinite for a NaN or ±Inf component.
Primary source: Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §6.4, eq. (120) (unit quaternions) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, 2, 2, 4) → (0.2, 0.4, 0.4, 0.8) (hand-checked)
Variants: rotation matrices (Mat3) and axis-angle are equivalent representations.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/quat.go`.
Notes: Hamilton product; active rotation v ↦ qvq⁻¹; q and −q are the same rotation.

### QuatInverse      (ID reserved: NUM-LINALG-107)
Formula: QuatInverse returns q⁻¹ = q*/|q|², so that q·q⁻¹ = q⁻¹·q = 1.
Inputs and units: `(q Quat)`; dimensionless.
Output: `(Quat, error)`.
Errors: ErrZeroVector when q = 0; ErrNotFinite for a NaN or ±Inf component; ErrOverflow when a component of the inverse is outside the float64 range (q tiny).
Primary source: Shoemake, "Animating Rotation with Quaternion Curves", Computer Graphics (SIGGRAPH '85), 1985, 19(3):245–254, doi:10.1145/325165.325242, §3.2 "How quaternions rotate", pp. 247–248 (negate the vector part and divide by the magnitude squared) — and — Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §6.1, eq. (100) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (1, 1, 1, 1) → (0.25, −0.25, −0.25, −0.25) (hand-checked)
Variants: rotation matrices (Mat3) and axis-angle are equivalent representations.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/quat.go`.
Notes: Hamilton product; active rotation v ↦ qvq⁻¹; q and −q are the same rotation.

### QuatFromAxisAngle      (ID reserved: NUM-LINALG-108)
Formula: QuatFromAxisAngle returns the unit quaternion [cos(θ/2), sin(θ/2)·k] (k = axis/|axis|) of the rotation by angleRad about axis, right-hand rule. This is the "quaternion of rotation".
Inputs and units: `(axis Vec3, angleRad float64)`; axis in any unit, angleRad in radians; the result is dimensionless.
Output: `(Quat, error)`.
Errors: ErrZeroVector when axis = 0; ErrNotFinite for a NaN or ±Inf input.
Primary source: Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §6.12, eq. (175) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: axis ẑ, π/2 → (√2/2, 0, 0, √2/2) (hand-checked)
Variants: rotation matrices (Mat3) and axis-angle are equivalent representations.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/quat_axis_angle.go`.
Notes: Hamilton product; active rotation v ↦ qvq⁻¹; q and −q are the same rotation.

### QuatToAxisAngle      (ID reserved: NUM-LINALG-109)
Formula: QuatToAxisAngle returns the unit axis and the angle in [0, π] radians of the rotation represented by q (normalised first; q and −q give the same answer). The angle is 2·atan2(|q⃗|, |W|), accurate for small angles where Diebel's 2·acos(W) loses half its digits. For the identity rotation the axis is undefined and x̂ = (1, 0, 0) is returned with angle 0.
Inputs and units: `(q Quat)`; dimensionless quaternion; the axis is dimensionless, the angle in radians.
Output: `(axis Vec3, angleRad float64, err error)`.
Errors: ErrZeroVector when q = 0; ErrNotFinite for a NaN or ±Inf component.
Primary source: Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §6.13, eqs. (198)–(199) (with acos replaced by the equivalent atan2 of eq. (175)'s sine and cosine) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (√2/2, 0, 0, √2/2) → axis ẑ, angle π/2 (hand-checked)
Variants: rotation matrices (Mat3) and axis-angle are equivalent representations.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/quat_axis_angle.go`.
Notes: Hamilton product; active rotation v ↦ qvq⁻¹; q and −q are the same rotation.

### QuatFromMat3      (ID reserved: NUM-LINALG-110)
Formula: QuatFromMat3 returns the unit quaternion (with W ≥ 0) of the rotation matrix m (active, column vectors). Of the four candidate formulas it uses the one whose square root has the largest argument (1 + trace, or 1 + 2mᵢᵢ − trace), so no division is by a small number (Shepperd's method). m is assumed to be a rotation (orthogonal, det +1); it is not checked, and the result is normalised, so a matrix slightly off orthogonal from rounding gives the nearest-looking unit quaternion.
Inputs and units: `(m Mat3)`; dimensionless.
Output: `(Quat, error)`.
Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when sums of elements overflow (elements near the float64 maximum, which no rotation matrix has).
Primary source: Shepperd, "Quaternion from Rotation Matrix", Journal of Guidance and Control, 1978, 1(3):223–224, doi:10.2514/3.55767b — and — Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §6.5, eqs. (131)–(145) (for Diebel's frame matrix; the active matrix here is its transpose) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: Mat3RotationZ(π/2) → (√2/2, 0, 0, √2/2) (hand-checked)
Variants: rotation matrices (Mat3) and axis-angle are equivalent representations.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/quat_matrix.go`.
Notes: Hamilton product; active rotation v ↦ qvq⁻¹; q and −q are the same rotation.

### QuatToMat3      (ID reserved: NUM-LINALG-111)
Formula: QuatToMat3 returns the active rotation matrix of q (normalised first, so any nonzero q works): with s = 2/|q|², [1 − s(Y²+Z²)  s(XY − WZ)  s(XZ + WY); s(XY + WZ)  1 − s(X²+Z²)  s(YZ − WX); s(XZ − WY)  s(YZ + WX)  1 − s(X²+Y²)]. It is the transpose of Diebel's frame-rotation matrix for the same q.
Inputs and units: `(q Quat)`; dimensionless.
Output: `(Mat3, error)`.
Errors: ErrZeroVector when q = 0; ErrNotFinite for a NaN or ±Inf component.
Primary source: Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §6.4, eq. (125) (transposed: Diebel's matrix is the frame rotation) — and — Shoemake, "Animating Rotation with Quaternion Curves", Computer Graphics (SIGGRAPH '85), 1985, 19(3):245–254, doi:10.1145/325165.325242, §3.2 "How quaternions rotate", p. 248 (the same matrix, for row vectors) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (√2/2, 0, 0, √2/2) → [0 −1 0; 1 0 0; 0 0 1] (to rounding; hand-checked)
Variants: rotation matrices (Mat3) and axis-angle are equivalent representations.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/quat_matrix.go`.
Notes: Hamilton product; active rotation v ↦ qvq⁻¹; q and −q are the same rotation.

### QuatRotateVec3      (ID reserved: NUM-LINALG-112)
Formula: QuatRotateVec3 returns v rotated by q: v′ = q v q⁻¹ with v taken as the pure quaternion [0, v] (active rotation, right-hand rule; q is normalised first, so any nonzero q works). For a unit q this is computed as t = 2(q⃗ × v), v′ = v + W·t + q⃗ × t, which costs two cross products. Shoemake and Diebel write the frame rotation q⁻¹ v q; this is its inverse.
Inputs and units: `(q Quat, v Vec3)`; q dimensionless, v in any unit; the result in v's unit.
Output: `(Vec3, error)`.
Errors: ErrZeroVector when q = 0; ErrNotFinite for a NaN or ±Inf input; ErrOverflow when a component of the result is outside the float64 range.
Primary source: Shoemake, "Animating Rotation with Quaternion Curves", Computer Graphics (SIGGRAPH '85), 1985, 19(3):245–254, doi:10.1145/325165.325242, §3.2 "How quaternions rotate", p. 247 (rotation by conjugation, x′ = q⁻¹xq; the active form is its inverse) — and — Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §6.4, eqs. (121)–(127) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: q of π/2 about ẑ rotates x̂ to ŷ (hand-checked)
Variants: rotation matrices (Mat3) and axis-angle are equivalent representations.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/quat_rotate.go`.
Notes: Hamilton product; active rotation v ↦ qvq⁻¹; q and −q are the same rotation.

### QuatSlerp      (ID reserved: NUM-LINALG-113)
Formula: QuatSlerp returns the spherical linear interpolation between the rotations q0 and q1 at t ∈ [0, 1]: Slerp = sin((1−t)θ)/sin θ · q0 + sin(tθ)/sin θ · q1, with cos θ = q0·q1, which turns at a constant angular rate along the great arc. Both inputs are normalised first, and q1 is negated when q0·q1 < 0 so the interpolation takes the shorter of the two arcs (q and −q are the same rotation); at t = 1 the result may therefore be −q1. θ is computed as 2·atan2(|q0 − q1|, |q0 + q1|) (accurate for nearly equal rotations), and the weights as s·sinc(sθ)/sinc(θ), which equal (1−t, t) exactly when θ = 0, so no threshold is needed. At t = 0 the result is exactly the normalised q0.
Inputs and units: `(q0, q1 Quat, t float64)`; dimensionless quaternions, t dimensionless.
Output: `(Quat, error)`.
Errors: ErrOutOfRange when t is outside [0, 1]; ErrZeroVector when q0 or q1 is zero; ErrNotFinite for a NaN or ±Inf input.
Primary source: Shoemake, "Animating Rotation with Quaternion Curves", Computer Graphics (SIGGRAPH '85), 1985, 19(3):245–254, doi:10.1145/325165.325242, §3.3 "Great arc in-betweening", p. 248 (Slerp) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: identity to π/2 about ẑ at t = 0.5 → π/4 about ẑ (hand-checked)
Variants: rotation matrices (Mat3) and axis-angle are equivalent representations.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/quat_slerp.go`.
Notes: Hamilton product; active rotation v ↦ qvq⁻¹; q and −q are the same rotation.

### QuatFromEulerZYX      (ID reserved: NUM-LINALG-114)
Formula: QuatFromEulerZYX returns the unit quaternion of the yaw-pitch-roll rotation: W = cφcθcψ + sφsθsψ, X = sφcθcψ − cφsθsψ, Y = cφsθcψ + sφcθsψ, Z = cφcθsψ − sφsθcψ, where c and s are the cosine and sine of half of each angle.
Inputs and units: `(yawRad, pitchRad, rollRad float64)`; angles in radians; the result is dimensionless.
Output: `(Quat, error)`.
Errors: ErrNotFinite for a NaN or ±Inf angle.
Primary source: Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §5.6.5, eq. (84) (sequence (1, 2, 3): φ roll, θ pitch, ψ yaw) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: yaw π/2, pitch 0, roll 0 → (√2/2, 0, 0, √2/2) (hand-checked)
Variants: rotation matrices (Mat3) and axis-angle are equivalent representations.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/quat_euler.go`.
Notes: Hamilton product; active rotation v ↦ qvq⁻¹; q and −q are the same rotation.

### QuatToEulerZYX      (ID reserved: NUM-LINALG-115)
Formula: QuatToEulerZYX returns the yaw ψ ∈ (−π, π], pitch θ ∈ [−π/2, π/2] and roll φ ∈ (−π, π] of the rotation q (normalised first), by Diebel's formulas: φ = atan2(2(WX + YZ), W² − X² − Y² + Z²), θ = asin(2(WY − XZ)), ψ = atan2(2(XY + WZ), W² + X² − Y² − Z²). Gimbal lock: at pitch ±π/2 roll and yaw turn about the same axis and only their difference (θ = +π/2) or sum (θ = −π/2) is defined. When 1 − |2(WY − XZ)| ≤ lockTol the pitch is set to ±π/2, roll to 0, and yaw to 2·atan2(Z, W), which reproduces the rotation exactly. lockTol < 0 selects the default 2⁻⁵⁰ (four rounding units of the sine of the pitch). Near the lock, roll and yaw separately lose accuracy as the lock is approached (Diebel §5.6.6); the rotation they describe does not.
Inputs and units: `(q Quat, lockTol float64)`; q dimensionless, lockTol dimensionless; the angles in radians.
Output: `(yawRad, pitchRad, rollRad float64, err error)`.
Errors: ErrZeroVector when q = 0; ErrNotFinite for a NaN or ±Inf input.
Primary source: Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §5.6.3, eq. (73) and §5.6.6 "Singularities" (sequence (1, 2, 3): φ roll, θ pitch, ψ yaw) [confirmed 2026-11-12; see doc/linalg_2026-11-12/README.md, "Citations"]
Worked example: (√2/2, 0, 0, √2/2) → yaw π/2, pitch 0, roll 0 (hand-checked)
Variants: rotation matrices (Mat3) and axis-angle are equivalent representations.
Best approach: as in Formula. Time O(1). Space O(1), no allocation. Code: `num/linalg/quat_euler.go`.
Notes: Hamilton product; active rotation v ↦ qvq⁻¹; q and −q are the same rotation.
