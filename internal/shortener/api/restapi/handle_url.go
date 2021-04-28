package restapi

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/IFtech-A/urlshortener/internal/shortener/model"
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Server) urlCreate(c echo.Context) error {
	cc := c.(*CustomContext)

	body, err := io.ReadAll(cc.Request().Body)
	if err != nil {
		logrus.Error(err.Error())
		return cc.NoContent(http.StatusBadRequest)
	}

	URL := &model.URL{
		UserID: cc.User.ID,
	}
	err = json.Unmarshal(body, URL)
	if err != nil {
		logrus.Error(err.Error())
		return cc.NoContent(http.StatusBadRequest)
	}

	if !govalidator.IsURL(URL.RealURL) {
		logrus.Error("invalid url")
		return cc.NoContent(http.StatusBadRequest)
	}

	fullUrl, err := url.Parse(URL.RealURL)
	if err == nil && fullUrl.Scheme != "https" && fullUrl.Scheme != "http" {
		logrus.Error("bad scheme: only http or https scheme allowed")
		return cc.NoContent(http.StatusBadRequest)
	}

	s.store.URL().Create(URL)

	return cc.JSON(http.StatusCreated, URL)
}

func (s *Server) urlRead(c echo.Context) error {

	short := c.Param("shortURL")

	url, err := s.store.URL().Get(short)
	if err != nil {
		logrus.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, url)
}
