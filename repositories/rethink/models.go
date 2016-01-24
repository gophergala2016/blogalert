package rethink

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gophergala2016/blogalert"
)

type blog struct {
	URL   string `gorethink:"id"`
	Title string `gorethink:"title"`
}

type article struct {
	BlogURL   string `gorethink:"blog"`
	URL       string `gorethink:"id"`
	Title     string `gorethink:"title"`
	MD5       string `gorethink:"md5"`
	Timestamp int64  `gorethink:"ts"`
}

type subscription struct {
	ID      string `gorethink:"id"`
	UID     string `gorethink:"uid"`
	BlogURL string `gorethink:"blog"`
}

type articleRead struct {
	ID         string `gorethink:"id"`
	UID        string `gorethink:"uid"`
	BlogURL    string `gorethink:"blog"`
	ArticleURL string `gorethink:"article"`
	Timestamp  int64  `gorethink:"ts"`
}

func newBlog(b *blogalert.Blog) *blog {
	return &blog{
		URL:   b.URL.String(),
		Title: b.Title,
	}
}

func newArticle(b *blogalert.Article) *article {
	return &article{
		Timestamp: b.Timestamp.Unix(),
		BlogURL:   b.Blog.URL.String(),
		URL:       b.URL.String(),
		Title:     b.Title,
		MD5:       b.MD5,
	}
}

func newSubscription(uid string, b *blogalert.Blog) *subscription {
	return &subscription{
		ID:      fmt.Sprintf("%s_%s", uid, b.URL.String()),
		UID:     uid,
		BlogURL: b.URL.String(),
	}
}

func newArticleRead(uid string, a *blogalert.Article) *articleRead {
	return &articleRead{
		ID:         fmt.Sprintf("%s_%s", uid, a.URL.String()),
		UID:        uid,
		ArticleURL: a.URL.String(),
		BlogURL:    a.Blog.URL.String(),
		Timestamp:  a.Timestamp.Unix(),
	}
}

func (b *blog) ToBlog() (*blogalert.Blog, error) {
	u, err := url.Parse(b.URL)
	if err != nil {
		return nil, err
	}

	return &blogalert.Blog{
		URL:   u,
		Title: b.Title,
	}, nil
}

func (a *article) ToArticle(blog *blogalert.Blog) (*blogalert.Article, error) {
	u, err := url.Parse(a.URL)
	if err != nil {
		return nil, err
	}

	return &blogalert.Article{
		Blog:      blog,
		URL:       u,
		Title:     a.Title,
		MD5:       a.MD5,
		Timestamp: time.Unix(a.Timestamp, 0),
	}, nil
}
