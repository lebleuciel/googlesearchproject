package gateway

import (
	"time"

	"github.com/pkg/errors"

	"github.com/lebleuciel/maani/gateway/forwarder"
	"github.com/lebleuciel/maani/gateway/server"
	"github.com/lebleuciel/maani/models"
	"github.com/lebleuciel/maani/pkg/database"
	"github.com/lebleuciel/maani/pkg/repository/user"
	"github.com/lebleuciel/maani/pkg/services/auth"
	"github.com/lebleuciel/maani/pkg/settings"
)

func NewGatewayServer(settings settings.Settings, database database.Database, realm string, secretKey string, tokenTimeout, refreshTokenTimeout time.Duration) (*server.Server, error) {
	// Initialize Repositories
	userRepo, err := user.NewUserRepository(database)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize user repository")
	}

	// Initialize API Modules
	authModule, err := auth.NewAuth(userRepo, secretKey, models.IdentityKey, realm, tokenTimeout, refreshTokenTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize auth module")
	}

	fileModule, err := forwarder.NewForwarderModule(
		authModule,
		settings.GatewayServer.StoreHost,
		settings.Global.AdminPort,
		settings.Global.BackendPort,
		settings.GatewayServer.UserIdHeaderKey,
		true,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Could not initialize new file module")
	}

	srv, err := server.NewServer(authModule, fileModule)
	if err != nil {
		return nil, errors.Wrap(err, "Could not initialize new gateway server")
	}
	return srv, nil
}
