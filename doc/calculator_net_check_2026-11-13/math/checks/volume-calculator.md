# volume-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope (geometry; their length-unit selects are unit conversion, not compared).
FIELDS (1 load): eleven unnamed forms: ballradius; coneradius, coneheight; cubeedge; ctankradius, ctanklength; rtanklength, rtankwidth, rtankheight; rcapsuleradius, rcapsuleheight; ccapbaseradius, ccapballradius, ccapheight; rconicalfrustumtopradius, …bottomradius, …height; rellipsoidaxis1..3; csquarepyramidedge, …height; ctubed1, ctubed2, ctubeheight; each with a "<field>unit" select (default feet).
MATCH: HAVE NUM-GEOM-063 SphereVolume, -058 ConeVolume, -049 CylinderVolume, -073 CubeVolume, -054 PyramidVolume, ConeFrustumVolume, EllipsoidVolume, SphericalCapVolume (num/geometry).
COMPARE (scratch/volume-calculator.go; 5 submissions)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| sphere r 2 | 33.51032163829113 | 33.510321638291 | 0 at 12 dp | match |
| cone r 3, h 7.5 | 70.68583470577036 | not captured | n/a | DRIVER |
| frustum r 2, R 5, h 4 | 163.36281798666926 | not captured | n/a | DRIVER |
| ellipsoid 1, 2, 3.5 | 29.321531433504738 | not captured | n/a | DRIVER |
| spherical cap R 5, h 2 | 54.454272662223076 | not captured | n/a | DRIVER |

DRIVER: volume-calculator (and likely surface-area-calculator): the forms have no name; `run` fills the field and reports `sent`, but the result text comes back empty for every form after the first (sphere). The driver needs a way to pick the form by the field it contains, or to submit the form that owns the filled field.
Outcome: MATCH on the one captured case; the rest left to a driver fix.

## Session 2 retry (2 submissions, total 7)
Cone (3, 7.5) and ellipsoid (1, 2, 3.5) with form `''`: `sent` correct, result text empty both times. The driver clicks the first Calculate button on the page (the sphere's), so a later shape's result is never produced. DRIVER persists: the button must be the one that belongs to the filled fields. One submission left; not used.

## Session 3 (fixed driver; 1 submission, total 8)
Ellipsoid semi-axes 1, 2, 3.5: ours 29.321531433504738, theirs 29.321531433505 — match at 12 dp. Cone not retried (limit reached); 2 of 5 shapes compared, both match.
