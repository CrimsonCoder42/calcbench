package invest

import "testing"

// The SEC staff paper defines the relative quoted spread without a worked number; hand-checked:
// bid $24.19, ask $24.21 → (2/((2419 + 2421)/2)) = 4/4840 = 0.000826446.
func TestBidAskSpread_handChecked(t *testing.T) {
	got, err := BidAskSpread(2419, 2421)
	if err != nil || !near(got, 4.0/4840.0, 1e-18) {
		t.Errorf("BidAskSpread = %v, %v; want %v", got, err, 4.0/4840.0)
	}
}

func TestBidAskSpread_cases(t *testing.T) {
	for _, c := range []struct {
		name     string
		bid, ask int64
		want     float64
	}{
		{"zero bid", 0, 100, 2},
		{"one cent wide", 99, 101, 0.02},
		{"extremes", maxCents - 2, maxCents, 4 / float64(2*maxCents-2)},
	} {
		got, err := BidAskSpread(c.bid, c.ask)
		if err != nil || !near(got, c.want, 1e-15) {
			t.Errorf("%s: %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestBidAskSpread_errors(t *testing.T) {
	_, err := BidAskSpread(100, 100)
	wantErr(t, "locked", err, ErrDomain)
	_, err = BidAskSpread(101, 100)
	wantErr(t, "crossed", err, ErrDomain)
	_, err = BidAskSpread(-1, 100)
	wantErr(t, "bid < 0", err, ErrNegative)
	_, err = BidAskSpread(1, maxCents+1)
	wantErr(t, "ask > 2^53", err, ErrOutOfRange)
}

func TestBidAskSpread_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = BidAskSpread(2419, 2421) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkBidAskSpread(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkErr = BidAskSpread(2419, 2421)
	}
}
