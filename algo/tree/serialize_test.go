package tree

import (
	"errors"
	"math/rand/v2"
	"strconv"
	"testing"
)

func appendInt(dst []byte, v int) []byte { return strconv.AppendInt(dst, int64(v), 10) }

func parseInt(b []byte) (int, error) { return strconv.Atoi(string(b)) }

func appendString(dst []byte, v string) []byte { return append(dst, v...) }

func parseString(b []byte) (string, error) { return string(b), nil }

// The format's own worked examples (Serialize doc comment; netstrings per Bernstein 1997, "12:hello world!,").
func TestSerialize_formatExamples(t *testing.T) {
	tests := []struct {
		name string
		root *Node[int]
		want string
	}{
		{"empty", nil, "#"},
		{"leaf", &Node[int]{Value: 7}, "1:7,##"},
		{"1 → (2, 3)", &Node[int]{Value: 1, Left: &Node[int]{Value: 2}, Right: &Node[int]{Value: 3}}, "1:1,1:2,##1:3,##"},
		{"negative, two digits", &Node[int]{Value: -12, Right: &Node[int]{Value: 3456789012}}, "3:-12,#10:3456789012,##"},
		{"shared tree", sharedTree(), "1:1,1:2,1:4,##1:5,1:6,###1:3,#1:7,##"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Serialize(nil, tt.root, appendInt)
			if err != nil || string(got) != tt.want {
				t.Errorf("Serialize = (%q, %v), want %q", got, err, tt.want)
			}
		})
	}
}

func TestSerialize_netstringValue(t *testing.T) {
	got, _ := Serialize([]byte("prefix|"), &Node[string]{Value: "hello world!"}, appendString)
	if string(got) != "prefix|12:hello world!,##" {
		t.Errorf("got %q", got)
	}
}

// Values containing the format's own bytes ('#', ':', ',', digits, empty) survive the round trip.
func TestDeserialize_roundTripAwkwardValues(t *testing.T) {
	n := func(v string, l, r *Node[string]) *Node[string] { return &Node[string]{Value: v, Left: l, Right: r} }
	root := n("#", n("", n("1:2,", nil, nil), nil), n("##,:", nil, n("\x00\xff", nil, nil)))
	data, err := Serialize(nil, root, appendString)
	if err != nil {
		t.Fatal(err)
	}
	got, err := Deserialize(data, parseString)
	if err != nil {
		t.Fatalf("Deserialize(%q): %v", data, err)
	}
	again, _ := Serialize(nil, got, appendString)
	if string(again) != string(data) {
		t.Errorf("round trip changed %q to %q", data, again)
	}
}

func TestDeserialize_cases(t *testing.T) {
	for _, s := range []string{"#", "1:7,##", "1:1,1:2,##1:3,##", "1:1,1:2,1:4,##1:5,1:6,###1:3,#1:7,##"} {
		root, err := Deserialize([]byte(s), parseInt)
		if err != nil {
			t.Errorf("Deserialize(%q): %v", s, err)
			continue
		}
		if back, _ := Serialize(nil, root, appendInt); string(back) != s {
			t.Errorf("Deserialize(%q) re-serialized as %q", s, back)
		}
	}
	root, _ := Deserialize([]byte("1:1,1:2,1:4,##1:5,1:6,###1:3,#1:7,##"), parseInt)
	if !sameShape(root, sharedTree()) {
		t.Error("shared tree did not decode to its shape")
	}
}

func TestDeserialize_errors(t *testing.T) {
	malformed := []string{
		"",                             // nothing
		"##",                           // trailing bytes
		"1:7,#",                        // truncated
		"1:7,###",                      // trailing
		"x",                            // not a digit
		"01:7,##",                      // leading zero
		"1;7,##",                       // missing ':'
		"1:7;##",                       // missing ','
		"5:7,##",                       // length past the end
		"1:",                           // value missing
		"1",                            // colon missing
		"99999999999999999999999:x,##", // absurd length, no overflow
	}
	for _, s := range malformed {
		if root, err := Deserialize([]byte(s), parseInt); !errors.Is(err, ErrMalformed) || root != nil {
			t.Errorf("Deserialize(%q) = (%v, %v), want ErrMalformed", s, root, err)
		}
	}
	if _, err := Deserialize([]byte("1:x,##"), parseInt); err == nil || !errors.Is(err, strconv.ErrSyntax) {
		t.Errorf("parse error not wrapped: %v", err)
	}
	if _, err := Deserialize([]byte("#"), (func([]byte) (int, error))(nil)); !errors.Is(err, ErrNilFunc) {
		t.Errorf("nil parseValue: %v", err)
	}
	dst := []byte("keep")
	if got, err := Serialize(dst, sharedTree(), nil); !errors.Is(err, ErrNilFunc) || string(got) != "keep" {
		t.Errorf("nil appendValue: (%q, %v)", got, err)
	}
}

func TestDeserialize_zeroLengthValue(t *testing.T) {
	root, err := Deserialize([]byte("0:,##"), parseString)
	if err != nil || root == nil || root.Value != "" {
		t.Errorf("Deserialize(\"0:,##\") = (%v, %v), want one empty-string node", root, err)
	}
}

func TestSerialize_deepChain(t *testing.T) {
	for _, left := range []bool{true, false} {
		root := chain(deepN, left)
		data, err := Serialize(nil, root, appendInt)
		if err != nil {
			t.Fatal(err)
		}
		got, err := Deserialize(data, parseInt)
		if err != nil || Height(got) != deepN-1 {
			t.Errorf("chain left=%v: err %v, height %d", left, err, Height(got))
		}
	}
}

// Oracle: random shapes round-trip to the same shape and values.
func TestDeserialize_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		root := randomTree(rng, rng.IntN(60))
		data, err := Serialize(nil, root, appendInt)
		if err != nil {
			t.Fatal(err)
		}
		got, err := Deserialize(data, parseInt)
		if err != nil || !sameShape(got, root) {
			t.Fatalf("trial %d: round trip of %q failed: %v", trial, data, err)
		}
		// Every proper prefix is malformed (the format is a prefix code).
		cut := rng.IntN(len(data))
		if _, err := Deserialize(data[:cut], parseInt); !errors.Is(err, ErrMalformed) {
			t.Fatalf("trial %d: prefix %q accepted: %v", trial, data[:cut], err)
		}
	}
}

func BenchmarkSerialize(b *testing.B) {
	for _, size := range benchSizes {
		root := balancedTree(size.n)
		buf := make([]byte, 0, 16*size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				buf, _ = Serialize(buf[:0], root, appendInt)
			}
		})
	}
}

func BenchmarkDeserialize(b *testing.B) {
	for _, size := range benchSizes {
		data, _ := Serialize(nil, balancedTree(size.n), appendInt)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = Deserialize(data, parseInt)
			}
		})
	}
}
