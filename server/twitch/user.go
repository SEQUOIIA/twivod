package twitch

import (
	"fmt"
	"encoding/json"
	"github.com/sequoiia/twivod/server/twitch/model"
	"strings"
)


type userByUsernameResult struct {
	Id * string `json:"_id"`
	Bio * string `json:"bio"`
	CreatedAt * string `json:"created_at"`
	DisplayName * string `json:"display_name"`
	Logo * string `json:"logo"`
	Name * string `json:"name"`
	Type * string `json:"type"`
	UpdatedAt * string `json:"updated_at"`
}

type getUsersByUsernameResult struct {
	Total int `json:"_total"`
	Users []userByUsernameResult `json:"users"`
}

const debugM string = "twiVod.twitch.Twitch.GetUsersByUsername ->"

// Endpoint: /kraken/users
// Documentation URL: https://dev.twitch.tv/docs/v5/guides/using-the-twitch-api/#translating-from-user-names-to-user-ids
// Query parameters:
//   login : comma separated usernames
func (t * Twitch) GetUsersByUsername(users []string) ([]model.User, error) {
	req, err := t.newRequest("GET", fmt.Sprintf("%s%s", OfficialTwitchAPIEndpoint, "/kraken/users"), nil)
	if err != nil {
		return nil, err
	}

	queryParameters := req.URL.Query()

	// Iterate through all requested users and add them to the query parameter
	var usersFormatted string = ""
	for i := 0; i < len(users); i++ {
		usersFormatted = fmt.Sprintf("%s%s", usersFormatted, users[i])

		if (i + 1) != len(users) {
			usersFormatted = fmt.Sprintf("%s,", usersFormatted)
		}
	}

	if strings.Compare("", usersFormatted) != 0 {
		t.debugPrint(debugM, usersFormatted)
	}

	queryParameters.Add("login", usersFormatted)
	req.URL.RawQuery = queryParameters.Encode()

	resp, err := t.httpCli.Do(req)
	if err != nil {
		return nil, err
	}

	var result getUsersByUsernameResult
	var payload []model.User
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(result.Users); i++ {
		r := result.Users[i]
		u := model.User{}
		u.Name = t.nilCheckerString(u.Name, r.Name)
		u.Bio = t.nilCheckerString(u.Bio, r.Bio)
		u.CreatedAt = t.nilCheckerString(u.CreatedAt, r.CreatedAt)
		u.DisplayName = t.nilCheckerString(u.DisplayName, r.DisplayName)
		u.Id = t.nilCheckerString(u.Id, r.Id)
		u.Logo = t.nilCheckerString(u.Logo, r.Logo)
		u.Type = t.nilCheckerString(u.Type, r.Type)
		u.UpdatedAt = t.nilCheckerString(u.UpdatedAt, r.UpdatedAt)
		payload = append(payload, u)

	}
	return payload, nil
}

func (t * Twitch) nilCheckerString(a string, b * string) string{
	if b != nil {
		t.debugPrint(debugM, "b does not equal nil")
		a = *b
	} else {
		t.debugPrint(debugM, "b equals nil")
		a = ""
	}

	return a
}