package blogalert

// Repository defines the blog repository
type Repository interface {
	// Get blog by URL, returns nil, nil on no value
	GetBlog(URL string) (*Blog, error)
	// Get all blogs
	GetAllBlogs() ([]*Blog, error)
	// Insert new blog
	InsertBlog(*Blog) error

	// Get Article by URL, returns nil, nil on no value
	GetArticle(URL string) (*Article, error)
	// Get all articles
	GetAllArticles() ([]*Article, error)
	// Insert new article
	InsertArticle(*Article) error

	// Get blogs a user is subscribed to
	GetUserSubscriptions(UID string) ([]*Blog, error)
	// Adds subscription to blog
	AddUserSubscription(UID string, blog *Blog) error
	// Deleres subscription to blog
	DeleteUserSubscription(UID string, blog *Blog) error
	// Get articles a user has read
	GetUserArticlesRead(UID string, blog *Blog) ([]*Article, error)
	// Set an article as read
	SetUserArticleAsRead(UID string, article *Article) error
}
