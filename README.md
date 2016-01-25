# blogalert

http://blogalert.adamtalbot.me/

Blogalert crawls blog sites and alerts you when a subscription has new content

It is devided into a crawler, an api and a frontend. All of these can be hosted seperatly as long as they share a database.

## TODO

* Unsubscribe (in progess)
* URL autocomplete
* Blacklist so people cant crawl sites like twitter (there is already a limit on crawling a site)
* Preproccess unread articles

## Wishlist

 * Content extraction using diff
 * Search tool to find intersting articles
 * Analytics to procuce a trending section
 * Notifications (email/push)

## How it works

### Spider
The spider uses a worker pool, when a blog is crawled the initial page is scanned. Following this all links (within the same domain) are queued to be scanned. This repeats up to a depth of 5 or a total of 200 pages, whatever happenes first.

The MD5 value of a page is stored as well, if this has not changed then the page is not proccessed. This is to stop articles being proccessed multiple times. This also means that if the initial page has not changed then no other pages are scanned. It is assumed new blog articles get featured on this page.

### Frontend

This is a very small html server serving a html tempalte with some js and css. 

### API

This handles requests from the frontend. Currently new articles are calculated on request, but in time the hope is to pre proccess this so the api can just call down to pre rendered data in the db. Pre rendering would also allow for sending notifications to users.



