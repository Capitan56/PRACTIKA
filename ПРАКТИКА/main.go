package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Server_ip   string `json:"server_ip"`
	Server_port string `json:"server_port"`
	Data_sorce  string `json:"data_sorce"`
}
type Data_Json struct {
	Id          int    `json:"id"`
	Request     string `json:"request"`
	Data_source string `json:"data_souce"`
}

func LoadConfiguration(filename string) (Config, error) {
	var config Config
	configFile, err := os.Open(filename)
	defer configFile.Close()
	if err != nil {
		return config, err
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	return config, err
}

func main() {
	config, err := LoadConfiguration("config.json")
	if err != nil {
		log.Fatal(err)
	}

	handler := http.NewServeMux()
	handler.HandleFunc("/handle_hook", func(w http.ResponseWriter, r *http.Request) {

	})
	s := &http.Server{
		ReadHeaderTimeout: 30 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		Handler:           handler,
		Addr:              config.Server_ip + ":" + config.Server_port,
	}

	log.Fatal(s.ListenAndServe())

}
