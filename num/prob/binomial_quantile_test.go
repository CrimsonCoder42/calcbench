package prob

import "testing"

// JKK ch. 3: n = 10, p = 0.5: F(4) = 386/1024 < 0.5 ≤ F(5) = 638/1024 → 5. Hand-checked.
func TestBinomialQuantile_sourceExample(t *testing.T) {
	if got, err := BinomialQuantile(0.5, 10, 0.5); err != nil || got != 5 {
		t.Errorf("BinomialQuantile(0.5, 10, 0.5) = (%v, %v), want 5", got, err)
	}
}

func TestBinomialQuantile_cases(t *testing.T) {
	cases := []struct {
		name string
		u    float64
		n    int
		p    float64
		want int
	}{
		{"u = 0", 0, 10, 0.3, 0},
		{"u = 1", 1, 10, 0.3, 10},
		{"n = 0", 0.7, 0, 0.3, 0},
		{"p = 0", 0.9, 10, 0, 0},
		{"p = 1", 0.1, 10, 1, 10},
		{"exactly F(5)", 638.0 / 1024, 10, 0.5, 5},
	}
	for _, c := range cases {
		if got, err := BinomialQuantile(c.u, c.n, c.p); err != nil || got != c.want {
			t.Errorf("%s: (%v, %v), want %v", c.name, got, err, c.want)
		}
	}
}

// Least k: F(k − 1) < u ≤ F(k) on a grid of u.
func TestBinomialQuantile_leastK(t *testing.T) {
	for _, n := range []int{1, 7, 100, 5000} {
		for _, u := range []float64{1e-12, 0.05, 0.5, 0.95, 1 - 1e-12} {
			k, err := BinomialQuantile(u, n, 0.37)
			hi, _ := BinomialCDF(k, n, 0.37)
			lo := 0.0
			if k > 0 {
				lo, _ = BinomialCDF(k-1, n, 0.37)
			}
			if err != nil || !(lo < u && u <= hi) {
				t.Errorf("n %v u %v: k %v, F(k−1) %v, F(k) %v, err %v", n, u, k, lo, hi, err)
			}
		}
	}
}

func TestBinomialQuantile_errors(t *testing.T) {
	_, err := BinomialQuantile(1.1, 10, 0.5)
	requireErr(t, "u > 1", err, ErrProbability)
	_, err = BinomialQuantile(0.5, -1, 0.5)
	requireErr(t, "n < 0", err, ErrOutOfRange)
	_, err = BinomialQuantile(nan, 10, 0.5)
	requireErr(t, "u NaN", err, ErrNotFinite)
}

func BenchmarkBinomialQuantile(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = BinomialQuantile(0.9, 1000, 0.3)
	}
}
