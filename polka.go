package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerMembership(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data struct{
			UserID string `json:"user_id"`
		}
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't get user id", err)
		return
	}

	_, err = cfg.db.UpgradeMembership(req.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "couldn't find user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}