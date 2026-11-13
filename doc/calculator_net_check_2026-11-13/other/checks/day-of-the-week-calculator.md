# Day of the Week Calculator (row 9)

CLASSIFY: in scope.
FIELDS (calcf): today_Month_ID select 0..11, today_Day_ID select 1..31, today_Year_ID text; defaults today.
MATCH: HAVE SCI-ASTRO-003 GregorianDayOfWeek (Meeus ch. 7, Example 7.e reproduced: 1954-06-30 Wednesday).
COMPARE (ours: checks/scratch/day-of-the-week-calculator.go):

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 2026-11-13 (their default; the driver could not set the year) | Friday | Friday | 0 | MATCH |
| 1954-06-30 (Meeus 7.e) | Wednesday | not settable | | source example reproduced |

DRIVER: today_Year_ID reported missing even alone. Submissions: 1.
CLOSE: MATCH (one point only; more blocked by the DRIVER gap).
