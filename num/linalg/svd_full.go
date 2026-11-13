package linalg

import "math"

// SVDFull returns the full singular value decomposition A = U·Σ·Vᵀ of an m×n matrix (numpy's svd with
// full_matrices=True): U is m×m and V is n×n, both orthogonal, row-major, and S holds the min(m, n)
// singular values in descending order (Σ is m×n with S on its diagonal). The thin factors come from
// SVDJacobi (NUM-LINALG-160); columns it leaves zero (for zero singular values) and the missing
// columns are completed to an orthonormal basis by modified Gram–Schmidt, applied twice, against the unit
// vectors e₁, e₂, … in order. The completion is not unique; any orthonormal completion satisfies
// A = U·Σ·Vᵀ.
// Errors: as SVDJacobi.
//
// ID:      NUM-LINALG-189
// Time:    O(m·n·min(m, n)) per SVD sweep, plus O(m³ + n³) for the completion
// Space:   O(m² + n²) for the result and O(m·n) for the SVD
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §2.4 "The Singular Value Decomposition" (Theorem 2.4.1, full U and V; §2.4.3 the thin SVD) and §5.2.8 (modified Gram–Schmidt)
// Tags:    matrix, n-dimensional, svd, full-svd, singular-values, orthogonal, factorization, linear-algebra
// Units:   a in unit A; singular values in A; U and V dimensionless
// Example: the 3×1 column (3, 4, 0) → S = (5), U's first column (0.6, 0.8, 0), V = [1] (hand-checked)
func SVDFull(a []float64, m, n int, tol float64) (SVDResult, error) {
	thin, err := SVDJacobi(a, m, n, tol, true)
	if err != nil {
		return SVDResult{}, err
	}
	k := min(m, n)
	thin.U = completeBasis(thin.U, m, k)
	thin.V = completeBasis(thin.V, n, k)
	return thin, nil
}

// completeBasis widens the p×k matrix q (orthonormal or zero columns) to a p×p orthogonal matrix: zero
// columns and columns k…p−1 are filled with unit vectors orthogonalised against the columns already set.
func completeBasis(q []float64, p, k int) []float64 {
	out := make([]float64, p*p)
	for i := range p {
		copy(out[i*p:i*p+k], q[i*k:i*k+k])
	}
	col := make([]float64, p)
	// Some eᵢ always keeps ‖residual‖ ≥ √(d/p) ≥ 1/√p while d columns are missing, so this cut always finds one;
	// a rejected eᵢ only shrinks as the span grows, so the scan never revisits it.
	cut := 0.5 / math.Sqrt(float64(p))
	set := make([]bool, p)
	for j := range k {
		s := 0.0
		for i := range p {
			s += out[i*p+j] * out[i*p+j]
		}
		set[j] = s > 0.5
	}
	cand := 0
	for j := range p {
		if set[j] {
			continue
		}
		for ; cand < p; cand++ {
			clear(col)
			col[cand] = 1
			for range 2 {
				for c := range p {
					if !set[c] {
						continue
					}
					d := 0.0
					for i := range p {
						d += out[i*p+c] * col[i]
					}
					for i := range p {
						col[i] -= d * out[i*p+c]
					}
				}
			}
			nrm := 0.0
			for _, x := range col {
				nrm += x * x
			}
			nrm = math.Sqrt(nrm)
			if nrm > cut {
				for i := range p {
					out[i*p+j] = col[i] / nrm
				}
				set[j] = true
				cand++
				break
			}
		}
	}
	return out
}
