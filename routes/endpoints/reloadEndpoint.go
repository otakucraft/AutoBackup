package endpoints

import (
	"backup/cfg"
	"backup/utils"
	"encoding/json"
	"net/http"
)

func ReloadConfigFile(writer http.ResponseWriter, request *http.Request) {
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
