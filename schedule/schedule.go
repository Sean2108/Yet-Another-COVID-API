package schedule

import "time"

// CallFunctionDaily : call functionToCall every day at hourToCallAt
func CallFunctionDaily(functionToCall func(), hourToCallAt int) {
	functionToCall()
	time.AfterFunc(getTimeTillUpdate(hourToCallAt, time.Now()), func() { CallFunctionDaily(functionToCall, hourToCallAt) })
}

func getTimeTillUpdate(hourToCallAt int, now time.Time) time.Duration {
	nextUpdate := time.Date(now.Year(), now.Month(), now.Day(), hourToCallAt, 0, 0, 0, time.UTC)
	if now.After(nextUpdate) {
		nextUpdate = nextUpdate.Add(24 * time.Hour)
	}
	return nextUpdate.Sub(now)
}
