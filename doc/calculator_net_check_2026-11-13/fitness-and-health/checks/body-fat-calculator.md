# body-fat-calculator

## FIELDS (form calform, 2026-11-13)
csex m|f (m) · cage 25 · cweightlbs 152 · cheightfeet 5 / cheightinch 10.5 · cneckfeet 1 / cneckinch 7.5 · cwaistfeet 3 / cwaistinch 1.5 ·
chipfeet 2 / chipinch 10.5 · metric: cweightkgs 70, cheightmeter 178, cneckmeter 50, cwaistmeter 96, chipmeter 92.
Results: Navy-method body fat % (0.1), a category, fat mass, lean mass, an "ideal" % by age, BMI-method %.

## MATCH
Was PARTIAL MED-ANTHRO-008 (Deurenberg BMI method), MED-ANTHRO-011 (RFM). Tape method: NONE → built.

## RESEARCH
DoD circumference equations (% fat directly, inches): men 86.010·log10(abdomen − neck) − 70.041·log10(height) + 36.76;
women 163.205·log10(waist + hip − neck) − 97.684·log10(height) − 78.387.
Read in: Foulis SA, et al. (co-author Friedl KE). Front Physiol 2023;14:1183836, doi:10.3389/fphys.2023.1183836, Methods (cm ÷ 2.54 form).
Origin: Hodgdon JA, Beckett MB, NHRC Reports 84-11 and 84-29 (1984); Hodgdon JA, Friedl KE, NHRC TD 99-2B (1999) — DTIC PDF returned
403 and a direct download was refused in this session, so not read first hand.

## BUILD
MED-ANTHRO-016 BodyFatPercentNavyMale, MED-ANTHRO-017 BodyFatPercentNavyFemale — med/anthropometry/body_fat_percent_navy.go;
tests body_fat_percent_navy_test.go (hand examples 17.51309, 28.55626; unit-log identities; errors; 0 allocs; benchmarks).

## COMPARE (scratch: checks/scratch/body-fat-calculator.go)
| inputs (in) | ours | theirs | diff | verdict |
|---|---|---|---|---|
| m, waist 37.5, neck 19.5, h 70.5 (defaults) | 15.2770 | 15.3 | 0 at 0.1 | MATCH |
| m, waist 34, neck 15, h 70 | 17.5131 | 17.4 | 0.11 | THEIRS DIFFERS |
| f, waist 30, hip 38, neck 13, h 65 | 28.5563 | 28.3 | 0.26 | THEIRS DIFFERS |

## DISCREPANCY: THEIRS DIFFERS
Their numbers are reproduced exactly by the body-DENSITY form (metric, cm) followed by Siri's 495/D − 450:
men D = 1.0324 − 0.19077·log10(waist − neck) + 0.15456·log10(height) → 17.436; women D = 1.29579 − 0.35004·log10(waist + hip − neck)
+ 0.22100·log10(height) → 28.258 (hand computation, this file). That is the Hodgdon–Beckett 1984 density formulation (not read first hand;
coefficients here are only a reconstruction of their output, not shipped). Ours is the DoD direct-%fat form as printed in Foulis et al. 2023.
The two forms differ by a few tenths of a percent. Ours stays. If Devin wants the density form too, it needs the 1984 NHRC reports.
Submissions: 1 fields + 3 runs.

## WITHDRAWN (same session, 2026-11-13)
The two-source bar of D-075/D-076 (two independent authoritative secondaries that visibly cite the primary, not by each other's
authors nor the primary's) is not met, so the code was moved to _delete/med/anthropometry/ (files kept for restore when the PDF arrives)
and the row is NEEDS SOURCE. The comparison above stays as the record of what the withdrawn code produced.
