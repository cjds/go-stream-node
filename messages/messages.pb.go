// Code generated by protoc-gen-go. DO NOT EDIT.
// source: messages.proto

package messages

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type String struct {
	Data                 string   `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *String) Reset()         { *m = String{} }
func (m *String) String() string { return proto.CompactTextString(m) }
func (*String) ProtoMessage()    {}
func (*String) Descriptor() ([]byte, []int) {
	return fileDescriptor_4dc296cbfe5ffcd5, []int{0}
}

func (m *String) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_String.Unmarshal(m, b)
}
func (m *String) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_String.Marshal(b, m, deterministic)
}
func (m *String) XXX_Merge(src proto.Message) {
	xxx_messageInfo_String.Merge(m, src)
}
func (m *String) XXX_Size() int {
	return xxx_messageInfo_String.Size(m)
}
func (m *String) XXX_DiscardUnknown() {
	xxx_messageInfo_String.DiscardUnknown(m)
}

var xxx_messageInfo_String proto.InternalMessageInfo

func (m *String) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

type BatteryState struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	IsCharging           bool     `protobuf:"varint,2,opt,name=is_charging,json=isCharging,proto3" json:"is_charging,omitempty"`
	TotalCapacity        float32  `protobuf:"fixed32,3,opt,name=total_capacity,json=totalCapacity,proto3" json:"total_capacity,omitempty"`
	CurrentCapacity      float32  `protobuf:"fixed32,4,opt,name=current_capacity,json=currentCapacity,proto3" json:"current_capacity,omitempty"`
	BatteryVoltage       float32  `protobuf:"fixed32,5,opt,name=battery_voltage,json=batteryVoltage,proto3" json:"battery_voltage,omitempty"`
	SupplyVoltage        float32  `protobuf:"fixed32,6,opt,name=supply_voltage,json=supplyVoltage,proto3" json:"supply_voltage,omitempty"`
	ChargerVoltage       float32  `protobuf:"fixed32,7,opt,name=charger_voltage,json=chargerVoltage,proto3" json:"charger_voltage,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BatteryState) Reset()         { *m = BatteryState{} }
func (m *BatteryState) String() string { return proto.CompactTextString(m) }
func (*BatteryState) ProtoMessage()    {}
func (*BatteryState) Descriptor() ([]byte, []int) {
	return fileDescriptor_4dc296cbfe5ffcd5, []int{1}
}

func (m *BatteryState) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BatteryState.Unmarshal(m, b)
}
func (m *BatteryState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BatteryState.Marshal(b, m, deterministic)
}
func (m *BatteryState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BatteryState.Merge(m, src)
}
func (m *BatteryState) XXX_Size() int {
	return xxx_messageInfo_BatteryState.Size(m)
}
func (m *BatteryState) XXX_DiscardUnknown() {
	xxx_messageInfo_BatteryState.DiscardUnknown(m)
}

var xxx_messageInfo_BatteryState proto.InternalMessageInfo

func (m *BatteryState) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *BatteryState) GetIsCharging() bool {
	if m != nil {
		return m.IsCharging
	}
	return false
}

func (m *BatteryState) GetTotalCapacity() float32 {
	if m != nil {
		return m.TotalCapacity
	}
	return 0
}

func (m *BatteryState) GetCurrentCapacity() float32 {
	if m != nil {
		return m.CurrentCapacity
	}
	return 0
}

func (m *BatteryState) GetBatteryVoltage() float32 {
	if m != nil {
		return m.BatteryVoltage
	}
	return 0
}

func (m *BatteryState) GetSupplyVoltage() float32 {
	if m != nil {
		return m.SupplyVoltage
	}
	return 0
}

func (m *BatteryState) GetChargerVoltage() float32 {
	if m != nil {
		return m.ChargerVoltage
	}
	return 0
}

