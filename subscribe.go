// Copyright 2018 Fetch Robotics, Inc.
// Author(s): Pavan Soundara

package main

import (
	"bytes"
	"context"
	"net/url"
	"os"
	"power_msgs"
	"std_msgs"
	"time"

	"github.com/akio/rosgo/ros"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 128000
)

// MsgMap contains a map for ros message types.
var MsgMap = map[ros.MessageType]string{
	power_msgs.MsgBatteryState: "battery",
	std_msgs.MsgString:         "string",
}

// SubscriberManager struct.
type SubscriberManager struct {
	ID      string
	Auth    *AuthManager
	Conn    *websocket.Conn
	Streams map[string]ros.MessageType
	Send    chan *bytes.Buffer
}

// NewSubscriber creates a new SubscriberManager and returns it.
func NewSubscriber(id string, a *AuthManager, conn *websocket.Conn, streams map[string]ros.MessageType) *SubscriberManager {
	sm := &SubscriberManager{
		ID:      id,
		Auth:    a,
		Conn:    conn,
		Streams: streams,
		Send:    make(chan *bytes.Buffer, 5),
	}
	return sm
}

// newNode creates and returns a new ros node.
func (sm *SubscriberManager) newNode(ctx context.Context) {
	node, err := ros.NewNode("/listener", os.Args)
	if err != nil {
		logrus.Info(err)
		return
	}

	node.Logger().SetSeverity(ros.LogLevelInfo)
	childCtx, cancel := context.WithCancel(ctx)

	for {
		select {
		case <-ctx.Done():
			logrus.Info("[Subscribe] Shutting down readPump, writePump, websocket and node and exiting")
			cancel()
			sm.Conn.Close()
			node.Shutdown()
			return
		case <-sm.Auth.Connect:
			sm.connectToSocket()
			go sm.readPump(childCtx)
			go sm.writePump(childCtx)
			go sm.createNewListeners(node)
		}
	}
}

// createNewListeners adds new listeners to the ros node for all the streams.
func (sm *SubscriberManager) createNewListeners(n ros.Node) {
	for k, v := range sm.Streams {
		// n.NewSubscriber creates new subscribers.
		// fun() implements the third parameter which is a callback interface{}.
		// Whenever a new message is emited, sm.readData routine is called 
		// with the message data.
		n.NewSubscriber("/"+k, v, sm.readData)
	}
	n.Spin()
}

// readData reads incoming data from robot and checks for token in cache before sending data.
func (sm *SubscriberManager) readData(msg interface{}) {
	var mb bytes.Buffer
	var err error

	switch msg.(type) {
	case *power_msgs.BatteryState:
		m := msg.(*power_msgs.BatteryState)
		err = m.Serialize(&mb)
	case *std_msgs.String:
		m := msg.(*std_msgs.String)
		err = m.Serialize(&mb)
	default:
		logrus.Info("[Subscribe] Unsupported message type")
		return
	}

	if err != nil {
		logrus.Error(err)
		return
	}

	sm.Send <- &mb
}

// readPump reads responses from websocket.
// This is essential for maintaining the websocket connection when
// there is no data being tranferred and also for checking unexpected
// closure of socket connection by streams server.
func (sm *SubscriberManager) readPump(ctx context.Context) {
	sm.Conn.SetReadLimit(maxMessageSize)
	sm.Conn.SetReadDeadline(time.Now().Add(pongWait))
	sm.Conn.SetPongHandler(func(string) error {
		sm.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		select {
		case <-ctx.Done():
			logrus.Info("[Subscribe] readPump cancelled")
			return
		default:
			_, _, err := sm.Conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					logrus.Fatal("Websocket has unexpectedly closed: ", err)
				}
				logrus.Error("[Subscribe] Websocket reading error:", err)
			}
		}
	}
}

// writePump upstreams payload data through socket.
// Ticker is used to ping streams server and is essential to
// maintain the websocket connection when no data is being sent.
func (sm *SubscriberManager) writePump(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logrus.Info("[Subscribe] writePump cancelled")
			return
		case message := <-sm.Send:
			sm.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := sm.Conn.WriteMessage(websocket.BinaryMessage, message.Bytes())
			if err != nil {
				return
			}
		case <-ticker.C:
			sm.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := sm.Conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				return
			}
		}
	}
}

// connectToSocket Creates the initital websocket connection to streams server.
func (sm *SubscriberManager) connectToSocket() {
	if sm.Conn != nil {
		sm.Conn.Close()
	}

	s := sm.getWebsocketURL()
	u, _ := url.Parse(s)

	token, err := sm.Auth.checkToken()
	if err != nil {
		logrus.Fatal(err)
	}

	sockURL := u.String() + sm.getNewStreamParam() + "&token=" + token
	logrus.Infof("[Subscribe] Connecting to %s", sockURL)

	(*sm).Conn, _, err = websocket.DefaultDialer.Dial(sockURL, nil)
	if err != nil {
		logrus.Fatal("[Subscribe] Error connecting to websocket:", err)
	}

	logrus.Info("[Subscribe] Connected to websocket")
}

// getWebsocketURL generates URL to reach stream_server websocket.
func (sm *SubscriberManager) getWebsocketURL() string {
	host := viper.GetString("streams_server.host")
	port := viper.GetString("streams_server.port")
	api := viper.GetString("streams_server.api_uri")
	return "ws://" + host + ":" + port + api + "?" + "id=" + sm.ID
}

// getNewStreamParam generates URL params for all the streams.
func (sm *SubscriberManager) getNewStreamParam() string {
	var sb bytes.Buffer
	for _, v := range sm.Streams {
		sb.WriteString("&streams=/" + sm.ID + "/sensor/" + MsgMap[v])
	}
	return sb.String()
}
