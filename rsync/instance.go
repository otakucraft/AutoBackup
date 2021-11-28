package rsync

import (
	"backup/cfg"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
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

	rs.zipAndRotate()
}

func (rs *Instance) zipAndRotate() {
	path := strings.Split(strings.TrimSuffix(rs.server.LocalPath, "/"), "/")
	zipDir := strings.Join(path[:len(path) - 1], "/")
	files, err := ioutil.ReadDir(zipDir)
	if err != nil {
		return
	}
	var zipFileList []string
	for _, file := range files {
		if hasValidFormat, _ := regexp.MatchString(fmt.Sprintf("^%s_\\d{4}-\\d{1,2}-\\d{1,2}.zip", rs.server.Name), file.Name()); hasValidFormat && !file.IsDir() {
			zipFileList = append(zipFileList, file.Name())
		}
	}

	if len(zipFileList) >= rs.server.NBackups {
		for _, oldZipFile := range zipFileList[:len(zipFileList) - (rs.server.NBackups - 1)] {
			_ = os.Remove(filepath.Join(zipDir, oldZipFile))
		}
	}

	target := fmt.Sprintf("%s_%s.zip", rs.server.Name, time.Now().Format("2006-01-02"))
	_, _ = exec.Command("tar", "-zcf", target, "-C", rs.server.LocalPath, ".").Output()
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
