package invest

import (
	"fmt"
	"math"
	"slices"
	"testing"
)

// tailReturns sorted: −0.05, −0.03, −0.01, 0, … Confidence 0.8: w = 2 → (0.05 + 0.03)/2 = 0.04.
// Confidence 0.85: w = 1.5 → (0.05 + 0.5·0.03)/1.5 = 0.0433….
func TestExpectedShortfall_handWorked(t *testing.T) {
	got, err := ExpectedShortfall(tailReturns, 0.8)
	if err != nil || !near(got, 0.04, 1e-15) {
		t.Errorf("0.8: %v, %v; want 0.04", got, err)
	}
	got, err = ExpectedShortfall(tailReturns, 0.85)
	if err != nil || !near(got, 0.065/1.5, 1e-15) {
		t.Errorf("0.85: %v, %v; want 0.04333", got, err)
	}
}

func TestExpectedShortfall_cases(t *testing.T) {
	got, err := ExpectedShortfall([]float64{0.02}, 0.99)
	if err != nil || !near(got, -0.02, 1e-15) {
		t.Errorf("one value: %v, %v; want -0.02", got, err)
	}
	// Duplicates at the quantile: −0.02 × 3, 0.01; confidence 0.5 → w = 2 → 0.02.
	got, err = ExpectedShortfall([]float64{-0.02, 0.01, -0.02, -0.02}, 0.5)
	if err != nil || !near(got, 0.02, 1e-15) {
		t.Errorf("ties: %v, %v; want 0.02", got, err)
	}
}

// Oracle: sort, then the source's formula term by term; also ES ≥ HistoricalVaR type 1.
func TestExpectedShortfall_oracle(t *testing.T) {
	r := genReturns(300)
	for _, n := range []int{1, 7, 50, 300} {
		s := slices.Clone(r[:n])
		slices.Sort(s)
		for _, c := range []float64{0.5, 0.9, 0.95, 0.975, 0.99} {
			w := float64(n) * (1 - c)
			m := int(math.Floor(w))
			sum := 0.0
			for i := range m {
				sum += s[i]
			}
			sum += (w - float64(m)) * s[m]
			want := -sum / w
			got, err := ExpectedShortfall(r[:n], c)
			if err != nil || !near(got, want, 1e-14) {
				t.Errorf("n=%d c=%v: %v, %v; want %v", n, c, got, err, want)
			}
			if v := -s[max(0, int(math.Ceil(w))-1)]; got < v-1e-14 {
				t.Errorf("n=%d c=%v: ES %v below VaR %v", n, c, got, v)
			}
		}
	}
}

func TestExpectedShortfall_errors(t *testing.T) {
	_, err := ExpectedShortfall(nil, 0.95)
	wantErr(t, "empty", err, ErrTooFew)
	_, err = ExpectedShortfall([]float64{0.01, math.Inf(-1)}, 0.95)
	wantErr(t, "Inf", err, ErrNotFinite)
	for _, bad := range []float64{0, 1, -1, math.NaN()} {
		_, err = ExpectedShortfall(tailReturns, bad)
		wantErr(t, "bad confidence", err, ErrDomain)
	}
	_, err = ExpectedShortfall([]float64{-1.7e308, -1.7e308, -1.7e308, -1.7e308}, 0.5)
	wantErr(t, "overflow", err, ErrOverflow)
}

func BenchmarkExpectedShortfall(b *testing.B) {
	for _, n := range benchSizes {
		r := genReturns(n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, sinkErr = ExpectedShortfall(r, 0.975)
			}
		})
	}
}
