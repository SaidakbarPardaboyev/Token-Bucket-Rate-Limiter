package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type RateLimiter interface {
	Run()
	Config() *rateLimiter
	SetConfig(RateLimiterConfig)
	RefillBucket()
	GetBucketStatus(w http.ResponseWriter, r *http.Request)
	RateLimitMiddleware(next http.Handler) http.Handler
}

type rateLimiter struct {
	RATE_LIMIT      int64
	REFILL_INTERVAL time.Duration
	tokenBucket     []int64
	mx              sync.Mutex
}

type RateLimiterConfig struct {
	RATE_LIMIT      int64
	REFILL_INTERVAL time.Duration
}

type BucketStatus struct {
	BucketLimit       int64
	CurrentBucketSize int64
	Bucket            []int64
}

func New() RateLimiter {
	return &rateLimiter{}
}

func (r *rateLimiter) Config() *rateLimiter {
	return r
}

func (r *rateLimiter) SetConfig(rateLimiter RateLimiterConfig) {
	r.RATE_LIMIT = rateLimiter.RATE_LIMIT
	r.REFILL_INTERVAL = rateLimiter.REFILL_INTERVAL
}

func (r *rateLimiter) RefillBucket() {
	r.mx.Lock()
	defer r.mx.Unlock()

	if int64(len(r.tokenBucket)) < r.RATE_LIMIT {
		r.tokenBucket = append(r.tokenBucket, time.Now().UnixNano())
	}
}

func (r *rateLimiter) GetBucketStatus(w http.ResponseWriter, request *http.Request) {
	r.mx.Lock()
	defer r.mx.Unlock()

	response := BucketStatus{
		BucketLimit:       r.RATE_LIMIT,
		CurrentBucketSize: int64(len(r.tokenBucket)),
		Bucket:            []int64{},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (r *rateLimiter) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		r.mx.Lock()
		defer r.mx.Unlock()

		if len(r.tokenBucket) > 0 {
			r.tokenBucket = r.tokenBucket[1:]
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", len(r.tokenBucket)))
			next.ServeHTTP(w, request)
		} else {
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("Retry-After", fmt.Sprintf("%f second", r.REFILL_INTERVAL.Seconds()))
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"message": "Too many requests",
			})
		}
	})
}

func (r *rateLimiter) Run() {
	ticker := time.NewTicker(r.REFILL_INTERVAL)

	go func() {
		for range ticker.C {
			r.RefillBucket()
		}
		defer ticker.Stop()
	}()
}

// Sample endpoint for testing rate limiting
func TestEndpoint(w http.ResponseWriter, r *http.Request) {
	rockPaperScissors := []string{"rock 🪨", "paper 📃", "scissors ✂️"}
	randomChoice := rockPaperScissors[time.Now().UnixNano()%3]

	response := map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("You got %s", randomChoice),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// func main() {
// 	// Initialize the rate limiter
// 	rateLimiter := New()
// 	rateLimiter.SetConfig(RateLimiterConfig{
// 		RATE_LIMIT:      1000,
// 		REFILL_INTERVAL: 2 * time.Second,
// 	})
// 	rateLimiter.Run()

// 	// Setup HTTP server and routes
// 	http.HandleFunc("/bucket", rateLimiter.GetBucketStatus)
// 	http.Handle("/test", rateLimiter.RateLimitMiddleware(http.HandlerFunc(testEndpoint)))

// 	// Start the server
// 	fmt.Println("Server running on port 5000")
// 	if err := http.ListenAndServe(":5000", nil); err != nil {
// 		fmt.Println("Error starting server:", err)
// 	}
// }
