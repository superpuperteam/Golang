package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type userIDResponse struct {
	UserID int `json:"user_id"`
}

type createUserRequest struct {
	Name string `json:"name"`
}

type createUserResponse struct {
	Created string `json:"created"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		getUserHandler(w, r)
		return
	}
	if r.Method == http.MethodPost {
		postUserHandler(w, r)
		return
	}
	w.Header().Set("Allow", http.MethodGet+", "+http.MethodPost)
	WriteJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || idStr == "" {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid id"})
		return
	}
	WriteJSON(w, http.StatusOK, userIDResponse{UserID: id})
}

func postUserHandler(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid name"})
		return
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid name"})
		return
	}
	WriteJSON(w, http.StatusCreated, createUserResponse{Created: name})
}
