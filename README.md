# URL shortener

## Branches
- master - url shortener realization with go and postgres
- app_replicas_docker - extended master branch with nginx as reversed proxy and several replicas of app

## Requirements
- Shorneting urls. Endpoint for creating short url.
    - Check url for availability. Create API client to do it.
    - Get or create row in db using transaction.
    - Return ready new short url from service layer
- Url redirection. Endpoint for redirecting short url to original one.
    - Get or update given short url count.
    - Redirect ot original url from db.
- Analytics. Store number of redirections.
    - Add count field to table

## Request examples

- Create short url
```
curl -v "http://localhost:8000/api/short_url" -d '{"url": "https://gorm.io/docs/error_handling.html"}' -H "Content-Type: application/json"
```

- Redirect url
```
curl -v http://localhost:8000/api/x3UlUaD
```
