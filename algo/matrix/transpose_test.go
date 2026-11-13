package matrix

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"testing"
)

// Cormen et al., 3rd ed., Appendix D.1: the transpose of a 2×3 matrix is 3×2 with rows and columns
// exchanged (our own numbers).
func TestTranspose_definition(t *testing.T) {
	got, err := Transpose(nil, [][]int{{1, 2, 3}, {4, 5, 6}})
	if want := [][]int{{1, 4}, {2, 5}, {3, 6}}; err != nil || !gridEqual(got, want) {
		t.Errorf("got %v, %v; want %v", got, err, want)
	}
}

func TestTranspose_cases(t *testing.T) {
	tests := []struct {
		name string
		in   [][]int
		want [][]int
	}{
		{"empty", nil, nil},
		{"one cell", [][]int{{7}}, [][]int{{7}}},
		{"row", [][]int{{1, 2, 3}}, [][]int{{1}, {2}, {3}}},
		{"column", [][]int{{1}, {2}}, [][]int{{1, 2}}},
		{"rows of zero length", [][]int{{}, {}}, [][]int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Transpose(nil, tt.in)
			if err != nil || len(got) != len(tt.want) || (len(got) > 0 && !gridEqual(got, tt.want)) {
				t.Errorf("got %v, %v; want %v", got, err, tt.want)
			}
		})
	}
}

func TestTranspose_reusesDst(t *testing.T) {
	dst := [][]int{{0, 0}, {0, 0}, {0, 0}}
	got, err := Transpose(dst, [][]int{{1, 2, 3}, {4, 5, 6}})
	if err != nil || &got[0][0] != &dst[0][0] {
		t.Fatalf("dst of the right shape was not reused")
	}
	src := [][]int{{1, 2, 3}, {4, 5, 6}}
	if allocs := testing.AllocsPerRun(50, func() { _, _ = Transpose(dst, src) }); allocs != 0 {
		t.Errorf("Transpose with a reusable dst allocates %v times, want 0", allocs)
	}
	got, _ = Transpose([][]int{{0}}, [][]int{{1, 2}})
	if !gridEqual(got, [][]int{{1}, {2}}) {
		t.Errorf("wrong-shape dst: %v", got)
	}
}

func TestTranspose_errors(t *testing.T) {
	if _, err := Transpose(nil, [][]int{{1, 2}, {3}}); !errors.Is(err, ErrRagged) {
		t.Errorf("err = %v, want ErrRagged", err)
	}
}

// TestTranspose_oracle: transposing twice gives the input back, and out[c][r] = in[r][c].
func TestTranspose_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		m, n := 1+rng.IntN(8), 1+rng.IntN(8)
		g := randomGrid(rng, m, n, -9, 9)
		before := cloneGrid(g)
		tr, _ := Transpose(nil, g)
		for r := range m {
			for c := range n {
				if tr[c][r] != g[r][c] {
					t.Fatalf("trial %d: tr[%d][%d] = %d, want %d", trial, c, r, tr[c][r], g[r][c])
				}
			}
		}
		back, _ := Transpose(nil, tr)
		if !gridEqual(back, g) || !gridEqual(g, before) {
			t.Fatalf("trial %d: double transpose or input changed", trial)
		}
	}
}

func BenchmarkTranspose(b *testing.B) {
	for _, n := range []int{32, 100, 316} {
		g := randomGrid(rand.New(rand.NewPCG(1, 2)), n, n, 0, 100)
		dst, _ := Transpose(nil, g)
		b.Run(fmt.Sprint(n*n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				dst, _ = Transpose(dst, g)
			}
			benchSink = dst[0][0]
		})
	}
}
