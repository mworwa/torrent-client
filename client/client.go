package client

import (
	"github.com/mworwa/bittorrent/peers"
	"github.com/mworwa/bittorrent/torrentfile"
)

type Client struct {
	file  *torrentfile.TorrentFile
	peers []peers.Peer
}

func (c *Client) DownloadTorrent(torrentFilePath string, targetFilePath string) error {
	torrentFile, err := torrentfile.Open(torrentFilePath)

	if err != nil {
		return err
	}

	err = torrentFile.DownloadToFile(targetFilePath)
	if err != nil {
		return err
	}

	return nil
}
