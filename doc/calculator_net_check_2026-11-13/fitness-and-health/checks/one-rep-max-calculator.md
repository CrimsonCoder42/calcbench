# one-rep-max-calculator

## FIELDS (form calform, 2026-11-13)
clift 100 · cliftunit kg|lb (kg) · crepeat 5 · coutunit kg|lb (kg) · cformula e (Epley, default) | b (Brzycki) | l (Lombardi).
Results: 1 RM (0.1), a table of reps and % of 1 RM.

## MATCH
HAVE LIF-FIT-012 OneRepMaxBrzycki (LIF-FIT-014 Wathen not offered by them). Epley and Lombardi not built: NEEDS SOURCE
(Epley 1985 "Poundage chart", Boyd Epley Workout; Lombardi 1989, Beginning weight training — neither read). LIF-FIT-007 covers the
% of 1 RM table.

## COMPARE (scratch: checks/scratch/one-rep-max-calculator.go)
| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 100 kg × 5, Brzycki | 112.511 | 112.5 | 0 at 0.1 | MATCH |
| 225 lb × 8, Brzycki | 279.364 | 279.3 | 0.064 | MATCH (they truncate to 0.1; half-up would show 279.4) |
Submissions: 1 fields + 2 runs.
