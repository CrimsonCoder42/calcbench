package market

import (
	"slices"
	"testing"
)

// Lee & Ready 1991 §I–II: quote rule against the midpoint, tick test at the midpoint; hand-checked
// with the quote 100/102 (mid 101).
func TestLeeReadySign_example(t *testing.T) {
	p := []int64{101, 102, 101, 100, 101}
	bid := []int64{100, 100, 100, 100, 100}
	ask := []int64{102, 102, 102, 102, 102}
	got, err := LeeReadySign(nil, p, bid, ask)
	want := []int8{0, 1, -1, -1, 1}
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("LeeReadySign = %v, %v; want %v", got, err, want)
	}
}

func TestLeeReadySign_cases(t *testing.T) {
	// odd spread: mid 100.5 is never hit; outside-quote trades follow the quote rule
	got, err := LeeReadySign(nil, []int64{100, 101, 99, 105}, []int64{100, 100, 100, 100}, []int64{101, 101, 101, 101})
	want := []int8{-1, 1, -1, 1}
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("odd spread: %v, %v; want %v", got, err, want)
	}
	// locked quote: every trade at the quote is at the mid, so the tick test decides
	got, err = LeeReadySign(nil, []int64{50, 50, 51}, []int64{50, 50, 51}, []int64{50, 50, 51})
	want = []int8{0, 0, 1}
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("locked: %v, %v; want %v", got, err, want)
	}
}

func TestLeeReadySign_errors(t *testing.T) {
	_, err := LeeReadySign(nil, nil, nil, nil)
	wantErr(t, "empty", err, ErrEmpty)
	_, err = LeeReadySign(nil, []int64{1}, []int64{1, 2}, []int64{1})
	wantErr(t, "lengths", err, ErrLengthMismatch)
	_, err = LeeReadySign(nil, []int64{1}, []int64{3}, []int64{2})
	wantErr(t, "crossed", err, ErrDomain)
	_, err = LeeReadySign(nil, []int64{-1}, []int64{1}, []int64{2})
	wantErr(t, "negative", err, ErrNegative)
	_, err = LeeReadySign(nil, []int64{maxTicks + 1}, []int64{1}, []int64{2})
	wantErr(t, "too large", err, ErrOutOfRange)
}

func BenchmarkLeeReadySign(b *testing.B) {
	const n = 1000
	p, bid, ask := make([]int64, n), make([]int64, n), make([]int64, n)
	for i := range n {
		p[i], bid[i], ask[i] = int64(100+i%3), 100, 102
	}
	dst := make([]int8, n)
	b.ReportAllocs()
	for b.Loop() {
		sinkSigns, sinkErr = LeeReadySign(dst, p, bid, ask)
	}
}
