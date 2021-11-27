package cfg

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

var _SERVERS []Server

type Server struct {
	Name 		  string `json:"name"`  // Unique name ID.
	SSHRemotePath string `json:"ssh_remote_path"`  // SSH path to download where the world/ dir is located.
	SSHUser       string `json:"ssh_user"`  // SSH user.
	SSHPass       string `json:"ssh_pass"`  // SSH pass.
	LocalPath	  string `json:"local_path"`  // Path where we want the copy to be stored.
	NBackups      uint   `json:"n_backups"`  // Number of maximum backups that can be stored.
}

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

	_SERVERS = servers

	return nil
}

func GetServer(name string) (Server, error) {
	for _, server := range _SERVERS {
		if server.Name == name {
			return server, nil
		}
	}
	return Server{}, errors.New("server not found")
}

func CreateSample(path string) {
	_SERVERS = []Server{
		{
			Name:          "test",
			SSHRemotePath: "1.2.3.4:/home/test/bck/",
			SSHUser:       "user",
			SSHPass:       "pass",
			LocalPath:     "/home/bck/",
			NBackups:      5,
		},
	}
	jsonString, err := json.MarshalIndent(_SERVERS, "", "    ")
	if err != nil {
		return
	}
	err = ioutil.WriteFile(path, jsonString, 0644)  // nolint:gosec
	if err != nil {
		return
	}
}
