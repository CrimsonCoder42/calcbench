package prefix

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// grid3x4 is a hand-made grid for the rectangle identity of Crow (1984) and Viola–Jones (2004) §2.1:
// sum = S(r1,c1) − S(r0,c1) − S(r1,c0) + S(r0,c0). Expected sums below are hand-derived.
var grid3x4 = [][]int{
	{1, 2, 3, 4},
	{5, 6, 7, 8},
	{9, 10, 11, 12},
}

func TestPrefixSums2D_worked(t *testing.T) {
	table, err := PrefixSums2D(grid3x4)
	if err != nil {
		t.Fatal(err)
	}
	if table.Rows() != 3 || table.Cols() != 4 {
		t.Fatalf("table is %d × %d, want 3 × 4", table.Rows(), table.Cols())
	}
	// The full table, row 0 and column 0 zero.
	want := []int{
		0, 0, 0, 0, 0,
		0, 1, 3, 6, 10,
		0, 6, 14, 24, 36,
		0, 15, 33, 54, 78,
	}
	for i, w := range want {
		if table.sums[i] != w {
			t.Fatalf("sums[%d] = %d, want %d (table %v)", i, table.sums[i], w, table.sums)
		}
	}
}

func TestRectangleSum_worked(t *testing.T) {
	table, _ := PrefixSums2D(grid3x4)
	tests := []struct {
		name                   string
		row0, col0, row1, col1 int
		want                   int
	}{
		{"whole grid", 0, 0, 3, 4, 78},
		{"centre 2×2", 1, 1, 3, 3, 6 + 7 + 10 + 11},
		{"one cell", 2, 3, 3, 4, 12},
		{"one row", 1, 0, 2, 4, 26},
		{"one column", 0, 2, 3, 3, 21},
		{"empty rows", 2, 0, 2, 4, 0},
		{"empty columns", 0, 4, 3, 4, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := table.RectangleSum(tt.row0, tt.col0, tt.row1, tt.col1)
			if err != nil || got != tt.want {
				t.Errorf("RectangleSum(%d,%d,%d,%d) = (%d, %v), want (%d, nil)", tt.row0, tt.col0, tt.row1, tt.col1, got, err, tt.want)
			}
		})
	}
}

func TestPrefixSums2D_cases(t *testing.T) {
	for _, grid := range [][][]int{nil, {}, {{}, {}}} {
		table, err := PrefixSums2D(grid)
		if err != nil {
			t.Fatalf("PrefixSums2D(%v) error %v", grid, err)
		}
		if got, err := table.RectangleSum(0, 0, table.Rows(), 0); err != nil || got != 0 {
			t.Errorf("empty grid %v: RectangleSum = (%d, %v), want (0, nil)", grid, got, err)
		}
	}
	one, _ := PrefixSums2D([][]int{{-4}})
	if got, err := one.RectangleSum(0, 0, 1, 1); err != nil || got != -4 {
		t.Errorf("1×1 grid: RectangleSum = (%d, %v), want (-4, nil)", got, err)
	}
	ext, err := PrefixSums2D([][]int64{{math.MaxInt64, math.MinInt64}, {math.MinInt64, math.MaxInt64}})
	if err != nil {
		t.Fatal(err)
	}
	if got, err := ext.RectangleSum(0, 0, 2, 2); err != nil || got != -2 {
		t.Errorf("extremes: RectangleSum = (%d, %v), want (-2, nil)", got, err)
	}
}

func TestSummedAreaTable_zeroValue(t *testing.T) {
	var table SummedAreaTable[int]
	if got, err := table.RectangleSum(0, 0, 0, 0); err != nil || got != 0 {
		t.Errorf("zero table RectangleSum(0,0,0,0) = (%d, %v), want (0, nil)", got, err)
	}
	if _, err := table.RectangleSum(0, 0, 1, 1); !errors.Is(err, ErrRange) {
		t.Errorf("zero table RectangleSum(0,0,1,1) error = %v, want ErrRange", err)
	}
}

func TestPrefixSums2D_errors(t *testing.T) {
	if _, err := PrefixSums2D([][]int{{1, 2}, {3}}); !errors.Is(err, ErrRagged) {
		t.Errorf("ragged grid error = %v, want ErrRagged", err)
	}
	if _, err := PrefixSums2D([][]int8{{100, 100}}); !errors.Is(err, ErrOverflow) {
		t.Errorf("row overflow error = %v, want ErrOverflow", err)
	}
	if _, err := PrefixSums2D([][]int8{{100}, {100}}); !errors.Is(err, ErrOverflow) {
		t.Errorf("column overflow error = %v, want ErrOverflow", err)
	}
}

