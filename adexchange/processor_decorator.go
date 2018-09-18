package adexchange

import "github.com/easierway/service_decorators"

// DecoratedDealProcessor is the deal processor instance decorated with service deocorators
type DecoratedDealProcessor struct {
	serviceDecorators []service_decorators.Decorator
	orgProcessor      DealProcessor
	decFn             service_decorators.ServiceFunc
}

func (dp *DecoratedDealProcessor) decorateOrgProcessor() {
	dp.decFn = func(req service_decorators.Request) (service_decorators.Response, error) {
		return dp.orgProcessor.Process(req)
	}

	for _, dec := range dp.serviceDecorators {
		dp.decFn = dec.Decorate(dp.decFn)
	}
}

// Process is to implement the DealProcessor interface.
func (dp *DecoratedDealProcessor) Process(req AdsRequest) (AdsResponse, error) {
	return dp.decFn(req)
}

// DecorateDealProcessor is to leverage the project "https://github.com/easierway/service_decorators"
// By this method, you can decorate the deocrators on a deal processor instance
func DecorateDealProcessor(orgProcessor DealProcessor,
	serviceDecorators ...service_decorators.Decorator) DealProcessor {
	dp := DecoratedDealProcessor{orgProcessor: orgProcessor,
		serviceDecorators: serviceDecorators}
	dp.decorateOrgProcessor()
	return &dp
}
