# right-triangle-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope.
FIELDS (1 load): calform av (3), alphav, alphaunit {d, r}, bv, betav, betaunit {d, r}, cv (5), hv (altitude), areav, perimeterv — any two give the rest.
MATCH: HAVE NUM-TRIG-027..030 SolveRightTriangleLegs, …HypotenuseLeg, …LegAngle, …HypotenuseAngle. Altitude, area, perimeter, in/circumradius are one-line derivations from the solved sides (not compared). Solving from h, area or perimeter not built (noted in LEAD.md).
COMPARE (scratch/right-triangle-calculator.go; 3 submissions)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| a 3, c 5 | b 4, α 36.8698976458°, β 53.1301023542° | 4, 36.87°, 53.13° | 0 at displayed figs | match |
| a 7, α 25° | b 15.0115484436, c 16.5634110821, β 65° | 15.01155, 16.56341, 65° | 0 at 5 dp | match |
| c 10, α 0.3 rad | a 2.95520206661, b 9.55336489126, β 72.8112661461° | 2.9552, 9.55336, 72.811° | 0 at displayed figs | match |
