package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func Processor(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(rw, len(body))
}

func main() {

	handler := http.NewServeMux()
	handler.HandleFunc("/handleHook/Processoring", Processor)

	s := &http.Server{
		ReadHeaderTimeout: 30 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		Handler:           handler,
		Addr:              "127.0.0.1:3000",
	}
	log.Fatal(s.ListenAndServe())

}
