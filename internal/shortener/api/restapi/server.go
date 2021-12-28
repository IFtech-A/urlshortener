package restapi

import (
	"errors"
	"net/http"

	"github.com/IFtech-A/urlshortener/internal/shortener/store"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

const UserIDContextKey = "userID"
const UserEndpoint = "/user"
const LoginEndpoint = "/login"
const RestApiEndpoint = "/api"

const (
	Prod  = "production"
	Debug = "debug"
	Test  = "test"
)

type Config struct {
	Address     string `env:"HOST_ADDR"`
	StoreType   string `env:"STORE_TYPE"`
	TokenSecret string `env:"TOKEN_SECRET"`
	Mode        string `env:"MODE"`
}

func NewConfig() *Config {
	return &Config{
		Address:     "0.0.0.0:8080",
		StoreType:   "memory",
		TokenSecret: "hello world",
		Mode:        Prod,
	}
}

type Server struct {
	e           *echo.Echo
	store       store.Store
	cookieStore *sessions.CookieStore
	config      *Config
}

func New(conf *Config) *Server {
	return &Server{
		e:      echo.New(),
		config: conf,
	}
}

func (s *Server) configureRoutes() {
	s.configureAPIRoutes()
	s.configureRedirectRoutes()
	s.configureFrontendRoutes()

	s.e.GET("/health_check", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	if s.config.Mode == Debug {
		s.e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowCredentials: true,
		}))
		logrus.Debug("Running on DEBUG mode")
	}

}

func (s *Server) configureStore() error {
	s.store = store.CreateStore(s.config.StoreType)
	if s.store == nil {
		return errors.New("invalid store")
	}

	return nil
}

func (s *Server) Start() error {

	err := s.configureStore()
	if err != nil {
		return err
	}

	s.cookieStore = sessions.NewCookieStore([]byte(s.config.TokenSecret))

	s.configureRoutes()

	return s.e.Start(s.config.Address)
}
