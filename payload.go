// Copyright 2018 Fetch Robotics, Inc.
// Author(s): Pavan Soundara

package main

import ()

type Message struct {
	Type    string   `json:"type"`
	Payload *Payload `json:"payload"`
}

// Payload defines the data format for data coming from stream consumer and stream producer.
type Payload struct {
	StreamURL  string       `json:"stream"`
	Customer   string       `json:"customer"`
	ProducerID string       `json:"producer"`
	Data       *BatteryData `json:"data"`
}

// BatteryData defines the data format for battery data.
type BatteryData struct {
	Timestamp int     `json:"timestamp"`
	Percent   float32 `json:"percent"`
}
