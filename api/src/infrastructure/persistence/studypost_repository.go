package persistence

import (
	"database/sql"

	"github.com/code-wave/go-wave/domain/entity"
	"github.com/code-wave/go-wave/domain/repository"
	"github.com/code-wave/go-wave/infrastructure/errors"
	"github.com/code-wave/go-wave/infrastructure/helpers"
	"github.com/lib/pq"
)

type studyPostRepo struct {
	db *sql.DB
}

func NewStudyPostRepo(db *sql.DB) *studyPostRepo {
	return &studyPostRepo{db}
}

var _ repository.StudyPostRepository = &studyPostRepo{}

func (s *studyPostRepo) SavePost(studyPost *entity.StudyPost) (*entity.StudyPost, *errors.RestErr) {
	stmt, err := s.db.Prepare(`
		INSERT INTO study_post (user_id, title, topic, content, num_of_members, is_mentor, price, start_date, end_date, is_online, tech_stack, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id;
	`)
	if err != nil {
		return nil, errors.NewInternalServerError("database error" + err.Error())
	}

	currentTime := helpers.GetCurrentTimeForDB()

	row := stmt.QueryRow(studyPost.UserID, studyPost.Title, studyPost.Topic, studyPost.Content,
		studyPost.NumOfMembers, studyPost.IsMentor, studyPost.Price, studyPost.StartDate, studyPost.EndDate,
		studyPost.IsOnline, studyPost.TechStack, currentTime, currentTime)

	var lastInsertID int64

	err = row.Scan(&lastInsertID)
	if err != nil {
		return nil, errors.NewInternalServerError("row scan error " + err.Error())
	}

	studyPost.ID = lastInsertID

	return studyPost, nil
}

func (s *studyPostRepo) GetPost(id int64) (*entity.StudyPost, *errors.RestErr) {
	stmt, err := s.db.Prepare(`
		SELECT *
		FROM study_post
		WHERE id=$1;
	`)
	if err != nil {
		return nil, errors.NewInternalServerError("database error " + err.Error())
	}

	var studyPost entity.StudyPost

	err = stmt.QueryRow(id).Scan(&studyPost.ID, &studyPost.UserID, &studyPost.Title, &studyPost.Topic, &studyPost.Content, &studyPost.NumOfMembers,
		&studyPost.IsMentor, &studyPost.Price, &studyPost.StartDate, &studyPost.EndDate, &studyPost.IsOnline,
		pq.Array(&studyPost.TechStack), &studyPost.CreatedAt, &studyPost.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewBadRequestError("row doesn't exist " + err.Error())
		}
		return nil, errors.NewInternalServerError("database error" + err.Error())
	}

	return &studyPost, nil
}

func (s *studyPostRepo) GetPostsInLatestOrder(limit, offset int64) (entity.StudyPosts, *errors.RestErr) { // TODO: uint64 관련해서 js의 number는 64bit float형이라 데이터 받을때 string으로 받아야함
	stmt, err := s.db.Prepare(`
		SELECT *
		FROM study_post
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2;
	`)
	if err != nil {
		return nil, errors.NewInternalServerError("database error " + err.Error())
	}

	rows, err := stmt.Query(limit, offset)
	if err != nil {
		return nil, errors.NewInternalServerError("database error " + err.Error())
	}

	var studyPosts entity.StudyPosts

	for rows.Next() {
		var studyPost entity.StudyPost

		err := rows.Scan(&studyPost.ID, &studyPost.UserID, &studyPost.Title, &studyPost.Topic, &studyPost.Content, &studyPost.NumOfMembers,
			&studyPost.IsMentor, &studyPost.Price, &studyPost.StartDate, &studyPost.EndDate, &studyPost.IsOnline,
			pq.Array(&studyPost.TechStack), &studyPost.CreatedAt, &studyPost.UpdatedAt)
		if err != nil {
			return nil, errors.NewInternalServerError("database error " + err.Error())
		}
		studyPosts = append(studyPosts, studyPost)
	}

	if err = rows.Err(); err != nil { // 끝난 후에도 에러체크 한번
		return nil, errors.NewInternalServerError("database error " + err.Error())
	}

	return studyPosts, nil
}

