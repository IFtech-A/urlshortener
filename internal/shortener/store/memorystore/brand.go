package memorystore

import "github.com/IFtech-A/urlshortener/internal/shortener/model"

type brandRepo struct {
}

func (r *brandRepo) Create(u *model.Brand) error {
	return nil
}
func (r *brandRepo) Get(brandname string) (*model.Brand, error) {

	return nil, nil
}

func (r *brandRepo) Update(u *model.Brand) error {

	return nil
}
func (r *brandRepo) Delete(brandname string) error {

	return nil
}
