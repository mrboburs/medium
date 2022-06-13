package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"mediumuz/model"
	"mediumuz/package/repository"
	"mediumuz/util/logrus"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type AuthService struct {
	repo repository.Authorization
}

const (
	salt       = "hjqrhjqw124617aj564u564a654u65465aufhajs"
	signingKey = "qrkjk#4#%35FSFJlja#4353K554245987uaS?SFjH"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (service *AuthService) CreateUser(user model.User, logrus *logrus.Logger) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	logrus.Info("successfully password_hash")
	return service.repo.CreateUser(user, logrus)

}
func (s *AuthService) IsVerified(id int, logrus *logrus.Logger) (int, error) {
	count, err := s.repo.IsVerified(id, logrus)
	if err != nil {
		return 0, err
	}
	return int(count), nil

}
func (s *AuthService) GenerateToken(email, password string, logrus *logrus.Logger) (string, error) {
	user, err := s.repo.GetUser(email, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}
	return claims.UserId, nil
}
func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) GetUser(email, password string) (model.User, error) {
	return s.repo.GetUser(email, password)

}

func (s *AuthService) GetUserID(email string, logrus *logrus.Logger) (int, error) {
	return s.repo.GetUserID(email, logrus)
}
