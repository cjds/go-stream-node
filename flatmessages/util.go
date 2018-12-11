// Copyright 2018 Fetch Robotics, Inc.
// Author(s): Pavan Soundara

package flatmessages

import (
	"power_msgs"
	"std_msgs"
	"sensor_msgs"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/akio/rosgo/ros"
)


func MakeMsg(b *flatbuffers.Builder, msg ros.Message) []byte {
	switch msg.(type) {
	case *power_msgs.BatteryState:
		m := msg.(*power_msgs.BatteryState)
		return MakeBatteryState(b, m)
	case *sensor_msgs.LaserScan:
		m := msg.(*sensor_msgs.LaserScan)
		return MakeLaserScan(b, m)
	case *std_msgs.String:
		m := msg.(*std_msgs.String)
		return MakeMsgString(b, m)
	default:
		return nil
	}
}

func MakeMsgString(b *flatbuffers.Builder, msg *std_msgs.String) []byte {
	b.Reset()

        data := b.CreateByteString([]byte(msg.Data))
        StringStart(b)
        StringAddData(b, data)
        message_position := StringEnd(b)
        b.Finish(message_position)
	
	return b.Bytes[b.Head():]
}


func MakeBatteryState(b *flatbuffers.Builder, msg *power_msgs.BatteryState) []byte {
	b.Reset()
	
	name := b.CreateByteString([]byte(msg.Name))
	BatteryStateStart(b)
	BatteryStateAddName(b, name)
	BatteryStateAddChargeLevel(b, msg.ChargeLevel) 
	BatteryStateAddIsCharging(b, boolToByte(msg.IsCharging))
	BatteryStateAddTotalCapacity(b, msg.TotalCapacity)
	BatteryStateAddCurrentCapacity(b, msg.CurrentCapacity)
	BatteryStateAddBatteryVoltage(b, msg.BatteryVoltage)
	BatteryStateAddSupplyVoltage(b, msg.SupplyVoltage)
	BatteryStateAddChargerVoltage(b, msg.ChargerVoltage)
	message_position := BatteryStateEnd(b)
	b.Finish(message_position)

	return b.Bytes[b.Head():]
}

func MakeLaserScan(b *flatbuffers.Builder, msg *sensor_msgs.LaserScan) []byte {
	b.Reset()

        LaserScanStartRangesVector(b, len(msg.Ranges))
        for _, r := range msg.Ranges {
                b.PrependFloat32(r)
        }
	rangesArr := b.EndVector(len(msg.Ranges))

        LaserScanStartIntensitiesVector(b, len(msg.Ranges))
        for _, i := range msg.Ranges {
                b.PrependFloat32(i)
        }
        intensitiesArr := b.EndVector(len(msg.Intensities))

	LaserScanStart(b)
	LaserScanAddHeader(b, CreateHeader(b, msg.Header.Seq, msg.Header.Stamp.ToNSec()))
	LaserScanAddAngleMin(b, msg.AngleMin)
	LaserScanAddAngleMax(b, msg.AngleMax)
        LaserScanAddAngleIncrement(b, msg.AngleIncrement)
        LaserScanAddTimeIncrement(b, msg.TimeIncrement)
        LaserScanAddScanTime(b, msg.ScanTime)
        LaserScanAddRangeMin(b, msg.RangeMin)
        LaserScanAddRangeMax(b, msg.RangeMax)
	LaserScanAddRanges(b, rangesArr)
        LaserScanAddIntensities(b, intensitiesArr)
	message_position := LaserScanEnd(b)
	b.Finish(message_position)

	return b.Bytes[b.Head():]
}

func boolToByte(v bool) byte {
	if v {
		return 1
	} else {
		return 0
	}
}
