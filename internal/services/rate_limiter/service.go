package rate_limiter

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"time"
	"url_shortener/configs"
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
	config                  *configs.Config
	l                       logger.Logger
	userRateLimiterByUserId map[int]*UserRateLimiter
	mu                      *sync.Mutex
	done                    chan bool
}

func NewRateLimiterService(config *configs.Config, l logger.Logger) RateLimiterService {
	return &rateLimiterService{config, l, make(map[int]*UserRateLimiter), &sync.Mutex{}, make(chan bool)}
}

func (r *rateLimiterService) Check(userIdStr string) error {
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return fmt.Errorf(custom_errors.InvalidUserIdError)
	}
	r.mu.Lock()
	userRateLimiter, ok := r.userRateLimiterByUserId[userId]
	if !ok {
		userRateLimiter = &UserRateLimiter{limiter: DefaultAPILimiter(r.config, r.l)}
		r.userRateLimiterByUserId[userId] = userRateLimiter
		r.l.Debug("Create new limiter for", userId)
	}
	r.userRateLimiterByUserId[userId].lastSeen = time.Now()
	r.l.Debug(userId, userRateLimiter.lastSeen)
	if !userRateLimiter.limiter.Allow() {
		r.mu.Unlock()
		return fmt.Errorf(custom_errors.RateLimitError)
	}
	r.mu.Unlock()
	return nil
}

func leakMemory() {
	var leaked = make([][]byte, 0)
	for {
		// Only declare
		leaked = append(leaked, make([]byte, 1<<20))
		time.Sleep(100 * time.Millisecond)

		// Does store actual bytes
		// internal := make([]byte, 1<<20)
		// for i := 0; i < len(internal); i++ {
		// 	internal[i] = 'A'
		// }
		// leaked = append(leaked, internal)
		// time.Sleep(100 * time.Millisecond)
	}
}

func instantleakMemory() {
	var leaked = make([][]byte, 0)
	for i := 0; i < 10000; i++ {
		leaked = append(leaked, make([]byte, 1<<20))
	}
	fmt.Println("Finished")
	select {}

}

func (r *rateLimiterService) cleanUnusedRateLimiters() {
	// go leakMemory()
	// go instantleakMemory()
	// go blockGoroutine()
	go func() {
		ticker := time.Tick(5 * time.Second)
		var mem runtime.MemStats
		for {
			t := <-ticker
			r.l.Info(t)
			runtime.ReadMemStats(&mem)
			r.l.Info(fmt.Sprintf("Initial HeapAlloc: %v KB", mem.HeapAlloc/1024))
			r.l.Info(fmt.Sprintf("Initial TotalAlloc: %v KB", mem.TotalAlloc/1024))
			r.l.Info(fmt.Sprintf("StackInuse: %.2f KB", float64(mem.StackInuse)/1024))
			r.l.Info(fmt.Sprintf("StackSys:   %.2f KB\n", float64(mem.StackSys)/1024))
		}
	}()
	ticker := time.NewTicker(r.config.RateLimiter.CleaningPeriod)
	r.l.Info("Start clean rate limiters goroutine")
	for {
		// Generate some users
		count := 500000
		for i := 0; i < count; i++ {
			userRateLimiter := &UserRateLimiter{lastSeen: time.Now(), limiter: DefaultAPILimiter(r.config, r.l)}
			r.userRateLimiterByUserId[i] = userRateLimiter
			if i%100000 == 0 {
				r.l.Info("Start", len(r.userRateLimiterByUserId))
			}
		}
		break
	}

	for {
		select {
		case <-r.done:
			r.l.Info("Stop clean rate limiter goroutine")
			return
		case t := <-ticker.C:
			r.l.Debug("Current iteration", t, len(r.userRateLimiterByUserId))
			r.mu.Lock()
			for userId, userRateLimiter := range r.userRateLimiterByUserId {
				if time.Since(userRateLimiter.lastSeen) > r.config.RateLimiter.ExpiresPeriod {
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
