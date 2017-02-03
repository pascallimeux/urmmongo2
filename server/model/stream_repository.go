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
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

// Check if a stream exist from st_id
func (m MongoContext) checkStreamID(st_id string) (bool, error) {
	log.Trace(log.Here(), "checkStreamID() : calling method - ", st_id)
	if !bson.IsObjectIdHex(st_id) {
		return false, errors.New("bad format for streamID!")
	}
	mongoSession := m.Session.Clone()
	defer mongoSession.Close()
	c := mongoSession.DB(m.MongoDbName).C("stream")
	nb, err := c.Find(bson.M{"_id": bson.ObjectIdHex(st_id)}).Count()
	if nb == 1 && err == nil {
		return true, nil
	} else {
		return false, nil
	}
}

// Create index for stream
func (m MongoContext) CreateStreamIndex() error {
	log.Trace(log.Here(), "CreateStreamIndex() : calling method -")
	mongoSession := m.Session.Clone()
	defer mongoSession.Close()
	c := mongoSession.DB(m.MongoDbName).C("stream")
	index := mgo.Index{
		Key:        []string{"Ds_id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		log.Error(log.Here(), err.Error())
		return err
	}
	return nil
}

// Create a stream for a datasouce
func (m MongoContext) Create_stream(ds_id string, stream *Stream) error {
	log.Trace(log.Here(), "Create_stream() : calling method -")
	ds_exist, err := m.checkDatasourceID(ds_id)
	if err != nil {
		log.Error(log.Here(), err.Error())
		return err
	}
	if m.Control {
		if !ds_exist {
			err = errors.New("DatasourceID does not exist!")
			log.Error(log.Here(), err.Error())
			return err
		}
	} else {
		if !ds_exist {
			var datasource DataSource
			datasource.Id = bson.ObjectIdHex(ds_id)
			datasource.Name = "autocreate"
			datasource.Description = "Datasource auto create"
			log.Trace(log.Here(), "Autocreate a datasource ID=", datasource.Id.String())
			err := m.Create_datasource(&datasource)
			if err != nil {
				return err
			}
		}
	}
	if !bson.IsObjectIdHex(stream.Id.Hex()) {
		stream.Id = bson.NewObjectId()
		log.Trace(log.Here(), "create stream with a new ID:", stream.Id.Hex())
	} else {
		log.Trace(log.Here(), "create stream with ID:", stream.Id.Hex())
	}
	stream.Ds_id = ds_id
	stream.Date_creation = time.Now()
	mongoSession := m.Session.Clone()
	defer mongoSession.Close()
	c := mongoSession.DB(m.MongoDbName).C("stream")
	err = c.Insert(stream)
	if err != nil {
		log.Error(log.Here(), err.Error())
		return err
	}
	return nil

}

// Get all streams for a datasource by using ds_id
func (m MongoContext) Get_streams(ds_id string) ([]Stream, error) {
	log.Trace(log.Here(), "Get_streams() : calling method -")
	var streams = make([]Stream, 0)
	ds_exist, err := m.checkDatasourceID(ds_id)
	if err != nil {
		log.Error(log.Here(), err.Error())
		return streams, err
	}
	if !ds_exist {
		err = errors.New("DatasourceID does not exist!")
		log.Error(log.Here(), err.Error())
		return streams, err
	}
	mongoSession := m.Session.Clone()
	defer mongoSession.Close()
	c := mongoSession.DB(m.MongoDbName).C("stream")
	err = c.Find(bson.M{"ds_id": ds_id}).All(&streams)
	if err != nil {
		log.Error(log.Here(), err.Error())
		return streams, err
	}
	nb := strconv.Itoa(len(streams))
	log.Info(log.Here(), "get "+nb+" streams")
	return streams, nil
}

// Get a stream from ds_id and st_id
func (m MongoContext) Get_stream(ds_id, st_id string) (Stream, error) {
	log.Trace(log.Here(), "Get_stream() : calling method -")
	var stream Stream
	ds_exist, err := m.checkDatasourceID(ds_id)
	if err != nil {
		log.Error(log.Here(), err.Error())
		return stream, err
	}
	if !ds_exist {
		err = errors.New("DatasourceID does not exist!")
		log.Error(log.Here(), err.Error())
		return stream, err
	}
	mongoSession := m.Session.Clone()
	defer mongoSession.Close()
	c := mongoSession.DB(m.MongoDbName).C("stream")
	err = c.Find(bson.M{"_id": bson.ObjectIdHex(st_id)}).One(&stream)
	if err != nil {
		log.Error(log.Here(), err.Error())
		return stream, err
	}
	return stream, nil
}
