// SPDX-License-Identifier: MIT
/*
   * Bad Kitty is a simple web server that can serve static files and reverse proxy requests to other servers.
   * It is designed to be a simple, easy to use, and easy to configure web server.

   * Contributors can add copyright here (not necessary - but a good idea).

   Copyright (c) 2024 - Caprica LLC
*/
package main

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

const BehaviorExploitScanning = "exploit_scanning"
const BehaviorBruteForce = "brute_force"
const BehaviorDDoS = "ddos"
const BehaviorPortScan = "port_scan"
const BehaviorOther = "other"
const BehaviorNoisy = "noisy"
const StatusActive = "active"
const StatusTimeout = "timeout"
const StatusBanned = "banned"

type BadKitty struct {
	IP          string    `json:"ip"`
	Port        int       `json:"port"`
	BadBehavior string    `json:"bad_behavior"`
	Status      string    `json:"status"`
	ParoleOver  time.Time `json:"parole_over"`
	LastSeen    time.Time `json:"last_seen"`
}

type Monitor struct {
	Version       string     `json:"version"`
	UpTimeSeconds int64      `json:"uptime_seconds"`
	Status        string     `json:"status"`
	Connections   int        `json:"connections"`
	BadKitties    []BadKitty `json:"bad_kitties"`
}

var MonitorData Monitor

func config2Monitor() {
	MonitorData.Version = version
	MonitorData.Status = StatusActive
	MonitorData.Connections = 0
	MonitorData.BadKitties = []BadKitty{}
}

func monitorGet(w http.ResponseWriter, r *http.Request) {
	connectionList.Set([]byte(r.RemoteAddr), []byte("1"), config.ConnectionRetentionSeconds)
	MonitorData.Connections = int(connectionList.EntryCount())
	MonitorData.BadKitties = []BadKitty{}
	MonitorData.UpTimeSeconds = int64(time.Since(startTime).Seconds())
	iter := badKittyList.NewIterator()
	for i := 0; int64(i) < badKittyList.EntryCount(); i++ {
		value := iter.Next()
		var bk BadKitty
		json.Unmarshal(value.Key, &bk)
		MonitorData.BadKitties = append(MonitorData.BadKitties, bk)
	}
	serialized, _ := json.Marshal(MonitorData)
	io.WriteString(w, string(serialized))
}

func monitorPort() {
	logger.Info("Starting Monitor Server ", zap.Int("port", config.MonitorPort))
	http.HandleFunc("/.bad_kitty/heartbeat", monitorGet)
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.MonitorPort), nil)
	if err != nil {
		logger.Fatal("error: ", zap.Error(err))
	}
}
