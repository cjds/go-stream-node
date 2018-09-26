package main

//go:generate gengo msg power_msgs/BatteryState
import (
        "fmt"
        "github.com/akio/rosgo/ros"
        "os"
        "power_msgs"
)
        
func callback(msg *power_msgs.BatteryState) {
        fmt.Println("Name:", msg.Name)
        fmt.Println("Charge_level", msg.ChargeLevel)
        fmt.Println("is_charging", msg.IsCharging)
        fmt.Println("remaining_time", msg.RemainingTime)
}

func main() {
        node, err := ros.NewNode("/listener", os.Args)
        if err != nil {
                fmt.Println(err)
                os.Exit(-1)
        }
        defer node.Shutdown()
        node.Logger().SetSeverity(ros.LogLevelDebug)
        node.NewSubscriber("/chatter", power_msgs.MsgBatteryState, callback)
        node.Spin()
}
