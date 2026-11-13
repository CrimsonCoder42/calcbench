// Package matrix holds algorithms on rectangular grids [][]T: rotation and transposition, a sparse
// (CSR) matrix product, a Game of Life generation, and the largest rectangle sum not above a bound.
//
// ID short form: MATRIX (IDs are ALG-MATRIX-NNN).
//
// Scope: these are grid algorithms, generic over the cell type or combinatorial. Float64 linear algebra
// (solving, decompositions, condition numbers) is num/linalg; searching a sorted matrix is algo/search
// (ALG-SEARCH-009, -010).
//
// Grids: grid[r][c], r the row. Every row must have the same length (else ErrRagged); a grid with no
// rows is empty. Inputs are never modified except by Rotate90InPlace, whose name says so. A function
// producing a grid reuses dst when it already has the result's shape and otherwise allocates once (one
// backing array, one row-header slice); dst must not share memory with the input.
//
// Numbers: integer sums and products are exact or ErrOverflow, never wrapped (a partial sum that leaves
// int64 is ErrOverflow even when the final answer would fit, as in algo/prefix). Float arithmetic follows
// IEEE 754: Inf and NaN propagate and are not errors.
//
// Functions:
//
//	ALG-MATRIX-001  Rotate90InPlace               rotate a square grid a quarter turn clockwise, in place
//	ALG-MATRIX-002  Transpose                     the transpose of a grid
//	ALG-MATRIX-003  NewCSR                        compressed sparse row form of a dense grid
//	ALG-MATRIX-004  SparseMultiply                the product of two CSR matrices (Gustavson)
//	ALG-MATRIX-005  GameOfLifeStep                one generation of Conway's Game of Life
//	ALG-MATRIX-006  MaxRectangleSumNoLargerThanK  the largest sum of a sub-rectangle that is ≤ k
package matrix
