package main

//go:generate gengo msg power_msgs/BatteryState
import (
	"fmt"
       	"time"
	"power_msgs"
	"github.com/patrickmn/go-cache"
	//log "github.com/sirupsen/logrus"
)
        

func main() {
	c := cache.New(5*time.Minute, 10*time.Minute)
	c.Set("token", 55555, cache.NoExpiration)

	topics := []string{"chatter"}
	ch := make(chan *power_msgs.BatteryState,10)
	
	NewSubscriber(topics,ch)
	for{
		msg := <- ch
		fmt.Println(msg)
	}
}
