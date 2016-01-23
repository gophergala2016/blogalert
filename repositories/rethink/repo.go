package rethink

import (
	"github.com/dancannon/gorethink"
	"github.com/gophergala2016/blogalert"
)

//Table names
const (
	Database          string = "blogalert"
	ArticleTable      string = "article"
	BlogTable         string = "blogs"
	SubscriptionTable string = "subscriptions"
	ArticleReadTable  string = "articles_read"
)

type repo struct {
	session *gorethink.Session
}

// NewRepo creates new repo
func NewRepo(config *Config) (blogalert.Repository, error) {
	sess, err := gorethink.Connect(config.getConnOps())
	if err != nil {
		return nil, err
	}

	return &repo{
		session: sess,
	}, nil
}

// NewRepoFromSession creates new repo
func NewRepoFromSession(sess *gorethink.Session) blogalert.Repository {
	return &repo{
		session: sess,
	}
}

func (r *repo) GetBlog(URL string) (*blogalert.Blog, error) {
	cursor, err := gorethink.DB(Database).Table(BlogTable).
		Get(URL).Run(r.session)
	if err != nil {
		return nil, err
	}

	b := &blog{}

	err = cursor.One(b)
	if err != nil {
		return nil, err
	}

	return b.ToBlog(r)
}

func (r *repo) GetAllBlogs() ([]*blogalert.Blog, error) {
	cursor, err := gorethink.DB(Database).Table(BlogTable).Run(r.session)
	if err != nil {
		return nil, err
	}

	b := []*blog{}

	err = cursor.All(b)
	if err != nil {
		return nil, err
	}

	blogs := make([]*blogalert.Blog, 0, len(b))

	for _, v := range b {
		if blog, err := v.ToBlog(r); err == nil {
			blogs = append(blogs, blog)
		}
	}

	return blogs, nil
}

func (r *repo) InsertBlog(b *blogalert.Blog) error {
	_, err := gorethink.DB(Database).Table(BlogTable).
		Insert(newBlog(b)).RunWrite(r.session)
	return err
}

func (r *repo) GetArticle(URL string) (*blogalert.Article, error) {
	cursor, err := gorethink.DB(Database).Table(ArticleTable).
		Get(URL).Run(r.session)
	if err != nil {
		return nil, err
	}

	a := &article{}

	err = cursor.One(a)
	if err != nil {
		return nil, err
	}

	return a.ToArticle(r)
}

func (r *repo) GetAllArticles() ([]*blogalert.Article, error) {
	cursor, err := gorethink.DB(Database).Table(ArticleTable).Run(r.session)
	if err != nil {
		return nil, err
	}

	a := []*article{}

	err = cursor.All(a)
	if err != nil {
		return nil, err
	}

	articles := make([]*blogalert.Article, 0, len(a))

	for _, v := range a {
		if article, err := v.ToArticle(r); err == nil {
			articles = append(articles, article)
		}
	}

	return articles, nil
}

func (r *repo) InsertArticle(a *blogalert.Article) error {
	_, err := gorethink.DB(Database).Table(ArticleTable).
		Insert(newArticle(a)).RunWrite(r.session)
	return err
}

func (r *repo) GetUserSubscriptions(UID string) ([]*blogalert.Blog, error) {
	cursor, err := gorethink.DB(Database).Table(SubscriptionTable).
		Filter(gorethink.Row.Field("uid").Eq(UID)).Run(r.session)
	if err != nil {
		return nil, err
	}

	s := []*subscription{}

	err = cursor.All(s)
	if err != nil {
		return nil, err
	}

	blogs := make([]*blogalert.Blog, 0, len(s))

	for _, v := range s {
		if blog, err := r.GetBlog(v.BlogURL); err == nil {
			blogs = append(blogs, blog)
		}
	}

	return blogs, nil
}

func (r *repo) AddUserSubscription(UID string, blog *blogalert.Blog) error {
	_, err := gorethink.DB(Database).Table(SubscriptionTable).
		Insert(newSubscription(UID, blog)).RunWrite(r.session)
	return err
}

func (r *repo) DeleteUserSubscription(UID string, blog *blogalert.Blog) error {
	_, err := gorethink.DB(Database).Table(SubscriptionTable).
		Filter(gorethink.Row.Field("uid").Eq(UID)).
		Filter(gorethink.Row.Field("blog").Eq(blog.URL.String())).
		Delete().RunWrite(r.session)
	return err
}

func (r *repo) GetUserArticlesRead(UID string, blog *blogalert.Blog) ([]*blogalert.Article, error) {
	cursor, err := gorethink.DB(Database).Table(ArticleReadTable).
		Filter(gorethink.Row.Field("uid").Eq(UID)).
		Filter(gorethink.Row.Field("blog").Eq(blog.URL.String())).
		Run(r.session)
	if err != nil {
		return nil, err
	}

	a := []*articleRead{}

	err = cursor.All(a)
	if err != nil {
		return nil, err
	}

	articles := make([]*blogalert.Article, 0, len(a))

	for _, v := range a {
		if article, err := r.GetArticle(v.ArticleURL); err == nil {
			articles = append(articles, article)
		}
	}

	return articles, nil
}

func (r *repo) SetUserArticleAsRead(UID string, article *blogalert.Article) error {
	_, err := gorethink.DB(Database).Table(ArticleReadTable).
		Insert(newArticleRead(UID, article)).RunWrite(r.session)
	return err
}
