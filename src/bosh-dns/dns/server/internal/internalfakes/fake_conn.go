// Code generated by counterfeiter. DO NOT EDIT.
package internalfakes

import (
	"net"
	"sync"
	"time"
)

type FakeConn struct {
	CloseStub        func() error
	closeMutex       sync.RWMutex
	closeArgsForCall []struct {
	}
	closeReturns struct {
		result1 error
	}
	closeReturnsOnCall map[int]struct {
		result1 error
	}
	LocalAddrStub        func() net.Addr
	localAddrMutex       sync.RWMutex
	localAddrArgsForCall []struct {
	}
	localAddrReturns struct {
		result1 net.Addr
	}
	localAddrReturnsOnCall map[int]struct {
		result1 net.Addr
	}
	ReadStub        func([]byte) (int, error)
	readMutex       sync.RWMutex
	readArgsForCall []struct {
		arg1 []byte
	}
	readReturns struct {
		result1 int
		result2 error
	}
	readReturnsOnCall map[int]struct {
		result1 int
		result2 error
	}
	RemoteAddrStub        func() net.Addr
	remoteAddrMutex       sync.RWMutex
	remoteAddrArgsForCall []struct {
	}
	remoteAddrReturns struct {
		result1 net.Addr
	}
	remoteAddrReturnsOnCall map[int]struct {
		result1 net.Addr
	}
	SetDeadlineStub        func(time.Time) error
	setDeadlineMutex       sync.RWMutex
	setDeadlineArgsForCall []struct {
		arg1 time.Time
	}
	setDeadlineReturns struct {
		result1 error
	}
	setDeadlineReturnsOnCall map[int]struct {
		result1 error
	}
	SetReadDeadlineStub        func(time.Time) error
	setReadDeadlineMutex       sync.RWMutex
	setReadDeadlineArgsForCall []struct {
		arg1 time.Time
	}
	setReadDeadlineReturns struct {
		result1 error
	}
	setReadDeadlineReturnsOnCall map[int]struct {
		result1 error
	}
	SetWriteDeadlineStub        func(time.Time) error
	setWriteDeadlineMutex       sync.RWMutex
	setWriteDeadlineArgsForCall []struct {
		arg1 time.Time
	}
	setWriteDeadlineReturns struct {
		result1 error
	}
	setWriteDeadlineReturnsOnCall map[int]struct {
		result1 error
	}
	WriteStub        func([]byte) (int, error)
	writeMutex       sync.RWMutex
	writeArgsForCall []struct {
		arg1 []byte
	}
	writeReturns struct {
		result1 int
		result2 error
	}
	writeReturnsOnCall map[int]struct {
		result1 int
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeConn) Close() error {
	fake.closeMutex.Lock()
	ret, specificReturn := fake.closeReturnsOnCall[len(fake.closeArgsForCall)]
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct {
	}{})
	stub := fake.CloseStub
	fakeReturns := fake.closeReturns
	fake.recordInvocation("Close", []interface{}{})
	fake.closeMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeConn) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

func (fake *FakeConn) CloseCalls(stub func() error) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = stub
}

func (fake *FakeConn) CloseReturns(result1 error) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = nil
	fake.closeReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeConn) CloseReturnsOnCall(i int, result1 error) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = nil
	if fake.closeReturnsOnCall == nil {
		fake.closeReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.closeReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeConn) LocalAddr() net.Addr {
	fake.localAddrMutex.Lock()
	ret, specificReturn := fake.localAddrReturnsOnCall[len(fake.localAddrArgsForCall)]
	fake.localAddrArgsForCall = append(fake.localAddrArgsForCall, struct {
	}{})
	stub := fake.LocalAddrStub
	fakeReturns := fake.localAddrReturns
	fake.recordInvocation("LocalAddr", []interface{}{})
	fake.localAddrMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeConn) LocalAddrCallCount() int {
	fake.localAddrMutex.RLock()
	defer fake.localAddrMutex.RUnlock()
	return len(fake.localAddrArgsForCall)
}

func (fake *FakeConn) LocalAddrCalls(stub func() net.Addr) {
	fake.localAddrMutex.Lock()
	defer fake.localAddrMutex.Unlock()
	fake.LocalAddrStub = stub
}

