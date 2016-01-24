package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gophergala2016/blogalert"
)

type ReadAllController struct {
	repo           blogalert.Repository
	tokenValidator *TokenValidator
}

func NewReadAllController(repo blogalert.Repository, tokenValidator *TokenValidator) *ReadAllController {
	return &ReadAllController{
		repo:           repo,
		tokenValidator: tokenValidator,
	}
}

func (rac *ReadAllController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var payload interface{}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	defer func() {
		e := json.NewEncoder(w)
		e.Encode(payload)
	}()

	token := r.FormValue("token")
	url := r.FormValue("url")

	valid, uid, err := rac.tokenValidator.ValidateToken(token)

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

	blog, err := rac.repo.GetBlog(url)

	err = rac.repo.SetUserBlogAsRead(uid, blog)
	if err != nil {
		payload = ErrorPayload{
			Error: err.Error(),
		}
		w.WriteHeader(500)
		return
	}

	payload = SuccessPayload{Success: true}
}
