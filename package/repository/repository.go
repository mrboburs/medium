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
	IsVerified(id int, logrus *logrus.Logger) (int64, error)
}
type User interface {
	GetUserById(id string, logrus *logrus.Logger) (model.UserFull, error)
	GetAllUsers(logrus *logrus.Logger) (array []model.UserFull, err error)
	DeleteUser(id string, logrus *logrus.Logger) error
	UpdateAccountImage(id int, filePath string, logrus *logrus.Logger) (int64, error)
	UpdateProfile(id int, user model.UserUpdate, logrus *logrus.Logger) (int64, error)
	CheckUserId(id int, logrus *logrus.Logger) (int, error)
}
type Post interface {
	GetCommentsPost(postID int, pagination model.Pagination, logrus *logrus.Logger) (comments []model.CommentFull, err error)
	CommentPost(input model.CommentPost, logrus *logrus.Logger) (int, error)
	ClickLike(userId int, postId int, logrus *logrus.Logger) error
	ViewPost(userID, postID int, logrus *logrus.Logger) (int, error)
	UpdatePost(id int, input model.Post, logrus *logrus.Logger) error
	CreatePost(post model.Post, logrus *logrus.Logger) (int, error)
	CreatePostUser(userId, postId int, logrus *logrus.Logger) (int, error)
	GetPostById(postId int, logrus *logrus.Logger) (post model.PostFull, err error)
	GetAllPosts(pagination model.Pagination, logrus *logrus.Logger) (posts []model.PostFull, err error)
	CheckPostId(id int, logrus *logrus.Logger) (int, error)
	UpdatePostImage(userId, postId int, input model.PostFull, filePath string, logrus *logrus.Logger) error
	PostDelete(userId, postId int) error
}

type Repository struct {
	Authorization
	User
	Post
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{Authorization: NewAuthDB(db), User: NewUserDB(db), Post: NewPostDB(db)}
}
