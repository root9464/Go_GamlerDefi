package conference_utils

func GetPeerID(pc *PeerConnection) string {
	stats := pc.PC.GetStats()
	var peerID string
	if connectionStats, ok := stats.GetConnectionStats(pc.PC); ok {
		peerID = connectionStats.ID
	}
	return peerID
}
