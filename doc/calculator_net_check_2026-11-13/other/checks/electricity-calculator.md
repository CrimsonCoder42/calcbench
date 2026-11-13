# Electricity Calculator (row 31)

CLASSIFY: in scope (energy = power × time × duty; cost = energy × tariff).
FIELDS (calform): appliance preset select (sets power/usage), power 1000, powerunit W (kW, btu, hp, ton),
capacity 100 (% duty), usage 24, usageunit hpd (min/day, h/day, h/week, h/month, days/week, days/month,
days/year, min/year), price 0.15 ($/kWh).
MATCH: NONE. The physics part is one product (E = P·t, OpenStax University Physics Vol. 2 §9.5; 1 kWh =
3.6 × 10⁶ J exactly, NIST SP 811 B.8); the result the user wants is a money amount (int64 cents, rounding policy),
which is fin/** (A-FIN). The hp and "ton" (refrigeration) inputs need unit/power (not in the taxonomy list: a
lead decision). No submissions.
CLOSE: HANDOFF A-FIN (energy cost in cents; physics part noted here).
