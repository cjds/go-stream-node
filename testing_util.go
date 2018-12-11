// Copyright 2018 Fetch Robotics, Inc.
// Author(s): Pavan Soundara

package main

import (
	"go-stream-node/messages"
	"math/rand"
	"os"
	"power_msgs"
	"sensor_msgs"
	"std_msgs"
	"testing"
	"time"

	"github.com/akio/rosgo/ros"
	"github.com/golang/protobuf/proto"
)

func newRosNode(t *testing.T, nodeType string) (ros.Node, error) {
	node, err := ros.NewNode("/"+nodeType, os.Args)
	if err != nil {
		t.Error(err)
		return node, err
	}

	node.Logger().SetSeverity(ros.LogLevelInfo)

	return node, nil
}

func publishRosMsg(t *testing.T, pub ros.Publisher, msgType ros.MessageType) {
	switch msgType {
	case power_msgs.MsgBatteryState:
		msg := generateBatteryState()
		pub.Publish(msg)

	case sensor_msgs.MsgLaserScan:
		msg := generateLaserScan()
		pub.Publish(msg)

	}
}

func recvProtoMsg(t *testing.T, sub *SubscriberManager) int {
	msgCount := 0
	timer := time.After(6 * time.Second)
loop:
	for {
		select {
		case <-timer:
			break loop
		case msg := <-sub.Send:
			messagePayLoad := &messages.PayLoad{}
			err := proto.Unmarshal(msg, messagePayLoad)
			if err != nil {
				t.Error(err)
				continue
			}

			msgCount++
		}
	}

	return msgCount
}

func makeJsonPayload(b *testing.B, msg ros.Message) *Payload {
	p := &Payload{}

	switch msg.(type) {
	case *power_msgs.BatteryState:
			m := msg.(*power_msgs.BatteryState)
			p.StreamURL = "/freight20/sensor/battery"
			p.Data = &BatteryState{
					Timestamp:       time.Now().UnixNano() / int64(time.Millisecond),
					Name:            m.Name,
					IsCharging:      m.IsCharging,
					TotalCapacity:   m.TotalCapacity,
					CurrentCapacity: m.CurrentCapacity,
					BatteryVoltage:  m.BatteryVoltage,
					SupplyVoltage:   m.SupplyVoltage,
					ChargerVoltage:  m.ChargerVoltage,
			}
	case *std_msgs.String:
			m := msg.(*std_msgs.String)
			p.StreamURL = "/freight20/sensor/string"
			p.Data = &StringData{
					Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
					Message:   m.Data,
			}
	case *sensor_msgs.LaserScan:
			m := msg.(*sensor_msgs.LaserScan)
			p.StreamURL = "/freight20/sensor/rawlaser"
			p.Data = &LaserScan{
					Timestamp:      time.Now().UnixNano() / int64(time.Millisecond),
					AngleMin:       m.AngleMin,
					AngleMax:       m.AngleMax,
					AngleIncrement: m.AngleIncrement,
					TimeIncrement:  m.TimeIncrement,
					ScanTime:       m.ScanTime,
					RangeMin:       m.RangeMin,
					RangeMax:       m.RangeMax,
					Ranges:         m.Ranges,
					Intensities:    m.Intensities,
			}
	default:
		return nil
	}

	return p
}

func makeProtoPayload(b *testing.B, msg ros.Message) *messages.PayLoad {
	pl := &messages.PayLoad{}
	timeNow := time.Now().UnixNano() / int64(time.Millisecond)
	ID := "freight20"

	switch msg.(type) {
	case *power_msgs.BatteryState:
		m := msg.(*power_msgs.BatteryState)
		pl.Stream = "/" + ID + "/sensor/battery"
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
		pl.Stream = "/" + ID + "/sensor/rawlaser"
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
		pl.Stream = "/" + ID + "/sensor/string"
		pl.Data = &messages.PayLoad_StringMessage{
			&messages.String{
				Timestamp: timeNow,
				Data:      m.Data,
			},
		}

	default:
		return nil
	}

	return pl
}

func generateBatteryState() *power_msgs.BatteryState {
	rand.Seed(time.Now().UnixNano())
	return &power_msgs.BatteryState{
		Name: "Test Battery",
		ChargeLevel: randFloat(0, 100),
		IsCharging: true,
		TotalCapacity: 193190.609375,
		CurrentCapacity: randFloat(0, 193190.609375),
		BatteryVoltage: randFloat(0, 50),
		SupplyVoltage: randFloat(0, 50),
		ChargerVoltage: randFloat(0, 50),
	}
}

func generateLaserScan() *sensor_msgs.LaserScan {
	rand.Seed(time.Now().UnixNano())
	return &sensor_msgs.LaserScan{
		AngleMin: randFloat(-1.5, 1.9),
		AngleMax: randFloat(1, 1.9),
		AngleIncrement: randFloat(-1.5, 1.9),
		TimeIncrement: randFloat(0, 6),
		ScanTime: randFloat(0, 1),
		RangeMin: randFloat(0, 0.1),
		RangeMax: 25.0,
		Ranges: randFloats(0.0, 25.0, 662),
		Intensities: randFloats(55.0, 250.0, 662),
	}
}

func randFloats(min, max float32, n int) []float32 {
    res := make([]float32, n)
    for i := range res {
        res[i] = min + rand.Float32() * (max - min)
    }
    return res
}

func randFloat(min, max float32) float32 {
	return min + rand.Float32() * (max - min)
}
