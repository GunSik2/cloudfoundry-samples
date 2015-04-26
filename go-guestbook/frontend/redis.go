package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/garyburd/redigo/redis"
)

func getRedisConnection() (redis.Conn, error) {
	service := env.GetService(redisServiceInstance)
	if service == nil {
		return nil, errors.New("Service is nil")
	}

	address := fmt.Sprintf("%v:%v", service.Credentials["hostname"], service.Credentials["port"])
	conn, err := redis.Dial("tcp", address)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Cannot connect to Redis[%v]: %v", address, err))
	}

	if _, err := conn.Do("AUTH", fmt.Sprintf("%v", service.Credentials["password"])); err != nil {
		conn.Close()
		return nil, errors.New(fmt.Sprintf("Redis authentication error: %v", err))
	}

	return conn, nil
}

type HitCounter struct{}

func (h *HitCounter) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	err := increaseHitCounter()
	if err != nil {
		fmt.Fprintf(w, "Redis error: %v", err)
		return
	}
	next(w, r)
}

func increaseHitCounter() error {
	c, err := getRedisConnection()
	if err != nil {
		return err
	}
	defer c.Close()

	_, err = c.Do("INCR", "hit-counter")
	if err != nil {
		return err
	}
	return nil
}

func getHitCounter() (int64, error) {
	c, err := getRedisConnection()
	if err != nil {
		return 0, err
	}
	defer c.Close()

	counter, err := redis.Int64(c.Do("GET", "hit-counter"))
	if err != nil {
		return 0, err
	}
	return counter, nil
}

func discoverBackends() ([]string, error) {
	c, err := getRedisConnection()
	if err != nil {
		return nil, err
	}
	defer c.Close()

	exists, err := redis.Bool(c.Do("EXISTS", "go-guestbook-backends"))
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}

	var backends []string
	reply, err := redis.Strings(c.Do("SMEMBERS", "go-guestbook-backends"))
	if err != nil {
		return backends, err
	}

	for _, r := range reply {
		exists, err := redis.Bool(c.Do(
			"EXISTS",
			fmt.Sprintf("go-guestbook-backend-%v", r)))
		if err != nil {
			return backends, err
		}

		if !exists {
			_, err := c.Do("SREM", "go-guestbook-backends", r)
			if err != nil {
				return backends, err
			}
		} else {
			backend, err := redis.String(c.Do(
				"GET",
				fmt.Sprintf("go-guestbook-backend-%v", r)))
			if err != nil {
				return backends, err
			}
			backends = append(backends, backend)
		}
	}

	return backends, nil
}
