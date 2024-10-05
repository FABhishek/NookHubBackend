package models

// User represents a user in the system
type RegisterUser struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"` // In practice, will store a hashed password (later)
	DisplayName string `json:"displayName"`
}

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
