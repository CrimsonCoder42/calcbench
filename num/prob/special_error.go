package prob

import "fmt"

// specialErr re-wraps an error from num/special so errors.Is matches both this package's ErrOutOfRange and
// the num/special sentinel (in practice special.ErrNoConvergence: the arguments were validated here first).
func specialErr(what string, err error) error {
	return fmt.Errorf("%w: %s: %w", ErrOutOfRange, what, err)
}
