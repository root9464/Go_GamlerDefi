package game_session_ws

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	game_session_contracts "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/contracts"
	game_session_entity "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/entity"
)

type MsgObj struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}

func (h *GameSessionHandler) sockerErr(ep *socketio.EventPayload, err error) {
	errByte := []byte(err.Error())
	ep.Kws.Emit(errByte)
}

func (h *GameSessionHandler) SetupSocketEventHandlers() {
	socketio.On("connect", func(ep *socketio.EventPayload) {
		h.conferenceHandler.Connect(ep)
		log.Infof("conect success: %v", ep.Kws.UUID)
	})

	socketio.On("disconnect", func(ep *socketio.EventPayload) {
		h.conferenceHandler.Disconect(ep)
		h.mapMu.Lock()
		defer h.mapMu.Unlock()

		session, ok := h.uuidToSession[ep.Kws.UUID]
		if !ok {
			return
		}

		userID := ep.Kws.GetAttribute("user_id").(string)
		session.RemovePlayer(userID)

		delete(h.uuidToSession, ep.Kws.UUID)

		h.logger.Infof("User %s (%s) disconnected from session  %s", userID, ep.Kws.UUID, session.ID)
	})

	socketio.On("message", func(ep *socketio.EventPayload) {
		message := new(MsgObj)
		if err := json.Unmarshal(ep.Data, message); err != nil {
			h.sockerErr(ep, err)
		}

		dataByte, err := json.Marshal(message.Data)
		if err != nil {
			h.sockerErr(ep, err)
		}

		if h.config.ConferenceEnabled {
			h.conferenceHandler.SetupSocketEventHandlers(ep)
		}

		if message.Event != "" {
			ep.Kws.Fire(message.Event, dataByte)
		}
	})

	socketio.On("game_action", func(ep *socketio.EventPayload) {
		h.mapMu.RLock()
		session, ok := h.uuidToSession[ep.Kws.UUID]
		h.mapMu.RUnlock()

		if !ok {
			h.sockerErr(ep, errors.New("session not found"))
			return
		}

		userID := ep.Kws.GetAttribute("user_id").(string)
		h.logger.Infof("user_id: %s", userID)

		gameAction := new(game_session_contracts.Action)
		if err := json.Unmarshal(ep.Data, gameAction); err != nil {
			h.sockerErr(ep, err)
			return
		}

		session.HandleMessage(userID, *gameAction)
	})
}

func (h *GameSessionHandler) GameSessionWSHandler(c *fiber.Ctx) error {
	return socketio.New(func(kws *socketio.Websocket) {
		sessionID := kws.Params("session_id")
		userID := kws.Params("user_id")
		gameName := kws.Params("game_name")
		mainColor := kws.Query("main_color")
		highlightColor := kws.Query("highlight_color")

		ctx := context.Background()

		session, err := h.hubManager.ActiveteSession(ctx, sessionID, userID, gameName)
		if err != nil {
			h.logger.Errorf("error creating session: %v", err)
			kws.Close()
			return
		}

		kws.SetAttribute("user_id", userID)
		kws.SetAttribute("session_id", sessionID)

		h.mapMu.Lock()
		h.uuidToSession[kws.UUID] = session
		h.mapMu.Unlock()

		isHost := session.HostID == userID
		conn := &game_session_entity.Connection{
			Kws:    kws,
			ISHost: isHost,
			UserID: userID,
		}

		session.AddPlayer(conn, isHost, mainColor, highlightColor)

		if h.config.ConferenceEnabled {
			h.conferenceHandler.ConferenceHandler(kws)
		}

		welcomeMsg := map[string]string{"message": "Добро пожаловать в игру!"}
		welcomeBytes, _ := json.Marshal(welcomeMsg)
		kws.Emit(welcomeBytes)
	})(c)
}
