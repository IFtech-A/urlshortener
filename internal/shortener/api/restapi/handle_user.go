package restapi

import (
	"net/http"
	"strconv"

	"github.com/IFtech-A/urlshortener/internal/shortener/model"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Server) userCreate(c echo.Context) error {

	user := new(model.User)
	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	passHash, err := user.GeneratePasswordHash()
	if err != nil {
		logrus.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}
	user.Password = string(passHash)

	err = s.store.User().Create(user)
	if err != nil {
		logrus.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}
	logrus.Debug(user)
	cUser := *user
	cUser.Password = ""

	return c.JSON(http.StatusCreated, cUser)
}

func (s *Server) userRead(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logrus.Error(err.Error())
		return c.NoContent(http.StatusNotFound)
	}

	user, err := s.store.User().Get(id)
	if err != nil {
		logrus.Error(err.Error())
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, user)
}

func (s *Server) login(c echo.Context) error {

	user := new(model.User)
	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	existingUser, err := s.store.User().GetByUsername(user.Username)
	if err != nil {
		logrus.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	if err := existingUser.CheckPassword(user.Password); err != nil {
		logrus.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	token, err := existingUser.GenerateToken([]byte(s.config.TokenSecret))
	if err != nil {
		logrus.Error(err)
		return c.NoContent(http.StatusBadGateway)
	}

	cUser := *existingUser
	cUser.Password = ""

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
		"user":  cUser,
	})
}
