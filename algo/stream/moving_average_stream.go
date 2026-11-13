package stream

import (
	"errors"
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/queue"
	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Errors returned by this package. They are wrapped with the offending value; test with errors.Is.
var (
	// ErrWindowLen: a window length is less than 1 or above maxWindowLen.
	ErrWindowLen = errors.New("stream: window length must be in [1, 1<<24]")
	// ErrNotFinite: a value is NaN or ±Inf.
	ErrNotFinite = errors.New("stream: value is NaN or infinite")
	// ErrOverflow: a sum, a count or a variance term would leave the range of its type.
	ErrOverflow = errors.New("stream: result exceeds the range of its type")
	// ErrUninitialized: a method was called on a zero value that needs its constructor.
	ErrUninitialized = errors.New("stream: zero value; use the constructor")
)

// maxWindowLen caps a window at 1<<24 values (128 MiB of float64), the run's bound on what one library
// call allocates silently (lead decisions L13, L35).
const maxWindowLen = 1 << 24

// MovingAverageStream is the mean of the last windowLen values added (all of them while fewer have
// arrived): the single moving average, kept one value at a time. The values live in a
// queue.CircularBuffer (ALG-QUEUE-002) that overwrites the oldest; the window sum is compensated
// (Neumaier) and recomputed from the buffer every windowLen additions, so its error depends only on the
// values of the last two windows, never on the length of the stream. The zero value is not usable:
// its Add returns ErrUninitialized. Not safe for concurrent use.
//
// Methods: Add(value) adds a value and returns the new mean. Mean returns the current mean, or
// (0, false) when no value has arrived. Len is the number of values in the window, WindowLen its length.
//
// ID:      ALG-STREAM-001
// Time:    O(1) amortized per Add (O(w) on every w-th, the recomputation); Mean, Len O(1)   w = windowLen
// Space:   O(w), one allocation in NewMovingAverageStream; Add never allocates
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §6.4.2.1 "Single Moving Average", https://www.itl.nist.gov/div898/handbook/pmc/section4/pmc421.htm (read 2026-11-12)
// Source:  Neumaier, "Rundungsfehleranalyse einiger Verfahren zur Summation endlicher Summen", ZAMM, 1974, 54(1):39–51, doi:10.1002/zamm.19740540106
// Source:  Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.2 (a queue in a circular array), via ALG-QUEUE-002
// Tags:    stream, moving-average, sliding-window, rolling-mean, ring-buffer, compensated-summation, time-series
type MovingAverageStream struct {
	window       *queue.CircularBuffer[float64]
	sum, comp    float64 // the window total is sum + comp
	sinceRestart int     // Adds since sum was last recomputed from the buffer
}

// NewMovingAverageStream returns an empty MovingAverageStream over the last windowLen values.
// Errors: ErrWindowLen when windowLen < 1 or windowLen > 1<<24.
func NewMovingAverageStream(windowLen int) (*MovingAverageStream, error) {
	if windowLen < 1 || windowLen > maxWindowLen {
		return nil, fmt.Errorf("%w: windowLen = %d", ErrWindowLen, windowLen)
	}
	buf, err := queue.NewCircularBuffer[float64](windowLen, queue.Overwrite)
	if err != nil {
		return nil, fmt.Errorf("%w: windowLen = %d: %w", ErrWindowLen, windowLen, err)
	}
	return &MovingAverageStream{window: buf}, nil
}

// Add appends value to the window (dropping the oldest when full) and returns the mean of the window.
// Errors: ErrNotFinite when value is NaN or ±Inf; ErrOverflow when the window's sum would exceed the
// float64 range; ErrUninitialized on the zero value. On error nothing changes. O(1) amortized.
func (m *MovingAverageStream) Add(value float64) (float64, error) {
	if m.window == nil {
		return 0, ErrUninitialized
	}
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, fmt.Errorf("%w: value = %v", ErrNotFinite, value)
	}
	sum, comp := m.sum, m.comp
	if m.window.Len() == m.window.Cap() {
		leaving, _ := m.window.Front()
		sum, comp = reduce.NeumaierAdd(sum, comp, -leaving)
	}
	sum, comp = reduce.NeumaierAdd(sum, comp, value)
	if total := sum + comp; math.IsInf(sum, 0) || math.IsNaN(total) || math.IsInf(total, 0) {
		return 0, fmt.Errorf("%w: window sum with value = %v", ErrOverflow, value)
	}
	m.window.Push(value) // Overwrite policy: never fails on a constructed buffer
	m.sum, m.comp = sum, comp
	m.sinceRestart++
	if m.sinceRestart >= m.window.Cap() {
		m.recompute()
	}
	mean, _ := m.Mean()
	return mean, nil
}

// recompute replaces the running sum by a fresh compensated sum of the buffer, dropping rounding left by
// values that have since left the window. If a partial sum of the fresh pass overflows (possible when
// large values of both signs cancel), the running sum is kept. O(windowLen).
func (m *MovingAverageStream) recompute() {
	m.sinceRestart = 0
	var sum, comp float64
	for i := range m.window.Len() {
		v, _ := m.window.At(i)
		sum, comp = reduce.NeumaierAdd(sum, comp, v)
		if math.IsInf(sum, 0) {
			return
		}
	}
	m.sum, m.comp = sum, comp
}

// Mean returns the mean of the values in the window, or (0, false) when no value has been added. O(1).
func (m *MovingAverageStream) Mean() (float64, bool) {
	if m.window == nil || m.window.Len() == 0 {
		return 0, false
	}
	return (m.sum + m.comp) / float64(m.window.Len()), true
}

// Len returns the number of values in the window: min(values added, WindowLen()). O(1).
func (m *MovingAverageStream) Len() int {
	if m.window == nil {
		return 0
	}
	return m.window.Len()
}

// WindowLen returns the window length given to NewMovingAverageStream (0 for the zero value). O(1).
func (m *MovingAverageStream) WindowLen() int {
	if m.window == nil {
		return 0
	}
	return m.window.Cap()
}
