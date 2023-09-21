package gphotos

import (
	"encoding/json"
)

type Profile struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	VerifyEmail bool   `json:"verify_email"`
	Picture     string `json:"picture"`
}

func (c *Client) GetProfile() (*Profile, error) {
	resp, err := c.Get("https://www.googleapis.com/oauth2/v1/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r := &Profile{}
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return nil, err
	}

	return r, nil
}
