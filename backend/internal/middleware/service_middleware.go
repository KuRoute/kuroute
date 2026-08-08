package middleware

import (
	"net/http"
	"context"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/KuRoute/kuroute/backend/package/jwt"
	"github.com/KuRoute/kuroute/backend/package/response"
)

type serviceContext string

const (
	ServiceContextKey serviceContext = "service"
)

type AuthService struct {
	Service	domain.ServiceName
}

func ServiceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
            next.ServeHTTP(w, r)
            return
        }

		tokenString := jwt.ExtractTokenService(r)
		if tokenString == "" {
			response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing authorization token")
			return
		}

		claims, err := jwt.ValidateTokenService(tokenString)
		if err != nil {
			response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or expired token")
			return
		}

		authService := &AuthService{
			Service: claims.Service,
		}

		ctx := context.WithValue(r.Context(), ServiceContextKey, authService)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NameServiceMiddleware(serviceName ...domain.ServiceName) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authService, ok := r.Context().Value(ServiceContextKey).(*AuthService)
			if !ok || authService == nil {
				response.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "Service name required")
				return
			}

			hasName := false
			for _, service := range serviceName {
				if authService.Service == service {
					hasName = true
					break
				}
			}

			if !hasName {
				response.Fail(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GetAuthService(r *http.Request) *AuthService {
	authService, ok := r.Context().Value(ServiceContextKey).(*AuthService)
	if !ok {
		return nil
	}
	return authService
}