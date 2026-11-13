package computing

import (
	"errors"
	"testing"
)

// Hand computation (RFC 4632 §3.1): 192.168.37.200 = 0xC0A825C8; /22 mask 0xFFFFFC00 = 255.255.252.0;
// network 0xC0A82400 = 192.168.36.0; broadcast 0xC0A827FF = 192.168.39.255; 2^10 = 1024 addresses, 1022 hosts.
func TestIPv4SubnetOf_slash22(t *testing.T) {
	got, err := IPv4SubnetOf([4]byte{192, 168, 37, 200}, 22)
	want := IPv4Subnet{
		Network: [4]byte{192, 168, 36, 0}, Broadcast: [4]byte{192, 168, 39, 255},
		Mask: [4]byte{255, 255, 252, 0}, Wildcard: [4]byte{0, 0, 3, 255},
		FirstHost: [4]byte{192, 168, 36, 1}, LastHost: [4]byte{192, 168, 39, 254},
		TotalAddresses: 1024, UsableHosts: 1022,
	}
	if err != nil || got != want {
		t.Errorf("got %+v %v\nwant %+v", got, err, want)
	}
}

func TestIPv4SubnetOf_edges(t *testing.T) {
	cases := []struct {
		name        string
		addr        [4]byte
		p           int
		net, bc     [4]byte
		first, last [4]byte
		total, use  uint64
	}{
		{"/0", [4]byte{10, 1, 2, 3}, 0, [4]byte{}, [4]byte{255, 255, 255, 255}, [4]byte{0, 0, 0, 1}, [4]byte{255, 255, 255, 254}, 1 << 32, 1<<32 - 2},
		{"/30", [4]byte{74, 71, 54, 244}, 30, [4]byte{74, 71, 54, 244}, [4]byte{74, 71, 54, 247}, [4]byte{74, 71, 54, 245}, [4]byte{74, 71, 54, 246}, 4, 2},
		{"/31 RFC 3021", [4]byte{10, 0, 0, 1}, 31, [4]byte{10, 0, 0, 0}, [4]byte{10, 0, 0, 1}, [4]byte{10, 0, 0, 0}, [4]byte{10, 0, 0, 1}, 2, 2},
		{"/32", [4]byte{10, 0, 0, 7}, 32, [4]byte{10, 0, 0, 7}, [4]byte{10, 0, 0, 7}, [4]byte{10, 0, 0, 7}, [4]byte{10, 0, 0, 7}, 1, 1},
	}
	for _, c := range cases {
		got, err := IPv4SubnetOf(c.addr, c.p)
		if err != nil || got.Network != c.net || got.Broadcast != c.bc || got.FirstHost != c.first ||
			got.LastHost != c.last || got.TotalAddresses != c.total || got.UsableHosts != c.use {
			t.Errorf("%s: got %+v %v", c.name, got, err)
		}
	}
}

func TestIPv4SubnetOf_errors(t *testing.T) {
	for _, p := range []int{-1, 33} {
		if _, err := IPv4SubnetOf([4]byte{1, 2, 3, 4}, p); !errors.Is(err, ErrPrefixLength) {
			t.Errorf("prefix %d: %v", p, err)
		}
	}
}

func BenchmarkIPv4SubnetOf(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = IPv4SubnetOf([4]byte{192, 168, 37, 200}, 22)
	}
}
