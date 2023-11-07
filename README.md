# url-shortner

URL shortner backend written in Go, redis cache and ratelimit.

Request model is, `short` and `expiry` is optional
```json
{
  "url": "http://your-url/",
  "short": "mycustomeurl",
  "expiry": 60
}
```

Endpoints
```
 POST /generate  - generate shorten url
 GET  /{short}   - redirect to original
```