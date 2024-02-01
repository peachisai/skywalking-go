// Licensed to Apache Software Foundation (ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation (ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package go_elasticsearchv7

import (
	"bytes"
	"context"
	"encoding/json"
	_ "github.com/apache/skywalking-go"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"log"
	"net/http"
)

var testCases []struct {
	caseName string
	fn       testFunc
}

type testFunc func(es *elasticsearch.Client) error

func init() {
	log.Println("-----start append testCases------")
	append(testCases, struct {
		caseName string
		fn       testFunc
	}{caseName: "testIndex", fn: testIndex})
	log.Println("-----finish append testCases------")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func executeHandler(w http.ResponseWriter, r *http.Request) {
	es := initClient()
	for _, testCase := range testCases {
		log.Printf("start execute test case: %s", testCase.caseName)
		err := testCase.fn(es)
		if err != nil {
			log.Fatalf("fail to execute test case,name:%s", testCase.caseName)
		}
		log.Printf("finish execute test case: %s", testCase.caseName)
	}
}

func initClient() *elasticsearch.Client {
	address := []string{"http://43.139.166.178:9200"}

	cfg := elasticsearch.Config{
		Addresses: address,
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	info, err := es.Info()
	defer info.Body.Close()
	log.Println("connect es success")
	return es
}

func main() {

	http.HandleFunc("/health", healthHandler)

	http.HandleFunc("/execute", executeHandler)

	err := http.ListenAndServe(":8000", nil)

	if err != nil {
		log.Fatalf(err.Error())
	}
}

func testIndex(es *elasticsearch.Client) error {

	data, err := json.Marshal(struct {
		name string
	}{
		name: "sw-go",
	})
	if err != nil {
		log.Fatalf("Error marshaling document: %s", err)
	}

	log.Println(string(data))

	req := esapi.IndexRequest{
		Index:      "sw-index",
		Body:       bytes.NewReader(data),
		DocumentID: "1",
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), es)

	if err != nil || res.IsError() {
		log.Fatalf("Error getting index request'response: %s", err)
	}

	defer res.Body.Close()
	return nil
}
