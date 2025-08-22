package conference_utils

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/contrib/socketio"
)

func GenerateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func WriteJSON(kws *socketio.Websocket, lock *sync.Mutex, v any) error {
	lock.Lock()
	defer lock.Unlock()
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	kws.Emit(data, socketio.TextMessage)
	return nil
}
