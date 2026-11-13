package geometry2d

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// uShape is a concave polygon: a 6×6 square with a 2-wide notch cut down from the top to y = 2.
var uShape = []Point{{0, 0}, {6, 0}, {6, 6}, {4, 6}, {4, 2}, {2, 2}, {2, 6}, {0, 6}}

// Worked by hand with the crossings test (Haines, Graphics Gems IV, 1994, pp. 24–46; O'Rourke,
// Computational Geometry in C, 2nd ed., §7.4.2) on the U shape. (1,2)'s ray runs along the notch's floor
// through vertices (2,2) and (4,2); the half-open rule counts the edges x = 2, x = 4 and x = 6: odd, inside.
func TestPointInPolygon_uShapeExample(t *testing.T) {
	tests := []struct {
		q    Point
		want Location
	}{
		{Point{1, 4}, Inside},
		{Point{3, 4}, Outside},
		{Point{3, 2}, OnBoundary},
		{Point{5, 6}, OnBoundary},
		{Point{1, 2}, Inside},
		{Point{5, 2}, Inside},
		{Point{3, 1}, Inside},
		{Point{7, 2}, Outside},
		{Point{-1, 2}, Outside},
	}
	for _, tt := range tests {
		got, err := PointInPolygon(uShape, tt.q)
		if err != nil || got != tt.want {
			t.Errorf("PointInPolygon(uShape, %v) = (%v, %v), want (%v, nil)", tt.q, got, err, tt.want)
		}
	}
}

func TestPointInPolygon_cases(t *testing.T) {
	const maxI, minI = math.MaxInt64, math.MinInt64
	square := []Point{{0, 0}, {4, 0}, {4, 4}, {0, 4}}
	clockwise := slices.Clone(square)
	slices.Reverse(clockwise)
	huge := []Point{{minI, minI}, {maxI, minI}, {maxI, maxI}, {minI, maxI}}
	diamond := []Point{{0, -2}, {2, 0}, {0, 2}, {-2, 0}}
	tests := []struct {
		name    string
		polygon []Point
		q       Point
		want    Location
	}{
		{"inside", square, Point{2, 2}, Inside},
		{"outside", square, Point{5, 2}, Outside},
		{"vertex", square, Point{0, 0}, OnBoundary},
		{"edge", square, Point{4, 1}, OnBoundary},
		{"top edge", square, Point{2, 4}, OnBoundary},
		{"outside, ray through two vertices", square, Point{-1, 4}, Outside},
		{"outside, ray through bottom edge's line", square, Point{-1, 0}, Outside},
		{"clockwise inside", clockwise, Point{1, 3}, Inside},
		{"clockwise on edge", clockwise, Point{0, 2}, OnBoundary},
		{"diamond, ray through the right vertex", diamond, Point{0, 0}, Inside},
		{"diamond, outside level with a vertex", diamond, Point{-3, 0}, Outside},
		{"diamond, beyond the top vertex", diamond, Point{0, 3}, Outside},
		{"repeated vertex", []Point{{0, 0}, {4, 0}, {4, 0}, {4, 4}, {0, 4}}, Point{3, 3}, Inside},
		{"extremes, inside", huge, Point{0, 0}, Inside},
		{"extremes, corner", huge, Point{maxI, maxI}, OnBoundary},
		{"extremes, edge", huge, Point{minI, 12345}, OnBoundary},
		{"degenerate triangle (a segment)", []Point{{0, 0}, {2, 2}, {4, 4}}, Point{1, 1}, OnBoundary},
		{"degenerate triangle, off it", []Point{{0, 0}, {2, 2}, {4, 4}}, Point{1, 2}, Outside},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PointInPolygon(tt.polygon, tt.q)
			if err != nil || got != tt.want {
				t.Errorf("PointInPolygon(%v, %v) = (%v, %v), want (%v, nil)", tt.polygon, tt.q, got, err, tt.want)
			}
		})
	}
}

func TestPointInPolygon_errors(t *testing.T) {
	for _, poly := range [][]Point{nil, {{0, 0}}, {{0, 0}, {1, 1}}} {
		if _, err := PointInPolygon(poly, Point{0, 0}); !errors.Is(err, ErrTooFewPoints) {
			t.Errorf("PointInPolygon(%v) error = %v, want ErrTooFewPoints", poly, err)
		}
	}
}

func TestLocation_String(t *testing.T) {
	for l, want := range map[Location]string{Outside: "Outside", Inside: "Inside", OnBoundary: "OnBoundary", 7: "Location(7)"} {
		if got := l.String(); got != want {
			t.Errorf("Location(%d).String() = %q, want %q", int(l), got, want)
		}
	}
}

// slowConvexLocation is the obvious test for a counterclockwise convex polygon: on an edge, or strictly
// left of every edge.
func slowConvexLocation(ccw []Point, q Point) Location {
	inside := true
	for i := range ccw {
		a, b := ccw[i], ccw[(i+1)%len(ccw)]
		o := Orientation(a, b, q)
		if o == 0 && inBox(a, b, q) {
			return OnBoundary
		}
		if o <= 0 {
			inside = false
		}
	}
	if inside {
		return Inside
	}
	return Outside
}

// Oracle: convex polygons (hulls of random points) in both orientations against the edge-side test, with
// query points on a small grid so that vertices and edges are hit often.
func TestPointInPolygon_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	trials := 0
	for trials < 300 {
		hull := ConvexHull(nil, randomPoints(rng, 3+rng.IntN(20), 8))
		if len(hull) < 3 {
			continue
		}
		trials++
		reversed := slices.Clone(hull)
		slices.Reverse(reversed)
		for range 30 {
			q := Point{rng.Int64N(21) - 10, rng.Int64N(21) - 10}
			want := slowConvexLocation(hull, q)
			for _, poly := range [][]Point{hull, reversed} {
				if got, err := PointInPolygon(poly, q); err != nil || got != want {
					t.Fatalf("PointInPolygon(%v, %v) = (%v, %v), want %v", poly, q, got, err, want)
				}
			}
		}
	}
}

func TestPointInPolygon_noAllocation(t *testing.T) {
	q := Point{1, 2}
	if allocs := testing.AllocsPerRun(100, func() { _, _ = PointInPolygon(uShape, q) }); allocs != 0 {
		t.Errorf("PointInPolygon allocated %v times per call, want 0", allocs)
	}
}

// The hull of random points has few vertices, so the benchmark polygon is a comb with about n vertices:
// teeth between y = 0 and y = 10, closed along y = -1. The query sits halfway up, crossing every tooth.
func BenchmarkPointInPolygon(b *testing.B) {
	for _, size := range benchSizes {
		teeth := size.n / 2
		poly := make([]Point, 0, 2*teeth+2)
		for k := range teeth {
			poly = append(poly, Point{int64(2 * k), 10}, Point{int64(2*k + 1), 0})
		}
		poly = append(poly, Point{int64(2 * teeth), -1}, Point{0, -1})
		q := Point{int64(teeth), 5}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = PointInPolygon(poly, q)
			}
		})
	}
}
