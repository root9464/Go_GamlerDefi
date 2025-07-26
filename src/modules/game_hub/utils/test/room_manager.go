package conference_utils

import (
	"sync"
	"time"

	"github.com/pion/webrtc/v4"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/utils/conference"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type IConferenceUsecase interface {
	AddPeer(pc *webrtc.PeerConnection, conn *conference_utils.ThreadSafeWriter)
	SignalPeers() error

	AddTrack(peer *conference_utils.PeerConnection, t *webrtc.TrackRemote) (*webrtc.TrackLocalStaticRTP, error)
	RemoveTrack(track *webrtc.TrackLocalStaticRTP)
}

type RoomManager struct {
	rooms   map[string]*Room // roomID -> Room
	usecase IConferenceUsecase
	logger  *logger.Logger
	mu      sync.RWMutex
}

func NewRoomManager(usecase IConferenceUsecase, logger *logger.Logger) *RoomManager {
	return &RoomManager{
		rooms:   make(map[string]*Room),
		usecase: usecase,
		logger:  logger,
	}
}

type Room struct {
	ID        string
	Peers     map[string]*conference_utils.PeerConnection // userID -> PeerConnection
	CreatedAt time.Time
	mu        sync.RWMutex
}

func (rm *RoomManager) GetOrCreateRoom(roomID string) *Room {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if room, exists := rm.rooms[roomID]; exists {
		return room
	}

	room := &Room{
		ID:        roomID,
		Peers:     make(map[string]*conference_utils.PeerConnection),
		CreatedAt: time.Now(),
	}
	rm.rooms[roomID] = room
	return room
}

func (rm *RoomManager) RemoveRoom(roomID string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if room, exists := rm.rooms[roomID]; exists {
		for _, peer := range room.Peers {
			peer.PC.Close()
		}
		delete(rm.rooms, roomID)
	}
}

func (rm *RoomManager) AddPeerToRoom(roomID, userID string, pc *webrtc.PeerConnection, conn *conference_utils.ThreadSafeWriter) {
	room := rm.GetOrCreateRoom(roomID)

	room.mu.Lock()
	defer room.mu.Unlock()

	peer := &conference_utils.PeerConnection{
		PC:     pc,
		Conn:   conn,
		UserID: userID,
	}
	room.Peers[userID] = peer

	// Используем оригинальный usecase
	rm.usecase.AddPeer(pc, conn)
}

func (rm *RoomManager) RemovePeerFromRoom(roomID, userID string) {
	room := rm.GetOrCreateRoom(roomID)

	room.mu.Lock()
	defer room.mu.Unlock()

	if peer, exists := room.Peers[userID]; exists {
		peer.PC.Close()
		delete(room.Peers, userID)
	}

	if len(room.Peers) == 0 {
		rm.RemoveRoom(roomID)
	}
}

func (rm *RoomManager) SignalPeersInRoom(roomID string) error {
	room := rm.GetOrCreateRoom(roomID)
	room.mu.RLock()
	defer room.mu.RUnlock()

	// Используем оригинальный usecase
	return rm.usecase.SignalPeers()
}
