package geometry2d

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Worked by hand on Andrew's monotone chain (Andrew 1979; CLRS 3rd ed. §33.3): a square with one point
// inside and one inside an edge. The hull is the four corners, counterclockwise from (0,0).
func TestConvexHull_squareExample(t *testing.T) {
	pts := []Point{{2, 2}, {4, 4}, {0, 0}, {2, 0}, {0, 4}, {4, 0}}
	want := []Point{{0, 0}, {4, 0}, {4, 4}, {0, 4}}
	if got := ConvexHull(nil, pts); !slices.Equal(got, want) {
		t.Errorf("ConvexHull(square) = %v, want %v", got, want)
	}
}

func TestConvexHull_cases(t *testing.T) {
	const maxI, minI = math.MaxInt64, math.MinInt64
	tests := []struct {
		name   string
		points []Point
		want   []Point
	}{
		{"empty", nil, nil},
		{"one", []Point{{5, -1}}, []Point{{5, -1}}},
		{"one, repeated", []Point{{5, -1}, {5, -1}, {5, -1}}, []Point{{5, -1}}},
		{"two", []Point{{3, 3}, {1, 1}}, []Point{{1, 1}, {3, 3}}},
		{"collinear keeps the extremes", []Point{{2, 2}, {0, 0}, {3, 3}, {1, 1}}, []Point{{0, 0}, {3, 3}}},
		{"vertical collinear", []Point{{0, 5}, {0, -2}, {0, 1}}, []Point{{0, -2}, {0, 5}}},
		{"triangle, clockwise input", []Point{{0, 0}, {0, 3}, {4, 0}}, []Point{{0, 0}, {4, 0}, {0, 3}}},
		{"duplicates of hull vertices", []Point{{0, 0}, {1, 0}, {0, 1}, {1, 0}, {0, 0}}, []Point{{0, 0}, {1, 0}, {0, 1}}},
		{"negative coordinates", []Point{{-3, -3}, {-1, -3}, {-2, -1}, {-2, -2}}, []Point{{-3, -3}, {-1, -3}, {-2, -1}}},
		{"extremes of int64", []Point{{minI, minI}, {maxI, minI}, {maxI, maxI}, {minI, maxI}, {0, 0}},
			[]Point{{minI, minI}, {maxI, minI}, {maxI, maxI}, {minI, maxI}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvexHull(nil, tt.points); !slices.Equal(got, tt.want) {
				t.Errorf("ConvexHull(%v) = %v, want %v", tt.points, got, tt.want)
			}
		})
	}
}

func TestConvexHull_appendsToDst(t *testing.T) {
	dst := []Point{{9, 9}}
	got := ConvexHull(dst, []Point{{0, 0}, {2, 0}, {1, 2}})
	if want := []Point{{9, 9}, {0, 0}, {2, 0}, {1, 2}}; !slices.Equal(got, want) {
		t.Errorf("ConvexHull(dst, triangle) = %v, want %v", got, want)
	}
	if got := ConvexHull(dst, nil); !slices.Equal(got, dst) {
		t.Errorf("ConvexHull(dst, nil) = %v, want dst unchanged", got)
	}
}

func TestConvexHull_doesNotModifyInput(t *testing.T) {
	pts := randomPoints(rand.New(rand.NewPCG(1, 2)), 50, 10)
	orig := slices.Clone(pts)
	_ = ConvexHull(nil, pts)
	if !slices.Equal(pts, orig) {
		t.Errorf("ConvexHull modified its input")
	}
}

// slowHullVertices is the obvious O(n³) version: a directed pair (p, q) of distinct points is a hull edge
// when every point is strictly left of p→q or on the closed segment pq; the vertices are the edges' ends.
func slowHullVertices(points []Point) []Point {
	var out []Point
	for _, p := range points {
		for _, q := range points {
			if p == q {
				continue
			}
			edge := true
			for _, r := range points {
				o := Orientation(p, q, r)
				if o < 0 || (o == 0 && !inBox(p, q, r)) {
					edge = false
					break
				}
			}
			if edge {
				out = append(out, p, q)
			}
		}
	}
	if len(out) == 0 && len(points) > 0 {
		out = append(out, points[0]) // every point equal
	}
	slices.SortFunc(out, comparePoints)
	return slices.Compact(out)
}

// Oracle: same vertex set as the O(n³) edge test, in strictly counterclockwise order from the least point.
// Small coordinates force many collinear and repeated points.
func TestConvexHull_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		pts := randomPoints(rng, rng.IntN(40), int64(1+rng.IntN(6)))
		got := ConvexHull(nil, pts)
		want := slowHullVertices(pts)
		sortedGot := slices.SortedFunc(slices.Values(got), comparePoints)
		if !slices.Equal(sortedGot, want) {
			t.Fatalf("trial %d: ConvexHull(%v) = %v, want vertices %v", trial, pts, got, want)
		}
		if len(got) > 0 && got[0] != want[0] {
			t.Fatalf("trial %d: hull starts at %v, want %v", trial, got[0], want[0])
		}
		if len(got) >= 3 {
			for i := range got {
				if Orientation(got[i], got[(i+1)%len(got)], got[(i+2)%len(got)]) != 1 {
					t.Fatalf("trial %d: hull %v is not strictly counterclockwise at %d", trial, got, i)
				}
			}
		}
	}
}

func BenchmarkConvexHull(b *testing.B) {
	for _, size := range benchSizes {
		pts := randomPoints(rand.New(rand.NewPCG(1, 2)), size.n, 1<<40)
		dst := make([]Point, 0, size.n+1)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				dst = ConvexHull(dst[:0], pts)
			}
		})
	}
}
