package middleware

import (
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

// entry wraps a rate limiter with the last time it was seen.
// Used to evict stale limiters and prevent unbounded memory growth.
type entry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// limiters holds per-key token bucket limiters.
//
// Each call to IPRateLimit or UserRateLimit creates one limiters instance
// with its own background cleanup goroutine. Keep the number of distinct
// middleware registrations small (we currently have 3: global IP, auth IP,
// user) to avoid accumulating unnecessary goroutines.
type limiters struct {
	mu      sync.RWMutex
	entries map[string]*entry
	rps     rate.Limit
	burst   int
}

func newLimiters(rps float64, burst int) *limiters {
	l := &limiters{
		entries: make(map[string]*entry),
		rps:     rate.Limit(rps),
		burst:   burst,
	}
	go l.cleanup(5 * time.Minute)
	return l
}

// get returns the limiter for key, creating one if it does not exist.
// Uses RLock for the common read path; upgrades to Lock only on creation.
func (l *limiters) get(key string) *rate.Limiter {
	l.mu.RLock()
	e, ok := l.entries[key]
	l.mu.RUnlock()

	if ok {
		l.mu.Lock()
		e.lastSeen = time.Now()
		l.mu.Unlock()
		return e.limiter
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Double-check after acquiring write lock — another goroutine may have
	// created the entry between our RUnlock and Lock.
	if e, ok = l.entries[key]; ok {
		e.lastSeen = time.Now()
		return e.limiter
	}

	e = &entry{
		limiter:  rate.NewLimiter(l.rps, l.burst),
		lastSeen: time.Now(),
	}
	l.entries[key] = e
	return e.limiter
}

// cleanup removes entries that have not been seen for ttl.
// Runs in a background goroutine — one per limiters instance.
// Worst-case stale lifetime is up to 2×ttl (entry created just after
// a cleanup sweep, then not seen again before the next one).
func (l *limiters) cleanup(ttl time.Duration) {
	ticker := time.NewTicker(ttl)
	defer ticker.Stop()
	for range ticker.C {
		l.mu.Lock()
		for key, e := range l.entries {
			if time.Since(e.lastSeen) > ttl {
				delete(l.entries, key)
			}
		}
		l.mu.Unlock()
	}
}

// retryAfterSeconds calculates how long the caller should wait using
// Reserve(). If the request is rejected, Cancel() is called to release
// the reservation so rejected requests do not drain limiter tokens.
func retryAfterSeconds(l *rate.Limiter) (blocked bool, retryAfter string) {
	r := l.Reserve()
	if !r.OK() {
		return true, "1"
	}
	delay := r.Delay()
	if delay <= 0 {
		// Token was available — request should not have been rejected.
		// Cancel and signal not blocked.
		r.Cancel()
		return false, ""
	}
	// Request exceeds rate — cancel reservation and report wait time.
	r.Cancel()
	secs := int(math.Ceil(delay.Seconds()))
	if secs < 1 {
		secs = 1
	}
	return true, strconv.Itoa(secs)
}

// IPRateLimit limits requests by client IP.
//
// IP extraction order (Echo default):
//  1. X-Real-IP header
//  2. First value of X-Forwarded-For header
//  3. RemoteAddr
//
// NOTE: if this server sits behind a trusted reverse proxy (nginx, Caddy,
// a CDN), set echo.IPExtractor on the Echo instance in server.go to pin
// extraction to the correct header and prevent IP spoofing via forged headers.
// e.g.: e.IPExtractor = echo.ExtractIPFromXFFHeader()
func IPRateLimit(rps float64, burst int) echo.MiddlewareFunc {
	l := newLimiters(rps, burst)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			blocked, retryAfter := retryAfterSeconds(l.get(ip))
			if blocked {
				c.Response().Header().Set("Retry-After", retryAfter)
				return echo.NewHTTPError(http.StatusTooManyRequests, "too many requests")
			}
			return next(c)
		}
	}
}

// UserRateLimit limits requests by authenticated user ID.
// Must be applied after Auth + UserContext middleware so that
// "user_id" is already set on the context.
//
// Falls back to IP-based limiting if user_id is not present,
// so misconfigured routes degrade gracefully instead of panicking.
func UserRateLimit(rps float64, burst int) echo.MiddlewareFunc {
	l := newLimiters(rps, burst)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := userKey(c)
			blocked, retryAfter := retryAfterSeconds(l.get(key))
			if blocked {
				c.Response().Header().Set("Retry-After", retryAfter)
				return echo.NewHTTPError(http.StatusTooManyRequests, "too many requests")
			}
			return next(c)
		}
	}
}

// userKey returns "user:<id>" when user_id is on the context,
// falling back to "ip:<addr>" so the limiter always has a valid key.
func userKey(c echo.Context) string {
	if id, ok := c.Get("user_id").(uint); ok && id != 0 {
		return "user:" + uintToString(id)
	}
	return "ip:" + c.RealIP()
}

// uintToString converts uint to string without importing strconv globally.
func uintToString(n uint) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}