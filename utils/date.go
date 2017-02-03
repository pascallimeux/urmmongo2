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
	"github.com/pascallimeux/urmmongo/utils/log"
	"time"
)

const DATEFORMAT1 = "2006-01-02T15:04:05.000Z"
const DATEFORMAT2 = "2006-01-02"
const DATEFORMAT3 = "2006-01-02 15:04:05 -0700 -0700"
const DATEFORMAT4 = "2006-01-02 15:04:05 -0700 MST"

func DateParse(datestr, layout string, all_the_day bool) (time.Time, error) {
	log.Trace(log.Here(), "DateParse() : calling method for: ", datestr)
	date, err := time.Parse(layout, datestr)
	if err != nil {
		log.Error(log.Here(), "DateParse error: ", err.Error())
	}
	if all_the_day && layout == DATEFORMAT2 {
		date = date.Add(time.Duration(86400) * time.Second)
	}
	return date, err
}
