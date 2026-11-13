# matrix-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope.
FIELDS (2 loads): form calcf: arow (4), acolumn (4), a00…a33 (empty), pa (2, power), ta (3, scalar); brow (4), bcolumn (4), b00…b33, pb (2), tb (3). Operations are separate buttons on the page (A+B, A−B, AB, det, inverse, transpose, power, …).
MATCH: HAVE NUM-LINALG-090 MatMul, -097 MatDet, -098 MatInverse, -122 MatPower (general n×n, row-major []float64).
COMPARE: none. `run` exit 4 (no button / no result block found) on a 4×4 A, B submission (1 submission). The page's operation buttons are script handlers with no single Calculate button.
DRIVER: matrix-calculator — per-operation JS buttons (no form submit); the driver needs a way to name the button to click.
Outcome: BLOCKED (driver gap), not a site block.

## Session 2 (2026-11-13 09:48, driver update 2)
1 field load + 1 submission (total 2 loads, 2 submissions). The page's operation buttons are labelled Transpose, Power of, Determinant, Inverse (labels read once, nothing else stored). `run … '' Determinant` on a 4×4 A returned `submit: "live"` and exit 4: none of the driver's button selectors (input submit/button by value, `<button>` by text) matched the "Determinant" element, so the driver fell back to live mode and found no result block. Multiply's label could not be targeted either.
DRIVER: matrix-calculator — the operation controls are not `input`/`button` elements; the 4th-argument label match needs to also try other clickable elements (e.g. an element with that visible text and an onclick).
Outcome: stays BLOCKED (driver gap). Our NUM-LINALG-090/097/098 not compared.

## Session 3 (fixed driver; 3 submissions, total 5)
A = [[2,1,0,3],[1,3,2,0],[0,2,4,1],[5,0,1,2]], B = [[1,2,0,−1],[0,1,3,2],[4,0,1,1],[−2,1,0,3]] (row-major in scratch/matrix-calculator.go).

| operation | ours | theirs | diff | verdict |
|---|---|---|---|---|
| det A (NUM-LINALG-097) | −117.00000000000001 | −117 | 1 ulp-scale | match |
| A⁻¹ (NUM-LINALG-098) | −0.11111111111111109 0.11111111111111109 −0.11111111111111112 0.22222222222222224 / 0.14529914529914528 0.39316239316239315 −0.16239316239316243 −0.13675213675213677 / −0.16239316239316237 −0.1452991452991453 0.2991452991452992 0.09401709401709403 / 0.3589743589743589 −0.2051282051282051 0.1282051282051282 −0.10256410256410256 | same to their 10 significant figures | 0 at displayed figs | match |
| AB (NUM-LINALG-090) | −4 8 3 9 / 9 5 11 7 / 14 3 10 11 / 5 12 1 2 | identical | 0 | match |

Outcome: MATCH (3/3).
