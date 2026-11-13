# fitness-and-health: progress (A-FIT)

One row per calculator. status: OPEN until closed as MATCH / FIXED (what, source) / BUILT <IDs> / SKIPPED (reason) / NEEDS SOURCE / HANDOFF <agent> / BLOCKED / UNRESOLVED. A closed row is skipped on relaunch.
first-pass match is from doc/calculator_net_map_2026-11-13 and must be re-checked by grep.

| # | name | slug | url | first-pass match | status | IDs | check file |
|---|---|---|---|---|---|---|---|
| 1 | BMI Calculator | bmi-calculator | https://www.calculator.net/bmi-calculator.html | HAVE MED-ANTHRO-001, MED-ANTHRO-002, MED-ANTHRO-003 | MATCH (labels: theirs WHO, ours CDC) | MED-ANTHRO-001, MED-ANTHRO-002, MED-ANTHRO-003 | checks/bmi-calculator.md |
| 2 | Calorie Calculator | calorie-calculator | https://www.calculator.net/calorie-calculator.html | PARTIAL LIF-NUTR-001, MED-ANTHRO-010 (activity multipliers and Katch-McArdle not built) | MATCH (Mifflin simplified; activity table, Revised HB, Katch-McArdle NEEDS SOURCE) | LIF-NUTR-002 | checks/calorie-calculator.md |
| 3 | Body Fat Calculator | body-fat-calculator | https://www.calculator.net/body-fat-calculator.html | PARTIAL MED-ANTHRO-008, MED-ANTHRO-011 (US Navy tape method not built) | NEEDS SOURCE (built, then withdrawn per D-075/D-076: both readable sources co-authored by Friedl; needs NHRC 84-11/84-29 or TD 99-2B PDF) | (016, 017 withdrawn) | checks/body-fat-calculator.md |
| 4 | BMR Calculator | bmr-calculator | https://www.calculator.net/bmr-calculator.html | PARTIAL LIF-NUTR-001, MED-ANTHRO-010 (Katch-McArdle not built) | MATCH (Mifflin simplified; Revised HB, Katch-McArdle NEEDS SOURCE) | LIF-NUTR-002 | checks/bmr-calculator.md |
| 5 | Ideal Weight Calculator | ideal-weight-calculator | https://www.calculator.net/ideal-weight-calculator.html | PARTIAL MED-ANTHRO-006 (Devine, Robinson, Miller, Hamwi not built) | NEEDS SOURCE (built, then withdrawn per D-075/D-076: one secondary only, Peterson AJCN 2016 Table 3; needs Devine 1974, Robinson 1983, Miller 1983, Hamwi 1964) | (018-021 withdrawn) | checks/ideal-weight-calculator.md |
| 6 | Pace Calculator | pace-calculator | https://www.calculator.net/pace-calculator.html | HAVE LIF-FIT-009, LIF-FIT-010 | MATCH | LIF-FIT-009 | checks/pace-calculator.md |
| 7 | Army Body Fat Calculator | army-body-fat-calculator | https://www.calculator.net/army-body-fat-calculator.html | NONE (US Army tape formula) | NEEDS SOURCE (Army Directive 2023-11 / AR 600-9 one-site equation text) | | checks/army-body-fat-calculator.md |
| 8 | Lean Body Mass Calculator | lean-body-mass-calculator | https://www.calculator.net/lean-body-mass-calculator.html | NONE (Boer, James, Hume) | NEEDS SOURCE (Hume 1966 full text, Boer 1984, James 1976 not read) | | checks/lean-body-mass-calculator.md |
| 9 | Healthy Weight Calculator | healthy-weight-calculator | https://www.calculator.net/healthy-weight-calculator.html | PARTIAL MED-ANTHRO-001, MED-ANTHRO-003 (BMI-range inversion not built) | BUILT MED-ANTHRO-022, MED-ANTHRO-023 | MED-ANTHRO-022, MED-ANTHRO-023 | checks/healthy-weight-calculator.md |
| 10 | Calories Burned Calculator | calories-burned-calculator | https://www.calculator.net/calories-burned-calculator.html | NONE (MET compendium) | BUILT LIF-FIT-015 (compare not possible: their MET hidden, mass-dependent) | LIF-FIT-015 | checks/calories-burned-calculator.md |
| 11 | One Rep Max Calculator | one-rep-max-calculator | https://www.calculator.net/one-rep-max-calculator.html | HAVE LIF-FIT-005, LIF-FIT-012, LIF-FIT-013, LIF-FIT-014 | MATCH (Brzycki; they truncate to 0.1; Epley, Lombardi NEEDS SOURCE) | LIF-FIT-012 | checks/one-rep-max-calculator.md |
| 12 | Target Heart Rate Calculator | target-heart-rate-calculator | https://www.calculator.net/target-heart-rate-calculator.html | PARTIAL MED-CARDIO-001, MED-CARDIO-029, MED-CARDIO-030 (Karvonen zones not built) | MATCH (max HR Tanaka, Nes; Karvonen and 220 − age NEEDS SOURCE) | MED-CARDIO-001, MED-CARDIO-030 | checks/target-heart-rate-calculator.md |
| 13 | Pregnancy Calculator | pregnancy-calculator | https://www.calculator.net/pregnancy-calculator.html | PARTIAL MED-OB-001 (Naegele due date not built) | BUILT MED-OB-003, MED-OB-004 (MATCH, 4 cases incl. non-default dates in session 2) | MED-OB-003, MED-OB-004 | checks/pregnancy-calculator.md |
| 14 | Pregnancy Weight Gain Calculator | pregnancy-weight-gain-calculator | https://www.calculator.net/pregnancy-weight-gain-calculator.html | NONE (IOM 2009 ranges) | BUILT MED-OB-005, MED-OB-006 (MATCH) | MED-OB-005, MED-OB-006 | checks/pregnancy-weight-gain-calculator.md |
| 15 | Pregnancy Conception Calculator | pregnancy-conception-calculator | https://www.calculator.net/pregnancy-conception-calculator.html | NONE | NEEDS SOURCE (conception = EDD − 266 / ovulation day-14 rule not citable; Wilcox BMJ 2000) | | checks/pregnancy-conception-calculator.md |
| 16 | Due Date Calculator | due-date-calculator | https://www.calculator.net/due-date-calculator.html | NONE | MATCH at 28-day cycle; THEIRS DIFFERS for other cycle lengths (they add cycle − 28 days; ACOG CO 700 states none) | MED-OB-003, MED-OB-004 | checks/due-date-calculator.md |
| 17 | Ovulation Calculator | ovulation-calculator | https://www.calculator.net/ovulation-calculator.html | NONE | NEEDS SOURCE (luteal 14-day rule contradicted by Wilcox BMJ 2000; Devin decision) | | checks/ovulation-calculator.md |
| 18 | Conception Calculator | conception-calculator | https://www.calculator.net/conception-calculator.html | NONE | NEEDS SOURCE (same as row 17) | | checks/conception-calculator.md |
| 19 | Period Calculator | period-calculator | https://www.calculator.net/period-calculator.html | NONE | NEEDS SOURCE (future periods and ovulation from fixed cycle length; same as row 17) | | checks/period-calculator.md |
| 20 | Macro Calculator | macro-calculator | https://www.calculator.net/macro-calculator.html | PARTIAL LIF-NUTR-003 | MATCH energy; ranges THEIRS DIFFERS (not AMDR) | LIF-NUTR-002, LIF-NUTR-003 | checks/macro-calculator.md |
| 21 | Carbohydrate Calculator | carbohydrate-calculator | https://www.calculator.net/carbohydrate-calculator.html | PARTIAL LIF-NUTR-003 | MATCH energy; grams THEIRS DIFFERS (3.75 vs 4 kcal/g) | LIF-NUTR-002, LIF-NUTR-003 | checks/carbohydrate-calculator.md |
| 22 | Protein Calculator | protein-calculator | https://www.calculator.net/protein-calculator.html | HAVE LIF-NUTR-007, LIF-NUTR-005 | THEIRS DIFFERS (same 10–35 % rule; their energy basis ≈ 2310 kcal hidden) | LIF-NUTR-003 | checks/protein-calculator.md |
| 23 | Fat Intake Calculator | fat-intake-calculator | https://www.calculator.net/fat-intake-calculator.html | PARTIAL LIF-NUTR-003 | MATCH energy; grams THEIRS DIFFERS (≈ 8.8 vs 9 kcal/g) | LIF-NUTR-002, LIF-NUTR-003 | checks/fat-intake-calculator.md |
| 24 | TDEE Calculator | tdee-calculator | https://www.calculator.net/tdee-calculator.html | PARTIAL LIF-NUTR-001 (activity multiplier not built) | MATCH (REE and BMI; activity multiplier NEEDS SOURCE) | LIF-NUTR-002, MED-ANTHRO-002 | checks/tdee-calculator.md |
| 25 | GFR Calculator | gfr-calculator | https://www.calculator.net/gfr-calculator.html | HAVE MED-RENAL-001, MED-RENAL-022, MED-RENAL-007 | MATCH MDRD, Schwartz; CKD-EPI THEIRS DIFFERS (2009 vs 2021) | MED-RENAL-001, MED-RENAL-007, MED-RENAL-022 | checks/gfr-calculator.md |
| 26 | Body Type Calculator | body-type-calculator | https://www.calculator.net/body-type-calculator.html | OUT OF SCOPE (shape heuristic, no primary formula) | SKIPPED (shape heuristic, no primary formula) | | checks/body-type-calculator.md |
| 27 | Body Surface Area Calculator | body-surface-area-calculator | https://www.calculator.net/body-surface-area-calculator.html | PARTIAL MED-ANTHRO-004 (Mosteller, DuBois not built) | MATCH (Haycock; Du Bois, Mosteller and 6 others NEEDS SOURCE) | MED-ANTHRO-004 | checks/body-surface-area-calculator.md |
| 28 | BAC Calculator | bac-calculator | https://www.calculator.net/bac-calculator.html | NONE (Widmark) | NEEDS SOURCE (Widmark 1932 not available; NHTSA DOT HS 805-563 to read) | | checks/bac-calculator.md |
| 29 | Tile Calculator | tile-calculator | https://www.calculator.net/tile-calculator.html | HANDOFF from A-OTHER (life/home) | MATCH (NUM-GEOM-022 + caller arithmetic; waste % a convention) | NUM-GEOM-022 | checks/tile-calculator.md |
| 30 | Mulch Calculator | mulch-calculator | https://www.calculator.net/mulch-calculator.html | HANDOFF from A-OTHER (life/garden) | MATCH | NUM-GEOM-048 | checks/mulch-calculator.md |
| 31 | Gravel Calculator | gravel-calculator | https://www.calculator.net/gravel-calculator.html | HANDOFF from A-OTHER (life/garden) | MATCH (density a caller input) | NUM-GEOM-035, NUM-GEOM-048 | checks/gravel-calculator.md |
| 32 | Height Calculator | height-calculator | https://www.calculator.net/height-calculator.html | HANDOFF from A-OTHER (med/pediatrics) | NEEDS SOURCE (Khamis–Roche tables; Tanner 1970 scan) | | checks/height-calculator.md |
| 33 | Tip Calculator | tip-calculator | https://www.calculator.net/tip-calculator.html | HANDOFF from A-OTHER (life/shopping) | BUILT LIF-SHOP-001 (split remainder policy for Devin) | LIF-SHOP-001 | checks/tip-calculator.md |
| 34 | Golf Handicap Calculator | golf-handicap-calculator | https://www.calculator.net/golf-handicap-calculator.html | HANDOFF from A-OTHER (life/sports) | BUILT LIF-SPORT-001 (Rule 5.1 differential and 6.1 course handicap NEEDS SOURCE) | LIF-SPORT-001 | checks/golf-handicap-calculator.md |
| 35 | Sleep Calculator | sleep-calculator | https://www.calculator.net/sleep-calculator.html | HANDOFF from A-OTHER (med/pediatrics; HAVE MED-PEDS-002) | NEEDS SOURCE (90-min cycle timing; hours by age HAVE MED-PEDS-002) | MED-PEDS-002 | checks/sleep-calculator.md |
| 36 | Fuel Cost Calculator | fuel-cost-calculator | https://www.calculator.net/fuel-cost-calculator.html | HANDOFF from A-OTHER (life/automotive) | BUILT LIF-AUTO-001, LIF-AUTO-002, LIF-AUTO-003 (MATCH; tie rounding half-even vs their half-up) | LIF-AUTO-001, LIF-AUTO-002, LIF-AUTO-003 | checks/fuel-cost-calculator.md |
| 37 | Gas Mileage Calculator | gas-mileage-calculator | https://www.calculator.net/gas-mileage-calculator.html | HANDOFF from A-OTHER (life/automotive) | MATCH (built in row 36) | LIF-AUTO-001, LIF-AUTO-002 | checks/gas-mileage-calculator.md |
| 38 | Mileage Calculator | mileage-calculator | https://www.calculator.net/mileage-calculator.html | HANDOFF from A-OTHER (life/automotive) | SKIPPED (route distance from map data; IRS rate is a yearly table) | | checks/mileage-calculator.md |
| 39 | Tire Size Calculator | tire-size-calculator | https://www.calculator.net/tire-size-calculator.html | HANDOFF from A-OTHER (life/automotive) | NEEDS SOURCE (ISO 4000-1 / ETRTO / T&RA not read; CFR lacks aspect ratio) | | checks/tire-size-calculator.md |

