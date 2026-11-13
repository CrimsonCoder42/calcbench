# S2: num/logs — the golden ratio

Status: SHIPPED 2026-11-12
Work folder: `doc/math_build_2026-11-12/`
Decisions: D-008 (eˣ−1, scalb, 10ˣ dropped: standard library), D-011 (golden ratio lives here) · Assumptions: —

## Functions

| ID | Name | Signature | Units |
|---|---|---|---|
| `NUM-LOGS-001` | `GoldenRatio` | `const GoldenRatio = 1.618…` (untyped, 63 digits) | dimensionless |

## Per function

### `NUM-LOGS-001` `GoldenRatio`
- **Source**: Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §1.2.8 (φ = (1+√5)/2);
  OEIS A001622, https://oeis.org/A001622 (read 2026-11-12) for the digits.
- **Best known**: a constant; O(1), no allocation.
- **Approach**: an untyped constant with 63 significant digits, so it keeps full precision in constant expressions
  and rounds to the nearest double (equal to `math.Phi`) as a float64. **Lost**: a function returning φ (adds a
  call for a fixed value); not shipping it and pointing callers at `math.Phi` (D-011 asked for a named, cited
  calcbench constant beside the other transcendental constants).
- **Input modified**: no inputs.
- **Worked example** (first test): φ² − (φ + 1) < 2⁻¹⁹⁰ in 256-bit big.Float; float64((1+√5)/2) == GoldenRatio ==
  math.Phi.
- **Errors**: none.

## Out of scope

eˣ, 10ˣ, eˣ−1, scalb (standard library, D-008). The golden-ratio sequence is `num/series` (geometry catalog
row 24, MOVE).
