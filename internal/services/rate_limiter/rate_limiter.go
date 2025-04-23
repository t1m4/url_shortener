package rate_limiter

import (
	"context"
	"sort"
	"url_shortener/configs"
	"url_shortener/internal/logger"

	"golang.org/x/time/rate"
)

type RateLimiter interface {
	Wait(context.Context) error
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

func DefaultAPILimiter(config *configs.Config, l logger.Logger) RateLimiter {
	limiters := make([]*rate.Limiter, 0)
	for _, limiterConfig := range config.RateLimiter.Limiters {
		limiter := rate.NewLimiter(Per(l, limiterConfig.EventCount, limiterConfig.Duration), limiterConfig.Burst)
		limiters = append(limiters, limiter)
	}
	return New(limiters...)
}

func (l *multiLimiter) Wait(ctx context.Context) error {
	for _, l := range l.limiters {
		if err := l.Wait(ctx); err != nil {
			return err
		}
	}
	return nil
}
func (l *multiLimiter) Allow() bool {
	for _, limiter := range l.limiters {
		isAllowed := limiter.Allow()
		if !isAllowed {
			return false
		}
	}
	return true
}
