# college-cost-calculator

FIELDS (form ''): todaycost (30,990/yr) + useaverage select (30990|65470|50920|21320 presets), costincrease (5 %/yr),
collegelength (4 y), savingpercent (35 %), balancenow (0), returnrate (5 %), interesttaxrate (25 %), startin (3 y).

MATCH: PARTIAL by composition: each year's cost = today × (1+g)^t (FIN-INT-001); the savings target is a share of
their sum; the saving needed per period is FIN-INT-019 AnnuityPaymentFromFutureValue at the after-tax return
r(1 − tax). The sum of growing costs could also be FIN-INT-023 (growing annuity). Nothing built.

COMPARE: not possible. DRIVER: after submit the result block holds only "Results:" (text and tables empty; two
runs). The figures are probably drawn into a chart or a late-loading element.
Submissions: 2.
