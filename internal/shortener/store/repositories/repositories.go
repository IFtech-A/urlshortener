package repositories

import "github.com/IFtech-A/urlshortener/internal/shortener/model"

type UserRepo interface {
	Create(*model.User) error
	Get(int64) (*model.User, error)
	GetByUsername(string) (*model.User, error)
	Update(*model.User) error
	Delete(int64) error
}

type BrandRepo interface {
	Create(*model.Brand) error
	Get(string) (*model.Brand, error)
	Update(*model.Brand) error
	Delete(string) error
}

type URLRepo interface {
	Create(*model.User, *model.URL) error
	Get(string) (*model.URL, error)
	Update(*model.URL) error
	Delete(string) error
	ReadUserLinks(*model.User) ([]*model.URL, error)
}
