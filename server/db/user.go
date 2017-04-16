package db

import (
	"github.com/sequoiia/twiVod/server/twitch/model"
	"database/sql"
	"errors"
	"strings"
)

var ErrNoRow error = errors.New("twitch db: No row found with the specified criterias")

func (d * Dbc) AddUser(u model.User) error {
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

func (d * Dbc) GetUserByTwitchID(twitchId string) (model.User, error) {
	if d.ps.getTwitchUser == nil {
		var err error
		d.ps.getTwitchUser, err = d.GetDbc().Prepare("SELECT twitch_id, name, bio, created_at, display_name, logo, type, updated_at from twitch_user where twitch_id=?")
		if err != nil {
			return model.User{}, err
		}
	}

	var u model.User
	result := d.ps.getTwitchUser.QueryRow(twitchId).Scan(&u.Id, &u.Name, &u.Bio, &u.CreatedAt, &u.DisplayName, &u.Logo, &u.Type, &u.UpdatedAt)
	if result != nil {
		if result == sql.ErrNoRows {
			return model.User{}, ErrNoRow
		} else {
			return model.User{}, result
		}
	}

	return u, nil
}

func (d * Dbc) GetUserByTwitchUsername(twitchUsername string) (model.User, error) {
	if d.ps.getTwitchUser == nil {
		var err error
		d.ps.getTwitchUser, err = d.GetDbc().Prepare("SELECT twitch_id, name, bio, created_at, display_name, logo, type, updated_at from twitch_user where name=?")
		if err != nil {
			return model.User{}, err
		}
	}

	var u model.User
	result := d.ps.getTwitchUser.QueryRow(strings.ToLower(twitchUsername))
	//.Scan(&u.Id, &u.Name, &u.Bio, &u.CreatedAt, &u.DisplayName, &u.Logo, &u.Type, &u.UpdatedAt)
	err := result.Scan(&u.Id)
	if err != nil {
		e := d.checkUserResultErr(err)
		return model.User{}, e
	}

	err = result.Scan(&u.Name)
	if err != nil {
		e := d.checkUserResultErr(err)
		return model.User{}, e
	}

	err = result.Scan(&u.Bio)
	if err != nil {
		e := d.checkUserResultErr(err)
		return model.User{}, e
	}

	err = result.Scan(&u.CreatedAt)
	if err != nil {
		e := d.checkUserResultErr(err)
		return model.User{}, e
	}

	err = result.Scan(&u.DisplayName)
	if err != nil {
		e := d.checkUserResultErr(err)
		return model.User{}, e
	}

	err = result.Scan(&u.Logo)
	if err != nil {
		e := d.checkUserResultErr(err)
		return model.User{}, e
	}

	err = result.Scan(&u.Type)
	if err != nil {
		e := d.checkUserResultErr(err)
		return model.User{}, e
	}

	err = result.Scan(&u.UpdatedAt)
	if err != nil {
		e := d.checkUserResultErr(err)
		return model.User{}, e
	}

	return u, nil
}

func (d * Dbc) checkUserResultErr(e error) error{
	switch e {
	case sql.ErrNoRows:
		return ErrNoRow
	default:
		return e
	}
}