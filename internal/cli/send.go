package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/tylermeekel/p2pfileshare/internal/network"
)

func send(stunIP, fp string) {
	ip := network.GetIP(stunIP)
	if ip == nil {
		log.Fatalln("ERROR: unable to get IP from STUN server")
	}
	fmt.Println(fp)
	file, err := os.Open(fp)
	if err != nil {
		log.Fatalln("Error opening file:", err)
	}
	defer file.Close()

	stats, err := file.Stat()
	if err != nil {
		fmt.Println("Error reading file stats", err)
	}

	fileName := stats.Name()
	modTime := stats.ModTime()
	fileSize := stats.Size()

	fmt.Printf("Name: %s | Modtime: %v\n", fileName, modTime)

	buf := make([]byte, fileSize)

	_, err = file.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(buf))
	fmt.Println("IP is", ip)
}
