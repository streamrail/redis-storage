// Copyright (c) 2015 streamrail

// The MIT License (MIT)

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package redistorage

import (
	"flag"
	"testing"
)

var (
	redisHostBad      = flag.String("redis_bad", "loc22alhost:6379", "Redis host and port. Eg: localhost:6379")
	redisHost         = flag.String("redis", "localhost:6379", "Redis host and port. Eg: localhost:6379")
	redisConnPoolSize = flag.Int("redisConnPoolSize", 5, "Redis connection pool size. Default: 5")
	redisPrefix       = flag.String("redisPrefix", "rl_", "Redis prefix to attach to keys")
	lruSize           = flag.Int("lruSize", 5, "LRU Cache Size. Default: 5")
)

func TestInt(t *testing.T) {
	r := NewRedisStorage(*redisHost, *redisConnPoolSize, *redisPrefix)
	err := r.Set("test", 123)
	if err != nil {
		t.Fatalf(err.Error())
	}
	res, err := r.GetInt("test")
	if err != nil {
		t.Fatalf(err.Error())
	} else {
		if res != 123 {
			t.Fatalf("expected result to be 123")
		}
	}
}

func TestInt32(t *testing.T) {
	r := NewRedisStorage(*redisHost, *redisConnPoolSize, *redisPrefix)
	err := r.Set("test", 123)
	if err != nil {
		t.Fatalf(err.Error())
	}
	res, err := r.GetInt32("test")
	if err != nil {
		t.Fatalf(err.Error())
	} else {
		if res != 123 {
			t.Fatalf("expected result to be 123")
		}
	}
}
func TestInt64(t *testing.T) {
	r := NewRedisStorage(*redisHost, *redisConnPoolSize, *redisPrefix)
	err := r.Set("test", 123)
	if err != nil {
		t.Fatalf(err.Error())
	}
	res, err := r.GetInt64("test")
	if err != nil {
		t.Fatalf(err.Error())
	} else {
		if res != 123 {
			t.Fatalf("expected result to be 123")
		}
	}
}
func TestString(t *testing.T) {
	r := NewRedisStorage(*redisHost, *redisConnPoolSize, *redisPrefix)
	err := r.Set("test", "123")
	if err != nil {
		t.Fatalf(err.Error())
	}
	res, err := r.GetString("test")
	if err != nil {
		t.Fatalf(err.Error())
	} else {
		if res != "123" {
			t.Fatalf("expected result to be '123'")
		}
	}
}

func TestFloat64(t *testing.T) {
	r := NewRedisStorage(*redisHost, *redisConnPoolSize, *redisPrefix)
	err := r.Set("test", 1.5)
	if err != nil {
		t.Fatalf(err.Error())
	}
	res, err := r.GetFloat64("test")
	if err != nil {
		t.Fatalf(err.Error())
	} else {
		if res != 1.5 {
			t.Fatalf("expected result to be 1.5")
		}
	}
}

func TestBadConnection(t *testing.T) {
	r := NewRedisStorage(*redisHostBad, *redisConnPoolSize, *redisPrefix)
	err := r.Set("test", 123)
	if err == nil {
		t.Fatalf("exepcted an error when setting a value with bad redis host")
	}
}
