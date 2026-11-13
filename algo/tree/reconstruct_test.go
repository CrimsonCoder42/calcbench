package tree

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// sameShape reports whether two trees have the same shape and values.
func sameShape(a, b *Node[int]) bool {
	return slices.Equal(refPre(a, nil), refPre(b, nil)) && slices.Equal(refIn(a, nil), refIn(b, nil)) &&
		Height(a) == Height(b)
}

// Laaksonen (2018) §14.4 p. 139: the shared tree is the only one with preorder 1 2 4 5 6 3 7 and inorder
// 4 2 6 5 1 3 7; postorder and inorder determine it too.
func TestFromPreIn_laaksonenExample(t *testing.T) {
	got, err := FromPreIn(sharedPre, sharedIn)
	if err != nil || !sameShape(got, sharedTree()) || !slices.Equal(refPost(got, nil), sharedPost) {
		t.Errorf("FromPreIn(shared) = (%v, %v), want the shared tree", refPre(got, nil), err)
	}
}

func TestFromInPost_laaksonenExample(t *testing.T) {
	got, err := FromInPost(sharedIn, sharedPost)
	if err != nil || !sameShape(got, sharedTree()) {
		t.Errorf("FromInPost(shared) = (%v, %v), want the shared tree", refPre(got, nil), err)
	}
}

func TestFromPreIn_cases(t *testing.T) {
	tests := []struct {
		name    string
		pre, in []int
	}{
		{"empty", nil, nil},
		{"one", []int{7}, []int{7}},
		{"two, left child", []int{1, 2}, []int{2, 1}},
		{"two, right child", []int{1, 2}, []int{1, 2}},
		{"negative and extremes", []int{0, math.MinInt, math.MaxInt}, []int{math.MinInt, 0, math.MaxInt}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromPreIn(tt.pre, tt.in)
			if err != nil || !slices.Equal(refPre(got, nil), tt.pre) || !slices.Equal(refIn(got, nil), tt.in) {
				t.Errorf("FromPreIn = (pre %v in %v, %v)", refPre(got, nil), refIn(got, nil), err)
			}
			post := refPost(got, nil)
			got2, err := FromInPost(tt.in, post)
			if err != nil || !sameShape(got, got2) {
				t.Errorf("FromInPost = (pre %v, %v)", refPre(got2, nil), err)
			}
		})
	}
}

func TestFromPreIn_errors(t *testing.T) {
	tests := []struct {
		name    string
		pre, in []int
		want    error
	}{
		{"length", []int{1, 2}, []int{1}, ErrLength},
		{"duplicate in inorder", []int{1, 1}, []int{1, 1}, ErrDuplicate},
		{"duplicate in preorder", []int{1, 1}, []int{1, 2}, ErrDuplicate},
		{"different values", []int{1, 3}, []int{1, 2}, ErrInconsistent},
		{"same set, impossible orders", []int{1, 2, 3}, []int{3, 1, 2}, ErrInconsistent},
		{"subtree orders swapped", []int{2, 1, 3}, []int{3, 2, 1}, ErrInconsistent},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := FromPreIn(tt.pre, tt.in); !errors.Is(err, tt.want) || got != nil {
				t.Errorf("FromPreIn(%v, %v) = (%v, %v), want %v", tt.pre, tt.in, got, err, tt.want)
			}
		})
	}
}

func TestFromInPost_errors(t *testing.T) {
	tests := []struct {
		name     string
		in, post []int
		want     error
	}{
		{"length", []int{1}, nil, ErrLength},
		{"duplicate in inorder", []int{2, 2}, []int{2, 2}, ErrDuplicate},
		{"duplicate in postorder", []int{1, 2}, []int{2, 2}, ErrDuplicate},
		{"different values", []int{1, 2}, []int{1, 3}, ErrInconsistent},
		{"same set, impossible orders", []int{1, 2, 3}, []int{3, 1, 2}, ErrInconsistent},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := FromInPost(tt.in, tt.post); !errors.Is(err, tt.want) || got != nil {
				t.Errorf("FromInPost(%v, %v) = (%v, %v), want %v", tt.in, tt.post, got, err, tt.want)
			}
		})
	}
}

func TestFromPreIn_NaN(t *testing.T) {
	nan := math.NaN()
	if _, err := FromPreIn([]float64{nan}, []float64{nan}); !errors.Is(err, ErrInconsistent) {
		t.Errorf("FromPreIn(NaN) error = %v, want ErrInconsistent", err)
	}
}

func TestFromPreIn_deepChains(t *testing.T) {
	for _, left := range []bool{true, false} {
		root := chain(deepN, left)
		pre, in, post := slices.Collect(Preorder(root)), slices.Collect(Inorder(root)), slices.Collect(Postorder(root))
		got, err := FromPreIn(pre, in)
		if err != nil || Height(got) != deepN-1 {
			t.Errorf("FromPreIn(chain left=%v): err %v, height %d", left, err, Height(got))
		}
		got, err = FromInPost(in, post)
		if err != nil || Height(got) != deepN-1 {
			t.Errorf("FromInPost(chain left=%v): err %v, height %d", left, err, Height(got))
		}
	}
}

func TestFromPreIn_doesNotModifyInput(t *testing.T) {
	pre, in := slices.Clone(sharedPre), slices.Clone(sharedIn)
	_, _ = FromPreIn(pre, in)
	_, _ = FromInPost(in, slices.Clone(sharedPost))
	if !slices.Equal(pre, sharedPre) || !slices.Equal(in, sharedIn) {
		t.Error("inputs modified")
	}
}

// Oracle: random trees round-trip through both reconstructions; random permutations are rejected unless
// they happen to be consistent, which the round-trip check confirms.
func TestFromPreIn_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		root := randomTree(rng, rng.IntN(50))
		pre, in, post := refPre(root, nil), refIn(root, nil), refPost(root, nil)
		got, err := FromPreIn(pre, in)
		if err != nil || !sameShape(got, root) {
			t.Fatalf("trial %d: FromPreIn failed: %v", trial, err)
		}
		got, err = FromInPost(in, post)
		if err != nil || !sameShape(got, root) {
			t.Fatalf("trial %d: FromInPost failed: %v", trial, err)
		}
		shuffled := slices.Clone(pre)
		rng.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })
		if got, err := FromPreIn(shuffled, in); err == nil {
			if !slices.Equal(refPre(got, nil), shuffled) || !slices.Equal(refIn(got, nil), in) {
				t.Fatalf("trial %d: accepted inconsistent input", trial)
			}
		} else if !errors.Is(err, ErrInconsistent) {
			t.Fatalf("trial %d: error %v, want ErrInconsistent", trial, err)
		}
	}
}

func BenchmarkFromPreIn(b *testing.B) {
	for _, size := range benchSizes {
		root := balancedTree(size.n)
		pre, in := slices.Collect(Preorder(root)), slices.Collect(Inorder(root))
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = FromPreIn(pre, in)
			}
		})
	}
}

func BenchmarkFromInPost(b *testing.B) {
	for _, size := range benchSizes {
		root := balancedTree(size.n)
		in, post := slices.Collect(Inorder(root)), slices.Collect(Postorder(root))
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = FromInPost(in, post)
			}
		})
	}
}
