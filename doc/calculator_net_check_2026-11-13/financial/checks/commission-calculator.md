# commission-calculator

FIELDS: form `calc` (flat): price (200,000), rate (3 %), commission (solve-for). Form `calc2` (graduated): price2
(200,000), hasbc base commission (y|n) + bc (1,000) + bcunit (p|d), varyc graduated (y|n), fc flat (3) + unit,
tiers t_k (lower bound) / c_k (rate or amount) / u_k (p|d), up to 8 tiers (default $0–20,000 at 3 %, above at 5 %).

MATCH: NONE. A flat commission is one product; the graduated form charges each tier's rate on the part of the price
inside that tier (marginal), the same shape as a progressive tax schedule.
RESEARCH: 26 U.S.C. § 1(a) rate table (read 2026-11-13 at law.cornell.edu/uscode/text/26/1): "$5,535, plus 28% of
the excess over $36,900"; each row's base is the lower tiers in full at their own rates (5,535 = 15 % × 36,900).
The schedule itself is caller-supplied, so no year-specific table is built in (BRIEF: "a progressive tax on a
supplied bracket slice" is in scope).
BUILD: FIN-BIZ-055 TieredMarginalAmount, fin/business/tiered_marginal_amount.go:26. Fixed-amount tiers (u_k = d) and
a base commission are additions the caller makes.

COMPARE (scratch: checks/scratch/commission-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 200,000 at 3 % flat | 6,000.00 | 6,000 | 0 | MATCH |
| 200,000; 3 % to 20,000, 5 % above | 9,600.00 | 9,600 (600 + 9,000) | 0 | MATCH |

Submissions: 2.
