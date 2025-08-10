package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/shu1411/chirpy/internal/database"

	// import only for side effects
	_ "github.com/lib/pq"
)

const port = "8080"
const root = "."

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
	polkaKey       string
}

func main() {
	dbURL := getDatabaseURL(".env")
	dbQueries := setupDBQueries(dbURL)

	server := setupServer(dbQueries)

	log.Printf("Serving files from %s on port: %s\n", root, port)
	log.Fatal(server.ListenAndServe())
}

func setupServer(dbQueries *database.Queries) *http.Server {
	mux := http.NewServeMux()

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       os.Getenv("PLATFORM"),
		jwtSecret:      os.Getenv("SECRET"),
		polkaKey:       os.Getenv("POLKA_KEY"),
	}

	// register handler functions
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(root)))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)

	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpByID)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirpByID)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerUpdatePolkaByID)

	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)

	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)

	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerResetCount)

	return server
}

func getDatabaseURL(filename string) string {
	err := godotenv.Load(filename)
	if err != nil {
		log.Fatalf("couldn't load .env: %s", err)
	}

	return os.Getenv("DB_URL")
}

func setupDBQueries(dbURL string) *database.Queries {
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("couldn't open database: %s", err)
	}
	return database.New(db)
}
