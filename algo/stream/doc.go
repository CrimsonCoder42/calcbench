// Package stream keeps summaries of a sequence that arrives one value at a time: a moving average over
// the last values, a hit count over the last span of time, the disjoint integer ranges seen so far, and
// running count, mean, variance, minimum and maximum.
//
// ID short form: STREAM (IDs are ALG-STREAM-NNN).
//
// The caller drives: every type is a plain value with methods. Nothing reads the clock (HitCounter takes
// now time.Time on every call), nothing starts a goroutine, and nothing locks: a type is not safe for
// concurrent use, so a caller that shares one guards it with its own mutex.
//
// Errors: a method that rejects its input (NaN or ±Inf, a negative count, a time outside the range
// time.Time.UnixNano can express, an overflow) returns the error and leaves the state unchanged.
//
// Batch statistics: num/stats computes the variance of a whole slice in two passes (another subject,
// another use). RunningStats here is the one-pass, mergeable, streaming form (Welford 1962; Chan, Golub
// and LeVeque 1979/1983); it never holds the values.
//
// Functions:
//
//	ALG-STREAM-001  MovingAverageStream  mean of the last windowLen values, O(1) amortized per value (queue.CircularBuffer)
//	ALG-STREAM-002  HitCounter           hits in the last window of time, in per-resolution buckets of a ring
//	ALG-STREAM-003  DisjointIntervals    integers added so far, as sorted disjoint closed ranges
//	ALG-STREAM-004  RunningStats         count, mean, population and sample variance, min, max; Merge combines two
//	ALG-STREAM-005  WindowedVariance     mean, sample and population variance, std of the last w values; sliding Welford (West), O(1) per Add
//	ALG-STREAM-006  WindowedCovariance   covariance, correlation, beta of the last w pairs; sliding update, O(1) per Add
//	ALG-STREAM-007  WindowedZScore       z-score of the newest value against the last w values, O(1) per Add
//	ALG-STREAM-008  OnlineCovariance     cumulative covariance, correlation, beta of pairs; Merge combines two (Chan, Golub, LeVeque)
//	ALG-STREAM-009  EWMAStream           exponential moving average one value at a time; agrees with ALG-WINDOW-014
//	ALG-STREAM-010  StreamingVWAP        cumulative volume-weighted average price; Merge combines two
//	ALG-STREAM-011  WindowedVWAP         volume-weighted average price of the last w trades, O(1) per Add
//	ALG-STREAM-012  PriceToTicks         nearest whole number of ticks for a price, ties away from zero
//	ALG-STREAM-013  TicksToPrice         price of a whole number of ticks, ticks × tickSize
//	ALG-STREAM-014  RoundToTick          price rounded to the nearest tick, ties away from zero
//	ALG-STREAM-015  BasisPointsToFraction  basis points to a fraction, bp/10 000
//	ALG-STREAM-016  FractionToBasisPoints  fraction to basis points, × 10 000
package stream
