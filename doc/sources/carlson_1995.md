# Carlson 1995

**Citation**: Carlson, "Numerical computation of real or complex elliptic integrals", *Numerical Algorithms*, 1995, 10(1):13–26, doi:10.1007/BF02198293 (preprint arXiv:math/9409227).

Taken: §2, the R_D algorithm eqs. (35)–(41); E(k) = 2R_G(1 − k², 1, 0) and the AGM algorithm eqs. (42)–(46); R_G(0, 0, z) = √z/2; §3 check values. Used by NUM-GEOM-046 EllipsePerimeter, NUM-GEOM-070 (ellipsoid surface area). Checked: `TestCarlson_checkValues` reproduces R_D(0,2,1), R_D(2,3,4), R_G(2,3,4), R_G(0,16,16), R_G(0,0.0796,4). Status: preprint PDF read 2026-11-12; journal record (vol. 10, issue 1, pp. 13–26) via Crossref.
