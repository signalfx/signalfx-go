package signalflow

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

const sampleDataMsgBase64 = "AQUBAGNoYW5uZWwtMQAAAAAAAAAfiwgAAAAAAAAAY2BgjDM0fuPAwMAgwAQkGLa3Hbjh0DrjDBCcBQtcPh/M4LDtCIjJABaYfOodp8PhzM3GQAAWqGwI9XYIB/MhAq90Xk10WC6QBgJggdXf6vQdWgsQZtx+ep7doXU3QovGzqWLHdoXIbTw8FivdNi2AiEQrvdtt8PSmwiHndJlvOXgloZQ4Xd05iyHKZEIQ513nVntcOAfQgW7n42tw/kGHrgZ+pvL9B3aryAcFp27rM1hMT9YCwAua3WrHAEAAA=="

func dataMsgForChannel(ch string) []byte {
	msg, _ := base64.StdEncoding.DecodeString(sampleDataMsgBase64)
	for i, j := 4, 0; i < 16; i, j = i+1, j+1 {
		if len(ch) < j {
			msg[i] = ch[j]
		} else {
			msg[i] = 0
		}
	}

	return msg
}

var upgrader = websocket.Upgrader{} // use default options

type mockHandler struct {
	received         []map[string]interface{}
	validAccessToken string
}

func (s *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	defer cancel()

	textMsgs := make(chan string)
	binMsgs := make(chan []byte)
	go func() {
		for {
			var err error
			select {
			case m := <-textMsgs:
				err = c.WriteMessage(websocket.TextMessage, []byte(m))
			case m := <-binMsgs:
				err = c.WriteMessage(websocket.BinaryMessage, m)
			case <-ctx.Done():
				return
			}
			if err != nil {
				log.Printf("Could not write message: %v", err)
			}
		}
	}()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read err:", err)
			break
		}

		var in map[string]interface{}
		if err := json.Unmarshal(message, &in); err != nil {
			log.Println("error unmarshalling: ", err)
		}
		s.received = append(s.received, in)

		typ, ok := in["type"].(string)
		if !ok {
			c.WriteMessage(mt, []byte(`{"type": "error"}`))
			continue
		}

		switch typ {
		case "authenticate":
			token, _ := in["token"].(string)
			if s.validAccessToken == "" || token == s.validAccessToken {
				textMsgs <- `{"type": "authenticated"}`
			} else {
				textMsgs <- `{"type": "error", "message": "Invalid auth token"}`
				return
			}
		case "execute":
			ch, _ := in["channel"].(string)
			resMs, _ := in["resolution"].(float64)
			if resMs == 0 {
				resMs = 1000
			}
			textMsgs <- fmt.Sprintf(`{"type": "control-message", "channel": "%s", "event": "STREAM_START"}`, ch)
			textMsgs <- fmt.Sprintf(`{"type": "control-message", "channel": "%s", "event": "JOB_START", "handle": "asdf"}`, ch)
			textMsgs <- fmt.Sprintf(`{"type": "message", "channel": "%s", "logicalTimestampMs": 1464736034000, "message": {"contents": {"resolutionMs" : %d}, "messageCode": "JOB_RUNNING_RESOLUTION", "timestampMs": 1464736033000}}`, ch, int64(resMs))
			go func() {
				t := time.NewTicker(1 * time.Second)
				for {
					select {
					case <-ctx.Done():
						return
					case <-t.C:
						binMsgs <- dataMsgForChannel(ch)
					}
				}
			}()
		}
	}
}

type closeFunc func()

func runMockServer() (*mockHandler, string, closeFunc) {
	handler := &mockHandler{}
	s := httptest.NewServer(handler)
	return handler, strings.Replace(s.URL, "http", "ws", 1), s.Close
}

func TestAuthenticationFlow(t *testing.T) {
	handler, url, closer := runMockServer()
	defer closer()

	handler.validAccessToken = "testing123"
	c, err := NewClient(StreamURL(url), AccessToken(handler.validAccessToken))
	require.Nil(t, err)

	comp, err := c.Execute(&ExecuteRequest{
		Program: "data('cpu.utilization').publish()",
	})
	require.Nil(t, err)

	require.Equal(t, 1*time.Second, comp.Resolution())

	require.Equal(t, []map[string]interface{}{
		{"type": "authenticate",
			"token": "testing123"},
		{"type": "execute",
			"channel":    "ch-1",
			"immediate":  false,
			"maxDelay":   0.,
			"program":    "data('cpu.utilization').publish()",
			"resolution": 0.,
			"start":      0.,
			"stop":       0.,
			"timezone":   ""},
	}, handler.received)
}

func TestMultipleComputations(t *testing.T) {
	handler, url, closer := runMockServer()
	defer closer()

	handler.validAccessToken = "testing123"
	c, err := NewClient(StreamURL(url), AccessToken(handler.validAccessToken))
	require.Nil(t, err)

	for i := 1; i < 50; i++ {
		comp, err := c.Execute(&ExecuteRequest{
			Program:    "data('cpu.utilization').publish()",
			Resolution: time.Duration(i) * time.Second,
		})
		require.Nil(t, err)

		require.Equal(t, time.Duration(i)*time.Second, comp.Resolution())
		require.Equal(t, fmt.Sprintf("ch-%d", i), comp.Channel().name)
	}

}
