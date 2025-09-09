package middleware

import (
	"fmt"
	"net/http"
	"runtime/trace"
	"time"
)

// Tracing middleware tracks request details and execution time using Go's runtime trace.
func Tracing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Start a new trace task for this request.
		// defer ensures task.End() runs when this function exits, even on panic or early return.
		ctx, task := trace.NewTask(r.Context(), fmt.Sprintf("%s %s", r.Method, r.URL.Path))
		defer task.End()

		// Log request details
		trace.Log(ctx, "method", r.Method)
		trace.Log(ctx, "path", r.URL.Path)

		// Measure request duration
		start := time.Now()
		// Pass the tracing context to downstream handlers
		next.ServeHTTP(w, r.WithContext(ctx))
		trace.Log(ctx, "duration_ms", fmt.Sprintf("%.2f", time.Since(start).Seconds()*1000))
	})
}
