// Copyright 2018 Fetch Robotics, Inc.
// Author(s): Pavan Soundara

package main

import (
	"fmt"
	"os"
	"power_msgs"
	
	"github.com/sirupsen/logrus"
	"github.com/akio/rosgo/ros"
)

// Subscribes to topics and reads data
type Subscribe struct {
	Topics	[]string
	Write	chan *power_msgs.BatteryState
	Stop	chan bool
}

func NewSubscriber(topics []string, ch chan *power_msgs.BatteryState) {
	logrus.Info("Subscribed!")
	sc := &Subscribe{
		Topics:	topics,
		Write: 	ch,
		Stop:	make(chan bool),
	}
	go sc.newListener()
}

func (sc *Subscribe) newListener(){
	logrus.Info("New Listener created")
	node, err := ros.NewNode("/listener", os.Args)
        if err != nil {
                fmt.Println(err)
                os.Exit(-1)
        }

        defer node.Shutdown()
        node.Logger().SetSeverity(ros.LogLevelDebug)
        node.NewSubscriber("/chatter", power_msgs.MsgBatteryState, sc.readData)
	node.Spin()
}

func (sc *Subscribe) readData(msg *power_msgs.BatteryState){
	//TODO: Add provision to check cache for token
	sc.Write <- msg
}



