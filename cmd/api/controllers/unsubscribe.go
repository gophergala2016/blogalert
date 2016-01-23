package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gophergala2016/blogalert"
)

type UnsubscribeController struct {
	repo           blogalert.Repository
	tokenValidator *TokenValidator
}

func NewUnsubscribeController(repo blogalert.Repository, tokenValidator *TokenValidator) *UnsubscribeController {
	return &UnsubscribeController{
		repo:           repo,
		tokenValidator: tokenValidator,
	}
}

func (usc *UnsubscribeController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var payload interface{}

	defer func() {
		w.Header().Set("Content-Type", "application/json")
		e := json.NewEncoder(w)
		e.Encode(payload)
	}()

	token := r.FormValue("token")
	url := r.FormValue("url")

	valid, uid, err := usc.tokenValidator.ValidateToken(token)

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

	blog, err := usc.repo.GetBlog(url)

	if err != nil {
		payload = ErrorPayload{
			Error: err.Error(),
		}
		w.WriteHeader(500)
		return
	}

	if blog != nil {
		err := usc.repo.DeleteUserSubscription(uid, blog)
		if err != nil {
			payload = ErrorPayload{
				Error: err.Error(),
			}
			w.WriteHeader(500)
			return
		}
	}

	payload = SuccessPayload{Success: true}
}
