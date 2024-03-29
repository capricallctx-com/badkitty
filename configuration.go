package main

type Config struct {
	IOMode                     string        `hcl:"io_mode"`
	HowlIP                     string        `hcl:"howl"`
	MonitorPort                int           `hcl:"monitor_port"`
	InsecurePort               int           `hcl:"insecure_port" default:"0"`
	ConnectionRetentionSeconds int           `hcl:"connection_retention_seconds"  default:"60"`
	StaticPath                 string        `hcl:"static_path"`
	Service                    ServiceConfig `hcl:"service,block"`
	ParkingMode                bool          `hcl:"parking_mode"`
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

// var t ConfigurationInfoStruct
var config Config
