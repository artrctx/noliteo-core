package walktalk

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func (wk *WalkTalkService) WSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("failed to initialize websocket", slog.Any("error", err))
		// might not accept http res
		http.Error(w, fmt.Sprintf("failed to initialize walkie talkie websocket: %v", err.Error()), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			slog.Error("failed to receive msg", slog.Any("error", err))
			http.Error(w, fmt.Sprintf("failed to receive msg: %v", err.Error()), http.StatusInternalServerError)
			break
		}
		log.Printf("websocket msg: %v, msgType: %v", msg, msgType)
	}
}
