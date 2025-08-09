package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

const port = "8080"
const root = "/"

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	mux := http.NewServeMux()

	server := &http.Server{
		Handler: mux,
		Addr: ":" + port,
	}

	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
	}

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerResetCount)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	log.Printf("Serving files from %s on port: %s\n", root, port)
	log.Fatal(server.ListenAndServe())
}