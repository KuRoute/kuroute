package middleware

import (
	"context"
	"net/http"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/package/jwt"
	"github.com/KuRoute/kuroute/backend/package/response"
	"github.com/google/uuid"
)

type contextKey string

const (
	UserContextKey contextKey = "user"
)

type AuthUser struct {
	UserID uuid.UUID
	Role   domain.UserRole
	HubID  uuid.UUID
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
            next.ServeHTTP(w, r)
            return
        }

		tokenString := jwt.ExtractToken(r)
		if tokenString == "" {
			response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing authorization token")
			return
		}

		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or expired token")
			return
		}

		authUser := &AuthUser{
			UserID: claims.UserID,
			Role:   claims.Role,
			HubID:  claims.HubID,
		}

		ctx := context.WithValue(r.Context(), UserContextKey, authUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RoleMiddleware(roles ...domain.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authUser, ok := r.Context().Value(UserContextKey).(*AuthUser)
			if !ok || authUser == nil {
				response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
				return
			}

			hasRole := false
			for _, role := range roles {
				if authUser.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				response.Fail(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GetAuthUser(r *http.Request) *AuthUser {
	authUser, ok := r.Context().Value(UserContextKey).(*AuthUser)
	if !ok {
		return nil
	}
	return authUser
}