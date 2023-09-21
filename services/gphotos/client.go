package gphotos

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/abibby/google-photos-backup/app/models"
	"github.com/abibby/google-photos-backup/config"
	"github.com/abibby/google-photos-backup/database"
	"github.com/abibby/salusa/database/model"
)

type Client struct {
	user       *models.User
	httpClient http.Client
}

func NewClient(u *models.User) *Client {
	return &Client{
		user:       u,
		httpClient: *http.DefaultClient,
	}
}

func (c *Client) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	if c.user.ExpiresAt.Before(time.Now()) {
		resp, err := Token(&TokenRequest{RefreshToken: c.user.RefreshToken})
		if err != nil {
			return nil, fmt.Errorf("refreshing token: %w", err)
		}

		c.user.AccessToken = resp.AccessToken
		c.user.RefreshToken = resp.RefreshToken
		c.user.ExpiresAt = time.Now().Add(time.Second * time.Duration(resp.ExpiresIn))

		err = model.Save(database.DB, c.user)
		if err != nil {
			return nil, fmt.Errorf("refreshing token: %w", err)
		}
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.user.AccessToken))

	return req, nil
}

func (c *Client) Get(url string) (*http.Response, error) {
	req, err := c.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}
func (c *Client) Post(url string, body io.Reader) (*http.Response, error) {
	req, err := c.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		gerr := &GError{}
		err = json.NewDecoder(resp.Body).Decode(gerr)
		if err != nil {
			return nil, err
		}
		return nil, gerr
	}

	return resp, nil
}

type TokenRequest struct {
	Code         string
	RefreshToken string
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

func Token(req *TokenRequest) (*TokenResponse, error) {
	body := url.Values{}

	body.Add("client_id", config.GoogleClientID)
	body.Add("client_secret", config.GoogleClientSecret)
	if req.Code != "" {
		body.Add("code", req.Code)
		body.Add("grant_type", "authorization_code")
		body.Add("redirect_uri", "http://localhost:6900/gauth")
	}
	if req.RefreshToken != "" {
		body.Add("refresh_token", req.RefreshToken)
		body.Add("grant_type", "refresh_token")
	}

	r, err := http.PostForm("https://oauth2.googleapis.com/token", body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode < 200 || r.StatusCode > 299 {
		gerr := map[string]string{}
		err = json.NewDecoder(r.Body).Decode(&gerr)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s: %s", gerr["error"], gerr["error_description"])
	}

	tr := &TokenResponse{}
	err = json.NewDecoder(r.Body).Decode(tr)
	if err != nil {
		return nil, err
	}

	return tr, nil
}
