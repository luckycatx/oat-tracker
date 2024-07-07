package info

import (
	"github.com/luckycatx/oat-tracker/internal/repo"
)

type counter interface {
	CountPeers(room, info_hash string) (uint, uint, uint)
}
type Info map[string]int

var Counter counter = (*repo.MemRepo)(nil)
var Infos = make(map[string]Info)

func Update(room string) {
	info := Infos[room]
	for ih := range info {
		a, b, c := Counter.CountPeers(room, ih)
		info[ih] = int(a + b + c)
	}
}
