package main

import (
	"fmt"
	"net/http"
)

// Atomically inrement request count and return a handler to the next request
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

// Display the current count of requests made
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(
		fmt.Sprintf(`
<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
</html>
`, cfg.fileserverHits.Load()),
	))
}

// Reset request count to 0
func (cfg *apiConfig) handlerResetCount(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request count was set to 0"))
}

/***************************
*	JSON DECODING		   *
****************************/
// func handler(w http.ResponseWriter, r *http.Request) {
// 	type parameters struct {
// 		Name string `json:"name"`
// 		Age  int    `json:"age"`
// 	}

// 	decoder := json.NewDecoder(r.Body)
// 	params := parameters{}
// 	err := decoder.Decode(&params)
// 	if err != nil {
// 		log.Printf("Error decoding parameters: %s", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	// params is populated with data
// }

/*********************
*	JSON ENCODING	 *
**********************/
// func handler(w http.ResponseWriter, r *http.Request) {
// 	type returnVals struct {
// 		CreatedAt time.Time `json:"created_at"`
// 		ID int `json:"id"`
// 	}
// 	respBody := returnVals{
// 		CreatedAt: time.Now(),
// 		ID: 123,
// 	}
// 	dat, err := json.Marshal(respBody)
// 	if err != nil {
// 		log.Printf("Error marshaling JSON: %s", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return 
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(dat)
// }
