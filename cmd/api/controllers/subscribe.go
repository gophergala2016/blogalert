package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gophergala2016/blogalert"
)

type SubscribeController struct {
	repo           blogalert.Repository
	tokenValidator *TokenValidator
}

func NewSubscribeController(repo blogalert.Repository, tokenValidator *TokenValidator) *SubscribeController {
	return &SubscribeController{
		repo:           repo,
		tokenValidator: tokenValidator,
	}
}

func (sc *SubscribeController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var payload interface{}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	defer func() {
		e := json.NewEncoder(w)
		e.Encode(payload)
	}()

	token := r.FormValue("token")
	url := r.FormValue("url")
	title := r.FormValue("title")

	valid, uid, err := sc.tokenValidator.ValidateToken(token)

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

	blog, err := sc.repo.GetBlog(url)
	if err != nil {
		payload = ErrorPayload{
			Error: err.Error(),
		}
		w.WriteHeader(500)
		return
	}

	if blog != nil {
		sc.repo.AddUserSubscription(uid, blog)
		payload = SuccessPayload{Success: true}
		return
	}

	// Create it
	if title == "" {
		blog, err := blogalert.LoadBlog(url)
		if err != nil {
			payload = ErrorPayload{
				Error: err.Error(),
			}
			w.WriteHeader(500)
			return
		}

		payload = ConfimTitlePayload{
			Title: blog.Title,
			URL:   blog.URL.String(),
		}
		return
	}

	blog, err = blogalert.NewBlog(url, title)

	if err != nil {
		payload = ErrorPayload{
			Error: err.Error(),
		}
		w.WriteHeader(500)
		return
	}
	err = sc.repo.InsertBlog(blog)
	if err != nil {
		payload = ErrorPayload{
			Error: err.Error(),
		}
		w.WriteHeader(500)
		return
	}

	err = sc.repo.AddUserSubscription(uid, blog)
	if err != nil {
		payload = ErrorPayload{
			Error: err.Error(),
		}
		w.WriteHeader(500)
		return
	}

	payload = SuccessPayload{Success: true}
}
