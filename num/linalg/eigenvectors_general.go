package linalg

import (
	"fmt"
	"math"
	"math/cmplx"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// MatEigenvectors returns the n eigenvalues of a general real n×n matrix (as MatEigenvalues,
// NUM-LINALG-169, in the same order) and, in vectorsDst (n×n, row-major), the matching eigenvectors as
// columns: column k is v with A·v = values[k]·v. The matrix (scaled by an exact power of two) is reduced
// to upper Hessenberg form H = QᵀAQ with Q accumulated; for each eigenvalue λ, inverse iteration solves
// (H − λI)·z = b by Gaussian elimination with adjacent-row pivoting (O(n²) on a Hessenberg matrix),
// three steps from a fixed start vector, and v = Q·z. A pivot smaller than ε₃ = ‖H‖₁·2⁻⁵² is replaced by
// ε₃, and a shift that repeats an earlier eigenvalue (within √ε·‖H‖₁) is moved by ε₃ per repeat and its
// vector orthogonalised against the earlier ones, so a repeated eigenvalue with a full eigenspace gets
// independent vectors; for a defective eigenvalue (no full eigenspace) the vectors are (near) parallel,
// as in LAPACK. The vector of the second eigenvalue of a conjugate pair is the conjugate of the first.
// Each vector has unit 2-norm and is rotated so that its largest-modulus entry is real and positive; for
// a real eigenvalue the vector is real (zero imaginary parts). a is not modified.
// Errors: ErrDimension when n < 1 or len(a) ≠ n²; ErrNotFinite for a NaN or ±Inf element;
// ErrNoConvergence when the eigenvalue iteration fails (see MatEigenvalues) or an inverse-iteration
// vector cannot be formed.
//
// ID:      NUM-LINALG-191
// Time:    O(n³)   n = order: the eigenvalues and the Hessenberg reduction with Q, then O(n²) per eigenvalue for the factorisation and three solves, and n³ for v = Q·z
// Space:   O(n²) extra (the Hessenberg matrix, Q, one complex factorisation and the z vectors); O(n + n²) for the results
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §7.6.1 "Selected Eigenvectors via Inverse Iteration" (inverse iteration on the Hessenberg form, v = Q·z) and §7.4.3 (Hessenberg reduction with Q accumulated)
// Source:  Anderson et al., LAPACK Users' Guide, 3rd ed., SIAM, 1999, routine DHSEIN/DLAEIN (pivot floor ε₃ = ‖H‖·ulp and perturbation of close eigenvalues)
// Tags:    matrix, n-dimensional, eigenvectors, eigenvalues, nonsymmetric, hessenberg, inverse-iteration, complex, linear-algebra
// Units:   a in unit A; eigenvalues in A; eigenvectors dimensionless (unit length)
// Example: [2 1; 0 3] → λ = 2, 3 with v = (1, 0) and (1, 1)/√2; [0 −1; 1 0] → λ = ±i with v = (1, ∓i)/√2 (hand-checked)
func MatEigenvectors(valuesDst, vectorsDst []complex128, a []float64, n int) (values, vectors []complex128, err error) {
	values, err = MatEigenvalues(valuesDst, a, n)
	if err != nil {
		return nil, nil, err
	}
	big, _ := maxAbs(a)
	_, e := pow2Scale(max(big, 0x1p-1000))
	work := make([]float64, 2*n*n)
	h, q := work[:n*n], work[n*n:]
	for i, x := range a {
		h[i] = math.Ldexp(x, -e)
	}
	for i := range n {
		q[i*n+i] = 1
	}
	hessenberg(h, q, n)
	hnorm := 0.0
	for j := range n {
		s := 0.0
		for i := range min(j+2, n) {
			s += math.Abs(h[i*n+j])
		}
		hnorm = max(hnorm, s)
	}
	eps3 := hnorm * 0x1p-52
	if eps3 == 0 {
		eps3 = 0x1p-52
	}
	clusterTol := hnorm * 0x1p-26
	cw := make([]complex128, n*n+2*n+n*n)
	lu, b, mult, z := cw[:n*n], cw[n*n:n*n+n], cw[n*n+n:n*n+2*n], cw[n*n+2*n:] // z: column k = zₖ (row-major n×n)
	swap := make([]bool, n)
	mu := make([]complex128, n)
	for k := range n {
		mu[k] = complex(math.Ldexp(real(values[k]), -e), math.Ldexp(imag(values[k]), -e))
		if k > 0 && imag(values[k]) != 0 && values[k] == cmplx.Conj(values[k-1]) {
			for i := range n {
				z[i*n+k] = cmplx.Conj(z[i*n+k-1])
			}
			continue
		}
		var mates []int
		for j := range k {
			if cmplx.Abs(mu[j]-mu[k]) <= clusterTol {
				mates = append(mates, j)
			}
		}
		shift := mu[k] + complex(float64(len(mates))*eps3, 0)
		hessenbergLU(lu, swap, mult, h, n, shift, eps3)
		for i := range n {
			b[i] = 1
			if len(mates) > 0 {
				b[i] = complex(math.Cos(float64(len(mates)*(i+1))), 0)
			}
		}
		for range 3 {
			hessenbergSolve(lu, swap, mult, b, n)
			before := cvecNorm(b)
			if !reduce.AllFinite(before) || before == 0 {
				return nil, nil, fmt.Errorf("%w: eigenvector %d", ErrNoConvergence, k)
			}
			cvecScale(b, complex(1/before, 0))
			if len(mates) > 0 {
				save := append([]complex128(nil), b...)
				for _, j := range mates {
					var dot complex128
					for i := range n {
						dot += cmplx.Conj(z[i*n+j]) * b[i]
					}
					for i := range n {
						b[i] -= dot * z[i*n+j]
					}
				}
				if after := cvecNorm(b); after < 1e-3 {
					copy(b, save) // defective: the eigenspace has no room for another vector
				} else {
					cvecScale(b, complex(1/after, 0))
				}
			}
		}
		for i := range n {
			z[i*n+k] = b[i]
		}
	}
	if cap(vectorsDst) >= n*n {
		vectorsDst = vectorsDst[:n*n]
	} else {
		vectorsDst = make([]complex128, n*n)
	}
	// v = Q·z, column by column; then unit norm and the largest entry real positive.
	for i := range n {
		for k := range n {
			var s complex128
			for j := range n {
				s += complex(q[i*n+j], 0) * z[j*n+k]
			}
			vectorsDst[i*n+k] = s
		}
	}
	for k := range n {
		norm, bigAbs, bigAt := 0.0, -1.0, 0
		for i := range n {
			x := vectorsDst[i*n+k]
			norm = math.Hypot(norm, cmplx.Abs(x))
			if m := cmplx.Abs(x); m > bigAbs*(1+1e-12) {
				bigAbs, bigAt = m, i
			}
		}
		if norm == 0 || !reduce.AllFinite(norm) {
			return nil, nil, fmt.Errorf("%w: eigenvector %d", ErrNoConvergence, k)
		}
		x := vectorsDst[bigAt*n+k]
		phase := cmplx.Conj(x) / complex(cmplx.Abs(x)*norm, 0)
		for i := range n {
			vectorsDst[i*n+k] *= phase
		}
		vectorsDst[bigAt*n+k] = complex(real(vectorsDst[bigAt*n+k]), 0)
		if imag(values[k]) == 0 {
			for i := range n {
				vectorsDst[i*n+k] = complex(real(vectorsDst[i*n+k]), 0)
			}
		}
	}
	return values, vectorsDst, nil
}

// hessenbergLU factors H − shift·I (H upper Hessenberg, real) into lu by Gaussian elimination with
// pivoting between adjacent rows; a pivot smaller than floor in modulus is replaced by floor (DLAEIN).
func hessenbergLU(lu []complex128, swap []bool, mult []complex128, h []float64, n int, shift complex128, floor float64) {
	for i := range n {
		for j := range n {
			lu[i*n+j] = 0
			if j >= i-1 {
				lu[i*n+j] = complex(h[i*n+j], 0)
			}
		}
		lu[i*n+i] -= shift
	}
	for k := 0; k+1 < n; k++ {
		swap[k] = cmplx.Abs(lu[(k+1)*n+k]) > cmplx.Abs(lu[k*n+k])
		if swap[k] {
			for j := k; j < n; j++ {
				lu[k*n+j], lu[(k+1)*n+j] = lu[(k+1)*n+j], lu[k*n+j]
			}
		}
		if cmplx.Abs(lu[k*n+k]) < floor {
			lu[k*n+k] = complex(floor, 0)
		}
		l := lu[(k+1)*n+k] / lu[k*n+k]
		mult[k] = l
		lu[(k+1)*n+k] = 0
		for j := k + 1; j < n; j++ {
			lu[(k+1)*n+j] -= l * lu[k*n+j]
		}
	}
	if cmplx.Abs(lu[n*n-1]) < floor {
		lu[n*n-1] = complex(floor, 0)
	}
}

// hessenbergSolve overwrites b with the solution of (H − shift·I)·x = b from hessenbergLU. The back
// substitution rescales by 2⁻⁵⁰⁰ whenever an entry passes 2⁵⁰⁰, so tiny pivots cannot overflow (only
// the direction of x is used).
func hessenbergSolve(lu []complex128, swap []bool, mult []complex128, b []complex128, n int) {
	for k := 0; k+1 < n; k++ {
		if swap[k] {
			b[k], b[k+1] = b[k+1], b[k]
		}
		b[k+1] -= mult[k] * b[k]
	}
	const huge, tiny = 0x1p500, 0x1p-500
	for i := n - 1; i >= 0; i-- {
		s := b[i]
		for j := i + 1; j < n; j++ {
			s -= lu[i*n+j] * b[j]
		}
		x := s / lu[i*n+i]
		if cmplx.Abs(x) > huge {
			x *= tiny
			for j := range n {
				if j != i {
					b[j] *= tiny
				}
			}
		}
		b[i] = x
	}
}

func cvecNorm(v []complex128) float64 {
	big := 0.0
	for _, x := range v {
		big = max(big, cmplx.Abs(x))
	}
	if big == 0 || !reduce.AllFinite(big) {
		return big
	}
	s := 0.0
	for _, x := range v {
		r := cmplx.Abs(x) / big
		s += r * r
	}
	return big * math.Sqrt(s)
}

func cvecScale(v []complex128, s complex128) {
	for i := range v {
		v[i] *= s
	}
}
