package invest

import "testing"

// Hand-checked (the OCC Handbook defines the ratio but gives no worked example): NOI = gross rent
// $24,000 − vacancy $1,200 − tax $3,000 − insurance $1,200 − maintenance $2,000 − other $500 = $16,100;
// $16,100 / $200,000 = 0.0805.
func TestCapitalizationRate_handChecked(t *testing.T) {
	got, err := CapitalizationRate(1610000, 20000000)
	noErr(t, "CapitalizationRate", err)
	if !near(got, 0.0805, 1e-15) {
		t.Errorf("CapitalizationRate(1610000, 20000000) = %v, want 0.0805", got)
	}
}

func TestCapitalizationRate_cases(t *testing.T) {
	cases := []struct {
		name       string
		noi, price int64
		want       float64
	}{
		{"zero NOI", 0, 100000, 0},
		{"negative NOI is a negative rate", -500000, 10000000, -0.05},
		{"NOI equals price", 100, 100, 1},
		{"one cent price", 1, 1, 1},
		{"largest exact inputs", maxCents, maxCents, 1},
		{"value from rate inverts", 700000, 10000000, 0.07},
	}
	for _, c := range cases {
		got, err := CapitalizationRate(c.noi, c.price)
		noErr(t, c.name, err)
		if !near(got, c.want, 1e-15) {
			t.Errorf("%s: got %v, want %v", c.name, got, c.want)
		}
	}
}

func TestCapitalizationRate_errors(t *testing.T) {
	_, err := CapitalizationRate(100, 0)
	wantErr(t, "zero price", err, ErrZeroDenominator)
	_, err = CapitalizationRate(100, -1)
	wantErr(t, "negative price", err, ErrNegative)
	_, err = CapitalizationRate(maxCents+1, 100)
	wantErr(t, "NOI beyond 2^53", err, ErrOutOfRange)
	_, err = CapitalizationRate(-maxCents-1, 100)
	wantErr(t, "NOI below −2^53", err, ErrOutOfRange)
	_, err = CapitalizationRate(100, maxCents+1)
	wantErr(t, "price beyond 2^53", err, ErrOutOfRange)
}

func TestCapitalizationRate_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = CapitalizationRate(1610000, 20000000) }); a != 0 {
		t.Errorf("CapitalizationRate allocates %v times, want 0", a)
	}
}

func BenchmarkCapitalizationRate(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkErr = CapitalizationRate(1610000, 20000000)
	}
}
