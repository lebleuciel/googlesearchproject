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
	adminapi "github.com/lebleuciel/maani/admin"
	backendapi "github.com/lebleuciel/maani/backend"
	"github.com/lebleuciel/maani/pkg/database"
	"github.com/lebleuciel/maani/pkg/database/postgres"
	"github.com/lebleuciel/maani/pkg/settings"
	"github.com/spf13/pflag"
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

var settingsPath string

func main() {
	logger.Infoln("Store is running")
	pflag.StringVar(&settingsPath, "settings", "/opt/maani/settings.yml", "Path to settings file")
	pflag.Parse()

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

	// init database
	db = initDatabase(st)

	logger.Infoln("Setup router")

	backendServer, adminServer := setupHttpServers(st, db)
	go runServer(backendServer, "backend_server")
	go runServer(adminServer, "admin_server")

	shutDown := make(chan os.Signal, 1)
	signal.Notify(shutDown, syscall.SIGHUP, syscall.SIGINT)
	signal := <-shutDown

	logger.Infow("Shutting down the server", "signal", signal)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := backendServer.Shutdown(ctx); err != nil {
		logger.Fatalw("Could not shutdown backend api server gracefully", "error", err.Error())
	}
	if err := adminServer.Shutdown(ctx); err != nil {
		logger.Fatalw("Could not shutdown backend api server gracefully", "error", err.Error())
	}
}

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

func setupHttpServers(settings settings.Settings, database database.Database) (*http.Server, *http.Server) {
	logger.Info("Initializing http servers")
	backendServerHandler, err := backendapi.NewBackendServer(settings, database)
	if err != nil {
		logger.Fatalw("Could not initialize Backend Server", "error", err.Error())
	}

	backendAddress := fmt.Sprintf(":%d", settings.Global.BackendPort)
	backendServer := &http.Server{
		Addr:              backendAddress,
		Handler:           backendServerHandler,
		ReadTimeout:       settings.Global.ReadTimeout,
		ReadHeaderTimeout: settings.Global.ReadHeaderTimeout,
		WriteTimeout:      settings.Global.WriteTimeout,
		IdleTimeout:       settings.Global.IdleTimeout,
		MaxHeaderBytes:    settings.Global.MaxHeaderBytes,
	}

	adminServerHandler, err := adminapi.NewAdminServer(settings, database)
	if err != nil {
		logger.Fatalw("Could not initialize Admin Server", "error", err.Error())
	}

	adminAddress := fmt.Sprintf(":%d", settings.Global.AdminPort)
	adminServer := &http.Server{
		Addr:              adminAddress,
		Handler:           adminServerHandler,
		ReadTimeout:       settings.Global.ReadTimeout,
		ReadHeaderTimeout: settings.Global.ReadHeaderTimeout,
		WriteTimeout:      settings.Global.WriteTimeout,
		IdleTimeout:       settings.Global.IdleTimeout,
		MaxHeaderBytes:    settings.Global.MaxHeaderBytes,
	}
	return backendServer, adminServer
}

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
