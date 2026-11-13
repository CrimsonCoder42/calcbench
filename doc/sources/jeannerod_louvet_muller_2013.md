# Jeannerod et al.

**Citation**: Jeannerod, Louvet, Muller, "Further analysis of Kahan's algorithm for the accurate computation of 2×2 determinants", Mathematics of Computation, 2013, 82(284):2245–2264, doi:10.1090/S0025-5718-2013-02679-8.

**Taken**: Kahan's FMA algorithm for ad − bc with relative error ≤ 2u: w = bc, e = fma(−b, c, w), f = fma(a, d, −w), result f + e.

**Used by**: NUM-LINALG-005, NUM-LINALG-006, NUM-LINALG-016, NUM-LINALG-057, NUM-LINALG-058, NUM-LINALG-060.

**Checked**: DOI confirmed 2026-11-12 (also cited by num/algebra). Test: Vec2Cross and Vec2Dot within 2 ulps of math/big on 2000 cancelling inputs.

## num/algebra (math build, 2026-11-12)

**Used by**: NUM-ALGEB-006 QuadraticDiscriminant (and NUM-ALGEB-007, -008, -009 through it). **Checked**: AMS landing page and abstract; oracle asserts ≤ 2u vs math/big on 600 inputs, including near-double roots.
