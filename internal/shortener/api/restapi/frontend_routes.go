package restapi

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

const staticRootDir = "./front"

func (s *Server) handleStaticFiles(c echo.Context) error {

	filePath := staticRootDir + c.Request().URL.Path

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.File(staticRootDir + "/index.html")
	}

	return c.File(filePath)
}

func (s *Server) configureFrontendRoutes() {
	s.e.GET("", func(c echo.Context) error {
		url := c.Request().URL
		url.Scheme = "http"
		if c.Request().TLS != nil {
			url.Scheme = "https"
		}

		if len(url.Host) == 0 {
			url.Host = c.Request().Host
		}

		return c.Redirect(http.StatusMovedPermanently, c.Request().URL.String()+"app")
	})

	s.e.GET("/static/*", s.handleStaticFiles)
	s.e.File("/app", staticRootDir+"/index.html")
	s.e.File("/manifest.json", staticRootDir)

	g := s.e.Group("/app")
	g.Static("/index.html", staticRootDir)
	g.Static("/", staticRootDir)
	g.Any("/*", s.handleStaticFiles)
}
