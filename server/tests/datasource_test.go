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

func TestDatasourceCreateAndGetNominal(t *testing.T) {
	datasourceId := testCreateDS(MOCK_DS, t)
	testGetDS(datasourceId, t)
}

func TestDatasourceGetAllNominal(t *testing.T) {
	DropDB(AppContext.Mongo.Session, AppContext.Mongo.MongoDbName)
	for i := 0; i < 10; i++ {
		testCreateDS(MOCK_DS, t)
	}
	datasources := testGetAllDS(t)
	if len(datasources) != 10 {
		t.Fatalf("Non-expected number of datasources: %v", len(datasources))
	}
}

func TestDatasourceCreateBadValues(t *testing.T) {
	postData := strings.NewReader(MOCK_BAD_DS)
	res, err := http.Post(httpServerTest.URL+"/datasources/", applicationJSON, postData)
	bytes, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	var datasource model.DataSource
	json.Unmarshal(bytes, &datasource)
	id := datasource.Id.Hex()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("Non-expected status code: %v\n\tbody: %v, data:%s\n", http.StatusBadRequest, res.StatusCode, id)
	}
}

func TestDatasourceGetBadID(t *testing.T) {
	datasourceId := "584bd567759b1421262bd9a0"
	res, err := http.Get(httpServerTest.URL + "/datasources/" + datasourceId)
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
