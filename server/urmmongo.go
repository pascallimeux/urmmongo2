/*
Copyright Pascal Limeux. 2016 All Rights Reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
		 http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"github.com/pascallimeux/urmmongo/server/api"
	"github.com/pascallimeux/urmmongo/server/model"
	"github.com/pascallimeux/urmmongo/utils"
	"github.com/pascallimeux/urmmongo/utils/log"
	"net/http"
	"os"
	"time"
)

func main() {

	// get arguments
	config_file := os.Getenv("URMMONGOCONFIGFILE")

	args := os.Args[1:]
	if len(args) == 1 {
		config_file = args[0]
	}

	if config_file == "" {
		panic(errors.New("No config file found!"))
	}

	// Init configuration
	configuration, err := utils.Readconf(config_file)
	if err != nil {
		panic(err.Error())
	}
	configuration.LogFileName = os.Getenv("URMMONGOLOGFILE")
	if configuration.LogFileName == "" {
		panic(errors.New("No logfile name defined!"))
	}

	// Init logger
	f := log.Init_log(configuration.LogFileName, configuration.Logger)
	defer f.Close()
	log.Info(log.Here(), configuration.To_string())

	// Init mongoDB
	mongoSession, err := model.GetMongoSession(configuration.MongoUrl, configuration.HandlerTimeout)
	if err != nil {
		log.Fatal(log.Here(), err.Error())
	}
	defer mongoSession.Close()

	// Init application context
	appContext := api.AppContext{}
	appContext.Mongo.Control = configuration.Control
	appContext.Mongo.Session = mongoSession
	appContext.Mongo.MongoDbName = configuration.MongoDbName
	appContext.Mongo.CreateIndex()

	// Start http server
	router := appContext.CreateRoutes()
	log.Info(log.Here(), "Listening on: ", configuration.HttpHostUrl)
	s := &http.Server{
		Addr:         configuration.HttpHostUrl,
		Handler:      router,
		ReadTimeout:  configuration.ReadTimeout * time.Second,
		WriteTimeout: configuration.WriteTimeout * time.Second,
	}
	log.Fatal(log.Here(), s.ListenAndServe().Error())
}
