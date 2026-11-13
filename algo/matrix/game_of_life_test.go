package matrix

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"testing"
)

// fromText reads a grid from rows of '#' (live) and '.' (dead).
func fromText(rows ...string) [][]bool {
	g := make([][]bool, len(rows))
	for i, r := range rows {
		g[i] = make([]bool, len(r))
		for j, ch := range r {
			g[i][j] = ch == '#'
		}
	}
	return g
}

// Gardner 1970: the blinker, a row of three live cells, becomes a column of three and back (period 2).
func TestGameOfLifeStep_gardnerBlinker(t *testing.T) {
	row := fromText(".....", ".....", ".###.", ".....", ".....")
	col := fromText(".....", "..#..", "..#..", "..#..", ".....")
	next, err := GameOfLifeStep(nil, row)
	if err != nil || !gridEqual(next, col) {
		t.Fatalf("generation 1: got %v, %v", next, err)
	}
	back, _ := GameOfLifeStep(nil, next)
	if !gridEqual(back, row) {
		t.Errorf("generation 2: got %v", back)
	}
}

// Gardner 1970: the 2×2 block is a still life.
func TestGameOfLifeStep_gardnerBlock(t *testing.T) {
	block := fromText("....", ".##.", ".##.", "....")
	if next, _ := GameOfLifeStep(nil, block); !gridEqual(next, block) {
		t.Errorf("block changed: %v", next)
	}
}

func TestGameOfLifeStep_cases(t *testing.T) {
	tests := []struct {
		name     string
		in, want [][]bool
	}{
		{"empty", nil, nil},
		{"lone cell dies", fromText("#"), fromText(".")},
		{"pair dies", fromText("##"), fromText("..")},
		{"birth in a corner from three neighbours", fromText("##", "#."), fromText("##", "##")},
		{"overcrowded centre dies", fromText("###", "###", "###"), fromText("#.#", "...", "#.#")},
		{"edge is dead, not wrapped", fromText("#..#"), fromText("....")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GameOfLifeStep(nil, tt.in)
			if err != nil || len(got) != len(tt.want) || (len(got) > 0 && !gridEqual(got, tt.want)) {
				t.Errorf("got %v, %v; want %v", got, err, tt.want)
			}
		})
	}
}

func TestGameOfLifeStep_errors(t *testing.T) {
	if _, err := GameOfLifeStep(nil, [][]bool{{true}, {true, false}}); !errors.Is(err, ErrRagged) {
		t.Errorf("err = %v, want ErrRagged", err)
	}
}

func TestGameOfLifeStep_doubleBufferNoAllocation(t *testing.T) {
	a := fromText(".....", ".....", ".###.", ".....", ".....")
	b, _ := GameOfLifeStep(nil, a)
	allocs := testing.AllocsPerRun(50, func() {
		b, _ = GameOfLifeStep(b, a)
		a, _ = GameOfLifeStep(a, b)
	})
	if allocs != 0 {
		t.Errorf("alternating buffers allocates %v times, want 0", allocs)
	}
}

// TestGameOfLifeStep_oracle: a padded-border version that counts neighbours by explicit offsets
// against the rules stated in Gardner 1970, and the input is unchanged.
func TestGameOfLifeStep_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		m, n := 1+rng.IntN(8), 1+rng.IntN(8)
		g := make([][]bool, m)
		for i := range g {
			g[i] = make([]bool, n)
			for j := range g[i] {
				g[i][j] = rng.IntN(2) == 1
			}
		}
		before := cloneGrid(g)
		live := func(r, c int) int {
			if r < 0 || r >= m || c < 0 || c >= n || !g[r][c] {
				return 0
			}
			return 1
		}
		got, _ := GameOfLifeStep(nil, g)
		for r := range m {
			for c := range n {
				k := live(r-1, c-1) + live(r-1, c) + live(r-1, c+1) + live(r, c-1) + live(r, c+1) + live(r+1, c-1) + live(r+1, c) + live(r+1, c+1)
				want := false
				if g[r][c] {
					want = k == 2 || k == 3
				} else {
					want = k == 3
				}
				if got[r][c] != want {
					t.Fatalf("trial %d: cell (%d,%d) with %d neighbours: got %v", trial, r, c, k, got[r][c])
				}
			}
		}
		if !gridEqual(g, before) {
			t.Fatalf("trial %d: input modified", trial)
		}
	}
}

func BenchmarkGameOfLifeStep(b *testing.B) {
	for _, n := range []int{32, 100, 316} {
		rng := rand.New(rand.NewPCG(1, 2))
		g := make([][]bool, n)
		for i := range g {
			g[i] = make([]bool, n)
			for j := range g[i] {
				g[i][j] = rng.IntN(3) == 0
			}
		}
		dst, _ := GameOfLifeStep(nil, g)
		b.Run(fmt.Sprint(n*n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				dst, _ = GameOfLifeStep(dst, g)
			}
			if dst[0][0] {
				benchSink++
			}
		})
	}
}
