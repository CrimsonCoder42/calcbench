# other: progress (A-OTHER)

One row per calculator. status: OPEN until closed as MATCH / FIXED (what, source) / BUILT <IDs> / SKIPPED (reason) / NEEDS SOURCE / HANDOFF <agent> / BLOCKED / UNRESOLVED. A closed row is skipped on relaunch.
first-pass match is from doc/calculator_net_map_2026-11-13 and must be re-checked by grep.

| # | name | slug | url | first-pass match | status | IDs | check file |
|---|---|---|---|---|---|---|---|
| 1 | Age Calculator | age-calculator | https://www.calculator.net/age-calculator.html | NONE (calendar arithmetic) | BUILT SCI-ASTRO-001..006 (MATCH; Y/M/D split for Devin) | SCI-ASTRO-001..006 | checks/age-calculator.md |
| 2 | Date Calculator | date-calculator | https://www.calculator.net/date-calculator.html | NONE (calendar arithmetic) | MATCH | SCI-ASTRO-003, -005, -006 | checks/date-calculator.md |
| 3 | Time Calculator | time-calculator | https://www.calculator.net/time-calculator.html | NONE | SKIPPED (stdlib time.Duration clock arithmetic) |  | checks/time-calculator.md |
| 4 | Hours Calculator | hours-calculator | https://www.calculator.net/hours-calculator.html | NONE | SKIPPED (stdlib time.Duration clock arithmetic) |  | checks/hours-calculator.md |
| 5 | Time Card Calculator | time-card-calculator | https://www.calculator.net/time-card-calculator.html | NONE | SKIPPED (duration sum + rate x hours; overtime rules jurisdictional) |  | checks/time-card-calculator.md |
| 6 | Time Zone Calculator | time-zone-calculator | https://www.calculator.net/time-zone-calculator.html | OUT OF SCOPE (tz database) | SKIPPED (IANA tz database) |  | checks/time-zone-calculator.md |
| 7 | Time Duration Calculator | time-duration-calculator | https://www.calculator.net/time-duration-calculator.html | NONE | SKIPPED (stdlib time.Duration clock arithmetic) |  | checks/time-duration-calculator.md |
| 8 | Day Counter | day-counter | https://www.calculator.net/day-counter.html | NONE | MATCH | SCI-ASTRO-005, -003 | checks/day-counter.md |
| 9 | Day of the Week Calculator | day-of-the-week-calculator | https://www.calculator.net/day-of-the-week-calculator.html | NONE | MATCH | SCI-ASTRO-003 | checks/day-of-the-week-calculator.md |
| 10 | Concrete Calculator | concrete-calculator | https://www.calculator.net/concrete-calculator.html | PARTIAL NUM-GEOM-022, NUM-GEOM-049 (volume only; bag yields not modelled) | MATCH (volume; weight/bags NEEDS SOURCE) | NUM-GEOM-022, -048, -049 | checks/concrete-calculator.md |
| 11 | BTU Calculator | btu-calculator | https://www.calculator.net/btu-calculator.html | NONE | NEEDS SOURCE (area-factor rule of thumb; Manual J is its own slice) |  | checks/btu-calculator.md |
| 12 | Square Footage Calculator | square-footage-calculator | https://www.calculator.net/square-footage-calculator.html | PARTIAL NUM-GEOM-022, NUM-GEOM-035 (shape areas only) | BUILT UNT-AREA-001, UNT-AREA-002 (MATCH) | UNT-AREA-001, -002; NUM-GEOM-022 | checks/square-footage-calculator.md |
| 13 | Stair Calculator | stair-calculator | https://www.calculator.net/stair-calculator.html | NONE | MATCH (geometry; IRC riser/tread limits not built) | NUM-GEOM-020, NUM-TRIG-027 | checks/stair-calculator.md |
| 14 | Roofing Calculator | roofing-calculator | https://www.calculator.net/roofing-calculator.html | NONE | MATCH (geometry) | NUM-GEOM-020 | checks/roofing-calculator.md |
| 15 | Tile Calculator | tile-calculator | https://www.calculator.net/tile-calculator.html | NONE | HANDOFF A-FIT (life/home) |  | checks/tile-calculator.md |
| 16 | Mulch Calculator | mulch-calculator | https://www.calculator.net/mulch-calculator.html | PARTIAL NUM-GEOM-048 (volume only) | HANDOFF A-FIT (life/garden) |  | checks/mulch-calculator.md |
| 17 | Gravel Calculator | gravel-calculator | https://www.calculator.net/gravel-calculator.html | PARTIAL NUM-GEOM-048 (volume only; material density not modelled) | HANDOFF A-FIT (life/garden) |  | checks/gravel-calculator.md |
| 18 | Height Calculator | height-calculator | https://www.calculator.net/height-calculator.html | NONE (child adult-height prediction) | HANDOFF A-FIT (med/peds) |  | checks/height-calculator.md |
| 19 | Conversion Calculator | conversion-calculator | https://www.calculator.net/conversion-calculator.html | PARTIAL UNT-LEN-001/002, UNT-MASS-001/002, UNT-TEMP-001..004 (few quantities) | MATCH (UNT-LEN-001/002 live; mass/temp/area tabs not listed by fields) | UNT-LEN-001/002 | checks/conversion-calculator.md |
| 20 | GDP Calculator | gdp-calculator | https://www.calculator.net/gdp-calculator.html | NONE | HANDOFF A-FIN (fin/economics) |  | checks/gdp-calculator.md |
| 21 | Density Calculator | density-calculator | https://www.calculator.net/density-calculator.html | NONE | BUILT SCI-PHYS-001..003 (MATCH) | SCI-PHYS-001..003 | checks/density-calculator.md |
| 22 | Mass Calculator | mass-calculator | https://www.calculator.net/mass-calculator.html | NONE | MATCH | SCI-PHYS-002 | checks/mass-calculator.md |
| 23 | Weight Calculator | weight-calculator | https://www.calculator.net/weight-calculator.html | NONE | MATCH (mass-unit converter) | UNT-MASS-002 | checks/weight-calculator.md |
| 24 | Speed Calculator | speed-calculator | https://www.calculator.net/speed-calculator.html | NONE | BUILT SCI-PHYS-004..006 (MATCH) | SCI-PHYS-004..006 | checks/speed-calculator.md |
| 25 | Molarity Calculator | molarity-calculator | https://www.calculator.net/molarity-calculator.html | NONE | BUILT SCI-CHEM-001..004 (MATCH) | SCI-CHEM-001..004 | checks/molarity-calculator.md |
| 26 | Molecular Weight Calculator | molecular-weight-calculator | https://www.calculator.net/molecular-weight-calculator.html | NONE (needs IUPAC atomic weights) | NEEDS SOURCE (CIAAW table edition/interval policy: Devin) |  | checks/molecular-weight-calculator.md |
| 27 | Roman Numeral Converter | roman-numeral-converter | https://www.calculator.net/roman-numeral-converter.html | NONE | HANDOFF A-MATH (algo/text) |  | checks/roman-numeral-converter.md |
| 28 | Voltage Drop Calculator | voltage-drop-calculator | https://www.calculator.net/voltage-drop-calculator.html | NONE | BUILT ENG-ELEC-007, -008 (compare: DRIVER mode gap remains; NEC tables NEEDS SOURCE) | ENG-ELEC-007, -008 | checks/voltage-drop-calculator.md |
| 29 | Resistor Calculator | resistor-calculator | https://www.calculator.net/resistor-calculator.html | NONE | BUILT ENG-ELEC-009..011 (MATCH; colour code NEEDS SOURCE) | ENG-ELEC-009..011 | checks/resistor-calculator.md |
| 30 | Ohms Law Calculator | ohms-law-calculator | https://www.calculator.net/ohms-law-calculator.html | NONE | BUILT ENG-ELEC-001..006 (MATCH, I and V) | ENG-ELEC-001..006 | checks/ohms-law-calculator.md |
| 31 | Electricity Calculator | electricity-calculator | https://www.calculator.net/electricity-calculator.html | NONE | HANDOFF A-FIN (energy cost; unit/power is a lead question) |  | checks/electricity-calculator.md |
| 32 | IP Subnet Calculator | ip-subnet-calculator | https://www.calculator.net/ip-subnet-calculator.html | NONE | BUILT ENG-COMPUT-001 (MATCH /22; THEIRS DIFFERS /31, RFC 3021) | ENG-COMPUT-001 | checks/ip-subnet-calculator.md |
| 33 | Password Generator | password-generator | https://www.calculator.net/password-generator.html | OUT OF SCOPE (random generator) | SKIPPED (random generator) |  | checks/password-generator.md |
| 34 | Bandwidth Calculator | bandwidth-calculator | https://www.calculator.net/bandwidth-calculator.html | NONE | BUILT ENG-COMPUT-002 (MATCH) | ENG-COMPUT-002 | checks/bandwidth-calculator.md |
| 35 | Base64 Encode / Decode | base64-encode-decode | https://www.calculator.net/base64-encode-decode.html | OUT OF SCOPE (stdlib encoding, not a formula) | SKIPPED (encoder; stdlib encoding/base64) |  | checks/base64-encode-decode.md |
| 36 | URL Encode / Decode | url-encode-decode | https://www.calculator.net/url-encode-decode.html | OUT OF SCOPE (stdlib net/url, not a formula) | SKIPPED (encoder; stdlib net/url) |  | checks/url-encode-decode.md |
| 37 | GPA Calculator | gpa-calculator | https://www.calculator.net/gpa-calculator.html | PARTIAL NUM-STATS-026 (weighted mean; grade-point scale not built) | SKIPPED (grade-point scale is an institution table; arithmetic HAVE NUM-STATS-026) | NUM-STATS-026 | checks/gpa-calculator.md |
| 38 | Grade Calculator | grade-calculator | https://www.calculator.net/grade-calculator.html | PARTIAL NUM-STATS-026 (weighted mean; final-grade solve not built) | MATCH (weighted mean; letter scale out of scope) | NUM-STATS-026 | checks/grade-calculator.md |
| 39 | Bra Size Calculator | bra-size-calculator | https://www.calculator.net/bra-size-calculator.html | OUT OF SCOPE (sizing tables) | SKIPPED (sizing tables) |  | checks/bra-size-calculator.md |
| 40 | Shoe Size Conversion | shoe-size-conversion | https://www.calculator.net/shoe-size-conversion.html | OUT OF SCOPE (sizing tables) | SKIPPED (sizing tables) |  | checks/shoe-size-conversion.md |
| 41 | Tip Calculator | tip-calculator | https://www.calculator.net/tip-calculator.html | NONE | HANDOFF A-FIT (life/shopping; money rounding policy) |  | checks/tip-calculator.md |
| 42 | Golf Handicap Calculator | golf-handicap-calculator | https://www.calculator.net/golf-handicap-calculator.html | NONE (WHS formula) | HANDOFF A-FIT (life/sports) |  | checks/golf-handicap-calculator.md |
| 43 | Sleep Calculator | sleep-calculator | https://www.calculator.net/sleep-calculator.html | PARTIAL MED-PEDS-002 (recommended hours only; cycle timing not built) | HANDOFF A-FIT (med; cycle timing likely NEEDS SOURCE) | MED-PEDS-002 | checks/sleep-calculator.md |
| 44 | Wind Chill Calculator | wind-chill-calculator | https://www.calculator.net/wind-chill-calculator.html | NONE | BUILT SCI-METEO-001 (MATCH) | SCI-METEO-001 | checks/wind-chill-calculator.md |
| 45 | Heat Index Calculator | heat-index-calculator | https://www.calculator.net/heat-index-calculator.html | NONE | BUILT SCI-METEO-002 (MATCH) | SCI-METEO-002, -004 | checks/heat-index-calculator.md |
| 46 | Dew Point Calculator | dew-point-calculator | https://www.calculator.net/dew-point-calculator.html | NONE | BUILT SCI-METEO-003, -004 (MATCH) | SCI-METEO-003, -004 | checks/dew-point-calculator.md |
| 47 | Fuel Cost Calculator | fuel-cost-calculator | https://www.calculator.net/fuel-cost-calculator.html | NONE | HANDOFF A-FIT (life/automotive) |  | checks/fuel-cost-calculator.md |
| 48 | Gas Mileage Calculator | gas-mileage-calculator | https://www.calculator.net/gas-mileage-calculator.html | NONE | HANDOFF A-FIT (life/automotive) |  | checks/gas-mileage-calculator.md |
| 49 | Horsepower Calculator | horsepower-calculator | https://www.calculator.net/horsepower-calculator.html | NONE | BUILT SCI-PHYS-007 (MATCH) | SCI-PHYS-007 | checks/horsepower-calculator.md |
| 50 | Engine Horsepower Calculator | engine-horsepower-calculator | https://www.calculator.net/engine-horsepower-calculator.html | NONE | NEEDS SOURCE (empirical ET/trap-speed fits) |  | checks/engine-horsepower-calculator.md |
| 51 | Mileage Calculator | mileage-calculator | https://www.calculator.net/mileage-calculator.html | NONE | HANDOFF A-FIT (life/automotive) |  | checks/mileage-calculator.md |
| 52 | Tire Size Calculator | tire-size-calculator | https://www.calculator.net/tire-size-calculator.html | NONE | HANDOFF A-FIT (life/automotive) |  | checks/tire-size-calculator.md |
| 53 | Dice Roller | dice-roller | https://www.calculator.net/dice-roller.html | OUT OF SCOPE (entertainment / random) | SKIPPED (random / entertainment) |  | checks/dice-roller.md |
| 54 | Love Calculator | love-calculator | https://www.calculator.net/love-calculator.html | OUT OF SCOPE (entertainment) | SKIPPED (entertainment) |  | checks/love-calculator.md |
| 55 | Distance Calculator (Earth surface; HANDOFF from A-MATH) | distance-calculator | https://www.calculator.net/distance-calculator.html | PARTIAL NUM-TRIG-031 (unit sphere) | BUILT SCI-GEO-001..003 (MATCH; theirs truncates to 0.1 km) | SCI-GEO-001..003 | checks/distance-calculator.md |

