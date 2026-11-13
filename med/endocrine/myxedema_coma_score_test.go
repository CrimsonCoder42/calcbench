package endocrine

import "testing"

var myxBase = MyxedemaComaFindings{TemperatureC: 36.5, HeartRateBpm: 75}

// Çetin 2025 (PMC12432182) / Amin 2023 (PMC10503452) table: every row, both sides of each edge.
func TestMyxedemaComaScore_table(t *testing.T) {
	cases := []struct {
		name string
		mod  func(*MyxedemaComaFindings)
		want int
	}{
		{"base", func(f *MyxedemaComaFindings) {}, 0},
		{"temp 35.1", func(f *MyxedemaComaFindings) { f.TemperatureC = 35.1 }, 0},
		{"temp 35", func(f *MyxedemaComaFindings) { f.TemperatureC = 35 }, 10},
		{"temp 32", func(f *MyxedemaComaFindings) { f.TemperatureC = 32 }, 10},
		{"temp 31.9", func(f *MyxedemaComaFindings) { f.TemperatureC = 31.9 }, 20},
		{"CNS somnolent", func(f *MyxedemaComaFindings) { f.CNS = MyxedemaCNSSomnolent }, 10},
		{"CNS obtunded", func(f *MyxedemaComaFindings) { f.CNS = MyxedemaCNSObtunded }, 15},
		{"CNS stupor", func(f *MyxedemaComaFindings) { f.CNS = MyxedemaCNSStupor }, 20},
		{"CNS coma", func(f *MyxedemaComaFindings) { f.CNS = MyxedemaCNSComaOrSeizures }, 30},
		{"GI anorexia", func(f *MyxedemaComaFindings) { f.GI = MyxedemaGIAnorexiaPainConstipation }, 5},
		{"GI motility", func(f *MyxedemaComaFindings) { f.GI = MyxedemaGIDecreasedMotility }, 15},
		{"GI ileus", func(f *MyxedemaComaFindings) { f.GI = MyxedemaGIParalyticIleus }, 20},
		{"precipitating", func(f *MyxedemaComaFindings) { f.PrecipitatingEvent = true }, 10},
		{"HR 60", func(f *MyxedemaComaFindings) { f.HeartRateBpm = 60 }, 0},
		{"HR 59", func(f *MyxedemaComaFindings) { f.HeartRateBpm = 59 }, 10},
		{"HR 50", func(f *MyxedemaComaFindings) { f.HeartRateBpm = 50 }, 10},
		{"HR 49", func(f *MyxedemaComaFindings) { f.HeartRateBpm = 49 }, 20},
		{"HR 40", func(f *MyxedemaComaFindings) { f.HeartRateBpm = 40 }, 20},
		{"HR 39", func(f *MyxedemaComaFindings) { f.HeartRateBpm = 39 }, 30},
		{"EKG", func(f *MyxedemaComaFindings) { f.OtherEKGChanges = true }, 10},
		{"effusion", func(f *MyxedemaComaFindings) { f.PericardialPleural = true }, 10},
		{"pulmonary edema", func(f *MyxedemaComaFindings) { f.PulmonaryEdema = true }, 15},
		{"cardiomegaly", func(f *MyxedemaComaFindings) { f.Cardiomegaly = true }, 15},
		{"hypotension", func(f *MyxedemaComaFindings) { f.Hypotension = true }, 20},
		{"hyponatremia", func(f *MyxedemaComaFindings) { f.Hyponatremia = true }, 10},
		{"hypoglycemia", func(f *MyxedemaComaFindings) { f.Hypoglycemia = true }, 10},
		{"hypoxemia", func(f *MyxedemaComaFindings) { f.Hypoxemia = true }, 10},
		{"hypercarbia", func(f *MyxedemaComaFindings) { f.Hypercarbia = true }, 10},
		{"GFR", func(f *MyxedemaComaFindings) { f.DecreasedGFR = true }, 10},
		{"doc example", func(f *MyxedemaComaFindings) {
			*f = MyxedemaComaFindings{TemperatureC: 33, CNS: MyxedemaCNSObtunded, GI: MyxedemaGIAnorexiaPainConstipation, PrecipitatingEvent: true, HeartRateBpm: 45, Hyponatremia: true}
		}, 70},
		{"maximum", func(f *MyxedemaComaFindings) {
			*f = MyxedemaComaFindings{30, MyxedemaCNSComaOrSeizures, MyxedemaGIParalyticIleus, true, 30, true, true, true, true, true, true, true, true, true, true}
		}, 230},
	}
	for _, c := range cases {
		f := myxBase
		c.mod(&f)
		got, err := MyxedemaComaScore(f)
		if err != nil || got != c.want {
			t.Errorf("%s: got %d, %v; want %d", c.name, got, err, c.want)
		}
	}
}

// parity: ledger myxedema_coma, source case 1: the consumer looks up temperature / CNS / GI / bradycardia
// pickers and adds 10/10/10/15/15/20 and 5 × 10; temperature picker 1 (10), stupor (20), ileus (20),
// HR picker 3 (30), hypotension (20) → 100.
func TestMyxedemaComaScore_parity(t *testing.T) {
	f := MyxedemaComaFindings{TemperatureC: 34, CNS: MyxedemaCNSStupor, GI: MyxedemaGIParalyticIleus, HeartRateBpm: 35, Hypotension: true}
	if got, err := MyxedemaComaScore(f); err != nil || got != 100 {
		t.Errorf("got %d, %v; want 100", got, err)
	}
}

func TestMyxedemaComaScore_errors(t *testing.T) {
	bad := func(mod func(*MyxedemaComaFindings)) func() error {
		return func() error { f := myxBase; mod(&f); _, err := MyxedemaComaScore(f); return err }
	}
	runErrCases(t, []errCase{
		{"temp NaN", bad(func(f *MyxedemaComaFindings) { f.TemperatureC = nan }), ErrNotFinite},
		{"temp +Inf", bad(func(f *MyxedemaComaFindings) { f.TemperatureC = pInf }), ErrNotFinite},
		{"HR -Inf", bad(func(f *MyxedemaComaFindings) { f.HeartRateBpm = nInf }), ErrNotFinite},
		{"temp 19.9", bad(func(f *MyxedemaComaFindings) { f.TemperatureC = 19.9 }), ErrOutOfRange},
		{"temp 45.1", bad(func(f *MyxedemaComaFindings) { f.TemperatureC = 45.1 }), ErrOutOfRange},
		{"HR negative", bad(func(f *MyxedemaComaFindings) { f.HeartRateBpm = -1 }), ErrOutOfRange},
		{"CNS 5", bad(func(f *MyxedemaComaFindings) { f.CNS = 5 }), ErrOutOfRange},
		{"CNS -1", bad(func(f *MyxedemaComaFindings) { f.CNS = -1 }), ErrOutOfRange},
		{"GI 4", bad(func(f *MyxedemaComaFindings) { f.GI = 4 }), ErrOutOfRange},
	})
}

func TestMyxedemaComaScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = MyxedemaComaScore(myxBase) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkMyxedemaComaScore(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkInt, _ = MyxedemaComaScore(myxBase)
	}
}
