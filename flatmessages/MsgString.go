// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package flatmessages

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type MsgString struct {
	_tab flatbuffers.Table
}

func GetRootAsMsgString(buf []byte, offset flatbuffers.UOffsetT) *MsgString {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &MsgString{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *MsgString) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *MsgString) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *MsgString) Data() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func MsgStringStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func MsgStringAddData(builder *flatbuffers.Builder, data flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(data), 0)
}
func MsgStringEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
