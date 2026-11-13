# IP Subnet Calculator (row 32)

CLASSIFY: in scope (IPv4 CIDR arithmetic). IPv6 form (calform2: c6subnet 1..128 default 64, c6ip) not built
this session. "IP class" and "IP type (private)" outputs are registry tables (RFC 1918 ranges), not built.
FIELDS (calform): cclass radio (any/A/B/C), csubnet select (prefix 1..32, default 30), cip text (default
74.71.54.244).
MATCH: NONE. BUILD eng/computing (planned COMPUT): ENG-COMPUT-001 IPv4SubnetOf. Sources: RFC 4632 §3.1, RFC 950
(all-zeros/all-ones host field), RFC 3021 §2 (/31 point-to-point).
COMPARE (ours: hand values in ipv4_subnet_test.go):

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 192.168.37.200/22 | net 192.168.36.0, bc 192.168.39.255, hosts .36.1-.39.254, 1024 total, 1022 usable, mask 255.255.252.0, wildcard 0.0.3.255 | same, all fields | 0 | MATCH |
| 10.0.0.1/31 | net 10.0.0.0, bc 10.0.0.1, 2 total, 2 usable (hosts 10.0.0.0-10.0.0.1) | 2 total, 0 usable, range NA | usable 2 vs 0 | THEIRS DIFFERS |

DISCREPANCY (/31 usable hosts): ours follows RFC 3021 §2 (Retana et al., 2000): on a 31-bit prefix both
addresses are host addresses (no network or broadcast address is reserved). Theirs applies the pre-RFC-3021
"total − 2" rule from RFC 950. Ours stays.
Submissions: 2.
CLOSE: BUILT ENG-COMPUT-001 (MATCH on /22; THEIRS DIFFERS on /31).
