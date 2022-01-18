package memorystore

import (
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/IFtech-A/urlshortener/internal/shortener/model"
	"github.com/sirupsen/logrus"
)

//A-Za-z0-9
const sizeMultiplier = 26 + 26 + 10

var currentDBSize int32 = 3
var urlGeneratorSize int32 = 1

type urlRepo struct {
	s  *Store
	mu sync.RWMutex
}

func init() {
	for i := int32(0); i < currentDBSize; i++ {
		urlGeneratorSize *= sizeMultiplier
	}
	rand.Seed(time.Now().UnixNano())
}

func indexLetter(index int) rune {

	if index < 0 {
		logrus.Fatal("invalid index")
		return '-'
	}
	if index < 26 {
		return rune(index + int('A'))
	} else if index < 52 {
		return rune(index - 26 + int('a'))
	} else if index < 62 {
		return rune(index - 52 + int('0'))
	} else {
		logrus.Error("index out of range. max 62")
	}
	return '-'
}

func generateShortURL(url string) string {

	short := strings.Builder{}
	for i := int32(0); i < currentDBSize; i++ {
		index := rand.Intn(sizeMultiplier)

		short.WriteRune(indexLetter(index))
	}

	return short.String()
}

func (r *urlRepo) Create(o *model.User, u *model.URL) error {

	if o != nil {
		u.UserID = o.ID
	}

	if int32(len(r.s.urlDB)) > urlGeneratorSize-1000 {
		atomic.AddInt32(&currentDBSize, 1)
		atomic.AddInt32(&urlGeneratorSize, urlGeneratorSize*sizeMultiplier)
	}

	//Premium url
	if u.ShortenedURL != "" {
		_, err := r.Get(u.ShortenedURL)
		if err != nil {
			if err == ErrNotFound {
				r.mu.Lock()
				defer r.mu.Unlock()
				r.s.urlDB[u.ShortenedURL] = u
				return nil
			} else {
				return err
			}
		} else {
			return ErrAlreadyExists
		}
	}

	var err error
	var short string
	for err == nil {
		short = generateShortURL(u.RealURL)
		_, err = r.Get(short)
	}
	if err != ErrNotFound {
		return err
	}

	r.mu.Lock()
	r.s.urlDB[short] = u
	u.ShortenedURL = short
	defer r.mu.Unlock()

	return nil
}
func (r *urlRepo) Get(shortenedURL string) (*model.URL, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()
	URL, exists := r.s.urlDB[shortenedURL]
	if !exists {
		return nil, ErrNotFound
	}

	return URL, nil
}

func (r *urlRepo) ReadUserLinks(u *model.User) ([]*model.URL, error) {
	urls := make([]*model.URL, 0)
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, v := range r.s.urlDB {
		if v.UserID == u.ID {
			urls = append(urls, v)
		}
	}
	return urls, nil
}

func (r *urlRepo) Update(u *model.URL) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	URL, exists := r.s.urlDB[u.ShortenedURL]
	if !exists {
		return ErrNotFound
	}
	URL.UserID = u.UserID
	URL.RealURL = u.RealURL
	URL.UpdatedAt = time.Now()

	return nil
}
func (r *urlRepo) Delete(urlname string) error {

	return nil
}
