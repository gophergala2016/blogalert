package blogalert

import (
	"net/url"
	"sync"
)

type context struct {
	lock sync.Mutex
	dup  map[string]struct{}
	blog *Blog
	url  *url.URL
}

func newContext(blog *Blog) *context {
	u, _ := url.Parse(blog.URL)
	return &context{
		blog: blog,
		dup:  make(map[string]struct{}),
		url:  u,
	}
}

func (ctx *context) Queue(u *url.URL) bool {
	ctx.lock.Lock()
	defer ctx.lock.Unlock()

	_, ok := ctx.dup[u.String()]
	if ok {
		return false
	}

	ctx.dup[u.String()] = struct{}{}
	return true
}

func (ctx *context) Blog() *Blog {
	return ctx.blog
}

func (ctx *context) URL() *url.URL {
	return ctx.url
}
