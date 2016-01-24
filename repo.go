package blogalert

// Repository defines the blog repository
type Repository interface {
	// Get blog by URL, no result returns nil, nil
	GetBlog(URL string) (*Blog, error)
	// Get all blogs
	GetAllBlogs() ([]*Blog, error)
	// Insert new blog
	InsertBlog(*Blog) error

	// Get Article by URL, no result returns nil, nil
	GetArticle(URL string) (*Article, error)
	// Get all articles in blog
	GetAllArticlesInBlog(*Blog) ([]*Article, error)
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
	// Set an blog as read
	SetUserBlogAsRead(UID string, blog *Blog) error
}
