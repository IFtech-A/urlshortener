package restapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/IFtech-A/urlshortener/internal/shortener/model"
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const SessionCookieName = "urlshortener-session"
const CookieURLs = "myurls"

var ErrNoSessionCookie = errors.New("no session cookie")
var ErrNoSessionCookieUrls = errors.New("no session cookie urls")

func (s *Server) urlCreate(c echo.Context) error {

	var URL *model.URL
	if url := c.Get("validated-url"); url != nil {
		URL = url.(*model.URL)
	}

	var user *model.User
	if uc := c.Get(UserContextKey); uc != nil {
		user = uc.(*model.User)
	}
	if user != nil {
		URL.UserID = user.ID
	}

	s.store.URL().Create(URL)

	// User does not exist
	if user == nil {
		urls, _ := s.readUrlsFromCookie(c.Request(), c.Response())
		urls = append(urls, URL)

		if len(urls) > 5 {
			urls = urls[len(urls)-5:]
		}

		err := s.saveUrlsToCookie(c.Request(), c.Response(), urls)
		if err != nil {
			logrus.Error(err)
		}
	} else {
		urls, err := s.readUrlsFromCookie(c.Request(), c.Response())
		if err == nil {
			// Cookie exists and have urls
			// merge them to users urls
			session, _ := s.cookieStore.Get(c.Request(), SessionCookieName)
			delete(session.Values, CookieURLs)

			// Add to user's urls
			for _, url := range urls {
				url.UserID = user.ID
				s.store.URL().Update(url)
			}

			session.Options.MaxAge = -1
			session.Save(c.Request(), c.Response())
		}
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
	if uc := c.Get(UserContextKey); uc != nil {
		user = uc.(*model.User)
	}

	existing := make([]*model.URL, 0)
	// If user data exists
	if user != nil {

	} else {
		urls, err := s.readUrlsFromCookie(c.Request(), c.Response())
		logrus.Debugf("Cookie URLs: %v", urls)

		if err == nil {
			for _, url := range urls {
				if _, err := s.store.URL().Get(url.ShortenedURL); err == nil {
					existing = append(existing, url)
				}
			}
			s.saveUrlsToCookie(c.Request(), c.Response(), existing)
		}

	}

	return c.JSON(http.StatusOK, existing)

}

func (s *Server) saveUrlsToCookie(r *http.Request, rw http.ResponseWriter, urls []*model.URL) error {
	session, _ := s.cookieStore.Get(r, SessionCookieName)

	//save back to session
	urlsBytes, err := json.Marshal(urls)
	if err != nil {
		logrus.Error(err)
		return err
	} else {
		session.Values[CookieURLs] = urlsBytes
	}
	session.Save(r, rw)

	return nil
}

func (s *Server) readUrlsFromCookie(r *http.Request, rw http.ResponseWriter) ([]*model.URL, error) {
	session, err := s.cookieStore.Get(r, SessionCookieName)
	if err != nil || session.IsNew {
		logrus.Debug(err)
		return nil, ErrNoSessionCookie
	}

	// Check if cookie store contains URLs
	urlsInterface, exists := session.Values[CookieURLs]
	if !exists {
		return nil, ErrNoSessionCookieUrls
	}

	// Return URLs from cookie

	urls := make([]*model.URL, 0)
	urlsBytes := urlsInterface.([]byte)

	err = json.Unmarshal(urlsBytes, &urls)
	if err != nil {
		logrus.Error(err.Error())
		// remove problemous cookie store entry
		delete(session.Values, CookieURLs)
		// save it to response so that it was not sent next time from client side
		session.Save(r, rw)

		return nil, err
	}

	return urls, nil
}

func validateURLMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			logrus.Error(err.Error())
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error":   true,
				"message": "invalid json body format",
			})
		}

		URL := &model.URL{}

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

		c.Set("validated-url", URL)

		return next(c)
	}
}
