package torrentfile

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mworwa/bittorrent/peers"
	"github.com/stretchr/testify/assert"
)

func TestBuildTrackerURL(t *testing.T) {
	torrentFile := TorrentFile{
		Announce: "https://google.com",
		Info: Info{
			Length: 10,
		},
	}

	var peerID [20]byte
	url, _ := torrentFile.buildTrackerURL(peerID, 1234)
	assert.Equal(t, "https://google.com?compact=1&downloaded=0&info_hash=%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00&left=10&peer_id=%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00&port=1234&uploaded=0", url)
}

func TestRequestPeer(t *testing.T) {
	var peerID [20]byte
	port := uint16(6881)

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := []byte(
			"d" +
				"8:interval" + "i900e" +
				"5:peers" + "12:" +
				string([]byte{
					192, 0, 2, 123, 0x1A, 0xE1, // 0x1AE1 = 6881
					127, 0, 0, 1, 0x1A, 0xE9, // 0x1AE9 = 6889
				}) + "e")
		w.Write(response)
	}))

	defer mockServer.Close()

	torrentFile := &TorrentFile{
		Announce: mockServer.URL,
	}

	requestedPeers, err := torrentFile.requestPeers(peerID, port)

	expected := []peers.Peer{
		{IP: net.IP{192, 0, 2, 123}, Port: uint16(6881)},
		{IP: net.IP{127, 0, 0, 1}, Port: uint16(6889)},
	}

	assert.Nil(t, err)
	assert.Equal(t, expected, requestedPeers)
}
