package models

import "time"

const (
	// CustomerType Access Type: only access to backend APIs
	CustomerType = "Customer"

	// AdminType Access Type: full access to backend and admin APIs
	AdminType = "Admin"

	// IdentityKey Authentication Identity Field Name
	IdentityKey = "email"
)

// User general object contains user details
type User struct {
	Id          int        `json:"id"`
	FirstName   string     `json:"firstName"`
	LastName    string     `json:"lastName"`
	Email       string     `json:"email"`
	AccessType  string     `json:"accessType"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	LastLoginAt *time.Time `json:"LastLoginAt"`
}

// UserWithPassword private object to retrieve user's full details
type UserWithPassword struct {
	Password    string     `json:"password"`
	Id          int        `json:"id"`
	FirstName   string     `json:"firstName"`
	LastName    string     `json:"lastName"`
	Email       string     `json:"email"`
	AccessType  string     `json:"accessType"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	LastLoginAt *time.Time `json:"LastLoginAt"`
}

// UserLoginCredentials login credential
type UserLoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserRegisterParameters input parameters for creating users
type UserRegisterParameters struct {
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

// UserCreationParameters input parameters for creating users
type UserCreationParameters struct {
	Password   string `json:"password"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	AccessType string `json:"accessType"`
}

// UserTokenResponse successful login response object for JWT token output
type UserTokenResponse struct {
	Code   int       `json:"code"`
	Expire time.Time `json:"expire"`
	Token  string    `json:"token"`
}
