package trauma

import "testing"

// Lee 2005 Table 1 (PMC2779435): every category's points. HAND-DERIVED totals.
func TestTLICS_table1(t *testing.T) {
	morph := map[TLICSMorphology]int{TLICSCompression: 1, TLICSBurst: 2, TLICSTranslationalRotational: 3, TLICSDistraction: 4}
	neuro := map[TLICSNeurologic]int{TLICSNeuroIntact: 0, TLICSNerveRoot: 2, TLICSCordConusIncomplete: 3,
		TLICSCordConusComplete: 2, TLICSCaudaEquina: 3}
	plc := map[TLICSPLC]int{TLICSPLCIntact: 0, TLICSPLCIndeterminate: 2, TLICSPLCInjured: 3}
	for m, mp := range morph {
		for n, np := range neuro {
			for p, pp := range plc {
				if got, err := TLICS(m, n, p); err != nil || got != mp+np+pp {
					t.Errorf("TLICS(%d, %d, %d) = (%d, %v), want %d", m, n, p, got, err, mp+np+pp)
				}
			}
		}
	}
	if got, _ := TLICS(TLICSCompression, TLICSNeuroIntact, TLICSPLCIntact); got != 1 {
		t.Errorf("minimum = %d, want 1", got)
	}
	if got, _ := TLICS(TLICSDistraction, TLICSCordConusIncomplete, TLICSPLCInjured); got != 10 {
		t.Errorf("maximum = %d, want 10", got)
	}
}

// Doc comment example: burst + intact + indeterminate = 4 (Table 1).
func TestTLICS_example(t *testing.T) {
	if got, err := TLICS(TLICSBurst, TLICSNeuroIntact, TLICSPLCIndeterminate); err != nil || got != 4 {
		t.Errorf("got (%d, %v), want 4", got, err)
	}
}

// parity: ledger tlics, source cases. The app's arrays (morphology 0–4 with "None" = 0, neuro [0,2,2,3,3] for
// intact, nerve root, complete cord, incomplete cord, cauda equina; PLC [0,2,3]) agree with Table 1 except the
// "None" morphology, which Table 1 does not have (refused here). Compression + complete cord + injured → 1+2+3 = 6.
func TestTLICS_parity(t *testing.T) {
	if got, err := TLICS(TLICSCompression, TLICSCordConusComplete, TLICSPLCInjured); err != nil || got != 6 {
		t.Errorf("got (%d, %v), want 6", got, err)
	}
	_, err := TLICS(0, TLICSNeuroIntact, TLICSPLCIntact)
	checkErr(t, "None morphology", err, ErrOutOfRange)
}

func TestTLICS_errors(t *testing.T) {
	_, err := TLICS(5, TLICSNeuroIntact, TLICSPLCIntact)
	checkErr(t, "morphology 5", err, ErrOutOfRange)
	_, err = TLICS(TLICSBurst, -1, TLICSPLCIntact)
	checkErr(t, "neuro -1", err, ErrOutOfRange)
	_, err = TLICS(TLICSBurst, 5, TLICSPLCIntact)
	checkErr(t, "neuro 5", err, ErrOutOfRange)
	_, err = TLICS(TLICSBurst, TLICSNeuroIntact, 3)
	checkErr(t, "plc 3", err, ErrOutOfRange)
	_, err = TLICS(TLICSBurst, TLICSNeuroIntact, -1)
	checkErr(t, "plc -1", err, ErrOutOfRange)
}

func TestTLICS_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = TLICS(TLICSBurst, TLICSNerveRoot, TLICSPLCInjured) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkTLICS(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = TLICS(TLICSBurst, TLICSNerveRoot, TLICSPLCInjured)
	}
}
