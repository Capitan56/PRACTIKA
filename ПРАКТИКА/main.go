package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Server struct {
		Server_ip   string `json:"server_ip"`
		Server_port string `json:"server_port"`
		Data_sorce  string `json:"data_sorce"`
	}
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
	handler := http.NewServeMux()
	handler.HandleFunc("/SERVER", func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Starting the application...")
		config, _ := LoadConfiguration("config.json")
		fmt.Println(config.Server.Data_sorce)
	})
	s := &http.Server{
		ReadHeaderTimeout: 30 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		Handler:           handler,
		Addr:              ":8080",
	}

	log.Fatal(s.ListenAndServe())

}
