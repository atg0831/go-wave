package interfaces

import (
	"encoding/json"
	"fmt"
	"github.com/code-wave/go-wave/application"
	"github.com/code-wave/go-wave/domain/entity"
	"github.com/code-wave/go-wave/infrastructure/errors"
	"github.com/code-wave/go-wave/infrastructure/helpers"
	"net/http"
)

type StudyPost struct {
	sp application.StudyPostInterface
	//ts application.StudyPostTechStackInterface
}

func NewStudyPostHandler(sp application.StudyPostInterface) *StudyPost {
	return &StudyPost{
		sp: sp,
		//ts: ts,
	}
}

func (s *StudyPost) GetPost(w http.ResponseWriter, r *http.Request) {
	helpers.SetJsonHeader(w)

	studyPostID, restErr := helpers.ExtractIntParam(r, "study_post_id")
	if restErr != nil {
		w.WriteHeader(restErr.Status)
		w.Write(restErr.ResponseJSON().([]byte))
		return
	}

	studyPost, restErr := s.sp.GetPost(studyPostID)
	if restErr != nil {
		w.WriteHeader(restErr.Status)
		w.Write(restErr.ResponseJSON().([]byte))
		return
	}

	sJson, restErr := studyPost.ResponseJSON()
	if restErr != nil {
		w.WriteHeader(restErr.Status)
		w.Write(restErr.ResponseJSON().([]byte))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(sJson)
}

func (s *StudyPost) GetPostsInLatestOrder(w http.ResponseWriter, r *http.Request) {
	helpers.SetJsonHeader(w)

	limit, err := helpers.ExtractIntParam(r, "limit")
	if err != nil {
		w.WriteHeader(err.Status)
		w.Write(err.ResponseJSON().([]byte))
		return
	}

	offset, err := helpers.ExtractIntParam(r, "offset")
	if err != nil {
		w.WriteHeader(err.Status)
		w.Write(err.ResponseJSON().([]byte))
		return
	}

	studyPosts, err := s.sp.GetPostsInLatestOrder(limit, offset)
	if err != nil {
		w.WriteHeader(err.Status)
		w.Write(err.ResponseJSON().([]byte))
		return
	}

	sJson, err := studyPosts.ResponseJSON()
	if err != nil {
		w.WriteHeader(err.Status)
		w.Write(err.ResponseJSON().([]byte))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(sJson)
}

func (s *StudyPost) SavePost(w http.ResponseWriter, r *http.Request) {
	helpers.SetJsonHeader(w)

	var studyPost entity.StudyPost

	err := json.NewDecoder(r.Body).Decode(&studyPost)
	if err != nil {
		restErr := errors.NewBadRequestError("invalid json body") // TODO: 에러명 나중에 수정 (아래 에러들도)
		w.WriteHeader(restErr.Status)
		w.Write(restErr.ResponseJSON().([]byte))
		return
	}

	fmt.Println("studypost: ", studyPost) // check용 나중에 삭제

	restErr := studyPost.Validate()
	if restErr != nil {
		w.WriteHeader(restErr.Status)
		w.Write(restErr.ResponseJSON().([]byte))
		return
	}

	restErr = s.sp.SavePost(&studyPost)
	if restErr != nil {
		w.WriteHeader(restErr.Status)
		w.Write(restErr.ResponseJSON().([]byte))
		return
	}

	//restErr = s.ts.SaveStudyPostTechStack(studyPost.ID, studyPost.TechStack) // (studyPostID, techStackID) 형태로 저장
	//if restErr != nil {
	//	w.WriteHeader(restErr.Status)
	//	w.Write(restErr.ResponseJSON().([]byte))
	//	return
	//}

	err = json.NewEncoder(w).Encode(&studyPost) // TODO: marshal로 할까? 나중에 ResponseJSON 만들어서 할까?
	if err != nil {
		http.Error(w, "encode error", 500)
	}
}
