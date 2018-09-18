package adexchange

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/easierway/service_decorators"
)

type demoDealProcessor struct {
	executionCnt int
}

func (f *demoDealProcessor) Process(data AdsRequest) (AdsResponse, error) {
	var (
		str string
		ok  bool
	)
	if str, ok = data.(string); !ok {
		panic("invalid input")
	}
	fmt.Println("input:", str)
	f.executionCnt++
	if f.executionCnt < 3 {
		return nil, errors.New("error occurred")
	}
	return f.executionCnt, nil
}

func TestDecoratedDealProcessor(t *testing.T) {
	var (
		retryDec *service_decorators.RetryDecorator
		err      error
	)
	retriableChecker := func(err error) bool {
		return true
	}

	if retryDec, err = service_decorators.CreateRetryDecorator(5, /*max retry times*/
		time.Second*1, time.Second*2, retriableChecker); err != nil {
		panic(err)
	}

	orgProcessor := &demoDealProcessor{}
	decoratedProcessor := DecorateDealProcessor(orgProcessor, retryDec)
	startT := time.Now()
	ret, err := decoratedProcessor.Process("Hello")
	timeSpent := time.Now().Sub(startT).Nanoseconds()
	fmt.Println("time escaped: ", timeSpent)
	if err != nil {
		t.Error("Unexpected error occurred.")
	}
	if ret.(int) != 3 {
		t.Errorf("Expected value is %d, but actual value is %d", 3, ret)
	}
}

func ExampleDecoratedDealProcessor() {
	// 1. Create the decorators
	var (
		retryDec *service_decorators.RetryDecorator
		err      error
	)
	retriableChecker := func(err error) bool {
		return true
	}

	if retryDec, err = service_decorators.CreateRetryDecorator(5, /*max retry times*/
		time.Second*1, time.Second*2, retriableChecker); err != nil {
		panic(err)
	}

	// 2. Decorate the deal processor with the decorators by calling DecorateFilter
	// Be careful of the order of the decorators, your processor will be decorated by the same order
	// and the decorators will be invoked by the same order
	// after that you will get a deocrated DealProcessor instance
	orgProcessor := &demoDealProcessor{}
	decoratedProcessor := DecorateDealProcessor(orgProcessor, retryDec)

	// 3. The decorated instance can be used as a normal Fiter instance
	decoratedProcessor.Process("Hello")
}
