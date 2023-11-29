package server

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lebleuciel/maani/admin/files"
	"github.com/lebleuciel/maani/admin/users"
)

type Server struct {
	enviroment string
	engine     *gin.Engine
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.engine.ServeHTTP(w, r)
}
func NewServer(files *files.Files, users *users.Users) (*Server, error) {
	if files == nil {
		return nil, ErrNilFileModule
	}
	if users == nil {
		return nil, ErrNilUserModule
	}

	gin.SetMode("release")
	engine := gin.New()
	engine.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           1 * time.Hour,
	}))

	v1 := engine.Group("/api")
	files.RegisterRoutes(v1)
	users.RegisterRoutes(v1)

	return &Server{
		enviroment: "release",
		engine:     engine,
	}, nil
}
