package main

import (
	"fmt"
	"net/http"
)

// Atomically inrement request count and return a handler to the next request
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

// Display the current count of requests made
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(
		fmt.Sprintf(`
<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
</html>
`, cfg.fileserverHits.Load()),
	))
}

// Reset request count to 0
func (cfg *apiConfig) handlerResetCount(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "only allowed in local development", nil)
		return 
	}
	cfg.fileserverHits.Store(0)
	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't delete users", err)
		return 
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request count was set to 0"))
}
