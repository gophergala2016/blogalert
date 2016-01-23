package controllers

import "github.com/gophergala2016/blogalert"

type ArticlePayload struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type SubscriptionPayload struct {
	Title   string            `json:"title"`
	URL     string            `json:"url"`
	Updates []*ArticlePayload `json:"updates"`
}

type UpdatePayload struct {
	Subscriptions []*SubscriptionPayload `json:"subscriptions"`
	Updates       []*ArticlePayload      `json:"updates"`
}

type ErrorPayload struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type ConfimTitlePayload struct {
	Success bool   `json:"success"`
	Title   string `json:"title"`
	URL     string `json:"url"`
}

type SuccessPayload struct {
	Success bool `json:"success"`
}

func newArticlePayload(articles []*blogalert.Article) []*ArticlePayload {
	payload := make([]*ArticlePayload, 0, len(articles))
	for _, article := range articles {
		payload = append(payload, &ArticlePayload{
			Title: article.Title,
			URL:   article.URL.String(),
		})
	}

	return payload
}

func newSubscriptionPayload(blog *blogalert.Blog, updates []*blogalert.Article) *SubscriptionPayload {
	return &SubscriptionPayload{
		Title:   blog.Title,
		URL:     blog.URL.String(),
		Updates: newArticlePayload(updates),
	}
}
