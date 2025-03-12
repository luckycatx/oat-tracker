package repo

import "github.com/luckycatx/oat-tracker/internal/pkg/bt"

type Repoer interface {
	PutPeer(room, infoHash string, peer *bt.Peer, seed bool)
	GetPeers(room, infoHash string, peer *bt.Peer, seed bool, numWant uint) []*bt.Peer
	DeletePeer(room, infoHash string, peer *bt.Peer, seed bool)
	GraduateLeecher(room, infoHash string, peer *bt.Peer)
	CountPeers(room, infoHash string) (numSeeders, numSnachers, numLeechers uint)
	Cleanup()
}

// Interface check
var _ Repoer = (*MemRepo)(nil)
