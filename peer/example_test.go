// Copyright (c) 2015 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package peer_test

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/peer"
	"github.com/btcsuite/btcd/wire"
)

// mockRemotePeer creates a basic inbound peer listening on the simnet port for
// use with Example_peerConnection.  It does not return until the listner is
// active.
func mockRemotePeer() error {
	// Configure peer to act as a simnet node that offers no services.
	peerCfg := &peer.Config{
		UserAgentName:    "peer",  // User agent name to advertise.
		UserAgentVersion: "1.0.0", // User agent version to advertise.
		ChainParams:      &chaincfg.SimNetParams,
	}

	// Accept connections on the simnet port.
	listener, err := net.Listen("tcp", "127.0.0.1:18555")
	if err != nil {
		return err
	}
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Accept: error %v\n", err)
			return
		}

		// Create and start the inbound peer.
		if _, err := peer.NewInboundPeer(peerCfg, conn); err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}

// This example demonstrates the basic process for initializing and creating an
// outbound peer.  Peers negotiate by exchanging version and verack messages.
// For demonstration, a simple handler for version message is attached to the
// peer.
func Example_newOutboundPeer() {
	// Ordinarily this will not be needed since the outbound peer will be
	// connecting to a remote peer, however, since this example is executed
	// and tested, a mock remote peer is needed to listen for the outbound
	// peer.
	if err := mockRemotePeer(); err != nil {
		fmt.Printf("mockRemotePeer: unexpected error %v\n", err)
		return
	}

	// Create an outbound peer that is configured to act as a simnet node
	// that offers no services and has listeners for the version and verack
	// messages.  The verack listener is used here to signal the code below
	// when the handshake has been finished by signalling a channel.
	verack := make(chan struct{}, 1)
	peerCfg := &peer.Config{
		UserAgentName:    "peer",  // User agent name to advertise.
		UserAgentVersion: "1.0.0", // User agent version to advertise.
		ChainParams:      &chaincfg.SimNetParams,
		Services:         0,
		Listeners: peer.MessageListeners{
			OnVersion: func(p *peer.Peer, msg *wire.MsgVersion) {
				fmt.Println("outbound: received version")
			},
			OnVerAck: func(p *peer.Peer, msg *wire.MsgVerAck) {
				verack <- struct{}{}
			},
		},
	}

	// Establish the connection to the peer address and mark it connected.
	conn, err := net.Dial("tcp", "127.0.0.1:18555")
	if err != nil {
		fmt.Printf("net.Dial: error %v\n", err)
		return
	}
	p, err := peer.NewOutboundPeer(peerCfg, conn)
	if err != nil {
		fmt.Printf("NewOutboundPeer: error %v\n", err)
		return
	}

	// Wait for the verack message or timeout in case of failure.
	select {
	case <-verack:
	case <-time.After(time.Second * 1):
		fmt.Printf("Example_peerConnection: verack timeout")
	}

	// Shutdown the peer.
	p.Shutdown()

	// Output:
	// outbound: received version
}
