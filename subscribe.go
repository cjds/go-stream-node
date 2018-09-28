// Copyright 2018 Fetch Robotics, Inc.
// Author(s): Pavan Soundara

package main

import (
	"fmt"
	"os"
	"power_msgs"
	
	"github.com/sirupsen/logrus"
	"github.com/akio/rosgo/ros"
	"github.com/patrickmn/go-cache"
)

// Subscribe to topics and reads data
type Subscribe struct {
	Topics	[]string
	Write	chan *power_msgs.BatteryState
	Stop	chan bool
	Store	*cache.Cache
}

// Create a new subscriber and return it
func NewSubscriber(store *cache.Cache) *Subscribe {
	logrus.Info("[Subscribe] New subscriber called")
	sc := &Subscribe{
		Topics:	make([] string, 5),
		Write: 	make(chan *power_msgs.BatteryState),
		Stop:	make(chan bool),
		Store:	store,
	}
	return sc
}

// Create new listener. TODO: divide listener and subscriber 
// into different methods
func (sc *Subscribe) newListener(){
	logrus.Info("[Subscribe] New Listener created")
	node, err := ros.NewNode("/listener", os.Args)
	if err != nil {
		logrus.Info(err)
		os.Exit(-1)
	}
	defer node.Shutdown()
	node.Logger().SetSeverity(ros.LogLevelDebug)
	node.NewSubscriber("/chatter", power_msgs.MsgBatteryState, (*sc).readData)
	node.Spin()
}

// Read data and check for token in cache before sending data
func (sc *Subscribe) readData(msg *power_msgs.BatteryState){
	logrus.Info("[Subscribe] Name:", msg.Name)
	logrus.Info("[Subscribe] Charge_level:", msg.ChargeLevel)
	logrus.Info("[Subscribe] is_charging:", msg.IsCharging)
	logrus.Info("[Subscribe] remaining_time:", msg.RemainingTime)
	t, found := (*sc).Store.Get("token")
	if found {
		logrus.Info("[Subscribe] Found a token in cache")
		logrus.Info("[Subscribe] Token is:", t)
		//Todo send data async via socket
	}
}

// Add topics to Topics value in the struct
func (sc *Subscribe) addTopic(t string){
	(*sc).Topics = append((*sc).Topics, t)
}

// Print all the topics in struct
func (sc *Subscribe) printAllTopics(){
	fmt.Println((*sc).Topics)
}

