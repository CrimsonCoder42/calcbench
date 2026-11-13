# fraction-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope (rational arithmetic).
FIELDS (`fields`, 1 load): calform t1/b1 (2/7), op select {+,-,*,/}, t2/b2 (3/8); calform5 c5num1 "-2 3/4", c5op, c5num2 "3 5/7" (mixed numbers); calform4 c4l1/c4t1/c4b1 (mixed simplify 2 21/98); calform2 c2d1 1.375 (decimal to fraction); calform3 c3t1/c3b1 2/7 (fraction to decimal); calform6 big-integer fractions.
MATCH: PARTIAL. NUM-INTEG-004 GCD, ALG-INTMTH-013 FractionToDecimal (calform3). No rational add/sub/mul/div/reduce. Decimal-to-fraction (calform2) and big-integer fractions (calform6, would need math/big) not built.
RESEARCH: Knuth, TAOCP Vol. 2, 3rd ed., 1998, §4.5.1 "Fractions", eqs. (2) product and (3) sum with gcd cancellation.
BUILD: NUM-INTEG-007..011 FractionReduce/Add/Sub/Mul/Div, num/integer/fraction.go; tests num/integer/fraction_test.go (hand examples, errors, math/big.Rat oracle, alloc, benchmarks). Mixed numbers: caller converts w n/d to (w·d ± n)/d.

COMPARE (scratch/fraction-calculator.go; 4 submissions)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 2/7 + 3/8 | 37/56 | 37/56 | 0 | match |
| 5/12 ÷ 10/9 | 3/8 | 3/8 | 0 | match |
| 1/6 − 5/6 | −2/3 | −2/3 | 0 | match |
| −2 3/4 + 3 5/7 (= −11/4 + 26/7) | 27/28 | 27/28 | 0 | match |

Not covered here: decimal-to-fraction (calform2) and arbitrary-size fractions (calform6); left as gaps in LEAD.md.
