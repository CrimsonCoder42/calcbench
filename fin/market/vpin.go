package market

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// VPIN returns the volume-synchronized probability of informed trading over n equal-volume buckets,
// VPIN = Σ_{τ=1}^{n} |V^S_τ − V^B_τ| / (n·V), where V is the bucket volume. Bucket τ's buy and sell
// volumes come from bulk volume classification (FIN-MKT-008 or FIN-MKT-009) summed over its bars; the
// caller forms the buckets and chooses n and V (the source's sample uses n = 50 buckets, V = one
// fiftieth of average daily volume; both are parameters here). The result is in [0, 1] when every
// bucket holds V.
// Errors: ErrEmpty if buyVolumes is empty; ErrLengthMismatch if the lengths differ; ErrNotFinite or
// ErrNegative for a bad volume; ErrOutOfRange if bucketVolume is not finite and > 0; ErrOverflow if a
// sum exceeds float64.
//
// ID:      FIN-MKT-010
// Time:    O(n)   n = number of buckets
// Space:   O(1), no allocation
// Source:  Easley, López de Prado, O'Hara, "Flow toxicity and liquidity in a high-frequency world", Review of Financial Studies 25(5):1457–1493, 2012, doi:10.1093/rfs/hhs053, §2 (VPIN = Σ|V^S_τ − V^B_τ|/(nV))
// Tags:    market, microstructure, vpin, order-flow-toxicity, informed-trading, volume-bucket
// Units:   volumes and bucketVolume in one size unit; result dimensionless
// Example: hand-checked: V = 100, buckets (buy, sell) = (60, 40), (30, 70) → (20 + 40)/(2·100) = 0.3
func VPIN(buyVolumes, sellVolumes []float64, bucketVolume float64) (float64, error) {
	n := len(buyVolumes)
	if n == 0 {
		return 0, fmt.Errorf("%w: buyVolumes", ErrEmpty)
	}
	if len(sellVolumes) != n {
		return 0, fmt.Errorf("%w: len(buyVolumes) = %d, len(sellVolumes) = %d", ErrLengthMismatch, n, len(sellVolumes))
	}
	if !reduce.AllFinite(bucketVolume) || bucketVolume <= 0 {
		return 0, fmt.Errorf("%w: bucketVolume = %v, want finite and > 0", ErrOutOfRange, bucketVolume)
	}
	var sum, comp float64
	for i := range n {
		if err := checkSize("buyVolumes", buyVolumes[i]); err != nil {
			return 0, fmt.Errorf("bucket %d: %w", i, err)
		}
		if err := checkSize("sellVolumes", sellVolumes[i]); err != nil {
			return 0, fmt.Errorf("bucket %d: %w", i, err)
		}
		d := sellVolumes[i] - buyVolumes[i]
		if d < 0 {
			d = -d
		}
		sum, comp = reduce.NeumaierAdd(sum, comp, d)
	}
	total := sum + comp
	v := total / float64(n) / bucketVolume
	if !reduce.AllFinite(total, v) {
		return 0, fmt.Errorf("%w: VPIN", ErrOverflow)
	}
	return v, nil
}
