package postgresql

import (
	"awesomeProject/internal/domain/models"
	"context"
	"log"
)

func (p *PostgreSQL) ListUsersPostgreSQL(
	ctx context.Context,
	minAge, maxAge *int,
	startDate, endDate *int64,
) ([]models.User, error) {
	const op = "ListUsersPostgreSQL"
	query := `
		SELECT id, first_name, last_name, age, recording_date 
		FROM users 
		WHERE 
			($1::int IS NULL OR age >= $1) AND
			($2::int IS NULL OR age <= $2) AND
			($3::bigint IS NULL OR recording_date >= $3) AND
			($4::bigint IS NULL OR recording_date <= $4)
	`

	rows, err := p.Db.QueryContext(ctx, query, minAge, maxAge, startDate, endDate)
	if err != nil {
		log.Printf("error with querying users in %s: %v", op, err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Age,
			&user.RecordingDate,
		)
		if err != nil {
			log.Printf("error with scanning user in %s: %v", op, err)
			continue //иначе при ошибке все равно продолжит работать
		}
		users = append(users, user)
	}

	return users, nil
}
