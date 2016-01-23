package blogalert

import "net/url"

// Article defines article structure
type Article struct {
	Blog *Blog

	URL   *url.URL
	Title string
	MD5   string

	Flag Flag
}
