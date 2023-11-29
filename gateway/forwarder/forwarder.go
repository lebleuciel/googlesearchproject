package forwarder

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lebleuciel/maani/models"
	"github.com/lebleuciel/maani/pkg/services/auth"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func init() {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	logger = zapLogger.Sugar()
}

type Forwarder struct {
	adminUrl       string
	backendUrl     string
	authMiddleware *auth.Auth
	authEnabled    bool
	userHeaderKey  string
}

func (u *Forwarder) RegisterRoutes(v1 *gin.RouterGroup) {
	fmt.Println("Registering related endpoints to gateway server")

	file := v1.Group("/file")
	user := v1.Group("/user")
	if u.authEnabled {
		file.Use(u.authMiddleware.Middleware())
		user.Use(u.authMiddleware.Middleware())
	}
	file.Any("", u.forward(u.backendUrl, false))
	file.Any("/list", u.forward(u.adminUrl, true))
	user.Any("/list", u.forward(u.adminUrl, true))
}

func (u *Forwarder) forward(url string, shouldBeAdmin bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if u.authEnabled {
			userData, err := u.checkAuthorizedRequest(ctx, shouldBeAdmin)
			if err != nil {
				return
			}

			// Create a new GET request to the other code
			req, err := http.NewRequest(ctx.Request.Method, url+ctx.FullPath(), ctx.Request.Body)
			if err != nil {
				logger.Errorw("can not create new request in gatewey", "error", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
				return
			}

			// Copy headers from the original request to the new request
			for key, values := range ctx.Request.Header {
				for _, value := range values {
					req.Header.Add(key, value)
				}
			}
			req.Header.Add(u.userHeaderKey, fmt.Sprint(userData.Id))

			// Perform the HTTP request
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				logger.Errorw("can not do request in gatewey", "error", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
				return
			}
			defer resp.Body.Close()
			for key, values := range resp.Header {
				for _, value := range values {
					ctx.Header(key, value)
				}
			}

			// Copy the response from the other code to the current response
			var buf bytes.Buffer
			buf.ReadFrom(resp.Body)
			ctx.String(resp.StatusCode, buf.String())
			return
		}
		ctx.JSON(http.StatusForbidden, gin.H{})
	}
}

// checkAuthorizedRequest checks for user access scope (separated for later RBAC implementation)
func (u *Forwarder) checkAuthorizedRequest(c *gin.Context, shouldBeAdmin bool) (models.UserWithPassword, error) {
	userData, err := auth.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request: " + err.Error(),
		})
		return models.UserWithPassword{}, err
	}
	if shouldBeAdmin && userData.AccessType != models.AdminType {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Forbidden",
		})
		return models.UserWithPassword{}, errors.New("forbidden")
	}
	return userData, nil
}

func NewForwarderModule(auth *auth.Auth, storeHost string, adminPort int, backendPort int, userHeaderKey string, authEnabled bool) (*Forwarder, error) {
	if storeHost == "" {
		return nil, ErrEmptyStoreHost
	}
	if adminPort == 0 {
		return nil, ErrEmptyAdminPort
	}
	if backendPort == 0 {
		return nil, ErrEmptyBackendPort
	}
	if userHeaderKey == "" {
		return nil, ErrEmptyUserHeaderKey
	}
	return &Forwarder{
		adminUrl:       fmt.Sprintf("%s:%d", storeHost, adminPort),
		backendUrl:     fmt.Sprintf("%s:%d", storeHost, backendPort),
		authMiddleware: auth,
		authEnabled:    authEnabled,
		userHeaderKey:  userHeaderKey,
	}, nil
}
