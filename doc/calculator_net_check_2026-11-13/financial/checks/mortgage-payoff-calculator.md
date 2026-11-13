# mortgage-payoff-calculator

FIELDS: form `calform` (original loan known): cloanamount (400,000), cloanterm (30 y), cinterestrate (6 %),
cremainingyear (25) + cremainingmonth (0), cpayoffoption radio (together|extra|biweekly|original, default extra),
cadditionalmonth (500), cadditionalyear (0), cadditionalonetime (0).
Form `calform2` (original loan unknown): cunpaidprincipal (230,000), cmonthlypayment2 (1,500), cinterestrate2 (6 %),
cpayoffoption2 (extra|biweekly|original), cadditionalmonth2 (500), cadditionalyear2 (0), cadditionalonetime2 (0).

MATCH: form 1 HAVE FIN-LOAN-002 (balance), FIN-LOAN-005 (extra-payment payoff). Form 2 PARTIAL: payoff time for a
given payment needs a number-of-periods solver, which did not exist.

RESEARCH: N from the uniform-series present-worth factor (P/A,i,N) = PV/PMT (NCEES FE Reference Handbook 10.0.1,
2020, Engineering Economics factor table p. 230), solved with logarithms: N = −ln(1 − i·PV/PMT)/ln(1+i). OpenStax
Principles of Finance §8.2 was checked and does not solve for N (no second source found this session).

BUILD: FIN-INT-043 AnnuityPeriods, fin/interest/annuity_periods.go:25 (tests in annuity_periods_test.go).

COMPARE (scratch: checks/scratch/mortgage-payoff-calculator.go; i = 0.5 %/month)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| form 1: 400k, 30 y, 6 %, 25 y left, +500/mo | pmt 2,398.20; balance 372,217.56; ⌈N⌉ = ⌈206.04⌉ = 207 mo; FIN-LOAN-005 207 mo, saved 122,306.92 | 2,398.20; 372,217.43; 17 y 3 mo = 207; saved 122,306.20 | pmt 0; months 0; balance +0.13; saved +0.72 | MATCH on payment and months; cents THEIRS DIFFERS (rounding convention, see mortgage-calculator) |
| form 2: 230k, 1,500/mo, 6 % | N = 291.78 → 292 mo | 24 y 4 mo = 292 | 0 | MATCH |
| form 2 + 500/mo | N = 171.56 → 172 mo | 14 y 4 mo = 172 | 0 | MATCH |

Their total-interest figures for form 2 (207,677.36; 113,122.63) come from a month-by-month schedule; we have no
schedule for a given payment (BalloonBalance gives balances only), so they were not compared.
Submissions: 2.
