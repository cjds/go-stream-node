// Copyright 2018 Fetch Robotics, Inc.
// Author(s): Pavan Soundara

package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"go-stream-node/flatmessages"
	"go-stream-node/messages"
	"math/rand"
	// "power_msgs"
	"sensor_msgs"
	"testing"
	"time"

	"github.com/akio/rosgo/ros"
	"github.com/golang/protobuf/proto"
	fb "github.com/google/flatbuffers/go"
)

func benchmarkJson(b *testing.B, c int) {
	b.StopTimer()
	// var bd []*power_msgs.BatteryState
	var ld []*sensor_msgs.LaserScan

	for i := 0; i < 500; i++ {
		// bd = append(bd, generateBatteryState())
		ld = append(ld, generateLaserScan())
	}

	b.ReportAllocs()
	b.StartTimer()
	rand.Seed(time.Now().UnixNano())

	for j := 0; j < b.N; j++ {
		for k := 0; k < c; k++ {
			// benchmarkJsonMarshal(b, bd[rand.Intn(len(bd))])
			benchmarkJsonMarshal(b, ld[rand.Intn(len(ld))])
		}
	}
}

func benchmarkJsonMarshal(b *testing.B, msg ros.Message) {
	pl := makeJsonPayload(b, msg)
	// pl.Customer = "NoCustomer"
	// pl.ProducerID = "freight20"

	message := Message{}
	// message.Type = "publish"
	message.Payload = pl

	_, _ = json.Marshal(message)
}

func BenchmarkJson_1(b *testing.B)      { benchmarkJson(b, 1) }
func BenchmarkJson_10(b *testing.B)     { benchmarkJson(b, 10) }
func BenchmarkJson_100(b *testing.B)    { benchmarkJson(b, 100) }
func BenchmarkJson_1000(b *testing.B)   { benchmarkJson(b, 1000) }
func BenchmarkJson_10000(b *testing.B)  { benchmarkJson(b, 10000) }
func BenchmarkJson_100000(b *testing.B) { benchmarkJson(b, 100000) }

func benchmarkGob(b *testing.B, c int) {
	b.StopTimer()
	// var bd []*power_msgs.BatteryState
	var ld []*sensor_msgs.LaserScan

	for i := 0; i < 500; i++ {
		// bd = append(bd, generateBatteryState())
		ld = append(ld, generateLaserScan())
	}

	b.ReportAllocs()
	b.StartTimer()
	rand.Seed(time.Now().UnixNano())

	gob.RegisterName("messages.PayLoad_LaserScan", messages.PayLoad_LaserScan{})
	gob.RegisterName("messages.PayLoad_BatteryState", messages.PayLoad_BatteryState{})

	var d bytes.Buffer
	enc := gob.NewEncoder(&d)

	for j := 0; j < b.N; j++ {
		for k := 0; k < c; k++ {
			benchmarkGobMarshal(b, ld[rand.Intn(len(ld))], enc, &d)
		}
	}
}

func benchmarkGobMarshal(b *testing.B, msg ros.Message, g *gob.Encoder, d *bytes.Buffer) {
	d.Reset()
	pl := makeProtoPayload(b, msg)
	err := g.Encode(pl)
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkGob_1(b *testing.B)      { benchmarkGob(b, 1) }
func BenchmarkGob_10(b *testing.B)     { benchmarkGob(b, 10) }
func BenchmarkGob_100(b *testing.B)    { benchmarkGob(b, 100) }
func BenchmarkGob_1000(b *testing.B)   { benchmarkGob(b, 1000) }
func BenchmarkGob_10000(b *testing.B)  { benchmarkGob(b, 10000) }
func BenchmarkGob_100000(b *testing.B) { benchmarkGob(b, 100000) }

func benchmarkFlatbuf(b *testing.B, c int) {
	b.StopTimer()
	// var bd []*power_msgs.BatteryState
	var ld []*sensor_msgs.LaserScan

	for i := 0; i < 500; i++ {
		// bd = append(bd, generateBatteryState())
		ld = append(ld, generateLaserScan())
	}

	b.ReportAllocs()
	b.StartTimer()
	rand.Seed(time.Now().UnixNano())
	f := fb.NewBuilder(0)

	for j := 0; j < b.N; j++ {
		for k := 0; k < c; k++ {
			// benchmarkJsonMarshal(b, bd[rand.Intn(len(bd))])
			benchmarkFlatbufMarshal(b, ld[rand.Intn(len(ld))], f)
		}
	}
}

func benchmarkFlatbufMarshal(b *testing.B, msg ros.Message, f *fb.Builder) {
	_ = flatmessages.MakeMsg(f, msg)
}

func BenchmarkFlatbuf_1(b *testing.B)      { benchmarkFlatbuf(b, 1) }
func BenchmarkFlatbuf_10(b *testing.B)     { benchmarkFlatbuf(b, 10) }
func BenchmarkFlatbuf_100(b *testing.B)    { benchmarkFlatbuf(b, 100) }
func BenchmarkFlatbuf_1000(b *testing.B)   { benchmarkFlatbuf(b, 1000) }
func BenchmarkFlatbuf_10000(b *testing.B)  { benchmarkFlatbuf(b, 10000) }
func BenchmarkFlatbuf_100000(b *testing.B) { benchmarkFlatbuf(b, 100000) }

func benchmarkProtobuf(b *testing.B, c int) {
	b.StopTimer()
	// var bd []*power_msgs.BatteryState
	var ld []*sensor_msgs.LaserScan

	for i := 0; i < 500; i++ {
		// bd = append(bd, generateBatteryState())
		ld = append(ld, generateLaserScan())
	}

	b.ReportAllocs()
	b.StartTimer()
	rand.Seed(time.Now().UnixNano())

	for j := 0; j < b.N; j++ {
		for k := 0; k < c; k++ {
			// benchmarkJsonMarshal(b, bd[rand.Intn(len(bd))])
			benchmarkProtobufMarshal(b, ld[rand.Intn(len(ld))])
		}
	}
}

func benchmarkProtobufMarshal(b *testing.B, msg ros.Message) {
	pl := makeProtoPayload(b, msg)
	_, _ = proto.Marshal(pl)
}

func BenchmarkProtobuf_1(b *testing.B)      { benchmarkProtobuf(b, 1) }
func BenchmarkProtobuf_10(b *testing.B)     { benchmarkProtobuf(b, 10) }
func BenchmarkProtobuf_100(b *testing.B)    { benchmarkProtobuf(b, 100) }
func BenchmarkProtobuf_1000(b *testing.B)   { benchmarkProtobuf(b, 1000) }
func BenchmarkProtobuf_10000(b *testing.B)  { benchmarkProtobuf(b, 10000) }
func BenchmarkProtobuf_100000(b *testing.B) { benchmarkProtobuf(b, 100000) }
