# binary-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope (arbitrary-length binary addition); base conversion is stdlib strconv/math/big.
FIELDS (1 load): calform number1 (10101010), c2op {+,-,x,/}, number2 (11001100); calform2 b2dnumber1 (binary→decimal); calform3 d2bnumber1 170 (decimal→binary). Their inputs are capped at 64 characters (a 70-digit input arrived as 64 in `sent`).
MATCH: PARTIAL ALG-INTMTH-014 AddBinaryStrings (any length). Subtraction, multiplication, division and base conversion: strconv.ParseInt/FormatInt(·, 2) and math/big; no wrapper built.
COMPARE (scratch/binary-calculator.go; 5 submissions, 2 re-runs to read the value line the first filter cut off)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 10101010 + 11001100 | 101110110 | 0101110110 | leading zero only | match |
| 1×64 ones + 1 | 1 followed by 64 zeros | 1 followed by 64 zeros | 0 | match |
| 0 + 0 | 0 | not read (filter) | n/a | — |

Outcome: MATCH on addition. Their decimal echo for 2⁶⁴ − 1 + 1 prints 1.844674407371E+19 on both sides (float display), not a defect in the binary result.
