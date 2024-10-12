package peers

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendMessage(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	go func() {
		err := SendMessage(client, PeerMessage{Id: MessageChoke})
		assert.NoError(t, err)
	}()

	expectedMessage := []byte{0, 0, 0, 1, byte(MessageChoke)}
	recivedMessage := make([]byte, len(expectedMessage))

	_, err := server.Read(recivedMessage)

	assert.NoError(t, err)
	assert.Equal(t, expectedMessage, recivedMessage)
}

func TestReadMessage(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	go func() {
		message := []byte{0, 0, 0, 1, byte(MessageUnchoke)}
		_, err := server.Write(message)
		assert.NoError(t, err)
	}()

	peerMessage, err := ReadMessage(client)

	assert.NoError(t, err)
	assert.Equal(t, MessageUnchoke, peerMessage.Id)
	assert.Equal(t, []byte{byte(MessageUnchoke)}, peerMessage.Payload)
}
