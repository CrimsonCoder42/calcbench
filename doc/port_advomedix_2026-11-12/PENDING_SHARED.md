# Pending rows for shared files (port run 2026-11-12)

The lead owns these files for this run, but other subjects are live in the tree, so the rows are staged here in
final form for the merge pass, as every other subject has done. Paste in this order.

## 1. INDEX.md: 209 rows plus 3 REMOVED rows

In [out/index_rows.md](out/index_rows.md), sorted by ID and generated from the doc comments (Time, Space, Source and
Tags verbatim; Does is the doc.go line). The 44 medicine-run rows are in doc/medicine_2026-11-12/PENDING_SHARED.md §1.

## 2. doc/taxonomy.md: new package rows (D-P4, D-P8)

| MED | `med/gastro` | GASTRO | gastrointestinal and pancreatic scores |
| MED | `med/endocrine` | ENDO | glucose, diabetes and thyroid indices, ADA categories, CGM metrics |
| MED | `med/psychiatry` | PSYCH | psychiatric and substance-use instruments (arithmetic and published bands) |
| MED | `med/critical` | CRIT | early-warning, sepsis and critical-care scores |
| MED | `med/hematology` | HEMAT | haematologic prognostic and bleeding scores |
| MED | `med/rheumatology` | RHEUM | classification criteria scores in rheumatology |
| MED | `med/infectious` | INFECT | infectious-disease scores |
| MED | `med/pediatrics` | PEDS | paediatric scores and reference values |
| MED | `med/thermoregulation` | THERMO | body-temperature staging |
| MED | `med/emergency` | EMERG | emergency-department scores not owned by a specialty |
| MED | `med/neuro` | NEURO | stroke, coma, head-injury and seizure scores |
| MED | `med/trauma` | TRAUMA | injury scores and clinical decision rules |
| MED | `med/obstetric` | OB | obstetric categories and scores |
| MED | `med/oncology` | ONCO | (planned; no function landed) |
| UNT | `unit/temperature` | TEMP | (already planned) °C, °F and K conversions (NIST SP 811) |
| LIF | `life/fitness` | FIT | (already planned) strength, load and pace |
| LIF | `life/nutrition` | NUTR | energy expenditure, DRI tables and %DV |
| NUM | `num/series` | SERIES | (already planned) time-series trend, windows and buckets |

## 3. doc/scope/SCOPE.md slice row (the slice number is Devin's)

| S? | MED/LIF/UNT/NUM | Port of the reference knowledgebase calculators: 209 functions in 23 packages (see doc/port_advomedix_2026-11-12/LEDGER.md) | [README](../port_advomedix_2026-11-12/README.md) | SHIPPED 2026-11-12 |

## 4. doc/DECISIONS.md and doc/ASSUMPTIONS.md

D-P1…D-P15 and A-P1 (closed) from [DECISIONS.md](DECISIONS.md), renumbered to the next free ids.

## 5. doc/CHANGELOG.md

`2026-11-12 — port of the reference knowledgebase: 209 functions in 23 packages (13 new med packages; life/fitness,
life/nutrition, unit/temperature, num/series new; additive in 7 medicine-run packages). 902 ledger rows: 171 PORTED,
44 REUSED, 4 REQUESTED, the rest NOT PORTED with reasons. Every formula seen in its primary (D-P5). Work:
doc/port_advomedix_2026-11-12/.`

## 6. Requests to other owners

| To | Request | From |
|---|---|---|
| `num/stats` | TukeyFences (Tukey 1977, Exploratory Data Analysis; hinges vs type-7 quartiles to be decided) | series/robust.go TukeyFences |
| `num/stats` | Welford streaming mean/variance with push/pop (Welford 1962, doi:10.1080/00401706.1962.10490022) | series/moments.go MomentsOf |
| `num/stats` | chi-square goodness of fit; frequency-weighted SD; weighted median | pending/ triage |
| `num/prob` | likelihood-ratio odds chain posterior; chi-square quantile | pending/ triage |
| `algo/prefix` (S1) | prefix sums over float64 with compensated summation | series/moments.go NewPrefixSums |
| `algo/heap` (S1) | TopK with a caller comparison (key, count with tie-break) | series/rank.go TopK |
| `algo/window` (S1) | co-occurring pairs within a window | pending/ triage |
| `num/series` (this run's package, next pass) | average real variability (Mena 2005, PMID 15716690: formula not in the abstract) | series/moments.go ARV; ranges |
| `unit/mass`, `unit/length` | kg⇄lb, cm⇄in (NIST SP 811) | pending/ triage |

## 7. doc/sources/

161 notes in [sources/](sources/), one per cited primary, to move with `mv -n` into doc/sources/.

APPLIED 2026-11-12 (lockdown pass; see doc/lockdown_2026-11-12/APPLY_LOG.md)
