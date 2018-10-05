// Copyright 2018 Fetch Robotics, Inc.
// Author(s): Pavan Soundara

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"power_msgs"
	"std_msgs"
	"time"

	"github.com/akio/rosgo/ros"
	"github.com/gorilla/websocket"
	"github.com/patrickmn/go-cache"
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

var MsgMap = map[ros.MessageType]string{
	power_msgs.MsgBatteryState: "battery",
	std_msgs.MsgString:         "string",
}

// Subscribe to topics and read data.
type SubscriberManager struct {
	ID        string
	Store     *cache.Cache
	Conn      *websocket.Conn
	Connected bool
	Send      chan []byte
	Streams   []string
	Stop      chan bool
	AddStream chan string
}

// Create a new subscriber and return it.
func NewSubscriber(id string, store *cache.Cache, conn *websocket.Conn) *SubscriberManager {
	sm := &SubscriberManager{
		ID:        id,
		Store:     store,
		Conn:      conn,
		Connected: false,
		Send:      make(chan []byte, 5),
		Streams:   make([]string, 0),
		Stop:      make(chan bool),
		AddStream: make(chan string),
	}
	go sm.connectToSocket()
	go sm.readPump()
	go sm.writePump()
	return sm
}

// Create new node and return it.
func (sm *SubscriberManager) newNode() (ros.Node, error) {
	node, err := ros.NewNode("/listener", os.Args)
	if err != nil {
		logrus.Info(err)
	}
	return node, err
}

// Create new listener.
func (sm *SubscriberManager) newListener(topic string, msgType ros.MessageType, n ros.Node) {
	n.Logger().SetSeverity(ros.LogLevelDebug)
	n.NewSubscriber("/"+topic, msgType, sm.readData)
	sm.Streams = append(sm.Streams, MsgMap[msgType])
	sm.AddStream <- MsgMap[msgType]
	n.Spin()
}

// Read data and check for token in cache before sending data.
// TODO: Constants for streamURLS
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

// TODO: Resolve race condition.
// Race condition exists when sm.Connected is not used
// to check the websocket connection status.
func (sm *SubscriberManager) readPump() {
	for {
		if sm.Connected {
			_, _, err := sm.Conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					logrus.Fatal("Websocket has unexpectedly closed: ", err)
				}
				logrus.Warn("[Subscribe] Websocket reading error:", err)
			}
		}
	}
}

// Race condition exists when sm.Connected is not used
// to check the websocket connection status.
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

// A go routine to connect to websocket server when new streams are added.
func (sm *SubscriberManager) connectToSocket() {
	s := sm.getWebsocketURL()
	u, _ := url.Parse(s)

	var err error
	defer sm.Conn.Close()

	for range sm.AddStream {
		if sm.Connected {
			logrus.Info("[Subscribe] Closing socket connection to add more streams")
			sm.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			sm.Conn.Close()
			sm.Connected = false
		}

		sockURL := u.String() + sm.getNewStreamParam()
		logrus.Infof("[Subscribe] Connecting to %s", sockURL)
		sm.Conn, _, err = websocket.DefaultDialer.Dial(sockURL, nil)
		if err != nil {
			logrus.Fatal("[Subscribe] Error connecting to websocket:", err)
		}

		sm.Conn.SetReadLimit(maxMessageSize)
		sm.Conn.SetReadDeadline(time.Now().Add(pongWait))
		sm.Conn.SetPongHandler(func(string) error {
			sm.Conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})
		sm.Connected = true
		logrus.Info("[Subscribe] Connected to websocket")
	}
}

// Generates URL to reach stream_server websocket.
func (sm *SubscriberManager) getWebsocketURL() string {
	host := viper.GetString("streams_server.host")
	port := viper.GetString("streams_server.port")
	api := viper.GetString("streams_server.api_uri")
	return "ws://" + host + ":" + port + api + "?" + "id=" + sm.ID
}

// Generates url Param for all the subscribed streams.
func (sm *SubscriberManager) getNewStreamParam() string {
	var sb bytes.Buffer
	for _, stream := range sm.Streams {
		sb.WriteString("&streams=/" + sm.ID + "/sensor/" + stream)
	}
	return sb.String()
}

// Check for token in cache.
func (sm *SubscriberManager) checkToken() (string, error) {
	t, found := (*sm).Store.Get("token")
	if !found {
		return "", fmt.Errorf("[Subscribe] Token not found in cache.")
	}
	return t.(string), nil
}
