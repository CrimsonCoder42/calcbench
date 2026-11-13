package business

import (
	"errors"
	"testing"
)

// Principles of Accounting Vol. 1 §11.4, London Hoops: sold for $10 million, net assets $7 million →
// goodwill $3 million. Exact cents.
func TestGoodwillFromAcquisition_sourceExample(t *testing.T) {
	got, err := GoodwillFromAcquisition(1000000000, 700000000)
	if err != nil || got != 300000000 {
		t.Fatalf("got %v, %v; want 300000000", got, err)
	}
}

func TestGoodwillFromAcquisition_table(t *testing.T) {
	cases := []struct {
		name            string
		price, netAsset int64
		want            int64
	}{
		{"price equals net assets", 500, 500, 0},
		{"price below net assets", 400, 500, -100},
		{"negative net assets", 100, -50, 150},
		{"max", maxExactCents, 0, maxExactCents},
		{"min", 0, maxExactCents, -maxExactCents},
	}
	for _, c := range cases {
		got, err := GoodwillFromAcquisition(c.price, c.netAsset)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestGoodwillFromAcquisition_errors(t *testing.T) {
	cases := []struct {
		name            string
		price, netAsset int64
		want            error
	}{
		{"price above 2^53", maxExactCents + 1, 0, ErrOutOfRange},
		{"net assets below -2^53", 0, -maxExactCents - 1, ErrOutOfRange},
		{"negative price", -1, 0, ErrNegative},
		{"result above 2^53", maxExactCents, -1, ErrOverflow},
	}
	for _, c := range cases {
		if _, err := GoodwillFromAcquisition(c.price, c.netAsset); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestGoodwillFromAcquisition_allocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = GoodwillFromAcquisition(1000000000, 700000000) }); n != 0 {
		t.Fatalf("allocs = %v, want 0", n)
	}
}

func BenchmarkGoodwillFromAcquisition(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = GoodwillFromAcquisition(1000000000, 700000000)
	}
}
