/*
Copyright 2022 Red Hat, Inc.

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

package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	cdr "github.com/red-hat-storage/managed-fusion-metering/mock/api/v1"

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

var verbose = true

func debug(s string) {
	if verbose {
		log.Println(s)
	}
}

func init() {
	utilruntime.Must(monitoringv1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func GetToken(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func SendConsumptionReport(url string, report cdr.Request) ([]byte, error) {
	payload, err := json.Marshal(report)
	if err != nil {
		return nil, err
	}
	debug(string(payload))

	payloadBuffer := bytes.NewBuffer(payload)
	resp, err := http.Post(url, "application/json", payloadBuffer)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func GetAcknowledgement(url string, report cdr.Request) ([]byte, error) {
	payload, err := json.Marshal(report)
	if err != nil {
		return nil, err
	}
	debug(string(payload))

	payloadBuffer := bytes.NewBuffer(payload)
	resp, err := http.Post(url, "application/json", payloadBuffer)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func main() {
	cdrServiceHost := os.Getenv("MOCK_CDR_SERVICE_SERVICE_HOST")
	cdrServicePort := os.Getenv("MOCK_CDR_SERVICE_SERVICE_PORT")
	cdrEndpoint := "consumption"
	cdrServiceUrl := "http://" + cdrServiceHost + ":" + cdrServicePort + "/" + cdrEndpoint
	log.Println("MOCK_CDR_SERVICE_SERVICE_HOST: " + cdrServiceHost)
	log.Println("MOCK_CDR_SERVICE_SERVICE_PORT: " + cdrServicePort)

	tokenBody, err := GetToken(cdrServiceUrl)
	if err != nil {
		log.Fatalln(err)
	}

	debug(string(tokenBody))

	cdrReport := cdr.Request{
		Action: "create",
		DataSet: cdr.DataSet{
			Results: []cdr.Data{
				{UID: "3"},
			},
		},
	}

	sendBody, err := SendConsumptionReport(cdrServiceUrl, cdrReport)
	if err != nil {
		log.Fatalln(err)
	}

	debug(string(sendBody))

	cdrReport.Action = "ack"

	ackBody, err := GetAcknowledgement(cdrServiceUrl, cdrReport)
	if err != nil {
		log.Fatalln(err)
	}

	debug(string(ackBody))

	var ackResp cdr.Request
	err = json.Unmarshal(ackBody, &ackResp)
	if err != nil {
		log.Fatalln(err)
	}

	for _, data := range ackResp.DataSet.Results {
		// TODO
		debug(string(data.Status))
	}

	debug("Done.")
	os.Exit(0)
}
