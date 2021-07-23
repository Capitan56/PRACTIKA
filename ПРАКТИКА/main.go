package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Cache struct {
	requests map[string][]byte
	mux      sync.Mutex
}

var requestsMap = Cache{
	requests: map[string][]byte{},
	mux:      sync.Mutex{},
}

type MapDelete struct {
	Id string `json:"Id"`
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

func insertCached(c *Cache, keyMap string, valueMap []byte) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.requests[keyMap] = valueMap
}

func deleteCached(c *Cache, Id string) {
	c.mux.Lock()
	defer c.mux.Unlock()
	delete(c.requests, Id)
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

	if value, ok := requestsMap.requests[idRequest]; ok == true {

		fmt.Fprint(rw, value)

	} else if resp, err := http.Post("http://127.0.0.1:3000/handleHook/Processoring", "application/json", bytes.NewBuffer(t.Request)); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(rw, err)
		log.Println(err)

	} else {

		body, err = ioutil.ReadAll(resp.Body)
		insertCached(&requestsMap, idRequest, body)
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

	for key := range requestsMap.requests {
		if strings.HasPrefix(key, d.Id) == true {
			deleteCached(&requestsMap, d.Id)
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
