package auth

import (
	"context"
	"net/http"
)

const DefaultCookieName = "access_token"

type contextKey string

const claimsContextKey contextKey = "auth_claims"

func RequireAuth(tokens *TokenService, cookieName string) func(http.Handler) http.Handler {
	if cookieName == "" {
		cookieName = DefaultCookieName
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(cookieName)
			if err != nil {
				// http.Error(w, "Unauthorized", http.StatusUnauthorized)
				http.Redirect(w, r, "/auth/refresh",  http.StatusTemporaryRedirect)
				return
			}

			claims, err := tokens.Verify(cookie.Value)
			if err != nil {
				//http.Error(w, "Unauthorized", http.StatusUnauthorized)
				http.Redirect(w, r, "/auth/refresh",  http.StatusTemporaryRedirect)
				return
			}

			ctx := context.WithValue(r.Context(), claimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFromContext(r.Context())
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if _, ok := allowed[claims.Role]; !ok {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey).(*Claims)
	return claims, ok
}

