package model

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID int64 `json:"id"`

	Username string `json:"username"`
	Password string `json:"password"`

	EnabledBrands bool `json:"enabled_brand"`

	Brands []*Brand `json:"brands"`
}

type ApiToken struct {
}

func (u *User) GeneratePasswordHash() ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

type Claims struct {
	ID string
	jwt.StandardClaims
}

func (u *User) GenerateToken(signKey []byte) (string, error) {

	claims := Claims{
		ID: strconv.FormatInt(u.ID, 10),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(signKey)
}
