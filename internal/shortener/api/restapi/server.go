package restapi

import (
	"errors"

	"github.com/IFtech-A/urlshortener/internal/shortener/store"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

const UserIDContextKey = "userID"
const UserEndpoint = "/user"
const LoginEndpoint = "/login"
const RestApiEndpoint = "/api"

type Config struct {
	Address     string
	StoreType   string
	TokenSecret string
}

func NewConfig() *Config {
	return &Config{
		Address:     "0.0.0.0:8080",
		StoreType:   "memory",
		TokenSecret: "hello world",
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
