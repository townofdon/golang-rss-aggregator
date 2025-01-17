package util

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithJSON(res http.ResponseWriter, statusCode int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON:\n%v\n", payload)
		res.WriteHeader(500)
		return
	}

	res.Header().Add("Content-Type", "Application")
	res.WriteHeader(statusCode)

	_, err = res.Write(data)
	if err != nil {
		log.Printf("Failed to write bytes to response")
		res.WriteHeader(500)
		return
	}
}