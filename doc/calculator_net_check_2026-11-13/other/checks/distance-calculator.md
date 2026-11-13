# Distance Calculator (row 55, HANDOFF from A-MATH)

CLASSIFY: in scope. The 2D/3D Euclidean forms are A-MATH's (MATCH, NUM-GEOM-001, NUM-LINALG-019; see
math/checks/distance-calculator.md). This row is the Earth-surface form (calform3).
FIELDS (from A-MATH's load): calform3 la1 38.8976, lo1 −77.0366, la2 39.9496, lo2 −75.1503 (decimal degrees),
or D° M′ S″ with N/S, E/W selects (lad/lam/las/lau, lod/lom/los/lou).
MATCH: PARTIAL NUM-TRIG-031 CentralAngle (unit sphere). BUILD sci/geoscience (planned GEO):
SCI-GEO-001 GreatCircleDistance (R·σ, any radius; MeanEarthRadiusMeters = GRS 80 R1 6 371 008.7714 m, Moritz),
SCI-GEO-002 VincentyInverseDistance (Vincenty, Survey Review 1975;23(176):88-93, doi:10.1179/sre.1975.23.176.88),
SCI-GEO-003 WGS84Distance (NIMA TR8350.2 Table 3.1: a = 6 378 137 m, 1/f = 298.257223563).
Source check: Vincenty's test line (a) on the Bessel ellipsoid reproduced, 14 110 526.1696 m vs 14 110 526.170 m
printed; equator 1° = a·π/180 exactly. Known limit: no convergence for nearly antipodal points (ErrNoConvergence;
Karney 2013, J Geodesy 87:43-55, doi:10.1007/s00190-012-0578-z, would remove it: a later build).
COMPARE (ours: checks/scratch/distance-calculator.go):

| inputs | ours WGS 84 | ours sphere R1 | theirs | verdict |
|---|---|---|---|---|
| (38.8976, −77.0366) → (39.9496, −75.1503) | 200.0501 km, 124.3054 mi | 199.8305 km | 200.0 km, 124.3 mi | MATCH (truncation) |
| (51.5074, −0.1278) → (40.7128, −74.0060) | 5585.2336 km, 3470.5032 mi | 5570.2298 km | 5,585.2 km, 3,470.5 mi | MATCH |
| (−33.8688, 151.2093) → (35.6762, 139.6503) | 7792.1748 km, 4841.8330 mi | 7825.8294 km | 7,792.1 km, 4,841.8 mi | MATCH (truncation) |

Convention: theirs is the ellipsoidal (WGS 84) geodesic, not a sphere (the sphere misses by 15-34 km), shown
truncated to 1 decimal (200.0501 → 200.0, 7792.17 → 7792.1), not rounded. 1 mi = 1609.344 m (NIST SP 811 B.8).
Submissions: 3 (session 3).
CLOSE: BUILT SCI-GEO-001..003 (MATCH; their display truncates to 0.1 km).
