package network

import (
	"log"
	"net"

	"github.com/pion/stun"
)

func GetIP(uri string) (ip net.IP) {
	u, err := stun.ParseURI(uri)
	if err != nil {
		log.Fatalln("err")
	}

	c, err := stun.DialURI(u, &stun.DialConfig{})
	if err != nil {
		log.Fatalln(err)
	}

	message := stun.MustBuild(stun.TransactionID, stun.BindingRequest)
	if err := c.Do(message, func(res stun.Event) {
		if res.Error != nil {
			log.Fatalln(res.Error)
		}

		var xorAddr stun.XORMappedAddress
		if err := xorAddr.GetFrom(res.Message); err != nil {
			log.Fatalln(err)
		}

		ip = xorAddr.IP
	}); err != nil {
		log.Fatalln(err)
	}

	return
}
