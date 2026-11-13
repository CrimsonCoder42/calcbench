# target-heart-rate-calculator

## FIELDS (form calform, 2026-11-13)
mm a (estimate from age, default) | t (test result) · age 30 · mhr (tested max HR) · rhr 70 · cformula h (Haskell & Fox 1971, default) |
t (Tanaka 2001) | n (Nes 2013) · cintensity k (Karvonen, default) | b (Borg RPE) | m (Borg CR10).
Results: target range at 50-85 % of heart-rate reserve, and five zones with bpm (integers).

## MATCH
HAVE MED-CARDIO-001 MaxHeartRateTanaka, MED-CARDIO-030 MaxHeartRateNes (and -029 Gulati, not offered by them).
NOT BUILT, NEEDS SOURCE: Karvonen heart-rate-reserve target (Karvonen MJ, Kentala E, Mustala O. Ann Med Exp Biol Fenn
1957;35(3):307-315 — not available online; no two independent authoritative secondaries found this session, D-075/D-076) and
Haskell & Fox 220 − age (Fox SM, Naughton JP, Haskell WL. Ann Clin Res 1971;3:404-432 — not read). Their zone labels and % bands
also need a source (ACSM position stand, Garber et al. MSSE 2011;43:1334 Table, not read).

## COMPARE (scratch: checks/scratch/target-heart-rate-calculator.go)
| inputs | ours | theirs (100 % HRR = HRmax) | verdict |
|---|---|---|---|
| age 30, Tanaka | 187.000 | 187 | MATCH |
| age 50, Nes | 179.000 | 179 | MATCH |
Their Karvonen figures are consistent with (HRmax − rest) × p + rest rounded half-up (e.g. 0.5 × 117 + 70 = 128.5 → 129). Not ours.
Submissions: 1 fields + 2 runs.
