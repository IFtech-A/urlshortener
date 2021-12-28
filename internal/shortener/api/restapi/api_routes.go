package restapi

import (
	"net/http"
	"strconv"
	"time"

	"github.com/IFtech-A/urlshortener/internal/shortener/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type CustomContext struct {
	echo.Context
	User        *model.User
	Fingerprint string
}

func (s *Server) configureAPIRoutes() {
	apiJWTConfig := middleware.DefaultJWTConfig
	apiJWTConfig.SigningKey = []byte(s.config.TokenSecret)
	apiJWTConfig.SuccessHandler = jwtAuthSuccessHandler
	apiJWTConfig.Skipper = userAuthSkipper()

	g := s.e.Group("/api")

	//Custom Context set
	contextMW := s.CustomContextMiddleware

	//jwt authentication
	jwtAuth := middleware.JWTWithConfig(apiJWTConfig)

	//request limiter
	limiter := middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20))

	//fingerprint reader
	sessionReader := echo.MiddlewareFunc(s.FingerprintMiddleware)

	//skip jwt auth

	// restrictedMWGroup := []echo.MiddlewareFunc{limiter, jwtAuth}
	// unrestrictedMWGroup := []echo.MiddlewareFunc{limiter}
	restrictedMWGroup := []echo.MiddlewareFunc{limiter, contextMW, sessionReader, jwtAuth}
	unrestrictedMWGroup := []echo.MiddlewareFunc{limiter, contextMW, sessionReader}

	/* restricted */
	/* user */
	g.GET(UserEndpoint+"/:id", s.userRead, restrictedMWGroup...)
	g.GET(UserEndpoint+"/:id/url", nil, restrictedMWGroup...)
	g.PUT(UserEndpoint, nil, restrictedMWGroup...)
	g.DELETE(UserEndpoint, nil, restrictedMWGroup...)

	// redirect server should directly read from database
	// g.GET("/url/:shortURL", s.urlRead, restrictedMWGroup...)

	/* unrestricted */
	/* user */
	g.POST(UserEndpoint, s.userCreate, unrestrictedMWGroup...)

	/* login */
	g.POST(LoginEndpoint, s.login, unrestrictedMWGroup...)

	/* url */
	g.POST("/url", s.urlCreate, unrestrictedMWGroup...)
	// return urls for the fingerprint
	g.GET("/url", s.urlReadHistory, unrestrictedMWGroup...)

	g.GET("", func(c echo.Context) error {
		cc := c.(*CustomContext)

		return c.String(http.StatusOK, "hello "+cc.User.Username)
	}, restrictedMWGroup...)

}

func (s *Server) CustomContextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := &CustomContext{
			Context: c,
		}

		return next(cc)
	}
}

func (s *Server) FingerprintMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		cookie, err := c.Cookie(SessionCookieName)
		if err != nil || cookie.Expires.After(time.Now()) {
			logrus.Warn("invalid session key")
			return next(c)
		} else {
			cc, ok := c.(*CustomContext)
			if ok {
				cc.Fingerprint = cookie.Value
			}

			return next(c)
		}

	}
}

func (s *Server) UserAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc, ok := c.(*CustomContext)
		if !ok {
			logrus.Error("incorrect context")
			return echo.ErrBadGateway
		}

		if cc.User == nil {
			logrus.Warn("no user authentication data")
			logrus.Debug("user nil")
			return next(c)
		}
		if cc.User.ID == 0 {
			//skipped jwt
			return next(c)
		}

		var err error
		cc.User, err = s.store.User().Get(cc.User.ID)
		if err != nil {
			logrus.Error(err.Error())
			return c.NoContent(http.StatusBadRequest)
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

	cc, ok := c.(*CustomContext)
	if ok {
		cc.User = &model.User{
			ID: userID,
		}
	}
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
