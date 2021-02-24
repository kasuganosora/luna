package flags

import (
	"flag"
	"github.com/kabukky/journey/logger"
)

var (
	Log         = ""
	CustomPath  = ""
	IsInDevMode = false
	HttpPort    = ""
	HttpsPort   = ""

	Settings map[string]interface{}
)

func init() {
	Settings = make(map[string]interface{})
	// Parse all flags
	parseFlags()
	if IsInDevMode {
		logger.Info("Starting Journey in developer mode...")
	}
}

func parseFlags() {
	// Check if the log should be output to a file
	flag.StringVar(&Log, "log", "", "Use this option to save to log output to a file. Note: Journey needs create, read, and write access to that file. Example: -log=path/to/log.txt")
	if Log != "" {
		Settings["log"] = Log
	}

	// Check if a custom content path has been provided by the user
	flag.StringVar(&CustomPath, "custom-path", "", "Specify a custom path to store content files. Note: Journey needs read and write access to that path. A theme folder needs to be located in the custon path under content/themes. Example: -custom-path=/absolute/path/to/custom/folder")
	if CustomPath != "" {
		Settings["custom_path"] = CustomPath
	}

	// Check if the dvelopment mode flag was provided by the user
	flag.BoolVar(&IsInDevMode, "dev", false, "Use this flag flag to put Journey in developer mode. Features of developer mode: Themes and plugins will be recompiled immediately after changes to the files. Example: -dev")
	Settings["is_in_dev_mode"] = IsInDevMode

	// Check if the http port that was set in the config was overridden by the user
	flag.StringVar(&HttpPort, "http-port", "", "Use this option to override the HTTP port that was set in the config.json. Example: -http-port=8080")
	if HttpPort != "" {
		Settings["http_port"] = HttpPort
	}
	// Check if the http port that was set in the config was overridden by the user
	flag.StringVar(&HttpsPort, "https-port", "", "Use this option to override the HTTPS port that was set in the config.json. Example: -https-port=8081")
	if HttpsPort != "" {
		Settings["https_port"] = HttpsPort
	}
	flag.Parse()
}
