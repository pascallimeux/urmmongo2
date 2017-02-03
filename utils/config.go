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

package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Configuration struct {
	Logger         string
	MongoUrl       string
	MongoDbName    string
	HttpHostUrl    string
	LogFileName    string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	HandlerTimeout time.Duration
	Control        bool
}

// Read configuration file and create Configuration
func Readconf(configFileName string) (Configuration, error) {
	configuration := Configuration{}
	if _, err := os.Stat(configFileName); err != nil {
		if os.IsNotExist(err) {
			return configuration, fmt.Errorf("Readconf(): config file does not exist!")
		}
	}
	file, _ := os.Open(configFileName)
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&configuration)
	if err != nil {
		return configuration, fmt.Errorf("Readconf(): config file error!", err)
	}
	return configuration, nil
}

func (c Configuration) To_string() string {
	mes := ("Start application: http server=" + c.HttpHostUrl + "  mongo server=" + c.MongoUrl + "  logger mode=" + c.Logger)
	return mes
}
