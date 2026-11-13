package trauma

import "testing"

// Stiell 1992 abstract: the malleolar set and the midfoot set. HAND-DERIVED.
func TestOttawaAnkleRule1992_source(t *testing.T) {
	tests := []struct {
		name      string
		age       float64
		f         OttawaAnkle1992Findings
		mall, mid bool
	}{
		{"malleolar pain, age 55", 55, OttawaAnkle1992Findings{PainNearMalleoli: true}, true, false},
		{"malleolar pain, age 54.9", 54.9, OttawaAnkle1992Findings{PainNearMalleoli: true}, false, false},
		{"posterior tenderness", 30, OttawaAnkle1992Findings{PainNearMalleoli: true, MalleolusPosteriorTenderness: true}, true, false},
		{"cannot bear weight", 30, OttawaAnkle1992Findings{PainNearMalleoli: true, UnableToBearWeight: true}, true, false},
		{"age without malleolar pain", 80, OttawaAnkle1992Findings{}, false, false},
		{"midfoot, fifth metatarsal", 30, OttawaAnkle1992Findings{PainInMidfoot: true, FifthMetatarsalBaseTender: true}, false, true},
		{"midfoot, cuboid", 30, OttawaAnkle1992Findings{PainInMidfoot: true, CuboidTender: true}, false, true},
		{"midfoot, navicular", 30, OttawaAnkle1992Findings{PainInMidfoot: true, NavicularTender: true}, false, true},
		{"midfoot pain only", 30, OttawaAnkle1992Findings{PainInMidfoot: true, UnableToBearWeight: true}, false, false},
		{"navicular without midfoot pain", 30, OttawaAnkle1992Findings{NavicularTender: true}, false, false},
	}
	for _, tc := range tests {
		r, err := OttawaAnkleRule1992(tc.age, tc.f)
		if err != nil || r.MalleolarSet != tc.mall || r.MidfootSet != tc.mid {
			t.Errorf("%s: (%+v, %v), want %v %v", tc.name, r, err, tc.mall, tc.mid)
		}
	}
}

// parity: ledger ottawa_ankle_rule. The app applies the 1993 refined rule (no age item, no cuboid, weight bearing
// also opens the foot series); here the 1992 derivation is asserted: midfoot pain with inability to bear weight
// alone is not in the 1992 midfoot set (app: "Foot x-ray indicated").
func TestOttawaAnkleRule1992_parity(t *testing.T) {
	r, _ := OttawaAnkleRule1992(30, OttawaAnkle1992Findings{PainInMidfoot: true, UnableToBearWeight: true})
	if r.MidfootSet {
		t.Errorf("got %+v, want MidfootSet false (1992 primary; app true)", r)
	}
}

func TestOttawaAnkleRule1992_errors(t *testing.T) {
	for _, v := range nonFinite {
		_, err := OttawaAnkleRule1992(v, OttawaAnkle1992Findings{})
		checkErr(t, "non-finite", err, ErrNotFinite)
	}
	_, err := OttawaAnkleRule1992(-1, OttawaAnkle1992Findings{})
	checkErr(t, "range", err, ErrOutOfRange)
}

func TestOttawaAnkleRule1992_noAlloc(t *testing.T) {
	f := OttawaAnkle1992Findings{PainNearMalleoli: true}
	if n := testing.AllocsPerRun(100, func() { r, _ := OttawaAnkleRule1992(60, f); sinkBool = r.MalleolarSet }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkOttawaAnkleRule1992(b *testing.B) {
	b.ReportAllocs()
	f := OttawaAnkle1992Findings{PainNearMalleoli: true}
	for b.Loop() {
		r, _ := OttawaAnkleRule1992(60, f)
		sinkBool = r.MalleolarSet
	}
}
