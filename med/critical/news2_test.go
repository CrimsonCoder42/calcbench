package critical

import (
	"errors"
	"math"
	"testing"
)

var (
	sinkInt    int
	sinkNEWS2  NEWS2Scores
	sinkString string
)

type news2In struct {
	rr, spo2       int
	scale2, oxygen bool
	sbp, pulse     int
	alert          bool
	tempC          float64
}

func (in news2In) run() (NEWS2Scores, error) {
	return NEWS2(in.rr, in.spo2, in.scale2, in.oxygen, in.sbp, in.pulse, in.alert, in.tempC)
}

// HAND-DERIVED from RCP 2017 Chart 1: every printed band of every parameter, both sides of each edge.
func TestNEWS2_bandEdges(t *testing.T) {
	base := news2In{rr: 16, spo2: 98, sbp: 120, pulse: 70, alert: true, tempC: 37}
	type tc struct {
		name  string
		mod   func(*news2In)
		field func(NEWS2Scores) int
		want  int
	}
	rr := func(s NEWS2Scores) int { return s.RespirationRate }
	sp := func(s NEWS2Scores) int { return s.SpO2 }
	bp := func(s NEWS2Scores) int { return s.SystolicBP }
	hr := func(s NEWS2Scores) int { return s.Pulse }
	te := func(s NEWS2Scores) int { return s.Temperature }
	setRR := func(v int) func(*news2In) { return func(i *news2In) { i.rr = v } }
	setS1 := func(v int) func(*news2In) { return func(i *news2In) { i.spo2 = v } }
	setS2 := func(v int, o bool) func(*news2In) {
		return func(i *news2In) { i.spo2, i.scale2, i.oxygen = v, true, o }
	}
	setBP := func(v int) func(*news2In) { return func(i *news2In) { i.sbp = v } }
	setHR := func(v int) func(*news2In) { return func(i *news2In) { i.pulse = v } }
	setT := func(v float64) func(*news2In) { return func(i *news2In) { i.tempC = v } }
	tests := []tc{
		{"rr 8", setRR(8), rr, 3}, {"rr 9", setRR(9), rr, 1}, {"rr 11", setRR(11), rr, 1}, {"rr 12", setRR(12), rr, 0},
		{"rr 20", setRR(20), rr, 0}, {"rr 21", setRR(21), rr, 2}, {"rr 24", setRR(24), rr, 2}, {"rr 25", setRR(25), rr, 3},
		{"s1 91", setS1(91), sp, 3}, {"s1 92", setS1(92), sp, 2}, {"s1 93", setS1(93), sp, 2}, {"s1 94", setS1(94), sp, 1},
		{"s1 95", setS1(95), sp, 1}, {"s1 96", setS1(96), sp, 0},
		{"s2 83", setS2(83, false), sp, 3}, {"s2 84", setS2(84, false), sp, 2}, {"s2 85", setS2(85, false), sp, 2},
		{"s2 86", setS2(86, false), sp, 1}, {"s2 87", setS2(87, false), sp, 1}, {"s2 88", setS2(88, false), sp, 0},
		{"s2 92 oxygen", setS2(92, true), sp, 0}, {"s2 93 air", setS2(93, false), sp, 0}, {"s2 100 air", setS2(100, false), sp, 0},
		{"s2 93 oxygen", setS2(93, true), sp, 1}, {"s2 94 oxygen", setS2(94, true), sp, 1}, {"s2 95 oxygen", setS2(95, true), sp, 2},
		{"s2 96 oxygen", setS2(96, true), sp, 2}, {"s2 97 oxygen", setS2(97, true), sp, 3},
		{"sbp 90", setBP(90), bp, 3}, {"sbp 91", setBP(91), bp, 2}, {"sbp 100", setBP(100), bp, 2}, {"sbp 101", setBP(101), bp, 1},
		{"sbp 110", setBP(110), bp, 1}, {"sbp 111", setBP(111), bp, 0}, {"sbp 219", setBP(219), bp, 0}, {"sbp 220", setBP(220), bp, 3},
		{"hr 40", setHR(40), hr, 3}, {"hr 41", setHR(41), hr, 1}, {"hr 50", setHR(50), hr, 1}, {"hr 51", setHR(51), hr, 0},
		{"hr 90", setHR(90), hr, 0}, {"hr 91", setHR(91), hr, 1}, {"hr 110", setHR(110), hr, 1}, {"hr 111", setHR(111), hr, 2},
		{"hr 130", setHR(130), hr, 2}, {"hr 131", setHR(131), hr, 3},
		{"t 35.0", setT(35.0), te, 3}, {"t 35.1", setT(35.1), te, 1}, {"t 36.0", setT(36.0), te, 1}, {"t 36.1", setT(36.1), te, 0},
		{"t 38.0", setT(38.0), te, 0}, {"t 38.1", setT(38.1), te, 1}, {"t 39.0", setT(39.0), te, 1}, {"t 39.1", setT(39.1), te, 2},
	}
	for _, c := range tests {
		in := base
		c.mod(&in)
		got, err := in.run()
		if err != nil || c.field(got) != c.want {
			t.Errorf("%s: got %+v, %v; want parameter score %d", c.name, got, err, c.want)
		}
	}
}

