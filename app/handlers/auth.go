package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/abibby/google-photos-backup/app/models"
	"github.com/abibby/google-photos-backup/database"
	"github.com/abibby/google-photos-backup/services/gphotos"
	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/request"
)

type LoginRequest struct {
	Code    string        `query:"code"`
	Request *http.Request `json:"-"`
}

// var Login = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 	b, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		panic(fmt.Errorf("read all: %w", err))
// 	}
// 	r.Body.Close()

// 	params, err := url.ParseQuery(string(b))
// 	if err != nil {
// 		panic(fmt.Errorf("parse: %w", err))
// 	}

// 	spew.Dump(params)
// 	// tx := request.UseTx(r)
// 	// u := &models.User{
// 	// 	Credentials: params["credential"][0],
// 	// }
// 	// err = model.Save(tx, u)
// 	// if err != nil {
// 	// 	panic(fmt.Errorf("save: %w", err))
// 	// }

// 	// w.Header().Add("Location", "/")
// 	// w.WriteHeader(http.StatusSeeOther)
// })

var Login = request.Handler(func(r *LoginRequest) (any, error) {
	t, err := gphotos.Token(&gphotos.TokenRequest{Code: r.Code})
	if err != nil {
		return nil, fmt.Errorf("token: %w", err)
	}

	tx := request.UseTx(r.Request)
	u := &models.User{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Second * time.Duration(t.ExpiresIn)),
	}

	p, err := gphotos.NewClient(u).GetProfile()
	if err != nil {
		return nil, fmt.Errorf("profile: %w", err)
	}

	u.Email = p.Email

	existingUser, err := models.UserQuery(r.Request.Context()).Where("email", "=", u.Email).First(database.DB)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		existingUser.AccessToken = u.AccessToken
		if u.RefreshToken != "" {
			existingUser.RefreshToken = u.RefreshToken
		}
		existingUser.ExpiresAt = u.ExpiresAt
		u = existingUser
	}

	err = model.Save(tx, u)
	if err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}

	return request.NewHTMLResponse([]byte(fmt.Sprintf(`<a href="/">home</a><br>%#v`, u))), nil
})

//  func() (any, error) {
// 	log.Print(r.Credentials)
// 	return request.NewResponse(http.NoBody).
// 		SetStatus(http.StatusSeeOther).
// 		AddHeader("Location", "/"), nil
// })
