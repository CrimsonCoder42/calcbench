package geometry2d

import (
	"errors"
	"fmt"
)

var (
	// ErrInvertedRect: a rectangle has MinX > MaxX or MinY > MaxY.
	ErrInvertedRect = errors.New("geometry2d: rectangle min exceeds max")
	// ErrOverflow: the result does not fit the return type.
	ErrOverflow = errors.New("geometry2d: result overflows")
)

// Rect is the closed axis-aligned rectangle [MinX, MaxX] × [MinY, MaxY]. A zero width or height is
// allowed (a segment or a point, area 0).
type Rect struct {
	MinX, MinY, MaxX, MaxY int64
}

// RectanglesUnionArea returns the area covered by r1 or r2 (their overlap counted once), in the square
// of the coordinate unit: area(r1) + area(r2) − area(r1 ∩ r2), computed exactly in 192 bits.
// Errors: ErrInvertedRect (wrapped with r1 or r2) when a rectangle has min > max on an axis;
// ErrOverflow when the area exceeds 2⁶⁴ − 1.
//
// ID:      ALG-GEO2D-007
// Time:    O(1)
// Space:   O(1) extra, no allocation
// Source:  inclusion–exclusion |A ∪ B| = |A| + |B| − |A ∩ B|: Stanley, Enumerative Combinatorics, Vol. 1, 2nd ed., Cambridge University Press, 2012, §2.1; the intersection of two axis-aligned rectangles is the product of the intersections of their x and y intervals (proof in doc/algorithms_core_2026-11-12/catalog/geometry2d.md, ALG-GEO2D-007)
// Source:  Klee, "Can the measure of ∪[aᵢ, bᵢ] be computed in less than O(n log n) steps?", American Mathematical Monthly 84(4):284–285, 1977, doi:10.1080/00029890.1977.11994336 (the n-rectangle generalisation)
// Tags:    geometry, rectangle, union-area, inclusion-exclusion, overlap
func RectanglesUnionArea(r1, r2 Rect) (uint64, error) {
	if err := r1.validate(); err != nil {
		return 0, fmt.Errorf("%w: r1 = %+v", err, r1)
	}
	if err := r2.validate(); err != nil {
		return 0, fmt.Errorf("%w: r2 = %+v", err, r2)
	}
	total := widen128(r1.area()).add(widen128(r2.area()))
	overlap := Rect{max(r1.MinX, r2.MinX), max(r1.MinY, r2.MinY), min(r1.MaxX, r2.MaxX), min(r1.MaxY, r2.MaxY)}
	if overlap.validate() == nil {
		// The overlap is inside both rectangles, so subtracting it cannot go below zero.
		total = total.sub(widen128(overlap.area()))
	}
	if total.top != 0 || total.hi != 0 {
		return 0, fmt.Errorf("%w: union area of %+v and %+v exceeds 2^64-1", ErrOverflow, r1, r2)
	}
	return total.lo, nil
}

func (r Rect) validate() error {
	if r.MinX > r.MaxX || r.MinY > r.MaxY {
		return ErrInvertedRect
	}
	return nil
}

// area returns (MaxX − MinX)·(MaxY − MinY) exactly; r must be valid.
func (r Rect) area() uint128 {
	w, _ := coordDiff(r.MaxX, r.MinX)
	h, _ := coordDiff(r.MaxY, r.MinY)
	return mul64(w, h)
}
