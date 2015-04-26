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

func registerBackend() error {
	c, err := getRedisConnection()
	if err != nil {
		return err
	}
	defer c.Close()

	_, err = c.Do(
		"SETEX",
		fmt.Sprintf("go-guestbook-backend-%v", env.Application.InstanceID),
		"30",
		env.InstanceAddress)
	if err != nil {
		return err
	}

	_, err = c.Do("SADD", "go-guestbook-backends", env.Application.InstanceID)
	if err != nil {
		return err
	}

	return nil
}
