package database

import (
	"context"
	"database/sql"

	"github.com/lebleuciel/maani/models"
)

const (
	None                    = "none"
	PostgresSQL             = "pgsql"
	ErrSerializationFailure = "40001"
)

type Database interface {
	TransactionMethods
	UsersDatabaseMethods
	FilesDatabaseMethods
}

type (
	TransactionMethods interface {
		NewSerializableTransaction(ctx context.Context) (Transaction, error)
		NewTransaction(ctx context.Context, isolation sql.IsolationLevel) (Transaction, error)
	}

	// UsersDatabaseMethods to manage Users Repository Methods
	UsersDatabaseMethods interface {
		GetUserByEmail(email string) (*models.UserWithPassword, error)
		CreateUser(spec models.UserCreationParameters) (models.User, error)
		UpdateUserLastLogin(userId int) error
	}

	// FilesDatabaseMethods to manage Files Repository Methods
	FilesDatabaseMethods interface {
		AddFileTypeIfNotExist(string) error
		GetFileTypes() ([]models.FileType, error)
		GetFilesSize() (int, error)
		SaveFile(models.File) error
		GetFile([]string, []string) (models.File, error)
	}
)

type Transaction interface {
	UsersDatabaseMethods
	FilesDatabaseMethods
	Commit() error
	Rollback() error
}
