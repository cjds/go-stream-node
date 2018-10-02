// Copyright 2018 Fetch Robotics, Inc.
// Author(s): Pavan Soundara

package main

import (
	"os"
	"power_msgs"
	"std_msgs"

	"github.com/akio/rosgo/ros"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

// Subscribe to topics and read data
type SubscriberManager struct {
	Topics []string
	Stop   chan bool
	Store  *cache.Cache
}

// Create a new subscriber and return it
func NewSubscriber(store *cache.Cache) *SubscriberManager {
	logrus.Info("[Subscribe] New subscriber called")
	sm := &SubscriberManager{
		Topics: make([]string, 5),
		Stop:   make(chan bool),
		Store:  store,
	}
	return sm
}

func (sm *SubscriberManager) newNode() ros.Node {
	node, err := ros.NewNode("/listener", os.Args)
	if err != nil {
		logrus.Info(err)
		os.Exit(-1)
	}
	return node
}

// Create new listener. TODO: divide listener and subscriber
// into different methods
func (sm *SubscriberManager) newListener(topic string, msgType ros.MessageType, n ros.Node) {
	n.Logger().SetSeverity(ros.LogLevelDebug)
	n.NewSubscriber("/"+topic, msgType, (*sm).readData)
	n.Spin()

	(*sm).Topics = append((*sm).Topics, topic)
	logrus.Info("[Subscribe] New Listener created")
}

// Read data and check for token in cache before sending data
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

//A go routine to connect to websocket server
func connect() {

}

func (sm *SubscriberManager) checkToken() {
	t, found := (*sm).Store.Get("token")
	if found {
		logrus.Info("[Subscribe] Found a token in cache")
		logrus.Info("[Subscribe] Token is:", t)
		//Todo send data async via socket
	}
}
