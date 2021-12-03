package endpoints

import (
	"backup/cfg"
	"backup/utils"
	"encoding/json"
	"net/http"
)

// ReloadConfigFile /reload endpoint, re-reads the config file and applies new changes to the cfg.Server array.
func ReloadConfigFile(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	err := cfg.ReadConfig("config/config.json")
	var res utils.Result
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		res = utils.Result{Code: http.StatusInternalServerError, Data: "Unable to reload file! " + err.Error()}
	} else {
		writer.WriteHeader(http.StatusOK)
		res = utils.Result{Code: http.StatusOK, Data: "Config file reloaded!"}
	}
	response, _ := json.Marshal(res)
	_, _ = writer.Write(response)
}
