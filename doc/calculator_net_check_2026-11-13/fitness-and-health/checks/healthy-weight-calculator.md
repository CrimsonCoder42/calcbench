# healthy-weight-calculator

## FIELDS (form calform, 2026-11-13)
cheightfeet 5 / cheightinch 10 · metric cheightmeter 180. Result: healthy weight range (integer lb) and a band of further weights.

## MATCH
Was PARTIAL MED-ANTHRO-001/-003 (no inversion). Built the inversion.

## BUILD
MED-ANTHRO-022 WeightAtBMI (kg from m), MED-ANTHRO-023 WeightAtBMIImperial (lb from in, ÷ 703) — med/anthropometry/weight_at_bmi.go;
tests weight_at_bmi_test.go. Source: CDC BMI FAQ formula solved for weight; range edges 18.5 and 25 from the CDC adult categories
(MED-ANTHRO-003). (Peterson et al. AJCN 2016 eq. 13/21 universal approximation was seen and NOT built — the exact inversion is enough.)

## COMPARE (scratch: checks/scratch/healthy-weight-calculator.go; in × 0.0254 m, kg ÷ 0.45359237 lb)
| height | ours WeightAtBMI 18.5–25 (lb) | ours Imperial ×703 (lb) | theirs | verdict |
|---|---|---|---|---|
| 5'2" | 101.148–136.686 | 101.158–136.700 | 101–137 | MATCH |
| 6'4" | 151.985–205.385 | 152.000–205.405 | 152–205 | MATCH |
| 5'10" (BMI page, row 1 run) | 128.935–174.236 | 128.947–174.253 | 128.9–174.2 | MATCH (metric); imperial form 174.25 → 174.3 at 0.1: the CDC 703 factor, documented |
Their range follows the metric inversion. Submissions: 1 fields + 2 runs.
