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

const SessionCookieName = "urlshortener-session"

func (s *Server) urlCreate(c echo.Context) error {
	cc, ok := c.(*CustomContext)
	var user *model.User
	if ok {
		user = cc.User
	}

	session, err := s.cookieStore.Get(c.Request(), SessionCookieName)
	if err != nil {
		logrus.Error(err.Error())
	}
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		logrus.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	URL := &model.URL{}
	if user != nil {
		URL.UserID = user.ID
	}

	err = json.Unmarshal(body, URL)

	if err != nil {
		logrus.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	if !govalidator.IsURL(URL.RealURL) {
		logrus.Error("invalid url")
		return c.NoContent(http.StatusBadRequest)
	}

	fullUrl, err := url.Parse(URL.RealURL)
	if err == nil && fullUrl.Scheme != "https" && fullUrl.Scheme != "http" {
		logrus.Error("bad scheme: only http or https scheme allowed")
		return c.NoContent(http.StatusBadRequest)
	}

	s.store.URL().Create(URL)

	if user == nil {
		urlsInterface, exists := session.Values["myurls"]
		var urlsBytes []byte
		urls := make([]*model.URL, 0)

		if exists {
			urlsBytes = urlsInterface.([]byte)
			logrus.Debugf("url history: %v", string(urlsBytes))
			err := json.Unmarshal(urlsBytes, &urls)
			if err != nil {
				logrus.Error(err.Error())
			}
		}

		urls = append(urls, URL)

		if len(urls) > 5 {
			urls = urls[len(urls)-5:]
		}

		//save back to session
		urlsBytes, err = json.Marshal(urls)
		if err != nil {
			logrus.Error(err.Error())
		} else {
			session.Values["myurls"] = urlsBytes
		}
		session.Save(c.Request(), c.Response())
	} else {
		urlsInterface, exists := session.Values["myurls"]
		if exists {
			delete(session.Values, "myurls")
			urls := make([]*model.URL, 0)
			urlsBytes := urlsInterface.([]byte)
			logrus.Debugf("url history: %v", string(urlsBytes))
			err := json.Unmarshal(urlsBytes, &urls)
			if err != nil {
				logrus.Error(err.Error())
			}

			// Add to user's urls
			for _, url := range urls {
				url.UserID = user.ID
				s.store.URL().Update(url)
			}
		}
		session.Save(c.Request(), c.Response())
	}

	return c.JSON(http.StatusCreated, URL)
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

func (s *Server) urlReadHistory(c echo.Context) error {
	var user *model.User
	cc, ok := c.(*CustomContext)
	if ok {
		user = cc.User
	}

	urls := make([]*model.URL, 0)

	// If user data exists
	if user != nil {

	} else {
		// Read from cookie
		session, err := s.cookieStore.Get(c.Request(), SessionCookieName)
		if err != nil {
			logrus.Error(err.Error())
		}
		urlsInterface, exists := session.Values["myurls"]
		var urlsBytes []byte

		if exists {
			urlsBytes = urlsInterface.([]byte)
			logrus.Debugf("url history: %v", string(urlsBytes))
			err := json.Unmarshal(urlsBytes, &urls)
			if err != nil {
				logrus.Error(err.Error())
			}
		}

	}

	return c.JSON(http.StatusOK, urls)

}
