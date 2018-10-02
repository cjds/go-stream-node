// Copyright 2018 Fetch Robotics, Inc.
// Author(s): Pavan Soundara

package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

// Authenticate struct that fetches token and sets in cache
type AuthManager struct {
	Store     *cache.Cache
	Connected chan bool
	Token     Token
}

type Token struct {
	T string
}

//Auth0 credentials struct
type AuthPayload struct {
	Username           string
	Password           string
	ClientId           string
	ClientSecret       string
	AuthUrl            string
	Realm              string
	DefaultHost        string
	DefaultPort        int
	UseDeprecatedLogin bool
}

//Create new AuthManager that takes cache as parameter
func NewAuthManager(store *cache.Cache) *AuthManager {
	logrus.Info("[Auth] NewAuthManager called")
	auth := &AuthManager{
		Store:     store,
		Connected: make(chan bool),
	}
	return auth
}

//Set fetched token from server in cache
func (auth *AuthManager) setTokenInCache() {
	for {
		logrus.Info("[Auth] setTokenInCache called")
		t := (*auth).getTokenFromServer()
		(*auth).Store.Set("token", t, cache.DefaultExpiration)
		logrus.Info("[Auth] Recieved token:", t)
		(*auth).Connected <- true
		time.Sleep(10 * time.Minute)
	}
}

//Acquire token from auth0 server and return it.
//(WIP) TODO:
//Acquire credentials from config.
//Better error management depending on response codes
//Seperate struct for auth0 config and token.
func (auth *AuthManager) getTokenFromServer() string {

	p, err := json.Marshal(map[string]string{
		"username":      "freight20",
		"password":      "fr-3538323734355101002e0022",
		"grant_type":    "http://auth0.com/oauth/grant-type/password-realm",
		"scope":         "profile email fetchcore:all_access openid offline_access",
		"audience":      "fetchcore",
		"client_id":     "efPjIzZrsTq2XDZSBAS8kbbEbo2ptQZj",
		"client_secret": "g9CjGEqON0k-xiBFCRl7BGCGdNjgVwFPIVmsCLUm0dykPUvawhhCVhXDYjaKC1F9",
		"realm":         "dev",
	})

	if err != nil {
		logrus.Info("Marshalling error")
	}

	req, err := http.NewRequest("POST", "https://hello-there.auth0.com/oauth/token/", bytes.NewReader(p))

	if err != nil {
		logrus.Info("Error sending a POST request")
	}

	req.Header.Add("Content-Type", "application/json")

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	if err != nil {
		logrus.Info("Request error")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Info("Reading error")
	}

	var res map[string]interface{}
	errUnmarshal := json.Unmarshal(body, &res)
	if errUnmarshal != nil {
		logrus.Info("Response unmarshall error")
	}

	return res["access_token"].(string)
}
