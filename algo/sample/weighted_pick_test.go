package sample

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// Devroye 1986, ch. III, §2 (inversion): weights [1 3] give index 1 three times as often as index 0.
// 40 000 picks: index 0 within 5σ of 10 000.
func TestWeightedPicker_inversionExample(t *testing.T) {
	p, err := NewWeightedPicker([]int{1, 3})
	if err != nil {
		t.Fatal(err)
	}
	rng := rand.New(rand.NewPCG(1, 2))
	const trials = 40_000
	zero := 0
	for range trials {
		i, err := p.Pick(rng)
		if err != nil {
			t.Fatal(err)
		}
		if i == 0 {
			zero++
		}
	}
	if !within5Sigma(zero, trials, 0.25) {
		t.Errorf("index 0 picked %d of %d, want about %d", zero, trials, trials/4)
	}
}

func TestWeightedPicker_cases(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	tests := []struct {
		name    string
		weights []int64
		allowed map[int]bool
	}{
		{"one weight", []int64{7}, map[int]bool{0: true}},
		{"zero weights skipped", []int64{0, 5, 0, 0, 2, 0}, map[int]bool{1: true, 4: true}},
		{"only last nonzero", []int64{0, 0, 0, 1}, map[int]bool{3: true}},
		{"only first nonzero", []int64{1, 0, 0}, map[int]bool{0: true}},
		{"total at the type maximum", []int64{math.MaxInt64 - 1, 1}, map[int]bool{0: true, 1: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewWeightedPicker(tt.weights)
			if err != nil {
				t.Fatal(err)
			}
			for range 2000 {
				i, err := p.Pick(rng)
				if err != nil || !tt.allowed[i] {
					t.Fatalf("Pick = %d, %v; allowed %v", i, err, tt.allowed)
				}
			}
		})
	}
}

func TestWeightedPicker_unsignedAndSmall(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	p, err := NewWeightedPicker([]uint8{200, 55})
	if err != nil {
		t.Fatal(err)
	}
	for range 1000 {
		if i, err := p.Pick(rng); err != nil || i < 0 || i > 1 {
			t.Fatalf("Pick = %d, %v", i, err)
		}
	}
	if _, err := NewWeightedPicker([]uint8{200, 56}); !errors.Is(err, ErrOverflow) {
		t.Errorf("uint8 total 256: err = %v, want ErrOverflow", err)
	}
}

func TestWeightedPicker_errors(t *testing.T) {
	tests := []struct {
		name    string
		weights []int
		want    error
	}{
		{"empty", nil, ErrZeroTotal},
		{"all zero", []int{0, 0}, ErrZeroTotal},
		{"negative", []int{3, -1}, ErrNegativeWeight},
		{"overflow", []int{math.MaxInt, 1}, ErrOverflow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := NewWeightedPicker(tt.weights); !errors.Is(err, tt.want) {
				t.Errorf("err = %v, want %v", err, tt.want)
			}
		})
	}
	p, _ := NewWeightedPicker([]int{1})
	if _, err := p.Pick(nil); !errors.Is(err, ErrNilRand) {
		t.Errorf("nil rng: err = %v", err)
	}
	var nilP *WeightedPicker[int]
	if _, err := nilP.Pick(rand.New(rand.NewPCG(1, 2))); !errors.Is(err, ErrNilPicker) {
		t.Errorf("nil picker: err = %v", err)
	}
	var zeroP WeightedPicker[int]
	if _, err := zeroP.Pick(rand.New(rand.NewPCG(1, 2))); !errors.Is(err, ErrNilPicker) {
		t.Errorf("zero picker: err = %v", err)
	}
}

func TestWeightedPicker_inputNotModified(t *testing.T) {
	w := []int{3, 1, 2}
	if _, err := NewWeightedPicker(w); err != nil {
		t.Fatal(err)
	}
	if w[0] != 3 || w[1] != 1 || w[2] != 2 {
		t.Errorf("weights modified: %v", w)
	}
}

func TestWeightedPicker_pickNoAllocation(t *testing.T) {
	p, _ := NewWeightedPicker([]int{1, 2, 3, 4})
	rng := rand.New(rand.NewPCG(1, 2))
	if allocs := testing.AllocsPerRun(100, func() { benchSink, _ = p.Pick(rng) }); allocs != 0 {
		t.Errorf("Pick allocates %v times per call, want 0", allocs)
	}
}

// TestWeightedPicker_oracle: the slow version maps r to an index by a linear scan of the cumulative
// weights; both see the same r stream (two generators with the same seed), so they must agree exactly.
func TestWeightedPicker_oracle(t *testing.T) {
	gen := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		w := make([]int, 1+gen.IntN(20))
		total := 0
		for i := range w {
			if gen.IntN(3) > 0 {
				w[i] = gen.IntN(10)
			}
			total += w[i]
		}
		if total == 0 {
			w[0], total = 1, 1
		}
		p, err := NewWeightedPicker(w)
		if err != nil {
			t.Fatal(err)
		}
		a := rand.New(rand.NewPCG(uint64(trial), 9))
		b := rand.New(rand.NewPCG(uint64(trial), 9))
		for range 50 {
			got, _ := p.Pick(a)
			r := int(b.Uint64N(uint64(total)))
			want, cum := -1, 0
			for i, x := range w {
				cum += x
				if r < cum {
					want = i
					break
				}
			}
			if got != want {
				t.Fatalf("trial %d: weights %v, r=%d: Pick = %d, want %d", trial, w, r, got, want)
			}
		}
	}
}

func BenchmarkWeightedPickerPick(b *testing.B) {
	for _, size := range benchSizes {
		w := make([]int, size.n)
		for i := range w {
			w[i] = 1 + i%7
		}
		p, _ := NewWeightedPicker(w)
		rng := rand.New(rand.NewPCG(1, 2))
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			sink := 0
			for b.Loop() {
				i, _ := p.Pick(rng)
				sink += i
			}
			benchSink = sink
		})
	}
}

func BenchmarkNewWeightedPicker(b *testing.B) {
	for _, size := range benchSizes {
		w := make([]int, size.n)
		for i := range w {
			w[i] = 1 + i%7
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				p, _ := NewWeightedPicker(w)
				benchSink = len(p.prefix)
			}
		})
	}
}
