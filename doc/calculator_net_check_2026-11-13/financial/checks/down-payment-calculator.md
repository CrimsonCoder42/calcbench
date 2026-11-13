# down-payment-calculator

FIELDS: three forms. `calform` (price from cash): ccashamount1 (100,000), cdownpaypct1 (20), cincludeclosingcost1
(true), cclosingcost1 (3) + unit (p|d), cinterest1 (7.263 %), cloanterm1 (30). `calform2` (cash from price):
c2homeprice (500,000), c2downpaypct (20), closing cost (3 %), rate, term. `calform3` (down % from price and cash):
c3homeprice (500,000), c3downpaypct = cash available (100,000), closing (3 %), rate, term.

MATCH: PARTIAL → composed. Down share = 1 − LTV (FIN-LOAN-006); payment FIN-INT-018. Price from cash is
cash / (down % + closing %), one division; not built.

COMPARE (scratch: checks/scratch/down-payment-calculator.go; one submission returned all three forms' results,
forms 2 and 3 at their defaults)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| form 1: cash 100k, 20 % down, 3 % closing, 7.263 %, 30 y | price 434,782.61; loan 347,826.09; pmt 2,375.85 | 434,783; 347,826; 2,376 | 0 at whole dollars | MATCH |
| form 2: 500k, 20 % down | loan 400,000; pmt 2,732.23 | 400,000; 2,732 | 0 | MATCH |
| form 3: 500k, cash 100k, 3 % closing | LTV 0.83 → down 17.0 %; pmt 2,834.69 | 17.0 %; 2,835 | 0 | MATCH |

Submissions: 1.
