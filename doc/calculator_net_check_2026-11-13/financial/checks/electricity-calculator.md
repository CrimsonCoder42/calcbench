# electricity-calculator (row 73, HANDOFF from A-OTHER; their note: other/checks/electricity-calculator.md)

FIELDS (form `calform`): appliance preset select, power (1000) + powerunit (W|kW|btu|hp|ton), capacity (100 % duty),
usage (24) + usageunit (mpd|hpd|hpw|hpm|dpw|dpm|dpy|mpy), price (0.15 $/kWh).

MATCH: the energy part (kWh = kW × duty × hours) is physics, noted by A-OTHER (E = P·t; 1 kWh = 3.6 MJ). The money
part is one product, kWh × price per kWh, rounded to the cent (half-even, D-022). No fin function built: there is no
formula beyond the product, and a generic "quantity × unit price in cents" helper is a naming/API decision for the
lead (proposal in LEAD.md). Compared by hand.

COMPARE (hand: kWh × cents per kWh)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 1,000 W, 100 %, 24 h/day, 15 ¢/kWh | 24 kWh/day → 360 ¢ = $3.60; 8,766 kWh/yr (365.25 d) → $1,314.90 | $3.60; $1,314.90 | 0 | MATCH |
| 1.5 kW, 60 %, 5 h/day, 17.3 ¢/kWh | 4.5 kWh/day → 77.85 ¢ → $0.78; week 31.5 kWh → $5.45; year 1,643.6 kWh → $284.35 | $0.78; $5.45; $284.35 | 0 | MATCH |

Their month uses 30.44 days (730.6 kWh, $109.58); 365.25/12 = 30.4375 gives 730.5 kWh and $109.575 → $109.58 (half-even)
too. hp and refrigeration-ton inputs need a power unit package (A-OTHER's note; lead decision).
Submissions: 2.
