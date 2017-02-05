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

//HTTP Get - /datasources/{ds_id}/streams
func (a *AppContext) getStreamsHandler(w http.ResponseWriter, r *http.Request) {
	log.Trace(log.Here(), "GetStreamsHandler() : calling method -")
	var streams []model.Stream
	var err error
	vars := mux.Vars(r)
	ds_id := vars["ds_id"]
	streams, err = a.Mongo.Get_streams(ds_id)
	if err != nil {
		sendError(log.Here(), w, err)
		return
	}
	j, err := json.Marshal(streams)
	if err != nil {
		sendError(log.Here(), w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

//HTTP Post - /datasources/{ds_id}/streams
func (a *AppContext) postStreamHandler(w http.ResponseWriter, r *http.Request) {
	log.Trace(log.Here(), "PostStreamHandler() : calling method -")
	var stream model.Stream
	vars := mux.Vars(r)
	ds_id := vars["ds_id"]
	err := json.NewDecoder(r.Body).Decode(&stream)
	if err != nil {
		sendError(log.Here(), w, err)
		return
	}
	err = a.Mongo.Create_stream(ds_id, &stream)
	if err != nil {
		sendError(log.Here(), w, err)
		return
	}
	j, err := json.Marshal(stream)
	if err != nil {
		sendError(log.Here(), w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

//HTTP Get - /datasources/{ds_id}/streams/{st_id}
func (a *AppContext) getStreamHandler(w http.ResponseWriter, r *http.Request) {
	log.Trace(log.Here(), "GetStreamHandler() : calling method -")
	var stream model.Stream
	vars := mux.Vars(r)
	var err error
	ds_id := vars["ds_id"]
	st_id := vars["st_id"]
	stream, err = a.Mongo.Get_stream(ds_id, st_id)
	if err != nil {
		sendError(log.Here(), w, err)
		return
	}
	j, err := json.Marshal(stream)
	if err != nil {
		sendError(log.Here(), w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
