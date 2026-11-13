# mortgage-calculator

FIELDS (from lead's proof run `out/fields_mortgage.json`, form `calform`): chouseprice (400,000), cdownpayment (20) +
cdownpaymentunit (p|d), cloanterm (30 years), cinterestrate (7.263 %), cstartmonth/cstartyear, caddoptional (true),
cpropertytaxes (1.2, %/yr) + unit, chomeins (1,500 $/yr) + unit, cpmi (0) + unit, choa (0) + unit, cothercost (4,000) + unit,
yearly increase % for tax/insurance/HOA/other, extra monthly/yearly/one-time payments, csbw (biweekly).

MATCH: HAVE FIN-INT-018 (payment), FIN-LOAN-001 (schedule, totals), FIN-LOAN-007 (PITI). Not covered: HOA and
"other costs" add-ons (plain sums), cost escalation, biweekly mode.

COMPARE (scratch: checks/scratch/mortgage-calculator.go; rate/12 per month, end-of-period payments)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 300k, 20 % down, 30 y, 6 % | pmt 1,438.92; total 518,012.55; interest 278,012.55 | 1,438.92; 518,011.65; 278,011.65 | pmt 0; total +0.90 | pmt MATCH; total THEIRS DIFFERS |
| 400k, 20 % down, 30 y, 7.263 %, tax 1.2 %/yr, ins 1,500/yr | pmt 2,185.79; PITI 2,710.79; total 786,879.68 | 2,185.79; 2,710.79; 786,883.11 | pmt 0; PITI 0; total −3.43 | pmt, PITI MATCH; total THEIRS DIFFERS |
| 150k, 10 % down, 15 y, 3.5 % | pmt 965.09; total 173,716.64 | 965.09; 173,716.46 | total +0.18 | pmt MATCH; total THEIRS DIFFERS |
| 200k, 0 down, 30 y, 0 % | pmt 555.56; total 200,000.00 | 555.56; 200,000.00 | 0 | MATCH |

Total of payments: theirs = n × unrounded payment (786,883.11 / 360 = 2,185.786...). Ours = the sum of the cents
actually paid on a schedule rounded each period, last payment absorbing the rounding (D-022; fin/loan doc.go).
Ours is the cash-paid figure Reg Z calls the "total of payments" (12 CFR 1026.18(h): the sum of the scheduled
payments). THEIRS DIFFERS, ours stays.
Submissions: 3 (case A from the lead's proof run).
