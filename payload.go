// Copyright 2018 Fetch Robotics, Inc.
// Author(s): Pavan Soundara

package main

import (
	"std_msgs"
)

// Message struct of message data sent through socket.
type Message struct {
	Type    string   `json:"type"`
	Payload *Payload `json:"payload"`
}

// Payload defines the data format for data coming from stream consumer and stream producer.
type Payload struct {
	StreamURL  string      `json:"stream"`
	Customer   string      `json:"customer"`
	ProducerID string      `json:"producer"`
	Data       interface{} `json:"data"`
}

// BatteryData defines the data format for battery data.
type BatteryState struct {
	Name		string  `json:"name"`
	Timestamp       int64   `json:"timestamp"`
	ChargeLevel     float32 `json:"charge_level"`
	IsCharging      bool    `json:"is_charging"`
	TotalCapacity   float32 `json:"total_capacity"`
	CurrentCapacity float32 `json:"current_capacity"`
	BatteryVoltage  float32 `json:"battery_voltage"`
	SupplyVoltage   float32 `json:"supply_voltage"`
	ChargerVoltage  float32 `json:"charger_voltage"`
}

// StringData define the data format for string data.
type StringData struct {
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
}

type LaserScan struct {
	Timestamp      int64           `json:"timestamp"`
	Header         std_msgs.Header `json:"header"`
	AngleMin       float32         `json:"angle_min"`
	AngleMax       float32         `json:"angle_max"`
	AngleIncrement float32         `json:"angle_increment"`
	TimeIncrement  float32         `json:"time_increment"`
	ScanTime       float32         `json:"scan_time"`
	RangeMin       float32         `json:"range_min"`
	RangeMax       float32         `json:"range_max"`
	Ranges         []float32       `json:"ranges"`
	Intensities    []float32       `json:"intensities"`
}
