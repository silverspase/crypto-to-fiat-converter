# Crypto to Fiat converter

GRPC service written in Golang which converts crypto amount to fiat(currency) value

## Features
- Config, logger and Dependency injection
- GRPC as transport layer
- In-memory store with partial invalidation, check out Cache section for details
- Clean architecture(I haven't implement dedicated structs for each layer, but the most interchangeable layers has its own structures)
- Unit test example(`intenral/service/cache/memory/memory_test.go`)
- Linter
- Docker setup
- Integration test example for `Convert` RPC 
- Github actions pipeline


## Getting started

1. Generate proto files:
```
make proto-gen
```
2. Fetch dependencies:
```
go mod tidy
```
3. Run service
a. As a golang app(you should have go installed)
```
go build ./cmd/main.go && ./main
```
b. As Docker container:
```
make docker-build && make docker-run
```

### To send request you should have GRPC client.
Find proto contracts in `./proto` directory

To convert crypto into fiat:
```
{
  "request": [
    {
      "fromToken": "bitcoin",
      "toCurrency": "usd",
      "amount": 1
    }
  ]
}
```

To get tokens list:
```
{
  "pageToken": 0,
  "pageSize": 5
}
```

## Cache

### Solution 1

We can store everything, and we don't need to handle expiration, we just update all data periodically. 
And implementation is relatively simple - in case of given task we can collect all data, and replace it at once with old one.
But this solution has two constrains:
- Storage size can be limited.
- Cache update for large amount of data can be an issue, we need to figure out the update mechanism, because we cannot lock everything during update

### Solution 2

We can define most frequently used token list and currencies and load them in memory on app start and keep updating it periodically.
Additionally, we can store data that is absent in cache on user request. 

During the request, there is no need to wait for Upsert(Insert/Update) completion, we can send the result to user immediately and finish the Upsert in gorutine

I implemented this solution, check the details below

#### Cache invalidation

All cache entries has expiration window and on request we check if data is expired, 
if yes, we request actual data from Price provider and store it in cache. 
This is valid for not predefined tokens and currencies, predefined one will be updated automatically.

##### Bottleneck

Sooner or later we have a risk of being run out of cache storage, because of data that was requested only once and stored in cache forever.
We can implement `Cleaner service` which will remove outdated data from cache.
We can use Queue for that, each cache entry goes also to Queue like this:
```
type QueueEntry struct {
	Token string
	Currency string
	ExirationDate time.Time
}
```
Cleaner service will go over the queue as long as Expiration date is in the past and removes from map data.

#### Predefined tokens and currencies

We can start from predefined list of tokens and currencies, 
but we can also collect requests metrics and cache actual most frequent tokens and currencies

```
type metrics struct { // nolint: unused
	tokens     map[string]int
	currencies map[string]int
}
```


## Pagination

Pagination is implemented only for `GetTokenList`.

## TODO
1. Retry-Backoff logic for Price Provider
2. Metrics collection for convert requests and composing most frequent tokens and currencies based on actual data
3. Cache invalidation for not predefined tokens and currencies

