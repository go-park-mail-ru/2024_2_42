package utils

import "time"

func YearsBetween(t1, t2 time.Time) int {
	if t1.Location() != t2.Location() {
		t2 = t2.In(t1.Location())
	}
	if t1.After(t2) {
		t1, t2 = t2, t1
	}
	y1, _, _ := t1.Date()
	y2, _, _ := t2.Date()

	return int(y2 - y1)
}
