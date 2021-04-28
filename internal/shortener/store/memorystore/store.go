package memorystore

import (
	"sync"

	"github.com/IFtech-A/urlshortener/internal/shortener/model"
	"github.com/IFtech-A/urlshortener/internal/shortener/store/repositories"
)

var userIDSeq int64
var userIDSeq_Lock sync.Mutex

var brandIDSeq int64
var brandIDSeq_Lock sync.Mutex

var urlIDSeq int64
var urlIDSeq_Lock sync.Mutex

type Store struct {
	userRepo *userRepo
	urlRepo  *urlRepo
	userDB   map[int64]*model.User
	brandDB  map[int64]*model.Brand
	urlDB    map[string]*model.URL
}

func New() *Store {

	return &Store{
		userDB:  make(map[int64]*model.User),
		brandDB: make(map[int64]*model.Brand),
		urlDB:   make(map[string]*model.URL),
	}
}

func (s *Store) User() repositories.UserRepo {
	if s.userRepo == nil {
		s.userRepo = &userRepo{
			s: s,
		}
	}
	return s.userRepo
}

func (s *Store) URL() repositories.URLRepo {
	if s.urlRepo == nil {
		s.urlRepo = &urlRepo{
			s: s,
		}
	}
	return s.urlRepo
}
