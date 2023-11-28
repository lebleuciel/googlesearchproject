package file

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"

	"github.com/lib/pq"

	"github.com/lebleuciel/maani/models"
	"github.com/lebleuciel/maani/pkg/database"
	"github.com/lebleuciel/maani/pkg/helpers"
	"github.com/lebleuciel/maani/pkg/settings"
	"go.uber.org/zap"
)

// logger is a global variable for logging using Zap.
var logger *zap.SugaredLogger

// init initializes the Zap logger.
func init() {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	logger = zapLogger.Sugar()
}

type FileRepository struct {
	st settings.Settings
	db database.Database
}

// Is file valid
func (f *FileRepository) IsValidFile(file *multipart.FileHeader) error {
	err := f.db.AddFileTypeIfNotExist(file.Header.Get("Content-Type"))
	if err != nil {
		logger.Errorw("can't add file types into database", "error", err)
		return errors.New("can't add file types into database")
	}

	filetypes, err := f.db.GetFileTypes()
	if err != nil {
		logger.Errorw("can't get file types from database", "error", err)
		return errors.New("can't get file types from database")
	}

	for _, types := range filetypes {
		if types.Name == file.Header.Get("Content-Type") {
			if types.IsBanned {
				return fmt.Errorf("can't send file with %s type, for filename: %s", file.Header.Get("Content-Type"), file.Filename)
			}
			if file.Size > int64(types.AllowedSize) {
				return fmt.Errorf("file size is not allowed, you can send %s file with maximum %d byets, for filename: %s", types.Name, types.AllowedSize, file.Filename)
			}
			return nil
		}
	}
	return fmt.Errorf("file type %s not found, for filename: %s", file.Header.Get("Content-Type"), file.Filename)
}

// Save a file
func (f *FileRepository) SaveEncryptedFile(file models.File) error {
	currentSize, err := f.db.GetFilesSize()
	if err != nil {
		logger.Errorw("can't get files size from database", "error", err)
		return err
	}
	if f.st.BackendServer.MaxFilesSizeByte < currentSize+file.Size {
		return fmt.Errorf("reach maximum amount of disk usage")
	}

	uid, err := helpers.SaveEncryptedFile(file.Content, f.st.BackendServer.FilePath, []byte(f.st.BackendServer.EncryptKey))
	if err != nil {
		logger.Errorw("can't saved encrypted file from file repository", "error", err)
		return err
	}
	file.UUID = uid
	err = f.db.SaveFile(file)
	if err != nil {
		logger.Errorw("can't saved file into database from file repository", "error", err)
		return err
	}
	return nil
}

func (f *FileRepository) GetEncryptedFile(name []string, tags []string) (database.Transaction, models.File, error) {
	ctx := context.Background()
	tx, err := f.db.NewSerializableTransaction(ctx)

	defer func() {
		if err != nil {
			if e, ok := err.(*pq.Error); !ok || e.Code != database.ErrSerializationFailure {
				rollbackErr := tx.Rollback()
				if rollbackErr != nil {
					err = fmt.Errorf("error while rolling back transaction. original error: %w", err)
				}
			}
		}
	}()

	if err != nil {
		return nil, models.File{}, err
	}

	file, err := tx.GetFile(name, tags)
	if err != nil {
		return nil, models.File{}, nil
	}
	return tx, file, nil
}

func NewFileRepository(st settings.Settings, db database.Database) (*FileRepository, error) {
	if db == nil {
		return nil, errors.New("db should not be nil")
	}
	return &FileRepository{
		st: st,
		db: db,
	}, nil
}
