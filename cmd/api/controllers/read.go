package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gophergala2016/blogalert"
)

type ReadController struct {
	repo           blogalert.Repository
	tokenValidator *TokenValidator
}

func NewReadController(repo blogalert.Repository, tokenValidator *TokenValidator) *ReadController {
	return &ReadController{
		repo:           repo,
		tokenValidator: tokenValidator,
	}
}

func (rc *ReadController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var payload interface{}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	defer func() {
		e := json.NewEncoder(w)
		e.Encode(payload)
	}()

	token := r.FormValue("token")
	url := r.FormValue("url")

	valid, uid, err := rc.tokenValidator.ValidateToken(token)

	if err != nil {
		payload = ErrorPayload{
			Error: err.Error(),
		}
		w.WriteHeader(500)
		return
	}

	if !valid {
		payload = ErrorPayload{
			Error: "Invalid token",
		}
		w.WriteHeader(401)
		return
	}

	article, err := rc.repo.GetArticle(url)
	if err != nil {
		payload = ErrorPayload{
			Error: err.Error(),
		}
		w.WriteHeader(500)
		return
	}

	if article == nil {
		payload = ErrorPayload{
			Error: "Article does not exist",
		}
		w.WriteHeader(500)
		return
	}

	err = rc.repo.SetUserArticleAsRead(uid, article)
	if err != nil {
		payload = ErrorPayload{
			Error: err.Error(),
		}
		w.WriteHeader(500)
		return
	}

	payload = SuccessPayload{Success: true}
}
