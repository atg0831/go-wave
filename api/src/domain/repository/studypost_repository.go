package repository

import (
	"github.com/code-wave/go-wave/domain/entity"
	"github.com/code-wave/go-wave/infrastructure/errors"
)

type StudyPostRepository interface {
	SavePost(post *entity.StudyPost) *errors.RestErr
	GetPost(id uint64) (*entity.StudyPost, *errors.RestErr)
	GetPostsInLatestOrder(limit, offset uint64) (entity.StudyPosts, *errors.RestErr)
	UpdatePost(post *entity.StudyPost) (*entity.StudyPost, *errors.RestErr)
	DeletePost(post *entity.StudyPost) *errors.RestErr
}
