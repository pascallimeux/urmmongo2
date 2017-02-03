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
	"encoding/json"
	"github.com/pascallimeux/urmmongo/server/api"
	"github.com/pascallimeux/urmmongo/server/model"
	"github.com/pascallimeux/urmmongo/utils"
	"github.com/pascallimeux/urmmongo/utils/log"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	setup(true)
	code := m.Run()
	shutdown()
	os.Exit(code)
}

var applicationJSON string = "application/json"
var AppContext api.AppContext
var httpServerTest *httptest.Server
var logfile *os.File
var MOCK_HR string

//const MOCK_HR = `[{\"subject\": {\"reference\": \"7566939136933237222\", \"display\": \"personurm\"}, \"status\": \"final\", \"id\": \"glucose-meter\", \"Structure\": \"JSON\", \"component\": [{\"valueQuantity\": {\"system\": \"urn:std:iso:11073:10101\", \"code\": \"268192\", \"unit\": \"Cel\", \"value\": \"36.5\"}, \"code\": {\"coding\": [{\"system\": \"urn:std:iso:11073:10101\", \"code\": \"150364\", \"display\": \"MDC_TEMP_BODY\"}]}, \"valueDateTime\": \"2016-12-09T08:41:24+02:00\"}], \"issued\": \"2016-12-09T08:41:24+02:00\", \"identifier\": [{\"system\": \"Postman\", \"value\": \"227\"}], \"code\": {\"coding\": [{\"system\": \"urn:std:iso:11073:10101\", \"code\": \"528392\", \"display\": \"MDC_DEV_SPEC_PROFILE_TEMP\"}]}, \"resourceType\": \"Observation\", \"device\": {\"manufacturer\": \"xxxx \", \"model\": \"xxxxx\", \"udi\": \"xxxx\", \"version\": \"xxxx\"}}]`
const MOCKFILENAME = "./datafiles/Fhir_CAR_V1.json"
const MOCK_DS = "{\"name\":\"HD_pascal\",\"description\":\"Health data for login: Pascal\",\"serial\":\"None\"}"
const MOCK_ST = "{\"name\":\"ST_heart Rate\",\"description\":\"Stream for user Pascal and type: Heart rate\"}"
const MOCK_BAD_ST = ""
const MOCK_BAD_DS = ""

func testReadFile(filename string) string {
	payload, err := utils.ReadFile(filename)
	if err != nil {
		log.Fatal(log.Here(), "Error file reading: %v", string(len(filename)), err.Error())
	}
	return payload
}

func DropDB(session *mgo.Session, dbname string) {
	err := session.DB(dbname).DropDatabase()
	if err != nil {
		log.Fatal(log.Here(), "error:", err.Error())
	}
}

func setup(isDropDB bool) {
	// Read configuration file
	configuration, err := utils.Readconf("../config/configtest.json")
	if err != nil {
		log.Fatal(log.Here(), "error:", err.Error())
	}

	// Init logger
	logfile = log.Init_log(configuration.LogFileName, configuration.Logger)

	// Init mongosession
	mongoSession, err := model.GetMongoSession(configuration.MongoUrl, configuration.HandlerTimeout)
	if err != nil {
		log.Fatal(log.Here(), "error:", err.Error())
	}

	// Drop database
	if isDropDB {
		DropDB(mongoSession, configuration.MongoDbName)
	}

	// Init application context
	AppContext = api.AppContext{}
	AppContext.Mongo.Session = mongoSession
	AppContext.Mongo.MongoDbName = configuration.MongoDbName
	AppContext.Mongo.CreateIndex()

	// Init http server
	router := AppContext.CreateRoutes()
	httpServerTest = httptest.NewServer(router)

	// Read mack file for value
	MOCK_HR = testReadFile(MOCKFILENAME)
}

func shutdown() {
	defer logfile.Close()
	defer httpServerTest.Close()
	defer AppContext.Mongo.Session.Close()
}

func build_payload(value, datestr string) string {
	//payload := `[{"value":` + "\"" + value + "\"" + `, "at":` + "\"" + datestr + "\"" + `}]`
	//utils.Display_json(payload)
	payload := `[{"value":` + value + `, "at":` + "\"" + datestr + "\"" + `}]`
	return payload
}

