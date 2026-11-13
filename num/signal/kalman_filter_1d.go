package signal

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// KalmanFilter1D is the scalar Kalman filter for a random-walk state observed with noise:
//
//	x_k = x_{k−1} + w_k,  w_k ~ (0, Q)      (process model, transition 1)
//	z_k = x_k + v_k,      v_k ~ (0, R)      (measurement model, H = 1)
//
// It keeps the estimate x̂ and its error variance P. Predict projects ahead (x̂⁻ = x̂, P⁻ = P + Q); Update
// folds in a measurement z with the gain K = P⁻/(P⁻ + R): x̂ = x̂⁻ + K(z − x̂⁻), P = (1 − K)P⁻. A typical
// loop is Predict then Update for each new measurement. It is the minimum-variance linear estimator for
// these models (optimal, not only best linear, for Gaussian w and v). There is no Merge: two filters'
// posteriors are not combinable without their shared history. A filter is used by one goroutine at a time.
// Use NewKalmanFilter1D; the zero value has R = 0 and P = 0 and is not a valid filter.
//
// ID:      NUM-SIG-016
// Time:    O(1) per Predict or Update
// Space:   O(1), no allocation
// Source:  Kalman, "A New Approach to Linear Filtering and Prediction Problems", Transactions of the ASME, Journal of Basic Engineering, 1960, 82(1):35–45, doi:10.1115/1.3662552
// Source:  Brown & Hwang, Introduction to Random Signals and Applied Kalman Filtering, 4th ed., Wiley, 2012, ch. 4 "Discrete Kalman Filter Basics", §4.2 (the recursive loop: gain, estimate update, error covariance update, projection ahead), specialised to φ = 1, H = 1
// Tags:    kalman-filter, state-estimation, filtering, smoothing, random-walk, streaming, signal-processing
// Units:   state and measurements in one unit u; processVariance, measurementVariance and the variance in u²; Gain dimensionless in [0, 1]
// Example: x̂₀ = 0, P₀ = 1, Q = 0, R = 1; Predict, Update(2) → K = 0.5, x̂ = 1, P = 0.5 (hand-checked)
type KalmanFilter1D struct {
	state, variance                      float64 // x̂ and P
	processVariance, measurementVariance float64 // Q and R
	gain                                 float64 // K of the last Update
}

// NewKalmanFilter1D returns a filter with initial estimate initialState, its error variance
// initialVariance, process noise variance processVariance (Q) and measurement noise variance
// measurementVariance (R).
// Errors: ErrNotFinite when an argument is NaN or ±Inf; ErrOutOfRange when initialVariance < 0,
// processVariance < 0 or measurementVariance ≤ 0 (R = 0 makes the gain 1 and the filter a copy of the
// measurements; a zero prior with zero R would divide by zero).
func NewKalmanFilter1D(initialState, initialVariance, processVariance, measurementVariance float64) (*KalmanFilter1D, error) {
	if !reduce.AllFinite(initialState, initialVariance, processVariance, measurementVariance) {
		return nil, fmt.Errorf("%w: initialState = %v, initialVariance = %v, processVariance = %v, measurementVariance = %v",
			ErrNotFinite, initialState, initialVariance, processVariance, measurementVariance)
	}
	if processVariance < 0 {
		return nil, fmt.Errorf("%w: processVariance = %v, want ≥ 0", ErrOutOfRange, processVariance)
	}
	if measurementVariance <= 0 {
		return nil, fmt.Errorf("%w: measurementVariance = %v, want > 0", ErrOutOfRange, measurementVariance)
	}
	k := &KalmanFilter1D{processVariance: processVariance, measurementVariance: measurementVariance}
	if err := k.Reset(initialState, initialVariance); err != nil {
		return nil, err
	}
	return k, nil
}

// Predict projects the estimate one step ahead: the state is unchanged (random walk) and the variance
// grows by the process variance, P⁻ = P + Q.
// Errors: ErrOverflow when P + Q exceeds the float64 range (the filter is left unchanged).
func (k *KalmanFilter1D) Predict() error {
	p := k.variance + k.processVariance
	if !reduce.AllFinite(p) {
		return fmt.Errorf("%w: predicted variance", ErrOverflow)
	}
	k.variance = p
	return nil
}

// Update folds in one measurement: K = P⁻/(P⁻ + R), x̂ = x̂⁻ + K(z − x̂⁻), P = (1 − K)P⁻, the latter
// computed as P⁻R/(P⁻ + R), which cannot go negative by rounding.
// Errors: ErrNotFinite when measurement is NaN or ±Inf; ErrOverflow when the innovation or the new
// estimate exceeds the float64 range (the filter is left unchanged).
func (k *KalmanFilter1D) Update(measurement float64) error {
	if !reduce.AllFinite(measurement) {
		return fmt.Errorf("%w: measurement = %v", ErrNotFinite, measurement)
	}
	s := k.variance + k.measurementVariance
	if !reduce.AllFinite(s) {
		return fmt.Errorf("%w: innovation variance", ErrOverflow)
	}
	gain := k.variance / s
	state := k.state + gain*(measurement-k.state)
	variance := k.variance / s * k.measurementVariance
	if !reduce.AllFinite(state, variance) {
		return fmt.Errorf("%w: updated estimate", ErrOverflow)
	}
	k.state, k.variance, k.gain = state, variance, gain
	return nil
}

// State returns the current estimate x̂ (after an Update, the posterior; after a Predict, the prior).
func (k *KalmanFilter1D) State() float64 { return k.state }

// Variance returns the current error variance P of the estimate.
func (k *KalmanFilter1D) Variance() float64 { return k.variance }

// Gain returns the Kalman gain K of the last Update, 0 before the first.
func (k *KalmanFilter1D) Gain() float64 { return k.gain }

// Reset restarts the filter at estimate state with error variance variance, keeping Q and R.
// Errors: ErrNotFinite when state or variance is NaN or ±Inf; ErrOutOfRange when variance < 0 (the
// filter is left unchanged).
func (k *KalmanFilter1D) Reset(state, variance float64) error {
	if !reduce.AllFinite(state, variance) {
		return fmt.Errorf("%w: state = %v, variance = %v", ErrNotFinite, state, variance)
	}
	if variance < 0 {
		return fmt.Errorf("%w: variance = %v, want ≥ 0", ErrOutOfRange, variance)
	}
	k.state, k.variance, k.gain = state, variance, 0
	return nil
}
