[Back](../actions.md)

# Web Actions

## Safari

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

## URLs

### `expandURL(url)`

Expand `url`. This is generally used for short urls.

---

### `getURLComponent(url,component)`

Get `component` from `url`.

#### Components

- Scheme
- User
- Password
- Host
- Port
- Path
- Query
- Fragment

### `getURLs(input)`

Get urls from `input`.

---

### `url(...url)`

Create url value of `url`. No limit on `url` arguments.

## Web Requests

### `downloadURL(url,method)`

Download `url` using HTTP method `method`.

_This action is currently incomplete due to it's complexity_

---

### `getWebpageContents(url)`

Get contents of webpage at `url`.
