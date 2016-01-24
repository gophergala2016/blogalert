package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gophergala2016/blogalert"
)

type UpdateController struct {
	repo           blogalert.Repository
	tokenValidator *TokenValidator
}

func NewUpdateController(repo blogalert.Repository, tokenValidator *TokenValidator) *UpdateController {
	return &UpdateController{
		repo:           repo,
		tokenValidator: tokenValidator,
	}
}

func (u *UpdateController) Diff(read, total []*blogalert.Article) []*blogalert.Article {

	new := make([]*blogalert.Article, 0, len(total))
	for _, t := range total {
		valid := true
		for _, r := range read {
			if r == nil {
				continue
			}

			if t.URL.String() == r.URL.String() {
				valid = false
				break
			}
		}
		if valid {
			new = append(new, t)
		}
	}
	return new
}

func (u *UpdateController) GetPayload(uid string) (interface{}, int) {
	blogs, err := u.repo.GetUserSubscriptions(uid)
	if err != nil {
		return ErrorPayload{
			Error: err.Error(),
		}, 500
	}

	payload := &UpdatePayload{}

	for _, blog := range blogs {
		total, err := u.repo.GetAllArticlesInBlog(blog)
		if err != nil {
			return ErrorPayload{
				Error: err.Error(),
			}, 500
		}

		read, err := u.repo.GetUserArticlesRead(uid, blog)
		if err != nil {
			return ErrorPayload{
				Error: err.Error(),
			}, 500
		}

		new := u.Diff(read, total)

		payload.Subscriptions = append(payload.Subscriptions, newSubscriptionPayload(blog, new))
	}

	for _, sub := range payload.Subscriptions {
		payload.Updates = append(payload.Updates, sub.Updates...)
	}

	return payload, 200
}

func (u *UpdateController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var payload interface{}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	defer func() {
		e := json.NewEncoder(w)
		e.Encode(payload)
	}()

	token := r.FormValue("token")
	valid, uid, err := u.tokenValidator.ValidateToken(token)

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

	payload, code := u.GetPayload(uid)
	w.WriteHeader(code)
}
