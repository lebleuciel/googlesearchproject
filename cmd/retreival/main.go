package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	gatewayapi "github.com/lebleuciel/maani/gateway"
	"github.com/lebleuciel/maani/pkg/database"
	"github.com/lebleuciel/maani/pkg/database/postgres"
	"github.com/lebleuciel/maani/pkg/settings"
	"github.com/spf13/pflag"
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

// settingsPath is a variable to store the path to the settings file.
var settingsPath string

// main is the main entry point of the application.
func main() {
	logger.Infoln("Retrieval is running")

	// Parse command line flags.
	pflag.StringVar(&settingsPath, "settings", "/opt/maani/settings.yml", "Path to settings file")
	pflag.Parse()

	// Read settings from the configuration file.
	var st settings.Settings
	var db database.Database
	err := cleanenv.ReadConfig(settingsPath, &st)
	if err != nil {
		logger.Fatalw("Could not read settings", "error", err.Error())
	}
	_, err = st.IsValid()
	if err != nil {
		logger.Fatalw("Setting file is not valid", "error", err.Error())
	}

	logger.Infoln("Initializing database")

	// Initialize the database.
	db = initDatabase(st)

	logger.Infoln("Setup router")

	// Setup HTTP servers and run them.
	gatewayServer := setupHttpServers(st, db)
	go runServer(gatewayServer, "gateway_server")

	// Handle shutdown signals.
	shutDown := make(chan os.Signal, 1)
	signal.Notify(shutDown, syscall.SIGHUP, syscall.SIGINT)
	signal := <-shutDown

	logger.Infow("Shutting down the server", "signal", signal)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := gatewayServer.Shutdown(ctx); err != nil {
		logger.Fatalw("Could not shutdown gateway API server gracefully", "error", err.Error())
	}
}

// initDatabase initializes the database based on the provided settings.
func initDatabase(settings settings.Settings) (db database.Database) {
	logger.Infow("Initializing database client", "database_type", settings.Database.Type)
	if settings.Database.Type == database.PostgresSQL {
		db, err := postgres.NewPostgresDatabase(postgres.PGOptions{
			SSLMode:             settings.Database.SSLMode,
			Host:                settings.Database.Host,
			User:                settings.Database.User,
			DBName:              settings.Database.DatabaseName,
			Password:            settings.Database.Password,
			Port:                settings.Database.Port,
			MaxOpenConnections:  settings.Database.MaxOpenConnections,
			MaxIdleConnections:  settings.Database.MaxIdleConnections,
			ConnMaxLifetime:     settings.Database.ConnMaxLifetime,
			ConnMaxIdleTime:     settings.Database.ConnMaxIdleTime,
			StatusCheckInterval: settings.Database.StatusCheckInterval,
			Timeout:             settings.Database.QueryTimeout,
			BaseContext:         context.Background(),
		}, true)
		if err != nil {
			logger.Fatalw("Could not create postgres database client", "error", err.Error())
		}
		return db
	}
	return nil
}

// setupHttpServers initializes and returns an HTTP server based on the provided settings and database.
func setupHttpServers(settings settings.Settings, database database.Database) *http.Server {
	logger.Info("Initializing HTTP servers.")
	gatewayServerHandler, err := gatewayapi.NewGatewayServer(settings, database, "gateway", settings.GatewayServer.SecretKey, settings.GatewayServer.TokenTimeout, settings.GatewayServer.RefreshTokenTimeout)
	if err != nil {
		logger.Fatalw("Could not initialize Gateway Server", "error", err.Error())
	}

	gatewayAddress := fmt.Sprintf(":%d", settings.Global.GatewayPort)
	gatewayServer := &http.Server{
		Addr:              gatewayAddress,
		Handler:           gatewayServerHandler,
		ReadTimeout:       settings.Global.ReadTimeout,
		ReadHeaderTimeout: settings.Global.ReadHeaderTimeout,
		WriteTimeout:      settings.Global.WriteTimeout,
		IdleTimeout:       settings.Global.IdleTimeout,
		MaxHeaderBytes:    settings.Global.MaxHeaderBytes,
	}
	return gatewayServer
}

// runServer starts the HTTP server and logs any errors.
func runServer(server *http.Server, serverName string) {
	logger.Infoln(serverName, "Starting listening on port", server.Addr)
	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		logger.Fatalw("Could not create "+serverName+" server listener", "error", err.Error())
	}
	err = server.Serve(ln)
	if err != nil {
		logger.Infow("Serving failed", "error", err.Error(), "serverName", serverName)
	}
}
