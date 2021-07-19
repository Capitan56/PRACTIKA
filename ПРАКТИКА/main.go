package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

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
	}
	log.Println(string(body))
	var t DataJson
	err = json.Unmarshal(body, &t)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, err)
	}
	log.Println(t.Request)

	resp, err := http.Post("http://127.0.0.1:3000/handleHook/Processoring", "application/json", bytes.NewBuffer(t.Request))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, err)
	}

	io.Copy(rw, resp.Body)

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

