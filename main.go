package main

import (
	"fmt"

	"github.com/mworwa/bittorrent/bencode"
)

func main() {
	// torrentfile, err := torrentfile.Open("file.torrent")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	encoded, err := bencode.Encode(map[string]interface{}(map[string]interface{}{"dict": map[string]interface{}{"123": "test", "456": "thing"}, "list": []interface{}{"list-item-1", "list-item-2"}, "number": 123456, "string": "value"}))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(encoded))
	// fmt.Println(torrentfile)
}
