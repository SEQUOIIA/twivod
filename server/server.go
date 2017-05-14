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
var TrackedStreams map[string]model.UserTracked

func main() {
	T = twitch.NewClient("whfl4lyxyzgp36d1el8v2yuyit0ge5")
	T.Debug()
	TrackedStreams = make(map[string]model.UserTracked)

	// Run in separate goroutine
	go func() {
		usersTotal := debugData()

		// Print
		for i := 0; i < len(usersTotal); i++ {
			user := usersTotal[i]
			fmt.Println(user.Name)

		}
	}()


	// Block primary goroutine
	select {

	}
}

func debugData() []model.User{
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
				log.Fatal(err)
			}
		}
		log.Println(fmt.Sprintf("Db: %s", twitchUsername))
		usersTotal = append(usersTotal, u)
	}

	// Get users that do not exist in DB

	users, err := T.GetUsersByUsername(usersToGet)
	if err != nil {
		log.Fatal(err)
	}

	// Add those users to the DB

	for i := 0; i < len(users); i++ {
		user := users[i]
		err = T.Db.AddUser(user)
		if err != nil {
			log.Fatal(err)
		}
	}

	usersTotal = append(usersTotal, users...)

	return usersTotal
}