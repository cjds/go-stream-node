// Copyright 2018 Fetch Robotics, Inc.
// Author(s): Pavan Soundara

package main

import (
	"bytes"
	"encoding/json"
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
	pingPeriod = (pongWait) / 10

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
	Send    chan []byte
	Stop    chan bool
}

// NewSubscriber creates a new SubscriberManager and returns it.
func NewSubscriber(id string, a *AuthManager, conn *websocket.Conn, streams map[string]ros.MessageType) *SubscriberManager {
	sm := &SubscriberManager{
		ID:      id,
		Auth:    a,
		Conn:    conn,
		Streams: streams,
		Send:    make(chan []byte, 5),
		Stop:    make(chan bool),
	}
	return sm
}

// newNode creates and returns a new ros node.
func (sm *SubscriberManager) newNode() (ros.Node, error) {
	node, err := ros.NewNode("/listener", os.Args)
	if err != nil {
		logrus.Info(err)
		return nil, err
	}
	sm.connectToSocket()
	go sm.readPump()
	go sm.writePump()
	sm.createNewListeners(node)
	return node, err
}

// createNewListeners adds new listeners to the ros node for all the streams.
func (sm *SubscriberManager) createNewListeners(n ros.Node) {
	for k, v := range sm.Streams {
		go sm.newListener(k, v, n)
	}
}

// newListener adds individual listeners to the ros node.
func (sm *SubscriberManager) newListener(topic string, msgType ros.MessageType, n ros.Node) {
	n.Logger().SetSeverity(ros.LogLevelDebug)
	n.NewSubscriber("/"+topic, msgType, sm.readData)
	n.Spin()
}

// readData reads incoming data from robot and checks for token in cache before sending data.
func (sm *SubscriberManager) readData(msg interface{}) {
	payload := Payload{}

	switch msg.(type) {
	case *power_msgs.BatteryState:
		m := msg.(*power_msgs.BatteryState)
		payload.StreamURL = "/" + sm.ID + "/sensor/battery"
		payload.Data = &BatteryData{
			Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
			Percent:   m.ChargeLevel,
		}
	case *std_msgs.String:
		m := msg.(*std_msgs.String)
		payload.StreamURL = "/" + sm.ID + "/sensor/string"
		payload.Data = &StringData{
			Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
			Message:   m.Data,
		}
	default:
		logrus.Info("[Subscribe] Unsupported message type")
		return
	}

	payload.Customer = "NoCustomer"
	payload.ProducerID = sm.ID

	message := Message{}
	message.Type = "publish"
	message.Payload = &payload
	m, err := json.Marshal(message)
	if err != nil {
		logrus.Error("Error marshalling:", err)
	}

	sm.Send <- m
}

// readPump reads responses from websocket.
// This is essential for maintaining the websocket connection when
// there is no data being tranferred and also for checking unexpected
// closure of socket connection by streams server.
func (sm *SubscriberManager) readPump() {
	sm.Conn.SetReadLimit(maxMessageSize)
	sm.Conn.SetReadDeadline(time.Now().Add(pongWait))
	sm.Conn.SetPongHandler(func(string) error {
		sm.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := sm.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				logrus.Fatal("Websocket has unexpectedly closed: ", err)
			}
			logrus.Warn("[Subscribe] Websocket reading error:", err)
		}
	}
}

// writePump upstreams payload data through socket.
// Ticker is used to ping streams server and is essential to
// maintain the websocket connection when no data is being sent.
func (sm *SubscriberManager) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case message := <-sm.Send:
			sm.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := sm.Conn.WriteMessage(websocket.TextMessage, message)
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

// Connect to websocket with all subscribed streams as URL params.
func (sm *SubscriberManager) connectToSocket() {
	s := sm.getWebsocketURL()
	u, _ := url.Parse(s)
	var err error

	sockURL := u.String() + sm.getNewStreamParam()
	logrus.Infof("[Subscribe] Connecting to %s", sockURL)
	(*sm).Conn, _, err = websocket.DefaultDialer.Dial(sockURL, nil)
	if err != nil {
		logrus.Fatal("[Subscribe] Error connecting to websocket:", err)
	}
	logrus.Info("[Subscribe] Connected to websocket")
}

// Generate URL to reach stream_server websocket.
func (sm *SubscriberManager) getWebsocketURL() string {
	host := viper.GetString("streams_server.host")
	port := viper.GetString("streams_server.port")
	api := viper.GetString("streams_server.api_uri")
	return "ws://" + host + ":" + port + api + "?" + "id=" + sm.ID
}

// Generate URL Params for all the subscribed streams.
func (sm *SubscriberManager) getNewStreamParam() string {
	var sb bytes.Buffer
	for _, v := range sm.Streams {
		sb.WriteString("&streams=/" + sm.ID + "/sensor/" + MsgMap[v])
	}
	return sb.String()
}
