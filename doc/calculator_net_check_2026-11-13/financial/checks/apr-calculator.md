# apr-calculator

FIELDS: form `calform` (general): cloanamount (100,000), cloanterm (10 y) + cloantermmonth (0), cinterestrate (6 %),
ccompound (annually|semiannually|quarterly|monthly|semimonthly|biweekly; monthly), cpayback (daily…year; month),
cloanedfees (0, added to the loan), cfrontfees (2,500, paid upfront).
Form `calform2` (mortgage): c2housevalue (350,000), c2downpayment (20 %), c2loanterm (30), c2interestrate (6.2 %),
c2loanfees (3,500), c2points (0.5), c2pmi (0).

MATCH: HAVE FIN-INT-032 AnnualPercentageRate (Reg Z Appendix J, equal payments) with FIN-INT-018 for the payment.
Amount financed = note amount − prepaid finance charges (financed fees, upfront fees, points). Compounding other
than the payment period is not covered (all cases here use monthly = payment period).

COMPARE (scratch: checks/scratch/apr-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 100k, 10 y, 6 %, upfront 2,500 | pmt 1,110.21; APR 6.56279 % | 1,110.21; 6.563 % | 0 | MATCH |
| 25k + 750 financed fee, 5 y, 8.5 %, upfront 300 | pmt 528.30; APR 10.28726 % | 528.30; 10.288 % | 0.0007 pp | THEIRS DIFFERS (solver precision) |
| 280k (350k, 20 % down), 30 y, 6.2 %, fees 3,500, 0.5 point | pmt 1,714.91; APR 6.36746 % | 1,714.91; 6.367 % | 0 | MATCH |

Case 2: an independent bisection (hand script) gives 10.28726 % with the cent payment and 10.28732 % with the
unrounded one; both round to 10.287 %. Theirs is 0.0007 percentage points higher, well inside Reg Z's 1/8-point
tolerance (12 CFR 1026.22(a)(2)), so it reads as a coarser solver stop. Ours stays.
Submissions: 3.
