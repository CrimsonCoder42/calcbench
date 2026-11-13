package matrix

import (
	"fmt"
	"math"
	"slices"
)

// number is the cell types of a sparse product: signed integers (overflow-checked) and floats (IEEE).
type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

// CSR is a sparse matrix in compressed sparse row form: the nonzeros of row i are
// Values[RowStart[i]:RowStart[i+1]], in the columns ColIndex[RowStart[i]:RowStart[i+1]]. A valid CSR has
// Rows, Cols ≥ 0, len(RowStart) = Rows+1, RowStart[0] = 0, RowStart non-decreasing, RowStart[Rows] =
// len(ColIndex) = len(Values), and every column in [0, Cols). Columns within a row may be unsorted or
// repeated (repeats add). NewCSR and SparseMultiply produce rows sorted by column without repeats or zeros.
type CSR[T number] struct {
	Rows, Cols int
	RowStart   []int
	ColIndex   []int
	Values     []T
}

// NewCSR returns the CSR form of the dense grid, keeping its nonzero cells (a NaN is nonzero; −0.0 is
// zero) row by row in column order. dense is read, never modified.
// Errors: ErrRagged.
//
// ID:      ALG-MATRIX-003
// Time:    O(mn)   m = rows, n = columns
// Space:   O(m + z) for the result   z = number of nonzero cells
// Source:  Gustavson, "Two Fast Algorithms for Sparse Matrices: Multiplication and Permuted Transposition", ACM Transactions on Mathematical Software 4(3):250–269, 1978, doi:10.1145/355791.355796 (the row-wise sparse storage)
// Tags:    sparse-matrix, csr, compressed-sparse-row
func NewCSR[T number](dense [][]T) (*CSR[T], error) {
	m, n, err := gridShape(dense)
	if err != nil {
		return nil, err
	}
	out := &CSR[T]{Rows: m, Cols: n, RowStart: make([]int, m+1)}
	for r, row := range dense {
		for c, v := range row {
			if v != 0 {
				out.ColIndex = append(out.ColIndex, c)
				out.Values = append(out.Values, v)
			}
		}
		out.RowStart[r+1] = len(out.Values)
	}
	return out, nil
}

// validate reports ErrInvalidCSR for a malformed matrix; name labels the parameter.
func (a *CSR[T]) validate(name string) error {
	switch {
	case a == nil:
		return fmt.Errorf("%w: %s is nil", ErrInvalidCSR, name)
	case a.Rows < 0 || a.Cols < 0:
		return fmt.Errorf("%w: %s is %d×%d", ErrInvalidCSR, name, a.Rows, a.Cols)
	case len(a.RowStart) != a.Rows+1:
		return fmt.Errorf("%w: %s.RowStart has %d entries, want Rows+1 = %d", ErrInvalidCSR, name, len(a.RowStart), a.Rows+1)
	case a.RowStart[0] != 0:
		return fmt.Errorf("%w: %s.RowStart[0] = %d", ErrInvalidCSR, name, a.RowStart[0])
	case a.RowStart[a.Rows] != len(a.ColIndex) || len(a.ColIndex) != len(a.Values):
		return fmt.Errorf("%w: %s has RowStart[Rows] = %d, %d columns, %d values", ErrInvalidCSR, name, a.RowStart[a.Rows], len(a.ColIndex), len(a.Values))
	}
	for i := range a.Rows {
		if a.RowStart[i] > a.RowStart[i+1] {
			return fmt.Errorf("%w: %s.RowStart decreases at row %d", ErrInvalidCSR, name, i)
		}
	}
	for p, c := range a.ColIndex {
		if c < 0 || c >= a.Cols {
			return fmt.Errorf("%w: %s.ColIndex[%d] = %d outside [0, %d)", ErrInvalidCSR, name, p, c, a.Cols)
		}
	}
	return nil
}

// isFloat reports whether T is a floating-point type (1/2 is nonzero only there).
func isFloat[T number]() bool {
	var half T = 1
	half /= 2
	return half != 0
}

// mulChecked returns x·y and whether it is exact for an integer T (always true for floats).
func mulChecked[T number](x, y T, float bool) (T, bool) {
	p := x * y
	if float || x == 0 {
		return p, true
	}
	// p/x != y catches every overflow except −1·min, where p wraps to min and min/−1 = min.
	if p/x != y || (x == -1 && y != 0 && p == y) {
		return p, false
	}
	return p, true
}

