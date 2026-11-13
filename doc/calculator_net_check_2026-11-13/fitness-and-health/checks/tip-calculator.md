# tip-calculator (HANDOFF from A-OTHER; notes ../other/checks/tip-calculator.md)

## FIELDS (2026-11-13)
calform: serviceprice 55 (a service-price guide). calform2: billamount 55 · tippct 15 · numberofpeople 1. Results: tip, total (and per person).

## BUILD
LIF-SHOP-001 TipCents — NEW PACKAGE life/shopping (taxonomy SHORT SHOP): doc.go, tip.go, tip_test.go. int64 cents × int64 basis points,
half-to-even to the cent (numeric.md default; the percent definition is NIST SP 811 §7.10.2). Split among people and its remainder:
not built (numeric-policy decision for Devin, LEAD.md).

## COMPARE (scratch: checks/scratch/tip-calculator.go)
| inputs | ours | theirs | verdict |
|---|---|---|---|
| $10.10, 15 % (exact tie 151.5 c) | 152 c, total 1162 c | $1.52, $11.62 | MATCH (half-even and half-up agree here) |
| $47.33, 18 % | 852 c, total 5585 c | $8.52, $55.85 | MATCH |
Submissions: 1 fields + 2 runs.
