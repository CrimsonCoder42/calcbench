package electrical

import (
	"errors"
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// ErrEmpty reports an empty list of resistors.
var ErrEmpty = errors.New("electrical: no resistors")

// SeriesResistance returns the equivalent resistance of resistors in series, R = R1 + R2 + … + Rn, summed with
// Neumaier compensation.
// Errors: ErrEmpty for no resistors; ErrNotFinite or ErrNegative for a bad element (its index is named);
// ErrOverflow when the sum is beyond float64.
//
// ID:      ENG-ELEC-009
// Time:    O(n)   n = len(resistancesOhms)
// Space:   O(1), no allocation
// Source:  Ling SJ, Sanny J, Moebs W. University Physics Volume 2. OpenStax, 2016, https://openstax.org/details/books/university-physics-volume-2, §10.2 "Resistors in Series and Parallel", eq. (10.2) (R_S = R1 + R2 + … + R_N)
// Tags:    resistors, series, equivalent-resistance, circuit, electrical
// Units:   resistancesOhms Ω each; result Ω
// Example: 10, 2, 38, 23, 38, 23, 21 Ω in series → 155 Ω
func SeriesResistance(resistancesOhms []float64) (float64, error) {
	if len(resistancesOhms) == 0 {
		return 0, ErrEmpty
	}
	var sum, comp float64
	for i, r := range resistancesOhms {
		if err := checkElement(i, r); err != nil {
			return 0, err
		}
		sum, comp = reduce.NeumaierAdd(sum, comp, r)
	}
	total := sum + comp
	if math.IsInf(total, 0) || math.IsNaN(total) {
		return 0, fmt.Errorf("%w: series sum", ErrOverflow)
	}
	return total, nil
}

// ParallelResistance returns the equivalent resistance of resistors in parallel, 1/R = 1/R1 + 1/R2 + … + 1/Rn,
// with the reciprocals summed with Neumaier compensation. A zero-ohm branch shorts the network: the result is 0;
// so does a branch whose conductance 1/R, or a total conductance, overflows float64 (subnormal resistances).
// Errors: ErrEmpty for no resistors; ErrNotFinite or ErrNegative for a bad element (its index is named).
//
// ID:      ENG-ELEC-010
// Time:    O(n)   n = len(resistancesOhms)
// Space:   O(1), no allocation
// Source:  Ling SJ, Sanny J, Moebs W. University Physics Volume 2. OpenStax, 2016, https://openstax.org/details/books/university-physics-volume-2, §10.2 "Resistors in Series and Parallel", eq. (10.3) (1/R_P = 1/R1 + 1/R2 + … + 1/R_N)
// Tags:    resistors, parallel, equivalent-resistance, circuit, electrical
// Units:   resistancesOhms Ω each; result Ω
// Example: 10, 2, 38, 23, 38, 23, 21 Ω in parallel → 1.27031366795 Ω; 6 Ω ∥ 3 Ω → 2 Ω
func ParallelResistance(resistancesOhms []float64) (float64, error) {
	if len(resistancesOhms) == 0 {
		return 0, ErrEmpty
	}
	var sum, comp float64
	short := false
	for i, r := range resistancesOhms {
		if err := checkElement(i, r); err != nil {
			return 0, err
		}
		if r == 0 {
			short = true
			continue
		}
		inv := 1 / r
		if math.IsInf(inv, 0) {
			// A subnormal resistance: its conductance overflows, a short to float64 precision.
			short = true
			continue
		}
		sum, comp = reduce.NeumaierAdd(sum, comp, inv)
		if math.IsInf(sum, 0) || math.IsNaN(sum+comp) {
			// The total conductance overflowed: a short to float64 precision.
			short = true
		}
	}
	if short {
		return 0, nil
	}
	return 1 / (sum + comp), nil
}

// checkElement validates one resistor, formatting its index only on failure (no allocation on success).
func checkElement(i int, r float64) error {
	if math.IsNaN(r) || math.IsInf(r, 0) || r < 0 {
		return checkInput(fmt.Sprintf("resistancesOhms[%d]", i), r)
	}
	return nil
}
