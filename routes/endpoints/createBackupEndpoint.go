package endpoints

import (
	"backup/cfg"
	"backup/rsync"
	"backup/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

// CreateBackup /create/{name} endpoint, will check if a server exists in the config by name, and start the threaded
// rsync process.
func CreateBackup(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	name := vars["name"]
	writer.Header().Set("Content-Type", "application/json")
	server, err := cfg.GetServer(name)
	var res utils.Result
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		res = utils.Result{Code: http.StatusInternalServerError, Data: err.Error()}
	} else {
		writer.WriteHeader(http.StatusOK)
		res = utils.Result{Code: http.StatusOK, Data: "Starting backup!"}
		backupInstance := rsync.NewRsyncInstance(server)
		rsync.RsyncExecutor.StartInstance(backupInstance)
	}
	response, _ := json.Marshal(res)
	_, _ = writer.Write(response)
}
