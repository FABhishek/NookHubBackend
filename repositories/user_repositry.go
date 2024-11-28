package repositories

import (
	"Nookhub/models"
	"database/sql"
	"fmt"
)

type UserRepository interface {
	CreateUser(user models.RegisterUser) (int, error)
	LoginUser(user models.LoginUser) (int, error)
	IsAvailable(input string, inputType string) (bool, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user models.RegisterUser) (int, error) {
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

func (r *userRepository) LoginUser(user models.LoginUser) (int, error) {
	userId := 0
	//call the procedure
	// return if  userId < 0 // if email isn't there then user doesn't exist plz singunp
	// in case if email exists then return the encrypted password with encryption key and decrypt in sv layer obv
	// if it matches then return success else "wrong password entered"\
	return userId, nil
}

// EmailChecker implements UserRepository.
func (r *userRepository) IsAvailable(input string, inputType string) (bool, error) {
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
