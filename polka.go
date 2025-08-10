package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/shu1411/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUpdatePolkaByID(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		}
	}

	headerAPIKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't get api key", err)
		return
	}

	if headerAPIKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "invalid api key", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = cfg.db.UpgradeUserPolkaByID(r.Context(), params.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "couldn't find user", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "couldn't update user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
