# personal-loan-calculator

FIELDS (form `calform`): cloanamount (20,000), cinterestrate (10 %), cyears (5) + cmonths (0), cstartmonth +
cstartyear, caddoptional (false), corigamount (5) + corigisa (percentage|amount), corigpaid (deduct|upfront),
cinsurance (0).

MATCH: PARTIAL → composed: FIN-INT-018 payment; FIN-INT-032 APR with the origination fee as a prepaid finance charge.

COMPARE (scratch: checks/scratch/personal-loan-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 20,000, 10 %, 5 y, 5 % fee deducted | pmt 424.94; APR 12.2387 % | 424.94; 12.239 % | 0 | MATCH |
| 8,000, 13.5 %, 3 y, 300 fee deducted | pmt 271.48; APR 16.2026 % | 271.48; 16.205 % | APR +0.0024 pp | pmt MATCH; APR THEIRS DIFFERS |

Case 2: an independent bisection gives 16.20261 % (cent payment) and 16.20321 % (unrounded); neither reaches
16.205 %. Within Reg Z's 1/8-point tolerance (12 CFR 1026.22(a)(2)); same pattern as apr-calculator. Ours stays.
Submissions: 2.
