
// Automatically generated from the message definition "power_msgs/BatteryState.msg"
package power_msgs
import (
    "bytes"
    "encoding/binary"
    "github.com/akio/rosgo/ros"
)


type _MsgBatteryState struct {
    text string
    name string
    md5sum string
}

func (t *_MsgBatteryState) Text() string {
    return t.text
}

func (t *_MsgBatteryState) Name() string {
    return t.name
}

func (t *_MsgBatteryState) MD5Sum() string {
    return t.md5sum
}

func (t *_MsgBatteryState) NewMessage() ros.Message {
    m := new(BatteryState)
	m.Name = ""
	m.ChargeLevel = 0.0
	m.IsCharging = false
	m.RemainingTime = ros.Duration{}
	m.TotalCapacity = 0.0
	m.CurrentCapacity = 0.0
	m.BatteryVoltage = 0.0
	m.SupplyVoltage = 0.0
	m.ChargerVoltage = 0.0
    return m
}

var (
    MsgBatteryState = &_MsgBatteryState {
        `# Name of the battery
string name

# Charge level of battery as percentage of maximum charge
float32 charge_level

# If true, the battery is being charged
bool is_charging

# When charging, this is the remaining time until fully charged.
# When discharging, this is the time until battery is empty.
# Non-zero values are considered valid.
duration remaining_time

# Total capacity of battery
float32 total_capacity

# Current capacity of battery
float32 current_capacity

# Voltage of battery
float32 battery_voltage

# Voltage of the supply breaker
float32 supply_voltage

# Voltage of the charger
float32 charger_voltage

`,
        "power_msgs/BatteryState",
        "ccaf6ecc0ffe8d97f762d3e343f15d67",
    }
)

type BatteryState struct {
	Name string `rosmsg:"name:string"`
	ChargeLevel float32 `rosmsg:"charge_level:float32"`
	IsCharging bool `rosmsg:"is_charging:bool"`
	RemainingTime ros.Duration `rosmsg:"remaining_time:duration"`
	TotalCapacity float32 `rosmsg:"total_capacity:float32"`
	CurrentCapacity float32 `rosmsg:"current_capacity:float32"`
	BatteryVoltage float32 `rosmsg:"battery_voltage:float32"`
	SupplyVoltage float32 `rosmsg:"supply_voltage:float32"`
	ChargerVoltage float32 `rosmsg:"charger_voltage:float32"`
}

func (m *BatteryState) Type() ros.MessageType {
	return MsgBatteryState
}

func (m *BatteryState) Serialize(buf *bytes.Buffer) error {
    var err error = nil
    binary.Write(buf, binary.LittleEndian, uint32(len([]byte(m.Name))))
    buf.Write([]byte(m.Name))
    binary.Write(buf, binary.LittleEndian, m.ChargeLevel)
    binary.Write(buf, binary.LittleEndian, m.IsCharging)
    binary.Write(buf, binary.LittleEndian, m.RemainingTime.Sec)
    binary.Write(buf, binary.LittleEndian, m.RemainingTime.NSec)
    binary.Write(buf, binary.LittleEndian, m.TotalCapacity)
    binary.Write(buf, binary.LittleEndian, m.CurrentCapacity)
    binary.Write(buf, binary.LittleEndian, m.BatteryVoltage)
    binary.Write(buf, binary.LittleEndian, m.SupplyVoltage)
    binary.Write(buf, binary.LittleEndian, m.ChargerVoltage)
    return err
}


func (m *BatteryState) Deserialize(buf *bytes.Reader) error {
    var err error = nil
    {
        var size uint32
        if err = binary.Read(buf, binary.LittleEndian, &size); err != nil {
            return err
        }
        data := make([]byte, int(size))
        if err = binary.Read(buf, binary.LittleEndian, data); err != nil {
            return err
        }
        m.Name = string(data)
    }
    if err = binary.Read(buf, binary.LittleEndian, &m.ChargeLevel); err != nil {
        return err
    }
    if err = binary.Read(buf, binary.LittleEndian, &m.IsCharging); err != nil {
        return err
    }
    {
        if err = binary.Read(buf, binary.LittleEndian, &m.RemainingTime.Sec); err != nil {
            return err
        }

        if err = binary.Read(buf, binary.LittleEndian, &m.RemainingTime.NSec); err != nil {
            return err
        }
    }
    if err = binary.Read(buf, binary.LittleEndian, &m.TotalCapacity); err != nil {
        return err
    }
    if err = binary.Read(buf, binary.LittleEndian, &m.CurrentCapacity); err != nil {
        return err
    }
    if err = binary.Read(buf, binary.LittleEndian, &m.BatteryVoltage); err != nil {
        return err
    }
    if err = binary.Read(buf, binary.LittleEndian, &m.SupplyVoltage); err != nil {
        return err
    }
    if err = binary.Read(buf, binary.LittleEndian, &m.ChargerVoltage); err != nil {
        return err
    }
    return err
}
