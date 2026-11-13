# Entry list W3_B10 (wave 2)

Source root (read-only): `/Users/devinanderson/Desktop/advo_production/advomedix-v2/backendAdvo/`. Paths are as wave 1 recorded them: `knowledgebase/calculators/<category>/...`, `<category>/...` (= knowledgebase/calculators/<category>/), `knowledgebase/ranges/...`, `calculators/<vital>/...`, or relative to `knowledgebase/series/` for series functions. The wave 1 reason is the starting point, not the answer.

| # | entry id | source path (wave 1) | wave 1 section | wave 1 outcome | wave 1 reason |
|---|---|---|---|---|---|
| 1 | ARV | moments.go:108 | W3_B9 | NOT PORTED | primary not seen: Mena L et al., J Hypertens 2005;23(3):505-511, doi:10.1097/01.hjh.0000160205.81652.5a, PMID 15716690 — the abstract names ARV ("inspired by the total variability concept") but states no formula; full text not open. The source cites none |
| 2 | ClassifyMean | classify.go:49 | W3_B9 | NOT PORTED | plumbing (as ClassifyValue) |
| 3 | ClassifyRules | classify.go:64 | W3_B9 | NOT PORTED | plumbing (as ClassifyValue) |
| 4 | ComputeBaseline | robust.go:60 | W3_B9 | NOT PORTED | the reference's own policy, no primary: BaselineRule windows, min points/days, and the composition "drop outside inner fences, then median and 1.4826·MAD" |
| 5 | E1RMSeries | fitness.go:199 | W3_B9 | NOT PORTED | composition, no formula of its own: per-local-day maximum of a caller-supplied estimator = NUM-SERIES-005 buckets + a max; the 1RM equations themselves are NOT PORTED (W2-B8: primaries not seen) |
| 6 | EdwardsTRIMP | fitness.go:242 | W3_B9 | NOT PORTED | as EdwardsZone (Edwards 1993 not seen) |
| 7 | EdwardsZone | fitness.go:229 | W3_B9 | NOT PORTED | primary not seen: Edwards S, The Heart Rate Monitor Book, 1993 (book, no open text); zone edges 50/60/70/80/90 % are in the source's tests only |
| 8 | FromSorted | series.go:57 | W3_B9 | NOT PORTED | as NewSeries: the O(n) order check is inside every num/series function that needs order |
| 9 | LocalDate | days.go:15 | W3_B9 | NOT PORTED | `t.In(loc).Date()` of the standard library; nothing of ours |
| 10 | MAGE | glucose.go:31 | W3_B9 | NOT PORTED | primary not seen: Service FJ et al., Diabetes 1970;19(9):644-655, doi:10.2337/diab.19.9.644, PMID 5469118 — PubMed has no abstract; full text not open. The turning-point and 1-SD rule is seen only in the source's comment |
| 11 | MomentsOf | moments.go:48 (Push :17, Pop :31, Variance/SD/CV :73-102) | W3_B9 | REUSED (batch) + REQUESTED (streaming) | Welford 1962 Technometrics 4(3):419-420 doi:10.1080/00401706.1962.10490022 is for the streaming Push/Pop only; MomentsOf itself is two-pass (Chan, Golub, LeVeque 1983, already cited by num/stats) |
| 12 | MovingTime | fitness.go:76 | W3_B9 | NOT PORTED | the reference's own policy, no primary: "moving" = segment average speed ≥ a caller threshold is a product definition (no source defines moving time) |
| 13 | NewPrefixSums | moments.go:122 (RangeSum :132) | W3_B9 | REQUESTED | prefix sums are an algo pattern (planned algo/prefix, owner S1); no domain words. Primary for the request: Bentley J, Programming Pearls 2nd ed. 2000, col. 8 (cumulative arrays), or CLRS |
| 14 | RollingUsualRange | fitness.go:269 | W3_B9 | NOT PORTED | the reference's own policy, no primary (window, minN, inner-fence inliers' min-max); fences REQUESTED above |
| 15 | SplitHalfTrend | trend.go:46 | W3_B9 | NOT PORTED | the reference's own policy, no primary: TrendRule (min readings, days, per-half, stable band 2 %) and the direction labels are product decisions; the split-at-time-midpoint comparison is defined by no source |
| 16 | TopK | rank.go:37 | W3_B9 | REQUESTED | ALG-HEAP-001 TopK ranks `[]T cmp.Ordered`; the source ranks (key, count) map entries by count desc then key asc, which ALG-HEAP-001 cannot express. To the algo/heap owner (S1): a TopKFunc with a caller's comparison (Sedgewick & Wayne 4th ed. §2.4, as ALG-HEAP-001) |
| 17 | TukeyFences | robust.go:29 | W3_B9 | REQUESTED | to num/stats owner: Tukey JW, Exploratory Data Analysis, Addison-Wesley, 1977 (book, not seen here). NOTE: Tukey's fences are built on hinges (fourths), not on type-7 quantiles as the source does; the owner decides which and cites it |
| 18 | bp_trend.go | calculators/blood_pressure/bp_trend.go:399 | W3_B10 | REUSED (SD, CV); REQUESTED (ARV) | computeSD (:399) is the population SD (divisor n) = NUM-STATS-022 PopulationStandardDeviation; computeCV (:421) is 100·SD/mean with that SD (NUM-STATS-030 uses the sample SD: a consumer must pick). computeARV (:435) Σ |

## Also yours: the PORT LATER candidates of doc/port_advomedix_2026-11-12/PENDING_TRIAGE.md (read its rows and "Later wave" table)
| # | candidate (pending/ file:line) | home | primary to find |
|---|---|---|---|
| C1 | calculators/bayesian.go:196 odds-form Bayes chain, LR+ = Se/(1−Sp), LR− = (1−Se)/Sp | num/prob | a probability/biostatistics textbook or NIST; check NUM-PROB-008…013 for reuse first |
| C2 | calculators/clustering.go:37 chi-square goodness of fit vs uniform | num/stats | NIST e-Handbook §1.3.5.15 |
| C3 | (same) chi-square critical value 7.815 (df 3, α .05) | num/prob | NIST e-Handbook §1.3.6.6.6 / §1.3.6.7.4 table |
| C4 | calculators/correlation.go:172 co-occurring pairs within a window | algo/window | a sliding-window / sweep textbook source |
| C5 | calculators/diurnal.go:34 hour-of-day window means and their difference | num/series | arithmetic mean definition (NIST §1.3.5.1); no clinical threshold |
| C6 | calculators/unit_convert.go:13 kg⇄lb, cm⇄in | unit/mass, unit/length (NEW packages, planned in doc/taxonomy.md: SHORTs MASS, LEN) | NIST SP 811 (2008) App. B: 1 lb = 0.453 592 37 kg, 1 in = 2.54 cm exactly |
| C7 | calculators/variability.go:42 average real variability | num/series | Mena 2005 J Hypertens 23:505 (corrects the source's Rothwell citation) |
| C8 | calculators/velocity.go:32 average rate of change per day | num/series | a calculus text (average rate of change); 7-day floor and 30.44-day month are app choices, not ported |
| C9 | calculators/weighted.go:44 frequency-weighted SD / CV | num/stats | NIST Dataplot weighted standard deviation |
| C10 | calculators/weighted.go:107 weighted median | num/stats | NIST Dataplot / Hyndman & Fan 1996 |
(The endocrine, respiratory and cardio candidates are with W3-B9a/B9b.)
