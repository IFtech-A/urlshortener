package store

import "github.com/IFtech-A/urlshortener/internal/shortener/store/repositories"

type Store interface {
	User() repositories.UserRepo
	// Brand() BrandRepo
	URL() repositories.URLRepo
}
