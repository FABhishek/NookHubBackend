package models

// User represents a user in the system
type RegisterUser struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"` // In practice, will store a hashed password (later)
	DisplayName string `json:"displayName"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Homies struct {
	Name    string
	Id      int
	Pfp     string
	IsAdmin bool
}
