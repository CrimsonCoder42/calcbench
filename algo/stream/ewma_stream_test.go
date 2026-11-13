package stream

import (
	"errors"
	"math"
	"testing"

	"github.com/CrimsonCoder42/calcbench/algo/window"
)

// Hand-derived, α = 0.5: 2, 4, 8 smooth to 2, 0.5·4 + 0.5·2 = 3, 0.5·8 + 0.5·3 = 5.5.
func TestEWMAStream_example(t *testing.T) {
	e, err := NewEWMAStream(0.5)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := e.Value(); ok {
		t.Error("empty: Value defined")
	}
	for i, x := range []float64{2, 4, 8} {
		got, err := e.Add(x)
		if want := []float64{2, 3, 5.5}[i]; err != nil || got != want {
			t.Errorf("Add(%v) = (%v, %v), want %v", x, got, err, want)
		}
	}
	if v, ok := e.Value(); !ok || v != 5.5 || e.Alpha() != 0.5 {
		t.Errorf("Value = (%v, %v), Alpha %v", v, ok, e.Alpha())
	}
	e.Reset()
	if _, ok := e.Value(); ok || e.Alpha() != 0.5 {
		t.Error("Reset kept the value or lost alpha")
	}
	if got, _ := e.Add(-7); got != -7 {
		t.Errorf("first Add after Reset = %v, want -7", got)
	}
}

// Agreement with ALG-WINDOW-014 ExponentialMovingAverage on the same 1,000-value series, value for value.
func TestEWMAStream_matchesExponentialMovingAverage(t *testing.T) {
	values := benchFloats(1000)
	for _, alpha := range []float64{1, 0.5, 0.1, 0.01, 1.0 / 3} {
		batch, err := window.ExponentialMovingAverage(nil, values, alpha)
		if err != nil {
			t.Fatal(err)
		}
		e, _ := NewEWMAStream(alpha)
		for i, x := range values {
			got, err := e.Add(x)
			// Same expression term for term, so equal; tolerance 1e-15 relative (of max(1, |want|)) only
			// because the Go spec lets the compiler fuse a multiply-add differently in two places.
			if err != nil || !closeTo(got, batch[i], 1e-15) {
				t.Fatalf("alpha %v step %d: stream %v, batch %v (%v)", alpha, i, got, batch[i], err)
			}
		}
	}
}

func TestEWMAStream_errors(t *testing.T) {
	for _, a := range []float64{0, -0.1, 1.5, math.NaN(), math.Inf(1)} {
		if _, err := NewEWMAStream(a); !errors.Is(err, ErrAlpha) {
			t.Errorf("NewEWMAStream(%v) error = %v, want ErrAlpha", a, err)
		}
	}
	var zero EWMAStream
	if _, err := zero.Add(1); !errors.Is(err, ErrUninitialized) {
		t.Errorf("zero.Add error = %v, want ErrUninitialized", err)
	}
	e, _ := NewEWMAStream(0.25)
	_, _ = e.Add(4)
	for _, x := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := e.Add(x); !errors.Is(err, ErrNotFinite) {
			t.Errorf("Add(%v) error = %v, want ErrNotFinite", x, err)
		}
	}
	if v, _ := e.Value(); v != 4 {
		t.Errorf("state changed on error: Value %v, want 4", v)
	}
	// ErrOverflow is documented defensively: a convex combination of finite values exceeds the range only
	// by rounding, and no input pair at MaxFloat64 was found that does (same check as ALG-WINDOW-014).
}

func TestEWMAStream_noAllocation(t *testing.T) {
	e, _ := NewEWMAStream(0.1)
	values := benchFloats(1000)
	allocs := testing.AllocsPerRun(20, func() {
		for _, x := range values {
			_, _ = e.Add(x)
		}
	})
	if allocs != 0 {
		t.Errorf("Add allocates %v times per 1000 values, want 0", allocs)
	}
}

func BenchmarkEWMAStream(b *testing.B) {
	for _, size := range benchSizes {
		b.Run(size.name, func(b *testing.B) {
			values := benchFloats(size.n)
			e, _ := NewEWMAStream(0.1)
			b.ReportAllocs()
			for b.Loop() {
				for _, x := range values {
					_, _ = e.Add(x)
				}
			}
		})
	}
}
