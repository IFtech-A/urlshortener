package restapi

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Server) redirectHandler(c echo.Context) error {
	shortURL := c.Param("shortURL")
	if shortURL == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	extendedURL, err := s.store.URL().Get(shortURL)
	if err != nil {
		logrus.Error(err)
		return c.NoContent(http.StatusNotFound)
	}

	return c.Redirect(http.StatusMovedPermanently, extendedURL.RealURL)
}

func (s *Server) configureRedirectRoutes() {
	s.e.GET("/:shortURL", s.redirectHandler)
}
