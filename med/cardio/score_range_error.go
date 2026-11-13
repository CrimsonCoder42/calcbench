package cardio

import "fmt"

// errScore wraps ErrOutOfRange for a score outside lo–hi.
func errScore(score, lo, hi int) error {
	return fmt.Errorf("%w: score = %d, want %d–%d", ErrOutOfRange, score, lo, hi)
}
