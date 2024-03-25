package handlers

import (
	"log"
	"net/http"
)

func SendResponse(rw http.ResponseWriter, statusCode int, message string) {
	rw.WriteHeader(statusCode)
	_, err := rw.Write([]byte(message))
	if err != nil {
		log.Printf("error during writing response %v", err)
	}
}
