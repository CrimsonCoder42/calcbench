package sorting

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// CLRS 3rd ed. §8.4 (Figure 8.4) runs BUCKET-SORT on ten values in [0, 1) with ten buckets. Our own ten
// values of that shape; the expected output is hand-derived.
func TestBucketSort_clrsShapedExample(t *testing.T) {
	values := []float64{0.71, 0.18, 0.42, 0.93, 0.27, 0.15, 0.66, 0.24, 0.05, 0.39}
	want := []float64{0.05, 0.15, 0.18, 0.24, 0.27, 0.39, 0.42, 0.66, 0.71, 0.93}
	got, err := BucketSort(nil, values)
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("BucketSort = (%v, %v), want %v", got, err, want)
	}
}

func TestBucketSort_cases(t *testing.T) {
	below1 := math.Nextafter(1, 0)
	tests := []struct {
		name   string
		values []float64
		want   []float64
	}{
		{"empty", []float64{}, []float64{}},
		{"one element", []float64{0.5}, []float64{0.5}},
		{"two elements reversed", []float64{0.9, 0.1}, []float64{0.1, 0.9}},
		{"duplicates", []float64{0.3, 0.3, 0.1}, []float64{0.1, 0.3, 0.3}},
		{"zero and negative zero", []float64{0.2, math.Copysign(0, -1), 0.1, 0}, []float64{0, 0, 0.1, 0.2}},
		{"largest float below 1 (index clamp)", []float64{below1, 0.5, below1}, []float64{0.5, below1, below1}},
		{"smallest positive", []float64{0.5, math.SmallestNonzeroFloat64}, []float64{math.SmallestNonzeroFloat64, 0.5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BucketSort(nil, tt.values)
			if err != nil || !slices.Equal(got, tt.want) {
				t.Errorf("BucketSort(%v) = (%v, %v), want %v", tt.values, got, err, tt.want)
			}
		})
	}
}

// A crowded bucket (every value in [0, 1/n)) takes the slices.Sort path instead of insertion sort.
func TestBucketSort_crowdedBucket(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	values := make([]float64, 200)
	for i := range values {
		values[i] = rng.Float64() / 1000
	}
	got, err := BucketSort(nil, values)
	if err != nil || !slices.Equal(got, slices.Sorted(slices.Values(values))) {
		t.Errorf("BucketSort(crowded) not sorted: %v", err)
	}
}

func TestBucketSort_errors(t *testing.T) {
	for name, bad := range map[string]float64{
		"one": 1, "above one": 1.5, "negative": -0.1, "NaN": math.NaN(), "+Inf": math.Inf(1), "-Inf": math.Inf(-1),
	} {
		t.Run(name, func(t *testing.T) {
			dst := []float64{9, 9}
			got, err := BucketSort(dst, []float64{0.5, bad})
			if !errors.Is(err, ErrOutOfRange) || !slices.Equal(got, []float64{9, 9}) {
				t.Errorf("BucketSort(%v) = (%v, %v), want (dst unwritten, ErrOutOfRange)", bad, got, err)
			}
		})
	}
}

func TestBucketSort_inPlaceAndInputKept(t *testing.T) {
	values := []float64{0.7, 0.2, 0.5}
	if _, err := BucketSort(nil, values); err != nil || !slices.Equal(values, []float64{0.7, 0.2, 0.5}) {
		t.Errorf("BucketSort modified its input: %v", values)
	}
	if _, err := BucketSort(values, values); err != nil || !slices.Equal(values, []float64{0.2, 0.5, 0.7}) {
		t.Errorf("in place: %v", values)
	}
}

func TestBucketSort_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		n := rng.IntN(300)
		values := make([]float64, n)
		skew := 1 + rng.IntN(4) // skew > 1 crowds the low buckets.
		for i := range values {
			values[i] = math.Pow(rng.Float64(), float64(skew))
		}
		want := slices.Sorted(slices.Values(values))
		got, err := BucketSort(nil, values)
		if err != nil || !slices.Equal(got, want) {
			t.Fatalf("trial %d: BucketSort differs from slices.Sort", trial)
		}
	}
}

func BenchmarkBucketSort(b *testing.B) {
	for _, size := range benchSizes {
		rng := rand.New(rand.NewPCG(1, 2))
		values := make([]float64, size.n)
		for i := range values {
			values[i] = rng.Float64()
		}
		dst := make([]float64, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				out, _ := BucketSort(dst, values)
				benchSink = int(out[0])
			}
		})
	}
}
