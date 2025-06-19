package controllers

import (
	"awesomeProject/internal/domain/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

type Request struct {
	User struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Age       int    `json:"age"`
	} `json:"user"`
}

type Response struct {
	Message string   `json:"message,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

func (h *Handler) SaveUserHandler(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.SaveUserHandler"
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var req Request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Log.Error("ошибка декодирования кода",
			zap.String("op", op),
			zap.Error(err),
		)
		responseWithError(w, http.StatusBadRequest, "неверный JSON")
		return //ОБЯЗАТЕЛЬНО
	}

	if err := Validation(req); err != nil {
		h.Log.Error("валидация не пройдена",
			zap.String("op", op),
			zap.Error(err),
		)
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user := &models.User{
		FirstName: req.User.FirstName,
		LastName:  req.User.LastName,
		Age:       req.User.Age,
	}

	err := h.UserService.SaveUser(ctx, user)
	if err != nil {
		h.Log.Error("юзер не сохранен",
			zap.String("op", op),
			zap.Error(err),
		)
		responseWithError(w, http.StatusInternalServerError, "ошибка при сохранении")
		return
	}

	responseWithJson(w, http.StatusCreated, Response{
		Message: fmt.Sprintf("пользователь %s успешно сохранен", user.FirstName),
	})

}

func Validation(req Request) error {
	var validationErrors []string

	firstName := strings.TrimSpace(req.User.FirstName)
	if firstName == "" {
		validationErrors = append(validationErrors, "укажите имя")
	}

	lastName := strings.TrimSpace(req.User.LastName)
	if lastName == "" {
		validationErrors = append(validationErrors, "укажите фаилию")
	}

	if req.User.Age <= 0 {
		validationErrors = append(validationErrors, "укажите ваш возраст")
	} else if req.User.Age > 120 {
		validationErrors = append(validationErrors, "неверный возраст")
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf(strings.Join(validationErrors, "\n"))
	}

	return nil
}
