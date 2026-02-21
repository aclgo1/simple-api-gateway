package ex

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/aclgo/simple-api-gateway/internal/auth"
	"github.com/aclgo/simple-api-gateway/internal/delivery/websocket/service"
	"github.com/gorilla/websocket"
)

type exService struct {
	Upgrader websocket.Upgrader
}

func NewExService() *exService {
	return &exService{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

type Client struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

func (c *Client) Send(ctx context.Context, msg any) error {

	m, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	err = c.conn.WriteMessage(websocket.TextMessage, m)
	if err != nil {
		return fmt.Errorf("conn.WriteMessage: %w", err)
	}
	return nil
}

func (s *exService) ExWs(ctxF context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := s.Upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Printf("Error ws.Upgrader.Upgrade: %v\n", err)
			return
		}

		defer conn.Close()

		client := Client{conn: conn}

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		stop := context.AfterFunc(ctxF, func() {
			client.conn.Close()
			cancel()
		})

		defer stop()

		paramsUserLogger, ok := r.Context().Value(auth.KeyCtxParamsToken).(*auth.ParamsToken)
		if !ok {
			return
		}

		r = r.WithContext(ctx)

		for {

			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("Error connection end: %v\n", err)
				return
			}

			var params service.ParamsWsMessage

			err = json.Unmarshal(msg, &params)
			if err != nil {
				if err := client.Send(ctx,
					&service.ParamsWsError{
						Error: fmt.Errorf("json.Unmarshal: %w", err).Error(),
					},
				); err != nil {
					fmt.Printf("Error ws.Send: %v\n", err)
				}

				continue
			}

			if params.Line == "ping" {
				if err := client.Send(ctx, "pong"); err != nil {
					fmt.Printf("Error ws.Send: %v\n", err)
				}

				continue
			}

			if err := client.Send(ctx,
				&service.ParamsOutput{
					Message: params.Line + paramsUserLogger.UserID,
				},
			); err != nil {
				fmt.Printf("Error ws.Send: %v\n", err)
			}
		}
	}
}
