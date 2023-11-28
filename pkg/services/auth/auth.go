package auth

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/lebleuciel/maani/models"
	"github.com/lebleuciel/maani/pkg/helpers"
	"github.com/lebleuciel/maani/pkg/repository/user"
	"github.com/pkg/errors"
)

type Auth struct {
	userRepository *user.Repository
	middleware     *jwt.GinJWTMiddleware
}

func (a *Auth) GetGinAuthMiddleware() *jwt.GinJWTMiddleware {
	return a.middleware
}

func (a *Auth) Middleware() gin.HandlerFunc {
	return a.middleware.MiddlewareFunc()
}

func (a *Auth) RegisterRoutes(group *gin.RouterGroup) {
	fmt.Println("Registering auth related endpoints")

	group.POST("/auth/login", a.middleware.LoginHandler)
	group.POST("/auth/logout", a.middleware.LogoutHandler)
	group.POST("/auth/refresh", a.middleware.RefreshHandler)
	group.POST("/auth/register", a.RegisterUserHandler())
}

func (a *Auth) RegisterUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		inputData := models.UserRegisterParameters{}
		err := c.ShouldBindJSON(&inputData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		userData, err := a.userRepository.CreateUser(models.UserCreationParameters{
			FirstName:  inputData.FirstName,
			LastName:   inputData.LastName,
			Email:      inputData.Email,
			Password:   helpers.Encrypt(inputData.Password),
			AccessType: models.CustomerType,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, userData)
	}
}

func getAuthenticator(userRepository *user.Repository) func(*gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var creds models.UserLoginCredentials
		err := c.ShouldBindJSON(&creds)
		if err != nil {
			return nil, jwt.ErrMissingLoginValues
		}

		u, err := userRepository.GetUserByEmail(creds.Email)
		if err != nil {
			return nil, jwt.ErrFailedAuthentication
		}

		err = userRepository.UpdateUserLastLogin(u.Id)
		if err != nil {
			return nil, errors.Wrap(err, "Could not updating user last login")
		}

		hash := md5.Sum([]byte(creds.Password))
		hashedPassword := hex.EncodeToString(hash[:])

		if u.Password != hashedPassword {
			return nil, jwt.ErrFailedAuthentication
		}
		return &models.UserWithPassword{
			Id:          u.Id,
			FirstName:   u.FirstName,
			LastName:    u.LastName,
			Email:       u.Email,
			AccessType:  u.AccessType,
			CreatedAt:   u.CreatedAt,
			UpdatedAt:   u.UpdatedAt,
			LastLoginAt: u.LastLoginAt,
			Password:    u.Password,
		}, nil
	}
}

func getAuthorizer(userRepository *user.Repository) func(interface{}, *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {
		if _, ok := data.(*models.UserWithPassword); ok {
			return true
		}

		return false
	}
}

func getPayloadFunc(userRepository *user.Repository) func(interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		if v, ok := data.(*models.UserWithPassword); ok {
			return jwt.MapClaims{
				"email": v.Email,
			}
		}
		return jwt.MapClaims{}
	}
}

func getUnauthorizedFunc(userRepository *user.Repository) func(*gin.Context, int, string) {
	return func(c *gin.Context, code int, message string) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized: " + message,
		})
	}
}

func getIdentityHandlerFunc(userRepository *user.Repository) func(*gin.Context) interface{} {
	return func(c *gin.Context) interface{} {
		claims := jwt.ExtractClaims(c)
		username := claims["email"].(string)

		u, err := userRepository.GetUserByEmail(username)
		if err != nil {
			return nil
		}

		return &models.UserWithPassword{
			Id:          u.Id,
			FirstName:   u.FirstName,
			LastName:    u.LastName,
			Email:       u.Email,
			AccessType:  u.AccessType,
			CreatedAt:   u.CreatedAt,
			UpdatedAt:   u.UpdatedAt,
			LastLoginAt: u.LastLoginAt,
			Password:    u.Password,
		}
	}
}

func GetUserFromContext(c *gin.Context) (models.UserWithPassword, error) {
	userDataRaw, found := c.Get("email")
	if !found {
		return models.UserWithPassword{}, ErrUserObjectNotFound
	}
	userData, ok := userDataRaw.(*models.UserWithPassword)
	if !ok {
		return models.UserWithPassword{}, ErrInvalidUserObjectType
	}
	return *userData, nil
}

func NewAuth(userRepository *user.Repository, secretKey, identityKey, realm string, timeout, maxRefresh time.Duration) (*Auth, error) {
	if userRepository == nil {
		return nil, ErrNilUserRepo
	}

	if secretKey == "" {
		return nil, ErrEmptySecretKey
	}
	middleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            realm,
		SigningAlgorithm: "HS512",
		Key:              []byte(secretKey),
		Timeout:          timeout,
		MaxRefresh:       maxRefresh,
		Authenticator:    getAuthenticator(userRepository),
		Authorizator:     getAuthorizer(userRepository),
		PayloadFunc:      getPayloadFunc(userRepository),
		Unauthorized:     getUnauthorizedFunc(userRepository),
		IdentityHandler:  getIdentityHandlerFunc(userRepository),
		IdentityKey:      identityKey,
		TokenLookup:      "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:    "Bearer",
		TimeFunc:         time.Now,
		SendCookie:       false,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Could not create jwt auth module")
	}
	return &Auth{
		userRepository: userRepository,
		middleware:     middleware,
	}, nil
}