func TestRectangleSum_errors(t *testing.T) {
	table, _ := PrefixSums2D(grid3x4)
	for _, q := range [][4]int{{-1, 0, 1, 1}, {0, -1, 1, 1}, {2, 0, 1, 1}, {0, 2, 1, 1}, {0, 0, 4, 1}, {0, 0, 1, 5}} {
		if _, err := table.RectangleSum(q[0], q[1], q[2], q[3]); !errors.Is(err, ErrRange) {
			t.Errorf("RectangleSum%v error = %v, want ErrRange", q, err)
		}
	}
	// Every origin rectangle fits int8 (-100, 0, 100), but grid[0:1][1:3] = 100 + 100 does not.
	t8, err := PrefixSums2D([][]int8{{-100, 100, 100}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := t8.RectangleSum(0, 1, 1, 3); !errors.Is(err, ErrOverflow) {
		t.Errorf("RectangleSum over 254 in int8 error = %v, want ErrOverflow", err)
	}
}

func TestRectangleSum_noAllocation(t *testing.T) {
	table, _ := PrefixSums2D(grid3x4)
	allocs := testing.AllocsPerRun(20, func() { _, _ = table.RectangleSum(1, 1, 3, 3) })
	if allocs != 0 {
		t.Errorf("RectangleSum allocates %v times per call, want 0", allocs)
	}
}

// Oracle: random int8 grids against direct int sums of every rectangle. Building fails exactly when
// some origin rectangle or row prefix leaves int8; a query fails exactly when its own sum does.
func TestRectangleSum_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	inRange := func(s int) bool { return s >= math.MinInt8 && s <= math.MaxInt8 }
	for trial := range 300 {
		rows, cols := rng.IntN(4), 1+rng.IntN(4)
		grid := make([][]int8, rows)
		for r := range grid {
			grid[r] = randomInt8s(rng, cols)
			for c := range grid[r] {
				grid[r][c] /= 4 // keep building mostly successful
			}
		}
		rect := func(r0, c0, r1, c1 int) int {
			s := 0
			for r := r0; r < r1; r++ {
				for c := c0; c < c1; c++ {
					s += int(grid[r][c])
				}
			}
			return s
		}
		buildOK := true
		for r := 0; r < rows; r++ {
			for c := 1; c <= cols; c++ {
				if !inRange(rect(0, 0, r+1, c)) || !inRange(rect(r, 0, r+1, c)) {
					buildOK = false
				}
			}
		}
		table, err := PrefixSums2D(grid)
		if !buildOK {
			if !errors.Is(err, ErrOverflow) {
				t.Fatalf("trial %d: PrefixSums2D(%v) error = %v, want ErrOverflow", trial, grid, err)
			}
			continue
		}
		if err != nil {
			t.Fatalf("trial %d: PrefixSums2D(%v) error = %v", trial, grid, err)
		}
		// A grid with no rows has no columns either (documented: Cols() = 0), so query only what the
		// table holds.
		cols = table.Cols()
		for r0 := 0; r0 <= rows; r0++ {
			for r1 := r0; r1 <= rows; r1++ {
				for c0 := 0; c0 <= cols; c0++ {
					for c1 := c0; c1 <= cols; c1++ {
						want := rect(r0, c0, r1, c1)
						got, err := table.RectangleSum(r0, c0, r1, c1)
						if !inRange(want) {
							if !errors.Is(err, ErrOverflow) {
								t.Fatalf("trial %d: rect %v: got (%d, %v), want ErrOverflow", trial, [4]int{r0, c0, r1, c1}, got, err)
							}
						} else if err != nil || int(got) != want {
							t.Fatalf("trial %d: rect %v of %v: got (%d, %v), want %d", trial, [4]int{r0, c0, r1, c1}, grid, got, err, want)
						}
					}
				}
			}
		}
	}
}

// benchGrid returns a fixed-seed side × side grid of values in [-1000, 1000).
func benchGrid(side int) [][]int64 {
	flat := benchInts(side * side)
	grid := make([][]int64, side)
	for r := range grid {
		grid[r] = flat[r*side : (r+1)*side]
	}
	return grid
}

// Sizes are cell counts ≈ 1e3, 1e4, 1e5 (sides 32, 100, 316).
var benchGridSides = []struct {
	name string
	side int
}{{"1e3", 32}, {"1e4", 100}, {"1e5", 316}}

func BenchmarkPrefixSums2D(b *testing.B) {
	for _, size := range benchGridSides {
		grid := benchGrid(size.side)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				table, _ := PrefixSums2D(grid)
				benchSink = table.Rows()
			}
		})
	}
}

func BenchmarkRectangleSum(b *testing.B) {
	for _, size := range benchGridSides {
		table, _ := PrefixSums2D(benchGrid(size.side))
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			var sum int64
			for b.Loop() {
				for r := 0; r < size.side; r++ {
					s, _ := table.RectangleSum(r/2, r/3, r+1, size.side)
					sum += s
				}
			}
			benchSink = int(sum)
		})
	}
}

func BenchmarkSummedAreaTableRows(b *testing.B) {
	table, _ := PrefixSums2D(grid3x4)
	b.ReportAllocs()
	for b.Loop() {
		benchSink = table.Rows()
	}
}

func BenchmarkSummedAreaTableCols(b *testing.B) {
	table, _ := PrefixSums2D(grid3x4)
	b.ReportAllocs()
	for b.Loop() {
		benchSink = table.Cols()
	}
}
