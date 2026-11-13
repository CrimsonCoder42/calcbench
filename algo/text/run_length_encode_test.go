package text

import (
	"bytes"
	"math/rand/v2"
	"testing"
)

// packBitsUnpacked and packBitsPacked: hand-derived from TIFF 6.0, Section 9, p. 42: AA×3 → replicate
// FE AA; 80 00 2A → literal 02 …; AA×4 → FD AA; 80 00 2A 22 → literal 03 …; AA×10 → F7 AA.
var (
	packBitsUnpacked = []byte{
		0xAA, 0xAA, 0xAA, 0x80, 0x00, 0x2A, 0xAA, 0xAA, 0xAA, 0xAA, 0x80, 0x00, 0x2A, 0x22,
		0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA,
	}
	packBitsPacked = []byte{0xFE, 0xAA, 0x02, 0x80, 0x00, 0x2A, 0xFD, 0xAA, 0x03, 0x80, 0x00, 0x2A, 0x22, 0xF7, 0xAA}
)

func TestRunLengthEncode_tiffExample(t *testing.T) {
	if got := RunLengthEncode(nil, packBitsUnpacked); !bytes.Equal(got, packBitsPacked) {
		t.Errorf("RunLengthEncode = % X, want % X", got, packBitsPacked)
	}
}

func TestRunLengthEncode_cases(t *testing.T) {
	long := bytes.Repeat([]byte{7}, 130)
	tests := []struct {
		name string
		src  []byte
		want []byte
	}{
		{"empty", nil, nil},
		{"one byte", []byte{5}, []byte{0x00, 5}},
		{"two equal, alone: replicate", []byte{5, 5}, []byte{0xFF, 5}},
		{"two different", []byte{5, 6}, []byte{0x01, 5, 6}},
		{"2-run between literals joins them", []byte{'a', 'b', 'b', 'c'}, []byte{0x03, 'a', 'b', 'b', 'c'}},
		{"2-run before a 2-run: replicate", []byte{'a', 'b', 'b', 'c', 'c'}, []byte{0x00, 'a', 0xFF, 'b', 0xFF, 'c'}},
		{"2-run at the end: replicate", []byte{'a', 'b', 'b'}, []byte{0x00, 'a', 0xFF, 'b'}},
		{"2-run with no literal before: replicate", []byte{'b', 'b', 'c'}, []byte{0xFF, 'b', 0x00, 'c'}},
		{"128-run is one header", bytes.Repeat([]byte{9}, 128), []byte{0x81, 9}},
		{"130-run splits 128 + 2", long, []byte{0x81, 7, 0xFF, 7}},
		{"129-run leaves one literal", long[:129], []byte{0x81, 7, 0x00, 7}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RunLengthEncode(nil, tt.src); !bytes.Equal(got, tt.want) {
				t.Errorf("RunLengthEncode(% X) = % X, want % X", tt.src, got, tt.want)
			}
		})
	}
}

func TestRunLengthEncode_longLiteralSplits(t *testing.T) {
	src := make([]byte, 300)
	for i := range src {
		src[i] = byte(i % 251) // no two neighbours equal
	}
	got := RunLengthEncode(nil, src)
	if len(got) != 300+3 || got[0] != 127 || got[129] != 127 || got[258] != 300-256-1 {
		t.Errorf("300 literal bytes: length %d, headers %d %d %d; want 303, 127 127 43", len(got), got[0], got[129], got[258])
	}
}

func TestRunLengthEncode_noAllocation(t *testing.T) {
	buf := make([]byte, 0, 64)
	if allocs := testing.AllocsPerRun(100, func() { buf = RunLengthEncode(buf[:0], packBitsUnpacked) }); allocs != 0 {
		t.Errorf("RunLengthEncode with a large enough dst allocates %v times, want 0", allocs)
	}
}

// Round trip and the size bound on generated inputs rich in short runs.
func TestRunLengthEncode_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 500 {
		src := []byte(randomText(rng, rng.IntN(400), "aab\x00"[:1+trial%4]))
		enc := RunLengthEncode(nil, src)
		if bound := len(src) + len(src)/3 + 1; len(enc) > bound {
			t.Fatalf("trial %d: encoded %d bytes into %d, above the bound %d", trial, len(src), len(enc), bound)
		}
		dec, err := RunLengthDecode(nil, enc, len(src))
		if err != nil || !bytes.Equal(dec, src) {
			t.Fatalf("trial %d: round trip of % X gave % X, %v", trial, src, dec, err)
		}
	}
}

func BenchmarkRunLengthEncode(b *testing.B) {
	rng := rand.New(rand.NewPCG(1, 2))
	for _, size := range benchSizes {
		src := []byte(randomText(rng, size.n, "aaab"))
		buf := make([]byte, 0, 2*size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				buf = RunLengthEncode(buf[:0], src)
			}
			benchSink = len(buf)
		})
	}
}
