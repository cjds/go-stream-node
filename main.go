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

	"github.com/akio/rosgo/ros"
	"github.com/gorilla/websocket"
	"github.com/pkg/profile"
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
	mode := flag.String("profile.mode", "", "enable profiling mode, one of [cpu, memory]")
	flag.Parse()
	switch *mode {
	case "cpu":
		defer profile.Start(profile.CPUProfile, profile.ProfilePath("./profiles/cpu")).Stop()
	case "memory":
		defer profile.Start(profile.MemProfile, profile.ProfilePath("./profiles/memory")).Stop()
	default:
	}

	viper.AddConfigPath(*conf)
	if err := viper.ReadInConfig(); err != nil {
		logrus.Errorf("Config file not found: %s \n", err)
	}

	id := viper.GetString("auth.username")
	var conn *websocket.Conn

	// streams keys represent /topic and values represent message type.
	// Add items to this map to add listeners.
	streams := map[string]ros.MessageType{
		"string_messages": std_msgs.MsgString,
		"battery_state":   power_msgs.MsgBatteryState}

	a, ctx := NewAuthManager()
	s := NewSubscriber(id, a, conn, streams)
	s.Start(ctx)
}
