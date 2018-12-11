// Copyright 2018 Fetch Robotics, Inc.
// Author(s): Pavan Soundara

package main

import (
	"bytes"
	"compress/zlib"
	"context"
	"go-stream-node/messages"
	"net/url"
	"os"
	"power_msgs"
	"sensor_msgs"
	"std_msgs"
	"time"

	"github.com/akio/rosgo/ros"
	"github.com/golang/protobuf/proto"
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
	sensor_msgs.MsgLaserScan:   "raw_laser",
}

// SubscriberManager struct.
type SubscriberManager struct {
	ID      string
	Auth    *AuthManager
	Conn    *websocket.Conn
	Streams map[string]ros.MessageType
	Send    chan []byte
}

// NewSubscriber creates a new SubscriberManager and returns it.
func NewSubscriber(id string, a *AuthManager, conn *websocket.Conn, streams map[string]ros.MessageType) *SubscriberManager {
	sm := &SubscriberManager{
		ID:      id,
		Auth:    a,
		Conn:    conn,
		Streams: streams,
		Send:    make(chan []byte, 100),
	}

	return sm
}

// Start creates a new ros node, adds new listeners and starts socket routines.
// On token refreshal, restarts read-write pump and connectToSocket.
func (sm *SubscriberManager) Start(ctx context.Context) {
	node, err := ros.NewNode("/listener", os.Args)
	if err != nil {
		logrus.Info(err)
		return
	}

	defer node.Shutdown()
	node.Logger().SetSeverity(ros.LogLevelInfo)
	go sm.createNewListeners(node)

	var cancel context.CancelFunc
	var childCtx context.Context

	for {
		select {
		case <-ctx.Done():
			return

		case <-sm.Auth.Connect:
			if cancel != nil {
				cancel()
			}
			childCtx, cancel = context.WithCancel(ctx)
			sm.connectToSocket()
			go sm.readPump(childCtx)
			go sm.writePump(childCtx)
		}
	}
}

// createNewListeners adds new listeners to the ros node for all the streams.
func (sm *SubscriberManager) createNewListeners(n ros.Node) {
	for k, v := range sm.Streams {
		// n.NewSubscriber creates new subscribers.
		// func() implements the third parameter which is a callback interface{}.
		// Whenever a new message is emited, sm.readData routine is called
		// with the message data.
		n.NewSubscriber("/"+k, v, func(msg ros.Message) { go sm.readData(msg) })
	}

	defer n.Shutdown()
	n.Spin()
}

// readData reads incoming data from robot and checks for token in cache before sending data.
// Any ros.Message can not be directly converted or copied to a messages.{msgType} struct since
// messages.{msgType} struct contains extra fields for unknown values which is auto generated.
func (sm *SubscriberManager) readData(msg ros.Message) {
	pl := &messages.PayLoad{}
	timeNow := time.Now().UnixNano() / int64(time.Millisecond)

	switch msg.(type) {
	case *power_msgs.BatteryState:
		m := msg.(*power_msgs.BatteryState)
		pl.Stream = "/" + sm.ID + "/sensor/battery"
		pl.Data = &messages.PayLoad_BatteryState{
			&messages.BatteryState{
				Timestamp:       timeNow,
				Name:            m.Name,
				IsCharging:      m.IsCharging,
				TotalCapacity:   m.TotalCapacity,
				CurrentCapacity: m.CurrentCapacity,
				BatteryVoltage:  m.BatteryVoltage,
				SupplyVoltage:   m.SupplyVoltage,
				ChargerVoltage:  m.ChargerVoltage,
			},
		}

	case *sensor_msgs.LaserScan:
		m := msg.(*sensor_msgs.LaserScan)
		pl.Stream = "/" + sm.ID + "/sensor/raw_laser"
		pl.Data = &messages.PayLoad_LaserScan{
			&messages.LaserScan{
				Timestamp:      timeNow,
				AngleMin:       m.AngleMin,
				AngleMax:       m.AngleMax,
				AngleIncrement: m.AngleIncrement,
				TimeIncrement:  m.TimeIncrement,
				ScanTime:       m.ScanTime,
				RangeMin:       m.RangeMin,
				RangeMax:       m.RangeMax,
				Ranges:         m.Ranges,
				Intensities:    m.Intensities,
			},
		}

	case *std_msgs.String:
		m := msg.(*std_msgs.String)
		pl.Stream = "/" + sm.ID + "/sensor/string"
		pl.Data = &messages.PayLoad_StringMessage{
			&messages.String{
				Timestamp: timeNow,
				Data:      m.Data,
			},
		}

	default:
		logrus.Error("[Subscribe] Unsupported message type")
		return
	}

	message, err := proto.Marshal(pl)
	if err != nil {
		logrus.Warn(err)
		return
	}

	sm.Send <- message
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
// We initialize and reuse zlib compressor to compress the byte data.
func (sm *SubscriberManager) writePump(ctx context.Context) {
	var cd bytes.Buffer
	w, _ := zlib.NewWriterLevel(&cd, 9)

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case message := <-sm.Send:
			cd.Reset()
			w.Reset(&cd)
			w.Write(message)
			w.Close()

			sm.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := sm.Conn.WriteMessage(websocket.BinaryMessage, cd.Bytes())
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

// connectToSocket Creates the websocket connection to streams server.
// When called more than once, will close the existing connection and creates
// a new connection with the latest available token in cache.
func (sm *SubscriberManager) connectToSocket() {
	if sm.Conn != nil {
		sm.Conn.Close()
	}

	token, err := sm.Auth.checkToken()
	if err != nil {
		logrus.Fatal(err)
	}

	sockURL := sm.getWebsocketURL(token)
	logrus.Infof("[Subscribe] Connecting to %s", sockURL.String())

	sm.Conn, _, err = websocket.DefaultDialer.Dial(sockURL.String(), nil)
	if err != nil {
		logrus.Fatal("[Subscribe] Error connecting to websocket:", err)
	}

	logrus.Info("[Subscribe] Connected to websocket")
}

// getWebsocketURL generates URL to reach stream_server websocket.
func (sm *SubscriberManager) getWebsocketURL(token string) url.URL {
	host := viper.GetString("streams_server.host")
	port := viper.GetString("streams_server.port")
	api := viper.GetString("streams_server.api_uri")

	sockURL := url.URL{Scheme: "ws", Host: host + ":" + port, Path: api}
	query := sockURL.Query()
	query.Set("id", sm.ID)
	query.Set("token", token)

	for _, v := range sm.Streams {
		query.Add("streams", "/"+sm.ID+"/sensor/"+MsgMap[v])
	}

	sockURL.RawQuery = query.Encode()

	return sockURL
}
