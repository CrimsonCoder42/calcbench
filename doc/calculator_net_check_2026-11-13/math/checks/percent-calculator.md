# percent-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope.
FIELDS (1 load): calform cpar1 "of" cpar2 = cpar3 (any two → third); calc c21par1/c21par2 (P% of X), c22par1/c22par2 (X is what % of Y), c23par1/c23par2 (X is P% of what); calform3 c3par1/c3par2 (percent difference, "Value 1/2"); calform2 c2par1, c2type {increase, decrease}, c2par2 (% change applied) = c2par3. All defaults empty.
MATCH: NONE (grep INDEX.md percent/relative change: only NUM-STATS-030 %RSD and basis-point converters, different maths).
RESEARCH: Kirk, Contemporary Mathematics, OpenStax 2023, §6.1 (part = percent × total, Examples 6.4–6.6); Marecek et al., Prealgebra 2e, OpenStax 2020, §6.2 "Find Percent Increase and Percent Decrease"; Törnqvist, Vartia & Vartia, Am. Stat. 1985, 39(1):43–46, doi:10.1080/00031305.1985.10479385 (relative change with the arithmetic mean as reference).
BUILD: NUM-ALGEB-029 PercentOf, -030 PercentFromPart, -031 TotalFromPercent, -032 PercentChange, -033 ApplyPercentChange (num/algebra/percent.go); -034 PercentDifference (num/algebra/percent_difference.go); tests num/algebra/percent_test.go.

COMPARE (scratch/percent-calculator.go; 6 submissions)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 70 % of 3500 (Kirk Ex. 6.4) | 2450 | 2450 | 0 | match |
| 40 vs 50, percent difference | 22.22222222222222 | 22.222222222222 | 0 at their 12 dp | match |
| 500 decreased by 12.5 % | 437.5 | 437.5 | 0 | match |
| 155 is what % of 228 (Kirk Ex. 6.6: 67.98 %) | 67.98245614035088 | not captured | n/a | DRIVER |
| 4000 is 10 % of what (Kirk Ex. 6.5: 40 000) | 40000 | not captured | n/a | DRIVER |

DRIVER: form `calc` (the three phrase calculators c21/c22/c23) — `run` reports all fields missing and `sent` empty (twice; once re-run raw to look, not retried further). Ours matches the source's worked examples for both.
Their percent difference uses the mean (x+y)/2 as reference, the same as ours.

## Session 2 retry (1 submission, total 7)
`run … '{"c22par1":"155","c22par2":"228"}' ''` (whole page): `sent` now holds both fields and `missing` is empty, but the submit navigated through another form's Calculate button and the result block was empty (exit 4); `fieldsAfter` shows c21–c23 cleared. DRIVER: the phrase calculators need their own button (label unknown to the driver) or live evaluation; still not captured. One submission left; not used.

## Session 3 (form name `calc`; 1 submission, total 8)
`sent` came back empty ({}) and no result: the fields c21–c23 are still not reachable under form `calc`. Limit reached; the phrase forms stay uncompared. Ours already matches Kirk §6.1 Examples 6.5 and 6.6 in the tests.
