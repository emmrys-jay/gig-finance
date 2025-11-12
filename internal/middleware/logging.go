package middleware

import (
	"context"
	"net/http"
	"time"
)

type contextKey string

const StartTimeKey contextKey = "start_time"

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx := context.WithValue(r.Context(), StartTimeKey, start)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetStartTime(r *http.Request) time.Time {
	if start, ok := r.Context().Value(StartTimeKey).(time.Time); ok {
		return start
	}
	// Fallback to current time if not found (shouldn't happen with middleware)
	return time.Now()
}
