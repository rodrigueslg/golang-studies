package controller

import (
	"encoding/json"
	"gochallenges/pkg"
	"net/http"
)

const jsonContentType = "application/json"
const authHeader = "Authorization"

func authorizeRequest(w http.ResponseWriter, r *http.Request) bool {
	authHeader := r.Header.Get(authHeader)
	return authHeader == pkg.GetBearerToken()
}

func parseJsonBody(w http.ResponseWriter, r *http.Request, t interface{}) error {
	return json.NewDecoder(r.Body).Decode(&t)
}

func writeUnauthorizedResponse(w http.ResponseWriter, err error) {
	w.Header().Set("content-type", jsonContentType)
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(err.Error()))
}

func writeInternalErrorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("content-type", jsonContentType)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func writeBadRequestResponse(w http.ResponseWriter, err error) {
	w.Header().Set("content-type", jsonContentType)
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
}

func writeNotFoundResponse(w http.ResponseWriter) {
	w.Header().Set("content-type", jsonContentType)
	w.WriteHeader(http.StatusNotFound)
}

func writeOkResponse(w http.ResponseWriter, content any) {
	w.Header().Set("content-type", jsonContentType)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(content)
}

func writeCreatedResponse(w http.ResponseWriter, content any) {
	w.Header().Set("content-type", jsonContentType)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(content)
}
