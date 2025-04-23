package rate_limiter

import (
	"time"
	"url_shortener/internal/logger"

	"golang.org/x/time/rate"
)

func Per(l logger.Logger, eventCount int, duration time.Duration) rate.Limit {
	/*
		return rate.Limit - number of events per second
	*/
	l.Debug(eventCount, duration, time.Duration(eventCount), duration/time.Duration(eventCount), rate.Every(duration/time.Duration(eventCount)))
	return rate.Every(duration / time.Duration(eventCount))
}
