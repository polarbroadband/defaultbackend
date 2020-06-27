package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func init() {
	// config package level default logger
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)
}

func main() {
	// config http server
	r := mux.NewRouter()
	rru := r.PathPrefix("/api/test").Subrouter()
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{`*`}),
		handlers.AllowedHeaders([]string{"content-type"}),
		handlers.AllowCredentials(),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "HEAD"}),
	)

	rru.HandleFunc("", resp).Methods("GET", "POST")

	log.Fatal(http.ListenAndServe(":10009", cors(rru)))
}

func resp(w http.ResponseWriter, r *http.Request) {
	// assembly response data
	res := map[string]interface{}{"RequestHeader": r.Header}
	if r.Body != nil {
		var query interface{}
		json.NewDecoder(r.Body).Decode(&query)
		res["RequestBody"] = query
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.WithError(err).Error("encode error or fail to send response")
		w.WriteHeader(500)
	}
}
