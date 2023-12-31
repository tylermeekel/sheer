package program

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/pion/webrtc/v3"
)

type Offer struct {
	Description    webrtc.SessionDescription `json:"description"`
	Metadata       Metadata                  `json:"metadata"`
	NumberOfChunks int
}

func send(config *Config, filepath string) {
	if filepath == "" {
		fmt.Println("Please enter a filename!")
		os.Exit(1)
	}

	ClearTerminal()

	rtcConfig := webrtc.Configuration{
		ICEServers: config.Servers,
	}

	peerConnection, err := webrtc.NewPeerConnection(rtcConfig)
	if err != nil {
		panic(err)
	}
	defer func() {
		if cErr := peerConnection.Close(); cErr != nil {
			fmt.Println("Cannot close peer connection:", cErr.Error())
		}
	}()

	file, err := os.Open(filepath)
	if errors.Is(err, fs.ErrNotExist) {
		fmt.Printf("File \"%s\" does not exist.\n", filepath)
		os.Exit(1)
	} else if err != nil {
		fmt.Printf("Error opening file: \"%s\"\n", filepath)
		os.Exit(1)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file stats:", err.Error())
		os.Exit(1)
	}

	fileData := make([]byte, fileInfo.Size())
	_, err = file.Read(fileData)
	if err != nil {
		fmt.Println("Error reading file:", err.Error())
		os.Exit(1)
	}
	file.Close()

	fmt.Println("Configuring connection information...")

	chunks := splitBytesBySize(fileData, 65535)

	ConfigureSenderPeerConnection(peerConnection, chunks)
	ClearTerminal()

	//Serialize localDescription into JSON
	offer := Offer{
		Description: *peerConnection.LocalDescription(),
		Metadata: Metadata{
			FileName: fileInfo.Name(),
			FileSize: fileInfo.Size(),
		},
		NumberOfChunks: len(chunks),
	}

	jsonOffer, err := json.Marshal(offer)
	if err != nil {
		panic(err)
	}

	//Encode JSON representation of localDescription into Base64 to copy-paste more easily
	fmt.Println(getLineWithMessage("COPY AND SEND TO PEER"))
	encoder := base64.NewEncoder(base64.StdEncoding, os.Stdout)
	encoder.Write([]byte(jsonOffer))
	encoder.Close()

	fmt.Println("\n" + getLineWithMessage("") + "\n")

	//Accept Base64 encoded answer from other peer
	var encodedAnswer string
	fmt.Println(getLineWithMessage("PASTE RESPONSE AND HIT ENTER"))
	_, err = fmt.Scanln(&encodedAnswer)
	if err != nil {
		panic(err)
	}

	fmt.Println(getLineWithMessage("") + "\n")

	//Parse answer into JSON representation of Data
	jsonAnswer := make([]byte, base64.StdEncoding.DecodedLen(len(encodedAnswer)))
	n, err := base64.StdEncoding.Decode(jsonAnswer, []byte(encodedAnswer))
	if err != nil {
		fmt.Println("Unable to decode response:", err)
		os.Exit(1)
	}
	jsonAnswer = jsonAnswer[:n]

	//Unmarshal JSON-formatted answer into a SessionDescription object
	var answer webrtc.SessionDescription
	err = json.Unmarshal(jsonAnswer, &answer)
	if err != nil {
		panic(err)
	}

	//Set RemoteDescription and begin waiting for connection
	peerConnection.SetRemoteDescription(answer)

	//Block indefinitely
	select {}

}

func ConfigureSenderPeerConnection(peerConnection *webrtc.PeerConnection, chunks [][]byte) {
	ConfigureDefaultPeerConnection(peerConnection)

	//FILE DATA CHANNEL
	fileDataChannel, err := peerConnection.CreateDataChannel("filedata", nil)
	if err != nil {
		panic(err)
	}

	fileDataChannel.OnClose(func() {
		fmt.Println("Done!")
		peerConnection.Close()
		os.Exit(1)
	})

	fileDataChannel.OnOpen(func() {
		fmt.Println("Sending file data...")

		for _, chunk := range chunks {
			fileDataChannel.Send(chunk)
		}

		for fileDataChannel.BufferedAmount() > 0 {
		}
		fmt.Println("Done sending!")
		fileDataChannel.SendText("sent")
	})

	fileDataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		if string(msg.Data) == "done" {
			fmt.Println("Data received... send complete!")
			os.Exit(0)
		}
	})

	//Configure Local SDP
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		panic(err)
	}

	//Wait until ICE Candidates are gathered
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)
	<-gatherComplete
}
