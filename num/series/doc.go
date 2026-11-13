// Package series computes over time series and number sequences: the robust Theil–Sen trend slope,
// trailing sums over a window of time, and calendar bucketing of time-stamped points by local day, by
// local week and by a fixed stride.
//
// ID short form: SERIES (IDs are NUM-SERIES-NNN).
//
// Inputs: a Point is an instant (time.Time) and a value in the caller's unit; results keep that unit.
// Functions that walk a series in time order take points sorted by time (non-decreasing) and check it in
// O(n) (ErrUnsorted); none sorts or modifies its input. Instants, locations and the week's first day are
// always passed in: nothing reads the clock or the local zone. Calendar arithmetic is the proleptic
// Gregorian calendar of package time with the IANA rules of the *time.Location given, so a local day is
// 23, 24 or 25 hours. NaN or ±Inf in a value is ErrNotFinite, never a NaN result. Sums are compensated
// (Neumaier). Randomized algorithms take a caller's *rand.Rand and their result does not depend on it.
//
// Descriptive statistics (mean, variance, median, percentiles, MAD, least-squares slope, z-scores) live in
// num/stats and index-window sums in algo/window; this package does not repeat them. series imports
// algo/reduce, num/stats (LinearRegression for AR1Fit and HurstExponentRS) and num/prob
// (ChiSquareSurvival for LjungBox's p-value), which precede it in D-009's order.
//
// Functions:
//
//	NUM-SERIES-001  TheilSenSlope          median of the pairwise slopes, randomized slope selection
//	NUM-SERIES-002  TheilSenSlopePerDay    Theil–Sen slope of point values per day of elapsed time
//	NUM-SERIES-003  TrailingTimeWindowSum  sum and count of the values in (t − window, t] at each point
//	NUM-SERIES-004  LocalDayWindow         [from, to) covering the last N local calendar days
//	NUM-SERIES-005  LocalDayBuckets        index ranges of the points of each local calendar date
//	NUM-SERIES-006  FixedStrideBuckets     min, mean, max and count of the points in each fixed stride
//	NUM-SERIES-007  WeekStart              local midnight of the first day of an instant's week
//	NUM-SERIES-008  WeeklySums             sum and count of the values in each local week of a window
//	NUM-SERIES-009  AverageRealVariability     mean absolute change between consecutive readings (ARV)
//	NUM-SERIES-010  AverageRateOfChangePerDay  (last − first)/elapsed days between two points, the secant slope
//	NUM-SERIES-011  LocalHourWindowMean        mean and count of the points in a window of local clock hours
//	NUM-SERIES-012  Autocorrelation            sample autocorrelation r_k = c_k/c_0 at one lag (divisor N)
//	NUM-SERIES-013  ACF                        sample autocorrelations r_0…r_m into dst
//	NUM-SERIES-014  PACF                       partial autocorrelations φ_kk by the Durbin–Levinson recursion
//	NUM-SERIES-015  LjungBox                   Ljung–Box Q, its degrees of freedom and χ² p-value
//	NUM-SERIES-016  VarianceRatio              Lo–MacKinlay VR(q) and its homoskedastic z statistic
//	NUM-SERIES-017  HurstExponentRS            Hurst exponent by rescaled-range (R/S) analysis
//	NUM-SERIES-018  AR1Fit                     AR(1) intercept and φ by conditional least squares
//	NUM-SERIES-019  MeanReversionHalfLife      ln 2/(−ln φ), the half-life of an AR(1) deviation
//	NUM-SERIES-020  ArithmeticSequenceTerm     n-th term a₁ + (n − 1)d
//	NUM-SERIES-021  ArithmeticSeriesSum        sum of the first n terms of an arithmetic sequence
//	NUM-SERIES-022  GeometricSequenceTerm      n-th term a₁rⁿ⁻¹
//	NUM-SERIES-023  GeometricSeriesSum         sum of the first n terms of a geometric sequence
package series
