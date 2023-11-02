package main

import "net/http"

type apiError struct {
	Error string `json:"error"`
}

func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, apiError{Error: "Permission Denied"})
}

func invalidMethod(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, apiError{Error: "Invalid Method Used"})
}
