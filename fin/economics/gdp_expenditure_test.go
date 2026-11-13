package economics

import (
	"errors"
	"testing"
)

// OpenStax Principles of Macroeconomics 3e, §6.1 Table 6.1 (U.S. 2020, trillions of dollars):
// 14.0 + 3.6 + 3.9 + (2.1 − 2.7) = 20.9. In units of $0.1 trillion: 140 + 36 + 39 + 21 − 27 = 209.
func TestGDPExpenditure_sourceExample(t *testing.T) {
	got, err := GDPExpenditure(140, 36, 39, 21, 27)
	if err != nil || got != 209 {
		t.Fatalf("got %d, %v; want 209", got, err)
	}
}

func TestGDPExpenditure_cases(t *testing.T) {
	cases := []struct {
		name          string
		c, i, g, x, m int64
		want          int64
	}{
		{"all zero", 0, 0, 0, 0, 0, 0},
		{"negative investment", 1000, -200, 300, 0, 0, 1100},
		{"trade deficit larger than the rest", 0, 0, 0, 0, 500, -500},
		{"largest consumption", maxExactCents, 0, 0, 0, 0, maxExactCents},
	}
	for _, c := range cases {
		got, err := GDPExpenditure(c.c, c.i, c.g, c.x, c.m)
		if err != nil || got != c.want {
			t.Errorf("%s: got %d, %v; want %d", c.name, got, err, c.want)
		}
	}
}

func TestGDPExpenditure_errors(t *testing.T) {
	cases := []struct {
		name          string
		c, i, g, x, m int64
		want          error
	}{
		{"negative consumption", -1, 0, 0, 0, 0, ErrNegative},
		{"negative government", 0, 0, -1, 0, 0, ErrNegative},
		{"negative exports", 0, 0, 0, -1, 0, ErrNegative},
		{"negative imports", 0, 0, 0, 0, -1, ErrNegative},
		{"input beyond 2^53", maxExactCents + 1, 0, 0, 0, 0, ErrOutOfRange},
		{"investment below −2^53", 0, -maxExactCents - 1, 0, 0, 0, ErrOutOfRange},
		{"sum beyond 2^53", maxExactCents, 1, 0, 0, 0, ErrOverflow},
	}
	for _, c := range cases {
		if _, err := GDPExpenditure(c.c, c.i, c.g, c.x, c.m); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestGDPExpenditure_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { _, _ = GDPExpenditure(140, 36, 39, 21, 27) }); a != 0 {
		t.Errorf("allocations = %v, want 0", a)
	}
}

func BenchmarkGDPExpenditure(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = GDPExpenditure(140, 36, 39, 21, 27)
	}
}
