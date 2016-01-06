// This file was generated by counterfeiter
package fakes

import (
	"net"
	"sync"

	"github.com/cloudfoundry-incubator/diego-ssh/server"
)

type FakeConnectionHandler struct {
	HandleConnectionStub        func(net.Conn)
	handleConnectionMutex       sync.RWMutex
	handleConnectionArgsForCall []struct {
		arg1 net.Conn
	}
}

func (fake *FakeConnectionHandler) HandleConnection(arg1 net.Conn) {
	fake.handleConnectionMutex.Lock()
	fake.handleConnectionArgsForCall = append(fake.handleConnectionArgsForCall, struct {
		arg1 net.Conn
	}{arg1})
	fake.handleConnectionMutex.Unlock()
	if fake.HandleConnectionStub != nil {
		fake.HandleConnectionStub(arg1)
	}
}

func (fake *FakeConnectionHandler) HandleConnectionCallCount() int {
	fake.handleConnectionMutex.RLock()
	defer fake.handleConnectionMutex.RUnlock()
	return len(fake.handleConnectionArgsForCall)
}

func (fake *FakeConnectionHandler) HandleConnectionArgsForCall(i int) net.Conn {
	fake.handleConnectionMutex.RLock()
	defer fake.handleConnectionMutex.RUnlock()
	return fake.handleConnectionArgsForCall[i].arg1
}

var _ server.ConnectionHandler = new(FakeConnectionHandler)