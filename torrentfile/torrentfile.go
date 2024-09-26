package torrentfile

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mworwa/bittorrent/bencode"
)

type TorrentFile struct {
	Announce     string `json:"announce"`
	Comment      string `json:"comment"`
	CreationDate int    `json:"creation date"`
	Info         Info   `json:"info"`
	InfoHash     [20]byte
}

type Info struct {
	Length      int    `json:"length"`
	Name        string `json:"name"`
	PieceLength int    `json:"piece length"`
	Pieces      string `json:"pieces"`
}

func Open(path string) (TorrentFile, error) {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return TorrentFile{}, err
	}

	decodedValue, err := bencode.Decode(string(fileContent))
	if err != nil {
		return TorrentFile{}, err
	}

	jsonData, err := json.Marshal(decodedValue)
	if err != nil {
		return TorrentFile{}, fmt.Errorf("error marshalling json: %w", err)
	}
	var torrentFile TorrentFile
	err = json.Unmarshal(jsonData, &torrentFile)

	if err != nil {
		return TorrentFile{}, fmt.Errorf("error unmarshalling json: %w", err)
	}
	hash, err := torrentFile.Info.hashInfo()
	fmt.Println(hash, err)
	if err != nil {
		return TorrentFile{}, err
	}
	torrentFile.InfoHash = hash
	return torrentFile, nil
}

func (i *Info) hashInfo() ([20]byte, error) {
	info := map[string]interface{}{
		"name":         i.Name,
		"length":       i.Length,
		"piece length": i.PieceLength,
		"pieces":       i.Pieces,
	}
	encodedInfo, err := bencode.Encode(info)
	if err != nil {
		return [20]byte{}, err
	}
	hash := sha1.Sum(encodedInfo)

	return hash, nil
}

func (t *TorrentFile) DownloadToFile(path string) error {
	var peerID [20]byte
	_, err := rand.Read(peerID[:])
	if err != nil {
		return err
	}

	return nil
}
