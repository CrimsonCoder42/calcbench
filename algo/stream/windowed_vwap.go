package stream

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/queue"
	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// trade is one (price, volume) held in a window.
type trade struct{ price, volume float64 }

// WindowedVWAP is the volume-weighted average price of the last windowLen trades added (all of them
// while fewer have arrived), Σ priceᵢ·volumeᵢ / Σ volumeᵢ over the window. Each Add puts the new trade
// into two compensated (Neumaier) sums and, once the window is full, subtracts the oldest trade's
// price·volume and volume (recomputed from the held trade, so the same float64 product leaves as
// entered), in O(1), never re-summing the window. The trades live in a queue.CircularBuffer
// (ALG-QUEUE-002). The zero value is not usable: its Add returns ErrUninitialized. Not safe for
// concurrent use. No Merge (windows have no meaningful merge).
//
// Drift: the compensated sums keep the error at about 2u·|sum| plus a second-order term u²·Σ|terms|
// over every trade ever added and removed (Neumaier 1974; Higham 2002 §4.3), u = 2⁻⁵³; there is no
// recompute inside Add, so Add stays O(1) worst case. Whenever the window holds no trade of positive
// volume, both sums are set to exactly 0, which removes the drift.
//
// Methods: Add(price, volume) adds a trade. VWAP returns the window's average, or (0, false) while its
// volume is 0. TotalVolume is the window's volume. Len is the number of trades held, WindowLen the
// window length. Reset empties the window.
//
// ID:      ALG-STREAM-011
// Time:    O(1) worst case per Add and per reader; Reset O(w)   w = windowLen
// Space:   O(w), allocated only in NewWindowedVWAP; Add never allocates
// Source:  Berkowitz, Logue, Noser, "The total cost of transactions on the NYSE", The Journal of Finance, 1988, 43(1):97–112, doi:10.1111/j.1540-6261.1988.tb02591.x (the volume-weighted average price benchmark)
// Source:  Neumaier, "Rundungsfehleranalyse einiger Verfahren zur Summation endlicher Summen", ZAMM, 1974, 54(1):39–51, doi:10.1002/zamm.19740540106 (the compensated sums, via ALG-RED-002)
// Tags:    stream, vwap, rolling-vwap, volume-weighted-average-price, sliding-window, weighted-mean, execution, trading, ring-buffer
// Units:   price in any price unit (the result is in the same unit); volume in any non-negative quantity unit (shares, contracts, lots)
type WindowedVWAP struct {
	window        *queue.CircularBuffer[trade]
	pvSum, pvComp float64
	vSum, vComp   float64
	positive      int // trades in the window with volume > 0
}

// NewWindowedVWAP returns an empty WindowedVWAP over the last windowLen trades.
// Errors: ErrWindowLen when windowLen < 1 or windowLen > 1<<24.
func NewWindowedVWAP(windowLen int) (*WindowedVWAP, error) {
	buf, err := newWindow[trade](windowLen)
	if err != nil {
		return nil, err
	}
	return &WindowedVWAP{window: buf}, nil
}

// Add appends a trade of volume units at price, dropping the oldest trade when the window is full.
// O(1), no loop.
// Errors: ErrNotFinite when price or volume is NaN or ±Inf; ErrVolume when volume < 0; ErrOverflow
// when price·volume or a sum would leave the float64 range; ErrUninitialized on the zero value. On
// error nothing changes.
func (w *WindowedVWAP) Add(price, volume float64) error {
	if w.window == nil {
		return ErrUninitialized
	}
	if err := checkTrade(price, volume); err != nil {
		return err
	}
	pv := price * volume
	pvSum, pvComp, vSum, vComp, positive := w.pvSum, w.pvComp, w.vSum, w.vComp, w.positive
	if w.window.Len() == w.window.Cap() {
		old, _ := w.window.Front()
		pvSum, pvComp = reduce.NeumaierAdd(pvSum, pvComp, -(old.price * old.volume))
		vSum, vComp = reduce.NeumaierAdd(vSum, vComp, -old.volume)
		if old.volume > 0 {
			positive--
		}
	}
	if volume > 0 {
		positive++
	}
	if positive == 0 {
		pvSum, pvComp, vSum, vComp = 0, 0, 0, 0
	} else {
		pvSum, pvComp = reduce.NeumaierAdd(pvSum, pvComp, pv)
		vSum, vComp = reduce.NeumaierAdd(vSum, vComp, volume)
	}
	if !reduce.AllFinite(pv, pvSum, pvSum+pvComp, vSum, vSum+vComp) {
		return fmt.Errorf("%w: price·volume = %v·%v", ErrOverflow, price, volume)
	}
	w.window.Push(trade{price, volume}) // Overwrite policy: never fails on a constructed buffer
	w.pvSum, w.pvComp, w.vSum, w.vComp, w.positive = pvSum, pvComp, vSum, vComp, positive
	return nil
}

// VWAP returns the volume-weighted average price of the window, or (0, false) while its volume is 0.
// O(1).
func (w *WindowedVWAP) VWAP() (float64, bool) {
	if w.positive == 0 {
		return 0, false
	}
	return vwap(w.pvSum+w.pvComp, w.vSum+w.vComp)
}

// TotalVolume returns Σ volume over the trades in the window. O(1).
func (w *WindowedVWAP) TotalVolume() float64 { return w.vSum + w.vComp }

// Reset empties the window and keeps its length. O(w).
func (w *WindowedVWAP) Reset() {
	if w.window == nil {
		return
	}
	clearWindow(w.window)
	w.pvSum, w.pvComp, w.vSum, w.vComp, w.positive = 0, 0, 0, 0, 0
}

// Len returns the number of trades in the window: min(trades added since Reset, WindowLen()). O(1).
func (w *WindowedVWAP) Len() int {
	if w.window == nil {
		return 0
	}
	return w.window.Len()
}

// WindowLen returns the window length given to NewWindowedVWAP (0 for the zero value). O(1).
func (w *WindowedVWAP) WindowLen() int {
	if w.window == nil {
		return 0
	}
	return w.window.Cap()
}
