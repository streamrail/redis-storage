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
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/garyburd/redigo/redis"
	"sync"
	"time"
)

var ErrNil = redis.ErrNil
var ErrPoolExhausted = redis.ErrPoolExhausted

// RedisStorage holds a redis pool and a prefix and is used to Set and Get items from Redis
type RedisStorage struct {
	pool   *redis.Pool
	prefix string

	// pubsub client
	conn redis.Conn
	redis.PubSubConn
	sync.Mutex
}

type Message struct {
	Type    string
	Channel string
	Data    string
}

// Get a new instance of RedisStorage
func NewRedisStorage(redisHost string, redisConnPoolSize int, redisPrefix string) *RedisStorage {
	pool := newRedisConnectionPool(redisHost, redisConnPoolSize)
	client := &RedisStorage{
		pool,
		redisPrefix,
		pool.Get(),
		redis.PubSubConn{pool.Get()},
		sync.Mutex{},
	}
	go func() {
		for {
			time.Sleep(200 * time.Millisecond)
			client.Lock()
			client.conn.Flush()
			client.Unlock()
		}
	}()
	return client
}

// Get an item from Redis
func (rs *RedisStorage) Get(key string) ([]uint8, error) {
	conn := rs.connection()
	defer conn.Close()
	data, err := redis.Bytes(conn.Do("GET", rs.prefix+key))
	if err == redis.ErrNil {
		return nil, redis.ErrNil
	} else if err != nil {
		return nil, err
	}
	return data, nil
}

// Get string value from cached item
func (rs *RedisStorage) GetString(key string) (string, error) {
	buffer, err := rs.Get(key)
	if err != nil {
		return "", err
	}
	result := new(string)
	dec := gob.NewDecoder(bytes.NewBuffer(buffer))
	if err := dec.Decode(result); err != nil {
		return "", err
	} else {
		return *result, nil
	}
}

// Get int32 value from cached item
func (rs *RedisStorage) GetInt32(key string) (int32, error) {
	buffer, err := rs.Get(key)
	if err != nil {
		return 0, err
	}
	result := new(int32)
	dec := gob.NewDecoder(bytes.NewBuffer(buffer))
	if err := dec.Decode(result); err != nil {
		return 0, err
	} else {
		return *result, nil
	}
}

// Get int64 value from cached item
func (rs *RedisStorage) GetInt64(key string) (int64, error) {
	buffer, err := rs.Get(key)
	if err != nil {
		return 0, err
	}
	result := new(int64)
	dec := gob.NewDecoder(bytes.NewBuffer(buffer))
	if err := dec.Decode(result); err != nil {
		return 0, err
	} else {
		return *result, nil
	}
}

// Get int value from cached item
func (rs *RedisStorage) GetInt(key string) (int, error) {
	buffer, err := rs.Get(key)
	if err != nil {
		return 0, err
	}
	result := new(int)
	dec := gob.NewDecoder(bytes.NewBuffer(buffer))
	if err := dec.Decode(result); err != nil {
		return 0, err
	} else {
		return *result, nil
	}
}

// Get float64 value from cached item
func (rs *RedisStorage) GetFloat64(key string) (float64, error) {
	buffer, err := rs.Get(key)
	if err != nil {
		return 0, err
	}
	result := new(float64)
	dec := gob.NewDecoder(bytes.NewBuffer(buffer))
	if err := dec.Decode(result); err != nil {
		return 0, err
	} else {
		return *result, nil
	}
}

// Set an item on Redis, with an expiration duration
func (rs *RedisStorage) Setex(key string, val interface{}, duration time.Duration) error {
	conn := rs.connection()
	defer conn.Close()
	var buffer = bytes.NewBuffer(nil)
	toStore := []byte{}
	if val != nil {
		enc := gob.NewEncoder(buffer)
		enc.Encode(val)
		toStore = buffer.Bytes()
	} else {
		toStore = nil
	}
	result, err := redis.String(conn.Do("SETEX", rs.prefix+key, int64(duration.Seconds()), toStore))
	if err != nil {
		return err
	}
	if result != "OK" {
		return errors.New("redis: SETEX call failed")
	}
	return nil
}

// Set an item on Redis
func (rs *RedisStorage) Set(key string, val interface{}) error {
	conn := rs.connection()
	defer conn.Close()
	var buffer = bytes.NewBuffer(nil)
	toStore := []byte{}
	if val != nil {
		enc := gob.NewEncoder(buffer)
		enc.Encode(val)
		toStore = buffer.Bytes()
	} else {
		toStore = nil
	}
	result, err := redis.String(conn.Do("SET", rs.prefix+key, toStore))
	if err != nil {
		return err
	}
	if result != "OK" {
		return errors.New("redis: SET call failed")
	}
	return nil
}

// Remove an item from Redis
func (rs *RedisStorage) Delete(key string) error {
	conn := rs.connection()
	defer conn.Close()
	result, err := redis.Int(conn.Do("DEL", rs.prefix+key))
	if err != nil {
		return err
	}
	if result != 1 {
		return errors.New("redis: DEL call failed for key (inc. prefix) " + rs.prefix + key)
	}
	return nil
}

func newRedisConnectionPool(host string, poolSize int) *redis.Pool {
	return redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", host)
		if err != nil {
			return nil, err
		}
		return c, err
	}, poolSize)
}

func (r *RedisStorage) connection() redis.Conn {
	return r.pool.Get()
}

func (r *RedisStorage) Pool() *redis.Pool {
	return r.pool
}

func (client *RedisStorage) Receive(psc redis.PubSubConn) *Message {
	switch message := psc.Receive().(type) {
	case redis.Message:
		return &Message{"message", message.Channel, string(message.Data)}
	}
	return nil
}

func (client *RedisStorage) NewPubSubConn() redis.PubSubConn {
	return redis.PubSubConn{client.Pool().Get()}
}

func (client *RedisStorage) Publish(channel, data string) {
	client.Lock()
	client.conn.Send("PUBLISH", channel, data)
	client.Unlock()
}
