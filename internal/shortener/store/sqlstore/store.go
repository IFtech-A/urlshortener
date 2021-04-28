package sqlstore

import "github.com/IFtech-A/urlshortener/internal/shortener/store/repositories"

type Store struct {
	userRepo *userRepo
	urlRepo  *urlRepo
}

func New() *Store {
	return &Store{}
}

func (s *Store) User() repositories.UserRepo {
	if s.userRepo == nil {
		s.userRepo = &userRepo{}
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
