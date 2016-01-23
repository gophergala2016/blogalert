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
}
