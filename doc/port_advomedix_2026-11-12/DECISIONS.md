# Decisions of the port run (lead, under the request's no-stop instruction)

D-P ids are provisional; the merge pass renumbers them into doc/DECISIONS.md.

## D-P1 The medicine run's packages are additive-only in this run
The medicine agent finished: no change in med/ after 14:19, and its report was delivered. Its seven packages take new
files and new IDs (after the last), but none of its 44 files is edited, and the lead places the doc.go lines.
Lost: treating them as fully off-limits and pushing every cardio/renal/respiratory/hepatic port into a follow-up wave.
That would idle two of the four W1 builders, while the agent the rule protects is no longer writing.

## D-P2 The source path lives only in the ledger; parity tests cite the ledger row
Parity tests carry `// parity: ledger <entry id>, source case N`, never a file path or a product name.
Lost: the source file:line in the test comment. The request also says nothing in calcbench may name the product,
and the ledger is the only place the path may appear; the ledger row has the file:line.

## D-P3 Routing by home package, not by source folder
The source folders (general, cardiac…) mix specialties. Each entry is routed to the builder that owns its calcbench
home package, so no two builders share a package within a wave. The waves keep the request's order.
Lost: one builder per source folder writing into whatever package fits. Two builders would write into med/gastro
(Rome criteria live in general/ and metabolic/), which breaks "builders never touch another builder's package".

## D-P4 New packages and SHORTs
med/gastro GASTRO · med/endocrine ENDO · med/psychiatry PSYCH · med/critical CRIT · med/hematology HEMAT ·
med/oncology ONCO · med/rheumatology RHEUM · med/infectious INFECT · med/pediatrics PEDS ·
med/thermoregulation THERMO · med/emergency EMERG · med/neuro NEURO · med/trauma TRAUMA · med/obstetric OB.
Also, from doc/taxonomy.md: unit/temperature TEMP and life/fitness FIT. life/nutrition NUTR is new.
Lost: a single med/scales holding every scored instrument. Packages are by subject (CLAUDE.md), and a scale is a
shape, not a subject.

## D-P5 Citation bar carried over (medicine D-M2) and policy rows
A formula is ported only when it is SEEN in the primary. Dose and amount-to-administer rows, decision pathways and
administrative rules are NOT PORTED, with the reason recorded.
Lost: porting on the reference's citation. The request calls it a hypothesis, and the medicine run found several
wrong ones (PEF, Liebermeister, Harris–Benedict).

## D-P6 VTE risk scores live in med/cardio
IMPEDE VTE (myeloma) and TRiP(cast) (lower-limb immobilisation) are venous-thromboembolism risk scores, and VTE
(Wells, Geneva, PERC, D-dimer) is homed in med/cardio. Lost: med/hematology for IMPEDE and med/trauma for TRiP(cast);
that would split the VTE family across three packages by the population each one studied, when the subject is the same.

## D-P7 MED-CARDIO-009 withdrawn before release
The Canadian Syncope Risk Score point table obtained (sum −3…9) contradicts the paper's stated range (−3…11), so the
builder withdrew the function. The ID is not reused. Its INDEX row reads "MED-CARDIO-009 | REMOVED before release |
Canadian Syncope Risk Score: point table not confirmed".

## D-P8 Series kernel: num/series is ours; num/stats and algo are reused, never written
The reference's series/ math goes to num/series, a new package. taxonomy.md plans it as "sequences and series,
trend and smoothing", and no subject has claimed it. num/stats belongs to the statistics subject (29 functions
built; another subject is still live in num/), and algo/* belongs to S1. A series function that duplicates one of
theirs (median, MAD, percentile, TopK, interval merge, rolling sums) is REUSED. A statistic they lack is a
PENDING_SHARED request to their owner, and is never written here.
Lost: writing the robust statistics into num/stats ourselves. That edits another subject's package while it may
be live.

## D-P9 Fitness loads take one unit, any (`…KgOrLb`), no kg/lb twins
1RM, load-at-percent and volume-load equations are homogeneous in load: the result is in the unit of the input.
Lost: kg and lb twins. They would be identical code with two IDs; the FENa precedent (MED-RENAL-009) already takes
"one unit, any" for a ratio.

## D-P10 LIF-FIT-004/005 (Mayhew 1992) withdrawn before release
The only text of the abstract seen was the reference's own copy of the publisher page, and the live page returns
402. The reference's copy is the hypothesis, not the primary (D-P5). Files moved to _delete/life/fitness/. IDs
retired, not reused.
Lost: accepting the copied abstract. It would make the reference its own witness.

## D-P11 "patient" is allowed only inside a verbatim title in a Source: line
A citation's title is quoted exactly. Everywhere else (identifiers, prose, tests, doc.go) the identity rule holds.
Lost: rewording titles (a corrupted citation) or dropping titles (an incomplete citation, against sources.md).

## D-P12 %DV rounding is half-even
21 CFR 101.9(d)(7)(ii) says "nearest whole percent" and states no tie rule, so .claude/rules/numeric.md's default
(half-even) applies. The one parity case (2.5 g of 20 g DV: source 13, ours 12) is a recorded disagreement for
Devin.
Lost: half-up. It matches the reference, but no cited rule supports it.

## A-P1 (assumption) AKIN creatinine comparisons round to 0.001 mg/dL first
The AKIN stage compares a creatinine ratio with the cut-offs 1.5, 2.0 and 3.0. Rounding both creatinines to
0.001 mg/dL (the laboratory reporting precision is coarser) keeps a printed ratio such as 3.5/1.75 = 2.0 from
falling on the wrong side of the cut-off through float error.
Alternative: compare cr ≥ k × baseline by multiplication, which needs no division. The reviewer should check
whether that removes the need for rounding, and if it does, prefer it.

## D-P13 A definition from the SI Brochure is a source for pace
Pace is the reciprocal of speed. The BIPM SI Brochure (9th ed., 2019) defines speed, so LIF-FIT-009 cites it,
as LIF-FIT-007 cites the definition of "% of 1RM".
Lost: NOT PORTED for want of a paper. A unit definition is the primary for a unit relation.

## D-P14 The NEWS2 SpO₂ bands stay in two packages, guarded by a cross-check test
MED-CRIT-001 NEWS2 keeps its own SpO₂ switch, and MED-RESP-009/010 export the same bands. med/* packages do not
import each other in production code (medicine CONVENTIONS). An external test,
med/critical/news2_spo2_crosscheck_test.go, compares the two over every SpO₂ value 70–100 on both scales, on air and
on oxygen, so a drift fails the build.
Lost: having NEWS2 call med/respiratory, which would be the first production import between med packages.

## D-P15 Wording-only edits to three medicine-run files (exception to D-P1)
The request's identity rule applies to all of calcbench, so "patients" describing study populations outside
quotation marks was changed to "persons" or the population's own words. Files: med/hepatic/fib4.go (2), the
Cockcroft–Gault file (1) and the Epworth score file (1). No code changed. The medicine agent had finished, so the
reason for D-P1 does not apply to a comment fix.
Lost: leaving them. That breaks the identity rule, which the reviewer checks across the tree.

## A-P1 closed (reviewer): keep the thousandths rounding in AKIN
The reviewer probed 0.01 mg/dL pairs sitting exactly on the ×1.5, ×2 and ×3 cut-offs. Comparing cr ≥ k × baseline in
float64 misjudged 480 pairs, and the division form misjudged 375. Comparing whole thousandths, as the code does, is
exact. Overflow of the thousandths now returns ErrDomain (review MAJOR 3).
