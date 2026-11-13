//go:build ignore

package main

import (
	"fmt"
	"time"

	"github.com/CrimsonCoder42/calcbench/med/obstetric"
)

func main() {
	today := time.Date(2026, 11, 13, 0, 0, 0, 0, time.UTC)
	for _, lmp := range []time.Time{
		time.Date(2026, 9, 14, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 7, 4, 0, 0, 0, 0, time.UTC),
	} {
		w, d, err := obstetric.GestationalAgeFromLMP(lmp, today)
		fmt.Println("LMP", lmp.Format(time.DateOnly), "EDD", obstetric.DueDateFromLMP(lmp).Format(time.DateOnly), "GA on 2026-11-13", w, "w", d, "d", err)
	}
}
