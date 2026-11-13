package computing

import (
	"errors"
	"fmt"
	"math"
)

var (
	// ErrNotFinite reports a NaN or infinite input.
	ErrNotFinite = errors.New("computing: value is NaN or infinite")
	// ErrNegative reports a negative input.
	ErrNegative = errors.New("computing: value is negative")
	// ErrZero reports a zero rate.
	ErrZero = errors.New("computing: rate is zero")
	// ErrOverflow reports a result beyond float64.
	ErrOverflow = errors.New("computing: result overflows float64")
)

// TransmissionTime returns the time to push sizeBits onto a link of rateBitsPerSecond, the transmission delay
// L / R. Propagation, queueing and protocol overhead are not included.
// Errors: ErrNotFinite; ErrNegative for a negative size or rate; ErrZero for rate 0; ErrOverflow.
//
// ID:      ENG-COMPUT-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Kurose JF, Ross KW. Computer Networking: A Top-Down Approach, 7th ed., Pearson, 2017, §1.4.1 "Overview of Delay in Packet-Switched Networks", transmission delay L/R (L bits, R bits/s)
// Tags:    bandwidth, transmission-delay, download-time, data-rate, network, computing
// Units:   sizeBits bit; rateBitsPerSecond bit/s; result s
// Example: 500 MB = 4 × 10⁹ bit at 5 Mbit/s → 800 s (13 min 20 s)
func TransmissionTime(sizeBits, rateBitsPerSecond float64) (float64, error) {
	for _, p := range [2]struct {
		name string
		v    float64
	}{{"sizeBits", sizeBits}, {"rateBitsPerSecond", rateBitsPerSecond}} {
		if math.IsNaN(p.v) || math.IsInf(p.v, 0) {
			return 0, fmt.Errorf("%w: %s = %v", ErrNotFinite, p.name, p.v)
		}
		if p.v < 0 {
			return 0, fmt.Errorf("%w: %s = %v", ErrNegative, p.name, p.v)
		}
	}
	if rateBitsPerSecond == 0 {
		return 0, fmt.Errorf("%w: rateBitsPerSecond = 0", ErrZero)
	}
	t := sizeBits / rateBitsPerSecond
	if math.IsInf(t, 0) {
		return 0, fmt.Errorf("%w: sizeBits / rateBitsPerSecond", ErrOverflow)
	}
	return t, nil
}
