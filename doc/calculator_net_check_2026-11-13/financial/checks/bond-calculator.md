# bond-calculator

FIELDS: form `calform` (leave one blank to solve for it): cprice (blank), cface (100), cyield (6 %), cyears (3),
ccoupon (5) + cunit (p|d), cfrequency (a|s|q|m). Form `calform2` (dated): cface2, cyield2, ccoupon2 + unit,
cfrequency2, cmdate (maturity), csdate (settlement), cset day count radio (b|c|n|a).

MATCH: HAVE FIN-INVEST-002 BondPrice, FIN-INVEST-003 YieldToMaturity. PARTIAL for form 2: dirty/clean price and
accrued interest between coupon dates are not built (FIN-INT-036/037 give the year fractions). Not compared.

COMPARE (scratch: checks/scratch/bond-calculator.go; face 1,000 so the cents carry their 4 decimals)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| face 1,000, 5 % annual coupon, 3 y, yield 6 % | price 973.27 | 973.2699 | 0 at cents | MATCH |
| price 950, face 1,000, 4 % semiannual, 10 y | YTM 4.63032 % | 4.6302 % | −0.0001 pp | THEIRS DIFFERS |

YTM: an independent bisection gives 4.630325 %; at their 4.6302 % the price is 950.0096, about a cent off, so theirs
reads as a coarse solver stop. Source for ours: QIA 3rd ed. ch. 1 §6.1 Eq. 11 with §7.1 (solving for r); OpenStax
Principles of Finance §10.2.
Submissions: 2.
