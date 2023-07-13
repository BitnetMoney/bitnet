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

package rpc

import (
	"net"
	"strings"

	"github.com/ethereum/go-ethereum/log"
)

// StartIPCEndpoint starts an IPC endpoint.
func StartIPCEndpoint(ipcEndpoint string, apis []API) (net.Listener, *Server, error) {
	// Register all the APIs exposed by the services.
	var (
		handler    = NewServer()
		regMap     = make(map[string]struct{})
		registered []string
	)
	for _, api := range apis {
		if err := handler.RegisterName(api.Namespace, api.Service); err != nil {
			log.Info("IPC registration failed", "namespace", api.Namespace, "error", err)
			return nil, nil, err
		}
		if _, ok := regMap[api.Namespace]; !ok {
			registered = append(registered, api.Namespace)
			regMap[api.Namespace] = struct{}{}
		}
	}
	log.Debug("IPCs registered", "namespaces", strings.Join(registered, ","))
	// All APIs registered, start the IPC listener.
	listener, err := ipcListen(ipcEndpoint)
	if err != nil {
		return nil, nil, err
	}
	go handler.ServeListener(listener)
	return listener, handler, nil
}
