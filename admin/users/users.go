package users

import (
	"fmt"

	"github.com/gin-gonic/gin"
	userRepository "github.com/lebleuciel/maani/pkg/repository/user"
	userService "github.com/lebleuciel/maani/pkg/services/user"
)

type Users struct {
	repository *userRepository.UserRepository
	service    *userService.UserService
	// authMiddleware *auth.auth
	authEnabled bool
}

func (u *Users) RegisterRoutes(v1 *gin.RouterGroup) {
	fmt.Println("registering user related endpoints to admin server")
	users := v1.Group("/user")
	users.GET("/list", u.getUserList())
}
func (u *Users) getUserList() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		u.service.GetUserList(ctx, true)
	}
}
func NewUserModule(userService *userService.UserService, userRepo *userRepository.UserRepository, authEnabled bool) (*Users, error) {
	if userService == nil {
		return nil, ErrNilUserService
	}
	if userRepo == nil {
		return nil, ErrNilUserRepo
	}
	return &Users{
		repository:  userRepo,
		service:     userService,
		authEnabled: authEnabled,
	}, nil
}
