# calories-burned-calculator

## FIELDS (2026-11-13)
calform: activity (6 groups) · activity2 radio, 106 named activities (default "Walking: slow") · chour · cmin 45 · cweight 150 ·
cweightunit p|k (80-350 lb or 35-160 kg). calform2: c2activity walk|run|cycle · c2speed 3 + c2speedunit · c2distance 5 + unit · c2weight 150.
Result: one kcal figure (integer). Their MET value per activity is NOT shown.

## MATCH
NONE → built the stable formula. The activity→MET table is the Compendium's data table (not held here; caller supplies MET).

## BUILD
LIF-FIT-015 EnergyExpenditureMET — life/fitness/energy_expenditure_met.go; tests energy_expenditure_met_test.go.
Source: Herrmann SD, et al. 2024 Adult Compendium of Physical Activities. J Sport Health Sci 2024;13(1):6-12,
doi:10.1016/j.jshs.2023.10.010, PMC10818145, §2.5: 1 MET = "3.5 mL/kg/min or approximately 1.0 kcal/kg/h". Primary (the Compendium itself).

## COMPARE (scratch: checks/scratch/calories-burned-calculator.go)
| inputs | ours (MET 2.0 assumed) | theirs ("Walking: slow") | implied kcal/kg/h | verdict |
|---|---|---|---|---|
| 150 lb (68.039 kg), 45 min | 102.058 (107.161 at ×1.05) | 108 | 2.117 | not comparable |
| 100 kg, 60 min | 200.000 (210.000 at ×1.05) | 206 | 2.060 | not comparable |
Their result per kg·h changes with body mass (2.117 vs 2.060), so they do not use plain MET × kg × h with a fixed MET; a mass-dependent
resting rate (a "corrected MET") is the likely convention. Their MET for "Walking: slow" is not displayed, so no common input exists.
Verdict: no comparison possible; ours stays on the Compendium definition. Submissions: 1 fields + 2 runs.
