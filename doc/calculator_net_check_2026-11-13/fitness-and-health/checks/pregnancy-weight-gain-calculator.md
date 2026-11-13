# pregnancy-weight-gain-calculator

## FIELDS (form calform, 2026-11-13)
pstage Week 1-40 (20) · twins 0|1 (0) · cheightfeet 5 / cheightinch 6 · cpoundbefore 120 · cpoundnow 130 · metric cheightmeter 165,
ckgbefore 50, ckgnow 55. Results: recommended weight range at delivery (lb, 0.1), prepregnancy BMI and category, a week-by-week table.

## MATCH
NONE → built.

## RESEARCH
Institute of Medicine and National Research Council. Weight Gain During Pregnancy: Reexamining the Guidelines. National Academies
Press, 2009, doi:10.17226/12584, Summary, Table S-1 and the provisional twin ranges (nationalacademies.org/read/12584/chapter/2,
read 2026-11-13). Primary (the report itself). Note: their week-by-week table is an interpolation not in Table S-1 (first-trimester
0.5-2 kg assumption + weekly rate); not built, not compared.

## BUILD
MED-OB-005 GestationalWeightGainIOM, MED-OB-006 GestationalWeightGainTwinsIOM — med/obstetric/gestational_weight_gain_iom.go;
tests gestational_weight_gain_iom_test.go (both sides of 18.5, 25, 30).

## COMPARE (scratch: checks/scratch/pregnancy-weight-gain-calculator.go; BMI from MED-ANTHRO-002)
| inputs | ours | theirs | verdict |
|---|---|---|---|
| 5'6", 120 lb, singleton, week 40 | BMI 19.37 "Normal weight", +25–35 → 145.0–155.0 lb | 19.4 (Normal Weight), 145.0–155.0 lb | MATCH |
| 5'4", 190 lb, twins, week 40 | BMI 32.61 "Obese", +25–42 → 215.0–232.0 lb | 32.6 (Obese), 215.0–232.0 lb | MATCH |
Submissions: 1 fields + 2 runs.
