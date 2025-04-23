package rate_limiter

import (
	"fmt"
	"strconv"
	"sync"
	"time"
	"url_shortener/internal/custom_errors"
	"url_shortener/internal/logger"
)

type RateLimiterService interface {
	Check(string) error
	Start()
	Stop()
}

type UserRateLimiter struct {
	lastSeen time.Time
	limiter  RateLimiter
}
type rateLimiterService struct {
	l                       logger.Logger
	userRateLimiterByUserId map[int]*UserRateLimiter
	mu                      *sync.Mutex
	done                    chan bool
}

func NewRateLimiterService(l logger.Logger) RateLimiterService {
	return &rateLimiterService{l, make(map[int]*UserRateLimiter), &sync.Mutex{}, make(chan bool)}
}

func (r *rateLimiterService) Check(userIdStr string) error {
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return fmt.Errorf(custom_errors.InvalidUserIdError)
	}
	r.mu.Lock()
	userRateLimiter, ok := r.userRateLimiterByUserId[userId]
	if !ok {
		userRateLimiter = &UserRateLimiter{limiter: DefaultAPILimiter(r.l)}
		r.userRateLimiterByUserId[userId] = userRateLimiter
		r.l.Debug("Create new limiter for", userId)
	}
	r.userRateLimiterByUserId[userId].lastSeen = time.Now()
	r.l.Debug(userId, userRateLimiter.lastSeen, userRateLimiter.limiter.Limit(), userRateLimiter.limiter.Tokens())
	if !userRateLimiter.limiter.Allow() {
		r.mu.Unlock()
		return fmt.Errorf(custom_errors.RateLimitError)
	}
	r.mu.Unlock()
	return nil
}

func (r *rateLimiterService) cleanUnusedRateLimiters() {
	// TODO move to the config
	ticker := time.NewTicker(5 * time.Second)
	r.l.Info("Start clean rate limiters goroutine")
	for {
		select {
		case <-r.done:
			r.l.Info("Stop clean rate limiter goroutine")
			return
		case t := <-ticker.C:
			r.l.Debug("Current iteration", t, len(r.userRateLimiterByUserId))
			r.mu.Lock()
			for userId, userRateLimiter := range r.userRateLimiterByUserId {
				// TODO move to the config
				if time.Since(userRateLimiter.lastSeen) > 30*time.Second {
					r.l.Debug("Delete userId", userId, userRateLimiter.lastSeen, t)
					delete(r.userRateLimiterByUserId, userId)
				}
			}
			r.mu.Unlock()
		}

	}
}

func (r *rateLimiterService) Start() {
	go r.cleanUnusedRateLimiters()
}
func (r *rateLimiterService) Stop() {
	close(r.done)
}
