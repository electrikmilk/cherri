[Back](../actions.md)

# Web Actions

## Articles

### `getArticle(webpage)`

Get article from `webpage` in article.

---

### `getArticleDetail(article,detail)`

Get `detail` from `article`.

## Giphy

### `searchGiphy(query)`

Search Giphy for `query`.

---

### `getGifs(query,gifs)`

Get `gifs` number of gifs from Giphy for `query`.

## RSS

### `getRSS(items,url)`

Get `items` number of items from RSS feed at `url`.

---

### `getRSSFeeds(urls)`

Get RSS feeds from urls.

## Safari

### `addToReadingList(url)`

Add `url` to the users reading list.

---

### `getCurrentURL()`

Get current web page url.

---

### `getWebPageDetail(webpage,detail)`

Get `detail` of `webpage`.

#### Details

- Page Contents
- Page Selection
- Page URL
- Name

---

### `openURLs(urls)`

Open `urls` in default browser.

---

### `runJavaScriptOnWebpage(script)`

Run `script` on currently active tab as JavaScript.

---

### `searchWeb(engine,query)`

Search the web for `query` using `engine`.

#### Engines

- Amazon
- Bing
- DuckDuckGo
- eBay
- Google
- Reddit
- Twitter
- Yahoo!
- YouTube

---

### `showWebpage(url,useReader)`

Show `url`. `urlReader` is a boolean, the default is `true`. 

## URLs

### `expandURL(url)`

Expand `url`. This is generally used for short urls.

---

### `getURLDetail(url,detail)`

Get `detail` from `url`.

#### Details

- Scheme
- User
- Password
- Host
- Port
- Path
- Query
- Fragment

---

### `getURLs(input)`

Get urls from `input`.

---

### `url(...url)`

Create url value of `url`. No limit on `url` arguments.

## Web Requests

### `downloadURL(url)`

Download contents of `url`.

---

### `httpRequest(url,method,body,bodyType,headers)`

Download `url` using HTTP method `method`.

`body`, `bodyType` and `headers` are optional.

#### Body Types

- Form (default)
- JSON
- File

---

### `getWebpageContents(url)`

Get contents of webpage at `url`.

---

### `getURLHeaders(url)`

Get headers for `url`.
