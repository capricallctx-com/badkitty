// SPDX-License-Identifier: MIT
/*
   * Bad Kitty is a simple web server that can serve static files and reverse proxy requests to other servers.
   * It is designed to be a simple, easy to use, and easy to configure web server.

   * Contributors can add copyright here (not necessary - but a good idea).

   Copyright (c) 2024 - Caprica LLC
*/

package main

import (
	"fmt"
	"github.com/coocood/freecache"
	"github.com/gorilla/websocket"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"go.uber.org/zap"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var data []byte
var logger *zap.Logger
var badKittyList *freecache.Cache
var connectionList *freecache.Cache
var startTime time.Time

// main entry point
func main() {
	AmIRoot()        // need root powers - this is made for container use
	config2Monitor() // export configuration back
	// allocate TTL caches
	badKittyList = freecache.NewCache(1024 * 1024)
	connectionList = freecache.NewCache(1024 * 1024)
	startTime = time.Now()

	logger = zap.Must(zap.NewProduction())
	defer logger.Sync()

	err := hclsimple.DecodeFile("badkitty.hcl", nil, &config)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}
	log.Printf("Configuration is %#v", config)
	if config.MonitorPort != 0 {
		logger.Info("Starting Monitor Server ", zap.Int("port", config.MonitorPort))
		go monitorPort()
	}
	if config.HowlIP != "" {
		logger.Info("howl starting...")
		go howlLoop()
	}
	if config.ParkingMode {
		logger.Info("Parking Mode is enabled - every route will return the same response")
	}
	go serverInsecure()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Bad Kitty is running...")
	<-done // Will block here until user hits ctrl+c
}

func howlLoop() {
	url := "ws://localhost:8989" // Your WebSocket endpoint
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket server:", err)
	}
	defer conn.Close()

	// Send a message (you'll need something listening on the server side)
	err = conn.WriteMessage(websocket.TextMessage, []byte("Hello from Go client!"))
	if err != nil {
		log.Println("Error sending message:", err)
		return
	}

	// Receive messages in a loop
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error receiving message:", err)
			break
		}
		fmt.Printf("Received from server: %s\n", message)

		// Send another message after a short delay
		time.Sleep(2 * time.Second)
		err = conn.WriteMessage(websocket.TextMessage, []byte("Another message!"))
		if err != nil {
			log.Println("Error sending message:", err)
			return
		}
	}
}

func serverInsecure() {
	if config.InsecurePort == 0 {
		logger.Info("Insecure server is disabled")
		return
	}
	logger.Warn("WARNING: Insecure mode is enabled. This can be a security risk.")
	logger.Info("Starting Insecure Server ", zap.Int("port", config.InsecurePort))
	if !config.ParkingMode {
		for _, route := range config.Service.Processes {
			fmt.Println(route.Target, route.Type)
		}
		if IsNotEmpty(config.StaticPath) {
			// makes no sense to have a static file server that's not /
			http.Handle("/", http.FileServer(http.Dir(config.StaticPath)))
		}

	} else {
		http.HandleFunc("/", parkingHandler)
	}
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.InsecurePort), nil)
	if err != nil {
		logger.Fatal("error: ", zap.Error(err))
	}

}

func parkingHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new template.
	tmpl := template.New("index.html")

	// Parse the template.
	tmpl, err := tmpl.Parse(`
        <h1>Hello, {{.Version}}!</h1>
    `)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, MonitorData)
	if err != nil {
		panic(err)
	}
	//w.WriteHeader(http.StatusOK)
	//w.Write([]byte("Parking Mode"))
}