func testCreateValue(datasourceId, streamId, valuestr, datestr string, t *testing.T) {
	payload := build_payload(valuestr, datestr)
	postData := strings.NewReader(payload)
	//log.Trace(log.Here(), "postdata: ", payload)
	res, err := http.Post(httpServerTest.URL+"/datasources/"+datasourceId+"/streams/"+streamId+"/values/", applicationJSON, postData)
	bytes, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	var value model.Value
	json.Unmarshal(bytes, &value)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Non-expected status code: %v\n\tbody: %v\n", http.StatusCreated, res.StatusCode)
	}
}

func testCreateST(datasourceId, data string, t *testing.T) string {
	postData := strings.NewReader(data)
	res, err := http.Post(httpServerTest.URL+"/datasources/"+datasourceId+"/streams/", applicationJSON, postData)
	bytes, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	var stream model.Stream
	json.Unmarshal(bytes, &stream)
	id := stream.Id.Hex()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Non-expected status code: %v\n\tbody: %v, data:%s\n", http.StatusCreated, res.StatusCode, id)
	}
	return id
}

func testCreateDS(data string, t *testing.T) string {
	postData := strings.NewReader(data)
	res, err := http.Post(httpServerTest.URL+"/datasources/", applicationJSON, postData)
	bytes, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	var datasource model.DataSource
	json.Unmarshal(bytes, &datasource)
	id := datasource.Id.Hex()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Non-expected status code: %v\n\tbody: %v, data:%s\n", http.StatusCreated, res.StatusCode, id)
	}
	return id
}

func testGetDS(datasourceId string, t *testing.T) {
	res, err := http.Get(httpServerTest.URL + "/datasources/" + datasourceId)
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	body := string(data)
	if res.StatusCode != http.StatusOK {
		t.Fatal("Non-expected status code: %v\n\tbody: %v, data:%s\n", http.StatusCreated, res.StatusCode, body)
	}
	if !strings.Contains(body, "{\"id\":\""+datasourceId+"\",") {
		t.Fatalf("Non-expected body content: %v", body)
	}
}

func testGetST(datasourceId, streamId string, t *testing.T) {
	res, err := http.Get(httpServerTest.URL + "/datasources/" + datasourceId + "/streams/" + streamId)
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	body := string(data)
	if res.StatusCode != http.StatusOK {
		t.Fatal("Non-expected status code: %v\n\tbody: %v, data:%s\n", http.StatusCreated, res.StatusCode, body)
	}
	if !strings.Contains(body, "{\"id\":\""+streamId+"\",") {
		t.Fatalf("Non-expected body content: %v", body)
	}
}

func testGetAllDS(t *testing.T) []model.DataSource {
	res, err := http.Get(httpServerTest.URL + "/datasources/")
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	body := string(data)
	if res.StatusCode != http.StatusOK {
		t.Fatal("Non-expected status code: %v\n\tbody: %v, data:%s\n", http.StatusCreated, res.StatusCode, body)
	}
	var dts []model.DataSource
	json.Unmarshal([]byte(body), &dts)
	return dts
}

func testGetAllST(datasourceId string, t *testing.T) []model.Stream {
	res, err := http.Get(httpServerTest.URL + "/datasources/" + datasourceId + "/streams")
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	body := string(data)
	if res.StatusCode != http.StatusOK {
		t.Fatal("Non-expected status code: %v\n\tbody: %v, data:%s\n", http.StatusCreated, res.StatusCode, body)
	}
	var sts []model.Stream
	json.Unmarshal([]byte(body), &sts)
	return sts
}

func testGetValues(datasourceId, streamId, params string, t *testing.T) []model.Value {
	res, err := http.Get(httpServerTest.URL + "/datasources/" + datasourceId + "/streams/" + streamId + "/values/" + params)
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	body := string(data)
	if res.StatusCode != http.StatusOK {
		t.Fatal("Non-expected status code: %v\n\tbody: %v, data:%s\n", http.StatusCreated, res.StatusCode, body)
	}
	var values []model.Value
	json.Unmarshal([]byte(body), &values)
	return values
}
