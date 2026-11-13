package stream

import (
	"errors"
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// ErrVolume: a trade's volume is negative.
var ErrVolume = errors.New("stream: volume must not be negative")

// StreamingVWAP is the volume-weighted average price of the trades added so far,
// VWAP = Σ priceᵢ·volumeᵢ / Σ volumeᵢ, kept one trade at a time in two compensated (Neumaier) sums, so
// its error does not grow with the number of trades. Merge combines two summaries of separate parts of
// a stream (the sums add). The zero value is an empty summary, ready to use. Not safe for concurrent use.
//
// Methods: Add(price, volume) adds a trade. VWAP returns the average, or (0, false) while the total
// volume is 0. TotalVolume and Count read the sums. Merge(other) adds every trade other summarises.
// Reset empties the summary.
//
// ID:      ALG-STREAM-010
// Time:    O(1) per Add, per Merge and per reader
// Space:   O(1), no allocation
// Source:  Berkowitz, Logue, Noser, "The total cost of transactions on the NYSE", The Journal of Finance, 1988, 43(1):97–112, doi:10.1111/j.1540-6261.1988.tb02591.x (the volume-weighted average price benchmark)
// Source:  Neumaier, "Rundungsfehleranalyse einiger Verfahren zur Summation endlicher Summen", ZAMM, 1974, 54(1):39–51, doi:10.1002/zamm.19740540106 (the compensated sums, via ALG-RED-002)
// Tags:    stream, vwap, volume-weighted-average-price, weighted-mean, execution, benchmark, trading, merge
// Units:   price in any price unit (the result is in the same unit); volume in any non-negative quantity unit (shares, contracts, lots)
type StreamingVWAP struct {
	pvSum, pvComp float64 // Σ price·volume = pvSum + pvComp
	vSum, vComp   float64 // Σ volume = vSum + vComp
	n             int64
}

// Add adds a trade of volume units at price. A zero-volume trade is counted and does not move the VWAP.
// O(1).
// Errors: ErrNotFinite when price or volume is NaN or ±Inf; ErrVolume when volume < 0; ErrOverflow
// when price·volume, a sum or the count would leave its range. On error nothing changes.
func (s *StreamingVWAP) Add(price, volume float64) error {
	if err := checkTrade(price, volume); err != nil {
		return err
	}
	if s.n == math.MaxInt64 {
		return fmt.Errorf("%w: count", ErrOverflow)
	}
	pv := price * volume
	pvSum, pvComp := reduce.NeumaierAdd(s.pvSum, s.pvComp, pv)
	vSum, vComp := reduce.NeumaierAdd(s.vSum, s.vComp, volume)
	if !reduce.AllFinite(pv, pvSum, pvSum+pvComp, vSum, vSum+vComp) {
		return fmt.Errorf("%w: price·volume = %v·%v", ErrOverflow, price, volume)
	}
	s.pvSum, s.pvComp, s.vSum, s.vComp = pvSum, pvComp, vSum, vComp
	s.n++
	return nil
}

// checkTrade returns ErrNotFinite or ErrVolume for an invalid trade, else nil.
func checkTrade(price, volume float64) error {
	if !reduce.AllFinite(price, volume) {
		return fmt.Errorf("%w: (price, volume) = (%v, %v)", ErrNotFinite, price, volume)
	}
	if volume < 0 {
		return fmt.Errorf("%w: volume = %v", ErrVolume, volume)
	}
	return nil
}

// Merge adds to s every trade other summarises. other is not changed. O(1).
// Errors: ErrOverflow when the count or a sum would leave its range. On error nothing changes.
func (s *StreamingVWAP) Merge(other StreamingVWAP) error {
	if other.n > math.MaxInt64-s.n {
		return fmt.Errorf("%w: count %d + %d", ErrOverflow, s.n, other.n)
	}
	pvSum, pvComp := reduce.NeumaierAdd(s.pvSum, s.pvComp, other.pvSum)
	pvSum, pvComp = reduce.NeumaierAdd(pvSum, pvComp, other.pvComp)
	vSum, vComp := reduce.NeumaierAdd(s.vSum, s.vComp, other.vSum)
	vSum, vComp = reduce.NeumaierAdd(vSum, vComp, other.vComp)
	if !reduce.AllFinite(pvSum, pvSum+pvComp, vSum, vSum+vComp) {
		return fmt.Errorf("%w: merged sums", ErrOverflow)
	}
	s.pvSum, s.pvComp, s.vSum, s.vComp = pvSum, pvComp, vSum, vComp
	s.n += other.n
	return nil
}

// VWAP returns Σ price·volume / Σ volume, or (0, false) while the total volume is 0. O(1).
func (s *StreamingVWAP) VWAP() (float64, bool) {
	return vwap(s.pvSum+s.pvComp, s.vSum+s.vComp)
}

// vwap returns pv/v, or false when v is not positive or the quotient is not finite.
func vwap(pv, v float64) (float64, bool) {
	if !(v > 0) {
		return 0, false
	}
	p := pv / v
	if math.IsInf(p, 0) || math.IsNaN(p) {
		return 0, false
	}
	return p, true
}

// TotalVolume returns Σ volume over the trades added. O(1).
func (s *StreamingVWAP) TotalVolume() float64 { return s.vSum + s.vComp }

// Count returns the number of trades added, zero-volume ones included. O(1).
func (s *StreamingVWAP) Count() int64 { return s.n }

// Reset empties the summary. O(1).
func (s *StreamingVWAP) Reset() { *s = StreamingVWAP{} }
