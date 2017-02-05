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
	"github.com/pascallimeux/urmmongo2/utils"
	"github.com/pascallimeux/urmmongo2/utils/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"strings"
	"time"
)

// Create index for value
func (m MongoContext) CreateValueIndex() error {
	log.Trace(log.Here(), "CreateValueIndex() : calling method -")
	mongoSession := m.Session.Clone()
	defer mongoSession.Close()
	c := mongoSession.DB(m.MongoDbName).C("value")
	index := mgo.Index{
		Key:        []string{"Ds_id", "St_id"},
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

// Create a value for a datasource and a stream
func (m MongoContext) Create_value(ds_id, st_id string, value *Value) error {
	log.Trace(log.Here(), "Create_value() : calling method -")
	value.Id = bson.NewObjectId()
	var ds_exist, st_exist bool
	var err error
	ds_exist, err = m.checkDatasourceID(ds_id)
	st_exist, err = m.checkStreamID(st_id)
	if err != nil {
		log.Error(log.Here(), err.Error())
		return err
	}
	if m.Control {
		if !ds_exist || !st_exist {
			log.Error(log.Here(), "DatasourceID or StreamID does not exist!")
			return errors.New("DatasourceID or StreamID does not exist!")
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
		if !st_exist {
			var stream Stream
			stream.Id = bson.ObjectIdHex(st_id)
			stream.Ds_id = ds_id
			stream.Name = "autocreate"
			stream.Description = "Stream auto create"
			log.Trace(log.Here(), "Autocreate a stream ID=", stream.Id.String())
			err := m.Create_stream(ds_id, &stream)
			if err != nil {
				return err
			}
		}
	}
	value.Ds_id = ds_id
	value.St_id = st_id

	at, err := utils.DateParse(value.At.String(), utils.DATEFORMAT3, false)
	if err != nil {
		at, err = utils.DateParse(value.At.String(), utils.DATEFORMAT4, false)
		if err != nil {
			log.Error(log.Here(), err.Error())
			return err
		}
	}
	value.At = at
	mongoSession := m.Session.Clone()
	defer mongoSession.Close()
	c := mongoSession.DB(m.MongoDbName).C("value")
	err = c.Insert(value)
	if err != nil {
		log.Error(log.Here(), err.Error())
		return err
	}
	return nil
}

// Get values for a datasource, a stream by using parameters like date and minimum interval between values
func (m MongoContext) Get_values(ds_id, st_id, date_params, interval string) ([]Value, error) {
	log.Trace(log.Here(), "Get_values() : calling method with dateparams:-", date_params)
	timer := utils.Timer{}
	timer.StartTimer()
	var values = make([]Value, 0)
	var err error
	var ds_exist, st_exist bool
	ds_exist, err = m.checkDatasourceID(ds_id)
	if err != nil {
		log.Error(log.Here(), err.Error())
		return values, err
	}
	st_exist, err = m.checkStreamID(st_id)
	if err != nil {
		log.Error(log.Here(), err.Error())
		return values, err
	}
	if !ds_exist || !st_exist {
		log.Error(log.Here(), "DatasourceID or StreamID does not exist!")
		return values, errors.New("DatasourceID or StreamID does not exist!")
	}

	mongoSession := m.Session.Clone()
	defer mongoSession.Close()
	c := mongoSession.DB(m.MongoDbName).C("value")
	if date_params != "" {
		fromDate, toDate, err := check_dates(date_params)
		if err != nil {
			log.Error(log.Here(), err.Error())
			return values, err
		}
		log.Trace(log.Here(), "request: from=", fromDate.String(), " to=", toDate.String())
		err = c.Find(
			bson.M{
				"ds_id": ds_id,
				"st_id": st_id,
				"at": bson.M{
					"$gte": fromDate,
					"$lte": toDate,
				},
			}).All(&values)
	} else {
		err = c.Find(bson.M{"ds_id": ds_id, "st_id": st_id}).All(&values)
	}
	if err != nil {
		log.Error(log.Here(), err.Error())
		return values, err
	}
	nb := len(values)
	if nb > 1 && interval != "" {
		inter, err := strconv.Atoi(interval)
		if err != nil {
			log.Error(log.Here(), err.Error())
			return values, err
		}
		values = check_interval(values, inter)
	}
	nb = len(values)
	timer.LogElapsed(log.Here(), "get "+strconv.Itoa(nb)+" values")
	return values, nil
}

// Check minimum interval between values
func check_interval(values []Value, interval int) []Value {
	log.Trace(log.Here(), "check_interval() : calling method -")
	previous_date := time.Time{}
	i := 0
	nb_all_values := len(values)
	for _, value := range values {
		if value.At.Sub(previous_date) > 0 {
			values[i] = value
			previous_date = value.At.Add(time.Duration(interval) * time.Second)
			i++
		}
	}
	values = values[:i]
	log.Info(log.Here(), "filter by interval:"+strconv.Itoa(len(values))+" of "+strconv.Itoa(nb_all_values))
	return values
}

// Check dates parameters
func check_dates(params string) (time.Time, time.Time, error) {
	log.Trace(log.Here(), "check_date_params() : calling method -")
	date1_TS := time.Time{}
	date2_TS := time.Now()
	var err error
	if strings.ContainsAny(params, ",") {
		index := strings.Index(params, ",")
		date1 := params[0:index]
		if date1 != "" {
			date1 = strings.TrimSpace(date1)
			date1_TS, err = utils.DateParse(date1, utils.DATEFORMAT2, false)
			if err != nil {
				return date1_TS, date2_TS, err
			}
		}
		date2 := params[index+1 : len(params)]
		if date2 != "" {
			date2 = strings.TrimSpace(date2)
			date2_TS, err = utils.DateParse(date2, utils.DATEFORMAT2, true)
			if err != nil {
				return date1_TS, date2_TS, err
			}
		}
		return date1_TS, date2_TS, nil
	} else {
		params = strings.TrimSpace(params)
		date1_TS, err = utils.DateParse(params, utils.DATEFORMAT2, false)
		if err != nil {
			return date1_TS, date2_TS, err
		}
		date2_TS = date1_TS.AddDate(0, 0, +1)
		return date1_TS, date2_TS, nil
	}
}