type LaserScan struct {
	AngleMin             float32   `protobuf:"fixed32,1,opt,name=angle_min,json=angleMin,proto3" json:"angle_min,omitempty"`
	AngleMax             float32   `protobuf:"fixed32,2,opt,name=angle_max,json=angleMax,proto3" json:"angle_max,omitempty"`
	AngleIncrement       float32   `protobuf:"fixed32,3,opt,name=angle_increment,json=angleIncrement,proto3" json:"angle_increment,omitempty"`
	TimeIncrement        float32   `protobuf:"fixed32,4,opt,name=time_increment,json=timeIncrement,proto3" json:"time_increment,omitempty"`
	ScanTime             float32   `protobuf:"fixed32,5,opt,name=scan_time,json=scanTime,proto3" json:"scan_time,omitempty"`
	RangeMin             float32   `protobuf:"fixed32,6,opt,name=range_min,json=rangeMin,proto3" json:"range_min,omitempty"`
	RangeMax             float32   `protobuf:"fixed32,7,opt,name=range_max,json=rangeMax,proto3" json:"range_max,omitempty"`
	Ranges               []float32 `protobuf:"fixed32,8,rep,packed,name=ranges,proto3" json:"ranges,omitempty"`
	Intensities          []float32 `protobuf:"fixed32,9,rep,packed,name=intensities,proto3" json:"intensities,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *LaserScan) Reset()         { *m = LaserScan{} }
func (m *LaserScan) String() string { return proto.CompactTextString(m) }
func (*LaserScan) ProtoMessage()    {}
func (*LaserScan) Descriptor() ([]byte, []int) {
	return fileDescriptor_4dc296cbfe5ffcd5, []int{2}
}

func (m *LaserScan) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LaserScan.Unmarshal(m, b)
}
func (m *LaserScan) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LaserScan.Marshal(b, m, deterministic)
}
func (m *LaserScan) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LaserScan.Merge(m, src)
}
func (m *LaserScan) XXX_Size() int {
	return xxx_messageInfo_LaserScan.Size(m)
}
func (m *LaserScan) XXX_DiscardUnknown() {
	xxx_messageInfo_LaserScan.DiscardUnknown(m)
}

var xxx_messageInfo_LaserScan proto.InternalMessageInfo

func (m *LaserScan) GetAngleMin() float32 {
	if m != nil {
		return m.AngleMin
	}
	return 0
}

func (m *LaserScan) GetAngleMax() float32 {
	if m != nil {
		return m.AngleMax
	}
	return 0
}

func (m *LaserScan) GetAngleIncrement() float32 {
	if m != nil {
		return m.AngleIncrement
	}
	return 0
}

func (m *LaserScan) GetTimeIncrement() float32 {
	if m != nil {
		return m.TimeIncrement
	}
	return 0
}

func (m *LaserScan) GetScanTime() float32 {
	if m != nil {
		return m.ScanTime
	}
	return 0
}

func (m *LaserScan) GetRangeMin() float32 {
	if m != nil {
		return m.RangeMin
	}
	return 0
}

func (m *LaserScan) GetRangeMax() float32 {
	if m != nil {
		return m.RangeMax
	}
	return 0
}

func (m *LaserScan) GetRanges() []float32 {
	if m != nil {
		return m.Ranges
	}
	return nil
}

func (m *LaserScan) GetIntensities() []float32 {
	if m != nil {
		return m.Intensities
	}
	return nil
}

type PayLoad struct {
	// Types that are valid to be assigned to Data:
	//	*PayLoad_StringMessage
	//	*PayLoad_BatteryState
	//	*PayLoad_LaserScan
	Data                 isPayLoad_Data `protobuf_oneof:"data"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *PayLoad) Reset()         { *m = PayLoad{} }
func (m *PayLoad) String() string { return proto.CompactTextString(m) }
func (*PayLoad) ProtoMessage()    {}
func (*PayLoad) Descriptor() ([]byte, []int) {
	return fileDescriptor_4dc296cbfe5ffcd5, []int{3}
}

func (m *PayLoad) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PayLoad.Unmarshal(m, b)
}
func (m *PayLoad) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PayLoad.Marshal(b, m, deterministic)
}
func (m *PayLoad) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PayLoad.Merge(m, src)
}
func (m *PayLoad) XXX_Size() int {
	return xxx_messageInfo_PayLoad.Size(m)
}
func (m *PayLoad) XXX_DiscardUnknown() {
	xxx_messageInfo_PayLoad.DiscardUnknown(m)
}

var xxx_messageInfo_PayLoad proto.InternalMessageInfo

type isPayLoad_Data interface {
	isPayLoad_Data()
}

type PayLoad_StringMessage struct {
	StringMessage *String `protobuf:"bytes,2,opt,name=string_message,json=stringMessage,proto3,oneof"`
}

type PayLoad_BatteryState struct {
	BatteryState *BatteryState `protobuf:"bytes,3,opt,name=battery_state,json=batteryState,proto3,oneof"`
}

type PayLoad_LaserScan struct {
	LaserScan *LaserScan `protobuf:"bytes,4,opt,name=laser_scan,json=laserScan,proto3,oneof"`
}

