# irr-calculator

FIELDS: form `calform` (yearly flows): bb initial investment (50,000), m1…m50 year-end flows (−10,000; 30,000;
50,000). Form `calform2` (lump sum with periodic deposits/withdrawals): iv (10,000), hy (2) + hm (6) holding,
ir ending balance (15,000), po (d|w), pa (100), pf frequency (a|sa|q|m|sm|bw|w), pt timing (b|e).

MATCH: HAVE FIN-INT-039 InternalRateOfReturn (equal periods). Form 2 maps to FIN-INT-042 AnnuityRate (PV, PMT, FV)
when the holding length is a whole number of payment periods; not compared this session.

COMPARE (scratch: checks/scratch/irr-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| −50,000; −10,000; 30,000; 50,000 | 12.4464 % | 12.446 % | 0 | MATCH |
| −100,000; 30,000; 40,000; 50,000; 20,000 | 15.3221 % | 15.322 % | 0 | MATCH |

Submissions: 2.
