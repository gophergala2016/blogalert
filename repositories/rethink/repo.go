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
	cache   *cache
}

// NewRepo creates new repo
func NewRepo(config *Config) (blogalert.Repository, error) {
	sess, err := gorethink.Connect(config.getConnOps())
	if err != nil {
		return nil, err
	}

	return NewRepoFromSession(sess), nil
}

// NewRepoFromSession creates new repo
func NewRepoFromSession(sess *gorethink.Session) blogalert.Repository {
	gorethink.DBCreate(Database).RunWrite(sess)
	gorethink.DB(Database).TableCreate(ArticleTable).RunWrite(sess)
	gorethink.DB(Database).TableCreate(BlogTable).RunWrite(sess)
	gorethink.DB(Database).TableCreate(SubscriptionTable).RunWrite(sess)
	gorethink.DB(Database).TableCreate(ArticleReadTable).RunWrite(sess)
	gorethink.DB(Database).Table(ArticleTable).IndexCreate("blog")
	gorethink.DB(Database).Table(ArticleTable).IndexCreate("ts")
	gorethink.DB(Database).Table(SubscriptionTable).IndexCreate("uid")
	gorethink.DB(Database).Table(SubscriptionTable).IndexCreate("blog")
	gorethink.DB(Database).Table(ArticleReadTable).IndexCreate("uid")
	gorethink.DB(Database).Table(ArticleReadTable).IndexCreate("blog")

	return &repo{
		session: sess,
		cache:   newCache(),
	}
}

func (r *repo) GetBlog(URL string) (*blogalert.Blog, error) {
	if b := r.cache.GetBlog(URL); b != nil {
		return b, nil
	}

	cursor, err := gorethink.DB(Database).Table(BlogTable).
		Get(URL).Run(r.session)
	if err != nil {
		return nil, err
	}

	b := &blog{}

	err = cursor.One(b)

	if err == gorethink.ErrEmptyResult {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	blog, err := b.ToBlog()

	if err != nil {
		r.cache.SetBlog(blog)
	}

	return blog, err
}

func (r *repo) GetAllBlogs() ([]*blogalert.Blog, error) {
	cursor, err := gorethink.DB(Database).Table(BlogTable).Run(r.session)
	if err != nil {
		return nil, err
	}

	b := []*blog{}

	err = cursor.All(&b)
	if err != nil {
		return nil, err
	}

	blogs := make([]*blogalert.Blog, 0, len(b))

	for _, v := range b {
		if blog, err := v.ToBlog(); err == nil {
			blogs = append(blogs, blog)
			r.cache.SetBlog(blog)
		}
	}

	return blogs, nil
}

func (r *repo) InsertBlog(b *blogalert.Blog) error {
	_, err := gorethink.DB(Database).Table(BlogTable).
		Insert(newBlog(b)).RunWrite(r.session)

	r.cache.SetBlog(b)
	return err
}

func (r *repo) GetArticle(URL string) (*blogalert.Article, error) {
	if a := r.cache.GetArticle(URL); a != nil {
		return a, nil
	}

	cursor, err := gorethink.DB(Database).Table(ArticleTable).
		Get(URL).Run(r.session)
	if err != nil {
		return nil, err
	}

	a := &article{}

	err = cursor.One(a)

	if err == gorethink.ErrEmptyResult {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	blog, err := r.GetBlog(a.BlogURL)

	if err != nil {
		return nil, err
	}

	article, err := a.ToArticle(blog)

	if err != nil {
		r.cache.SetArticle(article)
	}

	return article, err
}

func (r *repo) GetAllArticlesInBlog(blog *blogalert.Blog) ([]*blogalert.Article, error) {
	if blog == nil {
		return nil, nil
	}

	cursor, err := gorethink.DB(Database).Table(ArticleTable).
		Filter(gorethink.Row.Field("blog").Eq(blog.URL.String())).
		OrderBy(gorethink.OrderByOpts{Index: gorethink.Desc("ts")}).
		Limit(100).
		Run(r.session)

	if err != nil {
		return nil, err
	}

	a := []*article{}

	err = cursor.All(&a)
	if err != nil {
		return nil, err
	}

	articles := make([]*blogalert.Article, 0, len(a))

	for _, v := range a {
		if article, err := v.ToArticle(blog); err == nil {
			articles = append(articles, article)
			r.cache.SetArticle(article)
		}
	}

	return articles, nil
}

func (r *repo) InsertArticle(a *blogalert.Article) error {
	_, err := gorethink.DB(Database).Table(ArticleTable).
		Insert(newArticle(a), gorethink.InsertOpts{Conflict: "replace"}).RunWrite(r.session)
	r.cache.SetArticle(a)
	return err
}

func (r *repo) GetUserSubscriptions(UID string) ([]*blogalert.Blog, error) {
	cursor, err := gorethink.DB(Database).Table(SubscriptionTable).
		Filter(gorethink.Row.Field("uid").Eq(UID)).Run(r.session)
	if err != nil {
		return nil, err
	}

	s := []*subscription{}

	err = cursor.All(&s)
	if err != nil {
		return nil, err
	}

	blogs := make([]*blogalert.Blog, 0, len(s))

	for _, v := range s {
		if blog, err := r.GetBlog(v.BlogURL); err == nil && blog != nil {
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
	if blog == nil {
		return nil, nil
	}

	cursor, err := gorethink.DB(Database).Table(ArticleReadTable).
		Filter(gorethink.Row.Field("uid").Eq(UID)).
		Filter(gorethink.Row.Field("blog").Eq(blog.URL.String())).
		Run(r.session)
	if err != nil {
		return nil, err
	}

	a := []*articleRead{}

	err = cursor.All(&a)
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

func (r *repo) SetUserBlogAsRead(UID string, blog *blogalert.Blog) error {
	articles, err := r.GetAllArticlesInBlog(blog)
	if err != nil {
		return err
	}
	rows := make([]*articleRead, 0, len(articles))
	for _, article := range articles {
		rows = append(rows, newArticleRead(UID, article))
	}

	_, err = gorethink.DB(Database).Table(ArticleReadTable).
		Insert(rows).RunWrite(r.session)
	return err

}
