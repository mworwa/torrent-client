package client

import (
	"crypto/rand"
	"fmt"
	"github.com/mworwa/bittorrent/bitfield"
	"github.com/mworwa/bittorrent/peers"
	"github.com/mworwa/bittorrent/torrentfile"
	"math"
	"os"
)

type Client struct {
	file   *torrentfile.TorrentFile
	peers  []peers.Peer
	peerID []byte
}

type pieceWork struct {
	index  int
	hash   [20]byte
	lenght int
}

const BlockSize = 16384

func (c *Client) DownloadTorrent(torrentFilePath string, targetFilePath string) error {
	torrentFile, err := torrentfile.Open(torrentFilePath)
	if err != nil {
		return err
	}

	var peerID [20]byte
	_, err = rand.Read(peerID[:])
	if err != nil {
		return err
	}

	requestedPeers, err := torrentFile.RequestPeers(peerID, 6881)
	if err != nil {
		return err
	}
	for _, requestedPeer := range requestedPeers {
		conn, err := requestedPeer.Connect(torrentFile.InfoHash, peerID)
		defer conn.Close()
		if err != nil {
			continue
		}

		peerMessage, err := peers.ReadMessage(conn)
		if err != nil {
			return err
		}

		if peerMessage.Id != peers.MessageBitfield {
			return fmt.Errorf("expecting bitfield message got: %d", peerMessage.Id)
		}

		peerBitfield := bitfield.New(peerMessage.Payload)
		// localBitfield := bitfield.NewEmpty(torrentFile.Info.Length, torrentFile.Info.PieceLength)
		fmt.Printf("%0b\n", peerBitfield)

		peers.SendUnchoke(conn)
		peers.SendInterested(conn)
		peerMessage, err = peers.ReadMessage(conn)
		if err != nil {
			return err
		}

		if peerMessage.Id != peers.MessageUnchoke {
			return fmt.Errorf("expecting unchoke message got: %d", peerMessage.Id)
		}

		pieceLenght := torrentFile.Info.PieceLength

		// pieceData := make([]byte, pieceLenght)
		numBlocks := int(math.Ceil(float64(pieceLenght) / float64(BlockSize)))

		for blockIndex := 0; blockIndex < numBlocks; blockIndex++ {
			blockOffset := blockIndex * BlockSize
			blockSize := BlockSize

			if blockOffset+blockSize > pieceLenght {
				blockSize = pieceLenght - blockOffset
			}

			err := peers.SendRequest(conn, 0, 0, uint32(blockSize))
			if err != nil {
				return err
			}

			// Extract the message type (assuming the message ID is at byte 0)
			peerMessage, err = peers.ReadMessage(conn)
			if err != nil {
				return err
			}

			if peerMessage.Id != peers.MessagePiece {
				return fmt.Errorf("expecting piece message got: %d", peerMessage.Id)
			}
			fmt.Println(peerMessage)
			os.Exit(1)

		}
	}
	return nil
}
