package repository

import (
	// "errors"
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

func (repo *UserDB) GetUserData(id string, logrus *logrus.Logger) (model.UserFull, error) {
	var user model.UserFull
	query := fmt.Sprintf("SELECT  	id,	email,	firstname,	secondname,	city,	is_verified,	account_image_path,	phone,	rating,	post_views,	is_super_user	FROM %s WHERE id=$1 ", usersTable)
	err := repo.db.Get(&user, query, id)
	if err != nil {
		logrus.Errorf("ERROR: don't get users %s", err)
		return user, err
	}
	logrus.Info("DONE:get user data from psql")
	return user, nil
}

// func (repo *UserDB) UpdateUserVerified(id string, logrus *logrus.Logger) (effectedRowsNum int64, err error) {
// 	tm := time.Now()
// 	query := fmt.Sprintf("	UPDATE %s SET is_verified = true,verification_date=$1,updated_at=$2	WHERE id = $3  RETURNING id ", usersTable)
// 	rows, err := repo.db.Exec(query, tm, tm, id)

// 	if err != nil {
// 		logrus.Errorf("ERROR: Update verificationCode failed : %v", err)
// 		return 0, err
// 	}
// 	effectedRowsNum, err = rows.RowsAffected()
// 	if err != nil {
// 		logrus.Errorf("ERROR: Update verificationCode effectedRowsNum failed : %v", err)
// 		return 0, err
// 	}
// 	logrus.Info("DONE:  email verified")
// 	return effectedRowsNum, nil
// }

func (repo *UserDB) UpdateAccountImage(id int, filePath string, logrus *logrus.Logger) (int64, error) {
	tm := time.Now()
	query := fmt.Sprintf("	UPDATE %s SET account_image_path=$1,updated_at=$2	WHERE id = $3  RETURNING id ", usersTable)
	rows, err := repo.db.Exec(query, filePath, tm, id)

	if err != nil {
		logrus.Errorf("ERROR: Update verificationCode : %v", err)
		return 0, err
	}
	effectedRowsNum, err := rows.RowsAffected()
	if err != nil {
		logrus.Errorf("ERROR: Update verificationCode effectedRowsNum : %v", err)
		return 0, err
	}
	logrus.Info("DONE:Update verify email")
	return effectedRowsNum, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// REDIS

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// func (repo *UserDB) CheckCode(email, code string, logrus *logrus.Logger) error {
// 	saveCode, err := repo.redis.Get(email).Result()
// 	if err != nil {
// 		logrus.Errorf("ERROR:don't save code %s", err)
// 		return err
// 	}
// 	if saveCode != code {
// 		return errors.New("code not found ")
// 	}
// 	logrus.Info("DONE: verify code")
// 	return nil
// }
