package main

import (
	"encoding/json"
	"net/http"
)

func handlerValidate(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return 
	}

	const maxBodyLength = 140
	if len(params.Body) > maxBodyLength {
		respondWithError(w, http.StatusBadRequest, "Message is too long", nil)
		return 
	}

	if containsProfaneWords(params.Body) {
		params.Body = cleanBody(params.Body)
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: params.Body,
	})
}