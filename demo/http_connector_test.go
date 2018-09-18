package demo

import (
	"adexchange"
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"testing"
	"time"
)

type DemoDealProcessor struct {
	resp string
}

func (processor *DemoDealProcessor) Process(req adexchange.AdsRequest) (
	adexchange.AdsResponse, error) {
	log.Println("process the Ads request")
	if reqStr, ok := req.(string); ok {
		return processor.resp + " Orignial request is " + reqStr, nil
	} else {
		panic("invalid request [DemoDealProcessor]")
	}

}

func askAdsFromAdx() string {

	const (
		MaxIdleConns        int = 100
		MaxIdleConnsPerHost int = 100
		IdleConnTimeout     int = 90
	)

	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        MaxIdleConns,
			MaxIdleConnsPerHost: MaxIdleConnsPerHost,
			IdleConnTimeout:     time.Duration(IdleConnTimeout) * time.Second,
		},
		Timeout: 20 * time.Second,
	}
	var endPoint = "http://localhost:8888"

	req, err := http.NewRequest("POST", endPoint, bytes.NewBuffer([]byte("Post this data")))
	if err != nil {
		log.Fatalf("Error Occured. %+v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// use httpClient to send request
	response, err := httpClient.Do(req)
	var ret string
	if err != nil && response == nil {
		log.Fatalf("Error sending request to API endpoint. %+v", err)
	} else {
		// Close the connection to reuse it
		defer response.Body.Close()

		// Let's check if the work actually is done
		// We have seen inconsistencies even when we get 200 OK response
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatalf("Couldn't parse response body. %+v", err)
		}
		ret = string(body)
		log.Println("Response Body:", ret)
	}
	return ret
}

func Example() {
	httpConnector := CreateHTTPConnector()
	dealProcessor := DemoDealProcessor{"Processed by DemoDealProcessor."}
	dealDispatcher := DealDispatcher{
		[]adexchange.DealProcessor{&dealProcessor},
	}
	ctx, cancel := context.WithCancel(context.Background())
	httpConnector.RegisterProcessor(&dealDispatcher)
	log.Println("Http connector is started.")
	go httpConnector.Start(ctx)
	askAdsFromAdx()
	time.Sleep(time.Second * 1)
	cancel()
	log.Println("Http connector is stoped")
	httpConnector.Stop(ctx)
	// ouput:
	// Register a deal processor
	// Http connector is started.
	// request is coming
	// select a proper deal processor
	// process the Ads request
	// Response Body:  Processed by DemoDealProcessor. Orignial request is /
	// Http connector is stoped
}

func TestDemo(t *testing.T) {
	Example()
}
