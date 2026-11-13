package shopping

import (
	"errors"
	"fmt"
	"math"
)

var (
	// ErrOutOfRange is returned when an amount or rate is negative.
	ErrOutOfRange = errors.New("shopping: input outside its range")
	// ErrOverflow is returned when an intermediate product does not fit in int64.
	ErrOverflow = errors.New("shopping: result does not fit in int64")
)

// TipCents returns the tip on a bill: billCents × tipRateBasisPoints / 10 000, rounded half to even to the
// whole cent (the house rounding rule; the source defines only the percentage). 1 basis point = 0.01 %, so
// 15 % is 1500. How a tip is split among people, and any remainder, is the caller's.
// Errors: ErrOutOfRange when billCents < 0 or tipRateBasisPoints < 0; ErrOverflow when billCents ×
// tipRateBasisPoints exceeds int64.
//
// ID:      LIF-SHOP-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Thompson A, Taylor BN. Guide for the Use of the International System of Units (SI). NIST Special Publication 811, 2008 ed., §7.10.2 "Percent": "the symbol % ... is used to represent the number 0.01"
// Tags:    shopping, tip, gratuity, percent, money
// Units:   billCents cents; tipRateBasisPoints basis points (0.01 %); result cents
// Range:   billCents ≥ 0; tipRateBasisPoints ≥ 0
// Example: $55.00 at 15 % → 5500 × 1500 / 10000 = 825 cents; $10.10 at 15 % → 151.5 → 152 cents (half to even); $10.30 at 15 % → 154.5 → 154 (HAND-DERIVED)
func TipCents(billCents, tipRateBasisPoints int64) (int64, error) {
	if billCents < 0 {
		return 0, fmt.Errorf("%w: billCents = %d, want ≥ 0", ErrOutOfRange, billCents)
	}
	if tipRateBasisPoints < 0 {
		return 0, fmt.Errorf("%w: tipRateBasisPoints = %d, want ≥ 0", ErrOutOfRange, tipRateBasisPoints)
	}
	if tipRateBasisPoints != 0 && billCents > math.MaxInt64/tipRateBasisPoints {
		return 0, fmt.Errorf("%w: billCents × tipRateBasisPoints = %d × %d", ErrOverflow, billCents, tipRateBasisPoints)
	}
	p := billCents * tipRateBasisPoints
	q, r := p/10000, p%10000
	if r > 5000 || (r == 5000 && q%2 != 0) {
		q++
	}
	return q, nil
}
