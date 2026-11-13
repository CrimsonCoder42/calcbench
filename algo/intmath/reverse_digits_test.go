package intmath

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"strconv"
	"testing"
)

func TestReverseDigits_cases(t *testing.T) {
	tests := []struct{ x, want int64 }{
		{0, 0}, {7, 7}, {1234, 4321}, {-560, -65}, {1200, 21}, {-1, -1},
		{1_000_000_000_000_000_000, 1}, {math.MinInt64, -8_085_774_586_302_733_229},
		{7_463_847_412, 2_147_483_647}, {math.MaxInt64, 7_085_774_586_302_733_229}, {1_085_774_586_302_733_229, 9_223_372_036_854_775_801},
	}
	for _, tt := range tests {
		if got, err := ReverseDigits(tt.x); err != nil || got != tt.want {
			t.Errorf("ReverseDigits(%d) = (%d, %v), want %d", tt.x, got, err, tt.want)
		}
	}
	if got, err := ReverseDigits[int32](-2_147_483_412); err != nil || got != -2_143_847_412 {
		t.Errorf("int32 ReverseDigits(-2147483412) = (%d, %v)", got, err)
	}
}

func TestReverseDigits_errors(t *testing.T) {
	if _, err := ReverseDigits[int32](1_000_000_009); !errors.Is(err, ErrOverflow) {
		t.Errorf("int32 1000000009: error = %v, want ErrOverflow", err)
	}
	if _, err := ReverseDigits[int32](math.MinInt32); !errors.Is(err, ErrOverflow) {
		t.Errorf("int32 MinInt32: error = %v, want ErrOverflow", err)
	}
	if _, err := ReverseDigits[int64](1_999_999_999_999_999_999); !errors.Is(err, ErrOverflow) {
		t.Errorf("int64 1999…9: error = %v, want ErrOverflow", err)
	}
	if _, err := ReverseDigits[int64](1_000_000_000_000_000_039); !errors.Is(err, ErrOverflow) {
		t.Errorf("int64 1000…039: error = %v, want ErrOverflow", err)
	}
	if _, err := ReverseDigits[int64](-1_000_000_000_000_000_039); !errors.Is(err, ErrOverflow) {
		t.Errorf("int64 −1000…039: error = %v, want ErrOverflow", err)
	}
	if _, err := ReverseDigits[int8](-128); !errors.Is(err, ErrOverflow) {
		t.Errorf("int8 −128: error = %v, want ErrOverflow", err)
	}
}

// Oracle: reverse the decimal string, parse with math/big, and check the int64 range.
func TestReverseDigits_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	lo, hi := big.NewInt(math.MinInt64), big.NewInt(math.MaxInt64)
	for range 600 {
		x := int64(rng.Uint64()) >> rng.IntN(64)
		s := strconv.FormatUint(magnitude(x), 10)
		rs := []byte(s)
		for i, j := 0, len(rs)-1; i < j; i, j = i+1, j-1 {
			rs[i], rs[j] = rs[j], rs[i]
		}
		want, _ := new(big.Int).SetString(string(rs), 10)
		if x < 0 {
			want.Neg(want)
		}
		got, err := ReverseDigits(x)
		if want.Cmp(lo) < 0 || want.Cmp(hi) > 0 {
			if !errors.Is(err, ErrOverflow) {
				t.Fatalf("ReverseDigits(%d) = (%d, %v), want ErrOverflow", x, got, err)
			}
			continue
		}
		if err != nil || got != want.Int64() {
			t.Fatalf("ReverseDigits(%d) = (%d, %v), want %s", x, got, err, want)
		}
	}
}

func BenchmarkReverseDigits(b *testing.B) {
	b.ReportAllocs()
	var s int64
	for i := int64(0); b.Loop(); i++ {
		v, _ := ReverseDigits(1_234_567_890_123_456 + i)
		s += v
	}
	benchSink = uint64(s)
}
