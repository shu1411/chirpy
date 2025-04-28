package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerRequestCount(respWriter http.ResponseWriter, req *http.Request) {
	respWriter.Header().Add("Content-Type", "text/html")
	respWriter.WriteHeader(http.StatusOK)
	respWriter.Write([]byte(fmt.Sprintf(`
<html>

<body>
	<h1>Welcome, Chirpy Admin<h1>
	<p> Chirpy has been visited %d times!</p>
</body>
</html>
	`, cfg.fileserverHits.Load())))
}