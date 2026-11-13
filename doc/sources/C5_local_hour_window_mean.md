# C5 (pending diurnal.go:34) → NUM-SERIES-011 LocalHourWindowMean (PORTED)
- Citation: NIST/SEMATECH e-Handbook §1.3.5.1 "Measures of Location", https://www.itl.nist.gov/div898/handbook/eda/section3/eda351.htm: "the mean is the sum of the data points divided by the number of data points. That is, Ȳ = Σ Yᵢ/N".
- Hours of the local wall clock (package time, IANA rules); window [from, to), wrapping when from > to. Source windows 05–12 and 17–24 are parameters in the tests; the 20 mmHg surge and "non-dipping" flags are not ported.
- Check: source-shaped case morning 140 (n 2), evening 126 (n 2), difference 14. Read 2026-11-12.
