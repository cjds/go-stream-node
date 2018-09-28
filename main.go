package main

//go:generate gengo msg power_msgs/BatteryState
//go:generate gengo msg std_msgs/String
import (
       	"time"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)
        
func main() {
	c := cache.New(1*time.Hour, 2*time.Hour)

	a:= NewAuthenticator(c)
	go a.setTokenInCache()

	t := <- a.Connected
	if !t {
		logrus.Info("[Main] Authentication unsuccessful")
		return
	}

 	s:= NewSubscriber(c)
	s.addTopic("chatter")
	s.printAllTopics()
	go s.newListener()	

	for t := range a.Connected {
		//TODO: Replace generic error info to more meaninful 
		logrus.Info("[Main] Refresh token status")
		if t {
			logrus.Info("[Main] Fetched token successfully")
		} else {
			logrus.Info("[Main] Fetching token was unsuccessful")
			return
		}
	}
}

