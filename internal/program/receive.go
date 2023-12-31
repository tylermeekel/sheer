package program

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"sync"

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
		fmt.Println("Error reading offer:", err)
	}
	fmt.Println(getLineWithMessage(""))

	jsonOffer := make([]byte, base64.StdEncoding.DecodedLen(len(encodedOffer)))
	n, err := base64.StdEncoding.Decode(jsonOffer, []byte(encodedOffer))
	if err != nil {
		fmt.Println("Error decoding offer:", err)
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

	var wg sync.WaitGroup
	mux := &sync.Mutex{}
	var data []byte

	fmt.Println("numchunks:", offer.NumberOfChunks)
	wg.Add(offer.NumberOfChunks)

	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {

		d.OnClose(func() {
			os.Exit(0)
		})

		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			// Use string messages to signal when sending is complete.
			if msg.IsString {
				if string(msg.Data) == "sent" {
					//Wait for data to be added to slice, close datachannel, and then write file.
					wg.Wait()
					d.Close()
					err := os.WriteFile(offer.Metadata.FileName, data, 0644)
					if err != nil {
						fmt.Println("error writing file", err)
					}
					fmt.Println("File received!")
					os.Exit(0)
				}
			} else {
				go func() {
					defer wg.Done()
					mux.Lock()
					data = append(data, msg.Data...)
					mux.Unlock()
				}()
			}
		})
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
