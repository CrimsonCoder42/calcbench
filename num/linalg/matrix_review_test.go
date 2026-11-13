package linalg

import (
	"errors"
	"math"
	"testing"
)

// Review fixes: QR and least squares at extreme scales return the scaled exact answers (no silent skip of a
// reflection, no false overflow).
func TestQR_extremeScales(t *testing.T) {
	for _, s := range []float64{1e-170, 1e160, 1e-300, 1e300} {
		a := []float64{3 * s, 0, 4 * s, 5 * s}
		q, r, err := QRHouseholder(nil, nil, a, 2, 2)
		if err != nil {
			t.Fatalf("scale %g: %v", s, err)
		}
		if !eqSlice(r, []float64{-5 * s, -4 * s, 0, -3 * s}, 1e-14) || orthoErr(q, 2, 2) > 1e-15 {
			t.Errorf("scale %g: R = %v", s, r)
		}
		rs := []float64{r[0] / s, r[1] / s, r[2] / s, r[3] / s}
		qr, _ := MatMul(nil, q, 2, 2, rs, 2)
		if !eqSlice(qr, []float64{3, 0, 4, 5}, 1e-14) {
			t.Errorf("scale %g: QR ≠ A", s)
		}
		res, err := LeastSquaresQR(nil, []float64{s, 0, s, s, s, 2 * s}, 3, 2, []float64{6 * s, 0, 0}, -1)
		if err != nil || !eqSlice(res.X, []float64{5, -3}, 1e-13) || !near(res.ResidualNorm, math.Sqrt(6)*s, 1e-13) {
			t.Errorf("scale %g: LS = %+v, %v", s, res, err)
		}
	}
	// A tiny column next to large ones: the reflector of the tiny column is still applied.
	a := []float64{1, 1e-200, 1, 0, 0, 1e-200}
	q, r, err := QRHouseholder(nil, nil, a, 3, 2)
	if err != nil {
		t.Fatal(err)
	}
	qr, _ := MatMul(nil, q, 3, 2, r, 2)
	if !eqSlice(qr, a, 1e-14) || orthoErr(q, 3, 2) > 1e-14 {
		t.Errorf("mixed scales: QR = %v", qr)
	}
}

// Review fix: Mat2Eigen separates nearly equal eigenvalues and handles huge entries.
func TestMat2Eigen_accuracy(t *testing.T) {
	e, err := Mat2Eigen(Mat2{1, 1, 1e-17, 1})
	if err != nil || e.Defective {
		t.Fatalf("near-defective: %+v, %v", e, err)
	}
	d := math.Sqrt(1e-17)
	if !near(real(e.Values[0]), 1-d, 1e-15) || !near(real(e.Values[1]), 1+d, 1e-15) {
		t.Errorf("λ = %v, want 1 ± %g", e.Values, d)
	}
	h, err := Mat2Eigen(Mat2{1e300, 1e300, -1e300, 1e300})
	if err != nil || !near(real(h.Values[1]), 1e300, 1e-15) || !near(imag(h.Values[1]), 1e300, 1e-15) {
		t.Errorf("huge: %v, %v", h.Values, err)
	}
	if _, err := Mat2Eigen(Mat2{math.MaxFloat64, math.MaxFloat64, math.MaxFloat64, math.MaxFloat64}); !errors.Is(err, ErrOverflow) {
		t.Errorf("λ = 2·MaxFloat64: %v", err)
	}
}

// Review fix: the PSD check passes its tolerance to the eigenvalue sweep, so a matrix accepted as symmetric
// is never refused there.
func TestMatIsPositiveSemidefinite_tolerance(t *testing.T) {
	ok, err := MatIsPositiveSemidefinite([]float64{1, 1 + 1e-9, 1, 1}, 2, 1e-6)
	if err != nil || !ok {
		t.Errorf("= (%v, %v), want (true, nil)", ok, err)
	}
}

