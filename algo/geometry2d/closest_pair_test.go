package geometry2d

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"slices"
	"testing"
)

// bigSquaredDistance is |p − q|² in arbitrary precision.
func bigSquaredDistance(p, q Point) *big.Int {
	dx := new(big.Int).Sub(big.NewInt(p.X), big.NewInt(q.X))
	dy := new(big.Int).Sub(big.NewInt(p.Y), big.NewInt(q.Y))
	return new(big.Int).Add(new(big.Int).Mul(dx, dx), new(big.Int).Mul(dy, dy))
}

// slowClosest is the obvious O(n²) version: the least squared distance over all pairs.
func slowClosest(points []Point) *big.Int {
	var best *big.Int
	for i := range points {
		for j := i + 1; j < len(points); j++ {
			if d := bigSquaredDistance(points[i], points[j]); best == nil || d.Cmp(best) < 0 {
				best = d
			}
		}
	}
	return best
}

// Worked by hand for the divide and conquer of Shamos–Hoey 1975 / CLRS 3rd ed. §33.4: the closest pair
// is (10,10)–(11,12), at squared distance 5.
func TestClosestPair_example(t *testing.T) {
	pts := []Point{{0, 0}, {10, 10}, {3, 4}, {11, 12}, {20, 0}}
	i, j, err := ClosestPair(pts)
	if err != nil || i != 1 || j != 3 {
		t.Errorf("ClosestPair(example) = (%d, %d, %v), want (1, 3, nil)", i, j, err)
	}
}

func TestClosestPair_cases(t *testing.T) {
	const maxI, minI = math.MaxInt64, math.MinInt64
	tests := []struct {
		name   string
		points []Point
		want   [2]int
	}{
		{"two", []Point{{1, 1}, {4, 5}}, [2]int{0, 1}},
		{"three, last pair", []Point{{0, 0}, {10, 0}, {11, 0}}, [2]int{1, 2}},
		{"duplicate points", []Point{{0, 0}, {5, 5}, {9, 9}, {5, 5}}, [2]int{1, 3}},
		{"split across the dividing line", []Point{{0, 0}, {1, 100}, {2, 200}, {49, 50}, {51, 50}, {100, 0}, {101, 100}, {102, 200}}, [2]int{3, 4}},
		{"same x", []Point{{7, 0}, {7, 10}, {7, 3}, {7, 30}}, [2]int{0, 2}},
		{"negative", []Point{{-10, -10}, {-3, -3}, {-4, -3}, {5, 5}}, [2]int{1, 2}},
		{"extremes of int64", []Point{{minI, minI}, {maxI, 0}, {maxI, minI}}, [2]int{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i, j, err := ClosestPair(tt.points)
			if err != nil || [2]int{i, j} != tt.want {
				t.Errorf("ClosestPair(%v) = (%d, %d, %v), want (%d, %d, nil)", tt.points, i, j, err, tt.want[0], tt.want[1])
			}
		})
	}
}

func TestClosestPair_errors(t *testing.T) {
	for _, pts := range [][]Point{nil, {}, {{1, 1}}} {
		if _, _, err := ClosestPair(pts); !errors.Is(err, ErrTooFewPoints) {
			t.Errorf("ClosestPair(%v) error = %v, want ErrTooFewPoints", pts, err)
		}
	}
}

func TestClosestPair_doesNotModifyInput(t *testing.T) {
	pts := randomPoints(rand.New(rand.NewPCG(1, 2)), 100, 1000)
	orig := slices.Clone(pts)
	if _, _, err := ClosestPair(pts); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(pts, orig) {
		t.Errorf("ClosestPair modified its input")
	}
}

// Oracle: the returned pair's distance equals the O(n²) minimum, on small grids (many ties and
// duplicates) and on the full int64 range.
func TestClosestPair_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		n := 2 + rng.IntN(150)
		var pts []Point
		switch trial % 3 {
		case 0:
			pts = randomPoints(rng, n, int64(1+rng.IntN(20)))
		case 1:
			pts = randomPoints(rng, n, 1<<40)
		default:
			pts = make([]Point, n)
			for k := range pts {
				pts[k] = Point{randomCoord(rng), randomCoord(rng)}
			}
		}
		i, j, err := ClosestPair(pts)
		if err != nil || i < 0 || j >= n || i >= j {
			t.Fatalf("trial %d: ClosestPair = (%d, %d, %v), want valid indices i < j", trial, i, j, err)
		}
		if got, want := bigSquaredDistance(pts[i], pts[j]), slowClosest(pts); got.Cmp(want) != 0 {
			t.Fatalf("trial %d: ClosestPair(%v) gives squared distance %v, want %v", trial, pts, got, want)
		}
	}
}

func BenchmarkClosestPair(b *testing.B) {
	for _, size := range benchSizes {
		pts := randomPoints(rand.New(rand.NewPCG(1, 2)), size.n, 1<<40)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _, _ = ClosestPair(pts)
			}
		})
	}
}
