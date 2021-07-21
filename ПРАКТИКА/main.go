package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var requests = map[string][]byte{}
var Req string

type MapDelete struct {
	Id json.RawMessage
}

type DataJson struct {
	Id         string `json:"id"`
	Request    json.RawMessage
	DataSource string `json:"datasource"`
}

type Config struct {
	ServerIp   string `json:"server_ip"`
	ServerPort string `json:"server_port"`
	DataSource string `json:"datasource"`
}

func LoadConfiguration(filename string) (Config, error) {
	var config Config
	configFile, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	return config, err
}

func test(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, err)
		log.Println(err)
		return
	}
	log.Println(string(body))
	var t DataJson
	err = json.Unmarshal(body, &t)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, err)
		log.Println(err)
		return
	}
	log.Println(string(t.Request))

	var idRequest = t.Id + string(t.Request)

	if value, ok := requests[idRequest]; ok == true {

		fmt.Fprint(rw, value)

	} else if resp, err := http.Post("http://127.0.0.1:3000/handleHook/Processoring", "application/json", bytes.NewBuffer(t.Request)); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(rw, err)
		log.Println(err)

	} else {

		body, err = ioutil.ReadAll(resp.Body)
		requests[idRequest] = body
		fmt.Fprint(rw, body)
	}

}

func deleteMap(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, err)
		log.Println(err)
		return
	}
	log.Println(string(body))
	var d MapDelete
	err = json.Unmarshal(body, &d)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, err)
		log.Println(err)
		return
	}

	for key := range requests {
		if string([]rune(string(d.Id))[1]) == string([]rune(key)[0]) {
			delete(requests, key)
		}
	}
}

func main() {

	config, err := LoadConfiguration("config.json")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/handleHook", test)
	http.HandleFunc("/handleHook/delete_cached", deleteMap)
	err = http.ListenAndServe(config.ServerIp+":"+config.ServerPort, nil)
	log.Fatal(err)
}
