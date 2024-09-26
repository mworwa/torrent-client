package torrentfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpen(t *testing.T) {
	torrentFile, _ := Open("testdata/file.torrent")

	assert.Equal(t, "http://bttracker.debian.org:6969/announce", torrentFile.Announce)
	assert.Equal(t, "\"Debian CD from cdimage.debian.org\"", torrentFile.Comment)
	assert.Equal(t, 1725106229, torrentFile.CreationDate)
	assert.Equal(t, 551858176, torrentFile.Info.Length)
	assert.Equal(t, "debian-12.7.0-arm64-netinst.iso", torrentFile.Info.Name)
	assert.Equal(t, 262144, torrentFile.Info.PieceLength)
	assert.NotEmpty(t, torrentFile.Info.Pieces)
}

func TestOpenInvalidPath(t *testing.T) {
	_, err := Open("invalidPath")

	assert.Error(t, err)
}