func (*PayLoad_StringMessage) isPayLoad_Data() {}

func (*PayLoad_BatteryState) isPayLoad_Data() {}

func (*PayLoad_LaserScan) isPayLoad_Data() {}

func (m *PayLoad) GetData() isPayLoad_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *PayLoad) GetStringMessage() *String {
	if x, ok := m.GetData().(*PayLoad_StringMessage); ok {
		return x.StringMessage
	}
	return nil
}

func (m *PayLoad) GetBatteryState() *BatteryState {
	if x, ok := m.GetData().(*PayLoad_BatteryState); ok {
		return x.BatteryState
	}
	return nil
}

func (m *PayLoad) GetLaserScan() *LaserScan {
	if x, ok := m.GetData().(*PayLoad_LaserScan); ok {
		return x.LaserScan
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*PayLoad) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _PayLoad_OneofMarshaler, _PayLoad_OneofUnmarshaler, _PayLoad_OneofSizer, []interface{}{
		(*PayLoad_StringMessage)(nil),
		(*PayLoad_BatteryState)(nil),
		(*PayLoad_LaserScan)(nil),
	}
}

func _PayLoad_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*PayLoad)
	// data
	switch x := m.Data.(type) {
	case *PayLoad_StringMessage:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.StringMessage); err != nil {
			return err
		}
	case *PayLoad_BatteryState:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.BatteryState); err != nil {
			return err
		}
	case *PayLoad_LaserScan:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.LaserScan); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("PayLoad.Data has unexpected type %T", x)
	}
	return nil
}

func _PayLoad_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*PayLoad)
	switch tag {
	case 2: // data.string_message
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(String)
		err := b.DecodeMessage(msg)
		m.Data = &PayLoad_StringMessage{msg}
		return true, err
	case 3: // data.battery_state
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(BatteryState)
		err := b.DecodeMessage(msg)
		m.Data = &PayLoad_BatteryState{msg}
		return true, err
	case 4: // data.laser_scan
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(LaserScan)
		err := b.DecodeMessage(msg)
		m.Data = &PayLoad_LaserScan{msg}
		return true, err
	default:
		return false, nil
	}
}

func _PayLoad_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*PayLoad)
	// data
	switch x := m.Data.(type) {
	case *PayLoad_StringMessage:
		s := proto.Size(x.StringMessage)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *PayLoad_BatteryState:
		s := proto.Size(x.BatteryState)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *PayLoad_LaserScan:
		s := proto.Size(x.LaserScan)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*String)(nil), "String")
	proto.RegisterType((*BatteryState)(nil), "BatteryState")
	proto.RegisterType((*LaserScan)(nil), "LaserScan")
	proto.RegisterType((*PayLoad)(nil), "PayLoad")
}

func init() { proto.RegisterFile("messages.proto", fileDescriptor_4dc296cbfe5ffcd5) }

