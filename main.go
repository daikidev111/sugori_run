package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type response struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		res, err := json.Marshal(&response{
			Message: "Hello World!",
		})
		if err != nil {
			log.Printf("failed to marshal json. %+v\n", err)
			return
		}
		if _, err := w.Write(res); err != nil {
			log.Printf("failed to write response. %+v\n", err)
		}
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("failed to serve. %+v", err)
	}
}
