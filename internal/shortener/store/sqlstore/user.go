package sqlstore

import "github.com/IFtech-A/urlshortener/internal/shortener/model"

type userRepo struct {
}

func (r *userRepo) Create(u *model.User) error {
	return nil
}
func (r *userRepo) Get(id int64) (*model.User, error) {

	return nil, nil
}
func (r *userRepo) GetByUsername(username string) (*model.User, error) {

	return nil, nil
}

func (r *userRepo) Update(u *model.User) error {

	return nil
}
func (r *userRepo) Delete(id int64) error {

	return nil
}
