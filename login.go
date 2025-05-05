package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/shu1411/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password 			string 	`json:"password"`
		Email	 			string 	`json:"email"`
		ExpiresInSeconds	int		`json:"expires_in_seconds"`
	}

	type response struct {
		User
		Token		 string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	if auth.CheckPasswordHash(user.HashedPassword, params.Password) != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	expirationTime := time.Hour
	if params.ExpiresInSeconds > 0 && params.ExpiresInSeconds < 3600 {
		expirationTime = time.Duration(params.ExpiresInSeconds) * time.Second
	}

	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.Secret,
		expirationTime,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create access JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID: 		user.ID,
			CreatedAt: 	user.CreatedAt,
			UpdatedAt: 	user.UpdatedAt,
			Email: 		user.Email,

		},
		Token: accessToken,
	})
}