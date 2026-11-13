# Bandwidth Calculator (row 34)

CLASSIFY: in scope (transmission time L/R; data-size unit conversion; website bandwidth = views × size ×
redundancy).
FIELDS: calform: convertamount 500, convertunit MB (b, kb, mb, gb, tb, B, KB, MB, GB, TB). calform2:
downloadsize2 500, downloadsize2unit MB (B..TB), bandwidth2 5, bandwidth2unit mb (b..tb per second).
calform3: pageviews 5000 per day (second..year), pagesize 500 KB, redundancyf 2. calform4: usage 1000 GB,
bandwidth4 (empty) mb.
MATCH: NONE. BUILD eng/computing ENG-COMPUT-002 TransmissionTime (Kurose & Ross, Computer Networking 7th ed.,
§1.4.1, L/R). Their units observed as decimal SI (1 MB = 10⁶ B, 1 Mbit = 10⁶ bit), 1 B = 8 bit; unit/datasize
(planned DATA) would hold B↔bit and the SI/IEC prefixes (IEC 80000-13:2008); not built this session.
COMPARE (sizes converted by hand in decimal: 500 MB = 4e9 bit, 1 GB = 8e9 bit):

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 500 MB at 5 Mbit/s | 800 s | ~13 min 20 s | 0 | MATCH |
| 1 GB at 100 Mbit/s | 80 s | ~1 min 20 s | 0 | MATCH |

Submissions: 2.
CLOSE: BUILT ENG-COMPUT-002 (MATCH).
