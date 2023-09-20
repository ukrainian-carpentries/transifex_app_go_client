package transifex_app_go_client

import (
	"github.com/hemantasapkota/djangobot"
	"github.com/parnurzeal/gorequest"
	"github.com/sirupsen/logrus"
)

type TransifexAppClient struct {
	l     *logrus.Logger // An instance of the logrus logger
	agent *gorequest.SuperAgent
	bot   *djangobot.Bot
}

// The function returns a new instance of the transifex API client
// with the configured logger
func New(config *Config) (*TransifexAppClient, error) {

	// Create a transifex API client instance
	tr := &TransifexAppClient{
		l: logrus.New(), // create a logger instance
	}

	// Configure the logger
	tr.configureLogger(config)
	return tr, tr.Authenticate()
}
