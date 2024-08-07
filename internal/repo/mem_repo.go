package repo

import (
	"sync"
	"time"

	"github.com/luckycatx/oat-tracker/internal/pkg/bt"
	"github.com/luckycatx/oat-tracker/internal/pkg/conf"

	"github.com/minio/highwayhash"
)

var peer_lifetime = conf.Load().PeerLifetime

type MemRepo struct {
	shards []*shard
}

// Hash key is the hash of info_hash + room
// Use RWMap instead of sync.Map because
// there's a certain level of writing and deleting
type shard struct {
	swarms        map[string]swarm
	num_snatchers map[string]uint
	sync.RWMutex
}

// Use map to store peers and their last seen time
// int64 for storing Unix timestamp
type swarm struct {
	seeders  map[bt.Peer]int64
	leechers map[bt.Peer]int64
}

func NewMemRepo(shard_size int) *MemRepo {
	var mem_repo = &MemRepo{shards: make([]*shard, shard_size)}
	for i := range shard_size {
		mem_repo.shards[i] = &shard{
			swarms:        make(map[string]swarm),
			num_snatchers: make(map[string]uint),
		}
	}
	return mem_repo
}

func (mr *MemRepo) ShardIndex(room, info_hash string) int {
	var idx = int(highwayhash.Sum64([]byte(room+info_hash), make([]byte, 32)) % uint64(len(mr.shards)))
	return idx
}

func (mr *MemRepo) PutPeer(room, info_hash string, peer *bt.Peer, seed bool) {
	var idx = mr.ShardIndex(room, info_hash)
	mr.shards[idx].Lock()
	defer mr.shards[idx].Unlock()

	if _, ok := mr.shards[idx].swarms[info_hash]; !ok {
		mr.shards[idx].swarms[info_hash] = swarm{seeders: make(map[bt.Peer]int64), leechers: make(map[bt.Peer]int64)}
	}
	var sw = mr.shards[idx].swarms[info_hash]
	if seed {
		sw.seeders[*peer] = time.Now().Unix()
	} else {
		sw.leechers[*peer] = time.Now().Unix()
	}
}

func (mr *MemRepo) GetPeers(room, info_hash string, peer *bt.Peer, seed bool, num_want uint) []*bt.Peer {
	var idx = mr.ShardIndex(room, info_hash)
	mr.shards[idx].RLock()
	defer mr.shards[idx].RUnlock()

	var sw = mr.shards[idx].swarms[info_hash]
	var peers = make([]*bt.Peer, 0, num_want)
	// Seeder is not interested in other seeders
	if !seed {
		for p := range sw.seeders {
			peers = append(peers, &p)
			if uint(len(peers)) >= num_want {
				break
			}
		}
	}
	for p := range sw.leechers {
		// Skip the peer which is leecher itself
		if p == *peer {
			continue
		}
		peers = append(peers, &p)
		if uint(len(peers)) >= num_want {
			break
		}
	}

	return peers
}

func (mr *MemRepo) DeletePeer(room, info_hash string, peer *bt.Peer, seed bool) {
	var idx = mr.ShardIndex(room, info_hash)
	mr.shards[idx].Lock()
	defer mr.shards[idx].Unlock()

	sw, ok := mr.shards[idx].swarms[info_hash]
	if !ok {
		return
	}
	if seed {
		delete(sw.seeders, *peer)
	} else {
		delete(sw.leechers, *peer)
	}
}

func (mr *MemRepo) GraduateLeecher(room, info_hash string, peer *bt.Peer) {
	var idx = mr.ShardIndex(room, info_hash)
	mr.shards[idx].Lock()
	defer mr.shards[idx].Unlock()

	if _, ok := mr.shards[idx].swarms[info_hash]; !ok {
		mr.shards[idx].swarms[info_hash] = swarm{seeders: make(map[bt.Peer]int64), leechers: make(map[bt.Peer]int64)}
	}
	var sw = mr.shards[idx].swarms[info_hash]
	delete(sw.leechers, *peer)
	sw.seeders[*peer] = time.Now().Unix()
	mr.shards[idx].num_snatchers[info_hash]++
}

func (mr *MemRepo) CountPeers(room, info_hash string) (num_seeders, num_snachers, num_leechers uint) {
	var idx = mr.ShardIndex(room, info_hash)
	mr.shards[idx].RLock()
	defer mr.shards[idx].RUnlock()

	sw, ok := mr.shards[idx].swarms[info_hash]
	if !ok {
		return
	}
	num_seeders, num_snachers, num_leechers =
		uint(len(sw.seeders)),
		(mr.shards[idx].num_snatchers[info_hash]),
		uint(len(sw.leechers))
	return
}

func (mr *MemRepo) Cleanup() {
	for {
		time.Sleep(5 * time.Minute)
		for _, shard := range mr.shards {
			shard.Lock()
			for info_hash, sw := range shard.swarms {
				for p, last_seen := range sw.seeders {
					if time.Now().Unix()-last_seen > peer_lifetime {
						delete(sw.seeders, p)
					}
				}
				for p, last_seen := range sw.leechers {
					if time.Now().Unix()-last_seen > peer_lifetime {
						delete(sw.leechers, p)
					}
				}
				if len(sw.seeders) == 0 && len(sw.leechers) == 0 {
					delete(shard.swarms, info_hash)
				}
			}
			shard.Unlock()
		}
	}
}
