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
	"bytes"
	"encoding/json"
	"github.com/pascallimeux/urmmongo2/utils/log"
	"io/ioutil"
	"regexp"
	"strings"
)

func Display_json(jsonstr string) {
	var out bytes.Buffer
	json.Indent(&out, []byte(jsonstr), "", "  ")
	log.Trace(log.Here(), "Json object: ", out.String())
}

func Struct_toString(i interface{}) (string, error) {
	st := ""
	out, err := json.Marshal(i)
	if err != nil {
		return st, err
	}
	return string(out), nil
}

func ReadFile(filename string) (string, error) {
	var payload string
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error(log.Here(), err.Error())
		return payload, err
	}

	payload = string(raw)
	payload = strings.Replace(payload, "\r\n", " ", -1)
	//payload = strings.Replace(payload, "\"", "\\\"", -1)
	re := regexp.MustCompile("  +")
	payload = re.ReplaceAllString(payload, " ")
	return payload, nil
}
