package blogalert

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/Sirupsen/logrus"
)

// Extractor crawls blogs
type Extractor struct {
	repo Repository
	log  *logrus.Logger
	wp   *WorkerPool
}

// NewExtractor creates new extractor
func NewExtractor(repo Repository, wp *WorkerPool, log *logrus.Logger) *Extractor {
	return &Extractor{
		repo: repo,
		wp:   wp,
		log:  log,
	}
}

func (e *Extractor) getBody(res *http.Response) (content string, hash string) {
	var body []byte
	if res.Body != nil {
		body, _ = ioutil.ReadAll(res.Body)
	}
	res.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return string(body), fmt.Sprintf("%x", md5.Sum(body))
}

func (e *Extractor) links(u *url.URL, doc *goquery.Document) ([]*url.URL, error) {
	urls := make([]*url.URL, 0, 5)

	doc.Find("a[href]").Each(func(i int, sel *goquery.Selection) {
		val, _ := sel.Attr("href")
		u, err := u.Parse(val)

		if err != nil {
			e.log.WithError(err).Errorf("Error resolving URL %s", val)
			return
		}

		u.Fragment = ""

		if u.Path != "" && u.Path[len(u.Path)-1:] == "/" {
			u.Path = u.Path[:len(u.Path)-1]
		}

		urls = append(urls, u)
	})

	return urls, nil
}

func (e *Extractor) crawl(ctx *context, u *url.URL, depth int) Worker {
	if !ctx.Queue(u) || depth > 5 {
		return nil
	}

	return Worker(func(wp *WorkerPool) {
		e.log.Infof("Requesting %s", u)

		res, err := http.Get(u.String())
		defer res.Body.Close()
		if err != nil {
			e.log.WithError(err).Errorf("Error requesting `%s`", u)
			return
		}

		_, hash := e.getBody(res)
		e.log.Infof("Page %s has hash %s", u, hash)

		doc, err := goquery.NewDocumentFromResponse(res)
		if err != nil {
			e.log.WithError(err).Errorf("Error proccessing URL %s", u)
			return
		}

		a, _ := e.repo.GetArticle(u.String())
		if a != nil && (a.MD5 == hash || a.Flag == Ignore) {
			return
		}

		links, err := e.links(res.Request.URL, doc)

		if err != nil {
			e.log.WithError(err).Errorf("Error parsing document `%s`", u)
			return
		}

		for _, link := range links {
			if link.Host == ctx.URL().Host {
				wp.Do(e.crawl(ctx, link, depth+1))
			}
		}

		e.log.Infof("Proccessing %s", u)

		title := doc.Find("head>title").First().Text()

		if title == "" {
			e.log.Infof("Page %s does not have title - ignored", u)
			return
		}

		article, err := ctx.Blog().NewArticle(u.String(), title, hash)

		if err != nil {
			e.log.WithError(err).Errorf("Error creating article %s", u)
			return
		}

		e.log.Infof("Storing %s [%s]", article.Title, article.URL)
		err = e.repo.InsertArticle(article)
		if err != nil {
			e.log.WithError(err).Errorf("Error storing article %s", article.URL)
			return
		}
	})
}

// Crawl crawls blog for articles
func (e *Extractor) Crawl(blog *Blog) {
	ctx := newContext(blog)
	e.wp.Do(e.crawl(ctx, ctx.URL(), 0))
}
