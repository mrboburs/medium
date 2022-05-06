package repository

import (
	"mediumuz/model"
	"mediumuz/util/logrus"

	// "github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user model.User, logrus *logrus.Logger) (int, error)
	GetUser(email, password string) (model.User, error)
	GetUserID(username string, logrus *logrus.Logger) (int, error)
	VerifyEmail(id int, logrus *logrus.Logger) (int64, error)
}
type User interface {
	GetUserData(id int, logrus *logrus.Logger) (model.UserFull, error)
	UpdateAccountImage(id int, filePath string, logrus *logrus.Logger) (int64, error)
	UpdateProfile(id int, username string, city string, phone string, logrus *logrus.Logger) (int64, error)
}

type Repository struct {
	Authorization
	User
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{Authorization: NewAuthDB(db), User: NewUserDB(db)}
}
