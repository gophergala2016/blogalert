package rethink

import (
	"net/url"

	"github.com/gophergala2016/blogalert"
)

type blog struct {
	URL   string `gorethink:"id"`
	Title string `gorethink:"title"`
}

type article struct {
	BlogURL string `gorethink:"blog"`
	URL     string `gorethink:"id"`
	Title   string `gorethink:"title"`
	MD5     string `gorethink:"md5"`
}

type subscription struct {
	UID     string `gorethink:"uid"`
	BlogURL string `gorethink:"blog"`
}

type articleRead struct {
	UID        string `gorethink:"uid"`
	BlogURL    string `gorethink:"blog"`
	ArticleURL string `gorethink:"article"`
}

func newBlog(b *blogalert.Blog) *blog {
	return &blog{
		URL:   b.URL.String(),
		Title: b.Title,
	}
}

func newArticle(b *blogalert.Article) *article {
	return &article{
		BlogURL: b.Blog.URL.String(),
		URL:     b.URL.String(),
		Title:   b.Title,
		MD5:     b.MD5,
	}
}

func newSubscription(uid string, b *blogalert.Blog) *subscription {
	return &subscription{
		UID:     uid,
		BlogURL: b.URL.String(),
	}
}

func newArticleRead(uid string, a *blogalert.Article) *articleRead {
	return &articleRead{
		UID:        uid,
		ArticleURL: a.URL.String(),
		BlogURL:    a.Blog.URL.String(),
	}
}

func (b *blog) ToBlog(repo *repo) (*blogalert.Blog, error) {
	u, err := url.Parse(b.URL)
	if err != nil {
		return nil, err
	}

	return &blogalert.Blog{
		URL:   u,
		Title: b.Title,
	}, nil
}

func (a *article) ToArticle(repo *repo) (*blogalert.Article, error) {
	u, err := url.Parse(a.URL)
	if err != nil {
		return nil, err
	}

	blog, err := repo.GetBlog(a.BlogURL)
	if err != nil {
		return nil, err
	}

	return &blogalert.Article{
		Blog:  blog,
		URL:   u,
		Title: a.Title,
		MD5:   a.MD5,
	}, nil
}
