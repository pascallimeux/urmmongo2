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

package model

import (
	"encoding/json"
	"fmt"
	"github.com/pascallimeux/urmmongo2/utils/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type MongoContext struct {
	Session     *mgo.Session
	MongoDbName string
	Control     bool
}

type DataSource struct {
	Id            bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name          string        `json:"name" bson:"name"`
	Description   string        `json:"description" bson:"description"`
	Serial        string        `json:"serial" bson:"serial"`
	Date_creation time.Time     `json:"date_creation" bson:"date_creation"`
}

type Stream struct {
	Id            bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name          string        `json:"name" bson:"name"`
	Description   string        `json:"description" bson:"description"`
	Date_creation time.Time     `json:"date_creation" bson:"date_creation"`
	Ds_id         string        `json:"ds_id" bson:"ds_id"`
}

type Value struct {
	Id    bson.ObjectId   `json:"id" bson:"_id,omitempty"`
	Value json.RawMessage `json:"value" bson:"value"`
	At    time.Time       `json:"at" bson:"at"`
	Ds_id string          `json:"ds_id" bson:"ds_id"`
	St_id string          `json:"st_id" bson:"st_id"`
}

// get MongoDB session
func GetMongoSession(MongoUrl string, handlerTimeout time.Duration) (*mgo.Session, error) {
	log.Trace(log.Here(), "InitSession() : calling method -")
	session, err := mgo.Dial(MongoUrl)
	if err != nil {
		return nil, fmt.Errorf("MongoDB connection failed, with address '%s'.", MongoUrl)

	}
	session.SetMode(mgo.Monotonic, true)
	session.SetSocketTimeout(handlerTimeout * time.Second)
	return session, nil
}

// Init index
func (m MongoContext) CreateIndex() {
	m.CreateStreamIndex()
	m.CreateValueIndex()
}