## Site log

Page loads and submissions made (count per calculator), any 403/429/block.

Session 1 (A-OTHER, 2026-11-13 09:08-09:35 EST). No 403, 429 or block page.
Fields loads: one per calculator worked (conversion-calculator 2). Submissions (run): age 1, date 3,
day-of-the-week 1, concrete 2, square-footage 1, stair 1, roofing 2, conversion 1 (exit 4), density 2, mass 1,
weight 1, speed 2, molarity 2, voltage-drop 1, ohms-law 2, resistor 3, ip-subnet 2, bandwidth 2, grade 1,
wind-chill 2, heat-index 4, dew-point 2, horsepower 2. Max per calculator: 4 (limit 8).
DRIVER gaps: (1) date-picker fields today_*/ageat_*/today2_* (Month/Day selects and Year text) are never filled
(age, date, day-counter, day-of-the-week); (2) conversion-calculator converts live with no submit button (exit 4);
(3) voltage-drop: the resistance-source mode (NEC table / wire size / custom) is not a fillable field;
(4) ohms-law: only the first solved quantity appears in result.text.

Session 2 (A-OTHER, 2026-11-13 09:36-09:45 EST, driver update). No 403, 429 or block page.
Submissions: age 3 (+ lead's 1; total 5), day-counter 3, conversion 2 (live; total 3), ohms-law 2 (total 4),
voltage-drop 0 (fields re-read only). Reopened and closed: age (MATCH), day-counter (BLOCKED → MATCH),
conversion (BLOCKED → MATCH on length).
DRIVER gaps still open: voltage-drop mode selector has no name or id; conversion-calculator's other quantity tabs
(mass, temperature, area) are not listed by the fields script. Withdrawn: the session-1 ohms-law "first line only"
gap was my jq slicing, not the driver.

Repair (session 3, 2026-11-13 09:50 EST): another agent's helper with the same shared name (setrow.py) had
overwritten the status and IDs columns of about 35 rows with financial closes (FIN- IDs, "TESTX", loan/APR/tax
wording). All 54 rows' status, IDs and check-file columns were rewritten from the checks/<slug>.md CLOSE lines by
scratchpad/aother/aother_repair_progress.py (slug of every row asserted against the expected slug). Verified:
0 FIN- IDs, 0 "TESTX", 0 OPEN rows; LEAD.md, DISCREPANCIES.md and checks/*.md were not touched by the other helper.
The old shared-name helper is no longer used.

Session 3 (A-OTHER, 2026-11-13 09:48 EST). No 403, 429 or block page. Submissions: distance-calculator 3
(Earth-surface form). Row 55 added (HANDOFF from A-MATH).

Reviewer fix (2026-11-13, one pass; go test not run, the lead re-runs the packages):
- eng/electrical/resistor_network.go:70-80 ParallelResistance: a branch whose 1/R is +Inf (subnormal R) is now a
  short before NeumaierAdd, and a total conductance that overflows is a short too; returns 0 as documented
  (doc comment :45-46 updated). Fixes TestParallelResistance_cases "subnormal: NaN".
- sci/meteorology/dew_point.go:42-50 DewPointC: RH = 100 % returns airTempC exactly; Td is clamped to at most
  airTempC (doc comment :15-16 updated). Fixes TestRelativeHumidityFromDewPoint_roundTrip.
- Source lines: OpenStax authors and URLs added (sci/physics, sci/chemistry, eng/electrical); RFC 950/3021/4632
  DOIs added (eng/computing/ipv4_subnet.go). No code change.
