package protocol

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p-net"
	"github.com/libp2p/go-libp2p-host"
	"github.com/libp2p/go-libp2p-protocol"
	"github.com/libp2p/go-libp2p-peerstore"
	"github.com/multiformats/go-multiaddr"
	"github.com/libp2p/go-libp2p-peer"
	logging "github.com/ipfs/go-log"
)

var (
	log = logging.Logger("protocol")
)

type StreamHandler func(stream net.Stream, initiator bool)

type Protocol struct {
	Name string
	Version uint
	Run StreamHandler
}

type ExtendedHost struct {
	Host host.Host
	Protocol []Protocol
}

func SetProtocol(host host.Host, proto []Protocol) *ExtendedHost {
	for _, pr := range proto {
		pid := protocol.ID(pr.Name + string(pr.Version))
		host.SetStreamHandler(pid, func(stream net.Stream) {
			pr.Run(stream, false)
		})
	}

	return &ExtendedHost{
		Host: host,
		Protocol: proto,
	}
}

func (eh *ExtendedHost) AddPeerWithAddr(addr string) error {
	pi, err := AddrStringToPeerInfo(addr)
	if err != nil {
		return err
	}
	log.Debugf("Adding Peer: ", pi.ID)
	return eh.AddPeer(pi)
}

func (eh *ExtendedHost) AddPeer(info peerstore.PeerInfo) error {
	eh.Host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	for _, pr := range eh.Protocol {
		pid := protocol.ID(pr.Name + string(pr.Version))
		s, err := eh.Host.NewStream(context.Background(), info.ID, pid)
		if err != nil {
			return err
		}
		go func(stream net.Stream) {
			pr.Run(stream, true)
		}(s)
	}
	return nil
}

func (eh *ExtendedHost) GetCompleteAddr() string {
	return fmt.Sprintf("%s/ipfs/%s\n", eh.Host.Addrs()[0], eh.Host.ID().Pretty())
}

func AddrStringToPeerInfo(addr string) (peerstore.PeerInfo, error) {
	// The following code extracts target's the peer ID from the
	// given multiaddress
	ipfsaddr, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		return peerstore.PeerInfo{}, err
	}
	pid, err := ipfsaddr.ValueForProtocol(multiaddr.P_IPFS)
	if err != nil {
		return peerstore.PeerInfo{}, err
	}

	peerid, err := peer.IDB58Decode(pid)
	if err != nil {
		return peerstore.PeerInfo{}, err
	}

	// Decapsulate the /ipfs/<peerID> part from the target
	// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
	targetPeerAddr, _ := multiaddr.NewMultiaddr(
		fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))
	targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

	// We have a peer ID and a targetAddr so we add
	// it to the peerstore so LibP2P knows how to contact it
	//h.Peerstore().AddAddr(peerid, targetAddr, peerstore.PermanentAddrTTL)
	return peerstore.PeerInfo{
		ID: peerid,
		Addrs: []multiaddr.Multiaddr{targetAddr},
	}, nil
}


