# finance-calculator

FIELDS (after driver update 2, form ''): cyearsv N (10), cinterestratev I/Y (6 %/yr), cstartingprinciplev PV (20,000),
ccontributeamountv PMT (−2,000), ctargetamountv FV (−10,000), cpy P/Y (1), ccy C/Y (1), ciadditionat1 (beginning|end).
The page has tabs that solve for FV, PMT, I/Y, N or PV; only the default (FV) tab is reachable.

MATCH: HAVE. FV: FIN-INT-001 + FIN-INT-014/015 (sign convention: money in negative). Other tabs map to FIN-INT-019
(PMT), FIN-INT-042 (I/Y), FIN-INT-043 (N, built session 1), FIN-INT-002 + FIN-INT-016 (PV). P/Y ≠ C/Y needs
FIN-INT-029 then FIN-INT-044 to get the per-payment rate.

COMPARE (scratch: checks/scratch/finance-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| N 10, 6 %, PV 20,000, PMT −2,000 at end, P/Y = C/Y = 1 | FV −9,455.36 | −9,455.36 | 0 | MATCH |
| N 60, 4.5 %, PV 15,000, PMT −250 at beginning, P/Y = C/Y = 12 | FV −1,927.60 | −1,927.60 | 0 | MATCH |

DRIVER: the solve-for tabs are not selectable; `sent` does not list the radio ciadditionat1 although the beginning
setting took effect.
Submissions: 2.
