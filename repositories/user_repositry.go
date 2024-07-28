package repositories

import (
	"Nookhub/models"
	"database/sql"
	"fmt"
)

type UserRepository interface {
	Create(user models.User) (int, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user models.User) (int, error) {
	// Create a variable to store the OUT parameter
	var userId int

	// Prepare the call to the stored procedure
	stmt, err := r.db.Prepare("CALL sp_InsertUser(?, ?, ?, @userId)")
	if err != nil {
		return 0, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	// Execute the stored procedure
	_, err = stmt.Exec(user.Username, user.Email, user.Password)
	if err != nil {
		return 0, fmt.Errorf("error executing stored procedure: %w", err)
	}

	// Retrieve the OUT parameter value
	row := r.db.QueryRow("SELECT @userId")
	err = row.Scan(&userId)
	if err != nil {
		return 0, fmt.Errorf("error retrieving user ID: %w", err)
	}

	return userId, nil
}
