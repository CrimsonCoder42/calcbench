# Higham, Accuracy and Stability of Numerical Algorithms

**Citation**: Higham, *Accuracy and Stability of Numerical Algorithms*, 2nd ed., SIAM, 2002, §4.3.

Taken: the error behaviour of compensated summation, used to set the oracle tolerance in algo/window tests. Used by ALG-WINDOW-001, -002.

## num/algebra and num/prob (math build, 2026-11-12)

| Taken | Locator | Used by | How checked |
|---|---|---|---|
| Solving a quadratic equation (cancellation) | §1.8 "Solving a Quadratic Equation" | NUM-ALGEB-007 | Vieta and exact-discriminant-sign oracles |
| Computing the sample variance | §1.9 "Computing the Sample Variance" | NUM-PROB-019 | Hand-checked examples |
| γₖ = ku/(1 − ku) | §3.1 | tolerances in the NUM-ALGEB-001, -002, -003 tests | Oracles vs math/big within γ₂ₙ |
| Compensated summation | §4.3 | the unexported `compensatedSum` in num/prob (comment) | As ALG-WINDOW-001 |
| Horner's method, backward error γ₂ₙ | §5.1 "Horner's Method", p. 94 | NUM-ALGEB-001 | Oracle vs math/big within γ₂ₙ |
| Evaluating derivatives | §5.2 "Evaluating Derivatives" | NUM-ALGEB-002 | Oracle vs math/big |

Status: section titles and p. 94 confirmed by search snippets of the table of contents (SIAM page returned 403); §1.8's page NOT confirmed, so it is cited by section only; §1.9 confirmed from a library TOC.

## Added at the 2026-11-12 merge (linear algebra pass)

§6.1 "Vector Norms": definitions of the 1-, ∞- and p-norms; used by NUM-LINALG-042, -043, -044; section confirmed in the table of contents 2026-11-12.
