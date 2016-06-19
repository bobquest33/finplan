package main

import ()

// revenues and expenses

type revexp struct {
	Timeframe  timeframe
	Val        money
	Begin, End MonthTime
	Save       bool
}

func (revexp revexp) active(now MonthTime) bool {
	return now.After(revexp.Begin.Time) && now.Before(revexp.End.Time)
}

type revexps map[string]revexp

// netrev calculates the net revenue for the month now.
// now is used to determine which revexps had expired,
// to exlude from the calculation.
func (revexps revexps) netrev(now MonthTime) money {
	var ret money
	for _, revexp := range revexps {
		// count out this revexp if it had expired
		if !revexp.active(now) {
			continue
		}

		// calculate net income
		val := revexp.Val
		if revexp.Timeframe == annual {
			val /= 12
		}
		ret += val
	}
	return ret
}
