package signal

import "testing"

// TestBartlett_formula checks 1 − |2n/(M−1) − 1| (Harris 1978 eq. (20)) at every point for M = 7.
func TestBartlett_formula(t *testing.T) {
	w, err := Bartlett(nil, 7)
	if err != nil {
		t.Fatal(err)
	}
	equalFloats(t, "Bartlett(7)", w, []float64{0, 1.0 / 3, 2.0 / 3, 1, 2.0 / 3, 1.0 / 3, 0}, 1e-15)
}

func BenchmarkBartlett(b *testing.B) { benchWindow(b, Bartlett) }
