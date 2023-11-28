package user

import (
	"github.com/lebleuciel/maani/models"
	"github.com/lebleuciel/maani/pkg/database"
	"github.com/pkg/errors"
)

type Repository struct {
	db database.Database
}

func NewUserRepository(db database.Database) (*Repository, error) {
	if db == nil {
		return nil, ErrNilUserDatabase
	}
	return &Repository{
		db: db,
	}, nil
}

// GetUserByEmail get single user by email address
func (r *Repository) GetUserByEmail(email string) (*models.UserWithPassword, error) {
	user, err := r.db.GetUserByEmail(email)
	if err != nil {
		return nil, errors.Wrap(ErrGetUserByEmail, "Could not get User with given email address")
	}
	return user, nil
}

// CreateUser Creates new User
func (r *Repository) CreateUser(spec models.UserCreationParameters) (models.User, error) {
	user, err := r.db.CreateUser(spec)
	if err != nil {
		return models.User{}, errors.Wrap(err, "Could not create User with given specification")
	}
	return user, nil
}

func (r *Repository) UpdateUserLastLogin(userId int) error {
	err := r.db.UpdateUserLastLogin(userId)
	return err
}
