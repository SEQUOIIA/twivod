package db

import (
	"github.com/sequoiia/twiVod/server/twitch/model"
)

func (d * Dbc) AddUser(u model.User) error{
	if d.ps.addTwitchUser == nil {
		var err error
		d.ps.addTwitchUser, err = d.GetDbc().Prepare("INSERT INTO twitch_user (twitch_id, name, bio, created_at, display_name, logo, type, updated_at) VALUES (?,?,?,?,?,?,?,?)")
		if err != nil {
			return err
		}
	}

	_, err := d.ps.addTwitchUser.Exec(u.Id, u.Name, u.Bio, u.CreatedAt, u.DisplayName, u.Logo, u.Type, u.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}