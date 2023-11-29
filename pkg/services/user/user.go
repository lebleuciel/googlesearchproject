package user

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	repository "github.com/lebleuciel/maani/pkg/repository/user"
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

type UserService struct {
	st         settings.Settings
	repository *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository, st settings.Settings) (*UserService, error) {
	if repo == nil {
		return nil, ErrNilUserRepo
	}
	return &UserService{
		st:         st,
		repository: repo,
	}, nil
}

func (f *UserService) GetUserList(c *gin.Context, isAdmin bool) {
	users, err := f.repository.GetUserList()
	if err != nil {
		logger.Errorw("failed to get user list", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can not get users list"})
		return
	}
	c.JSON(http.StatusOK, users)
}
