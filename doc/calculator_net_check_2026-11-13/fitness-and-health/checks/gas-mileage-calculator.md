# gas-mileage-calculator (HANDOFF from A-OTHER; notes ../other/checks/gas-mileage-calculator.md)

## FIELDS (form calcf, 2026-11-13)
US: uscodreading 12360 · uspodreading 12000 · usgasputin 15 · usgasprice 3. Metric: mucodreading 18900 · mupodreading 18300 · mugasputin 50 ·
mugasprice 1. Result: mpg plus L/100 km and km/L.

## MATCH
LIF-AUTO-001 FuelEconomy, LIF-AUTO-002 FuelConsumptionPer100 (built in row 36). Distance = current − previous odometer (caller).

## COMPARE (scratch: checks/scratch/gas-mileage-calculator.go)
| inputs | ours | theirs | verdict |
|---|---|---|---|
| 12360 − 12000 mi, 15 gal | 24 mpg | 24 mpg | MATCH |
| 18900 − 18300 km, 50 L (computed, not submitted) | 12 km/L, 8.333 L/100 km | — | — |
Their 9.8 L/100 km / 10.2 km/L beside 24 mpg are unit conversions (unit/volume does not exist yet; not compared).
Submissions: 1 fields + 1 run.
