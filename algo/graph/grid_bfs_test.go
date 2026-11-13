package graph

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// toGrid turns rows of text into a byte grid.
func toGrid(rows ...string) [][]byte {
	g := make([][]byte, len(rows))
	for i, r := range rows {
		g[i] = []byte(r)
	}
	return g
}

func isDot(b byte) bool { return b == '.' }

// cells is a short way to write a list of cells.
func cells(rc ...int) []Cell {
	out := make([]Cell, 0, len(rc)/2)
	for i := 0; i+1 < len(rc); i += 2 {
		out = append(out, Cell{Row: rc[i], Col: rc[i+1]})
	}
	return out
}

// gridAsGraph builds the directed graph whose vertex r*cols+c is a passable cell, with arcs to its
// passable neighbours in nb's order: the explicit graph the grid searches walk implicitly.
func gridAsGraph(grid [][]byte, nb Neighborhood) *Graph {
	rows, cols := len(grid), len(grid[0])
	steps, _ := nb.offsets()
	g, _ := NewDirected(rows*cols, nil)
	for r := range rows {
		for c := range cols {
			if !isDot(grid[r][c]) {
				continue
			}
			for _, d := range steps {
				rr, cc := r+d.Row, c+d.Col
				if rr >= 0 && rr < rows && cc >= 0 && cc < cols && isDot(grid[rr][cc]) {
					_ = g.AddEdge(r*cols+c, rr*cols+cc, 0)
				}
			}
		}
	}
	return g
}

func randomGrid(rng *rand.Rand, rows, cols int, pDot float64) [][]byte {
	g := make([][]byte, rows)
	for r := range g {
		g[r] = make([]byte, cols)
		for c := range g[r] {
			g[r][c] = '#'
			if rng.Float64() < pDot {
				g[r][c] = '.'
			}
		}
	}
	return g
}

// Hand example (catalog/graph.md): BFS on the grid graph (CLRS §22.2) from the top-left, 4 neighbours.
func TestGridBFS_handExample(t *testing.T) {
	grid := toGrid(
		"..#",
		".##",
		"...",
	)
	seq, err := GridBFS(grid, Cell{}, isDot, FourNeighbors)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := slices.Collect(seq), cells(0, 0, 0, 1, 1, 0, 2, 0, 2, 1, 2, 2); !slices.Equal(got, want) {
		t.Errorf("GridBFS = %v, want %v", got, want)
	}
	seq8, _ := GridBFS(grid, Cell{}, isDot, EightNeighbors)
	if got, want := slices.Collect(seq8), cells(0, 0, 0, 1, 1, 0, 2, 1, 2, 0, 2, 2); !slices.Equal(got, want) {
		t.Errorf("GridBFS 8 = %v, want %v", got, want)
	}
}

func TestGridBFS_cases(t *testing.T) {
	grid := toGrid("#.", "..")
	if got := slices.Collect(must(GridBFS(grid, Cell{}, isDot, FourNeighbors))); len(got) != 0 {
		t.Errorf("blocked start: %v, want empty", got)
	}
	if got := slices.Collect(must(GridBFS(toGrid("."), Cell{}, isDot, EightNeighbors))); !slices.Equal(got, cells(0, 0)) {
		t.Errorf("1×1: %v", got)
	}
	seq := must(GridBFS(grid, Cell{Row: 1, Col: 1}, isDot, FourNeighbors))
	for range seq {
		break
	}
	if got := slices.Collect(seq); len(got) != 3 {
		t.Errorf("re-ranging: %d cells, want 3", len(got))
	}
	before := toGrid("#.", "..")
	if !slices.EqualFunc(grid, before, slices.Equal) {
		t.Error("grid modified")
	}
}

