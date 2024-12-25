package storage

import (
	"context"
	"database/sql"

	database_models "github.com/KonovalovIly/anki_pdf/database/model"
)

type UserStorage struct {
	db *sql.DB
}

func (s *UserStorage) GetUser(ctx context.Context, userID int64) (*database_models.UserDto, *database_models.DatabaseError) {
	query := `SELECT
        id,
        login,
        email,
        password
        FROM users WHERE id = $1
    `

	row := s.db.QueryRowContext(ctx, query, userID)
	var user database_models.UserDto
	err := row.Scan(&user.ID, &user.Login, &user.Email, &user.Password)

	if err != nil {
		return nil, database_models.ProcessErrorFromDatabase(err, "GetUser:Scan")
	}

	return &user, nil
}

func (s *UserStorage) SaveUser(ctx context.Context, userDto *database_models.UserDto) *database_models.DatabaseError {
	query := `INSERT INTO users
        (login, email, password)
        VALUES ($1, $2, $3) RETURNING id
    `

	ctx, cancel := context.WithTimeout(ctx, QueryRowTimeout)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		userDto.Login,
		userDto.Email,
		userDto.Password,
	).Scan(&userDto.ID)

	if err != nil {
		return database_models.ProcessErrorFromDatabase(err, "SaveUser:QueryRowContext")
	}

	return nil
}
