package peers

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createClientAndServer(t *testing.T) (clientConn, serverConn net.Conn, ln net.Listener) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.Nil(t, err)

	// net.Dial does not block, so we need this signalling channel to make sure
	// we don't return before serverConn is ready
	done := make(chan struct{})
	go func() {
		serverConn, err = ln.Accept()
		require.Nil(t, err)
		done <- struct{}{}
	}()
	clientConn, err = net.Dial("tcp", ln.Addr().String())
	<-done

	fmt.Printf("Server creaated: %s", ln.Addr().String())

	return clientConn, serverConn, ln
}

func TestEstablishConnection(t *testing.T) {
	clientConn, serverConn, ln := createClientAndServer(t)
	defer ln.Close()
	defer clientConn.Close()
	defer serverConn.Close()

	_, portStr, err := net.SplitHostPort(serverConn.LocalAddr().String())
	require.Nil(t, err)
	var port uint16
	_, err = fmt.Sscanf(portStr, "%d", &port)
	require.Nil(t, err)

	peer := &Peer{
		IP:   net.IP{127, 0, 0, 1},
		Port: port,
	}

	conn, err := peer.establishConnection()
	assert.Implements(t, (*net.Conn)(nil), conn)
	assert.Nil(t, err)

	if conn != nil {
		defer conn.Close()
	}
}

func TestCreateHandshakeMessage(t *testing.T) {
	peer := Peer{}
	infoHash := [20]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14}
	peerID := [20]byte{0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30, 0x31, 0x32, 0x33, 0x34}

	expectedHandshake := []byte{
		19,                                                                                            // protocol length
		'B', 'i', 't', 'T', 'o', 'r', 'r', 'e', 'n', 't', ' ', 'p', 'r', 'o', 't', 'o', 'c', 'o', 'l', // protocol name
		0, 0, 0, 0, 0, 0, 0, 0, // reserved bytes
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, // infoHash
		0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30, 0x31, 0x32, 0x33, 0x34, // peerID
	}

	handshake := peer.createHandshakeMessage(infoHash, peerID)

	assert.Equal(t, expectedHandshake, handshake, "The handshake message should match the expected byte array")
}

func TestPerfromHandshake(t *testing.T) {
	clientConn, serverConn, ln := createClientAndServer(t)
	defer ln.Close()
	defer clientConn.Close()
	defer serverConn.Close()

	infoHash := [20]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14}
	peerID := [20]byte{0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30, 0x31, 0x32, 0x33, 0x34}
	peer := Peer{}
	handshakeMessage := peer.createHandshakeMessage(infoHash, peerID)

	go func() {
		// Read handshake from client
		buffer := make([]byte, len(handshakeMessage))
		_, err := serverConn.Read(buffer)
		require.Nil(t, err)

		// Write the same handshake back to simulate a proper response
		_, err = serverConn.Write(buffer)
		require.Nil(t, err)
	}()
	err := peer.performHandshake(clientConn, handshakeMessage)

	assert.Nil(t, err)
}
