package repository

import (
	"errors"
	"fmt"
	"mediumuz/model"
	"mediumuz/util/logrus"
	"time"

	// "time"

	// "github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	// "github.com/lib/pq"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthDB(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (repo *AuthPostgres) CreateUser(user model.User, logrus *logrus.Logger) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (email,  first_name ,last_name, password, interests, bio, city, phone) values ($1, $2, $3, $4, $5,$6,$7,$8) RETURNING id", usersTable)

	row := repo.db.QueryRow(query, user.Email, user.FirstName, user.LastName, user.Password, user.Interests, user.Bio, user.City, user.Phone)

	if err := row.Scan(&id); err != nil {
		logrus.Infof("ERROR:PSQL Insert error %s", err.Error())
		return 0, err
	}
	logrus.Info("DONE: INSERTED Data PSQL")
	return id, nil
}

func (repo *AuthPostgres) IsVerified(id int, logrus *logrus.Logger) (int64, error) {

	tm := time.Now()
	query := fmt.Sprintf("	UPDATE %s SET is_verified = true,verification_date=$1,updated_at=$2	WHERE id = $3  RETURNING id ", usersTable)
	rows, err := repo.db.Exec(query, tm, tm, id)

	if err != nil {
		logrus.Errorf("ERROR: Update verificationCode failed : %v", err)
		return 0, err
	}
	effectedRowsNum, err := rows.RowsAffected()
	if err != nil {
		logrus.Errorf("ERROR: Update verificationCode effectedRowsNum failed : %v", err)
		return 0, err
	}
	logrus.Info("DONE:  email verified")
	return effectedRowsNum, nil
}

func (repo *AuthPostgres) GetUserID(email string, logrus *logrus.Logger) (int, error) {
	var id int
	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1 ", usersTable)
	err := repo.db.Get(&id, query, email)
	if err != nil {
		logrus.Errorf("ERROR: don't get users %s", err)
		return 0, errors.New("ERROR: user not found")
	}
	logrus.Info("DONE:get user data from psql")
	return id, nil
}

func (repo *AuthPostgres) GetUser(email, password string) (model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1 AND password=$2", usersTable)
	err := repo.db.Get(&user, query, email, password)

	return user, err
}
