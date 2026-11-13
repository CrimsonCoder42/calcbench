package neuro

import "testing"

// Edmiaston 2010: fail when any item is present; each item alone and none.
func TestAcuteStrokeDysphagiaScreenFailed_example(t *testing.T) {
	for _, it := range []AcuteStrokeDysphagiaScreenItems{
		{ReducedConsciousness: true}, {FacialAsymmetry: true}, {LingualAsymmetry: true},
		{PalatalAsymmetry: true}, {WaterSwallowFailed: true},
	} {
		if !AcuteStrokeDysphagiaScreenFailed(it) {
			t.Errorf("%+v: want failed", it)
		}
	}
	if AcuteStrokeDysphagiaScreenFailed(AcuteStrokeDysphagiaScreenItems{}) {
		t.Error("none: want passed")
	}
}

// parity: ledger bjh_stroke_dysphagia, source cases 1–2 (consumer: any of GCS, facial, tongue, palatal fails stage 1).
func TestAcuteStrokeDysphagiaScreenFailed_parity(t *testing.T) {
	if !AcuteStrokeDysphagiaScreenFailed(AcuteStrokeDysphagiaScreenItems{LingualAsymmetry: true}) { // parity: ledger bjh_stroke_dysphagia, source case 1
		t.Error("case 1")
	}
	if AcuteStrokeDysphagiaScreenFailed(AcuteStrokeDysphagiaScreenItems{}) { // parity: ledger bjh_stroke_dysphagia, source case 2
		t.Error("case 2")
	}
}

func TestAcuteStrokeDysphagiaScreenFailed_noAlloc(t *testing.T) {
	it := AcuteStrokeDysphagiaScreenItems{WaterSwallowFailed: true}
	if n := testing.AllocsPerRun(100, func() { sinkBool = AcuteStrokeDysphagiaScreenFailed(it) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkAcuteStrokeDysphagiaScreenFailed(b *testing.B) {
	b.ReportAllocs()
	it := AcuteStrokeDysphagiaScreenItems{WaterSwallowFailed: true}
	for b.Loop() {
		sinkBool = AcuteStrokeDysphagiaScreenFailed(it)
	}
}
