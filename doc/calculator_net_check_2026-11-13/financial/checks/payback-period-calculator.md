# payback-period-calculator

FIELDS: form `calform` (growing flow): initialinvestment1 (100,000), cashflow (30,000/yr), cashflowchange
(increase|decrease), cashflowchangerate (5 %/yr), years (5), discountrate (10 %).
Form `calform2` (listed flows): initialinvestment2 (100,000), discountrate2 (10 %), am1…am50 yearly flows
(5,000; 25,000; 35,000; 40,000; 30,000; 10,000).

MATCH: HAVE FIN-INVEST-022 PaybackPeriod, FIN-INVEST-023 DiscountedPaybackPeriod (interpolated within the year,
OpenStax Principles of Finance §16.1 and §16.4), FIN-INT-039 for their "cash flow return rate".

COMPARE (scratch: checks/scratch/payback-period-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| form 2 default flows, 10 % | payback 3.8750; discounted 5.4516; IRR 10.9637 % | 3.875; 5.452; 10.96 % | 0 | MATCH |
| form 1: −100,000; 30,000 growing 5 %/yr for 5 y, 10 % | payback 3.1562; discounted 3.9209; IRR 18.8617 % | 3.144; 3.943; 18.86 % | −0.012; +0.022; 0 | IRR MATCH; payback THEIRS DIFFERS |

Form 1: their table lists the same yearly flows (30,000; 31,500; 33,075; …) and the same cumulative balances, yet
the fractional year differs from straight-line interpolation (3 + 5,425/34,728.75 = 3.1562), which their own form 2
uses. Their within-year convention for form 1 could not be identified from the numbers. Ours follows OpenStax
§16.1 (payback with the last year's flow spread evenly) and §16.4 (discounted payback, same interpolation).
Submissions: 2.
