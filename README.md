# go-libp2p-protocol
Building Protocols on top of LibP2P

## Example

```go
// Create a new libp2p host
host, err := libp2p.New(context.Background(), libp2p.Defaults)
if(err != nil) {
	panic(err)
}

proto = protocol.Protocol{
	Name: "name_your_protocol",
	Version: 1,
	Run: func(stream net.Stream, init bool) {
		
		// use init to differentiate between dialler and listener peer.
        // init == true (peer who called AddPeer)

		fmt.Println("Stream opened to ", stream.Conn().RemotePeer())

		},
}

// This function will return you an extended host.
// It contains the host (xHost.Host) and a slice of protocols (xHost.Protocol)
xHost := protocol.SetProtocol(host, []protocol.Protocol{proto})

// Now you can add a peer.
// It will add the peer to the peerstore and open a stream.
// Protocol.Run function will will be invoked for both the peers.
xHost.AddPeerWithAddr("/ip4/127.0.0.1/tcp/51912/ipfs/QmbHvB9A2vLMno9LWkWnwnRHa5Fcjnogsq8nzgiooKfJTM")
// or
xHost.AddPeer(pi peerstore.PeerInfo)
```

And there you go. 

Look into the [examples](example) for more details.
