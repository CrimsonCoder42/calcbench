package pediatrics

import (
	_ "embed"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// ErrTable is returned when the embedded WHO table cannot be read; it means the data file in data/ was
// altered, never a bad input.
var ErrTable = errors.New("pediatrics: embedded reference table is malformed")

//go:embed data/who_weight_velocity_1mon_p.csv
var whoWeightVelocity1MonCSV string

// WHOIncrementPercentiles is one row of a WHO growth-velocity percentile table: the interval's label as the
// WHO file prints it and the 1st…99th percentiles of the increment over that interval.
type WHOIncrementPercentiles struct {
	Interval                                           string
	P1, P3, P5, P15, P25, P50, P75, P85, P95, P97, P99 float64
}

var (
	whoWeightVelocityOnce  sync.Once
	whoWeightVelocityTable [2][12]WHOIncrementPercentiles // [0] boys, [1] girls
	whoWeightVelocityErr   error
)

// parseWHOWeightVelocity fills whoWeightVelocityTable from the embedded CSV (header + 24 rows).
func parseWHOWeightVelocity() {
	lines := strings.Split(strings.TrimSpace(whoWeightVelocity1MonCSV), "\n")
	if len(lines) != 25 {
		whoWeightVelocityErr = fmt.Errorf("%w: %d lines, want 25", ErrTable, len(lines))
		return
	}
	var seen [2][12]bool
	for _, line := range lines[1:] {
		f := strings.Split(strings.TrimRight(line, "\r"), ",")
		if len(f) != 18 {
			whoWeightVelocityErr = fmt.Errorf("%w: %d fields in %q", ErrTable, len(f), line)
			return
		}
		sex := 0
		switch f[0] {
		case "boys":
		case "girls":
			sex = 1
		default:
			whoWeightVelocityErr = fmt.Errorf("%w: sex %q", ErrTable, f[0])
			return
		}
		k, err := strconv.Atoi(f[1])
		if err != nil || k < 0 || k > 11 || seen[sex][k] {
			whoWeightVelocityErr = fmt.Errorf("%w: interval_index %q", ErrTable, f[1])
			return
		}
		seen[sex][k] = true
		var p [11]float64
		for i := range p {
			if p[i], err = strconv.ParseFloat(f[7+i], 64); err != nil {
				whoWeightVelocityErr = fmt.Errorf("%w: %v", ErrTable, err)
				return
			}
		}
		whoWeightVelocityTable[sex][k] = WHOIncrementPercentiles{f[2], p[0], p[1], p[2], p[3], p[4], p[5], p[6], p[7], p[8], p[9], p[10]}
	}
}

// WHOWeightIncrement1MonthPercentiles returns the WHO Child Growth Standards percentiles (1st, 3rd, 5th,
// 15th, 25th, 50th, 75th, 85th, 95th, 97th, 99th) of the weight increment, in grams, over one 1-month
// interval between birth and 12 months, for boys or girls. intervalIndex picks the WHO file's row:
// 0 = "0 – 4 wks", 1 = "4 wks – 2 mo", k = "k – k+1 mo" for k = 2…11. Values are the file's, unrounded
// and uninterpolated; negative values (a weight loss over the interval) are the file's own.
// Errors: ErrOutOfRange for intervalIndex outside 0–11; ErrTable if the embedded file is malformed.
//
// ID:      MED-PEDS-008
// Time:    O(1) per call (a one-time parse of the fixed 24-row table on first use)
// Space:   O(1): a fixed 24-row table; no allocation after the first call
// Source:  WHO Multicentre Growth Reference Study Group. WHO Child Growth Standards: Growth velocity based on weight, length and head circumference: Methods and development. Geneva: World Health Organization; 2009. ISBN 978 92 4 154763 5
// Source:  World Health Organization. Weight velocity, 1-month increments, percentiles, boys and girls (ttt-weight-boys-1mon-p.xlsx, ttt-weight-girls-1mon-p.xlsx). https://www.who.int/tools/child-growth-standards/standards/weight-velocity, read 2026-11-12 (files, URLs and SHA-256 in med/pediatrics/data/README.md)
// Tags:    pediatrics, growth, velocity, weight, who, infant, percentile, reference table
// Units:   intervalIndex row number 0–11; female selects the girls' table; result grams per interval
// Range:   birth to 12 months, 1-month intervals (the WHO table's extent)
// Example: boys, "0 – 4 wks" → P50 1023 g, P3 369 g, P97 1575 g; girls, "11 – 12 mo" → P50 232 g, P1 −222 g (WHO files)
func WHOWeightIncrement1MonthPercentiles(female bool, intervalIndex int) (WHOIncrementPercentiles, error) {
	if intervalIndex < 0 || intervalIndex > 11 {
		return WHOIncrementPercentiles{}, fmt.Errorf("%w: intervalIndex = %d, want 0–11", ErrOutOfRange, intervalIndex)
	}
	whoWeightVelocityOnce.Do(parseWHOWeightVelocity)
	if whoWeightVelocityErr != nil {
		return WHOIncrementPercentiles{}, whoWeightVelocityErr
	}
	sex := 0
	if female {
		sex = 1
	}
	return whoWeightVelocityTable[sex][intervalIndex], nil
}
