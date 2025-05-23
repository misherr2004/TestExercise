package userservice

import (
	"awesomeProject/internal/domain/models"
	"go.uber.org/zap"
)

// интерфейс, который определяет контракт для работы с БД
type UserDB interface {
	SaveUser(user *models.User) error
	GetUserPostgreSQL(firstName, lastName string, age int) ([]models.User, error)
}

type UserService struct {
	db  UserDB
	Log *zap.Logger
}

func NewUserService(db UserDB, log *zap.Logger) *UserService {
	return &UserService{
		db:  db,
		Log: log,
	}
}
