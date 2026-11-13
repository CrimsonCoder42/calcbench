package stream

import (
	"errors"
	"math"
	"testing"
)

// Hand-derived, window 2: (10,100),(11,300) → 4300/400 = 10.75; (12,100) → (3300+1200)/400 = 11.25;
// (9,0) → window (12,100),(9,0) → 12; (8,0) → no volume, undefined, total volume exactly 0.
func TestWindowedVWAP_example(t *testing.T) {
	w, err := NewWindowedVWAP(2)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := w.VWAP(); ok {
		t.Error("empty: VWAP defined")
	}
	steps := []struct {
		price, volume, want float64
		ok                  bool
	}{{10, 100, 10, true}, {11, 300, 10.75, true}, {12, 100, 11.25, true}, {9, 0, 12, true}, {8, 0, 0, false}}
	for _, st := range steps {
		if err := w.Add(st.price, st.volume); err != nil {
			t.Fatal(err)
		}
		if g, ok := w.VWAP(); ok != st.ok || g != st.want {
			t.Errorf("after (%v, %v): VWAP (%v, %v), want (%v, %v)", st.price, st.volume, g, ok, st.want, st.ok)
		}
	}
	if w.TotalVolume() != 0 || w.Len() != 2 || w.WindowLen() != 2 {
		t.Errorf("TotalVolume %v Len %d WindowLen %d, want 0 2 2", w.TotalVolume(), w.Len(), w.WindowLen())
	}
	w.Reset()
	if w.Len() != 0 || w.TotalVolume() != 0 {
		t.Error("Reset left trades")
	}
}

// Oracle: at every step of 1,000 fixed-seed trades, against a direct sum over the same window.
func TestWindowedVWAP_oracle(t *testing.T) {
	prices, volumes := tradeSeries(1000)
	for _, n := range []int{1, 2, 7, 50, 1000, 1001} {
		w, _ := NewWindowedVWAP(n)
		for i := range prices {
			if err := w.Add(prices[i], volumes[i]); err != nil {
				t.Fatal(err)
			}
			lo := max(0, i+1-n)
			want, wok := directVWAP(prices[lo:i+1], volumes[lo:i+1])
			got, ok := w.VWAP()
			// Tolerance 1e-12 relative: prices ~100, each window sum within a few ulps (compensated) of
			// exact; the oracle's plain sum within 1000·u.
			if ok != wok || (ok && !relClose(got, want, 1e-12)) {
				t.Fatalf("n=%d step %d: (%v, %v), want (%v, %v)", n, i, got, ok, want, wok)
			}
		}
	}
}

func TestWindowedVWAP_errors(t *testing.T) {
	for _, n := range []int{0, maxWindowLen + 1} {
		if _, err := NewWindowedVWAP(n); !errors.Is(err, ErrWindowLen) {
			t.Errorf("NewWindowedVWAP(%d) error = %v, want ErrWindowLen", n, err)
		}
	}
	var zero WindowedVWAP
	if err := zero.Add(1, 1); !errors.Is(err, ErrUninitialized) {
		t.Errorf("zero.Add error = %v, want ErrUninitialized", err)
	}
	zero.Reset()
	if _, ok := zero.VWAP(); ok || zero.Len() != 0 || zero.WindowLen() != 0 {
		t.Error("zero value: VWAP defined or Len/WindowLen not 0")
	}
	w, _ := NewWindowedVWAP(3)
	_ = w.Add(10, 5)
	for _, tr := range [][2]float64{{math.NaN(), 1}, {1, math.Inf(1)}} {
		if err := w.Add(tr[0], tr[1]); !errors.Is(err, ErrNotFinite) {
			t.Errorf("Add(%v) error = %v, want ErrNotFinite", tr, err)
		}
	}
	if err := w.Add(10, -1); !errors.Is(err, ErrVolume) {
		t.Errorf("Add(10, -1) error = %v, want ErrVolume", err)
	}
	if err := w.Add(1e308, 10); !errors.Is(err, ErrOverflow) {
		t.Errorf("Add(1e308, 10) error = %v, want ErrOverflow", err)
	}
	if p, _ := w.VWAP(); w.Len() != 1 || p != 10 {
		t.Errorf("state changed on error: Len %d VWAP %v", w.Len(), p)
	}
}

func TestWindowedVWAP_noAllocation(t *testing.T) {
	prices, volumes := tradeSeries(1000)
	w, _ := NewWindowedVWAP(50)
	allocs := testing.AllocsPerRun(20, func() {
		for i := range prices {
			_ = w.Add(prices[i], volumes[i])
		}
	})
	if allocs != 0 {
		t.Errorf("Add allocates %v times per 1000 trades, want 0", allocs)
	}
}

func BenchmarkWindowedVWAP(b *testing.B) {
	for _, size := range benchSizes {
		b.Run(size.name, func(b *testing.B) {
			prices, volumes := tradeSeries(size.n)
			w, _ := NewWindowedVWAP(50)
			b.ReportAllocs()
			for b.Loop() {
				for i := range prices {
					_ = w.Add(prices[i], volumes[i])
				}
			}
		})
	}
}
