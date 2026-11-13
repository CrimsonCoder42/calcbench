# other: discrepancies (A-OTHER)

| slug | inputs | ours | theirs | verdict | citation |
|---|---|---|---|---|---|
| ip-subnet-calculator | 10.0.0.1/31 | 2 usable hosts (10.0.0.0, 10.0.0.1) | 0 usable, range NA | THEIRS DIFFERS | RFC 3021 §2 (31-bit prefixes on point-to-point links: both addresses usable); theirs applies the RFC 950 total − 2 rule |
