package userservice

import (
	"awesomeProject/internal/domain/models"
	"context"
	"fmt"
	"log"
)

func (u *UserService) SaveUser(ctx context.Context, user *models.User) error {
	const op = "SaveUser"

	if err := u.Db.SaveUser(ctx, user); err != nil {
		log.Printf("error with saving user in %s: %v", op, err)
		return fmt.Errorf("метод: %s, ошибка: %w", op, err)
	}
	log.Println("User created")

	return nil
}

func (u *UserService) GetUser(ctx context.Context, firstName, lastName string, age int) ([]models.User, error) {
	return u.Db.GetUserPostgreSQL(ctx, firstName, lastName, age)
}

func (u *UserService) ListUsers(
	ctx context.Context,
	minAge, maxAge *int,
	startDate, endDate *int64,
) ([]models.User, error) {
	return u.Db.ListUsersPostgreSQL(ctx, minAge, maxAge, startDate, endDate)
}

//func (u *UserService) UserDelete(ctx context.Context, user *models.User) error {
//	return u.Db.UpdateUser(ctx, user)
//}
