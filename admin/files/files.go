package files

import (
	"fmt"

	"github.com/gin-gonic/gin"
	fileRepository "github.com/lebleuciel/maani/pkg/repository/file"
	fileService "github.com/lebleuciel/maani/pkg/services/file"
)

type Files struct {
	repository *fileRepository.FileRepository
	service    *fileService.FileService
	// authMiddleware *auth.auth
	authEnabled bool
}

func (u *Files) RegisterRoutes(v1 *gin.RouterGroup) {
	fmt.Println("registering file related endpoints to admin server")
	files := v1.Group("/file")
	files.GET("/", u.getCurrentFile())
}
func (u *Files) getCurrentFile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		u.service.GetFile(ctx, false)
		return
		// ctx.JSON(200, "file data returned")
	}
}
func NewFileModule(fileService *fileService.FileService, fileRepo *fileRepository.FileRepository, authEnabled bool) (*Files, error) {
	if fileService == nil {
		return nil, ErrNilFileService
	}
	if fileRepo == nil {
		return nil, ErrNilFileRepo
	}
	return &Files{
		repository:  fileRepo,
		service:     fileService,
		authEnabled: authEnabled,
	}, nil
}
