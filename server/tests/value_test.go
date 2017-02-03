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

package tests

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

// example to start one test: go test -run "TestValueCreateNominal"
func TestValueCreateNominal(t *testing.T) {
	datasourceId := testCreateDS(MOCK_DS, t)
	streamId := testCreateST(datasourceId, MOCK_ST, t)
	datestr := "2016-12-09T08:41:24+01:00"
	testCreateValue(datasourceId, streamId, MOCK_HR, datestr, t)
}

func TestValueGetAllNominal(t *testing.T) {
	DropDB(AppContext.Mongo.Session, AppContext.Mongo.MongoDbName)
	datasourceId := testCreateDS(MOCK_DS, t)
	streamId := testCreateST(datasourceId, MOCK_ST, t)
	datestr := "2016-12-09T08:41:30+02:00"
	testCreateValue(datasourceId, streamId, MOCK_HR, datestr, t)
	datestr = "2016-12-09T08:41:31+02:00"
	testCreateValue(datasourceId, streamId, MOCK_HR, datestr, t)
	datestr = "2016-12-09T08:41:32+02:00"
	testCreateValue(datasourceId, streamId, MOCK_HR, datestr, t)
	values := testGetValues(datasourceId, streamId, "", t)
	if len(values) != 3 {
		t.Fatalf("Non-expected number of streams: %v", len(values))
	}
}

func TestValueGetForAdayNominal(t *testing.T) {
	DropDB(AppContext.Mongo.Session, AppContext.Mongo.MongoDbName)
	datasourceId := testCreateDS(MOCK_DS, t)
	streamId := testCreateST(datasourceId, MOCK_ST, t)
	datestr := "2016-12-10T08:41:32+02:00"
	testCreateValue(datasourceId, streamId, MOCK_HR, datestr, t)
	testCreateValue(datasourceId, streamId, MOCK_HR, datestr, t)
	datestr = "2016-12-11T08:41:30+02:00"
	testCreateValue(datasourceId, streamId, MOCK_HR, datestr, t)
	params := "?atInterval=2016-12-10"
	values := testGetValues(datasourceId, streamId, params, t)
	if len(values) != 2 {
		t.Fatalf("Non-expected number of streams: %v", len(values))
	}
}

func TestValueGetWithDateFromNominal(t *testing.T) {
	DropDB(AppContext.Mongo.Session, AppContext.Mongo.MongoDbName)
	datasourceId := testCreateDS(MOCK_DS, t)
	streamId := testCreateST(datasourceId, MOCK_ST, t)
	datestr := "2016-12-05T09:41:30+02:00"
	testCreateValue(datasourceId, streamId, MOCK_HR, datestr, t)
	datestr = "2016-12-07T10:41:31+02:00"
	testCreateValue(datasourceId, streamId, MOCK_HR, datestr, t)
	datestr = "2016-12-08T12:41:32+02:00"
	testCreateValue(datasourceId, streamId, MOCK_HR, datestr, t)
	params := "?atInterval=2016-12-08,"
	values := testGetValues(datasourceId, streamId, params, t)
	if len(values) != 1 {
		t.Fatalf("Non-expected number of streams: %v", len(values))
	}
}

func TestValueGetWithDateToNominal(t *testing.T) {
	DropDB(AppContext.Mongo.Session, AppContext.Mongo.MongoDbName)
	datasourceId := testCreateDS(MOCK_DS, t)
	streamId := testCreateST(datasourceId, MOCK_ST, t)
	datestr := "2016-12-05T09:41:30+02:00"
	testCreateValue(datasourceId, streamId, MOCK_HR, datestr, t)
	datestr = "2016-12-07T10:41:31+02:00"
	testCreateValue(datasourceId, streamId, MOCK_HR, datestr, t)
	datestr = "2016-12-08T12:41:32+02:00"
	testCreateValue(datasourceId, streamId, MOCK_HR, datestr, t)
	params := "?atInterval=,2016-12-07"
	values := testGetValues(datasourceId, streamId, params, t)
	if len(values) != 2 {
		t.Fatalf("Non-expected number of streams: %v", len(values))
	}
}

func TestValueGetWithPeriodNominal(t *testing.T) {
	DropDB(AppContext.Mongo.Session, AppContext.Mongo.MongoDbName)
	datasourceId := testCreateDS(MOCK_DS, t)
	streamId := testCreateST(datasourceId, MOCK_ST, t)
	datestr := "2016-12-10T09:41:30+02:00"
	testCreateValue(datasourceId, streamId, MOCK_HR, datestr, t)
	datestr = "2016-12-11T09:41:30+02:00"
	testCreateValue(datasourceId, streamId, MOCK_HR, datestr, t)
	datestr = "2016-12-08T09:41:30+02:00"
	testCreateValue(datasourceId, streamId, MOCK_HR, datestr, t)
	params := "?atInterval=2016-12-10,2016-12-12"
	values := testGetValues(datasourceId, streamId, params, t)
	if len(values) != 2 {
		t.Fatalf("Non-expected number of streams: %v", len(values))
	}
}

