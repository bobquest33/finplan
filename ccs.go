package main

import (
	"sort"
	"time"
)

type money int64    // monetary amount multiplied by 100
type percentage int // percentage multiplied by 100 for 2 decimal precision
type percentages []percentage

func (ps percentages) Len() int {
	return len(ps)
}
func (ps percentages) Less(i, j int) bool {
	return ps[i] < ps[j]
}
func (ps percentages) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

// aprs and amps define their values at a specific month and onwards until the next chronological entry
type cc struct {
	Name    string
	Balance money
	From    MonthTime
	APRs    map[string]percentage
	AMPs    map[string]percentage // annual minimum payment
}

func (cc cc) shouldpay(month time.Time) bool {
	return cc.Balance > 0.0 && normalize(time.Now()).Time.After(normalize(month).Time)
}

func (cc *cc) pay(amount money) {
	cc.Balance -= amount
}

// interest calculates interest for the month
func (cc cc) interest(balance money, now MonthTime) money {
	return money(cc.aprForMonth(now)) * balance / 12 / 10000
}

// minimum calculates the minimum for the month
func (cc cc) minimum(balance money, now MonthTime) money {
	return money(cc.ampForMonth(now)) * balance / 12 / 10000
}

// aprForMonth finds the most recent apr entry
func (cc cc) aprForMonth(now MonthTime) percentage {
	return cc.mostRecent(cc.APRs, now)
}

// ampForMonth finds the most recent amp entry
func (cc cc) ampForMonth(now MonthTime) percentage {
	return cc.mostRecent(cc.AMPs, now)
}

// mostRecent finds the most recent entry for target.
// This is for helping to find the current value for cc.aprs and cc.amps.
func (cc cc) mostRecent(target map[string]percentage, now MonthTime) percentage {
	var mostRecent MonthTime
	for month := range cc.APRs {
		parsedMonth := parseMonth(month)
		if parsedMonth.After(now.Time) {
			continue
		}
		if parsedMonth.After(mostRecent.Time) {
			mostRecent = parseMonth(month)
		}
	}
	return target[monthString(mostRecent)]
}

type ccs map[string]*cc

func (ccs ccs) add(cc *cc) {
	ccs[cc.Name] = cc
}

// byAPR sorts by APRs
func (ccs ccs) byAPR(month MonthTime) ([]percentage, map[percentage]*cc) {
	byAPR := make(map[percentage]*cc)
	aprs := make(percentages, 0, len(ccs))
	for _, cc := range ccs {
		apr := cc.aprForMonth(month)
		aprs = append(aprs, apr)
		byAPR[apr] = cc
	}
	sort.Sort(aprs)
	return aprs, byAPR
}

// highest finds the card with highest interest that still has a balance left over
func (ccs ccs) highest(month MonthTime) *cc {
	aprSorted, byAPR := ccs.byAPR(month)

	for i := len(aprSorted) - 1; i >= 0; i-- {
		rate := aprSorted[i] // this is the highest rate
		cc := byAPR[rate]    // this is the card with the highest rate
		if cc.Balance > 0 {
			return cc
		}
	}
	return &cc{}
}

// payments focuses payment on highest rate card that still has a balance
// while paying minimum on other cards
func (ccs ccs) payments(month MonthTime, available money) map[string]money {

	ret := make(map[string]money)

	highest := ccs.highest(month)

	for _, cc := range ccs {
		if cc.Balance == 0 || cc.Name == highest.Name {
			continue
		}
		min := cc.minimum(cc.Balance, month)
		available -= min
		ret[cc.Name] = min
	}

	ret[highest.Name] = min(highest.Balance, available)
		
	return ret
}

func min(a, b money) money {
	if a < b {
		return a
	}
	return b
}
