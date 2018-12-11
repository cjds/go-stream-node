// Copyright 2018 Fetch Robotics, Inc.
// Author(s): Pavan Soundara

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// AuthManager struct that fetches token and sets in cache.
type AuthManager struct {
	Store   *cache.Cache
	Connect chan bool
}

// NewAuthManager creates and returns new AuthManager.
func NewAuthManager() (*AuthManager, context.Context) {
	auth := &AuthManager{
		Store:   cache.New(1*time.Hour, 2*time.Hour),
		Connect: make(chan bool),
	}
	ctx, cancel := context.WithCancel(context.Background())
	go auth.setTokenInCache(cancel)

	return auth, ctx
}

// setTokenInCache sets fetched token from server in cache.
func (auth *AuthManager) setTokenInCache(cancel context.CancelFunc) {
	retries := 0
	maxRetries := viper.GetInt("auth.max_retries")
	onSuccessWait := viper.GetDuration("auth.on_success_wait")

	var err error
	var t string

	for retries < maxRetries {
		t, err = (*auth).getTokenFromServer()
		if err != nil {
			retries++
			waitTime := calcWaitTime(retries)
			logrus.Error(err)
			logrus.Infof("Waiting for %s before retrying.", waitTime)
			time.Sleep(waitTime)
		} else {
			(*auth).Store.Set("token", t, cache.DefaultExpiration)
			retries = 0
			auth.Connect <- true
			time.Sleep(onSuccessWait * time.Second)
		}
	}

	cancel()
}

// calcWaitTime calculates wait time and returns duration.
func calcWaitTime(retries int) time.Duration {
	r := float64(retries)
	retries = int(math.Pow(2, r))

	return time.Duration(retries) * time.Second
}

// loadConfig loads configuration from application.toml file and returns byte data.
func loadConfig() ([]byte, error) {
	p, err := json.Marshal(map[string]string{
		"username":      viper.GetString("auth.username"),
		"password":      viper.GetString("auth.password"),
		"grant_type":    viper.GetString("auth.grant_type"),
		"scope":         viper.GetString("auth.scope"),
		"audience":      viper.GetString("auth.audience"),
		"client_id":     viper.GetString("auth.client_id"),
		"client_secret": viper.GetString("auth.client_secret"),
		"realm":         viper.GetString("auth.realm"),
	})
	return p, err
}

// getTokenFromServer acquires token from auth0 server and returns it.
func (auth *AuthManager) getTokenFromServer() (string, error) {
	authEnabled := viper.GetBool("auth.enabled")
	if !authEnabled {
		return "testtoken", nil
	}

	p, err := loadConfig()
	if err != nil {
		return "", fmt.Errorf("Error marshalling config: %+v", err)
	}

	authURL := viper.GetString("auth.auth_url") + "/oauth/token/"
	req, err := http.NewRequest("POST", authURL, bytes.NewReader(p))
	if err != nil {
		return "", fmt.Errorf("Error creating new request: %+v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error sending request: %+v", err)
	}

	req.Close = true

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Server responded with an error: %+v", resp.Status)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading body: %+v", err)
	}

	var res map[string]interface{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", fmt.Errorf("Error unmarshalling body: %+v", err)
	}

	return res["access_token"].(string), nil
}

// checkToken checks for token in cache and returns it.
func (auth *AuthManager) checkToken() (string, error) {
	t, found := (*auth).Store.Get("token")
	if !found {
		return "", fmt.Errorf("[Subscribe] Token not found in cache")
	}

	return t.(string), nil
}
