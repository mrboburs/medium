package repository

import (
	// "errors"
	"database/sql"
	"fmt"
	"mediumuz/model"
	"mediumuz/util/logrus"
	"time"

	// "github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
)

type UserDB struct {
	db *sqlx.DB
}

func NewUserDB(db *sqlx.DB) *UserDB {
	return &UserDB{db: db}
}

func (repo *UserDB) GetAllUsers(logrus *logrus.Logger) (array []model.UserFull, err error) {
	rowsRs, err := repo.db.Query("SELECT id,email,first_name,last_name,city,is_verified,verification_date, interests,bio,account_image_path,phone,rating,post_views_count,follower_count,following_count,like_count,is_super_user,created_at FROM users")

	if err != nil {
		logrus.Infof("ERROR: not selecting data from sql %s", err.Error())
		// http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return array, err
	}

	Array := []model.UserFull{}
	defer rowsRs.Close()

	for rowsRs.Next() {
		mu := model.UserFull{}
		err = rowsRs.Scan(&mu.Id, &mu.Email, &mu.FirstName, &mu.Lastname, &mu.City, &mu.IsVerified, &mu.Verification_date, &mu.Interests, &mu.Bio, &mu.AccountImagePath, &mu.Rating, &mu.Rating, &mu.PostViewsCount, &mu.FollowerCount, &mu.FollowingCount, &mu.LikeCount, &mu.IsSuperUser, &mu.CreatedAt)
		if err != nil {
			logrus.Infof("ERROR: not scanning data from sql %s", err.Error())
			// log.Println(err)
			// http.Error(w, http.StatusText(500), 500)
			return array, err
		}
		Array = append(Array, mu)
	}

	if err = rowsRs.Err(); err != nil {

		return Array, err
	}
	return Array, nil
}
func (repo *UserDB) DeleteUser(id string, logrus *logrus.Logger) error {

	_, err := repo.db.Exec("DELETE from users WHERE id = $1", id)
	if err != nil {
		logrus.Errorf("ERROR: deleting user failed : %v", err)
		return err
	}
	return nil
}

func (repo *UserDB) UpdateProfile(id int, user model.UserUpdate, logrus *logrus.Logger) (int64, error) {
	tm := time.Now()
	query := fmt.Sprintf("	UPDATE %s SET first_name=$1, city=$2, phone=$3, updated_at=$4, last_name=$5,bio=$6,interests=$7 	WHERE id = $8  RETURNING id ", usersTable)
	rows, err := repo.db.Exec(query, user.FirstName, user.City, user.Phone, tm, user.LastName, user.Bio, user.Interests, id)

	if err != nil {
		logrus.Errorf("ERROR: Updating failed : %v", err)
		return 0, err
	}
	effectedRowsNum, err := rows.RowsAffected()
	if err != nil {
		logrus.Errorf("ERROR:  effectedRowsNum failed : %v", err)
		return 0, err
	}
	logrus.Info("DONE:    updated profile data saved")
	return effectedRowsNum, nil
}

func (repo *UserDB) GetUserById(id string, logrus *logrus.Logger) (model.UserFull, error) {
	var user model.UserFull
	query := fmt.Sprintf("SELECT id,email,first_name,last_name,city,is_verified,verification_date, interests,bio,account_image_path,phone,rating,post_views_count,follower_count,following_count,like_count,is_super_user,created_at		FROM %s WHERE id=$1 ", usersTable)

	err := repo.db.Get(&user, query, id)
	if err == sql.ErrNoRows {
		logrus.Info("invalid id")
	}
	// repo.db.

	if err != nil {
		logrus.Errorf("ERROR: don't get users %s", err)
		return user, err
	}
	logrus.Info("DONE:get user data from psql")
	return user, nil
}

func (repo *UserDB) UpdateAccountImage(id int, filePath string, logrus *logrus.Logger) (int64, error) {
	tm := time.Now()
	query := fmt.Sprintf("	UPDATE %s SET account_image_path=$1,updated_at=$2	WHERE id = $3  RETURNING id ", usersTable)
	rows, err := repo.db.Exec(query, filePath, tm, id)

	if err != nil {
		logrus.Errorf("ERROR: Update photo failed : %v", err)
		return 0, err
	}
	effectedRowsNum, err := rows.RowsAffected()
	if err != nil {
		logrus.Errorf("ERROR: Update  photo failed: %v", err)
		return 0, err
	}
	logrus.Info("DONE:Updated photo successfully")
	return effectedRowsNum, nil
}

func (repo *UserDB) CheckUserId(id int, logrus *logrus.Logger) (int, error) {
	var countID int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE id=$1", usersTable)
	err := repo.db.Get(&countID, query, id)

	if err != nil {
		logrus.Infof("ERROR:Email query error: %s", err.Error())
		return -1, err
	}

	return countID, nil
}
