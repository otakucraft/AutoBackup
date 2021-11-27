package rsync

import (
	"backup/cfg"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Status int

const (
	RUNNING = iota
	STOPPED
)

type Instance struct {
	server  cfg.Server
	status  Status
	process *os.Process
}

func NewRsyncInstance(server cfg.Server) *Instance {
	return &Instance{
		server:  server,
		status:  RUNNING,
		process: nil,
	}
}

func (rs *Instance) Run() {
	sshPass := fmt.Sprintf("sshpass -p %s ssh -l %s", rs.server.SSHPass, rs.server.SSHUser)
	cmd := exec.Command("rsync",
		"-avh",
		"--rsh", sshPass,
		rs.server.SSHRemotePath,
		rs.server.LocalPath,
	)
	log.Println(cmd)
	stdoutIn, _ := cmd.StdoutPipe()
	err := cmd.Start()
	if err != nil {
		log.Println(err)
		return
	}
	rs.process = cmd.Process
	rs.capture(stdoutIn)
	_ = cmd.Wait()
}

func (rs *Instance) Stop() {
	_ = rs.process.Kill()
}

func (rs *Instance) capture(r io.Reader) {
	buf := make([]byte, 1024)
	var out string
	for {
		n, err := r.Read(buf)
		if n > 0 {
			out = string(buf[:n])
			for _, v := range strings.Split(out, "\n") {
				fmt.Println(v)
			}
		}
		if err != nil {
			if err == io.EOF {
				rs.status = STOPPED
			}
			return
		}
	}
}
