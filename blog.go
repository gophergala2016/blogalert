package blogalert

// Blog defines blog structure
type Blog struct {
	URL   string
	Title string
}

// NewBlog creates a new blog
func NewBlog(url, title string) *Blog {
	return &Blog{
		URL:   url,
		Title: title,
	}
}

// NewArticle creates a new article in a blog
func (b *Blog) NewArticle(url, title, hash string) *Article {
	return &Article{
		Blog:  b,
		URL:   url,
		Title: title,
		MD5:   hash,
	}
}
