package restapi

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/IFtech-A/urlshortener/internal/shortener/model"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Server) userCreate(c echo.Context) error {

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		logrus.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	user := &model.User{}
	err = json.Unmarshal(body, user)
	if err != nil {
		logrus.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
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
	logrus.Debug(user.ID)

	return c.NoContent(http.StatusCreated)
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

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		logrus.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	user := &model.User{}
	err = json.Unmarshal(body, user)
	if err != nil {
		logrus.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
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

	// cookie := &http.Cookie{
	// 	Name:     s.config.AuthHeader,
	// 	Value:    token,
	// 	HttpOnly: true,
	// }

	// c.SetCookie(cookie)
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
