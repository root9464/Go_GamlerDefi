package conference_usecase

import (
	"sync"
	"time"

	"github.com/pion/webrtc/v4"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/utils/conference"
)

func (u *ConferenceUsecase) JoinHub(pc *conference_utils.PeerConnection) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	// Находим или создаем хаб
	hub, exists := u.hubs[pc.HubID]
	if !exists {
		hub = &Hub{
			ID:          pc.HubID,
			peers:       make([]*conference_utils.PeerConnection, 0),
			trackLocals: make(map[string]*webrtc.TrackLocalStaticRTP),
			mu:          sync.RWMutex{},
		}
		u.hubs[pc.HubID] = hub
	}

	// Добавляем пира в хаб
	hub.mu.Lock()
	defer hub.mu.Unlock()

	// Проверяем, не подключен ли уже этот пир
	for _, peer := range hub.peers {
		if peer.UserID == pc.UserID {
			return nil // Уже подключен
		}
	}

	hub.peers = append(hub.peers, pc)
	u.logger.Infof("Peer %s joined hub %s", pc.UserID, pc.HubID)

	u.startHubTicker(pc.HubID)

	return nil
}

func (u *ConferenceUsecase) LeaveHub(pc *conference_utils.PeerConnection) error {
	u.mu.RLock()
	hub, exists := u.hubs[pc.HubID]
	u.mu.RUnlock()

	if !exists {
		return nil // Хаб не существует
	}

	hub.mu.Lock()
	defer hub.mu.Unlock()

	// Удаляем пира из списка
	for i, peer := range hub.peers {
		if peer.UserID == pc.UserID {
			// Закрываем соединение
			if peer.PC != nil {
				peer.PC.Close()
			}

			// Удаляем из списка
			hub.peers = append(hub.peers[:i], hub.peers[i+1:]...)
			u.logger.Infof("Peer %s left hub %s", pc.UserID, pc.HubID)

			// Если хаб пуст - удаляем его
			if len(hub.peers) == 0 {
				u.mu.Lock()
				delete(u.hubs, pc.HubID)
				u.mu.Unlock()
				u.logger.Infof("Hub %s removed (no peers left)", pc.HubID)
			}

			return nil
		}
	}

	u.stopHubTicker(pc.HubID)
	return nil
}

func (u *ConferenceUsecase) startHubTicker(hubID string) {
	if _, exists := u.hubTickers[hubID]; exists {
		return
	}

	ticker := time.NewTicker(3 * time.Second)
	u.hubTickers[hubID] = ticker
	u.logger.Infof("start ticke for hub %s", hubID)

	go func(hubID string) {
		for range ticker.C {
			u.DispatchKeyFrames(hubID)
		}
	}(hubID)
}

func (u *ConferenceUsecase) stopHubTicker(hubID string) {
	if ticker, exists := u.hubTickers[hubID]; exists {
		ticker.Stop()
		delete(u.hubTickers, hubID)
	}
}