// Review fix: an impossible covariance matrix is ErrOverflow, never a silent Inf.
func TestCovarianceToCorrelation_overflow(t *testing.T) {
	if _, _, err := CovarianceToCorrelation(nil, nil, []float64{1e-300, 1e300, 1e300, 1e-300}, 2); !errors.Is(err, ErrOverflow) {
		t.Errorf("err = %v, want ErrOverflow", err)
	}
}

// Review fix: U has orthonormal columns for every nonzero singular value, tall and wide.
func TestSVD_leftVectors(t *testing.T) {
	r := newRNG()
	for range 40 {
		m, n := 1+r.IntN(6), 1+r.IntN(6)
		s, err := SVDJacobi(randMat(r, m, n), m, n, -1, true)
		if err != nil {
			t.Fatal(err)
		}
		if e := orthoErr(s.U, m, min(m, n)); e > 1e-12 {
			t.Fatalf("UᵀU ≠ I (%d×%d): %v", m, n, e)
		}
	}
}

// Review fix: the echelon form has the staircase shape (zeros below each pivot, pivots move right), and
// Mat3CharPoly agrees with det and trace computed independently.
func TestEchelon_staircaseAndCharPoly(t *testing.T) {
	r := newRNG()
	for range 60 {
		m, c := 1+r.IntN(6), 1+r.IntN(6)
		k := 1 + r.IntN(min(m, c))
		a, _ := MatMul(nil, randMat(r, m, k), m, k, randMat(r, k, c), c)
		e, piv, err := MatRowEchelon(nil, a, m, c, 1e-10)
		if err != nil {
			t.Fatal(err)
		}
		for row, p := range piv {
			if row > 0 && p <= piv[row-1] {
				t.Fatalf("pivots not increasing: %v", piv)
			}
			for i := row + 1; i < m; i++ {
				if e[i*c+p] != 0 {
					t.Fatalf("nonzero below pivot (%d, %d)", i, p)
				}
			}
			for j := 0; j < p; j++ {
				if e[row*c+j] != 0 {
					t.Fatalf("nonzero left of pivot in row %d", row)
				}
			}
		}
		for row := len(piv); row < m; row++ {
			for j := range c {
				if e[row*c+j] != 0 {
					t.Fatalf("nonzero in a row past the rank")
				}
			}
		}
	}
	for range 100 {
		m3 := randM3(r)
		p, _ := Mat3CharPoly(m3)
		d, _ := Mat3Det(m3)
		tr, _ := Mat3Trace(m3)
		if !near(p[0], -d, 1e-14) || !near(p[2], -tr, 1e-15) || p[3] != 1 {
			t.Fatalf("Mat3CharPoly %v vs det %v trace %v", p, d, tr)
		}
	}
}

// Review fix: the characteristics bundles allocate nothing.
func TestMatCharacteristics_noAllocation(t *testing.T) {
	m2, m3 := Mat2{4, 1, 2, 3}, Mat3{1, 2, 3, 0, 1, 4, 5, 6, 0}
	if n := testing.AllocsPerRun(100, func() { _, _ = Mat2Characteristics(m2, -1) }); n != 0 {
		t.Errorf("Mat2Characteristics allocs = %v", n)
	}
	if n := testing.AllocsPerRun(100, func() { _, _ = Mat3Characteristics(m3, -1) }); n != 0 {
		t.Errorf("Mat3Characteristics allocs = %v", n)
	}
	// rankFixed agrees with MatRank.
	r := newRNG()
	for range 100 {
		x := randM3(r)
		if r.IntN(3) == 0 {
			for j := range 3 {
				x[6+j] = x[j] + x[3+j]
			}
		}
		a, _ := rankFixed(x[:], 3, -1)
		b, _ := MatRank(x[:], 3, 3, -1)
		if a != b {
			t.Fatalf("rankFixed %d vs MatRank %d", a, b)
		}
	}
}