// HAND-DERIVED from Chart 1: oxygen adds 2; CVPU adds 3 and sets the red score.
func TestNEWS2_example(t *testing.T) {
	got, err := news2In{rr: 16, spo2: 98, sbp: 120, pulse: 70, alert: true, tempC: 37}.run()
	if err != nil || got != (NEWS2Scores{}) {
		t.Errorf("all-normal = %+v, %v; want zero scores", got, err)
	}
	got, err = news2In{rr: 16, spo2: 98, oxygen: true, sbp: 120, pulse: 70, alert: false, tempC: 37}.run()
	if err != nil || got.AirOrOxygen != 2 || got.Consciousness != 3 || got.Total != 5 || !got.RedScore {
		t.Errorf("oxygen + CVPU = %+v, %v; want AirOrOxygen 2, Consciousness 3, Total 5, RedScore", got, err)
	}
}

// parity: ledger news2, source cases 1–7. The source takes °F; the °C value here is (°F − 32) × 5/9 to
// two decimals, which falls in the same Chart 1 band.
func TestNEWS2_parity(t *testing.T) {
	tests := []struct {
		name string
		in   news2In
		want NEWS2Scores
	}{
		// source case 1 "normal": 98.6 °F = 37.0 °C
		{"normal", news2In{16, 98, false, false, 120, 70, true, 37.0}, NEWS2Scores{0, 0, 0, 0, 0, 0, 0, 0, false}},
		// source case 2 "medium": 100.5 °F = 38.06 °C
		{"medium", news2In{22, 94, false, false, 105, 95, true, 38.06}, NEWS2Scores{2, 1, 0, 1, 1, 0, 1, 6, false}},
		// source case 3 "max": 95 °F = 35.0 °C
		{"max", news2In{26, 90, false, true, 88, 135, false, 35.0}, NEWS2Scores{3, 3, 2, 3, 3, 3, 3, 20, true}},
		// source case 4 "scale 2 target": 90% on air, Scale 2
		{"scale 2 target", news2In{16, 90, true, false, 120, 70, true, 37.0}, NEWS2Scores{0, 0, 0, 0, 0, 0, 0, 0, false}},
		// source case 5 "scale 2 hyperoxia": 97% ON AIR, Scale 2; 102.5 °F = 39.17 °C. The source scores SpO2 3
		// (aggregate 7) because it ignores air vs oxygen (ledger news2); Chart 1 prints "≥93 on air" under 0, so
		// the primary's SpO2 score is 0 and the aggregate 4.
		{"scale 2 hyperoxia", news2In{9, 97, true, false, 215, 45, true, 39.17}, NEWS2Scores{1, 0, 0, 0, 1, 0, 2, 4, false}},
		// source case 6 "scale 2 low": 84% on oxygen; 96.8 °F = 36.0 °C; consciousness not alert
		{"scale 2 low", news2In{7, 84, true, true, 230, 120, false, 36.0}, NEWS2Scores{3, 2, 2, 3, 2, 3, 1, 16, true}},
		// source case 7 "edges": 96.9 °F = 36.06 °C
		{"edges", news2In{12, 93, false, false, 101, 51, true, 36.06}, NEWS2Scores{0, 2, 0, 1, 0, 0, 0, 3, false}},
	}
	for _, c := range tests {
		got, err := c.in.run()
		if err != nil || got != c.want {
			t.Errorf("%s: NEWS2 = %+v, %v; want %+v", c.name, got, err, c.want)
		}
	}
}

func TestNEWS2_errors(t *testing.T) {
	ok := news2In{rr: 16, spo2: 98, sbp: 120, pulse: 70, alert: true, tempC: 37}
	cases := []struct {
		name string
		mod  func(*news2In)
		want error
	}{
		{"temp NaN", func(i *news2In) { i.tempC = math.NaN() }, ErrNotFinite},
		{"temp +Inf", func(i *news2In) { i.tempC = math.Inf(1) }, ErrNotFinite},
		{"temp -Inf", func(i *news2In) { i.tempC = math.Inf(-1) }, ErrNotFinite},
		{"temp 9.99", func(i *news2In) { i.tempC = 9.99 }, ErrOutOfRange},
		{"temp 50.01", func(i *news2In) { i.tempC = 50.01 }, ErrOutOfRange},
		{"rr -1", func(i *news2In) { i.rr = -1 }, ErrOutOfRange},
		{"spo2 -1", func(i *news2In) { i.spo2 = -1 }, ErrOutOfRange},
		{"spo2 101", func(i *news2In) { i.spo2 = 101 }, ErrOutOfRange},
		{"sbp -1", func(i *news2In) { i.sbp = -1 }, ErrOutOfRange},
		{"pulse -1", func(i *news2In) { i.pulse = -1 }, ErrOutOfRange},
	}
	for _, c := range cases {
		in := ok
		c.mod(&in)
		if _, err := in.run(); !errors.Is(err, c.want) {
			t.Errorf("%s: error = %v, want %v", c.name, err, c.want)
		}
	}
	for _, in := range []news2In{
		{0, 0, false, false, 0, 0, true, 10}, {0, 100, true, true, 0, 0, false, 50},
	} {
		if _, err := in.run(); err != nil {
			t.Errorf("%+v: error = %v, want nil (range edge)", in, err)
		}
	}
}

func TestNEWS2_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkNEWS2, _ = NEWS2(22, 94, false, false, 105, 95, true, 38.06) }); n != 0 {
		t.Errorf("NEWS2 allocates %v times, want 0", n)
	}
}

func BenchmarkNEWS2(b *testing.B) {
	b.ReportAllocs()
	rr, spo2, sbp, hr, temp := 22, 94, 105, 95, 38.06
	for b.Loop() {
		sinkNEWS2, _ = NEWS2(rr, spo2, false, false, sbp, hr, true, temp)
	}
}
