package main

import (
	"github.com/IFtech-A/urlshortener/internal/shortener/api/restapi"
	"github.com/Netflix/go-env"
	"github.com/sirupsen/logrus"
)

func main() {

	logrus.SetLevel(logrus.DebugLevel)
	conf := restapi.NewConfig()

	_, err := env.UnmarshalFromEnviron(conf)
	if err != nil {
		logrus.Error("unmarshaling from environment failed: ", err)
	}

	serv := restapi.New(conf)

	err = serv.Start()
	if err != nil {
		logrus.Fatal(err.Error())
	}

}
