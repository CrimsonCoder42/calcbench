package critical

import (
	"errors"
	"math"
	"testing"
)

type scortenIn struct {
	age        float64
	malignancy bool
	hr, bsa    float64
	urea, gluc float64
	bicarb     float64
}

func (s scortenIn) run() (int, error) {
	return SCORTEN(s.age, s.malignancy, s.hr, s.bsa, s.urea, s.gluc, s.bicarb)
}

// HAND-DERIVED from Bastuji-Garin 2000, abstract: seven factors, each "above"/"below" its threshold, one point
// each. At the threshold itself a factor is absent; just past it, present.
func TestSCORTEN_example(t *testing.T) {
	none := scortenIn{40, false, 120, 10, 10, 14, 20}
	all := scortenIn{41, true, 121, 11, 10.1, 14.1, 19.9}
	if got, err := none.run(); err != nil || got != 0 {
		t.Errorf("all at threshold = (%d, %v), want 0", got, err)
	}
	if got, err := all.run(); err != nil || got != 7 {
		t.Errorf("all past threshold = (%d, %v), want 7", got, err)
	}
	one := []func(*scortenIn){
		func(s *scortenIn) { s.age = 40.01 }, func(s *scortenIn) { s.malignancy = true },
		func(s *scortenIn) { s.hr = 120.5 }, func(s *scortenIn) { s.bsa = 10.5 },
		func(s *scortenIn) { s.urea = 10.01 }, func(s *scortenIn) { s.gluc = 14.01 },
		func(s *scortenIn) { s.bicarb = 19.99 },
	}
	for i, f := range one {
		in := none
		f(&in)
		if got, err := in.run(); err != nil || got != 1 {
			t.Errorf("factor %d alone = (%d, %v), want 1", i, got, err)
		}
	}
}

func TestSCORTEN_errors(t *testing.T) {
	ok := scortenIn{50, false, 90, 20, 5, 6, 24}
	cases := []struct {
		name string
		mod  func(*scortenIn)
		want error
	}{
		{"age NaN", func(s *scortenIn) { s.age = math.NaN() }, ErrNotFinite},
		{"hr +Inf", func(s *scortenIn) { s.hr = math.Inf(1) }, ErrNotFinite},
		{"bsa -Inf", func(s *scortenIn) { s.bsa = math.Inf(-1) }, ErrNotFinite},
		{"urea NaN", func(s *scortenIn) { s.urea = math.NaN() }, ErrNotFinite},
		{"glucose NaN", func(s *scortenIn) { s.gluc = math.NaN() }, ErrNotFinite},
		{"bicarb NaN", func(s *scortenIn) { s.bicarb = math.NaN() }, ErrNotFinite},
		{"age -1", func(s *scortenIn) { s.age = -1 }, ErrOutOfRange},
		{"age 131", func(s *scortenIn) { s.age = 131 }, ErrOutOfRange},
		{"hr -1", func(s *scortenIn) { s.hr = -1 }, ErrOutOfRange},
		{"bsa 101", func(s *scortenIn) { s.bsa = 101 }, ErrOutOfRange},
		{"bsa -1", func(s *scortenIn) { s.bsa = -1 }, ErrOutOfRange},
		{"urea -1", func(s *scortenIn) { s.urea = -1 }, ErrOutOfRange},
		{"glucose -1", func(s *scortenIn) { s.gluc = -1 }, ErrOutOfRange},
		{"bicarb -1", func(s *scortenIn) { s.bicarb = -1 }, ErrOutOfRange},
	}
	for _, c := range cases {
		in := ok
		c.mod(&in)
		if _, err := in.run(); !errors.Is(err, c.want) {
			t.Errorf("%s: error = %v, want %v", c.name, err, c.want)
		}
	}
	for _, in := range []scortenIn{{0, false, 0, 0, 0, 0, 0}, {130, true, 300, 100, 50, 50, 50}} {
		if _, err := in.run(); err != nil {
			t.Errorf("%+v: error = %v, want nil (range edge)", in, err)
		}
	}
}

func TestSCORTEN_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = SCORTEN(55, true, 125, 30, 12, 8, 18) }); n != 0 {
		t.Errorf("SCORTEN allocates %v times, want 0", n)
	}
}

func BenchmarkSCORTEN(b *testing.B) {
	b.ReportAllocs()
	age, hr, bsa, urea, gluc, bic := 55.0, 125.0, 30.0, 12.0, 8.0, 18.0
	for b.Loop() {
		sinkInt, _ = SCORTEN(age, true, hr, bsa, urea, gluc, bic)
	}
}
