package controllers

import (
	"awesomeProject/internal/userservice"
	"encoding/json"
	"log"
	"net/http"

	"go.uber.org/zap"
)

type Handler struct {
	UserService userservice.IUserService
	Log         *zap.Logger
}

func NewHandler(userService userservice.IUserService, log *zap.Logger) *Handler {
	return &Handler{UserService: userService, Log: log}
}

func responseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("ошибка кодирования JSON: %v", err)
		http.Error(w, "ошибка формирования ответа", http.StatusInternalServerError)
	}
}

func responseWithError(w http.ResponseWriter, code int, message string) {
	responseWithJson(w, code, Response{Errors: []string{message}})
}
