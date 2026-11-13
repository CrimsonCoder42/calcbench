package logs

import (
	"errors"
	"fmt"
	"math"
)

var (
	// ErrNotFinite reports a NaN or +Inf input (−Inf is allowed: it is the logarithm of zero).
	ErrNotFinite = errors.New("logs: value is NaN or +Inf")
	// ErrEmpty reports an empty slice where at least one value is needed.
	ErrEmpty = errors.New("logs: empty input")
)

// checkLogValue returns ErrNotFinite, naming the parameter, when x is NaN or +Inf. −Inf passes: it is
// log 0, a valid log-domain value.
func checkLogValue(name string, x float64) error {
	if math.IsNaN(x) || math.IsInf(x, 1) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, x)
	}
	return nil
}

// LogAddExp returns log(eᵃ + eᵇ) without overflow or underflow: max(a, b) + log1p(e^−|a−b|). It adds two
// probabilities (or any positive quantities) held as natural logarithms. −Inf stands for log 0, so
// LogAddExp(−Inf, b) = b and LogAddExp(−Inf, −Inf) = −Inf.
// Errors: ErrNotFinite when a or b is NaN or +Inf.
//
// ID:      NUM-LOGS-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Blanchard, Higham, Higham, "Accurately computing the log-sum-exp and softmax functions", IMA Journal of Numerical Analysis, 2021, 41(4):2311–2330, doi:10.1093/imanum/draa038 (§2, eq. 2.3 shifted form with log1p, Algorithm 2 for n = 2)
// Tags:    logarithm, log-sum-exp, logaddexp, log-domain, probability, numerically-stable
// Units:   a, b natural logarithms of a positive quantity; result the natural logarithm of their sum
// Example: LogAddExp(log 1, log 3) = log 4 = 1.3862943611198906; LogAddExp(1000, 1000) = 1000 + log 2 (hand-checked)
func LogAddExp(a, b float64) (float64, error) {
	if err := checkLogValue("a", a); err != nil {
		return 0, err
	}
	if err := checkLogValue("b", b); err != nil {
		return 0, err
	}
	hi, lo := max(a, b), min(a, b)
	if math.IsInf(lo, -1) {
		return hi, nil // e^lo = 0 adds nothing; also covers both −Inf
	}
	return hi + math.Log1p(math.Exp(lo-hi)), nil
}

// LogAddExp2 returns log₂(2ᵃ + 2ᵇ) without overflow or underflow: max(a, b) + log1p(2^−|a−b|)/ln 2. It is
// LogAddExp in base 2 (bits instead of nats). −Inf stands for log₂ 0.
// Errors: ErrNotFinite when a or b is NaN or +Inf.
//
// ID:      NUM-LOGS-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Blanchard, Higham, Higham, "Accurately computing the log-sum-exp and softmax functions", IMA Journal of Numerical Analysis, 2021, 41(4):2311–2330, doi:10.1093/imanum/draa038 (§2, eq. 2.3; base change log₂ y = ln y/ln 2)
// Tags:    logarithm, log-sum-exp, logaddexp2, log-domain, base-2, bits, numerically-stable
// Units:   a, b base-2 logarithms of a positive quantity; result the base-2 logarithm of their sum
// Example: LogAddExp2(3, 3) = 4; LogAddExp2(0, 1) = log₂ 3 = 1.584962500721156 (hand-checked)
func LogAddExp2(a, b float64) (float64, error) {
	if err := checkLogValue("a", a); err != nil {
		return 0, err
	}
	if err := checkLogValue("b", b); err != nil {
		return 0, err
	}
	hi, lo := max(a, b), min(a, b)
	if math.IsInf(lo, -1) {
		return hi, nil
	}
	return hi + math.Log1p(math.Exp2(lo-hi))/math.Ln2, nil
}
