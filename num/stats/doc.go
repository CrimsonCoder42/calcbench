// Package stats computes descriptive statistics: location, scale, frequency tables, grouped-data
// estimates of the mean and median, percentiles, correlation, the least-squares line, and standardized scores.
//
// ID short form: STATS (IDs are NUM-STATS-NNN).
//
// Inputs: float64 slices in the caller's unit; a result keeps that unit (a variance is in unit²). Mode is
// generic over cmp.Ordered. No function modifies its input; Median, MedianAbsoluteDeviation, Percentile,
// Percentiles, Mode and SummaryStatistics reorder a copy. NaN or ±Inf in a float64 input is an error (ErrNotFinite),
// never a NaN result. Population and sample forms are separate functions (doc/DECISIONS.md D-010).
//
// Accuracy: sums are compensated (Neumaier); sums of squares and cross-products use the corrected
// two-pass formula, never the one-pass Σx² − (Σx)²/n. A sum that overflows, or one whose deviations are
// so small (below 2^-480) that their squares would underflow, is redone on values scaled by a power of two
// (exact), so spreads from subnormal (about 5e-324) to near MaxFloat64 are handled; ErrOverflow means the
// result itself is beyond the float64 range, and ErrZeroVariance means every value is exactly equal. A
// variance or sum of squares below the smallest float64 rounds to a subnormal or 0, as any float64 result
// does; the standard deviations are computed from the scaled sum and stay accurate.
//
// Percentile takes an explicit QuantileMethod, one of the nine Hyndman–Fan types; the documented default
// is DefaultQuantileMethod, type 7 (D-010).
//
// Functions:
//
//	NUM-STATS-001  ArithmeticMean                    sum of the values divided by their count
//	NUM-STATS-002  RelativeFrequency                 each class count divided by the total count
//	NUM-STATS-003  CoefficientOfVariation            sample standard deviation divided by the mean
//	NUM-STATS-004  CohenD                            difference of two means in standard deviations (Cohen's d)
//	NUM-STATS-005  GroupedMean                       estimated mean of grouped data from class marks
//	NUM-STATS-006  GroupedMedian                     estimated median of grouped data by interpolation
//	NUM-STATS-008  FrequencyDistribution             count of values in each class
//	NUM-STATS-009  LinearRegression                  intercept and slope of the least-squares line
//	NUM-STATS-010  MatthewsCorrelation               Matthews correlation coefficient of a 2×2 table
//	NUM-STATS-011  MeanAbsoluteDeviation             average absolute deviation about the mean
//	NUM-STATS-012  Median                            middle value, or mean of the two middle values
//	NUM-STATS-013  Midrange                          mean of the smallest and largest values
//	NUM-STATS-014  PearsonCorrelation                Pearson correlation coefficient of paired samples
//	NUM-STATS-015  Percentile                        p-quantile by a chosen Hyndman–Fan method
//	NUM-STATS-016  PooledStandardDeviation           pooled standard deviation of two samples
//	NUM-STATS-017  PopulationVariance                variance with divisor N
//	NUM-STATS-018  Range                             largest value minus smallest value
//	NUM-STATS-019  SampleStandardDeviation           standard deviation with divisor n − 1
//	NUM-STATS-020  SummaryStatistics                 count, mean, sample SD, min, max, median and range
//	NUM-STATS-021  SampleVariance                    variance with divisor n − 1
//	NUM-STATS-022  PopulationStandardDeviation       standard deviation with divisor N
//	NUM-STATS-023  StandardErrorOfMean               sample standard deviation divided by √n
//	NUM-STATS-024  Mode                              most frequent values, all ties, ascending
//	NUM-STATS-025  SumOfSquares                      sum of squared deviations about the mean
//	NUM-STATS-026  WeightedMean                      weighted mean Σw·x / Σw
//	NUM-STATS-028  ZScore                            standard score (x − mean)/sd
//	NUM-STATS-029  RawScoreFromZScore                value with a given standard score, mean + z·sd
//	NUM-STATS-030  PercentRelativeStandardDeviation  100·s/abs(mean), the percent relative standard deviation
//	NUM-STATS-031  MedianAbsoluteDeviation           median absolute deviation about the median
//	NUM-STATS-032  TukeyFences                                   quartiles, IQ and inner/outer fences Q1 − 1.5·IQ … Q3 + 3·IQ
//	NUM-STATS-033  ChiSquareGoodnessOfFit                        Pearson χ² = Σ (O − E)²/E against expected counts
//	NUM-STATS-034  ChiSquareGoodnessOfFitUniform                 Pearson χ² against equal expected counts N/k
//	NUM-STATS-035  FrequencyWeightedPopulationStandardDeviation  population SD of a frequency table, √(Σf(x − μ)²/N)
//	NUM-STATS-036  FrequencyWeightedSampleStandardDeviation      sample SD of a frequency table, √(Σf(x − x̄)²/(n − 1))
//	NUM-STATS-037  FrequencyWeightedMedian                       median of a frequency table without expanding it
//	NUM-STATS-038  SampleCovariance                              sample covariance Σ(x − x̄)(y − ȳ)/(n − 1)
//	NUM-STATS-039  PopulationCovariance                          population covariance Σ(x − μx)(y − μy)/N
//	NUM-STATS-040  CovarianceMatrix                              k×k sample covariance matrix of k variables (row-major)
//	NUM-STATS-041  CorrelationMatrix                             k×k Pearson correlation matrix of k variables (row-major)
//	NUM-STATS-042  Percentiles                                   several Hyndman–Fan quantiles from one sort
//	NUM-STATS-043  HistogramBinEdges                             equal-width class edges by Sturges, Scott, Freedman–Diaconis or Doane
//	NUM-STATS-044  FrequencyDensity                              histogram density count/(n·width) per class
//	NUM-STATS-045  FrequencyDistribution2D                       bivariate histogram counts on two edge sets (row-major)
//	NUM-STATS-046  Digitize                                      class index of each value among half-open intervals cut by edges
//	NUM-STATS-047  SampleSizeProportion                          sample size to estimate a proportion to ±d (Cochran, finite N)
//	NUM-STATS-048  MarginOfErrorProportion                       ±d of a proportion at a confidence, with finite-population correction
//	NUM-STATS-049  GeometricMean                                 (∏xᵢ)^(1/n) of non-negative values, by the mean of the logs
//	NUM-STATS-050  MeanConfidenceIntervalZ                       x̄ ± z·σ/√n, two-sided, σ known
//	NUM-STATS-051  MeanConfidenceIntervalT                       x̄ ± t·s/√n, two-sided, n − 1 degrees of freedom
//
// Reserved, not shipped: NUM-STATS-007 GroupedMode is held (no confirmed source; its code is not in this
// package) and NUM-STATS-027 WeightedGeometricMean is out (D-010, no confirmed locator). Both IDs stay
// taken and are never reused.
package stats
