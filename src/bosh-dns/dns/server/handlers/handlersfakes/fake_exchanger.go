// Code generated by counterfeiter. DO NOT EDIT.
package handlersfakes

import (
	"bosh-dns/dns/server/handlers"
	"sync"
	"time"

	"github.com/miekg/dns"
)

type FakeExchanger struct {
	ExchangeStub        func(*dns.Msg, string) (*dns.Msg, time.Duration, error)
	exchangeMutex       sync.RWMutex
	exchangeArgsForCall []struct {
		arg1 *dns.Msg
		arg2 string
	}
	exchangeReturns struct {
		result1 *dns.Msg
		result2 time.Duration
		result3 error
	}
	exchangeReturnsOnCall map[int]struct {
		result1 *dns.Msg
		result2 time.Duration
		result3 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeExchanger) Exchange(arg1 *dns.Msg, arg2 string) (*dns.Msg, time.Duration, error) {
	fake.exchangeMutex.Lock()
	ret, specificReturn := fake.exchangeReturnsOnCall[len(fake.exchangeArgsForCall)]
	fake.exchangeArgsForCall = append(fake.exchangeArgsForCall, struct {
		arg1 *dns.Msg
		arg2 string
	}{arg1, arg2})
	stub := fake.ExchangeStub
	fakeReturns := fake.exchangeReturns
	fake.recordInvocation("Exchange", []interface{}{arg1, arg2})
	fake.exchangeMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeExchanger) ExchangeCallCount() int {
	fake.exchangeMutex.RLock()
	defer fake.exchangeMutex.RUnlock()
	return len(fake.exchangeArgsForCall)
}

func (fake *FakeExchanger) ExchangeCalls(stub func(*dns.Msg, string) (*dns.Msg, time.Duration, error)) {
	fake.exchangeMutex.Lock()
	defer fake.exchangeMutex.Unlock()
	fake.ExchangeStub = stub
}

func (fake *FakeExchanger) ExchangeArgsForCall(i int) (*dns.Msg, string) {
	fake.exchangeMutex.RLock()
	defer fake.exchangeMutex.RUnlock()
	argsForCall := fake.exchangeArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeExchanger) ExchangeReturns(result1 *dns.Msg, result2 time.Duration, result3 error) {
	fake.exchangeMutex.Lock()
	defer fake.exchangeMutex.Unlock()
	fake.ExchangeStub = nil
	fake.exchangeReturns = struct {
		result1 *dns.Msg
		result2 time.Duration
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeExchanger) ExchangeReturnsOnCall(i int, result1 *dns.Msg, result2 time.Duration, result3 error) {
	fake.exchangeMutex.Lock()
	defer fake.exchangeMutex.Unlock()
	fake.ExchangeStub = nil
	if fake.exchangeReturnsOnCall == nil {
		fake.exchangeReturnsOnCall = make(map[int]struct {
			result1 *dns.Msg
			result2 time.Duration
			result3 error
		})
	}
	fake.exchangeReturnsOnCall[i] = struct {
		result1 *dns.Msg
		result2 time.Duration
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeExchanger) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.exchangeMutex.RLock()
	defer fake.exchangeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeExchanger) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ handlers.Exchanger = new(FakeExchanger)
