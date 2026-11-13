// Package computing holds the arithmetic of computer networks and systems: IPv4 subnetting by CIDR prefix and
// the transmission time of data over a link.
//
// ID short form: COMPUT (IDs are ENG-COMPUT-NNN).
//
// Units: IPv4 addresses are [4]byte in network order (from net/netip: Addr.As4()); a prefix length is a bit
// count 0..32; counts of addresses are uint64. A prefix length outside 0..32 returns ErrPrefixLength. Sizes are in bits and rates in bit/s; NaN or ±Inf returns
// ErrNotFinite, a negative value ErrNegative, a zero rate ErrZero.
//
// Functions:
//
//	ENG-COMPUT-001  IPv4SubnetOf      network, broadcast, mask, wildcard, host range and counts of an address/prefix
//	ENG-COMPUT-002  TransmissionTime  L / R, seconds to send L bits at R bit/s
package computing
