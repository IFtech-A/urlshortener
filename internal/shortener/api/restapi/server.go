package restapi

import (
	"net/http"
	"strconv"

	"github.com/IFtech-A/urlshortener/internal/shortener/model"
	"github.com/IFtech-A/urlshortener/internal/shortener/store"
	"github.com/IFtech-A/urlshortener/internal/shortener/store/memorystore"
	"github.com/IFtech-A/urlshortener/internal/shortener/store/sqlstore"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
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
	e      *echo.Echo
	store  store.Store
	config *Config
}

func New(conf *Config) *Server {
	return &Server{
		e:      echo.New(),
		config: conf,
	}
}

type CustomContext struct {
	echo.Context
	User *model.User
}

func (s *Server) configureRoutes() {

	apiJWTConfig := middleware.DefaultJWTConfig
	apiJWTConfig.SigningKey = []byte(s.config.TokenSecret)
	apiJWTConfig.SuccessHandler = jwtAuthSuccessHandler
	apiJWTConfig.Skipper = userAuthSkipper()

	g := s.e.Group("/api", middleware.JWTWithConfig(apiJWTConfig), s.UserAuthMiddleware)
	/* user */
	g.POST(UserEndpoint, s.userCreate)
	g.GET(UserEndpoint+"/:id", s.userRead)
	g.PUT(UserEndpoint, nil)
	g.DELETE(UserEndpoint, nil)

	/* login */
	g.POST(LoginEndpoint, s.login)

	g.GET("", func(c echo.Context) error {
		cc := c.(*CustomContext)

		return c.String(http.StatusOK, "hello "+cc.User.Username)
	})
	g.POST("/url", s.urlCreate)
	g.GET("/url/:shortURL", s.urlRead)
}

func (s *Server) UserAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, ok := c.Get(UserIDContextKey).(int64)
		if !ok {
			logrus.Error("incorrect context")
			return c.NoContent(http.StatusForbidden)
		}
		if userID == 0 {
			//skipped jwt
			return next(c)
		}

		user, err := s.store.User().Get(userID)
		if err != nil {
			logrus.Error(err.Error())
			return c.NoContent(http.StatusBadRequest)
		}

		cc := &CustomContext{
			Context: c,
			User:    user,
		}

		return next(cc)
	}
}

func jwtAuthSuccessHandler(c echo.Context) {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userIDStr, ok := claims["ID"].(string)
	if !ok {
		logrus.Error("jwt ID malformed")
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	c.Set(UserIDContextKey, userID)
}

func userAuthSkipper() middleware.Skipper {

	const userRegistrationEndpoint = RestApiEndpoint + UserEndpoint
	const loginEndpoint = RestApiEndpoint + LoginEndpoint
	return func(c echo.Context) bool {
		if c.Request().Method != http.MethodPost {
			return false
		}

		if c.Request().URL.Path != loginEndpoint && c.Request().URL.Path != userRegistrationEndpoint {
			return false
		}

		c.Set(UserIDContextKey, int64(0))
		return true
	}
}

func (s *Server) configureStore() error {

	switch s.config.StoreType {
	case "memory":
		s.store = memorystore.New()
	case "sql":
		s.store = sqlstore.New()

	}

	return nil
}

func (s *Server) Start() error {

	err := s.configureStore()
	if err != nil {
		return err
	}

	s.configureRoutes()

	return s.e.Start(s.config.Address)
}
