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

		h.logger.Infof("starting websocket connection handling: session_id=%s, user_id=%s, game_name=%s", sessionID, userID, gameName)

		ctx := context.Background()

		h.logger.Infof("activating session: session_id=%s, user_id=%s, game_name=%s", sessionID, userID, gameName)
		session, err := h.hubManager.ActiveteSession(ctx, sessionID, userID, gameName)
		if err != nil {
			h.logger.Errorf("session activation failed: session_id=%s, error=%v", sessionID, err)
			kws.Close()
			return
		}
		h.logger.Infof("session activated successfully: session_id=%s", sessionID)

		kws.SetAttribute("user_id", userID)
		kws.SetAttribute("session_id", sessionID)
		h.logger.Infof("connection attributes set: user_id=%s, session_id=%s", userID, sessionID)

		h.mapMu.Lock()
		h.uuidToSession[kws.UUID] = session
		h.mapMu.Unlock()
		h.logger.Infof("uuid to session mapping added: uuid=%s, session_id=%s", kws.UUID, sessionID)

		isHost := session.HostID == userID
		h.logger.Infof("connection type determined: user_id=%s, is_host=%t", userID, isHost)

		conn := &game_session_entity.Connection{
			Kws:    kws,
			ISHost: isHost,
			UserID: userID,
		}

		h.logger.Infof("adding player to session: session_id=%s, user_id=%s, is_host=%t", sessionID, userID, isHost)
		session.AddPlayer(conn, isHost, mainColor, highlightColor)
		h.logger.Infof("player added to session successfully: session_id=%s, user_id=%s", sessionID, userID)

		if h.config.ConferenceEnabled {
			h.logger.Infof("attaching conference handler for session: session_id=%s", sessionID)
			h.conferenceHandler.ConferenceHandler(kws)
		}

		welcomeMsg := map[string]string{"message": "Добро пожаловать в игру!"}
		welcomeBytes, _ := json.Marshal(welcomeMsg)
		kws.Emit(welcomeBytes)
		h.logger.Infof("welcome message sent to user: user_id=%s", userID)

		h.logger.Infof("websocket connection established successfully: session_id=%s, user_id=%s", sessionID, userID)
	})(c)
}
