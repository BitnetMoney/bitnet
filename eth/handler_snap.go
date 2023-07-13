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

package eth

import (
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth/protocols/snap"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

// snapHandler implements the snap.Backend interface to handle the various network
// packets that are sent as replies or broadcasts.
type snapHandler handler

func (h *snapHandler) Chain() *core.BlockChain { return h.chain }

// RunPeer is invoked when a peer joins on the `snap` protocol.
func (h *snapHandler) RunPeer(peer *snap.Peer, hand snap.Handler) error {
	return (*handler)(h).runSnapExtension(peer, hand)
}

// PeerInfo retrieves all known `snap` information about a peer.
func (h *snapHandler) PeerInfo(id enode.ID) interface{} {
	if p := h.peers.peer(id.String()); p != nil {
		if p.snapExt != nil {
			return p.snapExt.info()
		}
	}
	return nil
}

// Handle is invoked from a peer's message handler when it receives a new remote
// message that the handler couldn't consume and serve itself.
func (h *snapHandler) Handle(peer *snap.Peer, packet snap.Packet) error {
	return h.downloader.DeliverSnapPacket(peer, packet)
}
