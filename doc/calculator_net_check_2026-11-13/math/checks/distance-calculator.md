# distance-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope (2D and 3D Euclidean distance: num); distance on the Earth (calform3) needs an Earth model (radius or ellipsoid), which is geoscience: sci/**.
FIELDS (1 load): calform x11 (1), y11 (1), x12 (4), y12 (5); calform2 x21, y21, z21 (1, 1, 1), x22, y22, z22 (2, 2, 2); calform3 la1 (38.8976), lo1 (−77.0366), la2 (39.9496), lo2 (−75.1503) decimal degrees, or lad/lam/las/lau and lod/lom/los/lou (D° M′ S″, n/s, e/w).
MATCH: HAVE NUM-GEOM-001 PointDistance, NUM-LINALG-019 Vec3Distance; PARTIAL NUM-TRIG-031 CentralAngle (angle on a unit sphere; no Earth radius or ellipsoid in num, by design).
COMPARE (scratch/distance-calculator.go; 3 submissions)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| (−2.5, 7), (10, −3.25) | 16.165163160327207 | 16.165163160327 | 0 at 12 dp | match |
| (1, −2, 3.5), (−4, 6, 0) | 10.062305898749054 | 10.062305898749 | 0 at 12 dp | match |
| (38.8976, −77.0366), (39.9496, −75.1503) | central angle 0.03136559860849744 rad; × 6371.0088 km (IUGG mean radius, reference only) = 199.83 km | 200.0 km | 0.17 km | not compared: theirs appears to use an ellipsoid (the WGS 84 equatorial radius would give 200.05 km on a sphere, the mean radius 199.83) |

HANDOFF A-OTHER: Earth surface distance (geodesic on the WGS 84 ellipsoid, e.g. Vincenty 1975, Survey Review 23(176):88–93, doi:10.1179/sre.1975.23.176.88, or Karney 2013, J. Geodesy 87:43–55, doi:10.1007/s00190-012-0578-z; and a spherical great-circle wrapper on CentralAngle with a cited Earth radius) belongs in sci/geoscience. Their test pair and result above (200.0 km, 124.3 mi) are ready for that comparison.
