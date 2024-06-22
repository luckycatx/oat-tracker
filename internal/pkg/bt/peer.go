package bt

import (
	"net/netip"
)

type Peer struct {
	ID string
	netip.AddrPort
}

func New(id string, ip string, port uint16) *Peer {
	return &Peer{
		ID:       id,
		AddrPort: netip.AddrPortFrom(netip.MustParseAddr(ip), port),
	}
}

// Pack a peer to binary model (IP + Port)
func (p *Peer) PackToBin() []byte {
	// Since AddrPort.MarshalBinary() returns the port in little-endian format
	// we construct the binary model manually
	ip, port := p.AddrPort.Addr().AsSlice(), p.AddrPort.Port()
	return append(ip, byte(port>>8), byte(port))
}

// Pack a peer to dictionary model
func (p *Peer) PackToDict() map[string]any {
	return map[string]any{
		"peer id": p.ID,
		"ip":      p.AddrPort.Addr().String(),
		"port":    p.AddrPort.Port(),
	}
}
