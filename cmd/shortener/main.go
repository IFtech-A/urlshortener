package main

import (
	"github.com/IFtech-A/urlshortener/internal/shortener/api/restapi"
	"github.com/sirupsen/logrus"
)

func main() {

	logrus.SetLevel(logrus.DebugLevel)
	conf := restapi.NewConfig()

	serv := restapi.New(conf)

	err := serv.Start()
	if err != nil {
		logrus.Fatal(err.Error())
	}

}
