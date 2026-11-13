package text

import (
	"bytes"
	"errors"
	"math/rand/v2"
	"testing"
)

// TIFF 6.0, Section 9, p. 42, the unpacking loop applied to the hand-derived stream of
// run_length_encode_test.go.
func TestRunLengthDecode_tiffExample(t *testing.T) {
	got, err := RunLengthDecode(nil, packBitsPacked, 1<<10)
	if err != nil || !bytes.Equal(got, packBitsUnpacked) {
		t.Errorf("RunLengthDecode = % X, %v; want % X", got, err, packBitsUnpacked)
	}
}

func TestRunLengthDecode_cases(t *testing.T) {
	tests := []struct {
		name string
		src  []byte
		want []byte
	}{
		{"empty", nil, nil},
		{"one literal", []byte{0x00, 'x'}, []byte("x")},
		{"no-op header is skipped", []byte{0x80, 0x00, 'x', 0x80}, []byte("x")},
		{"replicate 2", []byte{0xFF, 'y'}, []byte("yy")},
		{"replicate 128", []byte{0x81, 'z'}, bytes.Repeat([]byte("z"), 128)},
		{"literal 128", append([]byte{0x7F}, bytes.Repeat([]byte("q"), 128)...), bytes.Repeat([]byte("q"), 128)},
		{"another encoder's choice: 2-run as literal", []byte{0x01, 'b', 'b'}, []byte("bb")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RunLengthDecode(nil, tt.src, 1<<10)
			if err != nil || !bytes.Equal(got, tt.want) {
				t.Errorf("RunLengthDecode(% X) = % X, %v; want % X", tt.src, got, err, tt.want)
			}
		})
	}
}

func TestRunLengthDecode_errors(t *testing.T) {
	tests := []struct {
		name   string
		src    []byte
		maxLen int
		want   error
	}{
		{"negative cap", []byte{0x00, 'a'}, -1, ErrMaxLen},
		{"literal header alone", []byte{0x02}, 10, ErrTruncated},
		{"literal short by one", []byte{0x02, 'a', 'b'}, 10, ErrTruncated},
		{"replicate header alone", []byte{0xFE}, 10, ErrTruncated},
		{"truncated after good data", []byte{0x00, 'a', 0xFE}, 10, ErrTruncated},
		{"replicate over the cap", []byte{0x81, 'a'}, 127, ErrTooLong},
		{"literal over the cap", []byte{0x01, 'a', 'b'}, 1, ErrTooLong},
		{"zero cap", []byte{0x00, 'a'}, 0, ErrTooLong},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := []byte("keep")
			got, err := RunLengthDecode(dst, tt.src, tt.maxLen)
			if !errors.Is(err, tt.want) || string(got) != "keep" {
				t.Errorf("RunLengthDecode(% X, %d) = %q, %v; want \"keep\", %v", tt.src, tt.maxLen, got, err, tt.want)
			}
		})
	}
}

func TestRunLengthDecode_capIsExact(t *testing.T) {
	if got, err := RunLengthDecode(nil, []byte{0x81, 'a'}, 128); err != nil || len(got) != 128 {
		t.Errorf("cap equal to the output: len %d, %v; want 128, nil", len(got), err)
	}
	// The cap counts only the bytes this call appends, not what dst already held.
	if got, err := RunLengthDecode([]byte("xyz"), []byte{0xFF, 'a'}, 2); err != nil || string(got) != "xyzaa" {
		t.Errorf("cap with a prefilled dst: %q, %v; want \"xyzaa\", nil", got, err)
	}
}

func TestRunLengthDecode_noAllocation(t *testing.T) {
	buf := make([]byte, 0, 64)
	if allocs := testing.AllocsPerRun(100, func() { buf, _ = RunLengthDecode(buf[:0], packBitsPacked, 64) }); allocs != 0 {
		t.Errorf("RunLengthDecode with a large enough dst allocates %v times, want 0", allocs)
	}
}

// Random byte streams must never panic: every outcome is data or a documented error.
func TestRunLengthDecode_randomStreamsDoNotPanic(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 500 {
		src := make([]byte, rng.IntN(40))
		for i := range src {
			src[i] = byte(rng.Uint32())
		}
		got, err := RunLengthDecode(nil, src, 300)
		if err != nil && !errors.Is(err, ErrTruncated) && !errors.Is(err, ErrTooLong) {
			t.Fatalf("trial %d: unexpected error %v", trial, err)
		}
		if len(got) > 300 {
			t.Fatalf("trial %d: %d bytes exceed the cap", trial, len(got))
		}
	}
}

func BenchmarkRunLengthDecode(b *testing.B) {
	rng := rand.New(rand.NewPCG(1, 2))
	for _, size := range benchSizes {
		src := []byte(randomText(rng, size.n, "aaab"))
		enc := RunLengthEncode(nil, src)
		buf := make([]byte, 0, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				buf, _ = RunLengthDecode(buf[:0], enc, size.n)
			}
			benchSink = len(buf)
		})
	}
}
