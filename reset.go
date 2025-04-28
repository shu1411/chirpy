package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerResetRequestCount(respWriter http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	respWriter.WriteHeader(http.StatusOK)
	respWriter.Write([]byte("Hits reset to 0"))
}