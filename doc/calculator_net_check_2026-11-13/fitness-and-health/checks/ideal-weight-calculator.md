# ideal-weight-calculator

## FIELDS (form calform, 2026-11-13)
cage 25 (2-80) · csex m|f (m) · cheightfeet 5 / cheightinch 10 · metric cheightmeter 180.
Results table (lb, 0.1): Robinson (1983), Miller (1983), Devine (1974), Hamwi (1964), Healthy BMI Range.

## MATCH
Was PARTIAL MED-ANTHRO-006 (Traub–Kichen, children). Adult equations: NONE → built.

## RESEARCH
Peterson CM, Thomas DM, Blackburn GL, Heymsfield SB. Am J Clin Nutr 2016;103(5):1197-1203, doi:10.3945/ajcn.115.121178 (PMC4841935,
read 2026-11-13), Table 3 — all four equations, with the original citations (originals not read first hand).

## BUILD
MED-ANTHRO-018 IdealBodyWeightDevine, -019 IdealBodyWeightRobinson, -020 IdealBodyWeightMiller (kg), -021 IdealBodyWeightHamwiLb (lb)
— med/anthropometry/ideal_body_weight.go; tests ideal_body_weight_test.go. Range heightIn ≥ 60.

## COMPARE (scratch: checks/scratch/ideal-weight-calculator.go; kg → lb ÷ 0.45359237)
| inputs | ours (lb) R / M / D / H | theirs (lb) R / M / D / H | verdict |
|---|---|---|---|
| man 5'10" | 156.53 / 154.98 / 160.94 / 166.00 | 156.5 / 155.0 / 160.9 / 165.3 | R, M, D MATCH; Hamwi THEIRS DIFFERS (0.7 lb) |
| woman 5'4" | 123.02 / 129.06 / 120.59 / 120.00 | 123.0 / 129.1 / 120.6 / 119.7 | R, M, D MATCH; Hamwi THEIRS DIFFERS (0.3 lb) |

## DISCREPANCY: THEIRS DIFFERS (Hamwi)
Their Hamwi values are reproduced by a metric restatement, men 48.0 kg + 2.7 kg/in (75.0 kg = 165.3 lb), women 45.5 kg + 2.2 kg/in
(54.3 kg = 119.7 lb). Ours is Hamwi's pound form, 106 + 6 lb/in (men), 100 + 5 lb/in (women), as printed in Peterson et al. 2016 Table 3.
Ours stays.
Submissions: 1 fields + 2 runs.

## WITHDRAWN (same session, 2026-11-13)
The two-source bar of D-075/D-076 (two independent authoritative secondaries that visibly cite the primary, not by each other's
authors nor the primary's) is not met, so the code was moved to _delete/med/anthropometry/ (files kept for restore when the PDF arrives)
and the row is NEEDS SOURCE. The comparison above stays as the record of what the withdrawn code produced.
