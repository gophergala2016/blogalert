package rethink

import (
	"sync"
	"time"

	"github.com/gophergala2016/blogalert"
)

type item struct {
	created time.Time
	data    interface{}
}

type cache struct {
	lock     sync.RWMutex
	articles map[string]*item
	blogs    map[string]*item
}

func newCache() *cache {
	c := &cache{}
	c.Clean()
	return c
}

func (c *cache) Clean() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.articles = make(map[string]*item)
	c.blogs = make(map[string]*item)
}

func (c *cache) GetArticle(url string) *blogalert.Article {
	c.lock.RLock()
	defer c.lock.RUnlock()
	item, ok := c.articles[url]
	if ok && time.Since(item.created) < time.Minute*5 {
		return item.data.(*blogalert.Article)
	}
	return nil
}

func (c *cache) SetArticle(article *blogalert.Article) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.articles[article.URL.String()] = &item{
		data:    article,
		created: time.Now(),
	}
}

func (c *cache) GetBlog(url string) *blogalert.Blog {
	c.lock.RLock()
	defer c.lock.RUnlock()
	item, ok := c.blogs[url]
	if ok && time.Since(item.created) < time.Minute*5 {
		return item.data.(*blogalert.Blog)
	}
	return nil
}

func (c *cache) SetBlog(blog *blogalert.Blog) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.blogs[blog.URL.String()] = &item{
		data:    blog,
		created: time.Now(),
	}
}
