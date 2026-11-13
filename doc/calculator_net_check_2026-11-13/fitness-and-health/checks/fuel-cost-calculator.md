# fuel-cost-calculator (HANDOFF from A-OTHER; notes ../other/checks/fuel-cost-calculator.md)

## FIELDS (form calform, 2026-11-13)
tripdistance 500 + unit miles|kilometers · fuelefficiency 25 + unit mpg|lp100k|kpl|lpm · gasprice 3 + unit gallon|liter.
Results: fuel needed, cost (cents), and a what-if table at other efficiencies.

## BUILD
NEW PACKAGE life/automotive (taxonomy SHORT AUTO): LIF-AUTO-001 FuelEconomy, LIF-AUTO-002 FuelConsumptionPer100, LIF-AUTO-003 FuelCostCents
— life/automotive/doc.go, fuel_economy.go, fuel_economy_test.go. Source 40 CFR 600.002 "Fuel economy" (read 2026-11-13, law.cornell.edu).
L/100 km is given to FuelCostCents as km/L = 100 / (L/100 km) in the scratch (caller arithmetic; no unit/volume package yet).

## COMPARE (scratch: checks/scratch/fuel-cost-calculator.go)
| inputs | ours | theirs | verdict |
|---|---|---|---|
| 1234 mi, 31 mpg, $3.59/gal | 14290.516 → 14291 c | $142.91 | MATCH |
| 380 km, 7.5 L/100 km, $1.89/L | 28.5 L × 189 = 5386.5 → 5386 c (half-even) | 28.5 L, $53.87 | MATCH once the rounding mode is aligned: they round the tie half-up, ours half-even (numeric.md default) |
| 500 mi, 25 mpg, $3/gal (their default; not submitted) | 6000 c | — | — |
Submissions: 1 fields + 2 runs.
