package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func Processor(rw http.ResponseWriter, req *http.Request) {
	_, err := ioutil.ReadAll(req.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, err)
		log.Println(err)
	}

	seconds := time.Now().Unix()
	rand.Seed(seconds)
	fmt.Fprint(rw, rand.Intn(100)+1)
}

func main() {
	http.HandleFunc("/handleHook/Processoring", Processor)
	err := http.ListenAndServe("127.0.0.1:3000", nil)
	log.Fatal(err)

}
