# bmi-calculator

## FIELDS (form calform, 2026-11-13)
cage text 25 (ages 2-120) · csex radio m|f (m) · cheightfeet 5 · cheightinch 10 · cpound 160 · (metric tab) cheightmeter 180 cm · ckg 65.
Results shown: BMI (0.1), category, healthy BMI range, healthy weight for the height, lb to gain/lose, BMI Prime, Ponderal Index.

## MATCH
HAVE MED-ANTHRO-002 BodyMassIndexImperial (CDC × 703), MED-ANTHRO-001 BodyMassIndex, MED-ANTHRO-003 BMICategoryCDCAdult.
Not ours: "healthy weight for the height" (row 9), BMI Prime (BMI/25) and Ponderal Index (kg/m³) — not in scope of this row; noted in LEAD.md.

## COMPARE (scratch: checks/scratch/bmi-calculator.go)
| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 160 lb, 5 ft 10 in | 22.9551, "Healthy Weight" | 23 (22.96→23.0), "Normal" | 0 at 0.1 | MATCH (value) |
| 250 lb, 5 ft 4 in | 42.9077, "Class 3 Obesity (Severe Obesity)" | 42.9, "Obese Class III" | 0 | MATCH (value) |
| 100 lb, 6 ft 0 in | 13.5610, "Underweight" | 13.6, "Severe thinness" | 0 | MATCH (value) |

Category labels: theirs uses WHO-style labels (incl. thinness grades below 18.5); ours returns the CDC adult labels
verbatim from the cited CDC page. Same cut points at 18.5/25/30/35/40. THEIRS DIFFERS (label set), ours stays.
Submissions: 4 (1 fields + 4 runs).
