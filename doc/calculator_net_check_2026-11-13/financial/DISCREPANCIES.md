# financial: discrepancies (A-FIN)

| slug | inputs | ours | theirs | verdict | citation |
|---|---|---|---|---|---|
| mortgage-calculator | 400k, 20 % down, 30 y, 7.263 % | total of payments 786,879.68 | 786,883.11 | THEIRS DIFFERS (they multiply the unrounded payment by n; ours sums the rounded schedule, last payment adjusted) | 12 CFR 1026.18(h) total of payments = sum of scheduled payments; D-022 |
| amortization-calculator | 200k, 15 y, 6 %, +200/month | interest saved 18,230.54; total interest 103,788.82 | 18,230.42; 103,788.46 | THEIRS DIFFERS (unrounded running cents vs per-period half-even rounding) | D-022; 12 CFR 1026.18(h) |
| refinance-calculator | balance 250,000, payment 1,800, 7 % | 285.69 → 286 payments (balance after 285 = 1,242.95) | 285 months, total 513,000 | THEIRS DIFFERS (they truncate the payment count and drop the final partial payment) | NCEES FE Reference Handbook 10.0.1 p. 230 (P/A,i,n); balance recursion FIN-LOAN-003 |
| apr-calculator | 25,750 note (750 financed fee), 300 upfront, 5 y, 8.5 % | APR 10.28726 % | 10.288 % | THEIRS DIFFERS (0.0007 pp; independent bisection agrees with ours) | 12 CFR 1026 Appendix J; 1026.22(a)(2) tolerance |
| bond-calculator | price 950, face 1,000, 4 % semiannual coupon, 10 y | YTM 4.63032 % | 4.6302 % | THEIRS DIFFERS (0.0001 pp; independent bisection 4.630325 %; their yield prices the bond at 950.0096) | QIA 3rd ed. ch. 1 §6.1 Eq. 11, §7.1; OpenStax Principles of Finance §10.2 |
| payback-period-calculator | −100,000; 30,000 growing 5 %/yr, 5 y; 10 % | payback 3.1562; discounted 3.9209 | 3.144; 3.943 | THEIRS DIFFERS (their growing-flow form uses an unidentified within-year rule; their listed-flow form matches ours exactly) | Dahlquist & Knight, Principles of Finance, OpenStax 2022, §16.1 and §16.4 |
| business-loan-calculator | 10,000, 10 %, 60 mo, fees 500 + 750 | APR 15.9326 % | 15.931 % | THEIRS DIFFERS (0.0016 pp; independent bisection 15.93265 %) | 12 CFR 1026 Appendix J; 1026.22(a)(2) |
| personal-loan-calculator | 8,000, 13.5 %, 36 mo, 300 fee deducted | APR 16.2026 % | 16.205 % | THEIRS DIFFERS (0.0024 pp; independent bisection 16.2026–16.2032 %) | 12 CFR 1026 Appendix J; 1026.22(a)(2) |
