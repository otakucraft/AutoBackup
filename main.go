package main

import (
	"backup/cfg"
	"backup/routes"
	"backup/rsync"
	"backup/utils"
	"flag"
)

// main starts the rsync scheduler and api router, you need to specify a port, or it will default to 8462.
// EX: ./backup -p 8888
func main() {
	var port = flag.String("p", "8462", "Service port")
	flag.Parse()

	// IO stuff, create everything related to config and give a sample if not exist as helper for the end user.
	if !utils.DirExists("config") {
		if !utils.TouchDir("config") {
			return
		}
	}

	if !utils.FileExists("config/config.json") {
		if !utils.TouchFile("config/config.json") {
			return
		}
		cfg.CreateSample("config/config.json")
	}
	err := cfg.ReadConfig("config/config.json")
	if err != nil {
		return
	}

	// Start the rsync scheduler.
	rsync.RsyncExecutor = rsync.NewExecutor()
	rsync.RsyncExecutor.Start()

	// Start the api router.
	routes.StartRouter(port)
}
