# Voltage Drop Calculator (row 28)

CLASSIFY: in scope for the formula with a caller-supplied conductor resistance. The NEC conductor tables (NFPA 70
Chapter 9 Tables 8 and 9: resistance and reactance by size, material and conduit) are a copyrighted standard's
data: out of scope here (NEEDS SOURCE / licence decision for Devin).
FIELDS (calform): necmaterial copper|aluminum, necwiresize 0..20 (14 AWG..1000 kcmil, default 6 = 3 AWG),
necconduit pvc|aluminum|steel, necpf 0.85; material copper|aluminum|silver|gold, wiresize (Ω/km values, default
0.4066 = 1 AWG), resistance 1.2, resistanceunit Ω/km (Ω/m, mΩ/m, Ω/1000ft, Ω/ft, mΩ/ft), voltage 120, phase
dc|ac|ac3, noofconductor 1..9 (parallel sets), distance 500, distanceunit feet|meters|miles|kilometers, amperes 1.
MATCH: NONE. BUILD eng/electrical (planned ELEC): ENG-ELEC-007 VoltageDropTwoWire (ΔV = I·R'·2L, Ohm's law on
the out-and-back loop, OpenStax University Physics Vol. 2 §9.3-9.4), ENG-ELEC-008 VoltageDropThreePhase
(√3·I·R'·L, Glover et al., Power System Analysis and Design 5th ed. §2.5). Reactance/power factor not modelled.
COMPARE: one run, 10 A, 500 ft, 120 V, AC single phase, with "custom" resistance 1.2 Ω/km requested. The driver
could not select the custom-resistance mode (missing: wiresize=custom) and the page computed from a mode we cannot
see: theirs 2.44 V (2.03 %), ours for 1.2 Ω/km 3.65760 V. Not comparable (different conductor resistance);
2.44 V corresponds to 0.244 Ω per 1000 ft, i.e. their NEC-table 3 AWG copper default, which we do not hold.
DRIVER: the resistance-source mode (NEC table / wire size / custom) is not a fillable field. Session 2 re-read the
fields with ids: still only necmaterial/necwiresize/necconduit/necpf, material/wiresize (14 AWG..40 AWG, no
"custom" option), resistance/resistanceunit, voltage, phase, noofconductor, distance, amperes; no element by
name or id selects the mode (likely an unnamed tab). Not re-submitted.
Submissions: 1 (session 1).
CLOSE: BUILT ENG-ELEC-007, ENG-ELEC-008 (compare blocked by DRIVER; NEC tables NEEDS SOURCE).
