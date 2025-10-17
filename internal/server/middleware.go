package server

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"go.uber.org/zap"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap response writer to capture status code
			wrapped := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK, // default status code
			}

			// Call the next handler
			next.ServeHTTP(wrapped, r)

			// Log the request
			elapsed := time.Since(start)
			logger.Info("HTTP request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
				zap.Int("status", wrapped.statusCode),
				zap.Duration("duration", elapsed),
			)
		})
	}
}

// RecoveryMiddleware recovers from panics and logs them
func RecoveryMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Log the panic with stack trace
					logger.Error("Panic recovered",
						zap.String("method", r.Method),
						zap.String("path", r.URL.Path),
						zap.String("remote_addr", r.RemoteAddr),
						zap.Any("panic", err),
						zap.String("stack", string(debug.Stack())),
					)

					// Return 500 Internal Server Error
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, `{"error": "Internal server error"}`)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// AuthMiddleware validates Bearer Token authentication
func AuthMiddleware(expectedToken string, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Extract Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Warn("Missing Authorization header",
					zap.String("remote_addr", r.RemoteAddr),
					zap.String("path", r.URL.Path),
				)
				sendUnauthorized(w, "Missing Authorization header")
				return
			}

			// 2. Validate Bearer format
			const bearerPrefix = "Bearer "
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				logger.Warn("Invalid Authorization header format",
					zap.String("remote_addr", r.RemoteAddr),
					zap.String("path", r.URL.Path),
				)
				sendUnauthorized(w, "Invalid Authorization header format")
				return
			}

			// 3. Extract token
			providedToken := strings.TrimPrefix(authHeader, bearerPrefix)

			// 4. Constant-time comparison (防止时序攻击)
			if !compareTokens(providedToken, expectedToken) {
				logger.Warn("Invalid token",
					zap.String("remote_addr", r.RemoteAddr),
					zap.String("path", r.URL.Path),
				)
				sendUnauthorized(w, "Invalid token")
				return
			}

			// 5. Authentication successful
			logger.Debug("Authentication successful",
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("path", r.URL.Path),
			)

			next.ServeHTTP(w, r)
		})
	}
}

// sendUnauthorized sends a 401 Unauthorized JSON error response
func sendUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	errorResponse := map[string]string{
		"error":   "unauthorized",
		"message": message,
	}

	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		// Log encoding errors (rare but possible)
		// Note: Cannot log here without logger, but error is unlikely
		_ = err
	}
}

// compareTokens performs constant-time comparison to prevent timing attacks
func compareTokens(provided, expected string) bool {
	providedBytes := []byte(provided)
	expectedBytes := []byte(expected)

	// ConstantTimeCompare returns 1 if equal, 0 otherwise
	// Requires equal length inputs, otherwise returns 0
	return subtle.ConstantTimeCompare(providedBytes, expectedBytes) == 1
}
