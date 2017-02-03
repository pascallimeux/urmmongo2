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
	"errors"
	"github.com/pascallimeux/urmmongo/utils/log"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

// Check if a datasource exist from ds_id
func (m MongoContext) checkDatasourceID(ds_id string) (bool, error) {
	log.Trace(log.Here(), "checkDatasourceID() : calling method -")
	if !bson.IsObjectIdHex(ds_id) {
		return false, errors.New("bad format for datasourceID!")
	}
	mongoSession := m.Session.Clone()
	defer mongoSession.Close()
	c := mongoSession.DB(m.MongoDbName).C("datasource")
	nb, err := c.Find(bson.M{"_id": bson.ObjectIdHex(ds_id)}).Count()
	if nb == 1 && err == nil {
		return true, nil
	} else {
		return false, nil
	}
}

// Create a datasource
func (m MongoContext) Create_datasource(datasource *DataSource) error {
	log.Trace(log.Here(), "Create_datasource() : calling method -")
	if !bson.IsObjectIdHex(datasource.Id.Hex()) {
		datasource.Id = bson.NewObjectId()
		log.Trace(log.Here(), "create datasource with a new ID:", datasource.Id.Hex())
	} else {
		log.Trace(log.Here(), "create datasource with ID:", datasource.Id.Hex())
	}
	datasource.Date_creation = time.Now()
	mongoSession := m.Session.Clone()
	defer mongoSession.Close()
	c := mongoSession.DB(m.MongoDbName).C("datasource")
	err := c.Insert(datasource)
	if err != nil {
		log.Error(log.Here(), err.Error())
		return err
	}
	return nil
}

// Get all datasources
func (m MongoContext) Get_datasources() ([]DataSource, error) {
	log.Trace(log.Here(), "Get_datasources() : calling method -")
	var datasources = make([]DataSource, 0)
	mongoSession := m.Session.Clone()
	defer mongoSession.Close()
	c := mongoSession.DB(m.MongoDbName).C("datasource")
	err := c.Find(nil).All(&datasources)
	if err != nil {
		log.Error(log.Here(), err.Error())
		return datasources, err
	}
	nb := strconv.Itoa(len(datasources))
	log.Info(log.Here(), "get "+nb+" datasources")
	return datasources, nil
}

// Get a datasource from ds_id
func (m MongoContext) Get_datasource(ds_id string) (DataSource, error) {
	log.Trace(log.Here(), "Get_datasource() : calling method -")
	var datasource DataSource
	_, err := m.checkDatasourceID(ds_id)
	if err != nil {
		log.Error(log.Here(), err.Error())
		return datasource, err
	}
	mongoSession := m.Session.Clone()
	defer mongoSession.Close()
	c := mongoSession.DB(m.MongoDbName).C("datasource")
	err = c.Find(bson.M{"_id": bson.ObjectIdHex(ds_id)}).One(&datasource)
	if err != nil {
		log.Error(log.Here(), err.Error())
		return datasource, err
	}
	return datasource, nil
}