func TestValueGetWithIntervalBetweenValuesNominal(t *testing.T) {
	DropDB(AppContext.Mongo.Session, AppContext.Mongo.MongoDbName)
	datasourceId := testCreateDS(MOCK_DS, t)
	streamId := testCreateST(datasourceId, MOCK_ST, t)
	datestr := "2016-12-08T09:41:00+02:00"
	testCreateValue(datasourceId, streamId, MOCK_HR, datestr, t)
	datestr = "2016-12-08T09:41:29+02:00"
	testCreateValue(datasourceId, streamId, MOCK_HR, datestr, t)
	datestr = "2016-12-08T09:41:49+02:00"
	testCreateValue(datasourceId, streamId, MOCK_HR, datestr, t)
	params := "?interval_between_values=30"
	values := testGetValues(datasourceId, streamId, params, t)
	if len(values) != 2 {
		t.Fatalf("Non-expected number of streams: %v", len(values))
	}
}

func TestValueGetBadStID(t *testing.T) {
	AppContext.Mongo.Control = true
	DropDB(AppContext.Mongo.Session, AppContext.Mongo.MongoDbName)
	datasourceId := testCreateDS(MOCK_DS, t)
	streamId := "584bd567749b1421462bf9a0"
	res, err := http.Get(httpServerTest.URL + "/datasources/" + datasourceId + "/streams/" + streamId + "/values/")
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	body := string(data)
	if res.StatusCode != http.StatusBadRequest {
		t.Fatal("Non-expected status code: %v\n\tbody: %v, data:%s\n", http.StatusBadRequest, res.StatusCode, body)
	}
}

func TestValueGetBadDsID(t *testing.T) {
	AppContext.Mongo.Control = true
	DropDB(AppContext.Mongo.Session, AppContext.Mongo.MongoDbName)
	datasourceId := "584bd567759b1421472bf9a0"
	streamId := "584bd567749b1421462bf9a0"
	res, err := http.Get(httpServerTest.URL + "/datasources/" + datasourceId + "/streams/" + streamId + "/values/")
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	body := string(data)
	if res.StatusCode != http.StatusBadRequest {
		t.Fatal("Non-expected status code: %v\n\tbody: %v, data:%s\n", http.StatusBadRequest, res.StatusCode, body)
	}
}

func TestValueCreateBadStID(t *testing.T) {
	AppContext.Mongo.Control = true
	datasourceId := testCreateDS(MOCK_DS, t)
	streamId := "222222222222222222222222"
	datestr := "2016-12-09T08:41:24+02:00"
	payload := build_payload(MOCK_HR, datestr)
	postData := strings.NewReader(payload)
	res, err := http.Post(httpServerTest.URL+"/datasources/"+datasourceId+"/streams/"+streamId+"/values/", applicationJSON, postData)
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	body := string(data)
	if res.StatusCode != http.StatusBadRequest {
		t.Fatal("Non-expected status code: %v\n\tbody: %v, data:%s\n", http.StatusBadRequest, res.StatusCode, body)
	}
}

func TestValueCreateBadDsID(t *testing.T) {
	AppContext.Mongo.Control = true
	datasourceId := testCreateDS(MOCK_DS, t)
	streamId := testCreateST(datasourceId, MOCK_ST, t)
	datasourceId = "333333333333333333333333"
	datestr := "2016-12-09T08:41:24+02:00"
	payload := build_payload(MOCK_HR, datestr)
	postData := strings.NewReader(payload)
	res, err := http.Post(httpServerTest.URL+"/datasources/"+datasourceId+"/streams/"+streamId+"/values/", applicationJSON, postData)
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	body := string(data)
	if res.StatusCode != http.StatusBadRequest {
		t.Fatal("Non-expected status code: %v\n\tbody: %v, data:%s\n", http.StatusBadRequest, res.StatusCode, body)
	}
}

func TestValueCreateWithNewStID(t *testing.T) {
	AppContext.Mongo.Control = false
	datasourceId := testCreateDS(MOCK_DS, t)
	streamId := "584bd567749b1421462bf9a0"
	datestr := "2016-12-09T08:41:24+02:00"
	testCreateValue(datasourceId, streamId, MOCK_HR, datestr, t)
}

func TestValueCreateWithNewDtID(t *testing.T) {
	AppContext.Mongo.Control = false
	datasourceId := "584bf567759b1421472bf9a0"
	streamId := testCreateST(datasourceId, MOCK_ST, t)
	datestr := "2016-12-09T08:41:24+02:00"
	testCreateValue(datasourceId, streamId, MOCK_HR, datestr, t)
}

func TestValueCreateWithNewStAndDtID(t *testing.T) {
	AppContext.Mongo.Control = false
	datasourceId := "584bd567749b1445462bf9a0"
	streamId := "584bd567349b1421462bf9a0"
	datestr := "2016-12-09T08:41:24+02:00"
	testCreateValue(datasourceId, streamId, MOCK_HR, datestr, t)
}