func TestGridBFS_errors(t *testing.T) {
	grid := toGrid("..", "..")
	if _, err := GridBFS(grid, Cell{}, isDot, Neighborhood(6)); !errors.Is(err, ErrNeighborhood) {
		t.Errorf("neighborhood 6: %v", err)
	}
	if _, err := GridBFS(toGrid("..", "."), Cell{}, isDot, FourNeighbors); !errors.Is(err, ErrNotRectangular) {
		t.Errorf("ragged: %v", err)
	}
	for _, c := range []Cell{{Row: -1}, {Col: 2}, {Row: 2}, {Col: math.MinInt}} {
		if _, err := GridBFS(grid, c, isDot, FourNeighbors); !errors.Is(err, ErrCellOutOfRange) {
			t.Errorf("start %v: %v", c, err)
		}
	}
	if _, err := GridBFS([][]byte{}, Cell{}, isDot, FourNeighbors); !errors.Is(err, ErrCellOutOfRange) {
		t.Errorf("empty grid: %v", err)
	}
	if seq, err := GridBFS(grid, Cell{}, nil, FourNeighbors); !errors.Is(err, ErrNilPredicate) || seq != nil {
		t.Errorf("nil predicate: %v", err)
	}
	// Zero-size cells: two rows of math.MaxInt cells take no memory but rows×cols overflows int.
	huge := [][]struct{}{make([]struct{}, math.MaxInt), make([]struct{}, math.MaxInt)}
	if _, err := GridBFS(huge, Cell{}, func(struct{}) bool { return true }, FourNeighbors); !errors.Is(err, ErrTooLarge) {
		t.Errorf("rows×cols overflow: %v", err)
	}
	// L41: one row of maxAllocBytes+1 zero-size cells does not overflow int, but its visited bitmap would
	// exceed the allocation limit.
	if wideLen := uint64(maxAllocBytes) + 1; wideLen <= math.MaxInt { // 64-bit only: on 32-bit no int holds it
		wide := [][]struct{}{make([]struct{}, int(wideLen))}
		if _, err := GridBFS(wide, Cell{}, func(struct{}) bool { return true }, FourNeighbors); !errors.Is(err, ErrTooLarge) {
			t.Errorf("rows×cols above the allocation limit: %v", err)
		}
	}
}

// Oracle: GridBFS order equals BFS on the explicit grid graph (same neighbour order).
func TestGridBFS_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for run := 0; run < 300; run++ {
		rows, cols := 1+rng.IntN(8), 1+rng.IntN(8)
		grid := randomGrid(rng, rows, cols, 0.6)
		nb := FourNeighbors
		if run%2 == 1 {
			nb = EightNeighbors
		}
		start := Cell{Row: rng.IntN(rows), Col: rng.IntN(cols)}
		var got []int
		for c := range must(GridBFS(grid, start, isDot, nb)) {
			got = append(got, c.Row*cols+c.Col)
		}
		var want []int
		if isDot(grid[start.Row][start.Col]) {
			want = slices.Collect(must(BFS(gridAsGraph(grid, nb), start.Row*cols+start.Col)))
		}
		if !slices.Equal(got, want) {
			t.Fatalf("run %d: GridBFS %v, graph BFS %v", run, got, want)
		}
	}
}

// openGrid returns a side×side grid of '.' cells: every cell is reachable, so the work grows with n.
func openGrid(side int) [][]byte {
	g := make([][]byte, side)
	for r := range g {
		g[r] = make([]byte, side)
		for c := range g[r] {
			g[r][c] = '.'
		}
	}
	return g
}

// gridBenchSides gives square grids of 1e4, 1e5 and 1e6 cells.
var gridBenchSides = []struct {
	name string
	side int
}{{"1e4", 100}, {"1e5", 316}, {"1e6", 1000}}

func BenchmarkGridBFS(b *testing.B) {
	for _, size := range gridBenchSides {
		grid := openGrid(size.side)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				seq, _ := GridBFS(grid, Cell{}, isDot, FourNeighbors)
				for range seq {
				}
			}
		})
	}
}
