package service

import (
	"mediumuz/model"
	"mediumuz/package/repository"
	"mediumuz/util/logrus"
	"mime/multipart"
)

type Authorization interface {
	CreateUser(user model.User, logrus *logrus.Logger) (int, error)
	GenerateToken(email, password string, logrus *logrus.Logger) (string, error)
	VerifyEmail(id int, logrus *logrus.Logger) (int, error)
	// SendMessageEmail(email string, userName string, logrus *logrus.Logger) error
	// CheckDataExists(email string, logrus *logrus.Logger) (int, error)
	ParseToken(token string) (int, error)
}

type User interface {
	GetUserData(id int, logrus *logrus.Logger) (user model.UserFull, err error)
	UpdateProfile(id int, username string, city string, phone string, logrus *logrus.Logger) (int, error)

	UploadAccountImage(file multipart.File, header *multipart.FileHeader, user model.UserFull, logrus *logrus.Logger) (filePath string, err error)
	UpdateAccountImage(id int, filePath string, logrus *logrus.Logger) (int64, error)
}
type Service struct {
	Authorization
	User
}

func NewService(repos *repository.Repository) *Service {
	return &Service{Authorization: NewAuthService(repos.Authorization), User: NewUserService(repos.User)}
}
