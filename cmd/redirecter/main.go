package main

import (
	"github.com/IFtech-A/urlshortener/internal/redirect"
	"github.com/sirupsen/logrus"
)

func main() {

	logrus.SetLevel(logrus.DebugLevel)
	conf := redirect.NewConfig()

	serv := redirect.New(conf)

	err := serv.Start()
	if err != nil {
		logrus.Fatal(err.Error())
	}

}
