package list

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"slices"
	"testing"
)

// digitsToBig reads least-significant-first digits in radix as a big.Int (the oracle).
func digitsToBig(digits []int, radix int) *big.Int {
	r := big.NewInt(int64(radix))
	sum := new(big.Int)
	for i := len(digits) - 1; i >= 0; i-- {
		sum.Mul(sum, r)
		sum.Add(sum, big.NewInt(int64(digits[i])))
	}
	return sum
}

// Worked example (Knuth Vol. 2 §4.3.1 Algorithm A, numbers hand-chosen): 342 + 465 = 807 in radix 10,
// written least significant first. The same sum and 99 + 1 are the problem-site pattern, in
// our words.
func TestAddTwoNumbers_workedExample(t *testing.T) {
	for _, tt := range []struct{ a, b, want []int }{
		{[]int{2, 4, 3}, []int{5, 6, 4}, []int{7, 0, 8}},
		{[]int{9, 9}, []int{1}, []int{0, 0, 1}},
	} {
		head, err := AddTwoNumbers(FromSlice(tt.a), FromSlice(tt.b), 10)
		if got := values(t, head); err != nil || !slices.Equal(got, tt.want) {
			t.Errorf("%v + %v = (%v, %v), want %v", tt.a, tt.b, got, err, tt.want)
		}
	}
}

func TestAddTwoNumbers_cases(t *testing.T) {
	tests := []struct {
		name    string
		a, b    []int
		radix   int
		want    []int
		wantNil bool
	}{
		{"nil + nil", nil, nil, 10, nil, true},
		{"nil + digits", nil, []int{3, 2}, 10, []int{3, 2}, false},
		{"zero + zero", []int{0}, []int{0}, 10, []int{0}, false},
		{"leading zeros kept", []int{1, 0, 0}, []int{1}, 10, []int{2, 0, 0}, false},
		{"binary carry chain", []int{1, 1, 1}, []int{1}, 2, []int{0, 0, 0, 1}, false},
		{"radix 16", []int{15}, []int{15}, 16, []int{14, 1}, false},
		{"largest radix", []int{math.MaxInt - 1}, []int{math.MaxInt - 1}, math.MaxInt, []int{math.MaxInt - 2, 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head, err := AddTwoNumbers(FromSlice(tt.a), FromSlice(tt.b), tt.radix)
			if err != nil || (head == nil) != tt.wantNil || !slices.Equal(values(t, head), tt.want) {
				t.Errorf("got (%v, %v), want %v", values(t, head), err, tt.want)
			}
		})
	}
}

func TestAddTwoNumbers_inputsUnchanged(t *testing.T) {
	a, b := FromSlice([]int{9, 9}), FromSlice([]int{9})
	if _, err := AddTwoNumbers(a, b, 10); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(values(t, a), []int{9, 9}) || !slices.Equal(values(t, b), []int{9}) {
		t.Error("input modified")
	}
}

func TestAddTwoNumbers_errors(t *testing.T) {
	one := FromSlice([]int{1})
	for _, r := range []int{1, 0, -10} {
		if h, err := AddTwoNumbers(one, one, r); !errors.Is(err, ErrRadix) || h != nil {
			t.Errorf("radix %d: err = %v", r, err)
		}
	}
	for _, bad := range []int{-1, 10} {
		if _, err := AddTwoNumbers(FromSlice([]int{1, bad}), one, 10); !errors.Is(err, ErrDigit) {
			t.Errorf("a digit %d: err = %v", bad, err)
		}
		if _, err := AddTwoNumbers(one, FromSlice([]int{bad}), 10); !errors.Is(err, ErrDigit) {
			t.Errorf("b digit %d: err = %v", bad, err)
		}
	}
	cyc, _ := cyclic([]int{1, 2}, 0)
	if _, err := AddTwoNumbers(cyc, one, 10); !errors.Is(err, ErrCycle) {
		t.Errorf("cyclic a: err = %v", err)
	}
	if _, err := AddTwoNumbers(one, cyc, 10); !errors.Is(err, ErrCycle) {
		t.Errorf("cyclic b: err = %v", err)
	}
}

func TestAddTwoNumbers_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		radix := 2 + rng.IntN(20)
		a := randInts(rng, rng.IntN(30), radix)
		b := randInts(rng, rng.IntN(30), radix)
		head, err := AddTwoNumbers(FromSlice(a), FromSlice(b), radix)
		if err != nil {
			t.Fatal(err)
		}
		got := values(t, head)
		want := new(big.Int).Add(digitsToBig(a, radix), digitsToBig(b, radix))
		if digitsToBig(got, radix).Cmp(want) != 0 || len(got) < max(len(a), len(b)) || len(got) > max(len(a), len(b))+1 {
			t.Fatalf("%v + %v (radix %d) = %v", a, b, radix, got)
		}
	}
}

func TestAddTwoNumbers_oneAllocation(t *testing.T) {
	a, b := FromSlice(countUpTo(10)), FromSlice(countUpTo(10))
	var sink *Node[int]
	if allocs := testing.AllocsPerRun(100, func() { sink, _ = AddTwoNumbers(a, b, 10) }); allocs != 1 || sink == nil {
		t.Errorf("allocs = %v, want 1", allocs)
	}
}

func BenchmarkAddTwoNumbers(b *testing.B) {
	for _, s := range benchSizes {
		rng := rand.New(rand.NewPCG(1, 2))
		x, y := FromSlice(randInts(rng, s.n, 10)), FromSlice(randInts(rng, s.n, 10))
		b.Run(s.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = AddTwoNumbers(x, y, 10)
			}
		})
	}
}
