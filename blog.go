package blogalert

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Blog defines blog structure
type Blog struct {
	URL   *url.URL
	Title string
}

// LoadBlog creates a blog by parsing the page title of its URL
func LoadBlog(address string) (*Blog, error) {
	u, err := url.Parse(address)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "" {
		u.Scheme = "http"
	}

	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, err
	}

	title := doc.Find("title").First().Text()

	if title == "" {
		return nil, fmt.Errorf("Invalid blog %s - no title", u)
	}

	return NewBlog(u.String(), title)
}

// NewBlog creates a new blog
func NewBlog(address, title string) (*Blog, error) {
	u, err := url.Parse(address)
	if err != nil {
		return nil, err
	}

	if u.Path != "" && u.Path[len(u.Path)-1:] == "/" {
		u.Path = u.Path[:len(u.Path)-1]
	}

	u.Scheme = "http"

	return &Blog{
		URL:   u,
		Title: title,
	}, nil
}

// NewArticle creates a new article in a blog
func (b *Blog) NewArticle(address, title, hash string) (*Article, error) {
	u, err := url.Parse(address)
	if err != nil {
		return nil, err
	}

	if u.Path != "" && u.Path[len(u.Path)-1:] == "/" {
		u.Path = u.Path[:len(u.Path)-1]
	}

	return &Article{
		Blog:      b,
		URL:       u,
		Title:     title,
		MD5:       hash,
		Timestamp: time.Now(),
	}, nil
}
