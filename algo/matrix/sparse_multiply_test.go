package matrix

import (
	"errors"
	"fmt"
	"maps"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// toDense expands a CSR matrix into a grid (repeated columns add).
func toDense[T number](a *CSR[T]) [][]T {
	g := make([][]T, a.Rows)
	for i := range g {
		g[i] = make([]T, a.Cols)
		for p := a.RowStart[i]; p < a.RowStart[i+1]; p++ {
			g[i][a.ColIndex[p]] += a.Values[p]
		}
	}
	return g
}

// denseMultiply is the slow obvious oracle: the triple loop. p, the column count of b, is passed
// explicitly because a 0×p grid has no row to read it from (m×0 · 0×p is m×p zeros).
func denseMultiply(a, b [][]int, p int) [][]int {
	m, n := len(a), len(b)
	c := make([][]int, m)
	for i := range c {
		c[i] = make([]int, p)
		for k := range n {
			for j := range p {
				c[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return c
}

func mustCSR[T number](t *testing.T, g [][]T) *CSR[T] {
	t.Helper()
	c, err := NewCSR(g)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Gustavson 1978: C = A·B row by row. Our own 2×3 by 3×2 example, checked by hand:
// row 0 = 1·[0 4] + 2·[6 0] = [12 4]; row 1 = 3·[5 0] = [15 0].
func TestSparseMultiply_workedExample(t *testing.T) {
	a := mustCSR(t, [][]int{{1, 0, 2}, {0, 3, 0}})
	b := mustCSR(t, [][]int{{0, 4}, {5, 0}, {6, 0}})
	c, err := SparseMultiply(a, b)
	if err != nil {
		t.Fatal(err)
	}
	if want := [][]int{{12, 4}, {15, 0}}; !gridEqual(toDense(c), want) {
		t.Errorf("got %v, want %v", toDense(c), want)
	}
	if !slices.Equal(c.RowStart, []int{0, 2, 3}) || !slices.Equal(c.ColIndex, []int{0, 1, 0}) {
		t.Errorf("CSR layout: RowStart %v ColIndex %v", c.RowStart, c.ColIndex)
	}
}

func TestNewCSR_cases(t *testing.T) {
	c := mustCSR(t, [][]float64{{0, math.Copysign(0, -1), 1.5}, {math.NaN(), 0, 0}})
	if c.Rows != 2 || c.Cols != 3 || !slices.Equal(c.RowStart, []int{0, 1, 2}) || !slices.Equal(c.ColIndex, []int{2, 0}) {
		t.Errorf("got %+v", c)
	}
	e := mustCSR[int](t, nil)
	if e.Rows != 0 || len(e.RowStart) != 1 {
		t.Errorf("empty: %+v", e)
	}
	if _, err := NewCSR([][]int{{1}, {1, 2}}); !errors.Is(err, ErrRagged) {
		t.Errorf("ragged: err = %v", err)
	}
}

func TestSparseMultiply_cases(t *testing.T) {
	tests := []struct {
		name string
		a, b [][]int
	}{
		{"one by one", [][]int{{3}}, [][]int{{-4}}},
		{"all zero", [][]int{{0, 0}}, [][]int{{0}, {0}}},
		{"cancelling terms dropped", [][]int{{1, 1}}, [][]int{{2}, {-2}}},
		{"identity", [][]int{{1, 0}, {0, 1}}, [][]int{{5, 6}, {7, 8}}},
		{"row times column", [][]int{{1, 2, 3}}, [][]int{{4}, {5}, {6}}},
		{"column times row", [][]int{{1}, {2}}, [][]int{{3, 4}}},
		{"extremes", [][]int{{math.MaxInt, 0}}, [][]int{{1}, {math.MinInt}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := SparseMultiply(mustCSR(t, tt.a), mustCSR(t, tt.b))
			if err != nil {
				t.Fatal(err)
			}
			if want := denseMultiply(tt.a, tt.b, len(tt.b[0])); !gridEqual(toDense(c), want) {
				t.Errorf("got %v, want %v", toDense(c), want)
			}
			for _, v := range c.Values {
				if v == 0 {
					t.Errorf("explicit zero kept: %+v", c)
				}
			}
		})
	}
}

func TestSparseMultiply_zeroDimensions(t *testing.T) {
	a := &CSR[int]{Rows: 2, Cols: 0, RowStart: []int{0, 0, 0}}
	b := &CSR[int]{Rows: 0, Cols: 3, RowStart: []int{0}}
	c, err := SparseMultiply(a, b)
	if err != nil || c.Rows != 2 || c.Cols != 3 || len(c.Values) != 0 {
		t.Errorf("got %+v, %v", c, err)
	}
}

func TestSparseMultiply_floats(t *testing.T) {
	a := mustCSR(t, [][]float64{{math.Inf(1), 1}})
	b := mustCSR(t, [][]float64{{2}, {math.NaN()}})
	c, err := SparseMultiply(a, b)
	if err != nil || len(c.Values) != 1 || !math.IsNaN(c.Values[0]) {
		t.Errorf("Inf·2 + 1·NaN: got %+v, %v; want one NaN entry", c, err)
	}
	half := mustCSR(t, [][]float64{{0.5}})
	c, _ = SparseMultiply(half, half)
	if c.Values[0] != 0.25 {
		t.Errorf("0.5·0.5 = %v", c.Values[0])
	}
}

// Unsorted and repeated columns in an input row are accepted; the output is canonical.
func TestSparseMultiply_unsortedRepeatedInput(t *testing.T) {
	a := &CSR[int]{Rows: 1, Cols: 2, RowStart: []int{0, 3}, ColIndex: []int{1, 0, 1}, Values: []int{1, 2, 3}}
	b := mustCSR(t, [][]int{{1, 0, 0}, {0, 0, 1}})
	c, err := SparseMultiply(a, b)
	if err != nil || !slices.Equal(c.ColIndex, []int{0, 2}) || !slices.Equal(c.Values, []int{2, 4}) {
		t.Errorf("got %+v, %v", c, err)
	}
}

// The map accumulator path: b.Cols far above 4·nnz(B)+64 must not allocate b.Cols cells.
func TestSparseMultiply_hugeColsUsesMap(t *testing.T) {
	const huge = 1 << 60
	a := &CSR[int]{Rows: 1, Cols: 1, RowStart: []int{0, 1}, ColIndex: []int{0}, Values: []int{3}}
	b := &CSR[int]{Rows: 1, Cols: huge, RowStart: []int{0, 2}, ColIndex: []int{huge - 1, 5}, Values: []int{2, 7}}
	c, err := SparseMultiply(a, b)
	if err != nil || !slices.Equal(c.ColIndex, []int{5, huge - 1}) || !slices.Equal(c.Values, []int{21, 6}) {
		t.Errorf("got %+v, %v", c, err)
	}
}

func TestSparseMultiply_errors(t *testing.T) {
	ok := &CSR[int]{Rows: 1, Cols: 1, RowStart: []int{0, 1}, ColIndex: []int{0}, Values: []int{1}}
	tests := []struct {
		name string
		a, b *CSR[int]
		want error
	}{
		{"nil a", nil, ok, ErrInvalidCSR},
		{"nil b", ok, nil, ErrInvalidCSR},
		{"negative rows", &CSR[int]{Rows: -1, RowStart: []int{}}, ok, ErrInvalidCSR},
		{"short RowStart", &CSR[int]{Rows: 1, Cols: 1, RowStart: []int{0}}, ok, ErrInvalidCSR},
		{"RowStart[0] not 0", &CSR[int]{Rows: 1, Cols: 1, RowStart: []int{1, 1}, ColIndex: []int{0}, Values: []int{1}}, ok, ErrInvalidCSR},
		{"length mismatch", &CSR[int]{Rows: 1, Cols: 1, RowStart: []int{0, 1}, ColIndex: []int{0}}, ok, ErrInvalidCSR},
		{"decreasing RowStart", &CSR[int]{Rows: 2, Cols: 1, RowStart: []int{0, 2, 1}, ColIndex: []int{0}, Values: []int{1}}, ok, ErrInvalidCSR},
		{"column out of range", &CSR[int]{Rows: 1, Cols: 1, RowStart: []int{0, 1}, ColIndex: []int{1}, Values: []int{1}}, ok, ErrInvalidCSR},
		{"negative column", &CSR[int]{Rows: 1, Cols: 1, RowStart: []int{0, 1}, ColIndex: []int{-1}, Values: []int{1}}, ok, ErrInvalidCSR},
		{"shape", ok, &CSR[int]{Rows: 2, Cols: 1, RowStart: []int{0, 0, 0}}, ErrShape},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := SparseMultiply(tt.a, tt.b); !errors.Is(err, tt.want) {
				t.Errorf("err = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestSparseMultiply_overflow(t *testing.T) {
	tests := []struct {
		name string
		a, b [][]int64
	}{
		{"product", [][]int64{{math.MaxInt64}}, [][]int64{{2}}},
		{"minus one times min", [][]int64{{-1}}, [][]int64{{math.MinInt64}}},
		{"min times minus one", [][]int64{{math.MinInt64}}, [][]int64{{-1}}},
		{"sum", [][]int64{{1, 1}}, [][]int64{{math.MaxInt64}, {1}}},
		{"negative sum", [][]int64{{1, 1}}, [][]int64{{math.MinInt64}, {-1}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := SparseMultiply(mustCSR(t, tt.a), mustCSR(t, tt.b)); !errors.Is(err, ErrOverflow) {
				t.Errorf("err = %v, want ErrOverflow", err)
			}
		})
	}
	// int8: 100·2 overflows, 50·2 does not.
	if _, err := SparseMultiply(mustCSR(t, [][]int8{{100}}), mustCSR(t, [][]int8{{2}})); !errors.Is(err, ErrOverflow) {
		t.Errorf("int8 100·2: err = %v", err)
	}
	if c, err := SparseMultiply(mustCSR(t, [][]int8{{-64}}), mustCSR(t, [][]int8{{2}})); err != nil || c.Values[0] != -128 {
		t.Errorf("int8 −64·2 = −128 fits: %+v, %v", c, err)
	}
}

func TestSparseMultiply_inputNotModified(t *testing.T) {
	a := mustCSR(t, [][]int{{1, 2}, {0, 3}})
	b := mustCSR(t, [][]int{{4, 0}, {5, 6}})
	ac, bc := *a, *b
	ac.Values, bc.Values = slices.Clone(a.Values), slices.Clone(b.Values)
	_, _ = SparseMultiply(a, b)
	if !slices.Equal(a.Values, ac.Values) || !slices.Equal(b.Values, bc.Values) {
		t.Errorf("inputs modified")
	}
}

// TestSparseMultiply_oracle: against the dense triple loop on random sparse integer matrices.
func TestSparseMultiply_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	sparse := func(m, n int) [][]int {
		g := randomGrid(rng, m, n, -5, 6)
		for _, row := range g {
			for j := range row {
				if rng.IntN(3) > 0 {
					row[j] = 0
				}
			}
		}
		return g
	}
	for trial := range 300 {
		m, n, p := rng.IntN(7), rng.IntN(7), rng.IntN(7)
		a, b := sparse(m, n), sparse(n, p)
		ca, cb := mustCSR(t, a), mustCSR(t, b)
		ca.Cols, cb.Cols = n, p // a grid with no rows still has its columns
		c, err := SparseMultiply(ca, cb)
		if err != nil {
			t.Fatalf("trial %d: %v", trial, err)
		}
		want := denseMultiply(a, b, p)
		got := toDense(c)
		for i := range want {
			if !slices.Equal(got[i], want[i]) {
				t.Fatalf("trial %d: row %d = %v, want %v", trial, i, got[i], want[i])
			}
		}
		for i := range c.Rows {
			row := c.ColIndex[c.RowStart[i]:c.RowStart[i+1]]
			if !slices.IsSorted(row) || len(slices.Compact(slices.Clone(row))) != len(row) {
				t.Fatalf("trial %d: row %d columns %v not strictly ascending", trial, i, row)
			}
		}
	}
}

// BenchmarkSparseMultiply: n×n matrices with about 5 nonzeros per row.
func BenchmarkSparseMultiply(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		rng := rand.New(rand.NewPCG(1, 2))
		mk := func() *CSR[float64] {
			c := &CSR[float64]{Rows: n, Cols: n, RowStart: make([]int, n+1)}
			for i := range n {
				cols := map[int]bool{}
				for range 5 {
					cols[rng.IntN(n)] = true
				}
				for _, j := range slices.Sorted(maps.Keys(cols)) {
					c.ColIndex = append(c.ColIndex, j)
					c.Values = append(c.Values, rng.Float64())
				}
				c.RowStart[i+1] = len(c.Values)
			}
			return c
		}
		x, y := mk(), mk()
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				c, _ := SparseMultiply(x, y)
				benchSink = len(c.Values)
			}
		})
	}
}
