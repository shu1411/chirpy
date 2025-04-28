package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

const filePathRoot = "."
const port = "8080"

func main() {
	serveMux := http.NewServeMux()
	server := &http.Server{
		Handler: serveMux,
		Addr: ":" + port,
	}

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	serveMux.HandleFunc("GET /api/healthz", handlerReadiness)
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.handlerRequestCount)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.handlerResetRequestCount)
	serveMux.HandleFunc("POST /api/validate_chirp", handlerValidate)
	
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))

	log.Fatal(server.ListenAndServe())
}