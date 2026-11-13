package intmath

import (
	"errors"
	"math/big"
	"math/rand/v2"
	"slices"
	"testing"
)

func TestIncrementDigits_cases(t *testing.T) {
	tests := []struct {
		name         string
		digits, want []uint8
	}{
		{"no carry", []uint8{1, 2, 3}, []uint8{1, 2, 4}},
		{"carry", []uint8{1, 2, 9}, []uint8{1, 3, 0}},
		{"all nines", []uint8{9, 9}, []uint8{1, 0, 0}},
		{"one nine", []uint8{9}, []uint8{1, 0}},
		{"zero", []uint8{0}, []uint8{1}},
		{"leading zero kept", []uint8{0, 9}, []uint8{1, 0}},
		{"leading zeros kept", []uint8{0, 0, 4}, []uint8{0, 0, 5}},
		{"inner nines", []uint8{8, 9, 9, 9}, []uint8{9, 0, 0, 0}},
	}
	for _, tt := range tests {
		before := slices.Clone(tt.digits)
		got, err := IncrementDigits(nil, tt.digits)
		if err != nil || !slices.Equal(got, tt.want) {
			t.Errorf("%s: IncrementDigits(%v) = (%v, %v), want %v", tt.name, tt.digits, got, err, tt.want)
		}
		if !slices.Equal(tt.digits, before) {
			t.Errorf("%s: input modified", tt.name)
		}
	}
	if got, _ := IncrementDigits([]uint8{7}, []uint8{4, 9}); !slices.Equal(got, []uint8{7, 5, 0}) {
		t.Errorf("append to dst: %v, want [7 5 0]", got)
	}
}

func TestIncrementDigits_errors(t *testing.T) {
	dst := []uint8{7}
	for _, tt := range []struct {
		digits  []uint8
		wantErr error
	}{{nil, ErrEmpty}, {[]uint8{}, ErrEmpty}, {[]uint8{1, 10}, ErrInvalidDigit}, {[]uint8{255}, ErrInvalidDigit}} {
		got, err := IncrementDigits(dst, tt.digits)
		if !errors.Is(err, tt.wantErr) || !slices.Equal(got, []uint8{7}) {
			t.Errorf("IncrementDigits(%v) = (%v, %v), want ([7], %v)", tt.digits, got, err, tt.wantErr)
		}
	}
}

// Oracle: math/big plus one, compared on value and on length (input length, or one more on carry out).
func TestIncrementDigits_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 400 {
		digits := make([]uint8, 1+rng.IntN(40))
		for i := range digits {
			digits[i] = uint8(rng.IntN(10))
			if rng.IntN(3) == 0 {
				digits[i] = 9
			}
		}
		s := make([]byte, len(digits))
		for i, d := range digits {
			s[i] = '0' + d
		}
		x, _ := new(big.Int).SetString(string(s), 10)
		want := x.Add(x, big.NewInt(1)).String()
		got, err := IncrementDigits(nil, digits)
		if err != nil {
			t.Fatal(err)
		}
		g := make([]byte, len(got))
		for i, d := range got {
			g[i] = '0' + d
		}
		if gv, _ := new(big.Int).SetString(string(g), 10); gv.String() != want || (len(got) != len(digits) && len(got) != len(digits)+1) {
			t.Fatalf("IncrementDigits(%v) = %v, want value %s", digits, got, want)
		}
	}
}

func TestIncrementDigits_noAllocationWithBuffer(t *testing.T) {
	digits := slices.Repeat([]uint8{9}, 100)
	buf := make([]uint8, 0, 101)
	if a := testing.AllocsPerRun(20, func() { buf, _ = IncrementDigits(buf[:0], digits) }); a != 0 {
		t.Errorf("IncrementDigits with a reused buffer allocates %v times, want 0", a)
	}
}

func BenchmarkIncrementDigits(b *testing.B) {
	for _, size := range benchSizes {
		digits := slices.Repeat([]uint8{9}, size.n) // worst case: the carry runs the whole length
		buf := make([]uint8, 0, size.n+1)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				buf, _ = IncrementDigits(buf[:0], digits)
			}
			benchSink = uint64(len(buf))
		})
	}
}
