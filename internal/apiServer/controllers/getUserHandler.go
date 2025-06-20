package controllers

import (
	"awesomeProject/internal/domain/models"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type GetUserResponse struct {
	Users  []models.User `json:"users"`
	Errors []string      `json:"errors,omitempty"`
}

type GetUserParams struct {
	firstName string
	lastName  string
	age       *int
}

func (h *Handler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.GetUserHandler"
	ctx, cansel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cansel()

	query := r.URL.Query()
	var params GetUserParams

	params.firstName = strings.TrimSpace(query.Get("first_name"))
	params.lastName = strings.TrimSpace(query.Get("last_name"))
	if ageStr := query.Get("age"); ageStr != "" {
		age, err := strconv.Atoi(ageStr)
		if err != nil {
			h.Log.Error("ошибка преобразования возраста",
				zap.String("op", op),
				zap.String("value", ageStr),
				zap.Error(err))
			responseWithError(w, http.StatusBadRequest, "некорректный возраст")
			return
		}
		params.age = &age
	}

	if err := ValidateGetUserParams(params); err != nil {
		h.Log.Error("валидация не пройдена",
			zap.String("op", op),
			zap.Error(err))
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	users, err := h.UserService.GetUser(ctx, params.firstName, params.lastName, *params.age)
	if err != nil {
		h.Log.Error("ошибка при получении пользователей",
			zap.String("op", op),
			zap.Error(err))
		responseWithError(w, http.StatusInternalServerError, "ошибка сервера")
		return
	}

	responseWithJson(w, http.StatusOK, map[string]interface{}{
		"answer": users,
	})
}

func ValidateGetUserParams(params GetUserParams) error {
	var validationErrors []string

	if params.firstName == "" {
		validationErrors = append(validationErrors, "укажите имя")
	}

	if params.lastName == "" {
		validationErrors = append(validationErrors, "фамилию напиши - не ленись")
	}

	if params.age != nil {
		if *params.age < 0 {
			validationErrors = append(validationErrors, "возраст не может быть отрицательным")
		} else if *params.age > 120 {
			validationErrors = append(validationErrors, "возраст не может быть больше 120")
		}
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf(strings.Join(validationErrors, "\n"))
	}

	return nil
}
