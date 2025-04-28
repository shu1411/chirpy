package main

import (
	"net/http"
)

func handlerReadiness(respWriter http.ResponseWriter, req *http.Request) {
	respWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
	respWriter.WriteHeader(http.StatusOK)
	respWriter.Write([]byte("OK"))
}