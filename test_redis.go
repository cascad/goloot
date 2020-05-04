package main

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"reflect"
)

func test_redis_main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	if err != nil {
		log.Fatal(err)
	}
	val, err := client.RPush("123", 1).Result()
	log.Println(reflect.TypeOf(val), val, err)

	val, err = client.Incr("122").Result()
	log.Println(reflect.TypeOf(val), val, err)

	lst, err := client.LRange("123", 0, -1).Result()
	log.Println(reflect.TypeOf(lst), lst)

	//fmt.Println(client.Keys("*"))
	log.Println(fmt.Sprintf("stcounter_%s", "1112"))
	log.Println(errors.New("123").Error())

}
