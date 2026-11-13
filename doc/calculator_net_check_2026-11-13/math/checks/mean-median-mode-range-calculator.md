# mean-median-mode-range-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope.
FIELDS (1 load): calform numberinputs (textarea; default "10, 2, 38, 23, 38, 23, 21").
MATCH: HAVE NUM-STATS-001 ArithmeticMean, -012 Median, -024 Mode, -018 Range; NUM-STATS-049 GeometricMean (built this session for statistics-calculator).
COMPARE (scratch/statistics-calculator.go; 2 submissions)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 10, 2, 38, 23, 38, 23, 21 | mean 22.1428571428571, median 23, range 36, modes 23, 38 ×2, GM 16.412764443111307, sum 155 | 22.142857142857, 23, 36, 38, 23 ×2, 16.412764443111, 155 | 0 (mode order: ours ascending, theirs unordered) | match |
| 1.5, −2.25, 7, 7, 3.125, −2.25, 0 | mean 2.01785714285714, median 1.5, range 9.25, modes −2.25, 7 ×2, GM ErrNegative | 2.0178571428571, 1.5, 9.25, −2.25, 7 ×2, GM 0 | GM only | THEIRS DIFFERS (GM, see statistics-calculator.md) |
