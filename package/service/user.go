package service

import (
	"fmt"
	"io"
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

func (s *UserService) ConfirmEmail(code string, logrus *logrus.Logger) {

}

func (service *UserService) GetUserData(id string, logrus *logrus.Logger) (user model.UserFull, err error) {
	user, err = service.repo.GetUserData(id, logrus)
	if err != nil {
		logrus.Error("ERROR: get user Data failed: %v", err)
		return user, err
	}
	return user, nil
}

func (service *UserService) UploadAccountImage(file multipart.File, header *multipart.FileHeader, user model.UserFull, logrus *logrus.Logger) (string, error) {

	filename := header.Filename
	folderPath := fmt.Sprintf("public/%s/", user.UserName)
	err := os.MkdirAll(folderPath, 0777)
	if err != nil {
		logrus.Errorf("ERROR: Failed to create folder %s: %v", folderPath, err)
		return "", err
	}
	err = os.Chmod(folderPath, 0777)
	if err != nil {
		logrus.Errorf("ERROR: Failed Accessss Permission denied %s", err)
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
	return filePath, nil
}

func (service *UserService) UpdateAccountImage(id int, filePath string, logrus *logrus.Logger) (int64, error) {
	return service.repo.UpdateAccountImage(id, filePath, logrus)
}
