package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/Namularbre/knowledgeKeeperApi/internal/auth/domain"
)

type ctxKey struct{}

var userIDKey = ctxKey{}

// RequireBearer wraps an http.Handler so it only executes when a valid
// access token is present in the Authorization header. The authenticated
// user ID is injected into the request context.
func RequireBearer(tokens domain.TokenIssuer, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			writeError(w, http.StatusUnauthorized, "missing_bearer_token")
			return
		}
		raw := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		userID, err := tokens.ParseAccessToken(raw)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid_access_token")
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserIDFrom extracts the authenticated user ID injected by RequireBearer.
// The second return value indicates whether the context carried an ID.
func UserIDFrom(ctx context.Context) (int64, bool) {
	v, ok := ctx.Value(userIDKey).(int64)
	return v, ok
}
