# boat-loan-calculator

FIELDS (form ''): csaleprice (35,000) / cmonthlypay (350), cloanterm (10) + cloantermunit (y|m), cinterestrate (7 %),
cdownpayment (7,000) + unit (d|p), ctradeinvalue (0), csaletax (4 %), ctitlereg fees (2,200), cttrinloan (false).

MATCH: PARTIAL → composed, as auto-loan-calculator: amount financed is arithmetic (price − down − trade-in + tax +
fees when financed), payment FIN-INT-018. The sales-tax rate is a jurisdiction input.

COMPARE (scratch: checks/scratch/boat-loan-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 35,000, 7,000 down, 4 % tax, 2,200 fees financed, 7 %, 10 y → loan 31,600 | 366.90 | 366.90; loan 31,600.00 | 0 | MATCH |

Submissions: 1.
