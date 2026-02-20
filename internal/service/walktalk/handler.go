package walktalk

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/artrctx/noliteo-core/internal/database/repository"
	"github.com/artrctx/noliteo-core/internal/hub"
	"github.com/artrctx/noliteo-core/internal/jwt"
	"github.com/artrctx/noliteo-core/internal/middleware"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type WSMessage struct {
	// register description
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

type RtcMsg struct {
	Sdp string `json:"sdp"`
	//"answer" | "offer" | "pranswer" | "rollback"
	Type repository.RtcType `json:"type"`
}

func (wk *WalkTalkService) WSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// errors.Is(err, http.ErrHijacked) or direct comparison returned false
		if err.Error() == http.ErrHijacked.Error() {
			log.Println("This request is already hijecked")
			return
		}
		slog.Error("failed to initialize websocket", slog.Any("error", err))
		http.Error(w, fmt.Sprintf("failed to initialize walkie talkie websocket: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	tknCtxVal := r.Context().Value(middleware.TokenCtxKey)
	token, ok := tknCtxVal.(jwt.Token)
	if !ok {
		slog.Error("invalid token context", slog.Any("token", tknCtxVal))
		http.Error(w, "Invalid token context provided", http.StatusUnauthorized)
		return
	}

	client := &hub.Client{ID: token.TID, Hub: wk.Hub, Conn: conn, Send: make(chan []byte)}
	client.Hub.Register <- client

	go client.Read()
	go client.Write()

	w.WriteHeader(http.StatusAccepted)
	// repo := repository.New(wk.DB)
	// defer func() {
	// 	if err := repo.DeleteRTCDescription(context.Background(), token.TID); err != nil {
	// 		slog.Error("Delete user rtc description failed", slog.Any("error", err))
	// 	}
	// }()

	// for {
	// 	_, msgData, err := conn.ReadMessage()
	// 	if err != nil {
	// 		slog.Error("failed to receive msg", slog.Any("error", err))
	// 		break
	// 	}

	// 	var msg RtcMsg
	// 	if err := json.Unmarshal(msgData, &msg); err != nil {
	// 		slog.Error("failed to unmarshall rtc msg", slog.Any("error", err))
	// 		// send error msg
	// 		continue
	// 	}

	// 	switch msg.Type {
	// 	case repository.RtcTypeAnswer, repository.RtcTypeOffer:
	// 		if err := repo.CreateRTCDescription(r.Context(), repository.CreateRTCDescriptionParams{
	// 			TokenID: token.TID,
	// 			Sdp:     msg.Sdp,
	// 			Type:    msg.Type,
	// 		}); err != nil {
	// 			slog.Error("failed to insert rtc description", slog.Any("error", err))
	// 		}
	// 		//send msgs out
	// 	default:
	// 		slog.Error("Invalid msg type", slog.Any("msg", msg))
	// 		// send error msg
	// 	}
	// }
}
