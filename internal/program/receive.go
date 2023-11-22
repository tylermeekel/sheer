package program

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pion/webrtc/v3"
)

func receive(config Config) {

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

	ConfigureDefaultPeerConnection(peerConnection)

	var encodedOffer string
	ClearTerminal()
	fmt.Println(getLineWithMessage("PASTE OFFER AND HIT ENTER"))
	_, err = fmt.Scanln(&encodedOffer)
	if err != nil {
		panic(err)
	}
	fmt.Println(getLineWithMessage(""))

	jsonOffer := make([]byte, base64.StdEncoding.DecodedLen(len(encodedOffer)))
	n, err := base64.StdEncoding.Decode(jsonOffer, []byte(encodedOffer))
	if err != nil {
		panic(err)
	}
	jsonOffer = jsonOffer[:n]

	var offer Offer
	err = json.Unmarshal(jsonOffer, &offer)
	if err != nil {
		panic(err)
	}

	ConfigureReceiverPeerConnection(peerConnection, offer)

	jsonAnswer, err := json.Marshal(*peerConnection.LocalDescription())
	if err != nil {
		panic(err)
	}

	ClearTerminal()
	fmt.Println(getLineWithMessage("COPY AND SEND TO PEER"))
	encoder := base64.NewEncoder(base64.StdEncoding, os.Stdout)
	encoder.Write(jsonAnswer)
	encoder.Close()
	fmt.Println("\n" + getLineWithMessage(""))

	select {}
}

func ConfigureReceiverPeerConnection(peerConnection *webrtc.PeerConnection, offer Offer) {

	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		if d.Label() == "filedata" {
			d.OnOpen(func() {
				fmt.Println("Parsing filedata...")
			})

			d.OnMessage(func(msg webrtc.DataChannelMessage) {
				os.WriteFile(filepath.Base(offer.Metadata.FileName), msg.Data, 0644)
			})

			d.OnClose(func() {
				fmt.Println("Closed data transfer.")
				os.Exit(0)
			})
		}
	})

	err := peerConnection.SetRemoteDescription(offer.Description)
	if err != nil {
		panic(err)
	}

	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	//Wait for ICE candidate gathering to complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)
	<-gatherComplete
}
