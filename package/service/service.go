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
	IsVerified(id int, logrus *logrus.Logger) (int, error)
	GetUser(email, password string) (model.User, error)
	GetUserID(email string, logrus *logrus.Logger) (int, error)
	ParseToken(token string) (int, error)
}

type User interface {
	GetUserById(id string, logrus *logrus.Logger) (user model.UserFull, err error)
	GetAllUsers(logrus *logrus.Logger) (array []model.UserFull, err error)
	UpdateProfile(id int, user model.UserUpdate, logrus *logrus.Logger) (int64, error)
	CheckUserId(id int, logrus *logrus.Logger) (int, error)
	UploadAccountImage(file multipart.File, header *multipart.FileHeader, logrus *logrus.Logger) (filePath string, err error)
	UpdateAccountImage(id int, filePath string, logrus *logrus.Logger) (int64, error)
	DeleteUser(id string, logrus *logrus.Logger) error
}

type Post interface {
	GetCommentsPost(postID int, pagination model.Pagination, logrus *logrus.Logger) (comments []model.CommentFull, err error)
	CommentPost(input model.CommentPost, logrus *logrus.Logger) (int, error)
	ViewPost(userID, postID int, logrus *logrus.Logger) (int, error)
	ClickLike(userId int, postId int, logrus *logrus.Logger) error
	UpdatePost(id int, input model.Post, logrus *logrus.Logger) error
	GetAllPosts(pagination model.Pagination, logrus *logrus.Logger) (posts []model.PostFull, err error)
	CreatePost(userId int, post model.Post, logrus *logrus.Logger) (int, error)
	GetPostById(postId int, logrus *logrus.Logger) (post model.PostFull, err error)
	CheckPostId(id int, logrus *logrus.Logger) (int, error)
	PostDelete(userId, postId int) error
	UpdatePostImage(userId, postId int, input model.PostFull, filePath string, logrus *logrus.Logger) error
	UploadPostImage(file multipart.File, header *multipart.FileHeader, post model.PostFull, logrus *logrus.Logger) (string, error)
}
type Service struct {
	Authorization
	User
	Post
}

func NewService(repos *repository.Repository) *Service {
	return &Service{Authorization: NewAuthService(repos.Authorization), User: NewUserService(repos.User), Post: NewPostService(repos.Post)}
}
