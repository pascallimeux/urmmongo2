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
	"github.com/pascallimeux/urmmongo2/server/model"
	"github.com/pascallimeux/urmmongo2/utils/log"
	"net/http"
)

//HTTP Get - /datasources
func (a *AppContext) getDatasourcesHandler(w http.ResponseWriter, r *http.Request) {
	log.Trace(log.Here(), "GetDatasourcesHandler() : calling method -")
	var datasources []model.DataSource
	var err error
	var j []byte
	datasources, err = a.Mongo.Get_datasources()
	if err != nil {
		sendError(log.Here(), w, err)
		return
	}
	j, err = json.Marshal(datasources)
	if err != nil {
		sendError(log.Here(), w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

//HTTP Post - /datasources
func (a *AppContext) postDatasourceHandler(w http.ResponseWriter, r *http.Request) {
	log.Trace(log.Here(), "PostDatasourceHandler() : calling method -")
	var datasource model.DataSource
	err := json.NewDecoder(r.Body).Decode(&datasource)
	if err != nil {
		sendError(log.Here(), w, err)
		return
	}
	err = a.Mongo.Create_datasource(&datasource)
	if err != nil {
		sendError(log.Here(), w, err)
		return
	}
	j, err := json.Marshal(datasource)
	if err != nil {
		sendError(log.Here(), w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

//HTTP Get - /datasources/{ds_id}
func (a *AppContext) getDatasourceHandler(w http.ResponseWriter, r *http.Request) {
	log.Trace(log.Here(), "GetDatasourceHandler() : calling method -")
	var datasource model.DataSource
	var err error
	vars := mux.Vars(r)
	ds_id := vars["ds_id"]
	datasource, err = a.Mongo.Get_datasource(ds_id)
	if err != nil {
		sendError(log.Here(), w, err)
		return
	}
	j, err := json.Marshal(datasource)
	if err != nil {
		sendError(log.Here(), w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