var fileDescriptor_4dc296cbfe5ffcd5 = []byte{
	// 445 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x93, 0xd1, 0x6a, 0xdb, 0x30,
	0x14, 0x86, 0x63, 0x37, 0x73, 0xec, 0xe3, 0xc6, 0x29, 0xba, 0x32, 0xeb, 0x60, 0x21, 0x6c, 0x34,
	0x63, 0x10, 0x46, 0xb6, 0x27, 0x48, 0x6f, 0x32, 0x68, 0x61, 0x38, 0x63, 0xb7, 0xe6, 0xc4, 0x15,
	0x9e, 0xc0, 0x96, 0x83, 0xa4, 0x8e, 0xe4, 0x25, 0xf6, 0x18, 0xbb, 0xdc, 0x33, 0x16, 0x1d, 0xc9,
	0xae, 0xef, 0x94, 0xef, 0x7c, 0x44, 0xd2, 0xaf, 0xdf, 0x90, 0xb5, 0x5c, 0x6b, 0xac, 0xb9, 0xde,
	0x9c, 0x54, 0x67, 0xba, 0xd5, 0x3b, 0x88, 0x0e, 0x46, 0x09, 0x59, 0x33, 0x06, 0xd3, 0x27, 0x34,
	0x98, 0x07, 0xcb, 0x60, 0x9d, 0x14, 0xb4, 0x5e, 0xfd, 0x0d, 0xe1, 0x7a, 0x87, 0xc6, 0x70, 0x75,
	0x39, 0x18, 0x34, 0xdc, 0x4a, 0x12, 0x5b, 0xde, 0x4b, 0x76, 0xcd, 0xde, 0x43, 0x2a, 0x74, 0x59,
	0xfd, 0x46, 0x55, 0x0b, 0x59, 0xe7, 0xe1, 0x32, 0x58, 0xc7, 0x05, 0x08, 0x7d, 0xef, 0x09, 0xfb,
	0x08, 0x99, 0xe9, 0x0c, 0x36, 0x65, 0x85, 0x27, 0xac, 0x84, 0xb9, 0xe4, 0x57, 0xcb, 0x60, 0x1d,
	0x16, 0x73, 0xa2, 0xf7, 0x1e, 0xb2, 0x4f, 0x70, 0x53, 0x3d, 0x2b, 0xc5, 0xa5, 0x79, 0x15, 0xa7,
	0x24, 0x2e, 0x3c, 0x1f, 0xd4, 0x3b, 0x58, 0x1c, 0xdd, 0xb1, 0xca, 0x3f, 0x5d, 0x63, 0xb0, 0xe6,
	0xf9, 0x1b, 0x32, 0x33, 0x8f, 0x7f, 0x39, 0x6a, 0xb7, 0xd6, 0xcf, 0xa7, 0x53, 0xf3, 0xea, 0x45,
	0x6e, 0x6b, 0x47, 0x7b, 0xed, 0x0e, 0x16, 0x74, 0x7e, 0xae, 0x06, 0x6f, 0xe6, 0xfe, 0xcf, 0x63,
	0x2f, 0xae, 0xfe, 0x87, 0x90, 0x3c, 0xa0, 0xe6, 0xea, 0x50, 0xa1, 0x64, 0xb7, 0x90, 0xa0, 0xac,
	0x1b, 0x5e, 0xb6, 0x42, 0x52, 0x24, 0x61, 0x11, 0x13, 0x78, 0x14, 0xe3, 0x21, 0x9e, 0x29, 0x94,
	0x61, 0x88, 0x67, 0xbb, 0xa1, 0x1b, 0x0a, 0x59, 0x29, 0xde, 0x72, 0x69, 0x7c, 0x26, 0x19, 0xe1,
	0xef, 0x3d, 0xa5, 0xec, 0x44, 0x3b, 0xf6, 0xa6, 0x3e, 0x3b, 0xd1, 0x8e, 0xb4, 0x5b, 0x48, 0x74,
	0x85, 0xb2, 0xb4, 0xd4, 0x47, 0x11, 0x5b, 0xf0, 0x53, 0xb4, 0xdc, 0x0e, 0x15, 0xca, 0xda, 0x1d,
	0xd3, 0xdd, 0x3f, 0x26, 0xe0, 0x8f, 0xe9, 0x87, 0x78, 0xf6, 0x97, 0xf6, 0x43, 0x3c, 0xb3, 0xb7,
	0x10, 0xd1, 0x5a, 0xe7, 0xf1, 0xf2, 0x6a, 0x1d, 0xee, 0xc2, 0x9b, 0xa0, 0xf0, 0x84, 0x7d, 0x80,
	0x54, 0x48, 0xc3, 0xa5, 0x16, 0x46, 0x70, 0x9d, 0x27, 0x83, 0x30, 0xc6, 0xab, 0x7f, 0x01, 0xcc,
	0x7e, 0xe0, 0xe5, 0xa1, 0xc3, 0x27, 0xf6, 0x05, 0x32, 0x4d, 0x5d, 0x2b, 0x7d, 0x09, 0x29, 0x96,
	0x74, 0x3b, 0xdb, 0xb8, 0x0a, 0xee, 0x27, 0xc5, 0xdc, 0x09, 0x8f, 0x6e, 0xce, 0xbe, 0xc1, 0xbc,
	0x7f, 0x67, 0x6d, 0xfb, 0x47, 0x21, 0xa5, 0xdb, 0xf9, 0x66, 0x5c, 0xca, 0xfd, 0xa4, 0xb8, 0x3e,
	0x8e, 0x4b, 0xfa, 0x19, 0xa0, 0xb1, 0x6f, 0x54, 0xda, 0x04, 0x28, 0xaf, 0x74, 0x0b, 0x9b, 0xe1,
	0xd9, 0xf6, 0x93, 0x22, 0x69, 0xfa, 0x1f, 0xbb, 0xc8, 0xd5, 0xfe, 0x18, 0xd1, 0xf7, 0xf0, 0xf5,
	0x25, 0x00, 0x00, 0xff, 0xff, 0x67, 0x3e, 0x33, 0xf6, 0x21, 0x03, 0x00, 0x00,
}