package service

import (
	// "fmt"
	"io"
	"mediumuz/configs"
	"mediumuz/model"
	"mediumuz/package/repository"
	"mediumuz/util/logrus"
	"mime/multipart"
	"os"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (service *UserService) UpdateProfile(id int, user model.UserUpdate, logrus *logrus.Logger) (int64, error) {
	count, err := service.repo.UpdateProfile(id, user, logrus)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (service *UserService) GetUserById(id string, logrus *logrus.Logger) (model.UserFull, error) {
	return service.repo.GetUserById(id, logrus)

}

func (service *UserService) UploadAccountImage(file multipart.File, header *multipart.FileHeader, logrus *logrus.Logger) (string, error) {

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

func (service *UserService) UpdateAccountImage(id int, filePath string, logrus *logrus.Logger) (int64, error) {
	return service.repo.UpdateAccountImage(id, filePath, logrus)
}

func (service *UserService) CheckUserId(id int, logrus *logrus.Logger) (int, error) {
	return service.repo.CheckUserId(id, logrus)
}
func (s *UserService) GetAllUsers(logrus *logrus.Logger) (array []model.UserFull, err error) {
	return s.repo.GetAllUsers(logrus)
}

func (s *UserService) DeleteUser(id string, logrus *logrus.Logger) error {
	return s.repo.DeleteUser(id, logrus)
}
