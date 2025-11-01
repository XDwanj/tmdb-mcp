package middleware

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// authMiddleware validates Bearer Token
func AuthMiddleware(token string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer "+token {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"unauthorized"}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// AuthMiddlewareWithLogger 采用常量时间比较并输出结构化安全日志
// 失败时返回统一 JSON 错误体，避免泄露敏感信息
func AuthMiddlewareWithLogger(logger *zap.Logger, token string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		got := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
		want := token

		// 将长度不匹配情况也走 ConstantTimeCompare，避免时序侧信道
		// 若长度不同，拼接到相同长度再比较（比较结果必然不等）
		paddedGot := got
		paddedWant := want
		if len(paddedGot) < len(paddedWant) {
			paddedGot = paddedGot + strings.Repeat(" ", len(paddedWant)-len(paddedGot))
		} else if len(paddedWant) < len(paddedGot) {
			paddedWant = paddedWant + strings.Repeat(" ", len(paddedGot)-len(paddedWant))
		}

		if subtle.ConstantTimeCompare([]byte(paddedGot), []byte(paddedWant)) != 1 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error":   "unauthorized",
				"message": "missing or invalid bearer token",
			})
			if logger != nil {
				logger.Warn("auth failed",
					zap.String("event", "auth_failed"),
					zap.String("addr", r.RemoteAddr),
					zap.String("path", r.URL.Path),
					zap.String("method", r.Method),
				)
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}
