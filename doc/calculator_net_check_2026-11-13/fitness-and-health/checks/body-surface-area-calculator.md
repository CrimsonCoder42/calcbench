# body-surface-area-calculator

## FIELDS (form calform, 2026-11-13)
csex m|f · bodyweight 155 · bodyweightunit pound|kilogram|gram · bodyheightfeet 5 / bodyheightinch 9 · bodyheight (cm, metric tab).
Results: BSA by Du Bois, Mosteller, Haycock, Gehan & George, Boyd, Fujimoto, Takahira, Schlich (m² 0.01, ft² 0.01, in² integer).

## MATCH
HAVE MED-ANTHRO-004 BodySurfaceAreaHaycock. Not built this session (time cap; primaries not read): Du Bois & Du Bois 1916 (Arch Intern Med
17:863), Mosteller 1987 (N Engl J Med 317:1098), Gehan & George 1970, Boyd 1935, Fujimoto 1968, Takahira 1925, Schlich 2010 — NEEDS SOURCE
for the next session (Du Bois 1916 is public domain and a good first target).

## COMPARE (scratch: checks/scratch/body-surface-area-calculator.go; lb × 0.45359237, in × 2.54; 1 in² = 0.00064516 m²)
| inputs | ours Haycock | theirs Haycock | verdict |
|---|---|---|---|
| 155 lb, 5'9" | 1.8522 m² (= 2871.0 in²) | 1.85 m², 2,871 in² | MATCH |
| 200 lb, 6'2" | 2.1841 m² (= 3385.4 in²) | 2.18 m², 3,385 in² | MATCH |
Submissions: 1 fields + 2 runs.
