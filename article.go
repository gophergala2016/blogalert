package blogalert

// Article defines article structure
type Article struct {
	Blog *Blog

	URL   string
	Title string
	MD5   string
}
