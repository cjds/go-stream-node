package main

import (
	"crypto/rand"
	"fmt"
	"time"
	
	"github.com/sirupsen/logrus"
	"github.com/patrickmn/go-cache"
)

// Authenticate struct that fetches token and sets in cache
type Authenticate struct {
	Store *cache.Cache
	Connected chan bool
}

//Create new authenticator that takes cache as parameter
func NewAuthenticator(store *cache.Cache) *Authenticate {
	logrus.Info("[Auth] NewAuthenticator called")
	auth := &Authenticate{
		Store:		store,
		Connected:	make(chan bool),	
	}
	return auth
}

//Set fetched token from server in cache
func (auth *Authenticate) setTokenInCache(){
	for {
		logrus.Info("[Auth] setTokenInCache called")
		t := (*auth).getTokenFromServer()
		(*auth).Store.Set("token", t, cache.DefaultExpiration)
		(*auth).Connected <- true
		time.Sleep(10 * time.Second)
	}
}

// Acquire token from auth0 server and return it
func (auth *Authenticate) getTokenFromServer() string{
	logrus.Info("[Auth] getTokenFromServer called")
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

