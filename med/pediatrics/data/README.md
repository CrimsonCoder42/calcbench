# med/pediatrics/data

Official data files of the publishing body, as downloaded, plus the one CSV the code embeds.

## WHO Child Growth Standards: weight velocity, 1-month increments, percentiles (birth to 12 months)

Publisher: World Health Organization, Department of Nutrition. Methods: WHO Multicentre Growth Reference Study
Group. WHO Child Growth Standards: Growth velocity based on weight, length and head circumference: Methods and
development. Geneva: World Health Organization; 2009. ISBN 978 92 4 154763 5.
Index page: https://www.who.int/tools/child-growth-standards/standards/weight-velocity (read 2026-11-12).

| file | URL | read | SHA-256 |
|---|---|---|---|
| who_weight_velocity_boys_1mon_p.xlsx | https://cdn.who.int/media/docs/default-source/child-growth/child-growth-standards/indicators/weight-velocity/ttt-weight-boys-1mon-p.xlsx?sfvrsn=9ecec326_5 | 2026-11-12 | a0a3e36df5879049c61db527892601ffb281dfc078ee1480df94dd7eaa2e5a14 |
| who_weight_velocity_girls_1mon_p.xlsx | https://cdn.who.int/media/docs/default-source/child-growth/child-growth-standards/indicators/weight-velocity/ttt-weight-girls-1mon-p.xlsx?sfvrsn=b6d51fd6_5 | 2026-11-12 | 30ed019f8a4310cdede1c1a190aa93079c9d348a9a9b243d78dcf8f4df191673 |
| who_weight_velocity_1mon_p.csv (embedded) | derived from the two files above | 2026-11-12 | 448da5ad3cfedb5d808e3804c7f72672f789e33f0b6b6d6d255a36172a947d4d |

The CSV is the two workbooks' only sheet, cell for cell (the cell text Excel stores, no rounding), with two
columns added in front: `sex` (boys / girls) and `interval_index` (0 = "0 – 4 wks", 1 = "4 wks – 2 mo",
k = "k – k+1 mo" for k = 2…11). Conversion: unzip `xl/worksheets/sheet1.xml` and `xl/sharedStrings.xml`
and write each row's cells in column order (Python standard library `zipfile` + `xml.etree`), then prefix the
two columns with awk. Columns: Interval, L, M, S, Delta, then the 1st, 3rd, 5th, 15th, 25th, 50th, 75th,
85th, 95th, 97th, 99th percentiles of the weight increment over the interval, in grams. The workbook header
"85th " carries a trailing space; it is kept.

`med/pediatrics/who_weight_velocity_1mon_test.go` checks that the embedded CSV's SHA-256 equals the value in
this table and that at least five published cells per sex read back unchanged.
