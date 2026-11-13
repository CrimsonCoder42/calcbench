package computing

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// ErrPrefixLength reports a CIDR prefix length outside 0..32.
var ErrPrefixLength = errors.New("computing: IPv4 prefix length out of range")

// IPv4Subnet describes the subnet an IPv4 address belongs to under a prefix length.
type IPv4Subnet struct {
	Network   [4]byte // address with the host bits cleared
	Broadcast [4]byte // address with the host bits set (directed broadcast; for /31 and /32 the last address)
	Mask      [4]byte // prefixLen one bits, then zeros
	Wildcard  [4]byte // bitwise complement of Mask
	FirstHost [4]byte // first assignable address
	LastHost  [4]byte // last assignable address
	// TotalAddresses is 2^(32 − prefixLen).
	TotalAddresses uint64
	// UsableHosts is TotalAddresses − 2 for prefixLen ≤ 30 (network and broadcast reserved), 2 for /31
	// (point-to-point, RFC 3021) and 1 for /32.
	UsableHosts uint64
}

// IPv4SubnetOf returns the subnet of addr under a CIDR prefix length: mask = the first prefixLen bits set,
// network = addr AND mask, broadcast = network OR NOT mask; hosts run from network+1 to broadcast−1, except
// /31, whose two addresses are both hosts (RFC 3021), and /32, a single host.
// Errors: ErrPrefixLength for prefixLen outside 0..32.
//
// ID:      ENG-COMPUT-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Fuller V, Li T. Classless Inter-domain Routing (CIDR). RFC 4632, IETF, 2006, doi:10.17487/RFC4632, §3.1 (prefix and mask); Mogul J, Postel J. Internet Standard Subnetting Procedure. RFC 950, 1985, doi:10.17487/RFC0950, Appendix ("all ones" host field = broadcast, "all zeros" = this network); Retana A et al. Using 31-Bit Prefixes on IPv4 Point-to-Point Links. RFC 3021, 2000, doi:10.17487/RFC3021, §2
// Tags:    ipv4, subnet, cidr, netmask, broadcast, network, computing
// Units:   addr IPv4 address; prefixLen bits 0..32; counts in addresses
// Example: 192.168.37.200/22 → network 192.168.36.0, broadcast 192.168.39.255, mask 255.255.252.0, wildcard 0.0.3.255, hosts 192.168.36.1-192.168.39.254, 1024 addresses, 1022 usable
func IPv4SubnetOf(addr [4]byte, prefixLen int) (IPv4Subnet, error) {
	if prefixLen < 0 || prefixLen > 32 {
		return IPv4Subnet{}, fmt.Errorf("%w: prefixLen = %d, want 0..32", ErrPrefixLength, prefixLen)
	}
	a := binary.BigEndian.Uint32(addr[:])
	var mask uint32
	if prefixLen > 0 {
		mask = ^uint32(0) << (32 - prefixLen)
	}
	network := a & mask
	broadcast := network | ^mask
	total := uint64(1) << (32 - prefixLen)
	first, last, usable := network+1, broadcast-1, total-2
	switch prefixLen {
	case 31:
		first, last, usable = network, broadcast, 2
	case 32:
		first, last, usable = network, network, 1
	}
	return IPv4Subnet{
		Network:        be(network),
		Broadcast:      be(broadcast),
		Mask:           be(mask),
		Wildcard:       be(^mask),
		FirstHost:      be(first),
		LastHost:       be(last),
		TotalAddresses: total,
		UsableHosts:    usable,
	}, nil
}

func be(v uint32) [4]byte {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], v)
	return b
}