## Site log

Page loads and submissions made (count per calculator), any 403/429/block.

- bmi-calculator: 1 fields load, 4 runs, no block.
- calorie-calculator: 1 fields load, 3 runs, no block.
- body-fat-calculator: 1 fields load, 3 runs, no block.
- bmr-calculator: 1 fields load, 2 runs, no block.
- ideal-weight-calculator: 1 fields load, 2 runs, no block.
- pace-calculator: 1 fields load, 3 runs, no block.
- army-body-fat-calculator: 1 fields load, 0 runs.
- lean-body-mass-calculator: 0 loads.
- healthy-weight-calculator: 1 fields load, 2 runs, no block.
- calories-burned-calculator: 1 fields load, 2 runs, no block.
- one-rep-max-calculator: 1 fields load, 2 runs, no block.
- target-heart-rate-calculator: 1 fields load, 2 runs, no block.
- pregnancy-calculator: 1 fields load, 1 run, no block. DRIVER: date-picker selects (…_Month_ID, …_Day_ID, …_Year_ID) not filled.
- pregnancy-weight-gain-calculator: 1 fields load, 2 runs, no block.
- due-date-calculator: 1 fields load, 1 run, no block.
- rows 15, 17, 18, 19: 0 loads.
- macro-calculator: 1 fields load, 1 run, no block.
- carbohydrate-calculator: 1 fields load, 1 run, no block.
- protein-calculator: 1 fields load, 1 run, no block.
- fat-intake-calculator: 1 fields load, 1 run, no block.
- tdee-calculator: 1 fields load, 1 run, no block.
- gfr-calculator: 1 fields load, 3 runs, no block.
- body-surface-area-calculator: 1 fields load, 2 runs, no block.
- bac-calculator, body-type-calculator: 0 loads.

### Session 2 (start 09:36)
- tile-calculator: 1 fields load, 1 run, no block.
- mulch-calculator, gravel-calculator: 1 fields + 1 run each, no block.
- height-calculator: 1 fields load, 0 runs.
- golf-handicap-calculator: 2 fields loads, 2 runs, no block. DRIVER: differentials behind a link not captured.
- tip-calculator: 1 fields load, 2 runs, no block. sleep-calculator: 0 loads.
- fuel-cost-calculator: 1 fields + 2 runs; gas-mileage-calculator: 1 fields + 1 run; mileage-calculator: 1 fields. No block.
- tire-size-calculator: 0 loads.
- session 2: pregnancy-calculator 1 fields + 3 runs (4 runs total); due-date-calculator 2 runs (3 total).
- session 2: pregnancy-conception-calculator 1 fields + 2 runs (nothing of ours to compare).
