# Token Bucket Rate Limiter

A Go package for implementing a token bucket rate limiter.

## Overview

This package provides a simple and efficient way to rate limit incoming requests using the token bucket algorithm. The token bucket algorithm is a widely used algorithm for rate limiting, which allows for a burst of requests followed by a steady rate of requests.

## Features

* Token bucket algorithm for rate limiting
* Configurable rate limit and refill interval
* Support for multiple rate limiters
* Simple and efficient implementation

## Usage

To use this package, you can create a new rate limiter instance and configure it with your desired rate limit and refill interval.

```go
package main

import (
    "github.com/SaidakbarPardaboyev/Token-Bucket-Rate-Limiter"
    "time"
)

func main() {
    // Create a new rate limiter instance
    rateLimiter := New()

    // Configure the rate limiter with a rate limit of 1000 requests per 2 seconds
    rateLimiter.SetConfig(RateLimiterConfig{
        RATE_LIMIT:      1000,
        REFILL_INTERVAL: 2 * time.Second,
    })

    // Run the rate limiter
    rateLimiter.Run()

    // Use the rate limiter to rate limit incoming requests
    http.HandleFunc("/test", rateLimiter.RateLimitMiddleware(http.HandlerFunc(  TestEndpoint)))
}