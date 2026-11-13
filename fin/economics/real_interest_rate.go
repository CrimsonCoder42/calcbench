package economics

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// checkRate returns ErrRate (also wrapping ErrNotFinite for NaN and ±Inf) unless rate is finite and > −1.
func checkRate(name string, rate float64) error {
	if !reduce.AllFinite(rate) {
		return fmt.Errorf("%w: %w: %s = %v", ErrRate, ErrNotFinite, name, rate)
	}
	if rate <= -1 {
		return fmt.Errorf("%w: %s = %v", ErrRate, name, rate)
	}
	return nil
}

// RealInterestRate returns the real rate of interest implied by a nominal rate and an inflation rate over
// the same period, from the exact Fisher relation 1 + i = (1 + r)(1 + π): r = (1 + i)/(1 + π) − 1,
// evaluated as (i − π)/(1 + π) so that nearly equal rates do not cancel. It is also the inflation-adjusted
// return of an investment that earned i. The approximation r ≈ i − π is not used.
// Errors: ErrRate if either rate is NaN, ±Inf (these also wrap ErrNotFinite) or ≤ −1; ErrOverflow if the
// result is beyond float64 (1 + inflationRate near 0 with a large nominal rate).
//
// ID:      FIN-ECON-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Fisher, I., Appreciation and Interest, Publications of the American Economic Association 11(4), Macmillan for the AEA, 1896, ch. II §3, eq. (2) "1 + j = (1 + a)(1 + i)" and eq. (3), p. 9, https://fraser.stlouisfed.org/files/docs/meltzer/appreciationinte00fish.pdf (read 2026-11-12)
// Tags:    economics, interest, inflation, real-rate, fisher, inflation-adjusted-return
// Units:   nominalRate, inflationRate and the result decimal per the same period
// Example: hand-checked from OpenStax Principles of Finance §3.4 ($1,000 at 6 % → $1,060; prices +2 % → $1,020): 0.04/1.02 = 0.0392156862745098
func RealInterestRate(nominalRate, inflationRate float64) (float64, error) {
	if err := checkRate("nominalRate", nominalRate); err != nil {
		return 0, err
	}
	if err := checkRate("inflationRate", inflationRate); err != nil {
		return 0, err
	}
	r := (nominalRate - inflationRate) / (1 + inflationRate)
	if err := checkResult("real rate", r); err != nil {
		return 0, err
	}
	return r, nil
}
