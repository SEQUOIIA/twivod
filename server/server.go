package main

import (
	"github.com/sequoiia/twiVod/server/twitch"
	"log"
	"fmt"
)

func main() {
	t := twitch.NewClient("whfl4lyxyzgp36d1el8v2yuyit0ge5")
	t.Debug()

	users, err := t.GetUsersByUsername([]string{"nalcs1", "sequoiia"})
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(users); i++ {
		user := users[i]
		fmt.Println(user.Name)
	}

	select {

	}
}