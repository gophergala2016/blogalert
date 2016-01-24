package blogalert

import (
	"net/url"
	"time"
)

// Article defines article structure
type Article struct {
	Blog *Blog

	URL       *url.URL
	Title     string
	MD5       string
	Timestamp time.Time

	Flag Flag
}
