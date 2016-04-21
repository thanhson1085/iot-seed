package main

import (
	"fmt"
	"gopkg.in/redis.v3"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "192.168.1.191:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pubsub, err := client.Subscribe("mongdb_oblog")
	if err != nil {
		panic(err)
	}

	msg, err := pubsub.ReceiveMessage()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(msg.Channel, msg.Payload)
	defer pubsub.Close()
}
