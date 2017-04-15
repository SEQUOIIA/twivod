package twitch

import "fmt"

type User struct {

}

// Endpoint: /kraken/users
// Documentation URL: https://dev.twitch.tv/docs/v5/guides/using-the-twitch-api/#translating-from-user-names-to-user-ids
// Query parameters:
//   login : comma separated usernames
func (t * Twitch) GetUsersByUsername(users []string) ([]User, error) {
	req, err := t.newRequest("GET", fmt.Sprintf("%s%s", OfficialTwitchAPIEndpoint, "/kraken/users"), nil)
	if err != nil {
		return nil, err
	}

	
}