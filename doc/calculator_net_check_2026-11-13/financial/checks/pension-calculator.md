# pension-calculator

FIELDS: form `calform` (lump sum vs monthly): cretirementage (65), clumpsumamount (800,000), creturn (5 %/yr),
cmonthlyincome (5,000/month), ccola (3.5 %/yr). Form `calform2` (single life vs joint survivor): cretirementage,
clifeexpectancy (77), cspousesage (62), cSpouseslifeexpectancy (82), csinglelifepension (5,000), cjointsurvivorpension
(3,000), creturn, ccola. Form `calform3` (early vs late start): two retirement ages and incomes, creturn, ccola.

MATCH: PARTIAL. With COLA 0 the break-even age is AnnuityPeriods (FIN-INT-043) at the monthly rate equivalent to the
annual return (FIN-INT-044). COLA stepping yearly under monthly payments, joint-survivor and early/late comparisons
are household-planning conventions (D-021, `fin/personal`) with no single primary source; not built.

COMPARE (scratch: checks/scratch/pension-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 800,000 vs 5,000/month from 65, 5 %, COLA 0 | break-even 259.52 months → age 86.63 (effective-rate convention, as in retirement-calculator); 87.02 at nominal 5 %/12 | "age 87 or older" | 0 at whole years | MATCH |

Submissions: 1.
