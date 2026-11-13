# A-FIT notes for the lead

## INDEX rows to add (regenerate from code)
- MED-ANTHRO-022 WeightAtBMI, MED-ANTHRO-023 WeightAtBMIImperial — med/anthropometry/weight_at_bmi.go
- LIF-FIT-015 EnergyExpenditureMET — life/fitness/energy_expenditure_met.go
- MED-OB-003 DueDateFromLMP, MED-OB-004 GestationalAgeFromLMP — med/obstetric/due_date.go (also updated the package sentence in med/obstetric/doc.go header, lead to check)
- MED-OB-005 GestationalWeightGainIOM, MED-OB-006 GestationalWeightGainTwinsIOM — med/obstetric/gestational_weight_gain_iom.go (new exported types GWGCategory, GWGRecommendation)
- LIF-SPORT-001 HandicapIndexWHS — NEW PACKAGE life/sports (SHORT SPORT from taxonomy; doc.go header for the lead to check)
- LIF-SHOP-001 TipCents — NEW PACKAGE life/shopping (SHORT SHOP; doc.go header for the lead)
- LIF-AUTO-001 FuelEconomy, LIF-AUTO-002 FuelConsumptionPer100, LIF-AUTO-003 FuelCostCents — NEW PACKAGE life/automotive (SHORT AUTO; doc.go header for the lead)

## Source notes
- MED-ANTHRO-016/017: formula read in Foulis et al. 2023 (Front Physiol, co-author Friedl, one of the DoD equation's developers).
  Original NHRC reports (Hodgdon & Beckett 1984; Hodgdon & Friedl TD 99-2B 1999) not read: DTIC PDF gave 403 to WebFetch and a
  curl download was refused by the permission classifier. If Devin can supply the DTIC PDF, add it as the second Source line.
  doc/sources/ entry not written (shared path?) — lead to decide.

## Not built, needs a source (Devin)
- Activity multipliers 1.2/1.375/1.465/1.55/1.725/1.9 (calorie, BMR, TDEE pages): no primary source found.
- Revised Harris-Benedict (Roza & Shizgal 1984, AJCN 40:168) and Katch-McArdle (370 + 21.6·LBM): primaries not seen.

- Army one-site body fat (row 7): need Army Directive 2023-11 or AR 600-9 (2023) text; Taylor et al. Obes Sci Pract 2024;10:e747 Table 2 coefficients differ from the commonly quoted policy ones. See checks/army-body-fat-calculator.md.

- Lean body mass (row 8): need first-hand Hume 1966 (J Clin Pathol 19:389, PMC473290 scan), Boer 1984 (Am J Physiol 247:F632), James 1976 (HMSO).

- 1 RM Epley (1985) and Lombardi (1989): primaries not read.

- Karvonen (1957) and Fox/Haskell 220 − age (1971): primaries not available; ACSM 2011 position stand (Garber, MSSE 43:1334) would be a candidate second source.

- Ovulation / conception / period rows (15, 17, 18, 19): every one rests on "ovulation = next menses − 14 days". Wilcox et al. BMJ 2000;321:1259 (PMC27529) names that as a guideline assumption its data contradict, and gives the fertile window as ovulation day − 5 .. ovulation day. Decision for Devin: (a) build FertileWindow(ovulationDate) only, from Wilcox, with ovulation date supplied by the caller (recommended: sourced, no luteal assumption); (b) also build the 14-day rule labelled as a convention, citing a guideline that states it (none found yet); lost for (b): no citable primary.
- ACOG CO 700 ART / conception day counts (266, 261, 263): need the full CO 700 text (Obstet Gynecol 2017;129:e150).

- GFR: their page uses CKD-EPI 2009 (race term) and the Mayo quadratic; ours has 2021 only. Decision for Devin whether to add 2009 (superseded) as a separate function.

- BAC (row 28): Widmark 1932 not online; read Burns & Moskowitz, NHTSA DOT HS 805-563 (1980) first hand next time.
- BSA (row 27): Du Bois 1916 / Mosteller 1987 and six others not built — next session.

- Tip split among N people: the remainder policy (who pays the extra cent) is a numeric-policy decision for Devin; not built.
- WHS: Rule 5.1 Score Differential and Rule 6.1 Course Handicap need a first-hand text (usga.org 403; randa.org page has no formula text). Handicap Index tie rounding: Rules say "nearest tenth" with no tie rule; I used half-even (house default) — confirm.

- Tire size (row 39): needs ISO 4000-1 or the ETRTO / T&RA manual; 49 CFR 571.109 / 571.139 do not define aspect ratio.
- Height (row 32): needs Khamis & Roche 1994 table + erratum, and the Tanner 1970 text for the ± 6.5 cm mid-parental adjustment.

## Extras seen on their pages, not built
- BMI page: BMI Prime (BMI/25) and Ponderal Index (kg/m³).

## Driver gaps
- DRIVER (session 1): date pickers not filled — FIXED by the lead (hidden MM/DD/YYYY field). Note: the fields driver does not list the hidden field; its name was the widget prefix (cLastPeriodDay, dueDateDay).
- DRIVER: golf handicap differentials (behind "Click here to view the differentials") not captured.

## Notes
- checks/scratch/body-fat-calculator.go and ideal-weight-calculator.go call the withdrawn functions; they are `//go:build ignore` and will only run again if the files are restored.
- doc/sources/ entries for CDC (WeightAtBMI), ACOG CO 700, IOM 2009, Herrmann 2024 Compendium not written (shared path?) — lead to add.

## _delete moves (path | replaced by | date)
- med/anthropometry/body_fat_percent_navy.go | none (withdrawn, D-075/D-076 two-source bar) | 2026-11-13
- med/anthropometry/body_fat_percent_navy_test.go | none (withdrawn) | 2026-11-13
- med/anthropometry/ideal_body_weight.go | none (withdrawn, D-075/D-076) | 2026-11-13
- med/anthropometry/ideal_body_weight_test.go | none (withdrawn) | 2026-11-13
(_delete/med/anthropometry/body_mass_index.go was already there before this session; not mine.)

## ID decision for the lead
MED-ANTHRO-016..021 were written into the withdrawn code and briefly into med/anthropometry/doc.go, never into INDEX.md.
Recommendation: treat them as given (REMOVED rows, "withdrawn before release: two-source bar, NEEDS PDF"), so restoring the files
later keeps their IDs; WeightAtBMI keeps 022/023. Alternative: renumber WeightAtBMI to 016/017 (lost: the IDs-are-never-reused rule).
- 016 BodyFatPercentNavyMale, 017 BodyFatPercentNavyFemale: read in Foulis et al. Front Physiol 2023;14:1183836 and Potter et al.
  Front Physiol 2022;13:868627 (both co-authored by Friedl). Needs Hodgdon & Beckett NHRC 84-11 / 84-29 (DTIC ADA143890, ADA146456)
  or Hodgdon & Friedl TD 99-2B (DTIC ADA370158): DTIC returns 403 to fetch and a curl download was refused by the permission classifier.
- 018-021 IdealBodyWeightDevine/Robinson/Miller/HamwiLb: read in Peterson et al. AJCN 2016;103:1197 Table 3 only (Green & Duffull,
  Br J Clin Pharmacol 2004, PMC1884581, gives Devine only, in a rounded cm form). Needs the four originals or a second table.
