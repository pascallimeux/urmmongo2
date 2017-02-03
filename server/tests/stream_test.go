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
	"github.com/pascallimeux/urmmongo/server/model"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestStreamCreateAndGetNominal(t *testing.T) {
	datasourceId := testCreateDS(MOCK_DS, t)
	streamId := testCreateST(datasourceId, MOCK_ST, t)
	testGetST(datasourceId, streamId, t)
}

func TestStreamCreateBadDsID(t *testing.T) {
	AppContext.Mongo.Control = true
	postData := strings.NewReader(MOCK_ST)
	res, err := http.Post(httpServerTest.URL+"/datasources/"+"584bd567759b1421262bd9a0"+"/streams/", applicationJSON, postData)
	bytes, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	var stream model.Stream
	json.Unmarshal(bytes, &stream)
	id := stream.Id.Hex()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("Non-expected status code: %v\n\tbody: %v, data:%s\n", http.StatusBadRequest, res.StatusCode, id)
	}
}

func TestStreamCreateNewDsID(t *testing.T) {
	AppContext.Mongo.Control = false
	datasourceId := "111111111111111111111111"
	streamId := testCreateST(datasourceId, MOCK_ST, t)
	testGetDS(datasourceId, t)
	testGetST(datasourceId, streamId, t)
}

func TestStreamCreateBadValues(t *testing.T) {
	AppContext.Mongo.Control = true
	datasourceId := testCreateDS(MOCK_DS, t)
	postData := strings.NewReader(MOCK_BAD_ST)
	res, err := http.Post(httpServerTest.URL+"/datasources/"+datasourceId+"/streams/", applicationJSON, postData)
	bytes, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	var stream model.Stream
	json.Unmarshal(bytes, &stream)
	id := stream.Id.Hex()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("Non-expected status code: %v\n\tbody: %v, data:%s\n", http.StatusBadRequest, res.StatusCode, id)
	}
}

func TestStreamGetAllNominal(t *testing.T) {
	DropDB(AppContext.Mongo.Session, AppContext.Mongo.MongoDbName)
	datasourceId := testCreateDS(MOCK_DS, t)
	for i := 0; i < 10; i++ {
		testCreateST(datasourceId, MOCK_ST, t)
	}
	streams := testGetAllST(datasourceId, t)
	if len(streams) != 10 {
		t.Fatalf("Non-expected number of streams: %v", len(streams))
	}
}

func TestStreamGetBadStID(t *testing.T) {
	AppContext.Mongo.Control = true
	datasourceId := testCreateDS(MOCK_DS, t)
	streamId := "584bd567759b1421262bd9a0"
	res, err := http.Get(httpServerTest.URL + "/datasources/" + datasourceId + "/streams/" + streamId)
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

func TestStreamGetBadDsID(t *testing.T) {
	AppContext.Mongo.Control = true
	datasourceId := testCreateDS(MOCK_DS, t)
	streamId := testCreateST(datasourceId, MOCK_ST, t)
	datasourceId = "584bd567759b1421462bf9a0"
	res, err := http.Get(httpServerTest.URL + "/datasources/" + datasourceId + "/streams/" + streamId)
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
