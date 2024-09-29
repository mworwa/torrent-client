package torrentfile

import (
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/mworwa/bittorrent/bencode"
	"github.com/mworwa/bittorrent/peers"
)

func (t *TorrentFile) buildTrackerURL(peerID [20]byte, port uint16) (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil {
		return "", err
	}

	requestParameters := url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.Info.Length)},
	}
	base.RawQuery = requestParameters.Encode()

	return base.String(), nil
}

func (t *TorrentFile) requestPeers(peerID [20]byte, port uint16) ([]peers.Peer, error) {
	url, err := t.buildTrackerURL(peerID, port)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 15 * time.Second}
	response, err := client.Get(url)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)

	decodedValue, err := bencode.Decode(string(body))
	if err != nil {
		return nil, err
	}

	decodedMap := decodedValue.(map[string]interface{})
	responsePeers := decodedMap["peers"].(string)

	return peers.Unmarshal([]byte(responsePeers))
}
