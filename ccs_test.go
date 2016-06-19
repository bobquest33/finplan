package main

import (
	"reflect"
	"testing"
	"time"
)

func TestCCMinimum(t *testing.T) {
	var zeroTime time.Time
	cc := cc{amps: map[time.Time]percentage{zeroTime: 1000}}

	if want, got := money(2500), cc.minimum(300000, time.Now()); got != want {
		t.Errorf("Wanted %d, got %d", want, got)
	}

}

func TestMostRecent(t *testing.T) {
	t1 := time.Now().AddDate(-2, 0, 0)
	t2 := time.Now().AddDate(-1, 0, 0)
	t3 := time.Now().AddDate(1, 0, 0)
	t4 := time.Now().AddDate(2, 0, 0)

	cc := cc{
		aprs: map[time.Time]percentage{t1: 1000, t2: 2000, t3: 3000, t4: 4000},
	}

	want := percentage(2000)
	got := cc.mostRecent(cc.aprs, time.Now())

	if got != want {
		t.Errorf("Wanted %f, got %f", want, got)
	}
}

func TestByApr(t *testing.T) {
	var zeroTime time.Time
	cc1 := cc{aprs: map[time.Time]percentage{zeroTime: 3000}}
	cc2 := cc{aprs: map[time.Time]percentage{zeroTime: 1000}}
	cc3 := cc{aprs: map[time.Time]percentage{zeroTime: 2000}}
	ccs := ccs{"cc1": cc1, "cc2": cc2, "cc3": cc3}

	got1, got2 := ccs.byAPR(time.Now())

	want1 := []percentage{1000, 2000, 3000}
	want2 := map[percentage]cc{1000: cc2, 2000: cc3, 3000: cc1}

	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("Wanted\n%#v\ngot\n%#v", want1, got1)
	}
	if !reflect.DeepEqual(want2, got2) {
		t.Errorf("Wanted\n%#v\ngot\n%#v", want2, got2)
	}
}

func TestHighest(t *testing.T) {
	var zeroTime time.Time
	cc1 := cc{name: "cc1", aprs: map[time.Time]percentage{zeroTime: 4000}}
	cc2 := cc{name: "cc2", balance: 100000, aprs: map[time.Time]percentage{zeroTime: 1000}}
	cc3 := cc{name: "cc3", balance: 100000, aprs: map[time.Time]percentage{zeroTime: 3000}}
	cc4 := cc{name: "cc4", balance: 100000, aprs: map[time.Time]percentage{zeroTime: 2000}}
	ccs := ccs{"cc1": cc1, "cc2": cc2, "cc3": cc3, "cc4": cc4}

	got := ccs.highest(time.Now())

	want := cc3

	if want.name != got.name {
		t.Errorf("Wanted  %s got %s", want.name, got.name)
	}
}

func TestPayments(t *testing.T) {
	var zeroTime time.Time
	cc1 := cc{name: "cc1", aprs: map[time.Time]percentage{zeroTime: 4000}}
	cc2 := cc{
		name:    "cc2",
		balance: 100000,
		aprs:    map[time.Time]percentage{zeroTime: 5},
		amps:    map[time.Time]percentage{zeroTime: 1000}}
	cc3 := cc{
		name:    "cc3",
		balance: 200000,
		aprs:    map[time.Time]percentage{zeroTime: 7},
		amps:    map[time.Time]percentage{zeroTime: 3000},
	}
	cc4 := cc{
		name:    "cc4",
		balance: 300000,
		aprs:    map[time.Time]percentage{zeroTime: 1},
		amps:    map[time.Time]percentage{zeroTime: 2000},
	}
	ccs := ccs{"cc1": cc1, "cc2": cc2, "cc3": cc3, "cc4": cc4}

	const available = money(300000)

	payments := ccs.payments(time.Now(), available)

	if _, ok := payments["cc1"]; ok {
		t.Errorf("Did not expect a payment for cc1.")
	}

	gotcc2 := payments["cc2"]

	if want := money(100000 * 1000 / 100 / 12 / 100); gotcc2 != want {
		t.Errorf("Wanted %d, but got %d", want, gotcc2)
	}

	gotcc4 := payments["cc4"]

	if want := money(300000 * 2000 / 100 / 12 / 100); gotcc4 != want {
		t.Errorf("Wanted %d, but got %d", want, gotcc4)
	}

	if got, want := payments["cc3"], money(200000); got != want {
		t.Errorf("Wanted %d, but got %d", want, got)
	}

}
