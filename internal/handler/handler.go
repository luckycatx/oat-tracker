package handler

import (
	"net/http"

	"github.com/anacrolix/torrent/bencode"
	"github.com/gin-gonic/gin"

	"github.com/luckycatx/oat-tracker/internal/pkg/bt"
	"github.com/luckycatx/oat-tracker/internal/pkg/conf"
	"github.com/luckycatx/oat-tracker/internal/pkg/info"
	"github.com/luckycatx/oat-tracker/internal/repo"
)

type Handler struct {
	repo repo.Repoer
	cfg  *conf.Config
}

func New(cfg *conf.Config) *Handler {
	r := repo.NewMemRepo(cfg.NumShard)
	// Use for count peers
	info.Counter = r
	return &Handler{
		cfg:  cfg,
		repo: r,
	}
}

func (h *Handler) Announce(c *gin.Context) {
	var req = &bt.AnnounceReq{
		Compact:   true,
		NumWant:   30,
		IP:        c.ClientIP(),
		TrackerID: h.cfg.TrackerID,
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var reqPeer = bt.New(req.PeerID, req.IP, req.Port)
	var room = c.Param("room")
	if room == "" {
		c.JSON(400, gin.H{"error": "room is required"})
	}
	// Record the room info
	if _, ok := info.Infos[room]; !ok {
		info.Infos[room] = make(info.Info)
	}
	if _, ok := info.Infos[room][req.InfoHash]; !ok {
		info.Infos[room][req.InfoHash] = 0
	}

	switch req.Event {
	// case "started": behaves the same as default
	case "stopped":
		h.repo.DeletePeer(room, req.InfoHash, reqPeer, *req.Left == 0)
	case "completed":
		h.repo.GraduateLeecher(room, req.InfoHash, reqPeer)
	default:
		h.repo.PutPeer(room, req.InfoHash, reqPeer, *req.Left == 0)
	}

	interval, minInterval := h.cfg.Interval, h.cfg.MinInterval

	numSeeders, _, numLeechers := h.repo.CountPeers(room, req.InfoHash)
	if numSeeders == 0 {
		interval /= 2
		minInterval /= 2
	}

	var peers = h.repo.GetPeers(room, req.InfoHash, reqPeer, *req.Left == 0, req.NumWant)
	var packedPeer any
	var packedPeer6 []byte
	if req.Compact {
		packedPeer = make([]byte, 0, len(peers))
		for _, peer := range peers {
			if peer.AddrPort.Addr().Is4() {
				packedPeer = append(packedPeer.([]byte), peer.PackToBin()...)
			} else {
				packedPeer6 = append(packedPeer6, peer.PackToBin()...)
			}
		}
	} else {
		packedPeer = make([]map[string]any, 0, len(peers))
		for _, peer := range peers {
			packedPeer = append(packedPeer.([]map[string]any), peer.PackToDict())
		}
	}

	var resp = &bt.AnnounceResp{
		Interval:    interval,
		MinInterval: minInterval,
		TrackerID:   req.TrackerID,
		Complete:    numSeeders,
		Incomplete:  numLeechers,
		Peers:       packedPeer,
		Peers6:      packedPeer6,
	}
	respData, err := bencode.Marshal(resp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.String(http.StatusOK, string(respData))
}

func (h *Handler) Scrape(c *gin.Context) {
	var req *bt.ScrapeReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var resp = &bt.ScrapeResp{Files: make(map[string]bt.Stats)}
	for _, ih := range req.InfoHashes {
		s, sn, l := h.repo.CountPeers("", ih)
		var stats = &bt.Stats{
			Complete:   s,
			Downloaded: sn,
			Incomplete: l,
		}
		resp.Files[ih] = *stats
	}
	respData, err := bencode.Marshal(resp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.String(http.StatusOK, string(respData))
}

func (h *Handler) Cleanup() {
	h.repo.Cleanup()
}
