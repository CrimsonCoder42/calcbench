# business-loan-calculator

FIELDS (form `calform`): cloanamount (10,000), cinterestrate (10 %), ccompound (monthly), cyears (5) + cmonths (0),
cpayback (month), coriginationfee (5 % of loan), cdocumentationfee (750), cotherfee (0).

MATCH: PARTIAL → composed: FIN-INT-018 payment; FIN-INT-032 APR with all fees as prepaid finance charges
(amount financed = loan − fees). Nothing new needed.

COMPARE (scratch: checks/scratch/business-loan-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 10,000, 10 %, 5 y, fees 5 % + 750 | pmt 212.47; APR 15.9326 % | 212.47; 15.931 % | APR −0.0016 pp | pmt MATCH; APR THEIRS DIFFERS |
| 150,000, 8.5 %, 7 y, fees 2 % + 500 + 250 | pmt 2,375.47; APR 9.3015 % | 2,375.47; 9.302 % | 0 | MATCH |

APR case 1: an independent bisection gives 15.93265 % (cent payment) and 15.93274 % (unrounded); both round to
15.933 %. Theirs is 0.0016 points lower, within Reg Z's 1/8-point tolerance (12 CFR 1026.22(a)(2)); same pattern
as apr-calculator. Ours stays.
Submissions: 2.
