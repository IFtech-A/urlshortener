package sqlstore

import "github.com/IFtech-A/urlshortener/internal/shortener/model"

type urlRepo struct {
	s *Store
}

func (r *urlRepo) Create(u *model.URL) error {
	return nil
}
func (r *urlRepo) Get(urlname string) (*model.URL, error) {

	return nil, nil
}

func (r *urlRepo) Update(u *model.URL) error {

	return nil
}
func (r *urlRepo) Delete(urlname string) error {

	return nil
}
