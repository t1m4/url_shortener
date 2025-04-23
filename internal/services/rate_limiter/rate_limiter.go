package rate_limiter

import (
	"context"
	"sort"
	"time"
	"url_shortener/internal/logger"

	"golang.org/x/time/rate"
)

type RateLimiter interface {
	Wait(context.Context) error
	Limit() rate.Limit
	Tokens() float64
	Allow() bool
}

type multiLimiter struct {
	limiters []*rate.Limiter
}

func New(limiters ...*rate.Limiter) RateLimiter {
	byLimit := func(i, j int) bool {
		return limiters[i].Limit() < limiters[j].Limit()
	}
	sort.Slice(limiters, byLimit)
	return &multiLimiter{limiters: limiters}
}

func DefaultAPILimiter(l logger.Logger) RateLimiter {
	return New(
		// TODO take from config
		rate.NewLimiter(Per(l, 2, time.Second), 1),
		rate.NewLimiter(Per(l, 10, time.Minute), 10),
		rate.NewLimiter(Per(l, 50, time.Hour), 50),
	)
}

func (l *multiLimiter) Wait(ctx context.Context) error {
	for _, l := range l.limiters {
		if err := l.Wait(ctx); err != nil {
			return err
		}
	}
	return nil
}
func (l *multiLimiter) Limit() rate.Limit {
	return l.limiters[0].Limit()
}
func (l *multiLimiter) Tokens() float64 {
	return l.limiters[0].Tokens()
}
func (l *multiLimiter) Allow() bool {
	return l.limiters[0].Allow()
}
