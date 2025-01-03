package shared

import (
	"encoding/json"
	"net/http"
	"os"
)

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin",
			getenv("ACCESS_CONTROL_ALLOW_ORIGIN", "*"))
		w.Header().Set("Access-Control-Allow-Headers",
			getenv("ACCESS_CONTROL_ALLOW_HEADERS", "X-Session-id,Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token"))
		w.Header().Set("Access-Control-Allow-Methods",
			getenv("ACCESS_CONTROL_ALLOW_METHODS", "DELETE, GET, HEAD, OPTIONS, PATCH, POST, PUT"))

		next.ServeHTTP(w, r)
	})
}

func ValidateIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.PathValue("id")) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ID_CANNOT_BE_EMPTY)
			return
		}

		if !IsUUID(r.PathValue("id")) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ID_NOT_IN_UUID_FORMAT)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
