package main

import (
	"github.com/sequoiia/twiVod/server/twitch"
	"log"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
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

		stmt, err := T.Db.GetDb().Prepare("INSERT INTO twitch_user (twitch_id, name, bio, created_at, display_name, logo, type, updated_at) VALUES (?,?,?,?,?,?,?,?)")
		for i := 0; i < len(users); i++ {
			user := users[i]
			fmt.Println(user.Name)
			_, err := stmt.Exec(user.Id, user.Name, user.Bio, user.CreatedAt, user.DisplayName, user.Logo, user.Type, user.UpdatedAt)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()


	// Block primary goroutine
	select {

	}
}