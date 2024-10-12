package torrentfile

import (
	"crypto/sha1"
	"fmt"
	"os"

	"github.com/mworwa/bittorrent/bencode"
)

type TorrentFile struct {
	Announce     string
	CreationDate int
	Info         Info
	InfoHash     [20]byte
	PieceHashes  [][20]byte
}

type Info struct {
	Length      int
	Name        string
	PieceLength int
	Pieces      string
}

func Open(path string) (TorrentFile, error) {
	fmt.Printf("Opening file from path: %s\n", path)
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return TorrentFile{}, err
	}

	decodedValue, err := bencode.Decode(string(fileContent))

	if err != nil {
		return TorrentFile{}, err
	}
	torrentFile, err := createTorrentFile(decodedValue)
	if err != nil {
		return TorrentFile{}, err
	}

	return torrentFile, nil
}

func createTorrentFile(decodedFile interface{}) (TorrentFile, error) {
	decodedMap := decodedFile.(map[string]interface{})
	announce := decodedMap["announce"].(string)
	creationDate := decodedMap["creation date"].(int)

	info := decodedMap["info"].(map[string]interface{})
	lenght := info["length"].(int)
	name := info["name"].(string)
	pieceLenght := info["piece length"].(int)
	pieces := info["pieces"].(string)
	hash, err := hashInfo(info)
	if err != nil {
		return TorrentFile{}, err
	}
	torrentFile := TorrentFile{
		Announce:     announce,
		CreationDate: creationDate,
		Info: Info{
			Length:      lenght,
			Name:        name,
			PieceLength: pieceLenght,
			Pieces:      pieces,
		},
		InfoHash: hash,
	}
	return torrentFile, nil
}

func hashInfo(info interface{}) ([20]byte, error) {
	encodedInfo, err := bencode.Encode(info)
	if err != nil {
		return [20]byte{}, err
	}
	hash := sha1.Sum(encodedInfo)

	return hash, nil
}

func splitPieceHashes(pieces string) ([][20]byte, error) {
	hashLen := 20
	buf := []byte(pieces)
	if len(buf)%hashLen != 0 {
		err := fmt.Errorf("Recived malformed pieces of lenght %d", len(buf))
		return nil, err
	}

	return nil, nil
}
