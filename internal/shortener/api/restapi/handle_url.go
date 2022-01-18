package restapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/IFtech-A/urlshortener/internal/shortener/model"
	"github.com/IFtech-A/urlshortener/internal/shortener/store/memorystore"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const SessionCookieName = "urlshortener-session"
const CookieURLs = "myurls"

var ErrNoSessionCookie = errors.New("no session cookie")
var ErrNoSessionCookieUrls = errors.New("no session cookie urls")

func (s *Server) createURL(c echo.Context) error {

	URL := new(model.URL)
	if err := c.Bind(URL); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(URL); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var user *model.User
	if uc := c.Get(UserContextKey); uc != nil {
		user = uc.(*model.User)
	}

	URL.CreatedAt = time.Now()
	URL.UpdatedAt = time.Now()

	err := s.store.URL().Create(user, URL)
	if err != nil {
		if err == memorystore.ErrAlreadyExists {
			logrus.Debug("Requested URL already exists")
		} else {
			logrus.Error(err)
		}
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	logrus.Debugf("createURL %v", user)

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

	}

	return c.JSON(http.StatusCreated, URL)
}

func (s *Server) readURL(c echo.Context) error {

	short := c.Param("shortURL")

	url, err := s.store.URL().Get(short)
	if err != nil {
		logrus.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, url)
}

func (s *Server) readUserURL(c echo.Context) error {
	var user *model.User
	if uc := c.Get(UserContextKey); uc != nil {
		user = uc.(*model.User)
	}
	logrus.Debugf("readUserURL %v", user)

	existing := make([]*model.URL, 0)
	// If user data exists
	if user != nil {
		urls, err := s.readUrlsFromCookie(c.Request(), c.Response())
		if err == nil {
			// Cookie exists and have urls
			// merge them to users urls
			session, _ := s.cookieStore.Get(c.Request(), SessionCookieName)
			delete(session.Values, CookieURLs)

			logrus.Debug("Adding cookie URLs to the user's URL list")
			// Add to user's urls
			for _, url := range urls {
				url.UserID = user.ID
				s.store.URL().Update(url)
			}

			session.Options.MaxAge = -1
			session.Save(c.Request(), c.Response())
		}

		urls, err = s.store.URL().ReadUserLinks(user)
		if err != nil {
			logrus.Error(err)
		}
		existing = urls
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
		} else {
			logrus.Error(err)
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
