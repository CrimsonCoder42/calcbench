//go:build ignore

package main

import (
	"fmt"
	"time"

	"github.com/CrimsonCoder42/calcbench/med/obstetric"
)

func main() {
	for _, l := range []time.Time{time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)} {
		fmt.Println("LMP", l.Format(time.DateOnly), "EDD", obstetric.DueDateFromLMP(l).Format(time.DateOnly))
	}
	w2, d2, _ := obstetric.GestationalAgeFromLMP(time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 11, 13, 0, 0, 0, 0, time.UTC))
	fmt.Println("LMP 2026-08-01 EDD", obstetric.DueDateFromLMP(time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)).Format(time.DateOnly), "on 2026-11-13:", w2, "w", d2, "d")
	w, d, err := obstetric.GestationalAgeFromLMP(time.Date(2026, 6, 6, 0, 0, 0, 0, time.UTC), time.Date(2026, 11, 13, 0, 0, 0, 0, time.UTC))
	fmt.Println("LMP 2026-06-06 on 2026-11-13:", w, "w", d, "d", err)
}
