package service

import (
	// "fmt"
	"io"
	"mediumuz/configs"
	"mediumuz/model"
	"mediumuz/package/repository"
	"mime/multipart"
	"os"

	// "mediumuz/util/convert"
	"mediumuz/util/logrus"
)

type PostService struct {
	repo repository.Post
}

func NewPostService(repo repository.Post) *PostService {
	return &PostService{repo: repo}
}

func (service *PostService) CreatePost(userId int, post model.Post, logrus *logrus.Logger) (int, error) {
	postId, err := service.repo.CreatePost(post, logrus)
	if err != nil {
		return 0, err
	}
	postUserId, err := service.repo.CreatePostUser(userId, postId, logrus)
	if err != nil {
		return 0, err
	}
	return postUserId, err
}
func (s *PostService) GetCommentsPost(postID int, pagination model.Pagination, logrus *logrus.Logger) (comments []model.CommentFull, err error) {
	return s.repo.GetCommentsPost(postID, pagination, logrus)
}
func (s *PostService) GetAllPosts(pagination model.Pagination, logrus *logrus.Logger) (posts []model.PostFull, err error) {
	return s.repo.GetAllPosts(pagination, logrus)
}
func (s *PostService) GetPostById(postId int, logrus *logrus.Logger) (post model.PostFull, err error) {
	// post.PostImagePath = convert.EmptyStringToNull()
	return s.repo.GetPostById(postId, logrus)

}
func (s *PostService) CommentPost(input model.CommentPost, logrus *logrus.Logger) (int, error) {
	return s.repo.CommentPost(input, logrus)
}
func (s *PostService) ViewPost(userID, postID int, logrus *logrus.Logger) (int, error) {
	return s.repo.ViewPost(userID, postID, logrus)
}
func (s *PostService) ClickLike(userId int, postId int, logrus *logrus.Logger) error {
	return s.repo.ClickLike(userId, postId, logrus)
}

func (s *PostService) UpdatePost(id int, input model.Post, logrus *logrus.Logger) error {
	return s.repo.UpdatePost(id, input, logrus)
}
func (s *PostService) CheckPostId(id int, logrus *logrus.Logger) (int, error) {
	return s.repo.CheckPostId(id, logrus)
}

// func (s *PostService) UpdatePostImage(id int, filePath string, logrus *logrus.Logger) (int64, error) {

// }

func (s *PostService) PostDelete(userId, postId int) error {
	return s.repo.PostDelete(userId, postId)
}

func (s *PostService) UpdatePostImage(userId, postId int, input model.PostFull, filePath string, logrus *logrus.Logger) error {
	return s.repo.UpdatePostImage(userId, postId, input, filePath, logrus)
}
func (s *PostService) UploadPostImage(file multipart.File, header *multipart.FileHeader, post model.PostFull, logrus *logrus.Logger) (string, error) {

	filename := header.Filename
	folderPath := "public/"
	err := os.MkdirAll(folderPath, 0777)
	if err != nil {
		logrus.Errorf("ERROR: Failed to create folder %s: %v", folderPath, err)
		return "", err
	}
	err = os.Chmod(folderPath, 0777)
	if err != nil {
		logrus.Errorf("ERROR: Failed Access Permission denied %s", err)
		return "", err
	}
	filePath := folderPath + filename
	out, err := os.Create(filePath)
	if err != nil {
		logrus.Errorf("ERROR: Failed CreateFile: %v", err)
		return "", err
	}

	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		logrus.Errorf("ERROR: Failed copy %s", err)
		return "", err
	}
	configs, err := configs.InitConfig()
	logrus.Infof("configs %v", configs)
	if err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}
	imageURL := configs.ServiceHost + "/" + filePath
	return imageURL, nil
}
