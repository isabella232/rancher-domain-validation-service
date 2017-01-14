package main

import (
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rancher/rancher-domain-validaiton-service/manager"
	"github.com/rancher/rancher-domain-validaiton-service/service"
	"github.com/urfave/cli"
)

//VERSION for Rancher Authantication Filter Service
var VERSION = "v0.1.0-dev"

func main() {

	///init parsing command line
	app := cli.NewApp()
	app.Name = "rancher-auth-filter-service"
	app.Version = "v0.1.0-dev"
	app.Usage = "Rancher Authantication Filter Service"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "rancherUrl",
			Value:  "http://54.255.182.226:8080/",
			Usage:  "Rancher server url",
			EnvVar: "RANCHER_SERVER_URL",
		},
		cli.StringFlag{
			Name:   "localport",
			Value:  "8092",
			Usage:  "Local server port ",
			EnvVar: "LOCAL_VALIDATION_FILTER_PORT",
		},
	}

	app.Action = func(c *cli.Context) error {

		logrus.Infof("Starting token validation service")
		logrus.Infof("Rancher server URL:" + manager.URL + ". The validation filter server running on local port:" + manager.Port)
		//create mux router
		router := service.NewRouter()
		http.Handle("/", router)
		serverString := ":" + manager.Port
		//start local server
		manager.CreateDatabase()
		logrus.Fatal(http.ListenAndServe(serverString, nil))
		return nil
	}

	app.Run(os.Args)

}
