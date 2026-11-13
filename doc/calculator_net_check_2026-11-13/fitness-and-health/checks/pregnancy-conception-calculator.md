# pregnancy-conception-calculator

## FIELDS
Not loaded (nothing built, nothing to compare; saves a page load). Their date inputs use the date pickers the driver cannot fill (LEAD.md).

## MATCH
NONE (MED-OB-003/004 cover LMP → due date and gestational age only).

## RESEARCH → NEEDS SOURCE (decision for Devin)
These calculators place ovulation 14 days before the next expected period (cycle length − 14 from the LMP) and derive conception,
the fertile window and future periods from it.
- Wilcox AJ, Dunson D, Baird DD. The timing of the "fertile window" in the menstrual cycle: day specific estimates from a prospective
  study. BMJ 2000;321(7271):1259-1262, doi:10.1136/bmj.321.7271.1259 (PMC27529, read 2026-11-13): defines the fertile window as "the
  five days before ovulation and the day of ovulation itself" and day 1 as the onset of bleeding, but names "ovulation occurs 14 days
  before the next menses" as a clinical-guideline assumption its data contradict (true in about 10 % of 28-day cycles).
- ACOG CO 700 (2017) copy read gives LMP + 280 only; the conception (266 d) and ART (261/263 d) counts were not visible.
No citable source states the 14-day luteal rule as the method, so nothing is built. Options for Devin in LEAD.md.

## Session 2 (dates via hidden fields dueDateDay / cLastPeriodDay, MM/DD/YYYY)
FIELDS (form calform): ccalcmethod due date (default) | last period | ultrasound · dueDateDay (default 2027-03-13) · cLastPeriodDay · cCycleLength 28 ·
cUltrasoundDate, cUltrasoundWeek 22, cUltrasoundDay 6. Results: "most probable" and "possible" conception and intercourse date windows only.
| inputs | ours | theirs | verdict |
|---|---|---|---|
| due date 2027-01-20 | LMP = EDD − 280 = 2026-04-15 (MED-OB-003 inverse; no conception function) | most probable conception Apr 27 – May 1, 2026; possible Apr 26 – May 6 | not comparable |
| LMP 2026-03-10, cycle 28 | EDD 2026-12-15; no conception function | most probable conception Mar 22 – 26; possible Mar 21 – 31 | not comparable |
Their windows sit around LMP + 14 days (12–16 "most probable", 11–21 "possible") — the ovulation-at-day-14 convention the row is NEEDS
SOURCE for (Wilcox BMJ 2000; LEAD.md). Row stays NEEDS SOURCE. Submissions: session 2: 1 fields + 2 runs.
