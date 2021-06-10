package interfaces

import (
	"encoding/json"
	"fmt"
	"github.com/code-wave/go-wave/application"
	"github.com/code-wave/go-wave/domain/entity"
	"github.com/code-wave/go-wave/infrastructure/helpers"
	"net/http"
)

type StudyPost struct {
	sp application.StudyPostInterface
}

func NewStudyPost(sp application.StudyPostInterface) *StudyPost {
	return &StudyPost{
		sp,
	}
}

func (s *StudyPost) SavePost(w http.ResponseWriter, r *http.Request) {
	var studyPost entity.StudyPost

	err := json.NewDecoder(r.Body).Decode(&studyPost)
	if err != nil {
		http.Error(w, "decode error", 500) // TODO: 에러명 나중에 수정 (아래 marshal error도)
	}

	fmt.Println(studyPost) // check용 나중에 삭제

	validateErr := studyPost.Validate() // TODO: 그냥 map[string]string말고 err로 할까?
	if len(validateErr) > 0 {
		jsonData, err := json.Marshal(validateErr)
		if err != nil {
			http.Error(w, "marshal error", 500)
			return
		}
		helpers.JsonHeader(w)
		_, err = w.Write(jsonData)
		if err != nil {
			http.Error(w, "write error", 500)
			return
		}
	}

	w.Write([]byte("success"))

}