func (fake *FakeConn) LocalAddrReturns(result1 net.Addr) {
	fake.localAddrMutex.Lock()
	defer fake.localAddrMutex.Unlock()
	fake.LocalAddrStub = nil
	fake.localAddrReturns = struct {
		result1 net.Addr
	}{result1}
}

func (fake *FakeConn) LocalAddrReturnsOnCall(i int, result1 net.Addr) {
	fake.localAddrMutex.Lock()
	defer fake.localAddrMutex.Unlock()
	fake.LocalAddrStub = nil
	if fake.localAddrReturnsOnCall == nil {
		fake.localAddrReturnsOnCall = make(map[int]struct {
			result1 net.Addr
		})
	}
	fake.localAddrReturnsOnCall[i] = struct {
		result1 net.Addr
	}{result1}
}

func (fake *FakeConn) Read(arg1 []byte) (int, error) {
	var arg1Copy []byte
	if arg1 != nil {
		arg1Copy = make([]byte, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.readMutex.Lock()
	ret, specificReturn := fake.readReturnsOnCall[len(fake.readArgsForCall)]
	fake.readArgsForCall = append(fake.readArgsForCall, struct {
		arg1 []byte
	}{arg1Copy})
	stub := fake.ReadStub
	fakeReturns := fake.readReturns
	fake.recordInvocation("Read", []interface{}{arg1Copy})
	fake.readMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeConn) ReadCallCount() int {
	fake.readMutex.RLock()
	defer fake.readMutex.RUnlock()
	return len(fake.readArgsForCall)
}

func (fake *FakeConn) ReadCalls(stub func([]byte) (int, error)) {
	fake.readMutex.Lock()
	defer fake.readMutex.Unlock()
	fake.ReadStub = stub
}

func (fake *FakeConn) ReadArgsForCall(i int) []byte {
	fake.readMutex.RLock()
	defer fake.readMutex.RUnlock()
	argsForCall := fake.readArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeConn) ReadReturns(result1 int, result2 error) {
	fake.readMutex.Lock()
	defer fake.readMutex.Unlock()
	fake.ReadStub = nil
	fake.readReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeConn) ReadReturnsOnCall(i int, result1 int, result2 error) {
	fake.readMutex.Lock()
	defer fake.readMutex.Unlock()
	fake.ReadStub = nil
	if fake.readReturnsOnCall == nil {
		fake.readReturnsOnCall = make(map[int]struct {
			result1 int
			result2 error
		})
	}
	fake.readReturnsOnCall[i] = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeConn) RemoteAddr() net.Addr {
	fake.remoteAddrMutex.Lock()
	ret, specificReturn := fake.remoteAddrReturnsOnCall[len(fake.remoteAddrArgsForCall)]
	fake.remoteAddrArgsForCall = append(fake.remoteAddrArgsForCall, struct {
	}{})
	stub := fake.RemoteAddrStub
	fakeReturns := fake.remoteAddrReturns
	fake.recordInvocation("RemoteAddr", []interface{}{})
	fake.remoteAddrMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeConn) RemoteAddrCallCount() int {
	fake.remoteAddrMutex.RLock()
	defer fake.remoteAddrMutex.RUnlock()
	return len(fake.remoteAddrArgsForCall)
}

func (fake *FakeConn) RemoteAddrCalls(stub func() net.Addr) {
	fake.remoteAddrMutex.Lock()
	defer fake.remoteAddrMutex.Unlock()
	fake.RemoteAddrStub = stub
}

func (fake *FakeConn) RemoteAddrReturns(result1 net.Addr) {
	fake.remoteAddrMutex.Lock()
	defer fake.remoteAddrMutex.Unlock()
	fake.RemoteAddrStub = nil
	fake.remoteAddrReturns = struct {
		result1 net.Addr
	}{result1}
}

func (fake *FakeConn) RemoteAddrReturnsOnCall(i int, result1 net.Addr) {
	fake.remoteAddrMutex.Lock()
	defer fake.remoteAddrMutex.Unlock()
	fake.RemoteAddrStub = nil
	if fake.remoteAddrReturnsOnCall == nil {
		fake.remoteAddrReturnsOnCall = make(map[int]struct {
			result1 net.Addr
		})
	}
	fake.remoteAddrReturnsOnCall[i] = struct {
		result1 net.Addr
	}{result1}
}

func (fake *FakeConn) SetDeadline(arg1 time.Time) error {
	fake.setDeadlineMutex.Lock()
	ret, specificReturn := fake.setDeadlineReturnsOnCall[len(fake.setDeadlineArgsForCall)]
	fake.setDeadlineArgsForCall = append(fake.setDeadlineArgsForCall, struct {
		arg1 time.Time
	}{arg1})
	stub := fake.SetDeadlineStub
	fakeReturns := fake.setDeadlineReturns
	fake.recordInvocation("SetDeadline", []interface{}{arg1})
	fake.setDeadlineMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeConn) SetDeadlineCallCount() int {
	fake.setDeadlineMutex.RLock()
	defer fake.setDeadlineMutex.RUnlock()
	return len(fake.setDeadlineArgsForCall)
}

func (fake *FakeConn) SetDeadlineCalls(stub func(time.Time) error) {
	fake.setDeadlineMutex.Lock()
	defer fake.setDeadlineMutex.Unlock()
	fake.SetDeadlineStub = stub
}

func (fake *FakeConn) SetDeadlineArgsForCall(i int) time.Time {
	fake.setDeadlineMutex.RLock()
	defer fake.setDeadlineMutex.RUnlock()
	argsForCall := fake.setDeadlineArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeConn) SetDeadlineReturns(result1 error) {
	fake.setDeadlineMutex.Lock()
	defer fake.setDeadlineMutex.Unlock()
	fake.SetDeadlineStub = nil
	fake.setDeadlineReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeConn) SetDeadlineReturnsOnCall(i int, result1 error) {
	fake.setDeadlineMutex.Lock()
	defer fake.setDeadlineMutex.Unlock()
	fake.SetDeadlineStub = nil
	if fake.setDeadlineReturnsOnCall == nil {
		fake.setDeadlineReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.setDeadlineReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeConn) SetReadDeadline(arg1 time.Time) error {
	fake.setReadDeadlineMutex.Lock()
	ret, specificReturn := fake.setReadDeadlineReturnsOnCall[len(fake.setReadDeadlineArgsForCall)]
	fake.setReadDeadlineArgsForCall = append(fake.setReadDeadlineArgsForCall, struct {
		arg1 time.Time
	}{arg1})
	stub := fake.SetReadDeadlineStub
	fakeReturns := fake.setReadDeadlineReturns
	fake.recordInvocation("SetReadDeadline", []interface{}{arg1})
	fake.setReadDeadlineMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeConn) SetReadDeadlineCallCount() int {
	fake.setReadDeadlineMutex.RLock()
	defer fake.setReadDeadlineMutex.RUnlock()
	return len(fake.setReadDeadlineArgsForCall)
}

func (fake *FakeConn) SetReadDeadlineCalls(stub func(time.Time) error) {
	fake.setReadDeadlineMutex.Lock()
	defer fake.setReadDeadlineMutex.Unlock()
	fake.SetReadDeadlineStub = stub
}

func (fake *FakeConn) SetReadDeadlineArgsForCall(i int) time.Time {
	fake.setReadDeadlineMutex.RLock()
	defer fake.setReadDeadlineMutex.RUnlock()
	argsForCall := fake.setReadDeadlineArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeConn) SetReadDeadlineReturns(result1 error) {
	fake.setReadDeadlineMutex.Lock()
	defer fake.setReadDeadlineMutex.Unlock()
	fake.SetReadDeadlineStub = nil
	fake.setReadDeadlineReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeConn) SetReadDeadlineReturnsOnCall(i int, result1 error) {
	fake.setReadDeadlineMutex.Lock()
	defer fake.setReadDeadlineMutex.Unlock()
	fake.SetReadDeadlineStub = nil
	if fake.setReadDeadlineReturnsOnCall == nil {
		fake.setReadDeadlineReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.setReadDeadlineReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeConn) SetWriteDeadline(arg1 time.Time) error {
	fake.setWriteDeadlineMutex.Lock()
	ret, specificReturn := fake.setWriteDeadlineReturnsOnCall[len(fake.setWriteDeadlineArgsForCall)]
	fake.setWriteDeadlineArgsForCall = append(fake.setWriteDeadlineArgsForCall, struct {
		arg1 time.Time
	}{arg1})
	stub := fake.SetWriteDeadlineStub
	fakeReturns := fake.setWriteDeadlineReturns
	fake.recordInvocation("SetWriteDeadline", []interface{}{arg1})
	fake.setWriteDeadlineMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeConn) SetWriteDeadlineCallCount() int {
	fake.setWriteDeadlineMutex.RLock()
	defer fake.setWriteDeadlineMutex.RUnlock()
	return len(fake.setWriteDeadlineArgsForCall)
}

func (fake *FakeConn) SetWriteDeadlineCalls(stub func(time.Time) error) {
	fake.setWriteDeadlineMutex.Lock()
	defer fake.setWriteDeadlineMutex.Unlock()
	fake.SetWriteDeadlineStub = stub
}

func (fake *FakeConn) SetWriteDeadlineArgsForCall(i int) time.Time {
	fake.setWriteDeadlineMutex.RLock()
	defer fake.setWriteDeadlineMutex.RUnlock()
	argsForCall := fake.setWriteDeadlineArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeConn) SetWriteDeadlineReturns(result1 error) {
	fake.setWriteDeadlineMutex.Lock()
	defer fake.setWriteDeadlineMutex.Unlock()
	fake.SetWriteDeadlineStub = nil
	fake.setWriteDeadlineReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeConn) SetWriteDeadlineReturnsOnCall(i int, result1 error) {
	fake.setWriteDeadlineMutex.Lock()
	defer fake.setWriteDeadlineMutex.Unlock()
	fake.SetWriteDeadlineStub = nil
	if fake.setWriteDeadlineReturnsOnCall == nil {
		fake.setWriteDeadlineReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.setWriteDeadlineReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeConn) Write(arg1 []byte) (int, error) {
	var arg1Copy []byte
	if arg1 != nil {
		arg1Copy = make([]byte, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.writeMutex.Lock()
	ret, specificReturn := fake.writeReturnsOnCall[len(fake.writeArgsForCall)]
	fake.writeArgsForCall = append(fake.writeArgsForCall, struct {
		arg1 []byte
	}{arg1Copy})
	stub := fake.WriteStub
	fakeReturns := fake.writeReturns
	fake.recordInvocation("Write", []interface{}{arg1Copy})
	fake.writeMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeConn) WriteCallCount() int {
	fake.writeMutex.RLock()
	defer fake.writeMutex.RUnlock()
	return len(fake.writeArgsForCall)
}

func (fake *FakeConn) WriteCalls(stub func([]byte) (int, error)) {
	fake.writeMutex.Lock()
	defer fake.writeMutex.Unlock()
	fake.WriteStub = stub
}

func (fake *FakeConn) WriteArgsForCall(i int) []byte {
	fake.writeMutex.RLock()
	defer fake.writeMutex.RUnlock()
	argsForCall := fake.writeArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeConn) WriteReturns(result1 int, result2 error) {
	fake.writeMutex.Lock()
	defer fake.writeMutex.Unlock()
	fake.WriteStub = nil
	fake.writeReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeConn) WriteReturnsOnCall(i int, result1 int, result2 error) {
	fake.writeMutex.Lock()
	defer fake.writeMutex.Unlock()
	fake.WriteStub = nil
	if fake.writeReturnsOnCall == nil {
		fake.writeReturnsOnCall = make(map[int]struct {
			result1 int
			result2 error
		})
	}
	fake.writeReturnsOnCall[i] = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeConn) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	fake.localAddrMutex.RLock()
	defer fake.localAddrMutex.RUnlock()
	fake.readMutex.RLock()
	defer fake.readMutex.RUnlock()
	fake.remoteAddrMutex.RLock()
	defer fake.remoteAddrMutex.RUnlock()
	fake.setDeadlineMutex.RLock()
	defer fake.setDeadlineMutex.RUnlock()
	fake.setReadDeadlineMutex.RLock()
	defer fake.setReadDeadlineMutex.RUnlock()
	fake.setWriteDeadlineMutex.RLock()
	defer fake.setWriteDeadlineMutex.RUnlock()
	fake.writeMutex.RLock()
	defer fake.writeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeConn) recordInvocation(key string, args []interface{}) {
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

var _ net.Conn = new(FakeConn)
