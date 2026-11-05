package hamoc

import (
	"log"
	"net/http"
	"time"
)

// Logger logs method, path, status, and duration.
func Logger() Middleware {
	return func(next Handler) Handler {
		return func(r *Request) Response {
			start := time.Now()
			resp := next(r)
			status := resp.Status
			if status == 0 {
				status = http.StatusOK
			}
			log.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, status, time.Since(start))
			return resp
		}
	}
}

// Recoverer recovers from panics and returns a 500.
func Recoverer() Middleware {
	return func(next Handler) Handler {
		return func(r *Request) Response {
			defer func() {
				if rec := recover(); rec != nil {
					log.Printf("panic: %v", rec)
				}
			}()
			return next(r)
		}
	}
}
