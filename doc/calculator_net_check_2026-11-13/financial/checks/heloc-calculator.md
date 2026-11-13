# heloc-calculator

FIELDS: form `calform`: cloanamount (50,000), cinterestrate (8 %), cdrawperiod (5 y), crepayperiod (15 y),
caddoptional (false), cclosecost (2,000) + cclosecostisa (percentage|amount), cclosecostpaid (upfront|deduct),
cannualfee (50). Form `calform2`: chousevalue, cmbalance, ltv (same as home-equity-loan-calculator).

MATCH: PARTIAL → composed: draw period (interest only, full line drawn at once) = FIN-INT-035 SimpleInterest for
1/12 year; repayment = FIN-INT-018 over the repayment term. No two-phase function is built: the composition is two
calls, and real HELOC draws vary, which the caller models.

COMPARE (scratch: checks/scratch/heloc-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 50k, 8 %, draw 5 y, repay 15 y | draw 333.33; repay 477.83; total 106,009.20 | 333.33; 477.83; 106,008.69 | total +0.51 | pmts MATCH; total THEIRS DIFFERS (unrounded × n, see mortgage-calculator) |
| 120k, 9.25 %, draw 10 y, repay 20 y | draw 925.00; repay 1,099.04; total 374,769.60 | 925.00; 1,099.04; 374,769.65 | total −0.05 | pmts MATCH; total convention |

Submissions: 2.
