package main

import (
	"time"
)

type month struct {
	T      MonthTime
	Netrev money
	CCs    map[string]ccmonth
}

type ccmonth struct {
	CC      *cc
	Payment money // the payment for a specific month
	Balance money
	RI      money // running interest
}

func (month *month) pay(cc *cc, amount money) {
	month.Netrev -= amount
	cc.pay(amount)
}

// ri calculates the running interest for the month
func (month month) RI() money {
	ret := money(0)
	for _, cc := range month.CCs {
		ret += cc.RI
	}
	return ret
}

func (month month) Format() string {
	return month.T.Format("Jan 06")
}

func project(revexps revexps, ccs ccs, nrOfMonths int) []month {

	ret := make([]month, 0, nrOfMonths)

	var t MonthTime
	var prevNetrev money

	// make a projection for each month

	for mi := 0; mi < nrOfMonths; mi++ {

		// create month and set date
		if mi == 0 {
			t = normalize(time.Now())
		} else {
			t = ret[mi-1].T
		}
		month := month{T: MonthTime{t.AddDate(0, 1, 0)}, CCs: make(map[string]ccmonth)}
		netrev := revexps.netrev(t)
		month.Netrev = prevNetrev + netrev

		// calculate payments to be made
		payments := ccs.payments(month.T, month.Netrev)

		// make cc payments and add monthcc projections
		for name, cc := range ccs {
			// make payment
			payment := payments[name]
			month.pay(cc, payment)

			// get previous month's running interest for this card
			ri := money(0)
			if mi > 0 {
				ri = ret[mi-1].CCs[name].RI
			}

			// add cc projection
			ccmonth := ccmonth{
				CC:      cc,
				Payment: payment,
				Balance: cc.Balance,
				RI:      ri + cc.interest(cc.Balance, month.T),
			}
			month.CCs[name] = ccmonth
		}

		prevNetrev = month.Netrev

		ret = append(ret, month)
	}
	return ret
}
