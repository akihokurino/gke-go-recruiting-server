package handler

import (
	"net/http"
	"os"

	"gke-go-recruiting-server/domain"

	"gke-go-recruiting-server/adapter"
)

const (
	authKey      = "Authorization"
	debugAuthKey = "X-Debug-User-Id"
)

func NewCros() adapter.Cros {
	return func(base http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Debug-User-Id, Authorization, X-Auth-Id")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "3600")
			if r.Method == "OPTIONS" {
				w.WriteHeader(200)
				return
			}
			base.ServeHTTP(w, r)
		})
	}
}

func NewAdminAuthenticate(cp adapter.ContextProvider, fc adapter.FirebaseAppFactory) adapter.AdminAuthenticate {
	return func(base http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			if os.Getenv("APP_ENV") == "dev" {
				uid := r.Header.Get(debugAuthKey)
				if uid != "" {
					newContext, _ := cp.WithFirebaseUserID(ctx, domain.FirebaseUserID(uid))
					base.ServeHTTP(w, r.WithContext(newContext))
					return
				}
			}

			client, err := fc(ctx).Auth(ctx)
			if err != nil {
				http.Error(w, "internal server error", 500)
				return
			}

			token := r.Header.Get(authKey)
			if len(token) <= 7 {
				http.Error(w, "unauthorized", 401)
				return
			}

			decoded, err := client.VerifyIDToken(ctx, token[7:])
			if err != nil {
				http.Error(w, "unauthorized", 401)
				return
			}

			newContext, _ := cp.WithFirebaseUserID(ctx, domain.FirebaseUserID(decoded.UID))
			base.ServeHTTP(w, r.WithContext(newContext))
		})
	}
}
