package restapi

import (
	"errors"
	"net/http"

	"github.com/IFtech-A/urlshortener/internal/shortener/store"
	"github.com/go-playground/validator"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const UserIDContextKey = "userID"

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

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func (s *Server) configureRoutes() {
	s.e.Validator = &CustomValidator{validator: validator.New()}

	s.configureAPIRoutes()
	s.configureRedirectRoutes()
	s.configureFrontendRoutes()

	s.e.GET("/health_check", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	if s.config.Mode == Debug {
		s.e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				origin := c.Request().Header.Get("Origin")
				if origin != "" {
					c.Response().Header().Add("Access-Control-Allow-Origin", origin)
					c.Response().Header().Add("Access-Control-Allow-Credentials", "true")
					c.Response().Header().Add("Access-Control-Allow-Methods", c.Request().Method)
					reqHeaders := c.Request().Header.Get("Access-Control-Request-Headers")
					if reqHeaders != "" {
						c.Response().Header().Add("Access-Control-Allow-Headers", reqHeaders)
					}
				}
				return next(c)
			}
		})
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
