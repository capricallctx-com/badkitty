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
	"github.com/hashicorp/hcl/v2/hclsimple"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var data []byte
var logger *zap.Logger

type Config struct {
	IOMode  string        `hcl:"io_mode"`
	Service ServiceConfig `hcl:"service,block"`
}

type ServiceConfig struct {
	Protocol   string          `hcl:"protocol,label"`
	Type       string          `hcl:"type,label"`
	ListenAddr string          `hcl:"listen_addr"`
	Processes  []ProcessConfig `hcl:"route,block"`
}

type ProcessConfig struct {
	Type   string `hcl:"type,label"`
	Target string `hcl:"target"`
}

type ConfigurationInfoStruct struct {
	Global struct {
		LoggingLevel    string `yaml:"logging_level"`
		Verbose         int    `yaml:"verbose"`
		PrivateKey      string `yaml:"private_key"`
		Certificate     string `yaml:"certificate"`
		RandomizeServer bool   `yaml:"randomize_server"`
		EnableInsecure  bool   `yaml:"enable_insecure"`
		InsecurePort    int    `yaml:"insecure_port"`
		EnableTLS       bool   `yaml:"enable_tls"`
		TLSPort         int    `yaml:"tls_port"`
		StaticFiles     string `yaml:"static_files"`
	} `yaml:"global"`

	Static struct {
		Enable               bool   `yaml:"enable"`
		Path                 string `yaml:"path"`
		Index                string `yaml:"index"`
		Error                string `yaml:"error"`
		Cors                 bool   `yaml:"cors"`
		CorsAllow            string `yaml:"cors_allow"`
		CorsMaxAge           int    `yaml:"cors_max_age"`
		CorsAllowHeaders     string `yaml:"cors_allow_headers"`
		CorsAllowMethods     string `yaml:"cors_allow_methods"`
		CorsExposeHeaders    string `yaml:"cors_expose_headers"`
		CorsAllowCredentials bool   `yaml:"cors_allow_credentials"`
		CorsDebug            bool   `yaml:"cors_debug"`
	} `yaml:"static"`
}

var t ConfigurationInfoStruct
var config Config

func PrintConfig() {
	log.Println("Logging Level:", t.Global.LoggingLevel)
	log.Println("Verbose:", t.Global.Verbose)
	log.Println("Private Key:", t.Global.PrivateKey)
	log.Println("Certificate:", t.Global.Certificate)
	log.Println("Randomize Server:", t.Global.RandomizeServer)
	log.Println("Enable Insecure:", t.Global.EnableInsecure)
	log.Println("Insecure Port:", t.Global.InsecurePort)
	log.Println("Enable TLS:", t.Global.EnableTLS)
	log.Println("TLS Port:", t.Global.TLSPort)
}

func main() {
	AmIRoot()
	logger = zap.Must(zap.NewProduction())
	defer logger.Sync()

	err := hclsimple.DecodeFile("badkitty.hcl", nil, &config)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}
	log.Printf("Configuration is %#v", config)
	data, err := os.ReadFile("badkitty.yml")
	if err != nil {
		logger.Fatal("error: ", zap.Error(err))
	}
	err = yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		logger.Fatal("error: ", zap.Error(err))
	}
	PrintConfig()
	go serverInsecure()
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Bad Kitty is running...")
	<-done // Will block here until user hits ctrl+c
}

func serverInsecure() {
	if t.Global.EnableInsecure {
		logger.Warn("WARNING: Insecure mode is enabled. This can be a security risk.")
	}
	logger.Info("Starting Insecure Server ", zap.Int("port", t.Global.InsecurePort))
	for _, route := range config.Service.Processes {
		fmt.Println(route.Target, route.Type)
	}
	if IsNotEmpty(t.Static.Path) {
		// makes no sense to have a static file server that's not /
		http.Handle("/", http.FileServer(http.Dir(t.Static.Path)))
	}
	err := http.ListenAndServe(fmt.Sprintf(":%d", t.Global.InsecurePort), nil)
	if err != nil {
		logger.Fatal("error: ", zap.Error(err))
	}

}
