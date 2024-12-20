package auth

import (
	"github.com/gorilla/sessions"
	"github.com/kidusshun/ecom_bot/config"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const (
	key    = "randomstring"
	MaxAge = 86400 * 30
	IsProd = false
)

func NewAuth() {
	googleClientId := config.GoogleEnvs.GoogleClientID
	googleClientSecret := config.GoogleEnvs.GoogleClientSecret

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(MaxAge)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = IsProd

	gothic.Store = store
	goth.UseProviders(
		google.New(googleClientId, googleClientSecret, "http://localhost:8080/auth/google/callback"),
	)
}
