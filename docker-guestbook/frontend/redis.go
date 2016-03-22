package main

import (
	"errors"
	"fmt"

	"github.com/garyburd/redigo/redis"
)

func getRedisConnection() (redis.Conn, error) {
	service := env.GetService(redisServiceInstance)
	if service == nil {
		return nil, errors.New("Service is nil")
	}

	address := fmt.Sprintf("%v:%v", service.Credentials["host"], service.Credentials["port"])
	conn, err := redis.Dial("tcp", address)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Cannot connect to Redis[%v]: %v", address, err))
	}

	// login to redis with credentials
	if _, err := conn.Do("AUTH", fmt.Sprintf("%v", service.Credentials["password"])); err != nil {
		conn.Close()
		return nil, errors.New(fmt.Sprintf("Redis authentication error: %v", err))
	}

	return conn, nil
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
