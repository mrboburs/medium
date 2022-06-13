package repository

import (
	"fmt"
	"mediumuz/model"
	"mediumuz/util/logrus"
	// "strings"
	// "time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type PostDB struct {
	db *sqlx.DB
}

func NewPostDB(db *sqlx.DB) *PostDB {
	return &PostDB{db: db}

}

func (repo *PostDB) GetAllPosts(pagination model.Pagination, logrus *logrus.Logger) (posts []model.PostFull, err error) {
	query := fmt.Sprintf("SELECT id , post_title ,post_image_path, post_body, post_views_count, post_like_count, post_rated, post_vote, post_tags,  post_date, is_new, is_top_read FROM  post WHERE deleted_at IS NULL ORDER BY  post_date DESC OFFSET $1 LIMIT $2 ")

	err = repo.db.Select(&posts, query, pagination.Offset, pagination.Limit)
	if err != nil {
		logrus.Errorf("ERROR: don't get users %s", err)
		return posts, err
	}
	logrus.Info("DONE:get user data from psql")
	return posts, err
}
func (repo *PostDB) CreatePost(post model.Post, logrus *logrus.Logger) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (post_title  , post_body , post_tags) VALUES ($1, $2, $3)  RETURNING id", postTable)

	row := repo.db.QueryRow(query, post.Title, post.Body, pq.Array(post.Tags))

	if err := row.Scan(&id); err != nil {
		logrus.Infof("ERROR:PSQL Insert error %s", err.Error())
		return 0, err
	}
	logrus.Info("DONE: INSERTED Data PSQL")
	return id, nil
}

func (repo *PostDB) CreatePostUser(userId, postId int, logrus *logrus.Logger) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (post_author_id , post_id ) VALUES ($1, $2)  RETURNING id", postUserTable)
	row := repo.db.QueryRow(query, userId, postId)
	if err := row.Scan(&id); err != nil {
		logrus.Infof("ERROR:PSQL Insert error %s", err.Error())
		return 0, err
	}
	logrus.Info("DONE: INSERTED Data PSQL")
	return id, nil
}
func (repo *PostDB) CommentPost(input model.CommentPost, logrus *logrus.Logger) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (reader_id,post_id,comment) VALUES ($1, $2, $3)  RETURNING id", "comment_post")

	row := repo.db.QueryRow(query, input.ReaderID, input.PostID, input.PostComment)

	if err := row.Scan(&id); err != nil {
		logrus.Infof("ERROR:PSQL Insert error %s", err.Error())
		return 0, err
	}
	logrus.Info("DONE: INSERTED Data PSQL")
	return id, nil
}

func (repo *PostDB) GetPostById(postId int, logrus *logrus.Logger) (post model.PostFull, err error) {

	query := fmt.Sprintf("SELECT id , post_title ,post_image_path, post_body, post_views_count, post_like_count, post_rated, post_vote, post_tags,  post_date, is_new, is_top_read FROM %s WHERE id = $1 AND deleted_at IS NULL", postTable)
	// convert.StringToNullString(post.PostImagePath)
	err = repo.db.Get(&post, query, postId)
	if err != nil {
		logrus.Errorf("ERROR: didn't get postData %s", err)
		return post, err
	}
	logrus.Info("DONE:get user data from psql")
	return post, err
}

func (repo *PostDB) CheckPostId(id int, logrus *logrus.Logger) (int, error) {
	var postNumber int
	queryID := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE id=$1 AND deleted_at IS NULL", postTable)
	err := repo.db.Get(&postNumber, queryID, id)
	if err != nil {
		logrus.Infof("ERROR:Email query error: %s", err.Error())
		return -1, err
	}
	return postNumber, nil
}
func (repo *PostDB) ViewPost(userID, postID int, logrus *logrus.Logger) (int, error) {

	var id int
	query := fmt.Sprintf("INSERT INTO %s (reader_id  , post_id ) VALUES ($1, $2)  RETURNING id", "viewed_post")

	row := repo.db.QueryRow(query, userID, postID)

	if err := row.Scan(&id); err != nil {
		logrus.Infof("ERROR:PSQL Insert VIEW error %s", err.Error())
		return 0, err
	}
	logrus.Info("DONE: INSERTED  VIEW Data PSQL")
	return id, nil
}
func (repo *PostDB) ClickLike(userId int, postId int, logrus *logrus.Logger) error {
	likeQuery := fmt.Sprintln("SELECT toggle_comment_like($1,$2)")

	row, err := repo.db.Exec(likeQuery, userId, postId)

	if err != nil {
		logrus.Info("DONE: ERROR  LIKE Data PSQL %s ", err)
		return err
	}
	_, err = row.RowsAffected()
	if err != nil {
		logrus.Info("DONE: ERROR  LIKE Data PSQL %s ", err)

		return err
	}
	logrus.Info("DONE: INSERTED  LIKE Data PSQL")
	return nil
}

func (repo *PostDB) UpdatePostImage(userId, postId int, input model.PostFull, filePath string, logrus *logrus.Logger) error {

	// tm := time.Now()
	query := fmt.Sprintf("	UPDATE post AS p SET post_image_path=$1  FROM post_user AS pu  WHERE p.id = pu.post_id and p.id=$2 ")
	rows, err := repo.db.Exec(query, filePath, postId)

	if err != nil {
		logrus.Errorf("ERROR: Update photo failed : %v", err)
		return err
	}
	_, err = rows.RowsAffected()
	if err != nil {
		logrus.Errorf("ERROR: Update  photo failed: %v", err)
		return err
	}
	logrus.Info("DONE:Updated photo successfully")
	return nil

}
func (repo *PostDB) UpdatePost(id int, input model.Post, logrus *logrus.Logger) error {

	// tm := time.Now()
	query := fmt.Sprintf("	UPDATE post AS p SET post_title = COALESCE(NULLIF($1,''),post_title) , post_body = COALESCE(NULLIF($2,''),post_body), post_tags = COALESCE(NULLIF($3,array[null]),post_tags)  FROM post_user AS pu  WHERE p.id = pu.post_id and p.id=$4 ")
	rows, err := repo.db.Exec(query, input.Title, input.Body, pq.Array(input.Tags), id)

	if err != nil {
		logrus.Errorf("ERROR: Update  failed : %v", err)
		return err
	}
	_, err = rows.RowsAffected()
	if err != nil {
		logrus.Errorf("ERROR: Update   failed: %v", err)
		return err
	}
	logrus.Info("DONE:Updated  successfully")
	return nil

}
func (repo *PostDB) PostDelete(userId, postId int) error {
	query := fmt.Sprintf("DELETE FROM %s p USING %s pu WHERE p.id = pu.post_id AND pu.post_author_id=$1 AND pu.post_id=$2",
		postTable, postUserTable)
	_, err := repo.db.Exec(query, userId, postId)

	return err
}

func (repo *PostDB) GetCommentsPost(postID int, pagination model.Pagination, logrus *logrus.Logger) (comments []model.CommentFull, err error) {
	query := fmt.Sprintf("SELECT 	cmt.id,		u.first_name, u.last_name,	u.account_image_path,cmt.comment	 FROM %s u INNER JOIN %s cmt on u.id =cmt.reader_id WHERE cmt.post_id = $1 AND cmt.deleted_at IS NULL OFFSET $2 LIMIT $3",
		usersTable, "comment_post")
	err = repo.db.Select(&comments, query, postID, pagination.Offset, pagination.Limit)
	return comments, err

}
