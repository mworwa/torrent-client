package peers

import (
	"encoding/binary"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	t.Run("Valid peers data", func(t *testing.T) {
		peersBin := make([]byte, 12)
		copy(peersBin[0:4], []byte{192, 168, 1, 1})
		binary.BigEndian.PutUint16(peersBin[4:6], 6881)

		copy(peersBin[6:10], []byte{10, 0, 0, 1})
		binary.BigEndian.PutUint16(peersBin[10:12], 6882)
		fmt.Println(peersBin)
		expectedPeers := []Peer{
			{IP: net.IP{192, 168, 1, 1}, Port: 6881},
			{IP: net.IP{10, 0, 0, 1}, Port: 6882},
		}

		peers, _ := Unmarshal(peersBin)

		assert.Equal(t, expectedPeers, peers)
	})
	t.Run("Malformed peers data", func(t *testing.T) {
		peersBin := make([]byte, 5)
		copy(peersBin[0:4], []byte{192, 168, 1, 1, 66})

		_, err := Unmarshal(peersBin)

		assert.Error(t, err)
	})
}

func TestString(t *testing.T) {
	peer := Peer{
		IP:   net.IP{192, 168, 1, 1},
		Port: 6881,
	}

	expeected := "192.168.1.1:6881"
	actual := peer.String()

	assert.Equal(t, expeected, actual)

}
