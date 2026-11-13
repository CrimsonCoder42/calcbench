package graph

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"

	"github.com/CrimsonCoder42/calcbench/algo/unionfind"
)

func isOne(b byte) bool { return b == '1' }

// Hand example (catalog/graph.md): four islands under 4 neighbours; the diagonals join them all under 8.
func TestNumberOfIslands_handExample(t *testing.T) {
	grid := toGrid(
		"11000",
		"11010",
		"00100",
		"00011",
	)
	if got, err := NumberOfIslands(grid, isOne, FourNeighbors); got != 4 || err != nil {
		t.Errorf("4 neighbours: (%d, %v), want (4, nil)", got, err)
	}
	if got, err := NumberOfIslands(grid, isOne, EightNeighbors); got != 1 || err != nil {
		t.Errorf("8 neighbours: (%d, %v), want (1, nil)", got, err)
	}
}

func TestNumberOfIslands_cases(t *testing.T) {
	tests := []struct {
		name  string
		grid  [][]byte
		four  int
		eight int
	}{
		{"no rows", [][]byte{}, 0, 0},
		{"rows of zero length", toGrid("", ""), 0, 0},
		{"1×1 water", toGrid("0"), 0, 0},
		{"1×1 land", toGrid("1"), 1, 1},
		{"all water", toGrid("000", "000"), 0, 0},
		{"all land", toGrid("111", "111"), 1, 1},
		{"checkerboard", toGrid("101", "010", "101"), 5, 1},
		{"ring around a lake", toGrid("111", "101", "111"), 1, 1},
	}
	for _, tt := range tests {
		f, err1 := NumberOfIslands(tt.grid, isOne, FourNeighbors)
		e, err2 := NumberOfIslands(tt.grid, isOne, EightNeighbors)
		if f != tt.four || e != tt.eight || err1 != nil || err2 != nil {
			t.Errorf("%s: (%d, %d) errs %v %v, want (%d, %d)", tt.name, f, e, err1, err2, tt.four, tt.eight)
		}
	}
}

func TestNumberOfIslands_errors(t *testing.T) {
	if _, err := NumberOfIslands(toGrid("1"), isOne, Neighborhood(-4)); !errors.Is(err, ErrNeighborhood) {
		t.Errorf("neighborhood: %v", err)
	}
	if _, err := NumberOfIslands(toGrid("11", "1"), isOne, FourNeighbors); !errors.Is(err, ErrNotRectangular) {
		t.Errorf("ragged: %v", err)
	}
	if _, err := NumberOfIslands(toGrid("1"), nil, FourNeighbors); !errors.Is(err, ErrNilPredicate) {
		t.Errorf("nil predicate: %v", err)
	}
	huge := [][]struct{}{make([]struct{}, math.MaxInt), make([]struct{}, math.MaxInt)}
	if _, err := NumberOfIslands(huge, func(struct{}) bool { return false }, FourNeighbors); !errors.Is(err, ErrTooLarge) {
		t.Errorf("rows×cols overflow: %v", err)
	}
}

// Oracle: under 4 neighbours the count equals unionfind.IslandsOnline's final count after adding every
// land cell; under 8 neighbours, ConnectedComponents of the explicit grid graph minus the water cells.
func TestNumberOfIslands_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for run := 0; run < 300; run++ {
		rows, cols := 1+rng.IntN(10), 1+rng.IntN(10)
		grid := randomGrid(rng, rows, cols, 0.5)
		var land []unionfind.Cell
		for r := range rows {
			for c := range cols {
				if isDot(grid[r][c]) {
					land = append(land, unionfind.Cell{Row: r, Col: c})
				}
			}
		}
		want4 := 0
		if counts, _ := unionfind.IslandsOnline(nil, rows, cols, land); len(counts) > 0 {
			want4 = counts[len(counts)-1]
		}
		if got, err := NumberOfIslands(grid, isDot, FourNeighbors); got != want4 || err != nil {
			t.Fatalf("run %d: 4-neighbour %d (%v), IslandsOnline %d", run, got, err, want4)
		}
		_, comps, _ := ConnectedComponents(gridAsGraph(grid, EightNeighbors))
		want8 := comps - (rows*cols - len(land)) // each water cell is a component of its own
		if got, err := NumberOfIslands(grid, isDot, EightNeighbors); got != want8 || err != nil {
			t.Fatalf("run %d: 8-neighbour %d (%v), components %d", run, got, err, want8)
		}
	}
}

func BenchmarkNumberOfIslands(b *testing.B) {
	for _, size := range gridBenchSides {
		grid := randomGrid(rand.New(rand.NewPCG(1, 2)), size.side, size.side, 0.5)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = NumberOfIslands(grid, isDot, FourNeighbors)
			}
		})
	}
}
