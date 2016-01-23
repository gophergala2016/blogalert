package blogalert

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// Blog defines blog structure
type Blog struct {
	URL   string
	Title string
}

// LoadBlog creates a blog by parsing the page title of its URL
func LoadBlog(URL string) (*Blog, error) {
	res, err := http.Get(URL)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, err
	}

	title := doc.Find("title").First().Text()

	if title == "" {
		return nil, fmt.Errorf("Invalid blog %s - no title", URL)
	}

	return NewBlog(URL, title), nil
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
