package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

var requests = map[string]int{}
var result string

type RequestTest struct {
	NameRequest string
	Variant     [][]byte
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
	}
	log.Println(string(body))
	var t DataJson
	err = json.Unmarshal(body, &t)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, err)
		log.Println(err)
	}
	result = string(t.Request)
	log.Println(result)

	resp, err := http.Post("http://127.0.0.1:3000/handleHook/Processoring", "application/json", bytes.NewBuffer(t.Request))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, err)
		log.Println(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	b := string(body)
	data, err := strconv.Atoi(b)
	if err != nil {
		log.Fatal(err)
	}

	if requests[result] != 0 {
		fmt.Fprint(rw, requests[result])
	} else if requests[result] == 0 {
		requests[result] = data
		fmt.Fprint(rw, requests[result])

	}

}

func main() {

	config, err := LoadConfiguration("config.json")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/handleHook", test)
	err = http.ListenAndServe(config.ServerIp+":"+config.ServerPort, nil)
	log.Fatal(err)

}
