package peers

import (
	"fmt"
	"net"
	"time"
)

func (p *Peer) Connect(infoHash [20]byte, peerID [20]byte) error {
	conn, err := p.establishConnection()
	if err != nil {
		return err
	}

	handshakeMessage := p.createHandshakeMessage(infoHash, peerID)
	err = performHandshake(conn, handshakeMessage)
	if err != nil {
		return nil
	}
	return nil
}

func (p *Peer) establishConnection() (net.Conn, error) {
	fmt.Printf("Connecting to peer: %s\n", p.String())
	conn, err := net.DialTimeout("tcp", p.String(), 3*time.Second)

	if err != nil {
		return nil, err
	}

	fmt.Println("Connection establised")
	return conn, nil
}

func (peer *Peer) createHandshakeMessage(infoHash [20]byte, peerID [20]byte) []byte {
	protocol := "BitTorrent protocol"
	protocolLen := byte(len(protocol))

	// The reserved bytes are typically 0
	reserved := make([]byte, 8)

	// Create the handshake message as a byte slice
	handshake := make([]byte, 49+len(protocol))
	handshake[0] = protocolLen
	copy(handshake[1:], protocol)
	copy(handshake[20:], reserved)
	copy(handshake[28:], infoHash[:])
	copy(handshake[48:], peerID[:])

	return handshake
}

func performHandshake(conn net.Conn, handshakeMessage []byte) error {
	fmt.Printf("Performing handshake for message: %s", string(handshakeMessage))
	_, err := conn.Write(handshakeMessage)
	if err != nil {
		return err
	}

	response := make([]byte, 68)
	_, err = conn.Read(response)
	if err != nil {
		return err
	}

	fmt.Println(response)

	return nil
}
