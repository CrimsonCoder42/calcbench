# Port wave 2 (2026-11-12)

Request: the wave 2 port request (lead + builders + one reviewer). Takes every NOT PORTED and REQUESTED row of
wave 1 (`../port_advomedix_2026-11-12/LEDGER.md`) and the 17 PORT LATER candidates of its PENDING_TRIAGE.md, and
ends each in one final state: PORTED · PORTED-SECONDARY · REUSED · NEEDS PDF · HELD BY POLICY · NOT A CALCULATION ·
SOURCE DOES NOT SUPPORT. Nothing ported in wave 1 is rewritten.

Binding: WAVE2_CONVENTIONS.md (over wave 1's PORT_CONVENTIONS.md and the medicine CONVENTIONS.md).

Waves (builders own disjoint packages within a wave; at most 4 at once):
- W1: B1 cardiac (med/cardio) · B2 respiratory (med/respiratory) · B3 metabolic and endocrine (med/endocrine,
  gastro, hepatic, chemistry, anthropometry) · B4 general, first half (med/hematology, rheumatology, infectious,
  emergency, critical, scales, pediatrics, oncology).
- W2: B5 general, second half (W1's B4 packages minus critical, plus cardio and respiratory) · B6 neuro and
  psychiatry · B7 renal, sepsis, trauma (med/renal, critical, trauma) · B8 fitness, nutrition, obstetric.
- W3: B9 lookup tables and percentile charts (ranges) · B10 the 17 candidates, the REQUESTED rows and the series leftovers.

Files: PROGRESS.md (one line per entry) · LEDGER.md (all sections) · CONFIRM.md · NEEDS_PDF.md · CONSUMER_DEFECTS.md ·
REPORT.md · briefs/ · ledger/ (one section per builder, LEAD.md for rows finalised from wave 1 evidence) · the per-entry source notes (moved to doc/sources/<entry>.md at the end) · papers/ (Devin drops <PMID>.pdf here; gitignored) ·
out/ (raw test and benchmark output).
