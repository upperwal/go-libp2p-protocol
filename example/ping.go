package main

import(
	"github.com/upperwal/go-libp2p-protocol"
	"github.com/libp2p/go-libp2p"
	"context"
	"github.com/libp2p/go-libp2p-net"
	"fmt"
	"flag"
)

const(
	PING_REQUEST = "PING_REQUEST"
	PING_RESPONSE = "PING_RESPONSE"
)

var protoPing = protocol.Protocol{
	Name: "ping",
	Version: 1,
	Run: func(stream net.Stream, init bool) {

		fmt.Println("Stream opened to ", stream.Conn().RemotePeer())

		var data [100]byte

		if init {
			_, err := stream.Write( []byte(PING_REQUEST) )
			if err != nil {
				panic(err)
			}


			n, err := stream.Read(data[:])
			if err != nil {
				panic(err)
			}

			if string(data[:n]) == PING_RESPONSE {
				fmt.Printf("%s Alive\n", stream.Conn().RemotePeer())
			}
		} else {

			n, err := stream.Read(data[:])
			if err != nil {
				panic(err)
			}

			if string(data[:n]) == PING_REQUEST {
				fmt.Println("Replying to the ping")
				_, err := stream.Write( []byte(PING_RESPONSE) )
				if err != nil {
					panic(err)
				}
			}
		}
	},
}

func main() {
	dest := flag.String("d", "", "Dest MultiAddr String")
	flag.Parse()

	host, err := libp2p.New(context.Background(), libp2p.Defaults)
	if(err != nil) {
		panic(err)
	}

	xHost := protocol.SetProtocol(host, []protocol.Protocol{protoPing})
	fmt.Printf(xHost.GetCompleteAddr())

	if *dest != "" {
		err := xHost.AddPeerWithAddr(*dest)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Waiting...")
	select{}
}