// addChecked returns x+y and whether it is exact for an integer T (always true for floats).
func addChecked[T number](x, y T, float bool) (T, bool) {
	s := x + y
	if float {
		return s, true
	}
	return s, (y > 0) == (s > x) || y == 0
}

// SparseMultiply returns the product C = A·B of CSR matrices, a.Rows × b.Cols, in CSR form with each
// row's columns ascending and without explicit zeros (an entry whose terms cancel to 0 is dropped). The
// method is Gustavson's: each row of C accumulates a_ik · (row k of B) over the nonzeros a_ik of row i of
// A. The accumulator is a dense array of b.Cols cells when b.Cols ≤ 4·nnz(B)+64, else a map, so memory
// stays proportional to the input. a and b are read, never modified.
// Integers: every product and partial sum is checked (ErrOverflow). Floats: IEEE arithmetic, Inf and NaN
// propagate (a NaN entry is kept).
// Errors: ErrInvalidCSR (nil or malformed operand), ErrShape (a.Cols ≠ b.Rows), ErrOverflow.
//
// ID:      ALG-MATRIX-004
// Time:    O(m + F + Σᵢ zᵢ log zᵢ)   m = a.Rows; F = Σ over nonzeros a_ik of nnz(row k of B) (the multiply-adds); zᵢ = nonzeros in row i of C (sorted)
// Space:   O(nnz(B) + max zᵢ) extra (the accumulator); O(m + nnz(C)) for the result
// Source:  Gustavson, "Two Fast Algorithms for Sparse Matrices: Multiplication and Permuted Transposition", ACM Transactions on Mathematical Software 4(3):250–269, 1978, doi:10.1145/355791.355796
// Tags:    sparse-matrix, csr, matrix-multiply, gustavson, sparse-accumulator
func SparseMultiply[T number](a, b *CSR[T]) (*CSR[T], error) {
	if err := a.validate("a"); err != nil {
		return nil, err
	}
	if err := b.validate("b"); err != nil {
		return nil, err
	}
	if a.Cols != b.Rows {
		return nil, fmt.Errorf("%w: a is %d×%d, b is %d×%d", ErrShape, a.Rows, a.Cols, b.Rows, b.Cols)
	}
	float := isFloat[T]()
	nnzB := len(b.Values)
	useDense := b.Cols <= 4*min(nnzB, (math.MaxInt-64)/4)+64
	var (
		dense   []T
		mark    []int // mark[j] = 1 + the last row of C that touched column j
		sparse  map[int]T
		touched []int
	)
	if useDense {
		dense = make([]T, b.Cols)
		mark = make([]int, b.Cols)
	} else {
		sparse = make(map[int]T)
	}
	c := &CSR[T]{Rows: a.Rows, Cols: b.Cols, RowStart: make([]int, a.Rows+1)}
	for i := range a.Rows {
		touched = touched[:0]
		for p := a.RowStart[i]; p < a.RowStart[i+1]; p++ {
			k, aik := a.ColIndex[p], a.Values[p]
			for q := b.RowStart[k]; q < b.RowStart[k+1]; q++ {
				j := b.ColIndex[q]
				prod, ok := mulChecked(aik, b.Values[q], float)
				if !ok {
					return nil, fmt.Errorf("%w: a[%d][%d]·b[%d][%d]", ErrOverflow, i, k, k, j)
				}
				var cur T
				if useDense {
					if mark[j] != i+1 {
						mark[j] = i + 1
						dense[j] = 0
						touched = append(touched, j)
					}
					cur = dense[j]
				} else {
					var seen bool
					if cur, seen = sparse[j]; !seen {
						touched = append(touched, j)
					}
				}
				sum, ok := addChecked(cur, prod, float)
				if !ok {
					return nil, fmt.Errorf("%w: partial sum of c[%d][%d]", ErrOverflow, i, j)
				}
				if useDense {
					dense[j] = sum
				} else {
					sparse[j] = sum
				}
			}
		}
		slices.Sort(touched)
		for _, j := range touched {
			var v T
			if useDense {
				v = dense[j]
			} else {
				v = sparse[j]
				delete(sparse, j)
			}
			if v != 0 {
				c.ColIndex = append(c.ColIndex, j)
				c.Values = append(c.Values, v)
			}
		}
		c.RowStart[i+1] = len(c.Values)
	}
	return c, nil
}
