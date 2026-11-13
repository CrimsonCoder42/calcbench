# TukeyFences (robust.go:29) → NUM-STATS-032 TukeyFences (PORTED)
- Citation: NIST/SEMATECH e-Handbook of Statistical Methods, §7.1.6 "What are outliers in the data?", https://www.itl.nist.gov/div898/handbook/prc/section1/prc16.htm.
- Transcribed: "lower inner fence: Q1 - 1.5*IQ", "upper inner fence: Q3 + 1.5*IQ", "lower outer fence: Q1 - 3*IQ", "upper outer fence: Q3 + 3*IQ"; Q1, Q3 "the lower and upper quartiles (defined as the 25th and 75th percentiles)"; "A point beyond an inner fence on either side is considered a mild outlier. A point beyond an outer fence is considered an extreme outlier."
- Worked example (N = 90): inner fences −39.0 and 1211.0, outer −507.75 and 1679.75 (→ Q1 429.75, Q3 742.25, IQ 312.5); tested via fencesFromQuartiles.
- Quartile method: caller's QuantileMethod (NUM-STATS-015); source parity uses type 7 (S3: 62.25, 64.75, 68.5, 72.25).
- Route: WebFetch, read 2026-11-12.
