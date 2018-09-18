package demo

import (
	"adexchange"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// ServerAdd is the server address
const ServerAdd = ":8888"

// HTTPConnector is to demonstrate SupplierConnector implementation
type HTTPConnector struct {
	processors []adexchange.DealProcessor
	httpServer *http.Server
}

func (con *HTTPConnector) createAdsRequest(req *http.Request) adexchange.AdsRequest {
	return req.RequestURI
}

func (con *HTTPConnector) response(resps []adexchange.AdsResponse, res http.ResponseWriter) {
	combinedResp := ""
	for _, resp := range resps {
		if str, ok := resp.(string); ok {
			combinedResp = combinedResp + " " + str
		}
	}
	fmt.Fprint(res, combinedResp)
}

func (con *HTTPConnector) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	log.Println("request is coming")
	adsReq := con.createAdsRequest(req)
	resps := make([]adexchange.AdsResponse, 0)
	for _, processor := range con.processors {
		resp, _ := processor.Process(adsReq) //for Demo, here, error has been ignored
		resps = append(resps, resp)
	}
	con.response(resps, res)
}

// RegisterProcessor is to implements SupplierConnector interface
func (con *HTTPConnector) RegisterProcessor(processors ...adexchange.DealProcessor) error {
	log.Println("Register a deal processor")
	for _, processor := range processors {
		con.processors = append(con.processors, processor)
	}
	return nil
}

// Start is to implements SupplierConnector interface
func (con *HTTPConnector) Start(ctx context.Context) error {
	con.httpServer.ListenAndServe()
	return nil
}

// Stop is to implements SupplierConnector interface
func (con *HTTPConnector) Stop(ctx context.Context) error {
	return con.httpServer.Close()
}

// CreateHTTPConnector is to create Demo SupplierConnector
func CreateHTTPConnector() *HTTPConnector {
	httpConnector := HTTPConnector{}
	server := &http.Server{
		Addr:           ServerAdd,
		Handler:        &httpConnector,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	httpConnector.httpServer = server
	return &httpConnector
}
