package redirect

import (
	"net/http"
	"strings"

	"github.com/IFtech-A/urlshortener/internal/shortener/store"
	"github.com/sirupsen/logrus"
)

type Server struct {
	c *Config
	s *store.Store
}

func New(conf *Config) *Server {
	return &Server{
		c: conf,
	}

}

func (s *Server) redirectHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		rw.WriteHeader(http.StatusBadGateway)
		return
	}

	pathDirs := strings.Split(r.URL.Path, "/")
	logrus.Debug(pathDirs)
	if len(pathDirs) > 2 {
		rw.WriteHeader(http.StatusBadGateway)
		return
	}

	http.Redirect(rw, r, "https://www.google.com", http.StatusMovedPermanently)
	return

}

func (s *Server) Start() error {

	logrus.SetLevel(logrus.DebugLevel)

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.redirectHandler)

	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}

	return server.ListenAndServe()
}
