# amortization-calculator

FIELDS (form `calform`): cloanamount (200,000), cloanterm (15 y) + cloantermmonth (0), cinterestrate (6 %),
caddoptional (false), cstartmonth/cstartyear, extra monthly cexma / yearly cexya / one-time cexoa with start
month-year, xa1..xa10 + xm/xy (dated one-off extras).

MATCH: HAVE FIN-LOAN-001 (schedule), FIN-INT-018 (payment), FIN-LOAN-005 (extra monthly payment: periods,
interest saved). Not covered: yearly / one-time / dated extras (FIN-LOAN-005 takes one level extra per period).

COMPARE (scratch: checks/scratch/amortization-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 200k, 15 y, 6 % | pmt 1,687.71; total 303,788.82; interest 103,788.82 | 1,687.71; 303,788.46; 103,788.46 | pmt 0; totals +0.36 | pmt MATCH; totals THEIRS DIFFERS (convention) |
| same + 200/month extra | 152 payments; interest saved 18,230.54; interest 85,558.28 | 152; 18,230.42; 85,558.04 | periods 0; saved +0.12; interest +0.24 | periods MATCH; cents THEIRS DIFFERS (convention) |
| 10k, 5 months, 12 % | pmt 2,060.40; total 10,302.00 | 2,060.40; 10,301.99 | total +0.01 | pmt MATCH; total THEIRS DIFFERS (convention) |

Convention: ours rounds each period's interest to the cent (half-even) and the last payment absorbs the residue
(D-022); the total is the cash actually paid (12 CFR 1026.18(h)). Theirs carries unrounded cents through the
totals. Same verdict as mortgage-calculator.
DRIVER: the per-period schedule table is not captured (result.tables holds only the summary table), so row-level
comparison was not possible.
Submissions: 3.
