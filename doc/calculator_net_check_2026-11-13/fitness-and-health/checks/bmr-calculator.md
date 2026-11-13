# bmr-calculator

## FIELDS (form calform, 2026-11-13)
cage 25 (15-80) · csex m|f (m) · cheightfeet 5 / cheightinch 10 · cpound 160 · metric cheightmeter 180, ckg 60 · coutunit c|j (c) ·
cformula m (Mifflin St Jeor, default) | h (Revised Harris-Benedict) | k (Katch-McArdle) · cfatpct 20.
Results: BMR (integer) and a table at six activity levels.

## MATCH
HAVE LIF-NUTR-002 (Mifflin simplified) — their default. PARTIAL: Revised Harris-Benedict, Katch-McArdle, activity table — NEEDS SOURCE
(see calorie-calculator.md and LEAD.md). MED-ANTHRO-010 is the ORIGINAL 1919 Harris-Benedict, not the revised one; not compared.

## COMPARE (scratch: checks/scratch/bmr-calculator.go; lb × 0.45359237, in × 2.54)
| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| m 25 y, 5'10", 160 lb, kcal | 1716.998 | 1,717 | 0 | MATCH |
| f 50 y, 5'4", 130 lb, kJ | 1194.670 kcal × 4.1868 = 5001.8 kJ | 5,002 | 0 | MATCH (they use the International Table calorie, 4.1868 J) |
Submissions: 1 fields + 2 runs.
