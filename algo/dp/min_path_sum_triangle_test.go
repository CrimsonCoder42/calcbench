package dp

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// trianglePathSlow tries both children recursively (2^(r−1) paths; r ≤ 14).
func trianglePathSlow(tri [][]int64, i, j int) int64 {
	if i == len(tri)-1 {
		return tri[i][j]
	}
	return tri[i][j] + min(trianglePathSlow(tri, i+1, j), trianglePathSlow(tri, i+1, j+1))
}

// Worked by hand over all eight paths: 3→2→4→2 = 11 is the least.
func TestMinPathSumTriangle_example(t *testing.T) {
	tri := [][]int64{{3}, {6, 2}, {1, 9, 4}, {8, 5, 7, 2}}
	if got, err := MinPathSumTriangle(tri); err != nil || got != 11 {
		t.Errorf("MinPathSumTriangle = (%d, %v), want (11, nil)", got, err)
	}
}

func TestMinPathSumTriangle_cases(t *testing.T) {
	tests := []struct {
		name string
		tri  [][]int64
		want int64
	}{
		{"one row", [][]int64{{7}}, 7},
		{"two rows", [][]int64{{1}, {5, 2}}, 3},
		{"negative", [][]int64{{-1}, {-2, -3}}, -4},
		{"zeros", [][]int64{{0}, {0, 0}, {0, 0, 0}}, 0},
		{"extreme fits", [][]int64{{math.MaxInt64}, {0, 1}}, math.MaxInt64},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := MinPathSumTriangle(tt.tri); err != nil || got != tt.want {
				t.Errorf("MinPathSumTriangle(%v) = (%d, %v), want %d", tt.tri, got, err, tt.want)
			}
		})
	}
}

func TestMinPathSumTriangle_errors(t *testing.T) {
	tests := []struct {
		name string
		tri  [][]int64
		want error
	}{
		{"no rows", nil, ErrEmpty},
		{"first row too long", [][]int64{{1, 2}}, ErrNotTriangle},
		{"second row too short", [][]int64{{1}, {2}}, ErrNotTriangle},
		{"overflow", [][]int64{{math.MaxInt64}, {1, 1}}, ErrOverflow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := MinPathSumTriangle(tt.tri); !errors.Is(err, tt.want) {
				t.Errorf("MinPathSumTriangle error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestMinPathSumTriangle_doesNotModifyInput(t *testing.T) {
	tri := [][]int64{{3}, {6, 2}, {1, 9, 4}}
	if _, err := MinPathSumTriangle(tri); err != nil {
		t.Fatal(err)
	}
	if tri[2][0] != 1 || tri[2][1] != 9 || tri[2][2] != 4 || tri[1][0] != 6 {
		t.Errorf("input modified: %v", tri)
	}
}

func TestMinPathSumTriangle_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		tri := make([][]int64, 1+rng.IntN(12))
		for i := range tri {
			tri[i] = make([]int64, i+1)
			for j := range tri[i] {
				tri[i][j] = rng.Int64N(41) - 20
			}
		}
		if got, err := MinPathSumTriangle(tri); err != nil || got != trianglePathSlow(tri, 0, 0) {
			t.Fatalf("trial %d: got (%d, %v), want %d", trial, got, err, trianglePathSlow(tri, 0, 0))
		}
	}
}

// Triangles of about 10³, 10⁴ and 10⁵ entries (45, 141 and 447 rows).
func BenchmarkMinPathSumTriangle(b *testing.B) {
	for _, size := range []struct {
		name string
		rows int
	}{{"1e3", 45}, {"1e4", 141}, {"1e5", 447}} {
		tri := make([][]int64, size.rows)
		for i := range tri {
			tri[i] = benchInt64s(i+1, 100)
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = MinPathSumTriangle(tri)
			}
		})
	}
}
