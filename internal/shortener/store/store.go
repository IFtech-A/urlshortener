package store

import (
	"github.com/IFtech-A/urlshortener/internal/shortener/store/memorystore"
	"github.com/IFtech-A/urlshortener/internal/shortener/store/repositories"
	"github.com/IFtech-A/urlshortener/internal/shortener/store/sqlstore"
)

type Store interface {
	User() repositories.UserRepo
	// Brand() BrandRepo
	URL() repositories.URLRepo
}

func CreateStore(storeType string) Store {

	switch storeType {
	case "memory":
		return memorystore.New()
	case "sql":
		return sqlstore.New()
	}

	return nil
}
