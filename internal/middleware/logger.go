package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	reset  = "\033[0m"
	gray   = "\033[90m"
	cyan   = "\033[36m"
	green  = "\033[32m"
	yellow = "\033[33m"
	red    = "\033[31m"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(body []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}
	return rw.ResponseWriter.Write(body)
}

func statusColor(status int) string {
	switch {
	case status >= 500:
		return red
	case status >= 400:
		return yellow
	case status >= 300:
		return cyan
	default:
		return green
	}
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w}
		next.ServeHTTP(rw, r)

		color := statusColor(rw.statusCode)
		fmt.Printf(
			"%s[Go]%s %s%d%s - %s%s%s  %sLOG%s %s[HTTP]%s %s%s%s %s→%s %s%d%s | %s%dms%s | %s%s%s | %s%s%s\n",
			cyan, reset,
			gray, os.Getpid(), reset,
			gray, time.Now().Format("01/02/2006, 3:04:05 PM"), reset,
			green, reset,
			cyan, reset,
			yellow, r.Method+" "+r.URL.RequestURI(), reset,
			gray, reset,
			color, rw.statusCode, reset,
			gray, time.Since(start).Milliseconds(), reset,
			gray, r.RemoteAddr, reset,
			yellow, r.UserAgent(), reset,
		)
	})
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				fmt.Printf("PANIC: %v\n", rec)
				http.Error(w, `{"success":false,"message":"Internal server error"}`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
