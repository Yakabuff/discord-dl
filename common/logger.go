package common

import (
	"os"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func NewErrLogger() (*log.Logger, error) {
	file, err := os.OpenFile("errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	var log = logrus.New()
	log.Out = file
	return log, nil
}

func NewWebLogger() (*log.Logger, error) {
	file, err := os.OpenFile("web.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	var log = logrus.New()
	log.Out = file
	return log, nil
}
