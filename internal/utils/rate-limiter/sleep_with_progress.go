package ratelimiter

import (
	"time"
)

type gradualWaitThreshold struct {
	maxRemaining time.Duration
	sleepFor     time.Duration
}

var waitThresholds = []gradualWaitThreshold{
	{
		maxRemaining: time.Duration(time.Minute),
		sleepFor:     time.Duration(10 * time.Second),
	},
	{
		maxRemaining: time.Duration(10 * time.Minute),
		sleepFor:     time.Duration(1 * time.Minute),
	},
	{
		maxRemaining: time.Duration(30 * time.Minute),
		sleepFor:     time.Duration(5 * time.Minute),
	},
	{
		maxRemaining: time.Duration(1 * time.Hour),
		sleepFor:     time.Duration(10 * time.Minute),
	},
}

func sleepWithProgressFunc(totalSleepTime time.Duration, progress func(sleepTime, remaining time.Duration)) {
	sleepTimeRemaining := totalSleepTime

	for sleepTimeRemaining > 0 {
		var sleepTime time.Duration = waitThresholds[len(waitThresholds)-1].sleepFor

		// select "local" sleep time based on remaining duration
		for _, threshold := range waitThresholds {
			if sleepTimeRemaining < threshold.maxRemaining {
				sleepTime = threshold.sleepFor
				break
			}
		}
		// so e.g. the minimum wait time from the table above is 10 seconds.
		// however if we only have to wait, say, 3 seconds in total - there's
		// no point waiting 10.
		if sleepTime > sleepTimeRemaining {
			sleepTime = sleepTimeRemaining
		}
		progress(sleepTime, sleepTimeRemaining)
		time.Sleep(sleepTime)
		sleepTimeRemaining = sleepTimeRemaining - sleepTime
	}
}
