package main

import (
	"fmt"

	"github.com/mworwa/bittorrent/client"
)

func main() {
	client := client.Client{}

	err := client.DownloadTorrent("file.torrent", "test")

	fmt.Println(err)
}
