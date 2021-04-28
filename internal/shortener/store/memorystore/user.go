package memorystore

import (
	"errors"
	"strings"
	"sync"

	"github.com/IFtech-A/urlshortener/internal/shortener/model"
)

type userRepo struct {
	s  *Store
	mu sync.Mutex
}

var ErrNotFound = errors.New("not found")
var ErrAlreadyExists = errors.New("already exists")

func (r *userRepo) Create(u *model.User) error {

	_, err := r.GetByUsername(u.Username)
	if err == nil {
		return ErrAlreadyExists
	}

	userIDSeq_Lock.Lock()
	userIDSeq++
	u.ID = userIDSeq
	userIDSeq_Lock.Unlock()

	r.mu.Lock()
	r.s.userDB[u.ID] = u
	r.mu.Unlock()

	return nil
}

func (r *userRepo) Get(id int64) (*model.User, error) {

	r.mu.Lock()
	defer r.mu.Unlock()
	user, exists := r.s.userDB[id]
	if !exists {
		return nil, ErrNotFound
	}

	return user, nil
}

func (r *userRepo) GetByUsername(username string) (*model.User, error) {

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, user := range r.s.userDB {
		if strings.EqualFold(user.Username, username) {
			return user, nil
		}
	}
	return nil, ErrNotFound
}

func (r *userRepo) Update(u *model.User) error {

	user, err := r.Get(u.ID)
	if err != nil {
		return err
	}

	r.mu.Lock()
	user.EnabledBrands = u.EnabledBrands
	r.mu.Unlock()

	return nil
}

func (r *userRepo) Delete(id int64) error {

	user, err := r.Get(id)
	if err != nil {
		return err
	}

	r.mu.Lock()
	delete(r.s.userDB, user.ID)
	r.mu.Unlock()

	return nil
}
