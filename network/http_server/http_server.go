package http_server

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	serverPort = ":9999"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// convert headers to json
		jsonData, err := json.Marshal(r.Header)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// set content-type
		w.Header().Set("Content-Type", "application/json")

		// send back the response
		w.Write(jsonData)
	})

	server := &http.Server{
		Addr:    serverPort,
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
