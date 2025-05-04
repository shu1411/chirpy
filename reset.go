package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerReset(respWriter http.ResponseWriter, req *http.Request) {
	if cfg.Platform != "dev" {
		respWriter.WriteHeader(http.StatusForbidden)
		respWriter.Write([]byte("Only allowed in a dev environment"))
		return
	}
	cfg.fileserverHits.Store(0)
	cfg.db.Reset(req.Context())
	respWriter.WriteHeader(http.StatusOK)
	respWriter.Write([]byte("Hits reset to 0 and database reset to intial state"))
}
