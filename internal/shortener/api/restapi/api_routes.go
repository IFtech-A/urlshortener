package restapi

import (
	"net/http"
	"strconv"

	"github.com/IFtech-A/urlshortener/internal/shortener/model"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

const UserContextKey = "user-context"
const UserCookieURLs = "user-cookie-urls"

const LoginEndpoint = "/login"
const ApiEndpoint = "/api"
const UserEndpoint = "/users"
const LinkEndpoint = "/links"

func (s *Server) configureAPIRoutes() {
	apiJWTConfig := middleware.DefaultJWTConfig
	apiJWTConfig.SigningKey = []byte(s.config.TokenSecret)
	apiJWTConfig.SuccessHandler = jwtAuthSuccessHandler
	apiJWTConfig.Skipper = userAuthSkipper()
	apiJWTConfig.Claims = &model.Claims{}

	api := s.e.Group(ApiEndpoint)
	userApi := api.Group(UserEndpoint)
	urlApi := api.Group(LinkEndpoint)

	//Custom Context set
	// contextMW := s.CustomContextMiddleware

	//jwt authentication
	jwtAuth := middleware.JWTWithConfig(apiJWTConfig)

	//request limiter
	limiter := middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20))

	restrictedMWGroup := []echo.MiddlewareFunc{limiter, jwtAuth}
	unrestrictedMWGroup := []echo.MiddlewareFunc{limiter}

	/* restricted */
	/* user */
	userApi.GET("/:id", s.userRead, restrictedMWGroup...)
	userApi.GET("/:id"+LinkEndpoint, s.readUserURL, restrictedMWGroup...)
	userApi.POST("/:id"+LinkEndpoint, s.createURL, restrictedMWGroup...)
	userApi.PUT("/:id", nil, restrictedMWGroup...)
	userApi.DELETE("/:id", nil, restrictedMWGroup...)

	// redirect server should directly read from database
	// g.GET("/url/:shortURL", s.urlRead, restrictedMWGroup...)

	/* unrestricted */
	/* user */
	userApi.POST("", s.userCreate, unrestrictedMWGroup...)

	/* login */
	api.POST(LoginEndpoint, s.login, unrestrictedMWGroup...)
	api.OPTIONS("*", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	/* url */
	urlApi.POST("", s.createURL, unrestrictedMWGroup...)
	urlApi.GET("", s.readUserURL, unrestrictedMWGroup...)
	urlApi.PUT("", nil, restrictedMWGroup...)
	urlApi.DELETE("", nil, restrictedMWGroup...)

	api.GET("", func(c echo.Context) error {
		var user *model.User
		if uc := c.Get(UserContextKey); uc != nil {
			user = uc.(*model.User)
		}

		return c.String(http.StatusOK, "hello "+strconv.FormatInt(user.ID, 10))
	}, restrictedMWGroup...)

}

func jwtAuthSuccessHandler(c echo.Context) {
	token := c.Get("user").(*jwt.Token)

	claims, ok := token.Claims.(*model.Claims)
	if !ok {
		logrus.Error("jwt ID malformed", claims)
		return
	}
	userID, err := strconv.ParseInt(claims.ID, 10, 64)
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	c.Set(UserContextKey, &model.User{
		ID: userID,
	})
}

func userAuthSkipper() middleware.Skipper {

	const userRegistrationEndpoint = ApiEndpoint + UserEndpoint
	const loginEndpoint = ApiEndpoint + LoginEndpoint
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