// GetPostsByUserID 특정 user가 쓴 게시글들을 최신순으로 return
func (s *studyPostRepo) GetPostsByUserID(userID, limit, offset int64) (entity.StudyPosts, *errors.RestErr) {
	stmt, err := s.db.Prepare(`
		SELECT *
		FROM study_post
		WHERE user_id=$1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3;
	`)
	if err != nil {
		return nil, errors.NewInternalServerError("database error " + err.Error())
	}

	rows, err := stmt.Query(userID, limit, offset)
	if err != nil {
		return nil, errors.NewInternalServerError("database error " + err.Error())
	}

	var studyPosts entity.StudyPosts

	for rows.Next() {
		var studyPost entity.StudyPost

		err := rows.Scan(&studyPost.ID, &studyPost.UserID, &studyPost.Title, &studyPost.Topic, &studyPost.Content, &studyPost.NumOfMembers,
			&studyPost.IsMentor, &studyPost.Price, &studyPost.StartDate, &studyPost.EndDate, &studyPost.IsOnline,
			pq.Array(&studyPost.TechStack), &studyPost.CreatedAt, &studyPost.UpdatedAt)
		if err != nil {
			return nil, errors.NewInternalServerError("database error " + err.Error())
		}
		studyPosts = append(studyPosts, studyPost)
	}

	if err = rows.Err(); err != nil { // 끝난 후에도 에러체크 한번
		return nil, errors.NewInternalServerError("database error " + err.Error())
	}

	return studyPosts, nil
}

func (s *studyPostRepo) UpdatePost(studyPost *entity.StudyPost) (*entity.StudyPost, *errors.RestErr) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, errors.NewInternalServerError("database error " + err.Error())
	}
	stmt, err := tx.Prepare(`
		UPDATE study_post
		SET title=$1, topic=$2, content=$3, num_of_members=$4, is_mentor=$5, price=$6,
		    start_date=$7, end_date=$8, is_online=$9, tech_stack=$10, updated_at=$11
		WHERE id=$12
		RETURNING *;
	`)
	//stmt, err := s.db.Prepare(`
	//	UPDATE study_post
	//	SET title=$1, topic=$2, content=$3, num_of_members=$4, is_mentor=$5, price=$6,
	//	    start_date=$7, end_date=$8, is_online=$9, tech_stack=$10, updated_at=$11
	//	WHERE id=$12
	//	RETURNING *;
	//`)
	if err != nil {
		return nil, errors.NewInternalServerError("database error " + err.Error())
	}

	now := helpers.GetCurrentTimeForDB()
	row := stmt.QueryRow(studyPost.Title, studyPost.Topic, studyPost.Content, studyPost.NumOfMembers, studyPost.IsMentor, studyPost.Price, studyPost.StartDate,
		studyPost.EndDate, studyPost.IsOnline, pq.Array(studyPost.TechStack), now, studyPost.ID)
	err = row.Err()
	if err != nil {
		return nil, errors.NewInternalServerError("query row error " + err.Error())
	}

	err = row.Scan(&studyPost.ID, &studyPost.UserID, &studyPost.Title, &studyPost.Topic, &studyPost.Content, &studyPost.NumOfMembers,
		&studyPost.IsMentor, &studyPost.Price, &studyPost.StartDate, &studyPost.EndDate, &studyPost.IsOnline, pq.Array(&studyPost.TechStack), &studyPost.CreatedAt, &studyPost.UpdatedAt)
	if err != nil {
		rErr := tx.Rollback() // 에러시 rollback
		if rErr != nil {
			return nil, errors.NewInternalServerError("rollback error " + err.Error())
		}
		return nil, errors.NewInternalServerError("database update error " + err.Error())
	}

	return studyPost, nil
}

func (s *studyPostRepo) DeletePost(studyPostID int64) *errors.RestErr {
	stmt, err := s.db.Prepare(`
		DELETE 
		FROM study_post
		WHERE id=$1;
	`)
	if err != nil {
		return errors.NewInternalServerError("database error " + err.Error())
	}

	res, err := stmt.Exec(studyPostID)
	if err != nil {
		return errors.NewInternalServerError("database error " + err.Error())
	}

	n, err := res.RowsAffected()
	if err != nil {
		return errors.NewInternalServerError("database error " + err.Error())
	}

	if n == 0 {
		return errors.NewBadRequestError("no rows to be deleted")
	}

	return nil
}
