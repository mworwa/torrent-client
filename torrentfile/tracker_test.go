package torrentfile

import (
	"testing"

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
