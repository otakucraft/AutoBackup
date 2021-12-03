package cfg

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

var Servers []Server

type Server struct {
	Name 		  string `json:"name"`  // Unique name ID.
	SSHRemotePath string `json:"ssh_remote_path"`  // SSH path to download where the world/ dir is located.
	SSHUser       string `json:"ssh_user"`  // SSH user.
	SSHPass       string `json:"ssh_pass"`  // SSH pass.
	LocalPath	  string `json:"local_path"`  // Path where we want the copy to be stored.
	NBackups      int    `json:"n_backups"`  // Number of maximum backups that can be stored.
}

// ReadConfig Reads and parses the config json file into an array of Server.
func ReadConfig(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var servers []Server
	err = json.Unmarshal(data, &servers)
	if err != nil {
		return err
	}

	Servers = servers

	return nil
}

// GetServer iterates the Server array to get the credentials of the server based on endpoint name.
func GetServer(name string) (Server, error) {
	for _, server := range Servers {
		if server.Name == name {
			return server, nil
		}
	}
	return Server{}, errors.New("server not found")
}

// CreateSample if the config json file does not exist, it will create a default one.
func CreateSample(path string) {
	Servers = []Server{
		{
			Name:          "test",
			SSHRemotePath: "1.2.3.4:/home/test/bck/",
			SSHUser:       "user",
			SSHPass:       "pass",
			LocalPath:     "/home/bck/",
			NBackups:      5,
		},
	}
	jsonString, err := json.MarshalIndent(Servers, "", "    ")
	if err != nil {
		return
	}
	err = ioutil.WriteFile(path, jsonString, 0644)  // nolint:gosec
	if err != nil {
		return
	}
}
