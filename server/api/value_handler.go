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

package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pascallimeux/urmmongo/server/model"
	//"github.com/pascallimeux/urmmongo/utils"
	"github.com/pascallimeux/urmmongo/utils/log"
	"io/ioutil"
	"net/http"
)

//HTTP Get - /datasources/{ds_id}/streams/{st_id}/values
func (a *AppContext) getValuesHandler(w http.ResponseWriter, r *http.Request) {
	log.Trace(log.Here(), "GetValuesHandler() : calling method -")
	var values []model.Value
	var err error
	vars := mux.Vars(r)
	ds_id := vars["ds_id"]
	st_id := vars["st_id"]
	at_params := r.URL.Query()["atInterval"]
	date := ""
	if len(at_params) > 0 {
		date = at_params[0]
	}
	inter_params := r.URL.Query()["interval_between_values"]
	interval := ""
	if len(inter_params) > 0 {
		interval = inter_params[0]
	}
	values, err = a.Mongo.Get_values(ds_id, st_id, date, interval)
	if err != nil {
		sendError(log.Here(), w, err)
		return
	}
	j, err := json.Marshal(values)
	if err != nil {
		sendError(log.Here(), w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

//HTTP Post - /datasources/{ds_id}/streams/{st_id}/values
func (a *AppContext) postValueHandler(w http.ResponseWriter, r *http.Request) {
	log.Trace(log.Here(), "PostValueHandler() : calling method -")
	var values []model.Value
	var value model.Value

	vars := mux.Vars(r)
	ds_id := vars["ds_id"]
	st_id := vars["st_id"]

	bytes, _ := ioutil.ReadAll(r.Body)
	//log.Trace(log.Here(), "Value:", string(bytes))
	err := json.Unmarshal(bytes, &values)
	if err != nil {
		sendError(log.Here(), w, err)
		return
	}

	//dec := json.NewDecoder(r.Body)
	//err = dec.Decode(&values)
	//if err != nil {
	//	sendError(log.Here(), w, err)
	//	return
	//}

	value = values[0]
	//value.At, err = utils.CheckDateFormat(value.At)
	//if err != nil {
	//	sendError(log.Here(), w, err)
	//	return
	//}
	log.Trace(log.Here(), "create value for DSID:", ds_id, "  STDI: ", st_id, "with date: ", value.At.String())
	err = a.Mongo.Create_value(ds_id, st_id, &value)
	if err != nil {
		sendError(log.Here(), w, err)
		return
	}
	j, err := json.Marshal(value)
	if err != nil {
		sendError(log.Here(), w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}
