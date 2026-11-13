package market

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Level1Imbalance returns the top-of-book queue imbalance I = Q^b/(Q^b + Q^a), in [0, 1]: 1 when only
// bids rest at the touch, 0 when only asks do, ½ when balanced. The signed form (Q^b − Q^a)/(Q^b + Q^a)
// is 2I − 1.
// Errors: ErrNotFinite if a size is NaN or ±Inf; ErrNegative if a size is negative; ErrDomain if both
// sizes are 0; ErrOverflow if Q^b + Q^a exceeds float64.
//
// ID:      FIN-MKT-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Stoikov, "The micro-price: a high-frequency estimator of future prices", Quantitative Finance 18(12):1959–1966, 2018, doi:10.1080/14697688.2018.1489139, §2 (I = Q^b/(Q^b + Q^a))
// Tags:    market, microstructure, imbalance, queue-imbalance, order-book, level-1
// Units:   bidSize, askSize in one size unit (shares, lots); result dimensionless in [0, 1]
// Example: hand-checked: Q^b = 300, Q^a = 100 → 0.75
func Level1Imbalance(bidSize, askSize float64) (float64, error) {
	if err := checkSize("bidSize", bidSize); err != nil {
		return 0, err
	}
	if err := checkSize("askSize", askSize); err != nil {
		return 0, err
	}
	total := bidSize + askSize
	if total == 0 {
		return 0, fmt.Errorf("%w: bidSize + askSize = 0", ErrDomain)
	}
	if !reduce.AllFinite(total) {
		return 0, fmt.Errorf("%w: bidSize + askSize", ErrOverflow)
	}
	return bidSize / total, nil
}
