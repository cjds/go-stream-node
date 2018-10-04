// Copyright 2018 Fetch Robotics, Inc.
// Author(s): Pavan Soundara

package main

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"power_msgs"
	"std_msgs"

	"github.com/akio/rosgo/ros"
	"github.com/gorilla/websocket"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var MsgMap = map[ros.MessageType]string{
	power_msgs.MsgBatteryState: "battery",
	std_msgs.MsgString:         "string",
}

// Subscribe to topics and read data.
type SubscriberManager struct {
	Streams   []string
	Stop      chan bool
	AddStream chan string
	Store     *cache.Cache
}

// Create a new subscriber and return it.
func NewSubscriber(store *cache.Cache) *SubscriberManager {
	logrus.Info("[Subscribe] New subscriber called")
	sm := &SubscriberManager{
		Streams:   make([]string, 0),
		Stop:      make(chan bool),
		AddStream: make(chan string),
		Store:     store,
	}
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
	n.NewSubscriber("/"+topic, msgType, (*sm).readData)
	(*sm).Streams = append((*sm).Streams, MsgMap[msgType])
	(*sm).AddStream <- MsgMap[msgType]
	logrus.Info("[Subscribe] New Listener created")
	n.Spin()
}

// Read data and check for token in cache before sending data.
func (sm *SubscriberManager) readData(msg interface{}) {
	switch msg.(type) {
	case *power_msgs.BatteryState:
		m := msg.(*power_msgs.BatteryState)
		logrus.Info("[Subscribe] Name:", m.Name)
	case *std_msgs.String:
		m := msg.(*std_msgs.String)
		logrus.Info("[Subscribe] String Message:", m)
	default:
		logrus.Info("[Subscribe] Unsupported message type")
	}
}

// A go routine to connect to websocket server.
// TODO: Move connected and c to struct. Change name of func.
// Do we actually require add stream?
func (sm *SubscriberManager) connect() {
	sockURI := getWebsocketURL()
	u, _ := url.Parse(sockURI)

	var c *websocket.Conn
	var err error
	connected := false
	defer c.Close()

	for range (*sm).AddStream {
		if connected {
			logrus.Info("Closing socket connection to add more streams")
			c.Close()
			connected = false
		}
		sockURI := u.String() + (*sm).getNewStreamParam()
		logrus.Infof("Connecting to %s", sockURI)
		c, _, err = websocket.DefaultDialer.Dial(sockURI, nil)
		if err != nil {
			logrus.Fatal("dial:", err)
		}
		connected = true
		logrus.Info("[Subscribe] Connected to web socket")
	}
}

// Generates URL to reach stream_server's websocket.
func getWebsocketURL() string {
	proId := viper.GetString("auth.username")
	host := viper.GetString("streams_server.host")
	port := viper.GetString("streams_server.port")
	api := viper.GetString("streams_server.api_uri")
	return "ws://" + host + ":" + port + api + "?" + "id=" + proId
}

// Generates url Param for all the subscribed streams.
func (sm *SubscriberManager) getNewStreamParam() string {
	proId := viper.GetString("auth.username")
	var sb bytes.Buffer
	for _, stream := range (*sm).Streams {
		sb.WriteString("&streams=/" + proId + "/sensor/" + stream)
	}
	return sb.String()
}

// Check for token in cache.
func (sm *SubscriberManager) checkToken() (string, error) {
	t, found := (*sm).Store.Get("token")
	if !found {
		return "", fmt.Errorf("Token not found in cache")
	}
	return t.(string), nil
}
