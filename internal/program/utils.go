package program

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/pion/webrtc/v3"
	"golang.org/x/term"
)

type Config struct {
	Servers []webrtc.ICEServer `json:"servers"`
}

type Metadata struct {
	FileName string `json:"filename"`
	FileSize int64  `json:"filesize"`
}

type FileData struct {
	Metadata Metadata `json:"metadata"`
	Data     []byte   `json:"data"`
}

func getLineWithMessage(message string) string {
	messageLength := len(message)
	width, _, err := term.GetSize(0)
	if err != nil {
		width = 100
	}

	if messageLength > width {
		return message
	}

	lineSize := (width / 2) - (messageLength / 2)
	line := strings.Repeat("-", lineSize) + message + strings.Repeat("-", lineSize)

	if len(line) < width {
		line += "-"
	} else if len(line) > width {
		line = line[:len(line)-1]
	}

	return line
}

func ClearTerminal() {
	fmt.Printf("\x1bc")
}

func ReadConfig() *Config {
	ex, err := os.Executable()
	if err != nil{
		fmt.Println("Error getting executable location")
	}

	exPath := filepath.Dir(ex)
	fmt.Println(exPath)
	
	configPath := exPath + "/sheerconfig.json"

	data, err := os.ReadFile(configPath)
	if errors.Is(err, fs.ErrNotExist) {
		return CreateDefaultConfig(configPath)
	}

	var config Config
	json.Unmarshal(data, &config)
	return &config
}

func CreateDefaultConfig(configPath string) *Config {
	fmt.Println("Config not found. Creating new default config.")

	ex, err := os.Executable()
	exPath := filepath.Dir(ex)
	fmt.Println(exPath)

	file, err := os.Create(configPath)
	if err != nil {
		panic(err)
	}

	defaultConfig := Config{
		Servers: []webrtc.ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
		},
	}

	jsonConfig, err := json.MarshalIndent(&defaultConfig, "", "\t")

	file.Write(jsonConfig)
	file.Close()
	return &defaultConfig
}

func ConfigureDefaultPeerConnection(peerConnection *webrtc.PeerConnection) {
	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {

		if s == webrtc.PeerConnectionStateConnected {
			fmt.Println("Connection established.")
		}

		if s == webrtc.PeerConnectionStateFailed {
			fmt.Println("Peer connection failed. Exiting now.")
			os.Exit(0)
		}

		if s == webrtc.PeerConnectionStateClosed {
			fmt.Println("Peer connection closed. Exiting now.")
			os.Exit(0)
		}
	})
}
