# Decisions and assumptions taken in the medicine run (2026-11-12), final form for doc/DECISIONS.md and doc/ASSUMPTIONS.md

The request said "no stops; every decision is yours, logged with the alternative and why it lost". The D-ids
below are provisional (D-M1…); the merge pass renumbers them into the next free D-ids and keeps the text.

## D-M1 Medicine packages and SHORTs (2026-11-12, lead under Devin's no-stop instruction)
Nine `med/<specialty>` packages: anthropometry ANTHRO, cardio CARDIO, respiratory RESP, renal RENAL, chemistry CHEM,
fluids FLUID, hepatic HEPAT, scales SCALE, pharmacology PHARM. Seven exist after this run (fluids and pharmacology
hold only HELD rows, so no folder was created: a package is created with its first function).
`med/chemistry` is new to the request's list: derived clinical-chemistry values (corrected calcium, osmolality,
eAG, LDL, anion gap) are one subject, not renal.
Lost: putting the lab-derived values in `med/renal` (corrected calcium, eAG and LDL are not kidney measures; one
function, one home by subject) and a `med/labs` grab-bag (CLAUDE.md forbids grab-bags).

## D-M2 Citation bar for medicine: the formula must be SEEN in the primary (2026-11-12)
A function ships only when its formula, coefficients and labels were read in the primary source (abstract, full
text, scan or official document). A confirmed citation whose formula was seen only in a secondary source is HELD
(LEFT FOR DEVIN), never shipped.
Lost: shipping "well-known" formulas (Du Bois, Mosteller, MAP one-third rule, Friedewald, MDRD, Parkland) on
secondary pointers. That would break `.claude/rules/sources.md` exactly where it was written to hold.

## D-M3 One function per unit system (2026-11-12)
Where clinicians use both systems, there are twin functions (`…SI`, `…Imperial`, `…MmolL`, `…GDL`). Each
converts once at entry with a named constant whose source is on the line above, and the tests assert that the
twins agree. Conversion constants come from the source itself where it gives one (ADAG 1.5944/2.594, the Martin
2013 caption factors 0.0259/0.0113, Dorwart's own /18 and /2.8). Otherwise they come from KDIGO 2024 p. S126
(creatinine 88.4, calcium 0.2495).
Lost: a unit parameter (`unit Unit`) on one function. It hides the unit from the name, and a wrong enum value is a
silent ×88 error.

## D-M4 Categories are typed strings carrying the source's labels verbatim (2026-11-12)
`type XCategory string`, whose constant values are the source's own words (CDC "Healthy Weight", KDIGO "Normal or
high", Rimel "severe"). KDIGO uses the code as the value plus a `Term()` method.
Lost: integer enums with a separate label table. The label is the contract, so it is the value.

## D-M5 Held back by policy: every output that is a drug dose or an amount to administer (2026-11-12, Devin's rule)
This covers the dosage family, pediatric acetaminophen, both dewormer rows and IV dose, plus fluid volumes to
administer: Parkland/ABA burn resuscitation and Holliday–Segar maintenance. The fluid volumes are amounts to
administer, so the same rule applies.
Lost: treating resuscitation volumes as "not a drug". The rule's purpose (Devin chooses the source before anything
tells a clinician how much to give) covers fluids.

## D-M6 Corrected calcium is Payne 1973's form only (2026-11-12)
Adjusted Ca = Ca − albumin + 4.0 (mg/dL, g/dL), from the Payne abstract. The common 0.8 mg/dL-per-g/dL form is not
built.
Lost: the 0.8 form. No primary source was found; a pointer attributes it to a "modified Orrell" not seen.

## D-M7 MELD: the UNOS/OPTN form as quoted in Kim 2021, plus MELDNa and MELD 3.0 (2026-11-12)
Lost: Kamath 2001's own coefficient form (0.957/0.378/1.120/0.643), which was not read in the paper.

## D-M8 Epworth: score and published bands only, no item text (2026-11-12)
The ESS is © Murray W Johns and the official site says a licence is needed to use it. calcbench ships the
arithmetic (the sum of eight 0–3 ratings) and the five band labels; the licensed caller supplies the item wording.
Lost: embedding the item text (copyright), and leaving Epworth out entirely (the request names it).
Flag for Devin: confirm whether the consumer's use needs a licence.

## D-M9 Moves (2026-11-12)
- Hagen–Poiseuille and Banana Equivalent Dose go to `sci/physics`.
- Diluting stock solutions (C₁V₁ = C₂V₂) goes to `sci/chemistry`.
- Direct Care FTEs goes to `fin/business`.
- Bayes for disease testing goes to `num/prob` (already NUM-PROB-010/-011).
- Starling stays in `med/cardio` (HELD).
- Moens–Korteweg stays in `med/cardio`, because arterial pulse-wave velocity is its use.
- Time of death and Weber fraction are LATER.
Lost: keeping Hagen–Poiseuille in cardio. It is general fluid physics, and cardio would call it.

## Assumptions (for doc/ASSUMPTIONS.md)
| ID | Assumption | Why | Where |
|---|---|---|---|
| A-M1 | Where a source reports a population but no hard input limits, the Range enforces only physical positivity and states "no limits stated by the source" | Inventing limits would be a number with no source | MED-RENAL-001/003/009/010, MED-CHEM-001…004, -007/-008, MED-HEPAT-001…004, MED-ANTHRO-001/-002/-004 |
| A-M2 | Ages are `ageYears float64` (decimal allowed) and age edges are inclusive at both ends of the source's stated range | The papers state ranges as "aged 18–92" | all age-ranged functions |
| A-M3 | KDIGO GFR categories are read as continuous lower bounds: ≥90 G1, ≥60 G2, ≥45 G3a, ≥30 G3b, ≥15 G4, else G5 | The table prints integer ranges; eGFR is continuous | MED-RENAL-011 |
| A-M4 | eAG Range is A1C 4–12% | ADAG's lowest reported band starts at 4% and its table runs 5–12; no upper limit stated | MED-CHEM-005/-006 |
| A-M5 | A Martin–Hopkins value belongs to the first stratum whose printed upper bound it does not exceed (TG 49.5 → "50-56") | Figure 2 strata are integers | MED-CHEM-009/-010 |
| A-M6 | MELD 3.0 "rounded to the nearest integer" uses math.Round (half away from zero) | The paper states no tie rule | MED-HEPAT-004 |
| A-M7 | MED-CARDIO-002 QTc = QT/∛RR (RR in s) is Fridericia's cube-root proportionality normalised to RR = 1 s | The translation states the proportionality; the printed equation is an image (constant 8.22 in 1/100 s units) | MED-CARDIO-002 |
| A-M8 | Builders' guard: an extreme in-range input whose result leaves float64 returns ErrDomain (anthropometry, renal, chemistry) or ErrNotFinite (cardio, respiratory) | numeric.md: never a silent NaN or Inf | see each doc comment |
