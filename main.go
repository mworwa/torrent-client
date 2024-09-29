package main

import (
	"fmt"

	"github.com/mworwa/bittorrent/torrentfile"
)

func main() {
	fmt.Println("Opening file")
	file, err := torrentfile.Open("file.torrent")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("File opened")
	fmt.Println("Downloading file")
	err = file.DownloadToFile("test")
	if err != nil {
		fmt.Println(err)
	}
}
