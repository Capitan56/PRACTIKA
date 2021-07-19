package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func Processor(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, err)
	}
	fmt.Fprint(rw, len(body))
}

func main() {
	http.HandleFunc("/handleHook/Processoring", Processor)
	err := http.ListenAndServe("127.0.0.1:3000", nil)
	log.Fatal(err)

}
