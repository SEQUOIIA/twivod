package main

import (
	"github.com/sequoiia/twiVod/server/twitch"
	"log"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sequoiia/twiVod/server/db"
	"github.com/sequoiia/twiVod/server/twitch/model"
)

var T *twitch.Twitch

func main() {
	T = twitch.NewClient("whfl4lyxyzgp36d1el8v2yuyit0ge5")
	T.Debug()

	// Run in separate goroutine
	go func() {
		var userRaw []string = []string{"nalcs1", "sequoiia"}
		var usersToGet []string = []string{}
		var usersTotal []model.User

		for i := 0; i < len(userRaw); i++ {
			twitchUsername := userRaw[i]
			u, err := T.Db.GetUserByTwitchUsername(twitchUsername)
			if err != nil {
				if err == db.ErrNoRow {
					log.Println(fmt.Sprintf("Remote: %s", twitchUsername))
					usersToGet = append(usersToGet, twitchUsername)
					continue
				} else {
					log.Println("Weird")
					log.Fatal(err)
				}
			}
			log.Println(fmt.Sprintf("Db: %s", twitchUsername))
			usersTotal = append(usersTotal, u)
		}

		users, err := T.GetUsersByUsername(usersToGet)
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < len(users); i++ {
			user := users[i]
			err = T.Db.AddUser(user)
			if err != nil {
				log.Fatal(err)
			}
		}

		usersTotal = append(usersTotal, users...)

		for i := 0; i < len(usersTotal); i++ {
			user := usersTotal[i]
			fmt.Println(user.Name)
			/*
			err = T.Db.AddUser(user)
			if err != nil {
				log.Fatal(err)
			}
			*/
		}
	}()


	// Block primary goroutine
	select {

	}
}