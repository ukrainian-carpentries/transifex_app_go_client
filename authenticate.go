package transifex_app_go_client

import (
	"errors"
	"os"

	"github.com/hemantasapkota/djangobot"
)

func (t *TransifexAppClient) Authenticate() error {
	var err error

	username := os.Getenv("TRANSIFEX_USER")
	if len(username) == 0 {
		return errors.New("the environmental variable TRANSIFEX_USER is not set")
	}

	password := os.Getenv("TRANSIFEX_PASSWORD")
	if len(password) == 0 {
		return errors.New("the environmental variable TRANSIFEX_PASSWORD is not set")
	}

	t.bot = djangobot.With("https://app.transifex.com/signin/?next=/home/").
		ForHost("app.transifex.com").
		SetUsername(username).
		SetPassword(password).
		LoadCookies()
	if t.bot.Error != nil {
		return t.bot.Error
	}

	t.agent, err = t.bot.Set("next", "/home/").
		X("csrfmiddlewaretoken", t.bot.Cookie("csrftoken").Value).
		X("identification", t.bot.Username).
		X("password", t.bot.Password).
		X("remember_me", "True").
		X("next", "/home/").
		Login()
	if err != nil {
		return err
	}

	sessionid := t.bot.Cookie("sessionid").Value
	if sessionid == "" {
		return errors.New("authentication failed")
	}

	return nil
}
