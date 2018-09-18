// Package adexchange is to define a extensiable architecture
// with a plugin mode.
package adexchange

import "context"

// AdsRequest supplier Ads request
type AdsRequest interface{}

// AdsResponse is the response of Ads request
type AdsResponse interface{}

// DealProcessor processes the request according to the deal type
type DealProcessor interface {
	Process(req AdsRequest) (AdsResponse, error)
}

// SupplierConnector is to receive the request from Supplier
// which provides the different accessing ways (eg. RPC, HTTP)
type SupplierConnector interface {
	Start(ctx context.Context) error
	// Register the processors to handle the Requests coming from suppliers.
	RegisterProcessor(processors ...DealProcessor) error
	Stop(ctx context.Context) error
}
