package demo

import (
	"adexchange"
	"log"
)

// DealDispatcher is the demo implementation of deal dispatches
// which is to dispatch the Ads request to the different deal processor
// for handling the deal with the different deal protocol, such as, RTB.
type DealDispatcher struct {
	processors []adexchange.DealProcessor
}

func (disp *DealDispatcher) selectDealProcessor(req adexchange.AdsRequest) adexchange.DealProcessor {
	return disp.processors[0]
}

// Process is to implement DealProcessor interface.
// DealDispatcher is also a DealProcessor (as Proxy pattern)
func (disp *DealDispatcher) Process(req adexchange.AdsRequest) (adexchange.AdsResponse, error) {
	log.Println("select a proper deal processor")
	return disp.selectDealProcessor(req).Process(req)
}
