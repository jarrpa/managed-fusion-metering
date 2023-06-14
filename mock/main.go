package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	cdr "github.com/red-hat-storage/managed-fusion-metering/mock/api/v1"

	"github.com/gorilla/mux"
)

var Reports map[string]cdr.Data

func getToken(w http.ResponseWriter, r *http.Request) {
	token := "token1234567890"
	json.NewEncoder(w).Encode(token)
}

func handleReports(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var req cdr.Request
	err := json.Unmarshal(reqBody, &req)
	if err != nil {
		log.Printf("ERROR: %v", err)
		json.NewEncoder(w).Encode(req)
	}

	switch action := req.Action; action {
	case "create":
		for i, data := range req.DataSet.Results {
			if d, ok := Reports[data.UID]; ok {
				data = d
				log.Printf("Record already exists: %v", data)
			} else {
				data.Status = cdr.DataStatusReceived
				Reports[data.UID] = data
				log.Printf("Record received: %v", data)
			}
			req.DataSet.Results[i] = data
		}
	case "ack":
		for i, data := range req.DataSet.Results {
			if d, ok := Reports[data.UID]; ok {
				data = d
				log.Printf("Record found: %v", data)
			} else {
				log.Printf("Record not found: %v", data)
			}
			req.DataSet.Results[i] = data
		}
	default:
		w.WriteHeader(http.StatusNotAcceptable)
		err := fmt.Errorf("Invalid action: %s", action)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(req)
}

func main() {
	Reports = map[string]cdr.Data{
		"1": {UID: "1"},
		"2": {UID: "2"},
	}

	cdrRouter := mux.NewRouter().StrictSlash(true)
	cdrRouter.HandleFunc("/consumption", getToken).Methods("GET")
	cdrRouter.HandleFunc("/consumption", handleReports).Methods("POST")
	log.Fatal(http.ListenAndServe(":10000", cdrRouter))
}
