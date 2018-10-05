// Copyright 2018 Fetch Robotics, Inc.
// Author(s): Pavan Soundara

package main

import (
	"bytes"
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

// Authenticate struct that fetches token and sets in cache.
type AuthManager struct {
	Store      *cache.Cache
	AuthStatus chan AuthStatus
}

// Connect struct to manage the connection status and errors.
type AuthStatus struct {
	Connected bool
	Err       error
}

// Create new AuthManager that takes cache as parameter.
func NewAuthManager(store *cache.Cache) *AuthManager {
	auth := &AuthManager{
		Store:      store,
		AuthStatus: make(chan AuthStatus),
	}
	return auth
}

// Set fetched token from server in cache.
func (auth *AuthManager) setTokenInCache() {
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
			auth.AuthStatus <- AuthStatus{
				Connected: true,
				Err:       nil,
			}
			logrus.Info("[Auth] Succesfully received token.")
			retries = 0
			time.Sleep(onSuccessWait * time.Minute)
		}
	}

	//Send error message through channel after maximum retries.
	auth.AuthStatus <- AuthStatus{
		Connected: false,
		Err:       err,
	}
}

// Calculate wait time and return duration.
func calcWaitTime(retries int) time.Duration {
	r := float64(retries)
	retries = int(math.Pow(2, r))
	return time.Duration(retries) * time.Second
}

// Load configeration from application.toml file.
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

// Acquire token from auth0 server and return it.
func (auth *AuthManager) getTokenFromServer() (string, error) {
	authEnabled := viper.GetBool("auth.enabled")
	if !authEnabled {
		return "testtoken", nil
	}

	p, err := loadConfig()
	if err != nil {
		return "", fmt.Errorf("Error marshalling config: %+v", err)
	}

	auth_url := viper.GetString("auth.auth_url") + "/oauth/token/"
	req, err := http.NewRequest("POST", auth_url, bytes.NewReader(p))
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
