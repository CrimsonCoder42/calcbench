# average-return-calculator

FIELDS: form `calform` (dated money-weighted return): bb starting balance (5,600) + bd date, eb ending balance
(18,000) + ed date, up to 50 cash flows a_k (d|w) + m_k amount + d_k date.
Form `calform2` (average of returns held for different lengths): r_k rate %/yr, y_k years, n_k months (up to 60 rows).

MATCH: form 2 HAVE by composition: FIN-INT-044 (monthly rate from the yearly rate), FIN-INT-001 (growth over the
months), FIN-INVEST-034 CumulativeReturn, FIN-INT-040 GrowthRatePerPeriod + FIN-INT-029 (annualise).
Form 1 PARTIAL: a money-weighted return on dated flows (actual-day IRR) is not built; FIN-INT-039 needs equal
periods. Candidate source: CFA Institute, GIPS 2020 (money-weighted return). Not compared.

COMPARE (scratch: checks/scratch/average-return-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 10 %/yr for 1 y 2 mo; −2 %/yr for 5 mo; 15 %/yr for 2 y 3 mo | cumulative 51.7769 %; average 11.4990 %/yr | 51.78 %; 11.50 %/yr | 0 | MATCH |

Submissions: 1.
