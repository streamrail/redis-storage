# redis-storage [![Circle CI](https://circleci.com/gh/streamrail/redis-storage.svg?style=svg)](https://circleci.com/gh/streamrail/redis-storage) [![GoDoc](https://godoc.org/github.com/streamrail/redis-storage?status.svg)](https://godoc.org/github.com/streamrail/redis-storage)

## Summary
Small helper package around [Redigo](https://github.com/garyburd/redigo) that extracts some very common functionality that we have on all projects using `Redis`. 

## Benefits

- Easy value fetching by primitive types:

```go
r := NewRedisStorage(*redisHost, *redisConnPoolSize, *redisPrefix)
	err := r.Set("test", 123)
	if err != nil {
		t.Fatalf(err.Error())
	}
	res, err := r.GetInt32("test")
```

- Easy pub/sub

## License
MIT (see [LICENSE](https://github.com/streamrail/lrcache/blob/master/LICENSE) file)
