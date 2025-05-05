package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/shu1411/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits 	atomic.Int32
	db             	*database.Queries
	Platform       	string
	Secret			string
}

const filePathRoot = "."
const port = "8080"

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	jwtSecret := os.Getenv("SECRET")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open database: %v\n", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	serveMux := http.NewServeMux()
	server := &http.Server{
		Handler: serveMux,
		Addr:    ":" + port,
	}

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		Platform:       platform,
		Secret: 		jwtSecret,
	}

	serveMux.HandleFunc("GET /api/healthz", handlerReadiness)

	serveMux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	serveMux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	serveMux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	serveMux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	serveMux.HandleFunc("PUT /api/users", apiCfg.handlerUpdate)

	serveMux.HandleFunc("POST /api/chirps", apiCfg.handlerPostChirp)
	serveMux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpByID)

	serveMux.HandleFunc("GET /admin/metrics", apiCfg.handlerRequestCount)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))

	log.Fatal(server.ListenAndServe())
}
