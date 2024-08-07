package bt

type AnnounceReq struct {
	InfoHash   string `form:"info_hash" binding:"required"`
	PeerID     string `form:"peer_id"`
	Port       uint16 `form:"port" binding:"required,number"`
	Uploaded   uint   `form:"uploaded"`
	Downloaded uint   `form:"downloaded"`                     // in bytes
	Left       *uint  `form:"left" binding:"required,number"` // in bytes | type *uint in order to accept zero value
	Compact    bool   `form:"compact"`
	NoPeerID   bool   `form:"no_peer_id"`
	Event      string `form:"event"`
	IP         string `form:"ip"`      // dotted quad format or rfc3513 hexed IPv6 format
	NumWant    uint   `form:"numwant"` // Number of peers wanted
	Key        string `form:"key"`
	TrackerID  string `form:"trackerid"`
}

// Response peers in binary model (compact response)
type AnnounceResp struct {
	Interval    uint   `bencode:"interval"`
	MinInterval uint   `bencode:"min_interval,omitempty"`
	TrackerID   string `bencode:"tracker_id"`
	Complete    uint   `bencode:"complete"`
	Incomplete  uint   `bencode:"incomplete"`
	Peers       any    `bencode:"peers"`            // Dictionary model or Binary model (both IPv4 and IPv6)
	Peers6      []byte `bencode:"peers6,omitempty"` // Binary model for IPv6
}

type Stats struct {
	Complete   uint `bencode:"complete"`   // Seeders count
	Downloaded uint `bencode:"downloaded"` // Snatchers count
	Incomplete uint `bencode:"incomplete"` // Leechers count
}

type ScrapeReq struct {
	InfoHashes []string `query:"info_hash"`
}

type ScrapeResp struct {
	Files map[string]Stats `bencode:"files"` // Map of infohash to Stat
}
