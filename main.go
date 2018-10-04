// Copyright 2018 Fetch Robotics, Inc.
// Author(s): Pavan Soundara

package main

// Generates message structs. Run `go generate` in the directory.
//go:generate gengo msg power_msgs/BatteryState
//go:generate gengo msg std_msgs/String

import (
	"flag"
	"power_msgs"
	"std_msgs"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	conf = flag.String(
		"conf",
		"conf/development",
		"Directory in which to find the application.toml file.",
	)
)

func init() {
	viper.SetConfigName("application")
	viper.SetConfigType("toml")
}

func main() {
	viper.AddConfigPath(*conf)
	if err := viper.ReadInConfig(); err != nil {
		logrus.Errorf("Config file not found: %s \n", err)
	}

	c := cache.New(1*time.Hour, 2*time.Hour)

	a := NewAuthManager(c)
	go a.setTokenInCache()

	t := <-a.Connect
	if !t.Connected {
		logrus.Error(t.Err)
		logrus.Fatal("Reached maximum number of retries.")
	}

	s := NewSubscriber(c)
	n, err := s.newNode()
	if err != nil {
		logrus.Fatal(err)
	}
	defer n.Shutdown()

	go s.newListener("chatteone", std_msgs.MsgString, n)
	go s.newListener("chattertwo", power_msgs.MsgBatteryState, n)

	for t := range a.Connect {
		if !t.Connected {
			logrus.Error(t.Err)
			logrus.Fatal("Reached maximum number of retries.")
		}
	}
}
