package market

import (
	"slices"
	"testing"
)

var sinkSigns []int8

// Lee & Ready 1991 §I describe the tick test (uptick, zero-uptick, downtick, zero-downtick); hand-checked.
func TestTickRuleSign_cases(t *testing.T) {
	for _, c := range []struct {
		name   string
		prices []int64
		want   []int8
	}{
		{"example", []int64{100, 101, 101, 100, 100, 102}, []int8{0, 1, 1, -1, -1, 1}},
		{"one trade", []int64{5}, []int8{0}},
		{"flat start", []int64{5, 5, 4, 4}, []int8{0, 0, -1, -1}},
	} {
		got, err := TickRuleSign(nil, c.prices)
		if err != nil || !slices.Equal(got, c.want) {
			t.Errorf("%s: %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestTickRuleSign_reusesDst(t *testing.T) {
	dst := make([]int8, 0, 8)
	got, err := TickRuleSign(dst, []int64{1, 2, 3})
	if err != nil || &got[0] != &dst[:1][0] {
		t.Errorf("dst not reused: %v, %v", got, err)
	}
}

func TestTickRuleSign_errors(t *testing.T) {
	_, err := TickRuleSign(nil, nil)
	wantErr(t, "empty", err, ErrEmpty)
	_, err = TickRuleSign(nil, []int64{1, -1})
	wantErr(t, "negative", err, ErrNegative)
	_, err = TickRuleSign(nil, []int64{maxTicks + 1})
	wantErr(t, "too large", err, ErrOutOfRange)
}

func BenchmarkTickRuleSign(b *testing.B) {
	p := make([]int64, 1000)
	for i := range p {
		p[i] = int64(100 + i%4)
	}
	dst := make([]int8, len(p))
	b.ReportAllocs()
	for b.Loop() {
		sinkSigns, sinkErr = TickRuleSign(dst, p)
	}
}
