package admin

import (
	"github.com/pkg/errors"

	"github.com/lebleuciel/maani/admin/files"
	"github.com/lebleuciel/maani/admin/server"
	"github.com/lebleuciel/maani/pkg/database"
	FileRepository "github.com/lebleuciel/maani/pkg/repository/file"
	FileService "github.com/lebleuciel/maani/pkg/services/file"
	"github.com/lebleuciel/maani/pkg/settings"
)

func NewAdminServer(setting settings.Settings, database database.Database) (*server.Server, error) {
	// Initialize Repositories
	fileRepo, err := FileRepository.NewFileRepository(setting, database)
	if err != nil {
		return nil, errors.Wrap(err, "Could not initialize file repository")
	}

	// Initialize Services
	fileService, err := FileService.NewFileService(fileRepo, setting)
	if err != nil {
		return nil, errors.Wrap(err, "Could not initialize file service")
	}

	// Initialize API Modules
	fileModule, err := files.NewFileModule(fileService, fileRepo, true)
	if err != nil {
		return nil, errors.Wrap(err, "Could not initialize new file module")
	}

	srv, err := server.NewServer(fileModule)
	if err != nil {
		return nil, errors.Wrap(err, "Could not initialize new admin server")
	}
	return srv, nil
}
