package main

import (
	"testing"
	"time"
)

func TestNetrev(t *testing.T) {
	got := revexps{
		"rev1": revexp{
			timeframe: annual,
			val:       12000,
			end:       time.Now().AddDate(1, 0, 0),
		},
		"exp": revexp{
			val: -10000,
			end: time.Now().AddDate(1, 0, 0),
		},
		"rev2": revexp{
			val: 2000,
			end: time.Now().AddDate(1, 0, 0),
		},
	}.netrev(time.Now())

	want := money(-7000)

	if got != want {
		t.Errorf("Wanted %d, got %d", want, got)
	}
}
