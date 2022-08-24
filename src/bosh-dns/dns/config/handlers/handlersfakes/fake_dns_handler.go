// Code generated by counterfeiter. DO NOT EDIT.
package handlersfakes

import (
	"sync"

	"github.com/miekg/dns"
)

type FakeDnsHandler struct {
	ServeDNSStub        func(dns.ResponseWriter, *dns.Msg)
	serveDNSMutex       sync.RWMutex
	serveDNSArgsForCall []struct {
		arg1 dns.ResponseWriter
		arg2 *dns.Msg
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeDnsHandler) ServeDNS(arg1 dns.ResponseWriter, arg2 *dns.Msg) {
	fake.serveDNSMutex.Lock()
	fake.serveDNSArgsForCall = append(fake.serveDNSArgsForCall, struct {
		arg1 dns.ResponseWriter
		arg2 *dns.Msg
	}{arg1, arg2})
	stub := fake.ServeDNSStub
	fake.recordInvocation("ServeDNS", []interface{}{arg1, arg2})
	fake.serveDNSMutex.Unlock()
	if stub != nil {
		fake.ServeDNSStub(arg1, arg2)
	}
}

func (fake *FakeDnsHandler) ServeDNSCallCount() int {
	fake.serveDNSMutex.RLock()
	defer fake.serveDNSMutex.RUnlock()
	return len(fake.serveDNSArgsForCall)
}

func (fake *FakeDnsHandler) ServeDNSCalls(stub func(dns.ResponseWriter, *dns.Msg)) {
	fake.serveDNSMutex.Lock()
	defer fake.serveDNSMutex.Unlock()
	fake.ServeDNSStub = stub
}

func (fake *FakeDnsHandler) ServeDNSArgsForCall(i int) (dns.ResponseWriter, *dns.Msg) {
	fake.serveDNSMutex.RLock()
	defer fake.serveDNSMutex.RUnlock()
	argsForCall := fake.serveDNSArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeDnsHandler) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.serveDNSMutex.RLock()
	defer fake.serveDNSMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeDnsHandler) recordInvocation(key string, args []interface{}) {
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
