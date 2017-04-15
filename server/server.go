package main

import (
	"github.com/sequoiia/twiVod/server/twitch"
	"log"
	"fmt"
)

var T *twitch.Twitch

func main() {
	T = twitch.NewClient("whfl4lyxyzgp36d1el8v2yuyit0ge5")
	T.Debug()

	// Run in separate goroutine
	go func() {
		users, err := T.GetUsersByUsername([]string{"nalcs1", "sequoiia"})
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < len(users); i++ {
			user := users[i]
			fmt.Println(user.Name)
		}
	}()


	// Block primary goroutine
	select {

	}
}