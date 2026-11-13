package market

import "testing"

var sinkInt int64

// Hand-checked: bid 2419, ask 2421 → 2 ticks (exact integer comparison).
func TestSpreadTicks_cases(t *testing.T) {
	for _, c := range []struct {
		name           string
		bid, ask, want int64
	}{
		{"example", 2419, 2421, 2},
		{"locked", 7, 7, 0},
		{"full range", 0, maxTicks, maxTicks},
	} {
		got, err := SpreadTicks(c.bid, c.ask)
		if err != nil || got != c.want {
			t.Errorf("%s: %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestSpreadTicks_errors(t *testing.T) {
	_, err := SpreadTicks(3, 2)
	wantErr(t, "crossed", err, ErrDomain)
	_, err = SpreadTicks(-1, 2)
	wantErr(t, "negative", err, ErrNegative)
	_, err = SpreadTicks(0, maxTicks+1)
	wantErr(t, "too large", err, ErrOutOfRange)
}

func BenchmarkSpreadTicks(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, sinkErr = SpreadTicks(2419, 2421)
	}
}
