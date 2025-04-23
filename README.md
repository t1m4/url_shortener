# URL shortener

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
```bash
curl -v "http://localhost:8000/api/short_url" -d '{"url": "https://gorm.io/docs/error_handling.html"}' -H "Content-Type: application/json"
```

- Redirect url
```bash
curl -v http://localhost:8000/api/x3UlUaD
```

# Branches
- master - url shortener realization with go and postgres
- app_replicas_docker - extended master branch with nginx as reversed proxy and several replicas of app

# Rate limiter
- Using rate library to create multi limiter struct for API or any other resourse or process allowing to limit accoss different units of time
- Create middleware and service to check user limit
- Create backgroup goroutine that will delete expired rate limiter from map using sync.Mutex 

## Current problems
- Background cleaning task can take a lot of time the there is a lot of users
- Storing in memory can cause several issues: 
    - there is not enough memory for a lot of user
    - data lose after restarts
    - duplication of rate limiter in distributed systems, if reverse proxy doesn't route using some hash function.

Solution: Try in-memory db(like redis) to fix this problems 

## Tests
How much memory userRateLimiterByUserId map takes? 
How log cleaning time for whole map?
  - 0 users. memory 30-35mb, average clean - 1ms    
  - 1000 users. memory 30-35mb, average clean - 2-3ms
  - 10000 users. memory 30-35mb, average clean - 4-6ms
  - 100000 users. memory 55-65mb, average clean - 20-28ms
  - 200000 users. memory 85-95mb, average clean - 40-50ms
  - 400000 users. memory 140-150mb, average clean - 75-80ms
  - 500000 users. memory 170-150mb, average clean - 85-90ms
  - 1000000 users. memory 315-325mb, average clean - 180-200ms
Each 100k users adds 20-30mb of memory and 20-25ms of cleaning time
So for 1 million it will be 200-300mb of memory and 200-300ms of cleaning time
So for 10 million it will be 2-3gb of memory and 2-3s of cleaning time


# Stop signals 
Inside docker only works using 
```
CMD ["./main"]
```
Or installing tini
```
# Add tini as init system
RUN apk add --no-cache tini
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["./main"]
```
