# surface-area-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope (unit selects not compared).
FIELDS (1 load): calform ballradius; calform2 coneradius, coneheight; calform3 cubeedge; calform4 ctankradius, ctanklength; calform5 rtanklength, rtankwidth, rtankheight; calform6 rcapsuleradius, rcapsuleheight; calform7 ccapbaseradius, ccapballradius, ccapheight; calform8 rconicalfrustumtop/bottomradius, height; calform9 rellipsoidaxis1..3; each with a unit select.
MATCH: HAVE NUM-GEOM-064 SphereSurfaceArea, -060 ConeSurfaceArea, -051 CylinderSurfaceArea, -074 CubeSurfaceArea, EllipsoidSurfaceArea (exact, Carlson R_G).
COMPARE (scratch/surface-area-calculator.go; 4 submissions)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| sphere r 2 | 50.26548245743669 | 50.265482457437 | 0 at 12 dp | match |
| cone r 3, h 7.5 | 104.40530776639636 | 104.4053077664 | 0 at 10 dp | match |
| cylinder r 1.5, h 4 | 51.83627878423159 | 51.836278784232 | 0 at 12 dp | match |
| ellipsoid 1, 2, 3.5 | 56.14575157027773 | ≈ 56.193275605353 | 0.0475 (0.085 %) | THEIRS DIFFERS |

DISCREPANCY: THEIRS DIFFERS. Their page shows the approximation 4π((aᵖbᵖ + aᵖcᵖ + bᵖcᵖ)/3)^(1/p), p = 1.6 (the "≈" and the formula are on their result; it reproduces 56.1932756 exactly), whose error is up to about 1 %. Ours is the exact area S = 3V·R_G(a⁻², b⁻², c⁻²): NIST DLMF eq. 19.33.1 with Carlson's R_D algorithm (Numerical Algorithms 1995, 10(1):13–26, doi:10.1007/BF02198293). Ours stays.
