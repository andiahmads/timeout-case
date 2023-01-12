package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-retryablehttp"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", requestServerWithHTTPStandart).Methods(http.MethodGet)
	r.HandleFunc("/test", requestWithretryablehttp).Methods(http.MethodGet)
	r.HandleFunc("/test-mock", requestWithretryablehttpMock).Methods(http.MethodGet)

	port := "8787"
	server := &http.Server{Addr: ":" + port, Handler: r}
	log.Printf("starting service on port %s", port)
	server.ListenAndServe()
}

func requestServerWithHTTPStandart(w http.ResponseWriter, r *http.Request) {
	c := &http.Client{Timeout: 2 * time.Second}

	res, err := c.Get("http://localhost:8888")
	if err != nil {
		log.Println(err)
	}

	defer res.Body.Close()

	if err != nil {
		log.Println(err)
	}
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}

func requestWithretryablehttp(w http.ResponseWriter, r *http.Request) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Transport = tr
	retryClient.RetryMax = 10
	retryClient.HTTPClient.Timeout = 1
	standardClient := retryClient.StandardClient()

	req, err := standardClient.Get("http://localhost:8888")

	if err != nil {
		log.Println(err)
	}
	defer req.Body.Close()
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))

}

func requestWithretryablehttpMock(w http.ResponseWriter, r *http.Request) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Transport = tr
	retryClient.RetryMax = 10
	// retryClient.HTTPClient.Timeout = 10
	standardClient := retryClient.StandardClient()

	req, err := standardClient.Get("http://localhost:3001/test/mock")

	if err != nil {
		log.Println(err)
	}
	defer req.Body.Close()
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))

}
