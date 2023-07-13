// Copyright 2023 Bitnet
// This file is part of the Bitnet library.
//
// This software is provided "as is", without warranty of any kind,
// express or implied, including but not limited to the warranties
// of merchantability, fitness for a particular purpose and
// noninfringement. In no even shall the authors or copyright
// holders be liable for any claim, damages, or other liability,
// whether in an action of contract, tort or otherwise, arising
// from, out of or in connection with the software or the use or
// other dealings in the software.

package server

import (
	"net"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/les/utils"
	"github.com/ethereum/go-ethereum/les/vflux"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/rlp"
)

type (
	// Server serves vflux requests
	Server struct {
		limiter         *utils.Limiter
		lock            sync.Mutex
		services        map[string]*serviceEntry
		delayPerRequest time.Duration
	}

	// Service is a service registered at the Server and identified by a string id
	Service interface {
		Handle(id enode.ID, address string, name string, data []byte) []byte // never called concurrently
	}

	serviceEntry struct {
		id, desc string
		backend  Service
	}
)

// NewServer creates a new Server
func NewServer(delayPerRequest time.Duration) *Server {
	return &Server{
		limiter:         utils.NewLimiter(1000),
		delayPerRequest: delayPerRequest,
		services:        make(map[string]*serviceEntry),
	}
}

// Register registers a Service
func (s *Server) Register(b Service, id, desc string) {
	srv := &serviceEntry{backend: b, id: id, desc: desc}
	if strings.Contains(srv.id, ":") {
		// srv.id + ":" will be used as a service database prefix
		log.Error("Service ID contains ':'", "id", srv.id)
		return
	}
	s.lock.Lock()
	s.services[srv.id] = srv
	s.lock.Unlock()
}

// Serve serves a vflux request batch
// Note: requests are served by the Handle functions of the registered services. Serve
// may be called concurrently but the Handle functions are called sequentially and
// therefore thread safety is guaranteed.
func (s *Server) Serve(id enode.ID, address string, requests vflux.Requests) vflux.Replies {
	reqLen := uint(len(requests))
	if reqLen == 0 || reqLen > vflux.MaxRequestLength {
		return nil
	}
	// Note: the value parameter will be supplied by the token sale module (total amount paid)
	ch := <-s.limiter.Add(id, address, 0, reqLen)
	if ch == nil {
		return nil
	}
	// Note: the limiter ensures that the following section is not running concurrently,
	// the lock only protects against contention caused by new service registration
	s.lock.Lock()
	results := make(vflux.Replies, len(requests))
	for i, req := range requests {
		if service := s.services[req.Service]; service != nil {
			results[i] = service.backend.Handle(id, address, req.Name, req.Params)
		}
	}
	s.lock.Unlock()
	time.Sleep(s.delayPerRequest * time.Duration(reqLen))
	close(ch)
	return results
}

// ServeEncoded serves an encoded vflux request batch and returns the encoded replies
func (s *Server) ServeEncoded(id enode.ID, addr *net.UDPAddr, req []byte) []byte {
	var requests vflux.Requests
	if err := rlp.DecodeBytes(req, &requests); err != nil {
		return nil
	}
	results := s.Serve(id, addr.String(), requests)
	if results == nil {
		return nil
	}
	res, _ := rlp.EncodeToBytes(&results)
	return res
}

// Stop shuts down the server
func (s *Server) Stop() {
	s.limiter.Stop()
}
