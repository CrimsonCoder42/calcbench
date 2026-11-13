# auto-loan-calculator

FIELDS (after driver update 2; pass form '' — with `calform` every field came back MISSING): csaleprice (50,000) /
cmonthlypay (750) (solve either way), cloanterm (60 months), cinterestrate (5 %), cincentive (0), cdownpayment (10,000),
ctradeinvalue (0), ctradeinowned (0), cstate (NY…), csaletax (4 %), ctitlereg (2,200), cttrinloan (false).

MATCH: PARTIAL → composed. Payment FIN-INT-018. Amount financed = price − incentive − down − trade-in + amount owed
on the trade-in (+ tax and fees when financed): arithmetic. Their sales tax is on price − trade-in (7 % × 27,000);
whether a trade-in reduces taxable price is a state rule (jurisdiction table), not built.

COMPARE (scratch: checks/scratch/auto-loan-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 50k, 10k down, 5 %, 60 mo (page default; my first two runs with form `calform` sent nothing) | 754.85 | 754.85 | 0 | MATCH |
| 32k, 1k incentive, 3k down, trade 5k owing 2k, 7 % tax, 900 fees financed, 6.9 %, 72 mo → loan 27,790 | 472.46 | 472.46; loan 27,790.00 | 0 | MATCH |

DRIVER: form `calform` is listed by `fields` but `run` with it sends nothing (all MISSING); form '' works.
Submissions: 3 (2 sent nothing).
