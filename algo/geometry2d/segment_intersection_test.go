package geometry2d

import (
	"math"
	"math/rand/v2"
	"testing"
)

// Worked by hand for SEGMENTS-INTERSECT (CLRS 3rd ed. §33.1): a proper crossing, a gap on one line, a
// collinear overlap, and two segments touching at an endpoint.
func TestSegmentsIntersect_clrsExample(t *testing.T) {
	tests := []struct {
		s, t Segment
		want bool
	}{
		{Segment{Point{0, 0}, Point{4, 4}}, Segment{Point{0, 4}, Point{4, 0}}, true},
		{Segment{Point{0, 0}, Point{2, 0}}, Segment{Point{3, 0}, Point{5, 0}}, false},
		{Segment{Point{0, 0}, Point{3, 0}}, Segment{Point{2, 0}, Point{5, 0}}, true},
		{Segment{Point{0, 0}, Point{2, 2}}, Segment{Point{2, 2}, Point{4, 0}}, true},
	}
	for _, tt := range tests {
		if got := SegmentsIntersect(tt.s, tt.t); got != tt.want {
			t.Errorf("SegmentsIntersect(%v, %v) = %v, want %v", tt.s, tt.t, got, tt.want)
		}
	}
}

func TestSegmentsIntersect_cases(t *testing.T) {
	const maxI, minI = math.MaxInt64, math.MinInt64
	seg := func(ax, ay, bx, by int64) Segment { return Segment{Point{ax, ay}, Point{bx, by}} }
	tests := []struct {
		name string
		s, t Segment
		want bool
	}{
		{"parallel, apart", seg(0, 0, 4, 0), seg(0, 1, 4, 1), false},
		{"T junction", seg(0, 0, 4, 0), seg(2, 0, 2, 5), true},
		{"near miss of a T", seg(0, 0, 4, 0), seg(2, 1, 2, 5), false},
		{"lines cross outside both segments", seg(0, 0, 1, 1), seg(3, 0, 2, 1), false},
		{"collinear, one inside the other", seg(0, 0, 10, 10), seg(3, 3, 4, 4), true},
		{"collinear, touching ends", seg(0, 0, 1, 1), seg(1, 1, 2, 2), true},
		{"identical", seg(1, 2, 3, 4), seg(1, 2, 3, 4), true},
		{"reversed copy", seg(1, 2, 3, 4), seg(3, 4, 1, 2), true},
		{"point on segment", seg(2, 2, 2, 2), seg(0, 0, 4, 4), true},
		{"point off segment, on its line", seg(5, 5, 5, 5), seg(0, 0, 4, 4), false},
		{"two equal points", seg(1, 1, 1, 1), seg(1, 1, 1, 1), true},
		{"two different points", seg(1, 1, 1, 1), seg(1, 2, 1, 2), false},
		{"negative coordinates", seg(-5, -1, -1, -5), seg(-5, -5, -1, -1), true},
		{"extremes, crossing diagonals", seg(minI, minI, maxI, maxI), seg(minI, maxI, maxI, minI), true},
		{"extremes, parallel one apart", seg(minI, minI, maxI, maxI), seg(minI+1, minI, maxI, maxI-1), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SegmentsIntersect(tt.s, tt.t); got != tt.want {
				t.Errorf("SegmentsIntersect(%v, %v) = %v, want %v", tt.s, tt.t, got, tt.want)
			}
			if got := SegmentsIntersect(tt.t, tt.s); got != tt.want {
				t.Errorf("SegmentsIntersect(%v, %v) (swapped) = %v, want %v", tt.t, tt.s, got, tt.want)
			}
		})
	}
}

// slowSegmentsIntersect is the obvious version for small coordinates: solve s.A + λ(s.B − s.A) =
// t.A + μ(t.B − t.A) by Cramer's rule in int64 and check λ, μ ∈ [0, 1]; when the directions are
// parallel (or a segment is a point) the segments meet exactly when an endpoint of one lies on the other.
func slowSegmentsIntersect(s, t Segment) bool {
	dsx, dsy := s.B.X-s.A.X, s.B.Y-s.A.Y
	dtx, dty := t.B.X-t.A.X, t.B.Y-t.A.Y
	wx, wy := t.A.X-s.A.X, t.A.Y-s.A.Y
	den := dsx*dty - dsy*dtx
	if den != 0 {
		lam := wx*dty - wy*dtx
		mu := wx*dsy - wy*dsx
		if den < 0 {
			den, lam, mu = -den, -lam, -mu
		}
		return 0 <= lam && lam <= den && 0 <= mu && mu <= den
	}
	on := func(p Point, seg Segment) bool {
		cross := (seg.B.X-seg.A.X)*(p.Y-seg.A.Y) - (seg.B.Y-seg.A.Y)*(p.X-seg.A.X)
		return cross == 0 && min(seg.A.X, seg.B.X) <= p.X && p.X <= max(seg.A.X, seg.B.X) &&
			min(seg.A.Y, seg.B.Y) <= p.Y && p.Y <= max(seg.A.Y, seg.B.Y)
	}
	return on(s.A, t) || on(s.B, t) || on(t.A, s) || on(t.B, s)
}

// Oracle: Cramer's rule on a 7×7 grid (many collinear, touching and degenerate cases).
func TestSegmentsIntersect_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	c := func() int64 { return rng.Int64N(7) - 3 }
	for trial := range 20000 {
		s := Segment{Point{c(), c()}, Point{c(), c()}}
		u := Segment{Point{c(), c()}, Point{c(), c()}}
		if got, want := SegmentsIntersect(s, u), slowSegmentsIntersect(s, u); got != want {
			t.Fatalf("trial %d: SegmentsIntersect(%v, %v) = %v, want %v", trial, s, u, got, want)
		}
	}
}

func TestSegmentsIntersect_noAllocation(t *testing.T) {
	s, u := Segment{Point{0, 0}, Point{4, 4}}, Segment{Point{0, 4}, Point{4, 0}}
	if allocs := testing.AllocsPerRun(100, func() { _ = SegmentsIntersect(s, u) }); allocs != 0 {
		t.Errorf("SegmentsIntersect allocated %v times per call, want 0", allocs)
	}
}

func BenchmarkSegmentsIntersect(b *testing.B) {
	pts := randomPoints(rand.New(rand.NewPCG(1, 2)), 4, math.MaxInt64/2)
	s, u := Segment{pts[0], pts[1]}, Segment{pts[2], pts[3]}
	b.Run("1", func(b *testing.B) {
		b.ReportAllocs()
		hits := 0
		for b.Loop() {
			if SegmentsIntersect(s, u) {
				hits++
			}
		}
		_ = hits
	})
}
