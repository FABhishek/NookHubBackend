package repositories

import (
	"Nookhub/models"
	"database/sql"
	"fmt"
)

type SignupRepository interface {
	CreateUser(user models.RegisterUser) (int, error)
	LoginUser(user models.LoginUser) (string, string, int, error)
	IsAvailable(input string, inputType string) (bool, error)
}

type signupRepository struct {
	db *sql.DB
}

func NewSignupRepository(db *sql.DB) *signupRepository {
	return &signupRepository{db: db}
}

func (r *signupRepository) CreateUser(user models.RegisterUser) (int, error) {
	// Create a variable to store the OUT parameter
	var userId int

	// Prepare the call to the stored procedure
	stmt, err := r.db.Prepare("SELECT func_InsertUser($1, $2, $3)")
	if err != nil {
		return 0, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	// Retrieve the OUT parameter value
	err = stmt.QueryRow(user.Username, user.Email, user.Password).Scan(&userId)
	if err != nil {
		return 0, fmt.Errorf("error executing function: %w", err)
	}

	return userId, nil
}

func (r *signupRepository) LoginUser(user models.LoginUser) (string, string, int, error) {

	var passwordHash string
	var username string
	userId := 0
	stmt, err := r.db.Prepare("Select * From func_getUserLoginData($1)")
	if err != nil {
		return passwordHash, username, userId, fmt.Errorf("error executing function: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(user.Email).Scan(&passwordHash, &username, &userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return passwordHash, username, userId, err
		}
		return passwordHash, username, userId, fmt.Errorf("error executing function: %w", err)
	}
	return passwordHash, username, userId, nil
}

// EmailChecker implements UserRepository.
func (r *signupRepository) IsAvailable(input string, inputType string) (bool, error) {
	var isValid bool

	var stmt *sql.Stmt
	var err error
	if inputType == "email" {
		stmt, err = r.db.Prepare("SELECT func_CheckEmail($1)")
	} else {
		stmt, err = r.db.Prepare("SELECT func_CheckUsername($1)")
	}

	if err != nil {
		return false, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(input).Scan(&isValid)
	if err != nil {
		return false, fmt.Errorf("error executing function: %w", err)
	}

	return isValid, nil
}
