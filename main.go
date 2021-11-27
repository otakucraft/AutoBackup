package main

import (
	"backup/cfg"
	"backup/routes"
	"backup/rsync"
	"backup/utils"
)

func main() {
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

	rsync.RsyncExecutor = rsync.NewExecutor()
	rsync.RsyncExecutor.Start()

	routes.StartRouter()
}
