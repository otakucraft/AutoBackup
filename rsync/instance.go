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

// Status enum to know the instance status.
type Status int

const (
	RUNNING = iota
	STOPPED
)

// Instance structure of an instance.
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

// Run start the rsync process, zipping and zip rotation
func (rs *Instance) Run() {
	// WARN needs rsync, sshpass and tar previously installed.
	sshPass := fmt.Sprintf("sshpass -p %s ssh -l %s", rs.server.SSHPass, rs.server.SSHUser)
	cmd := exec.Command("rsync",
		"-avh",
		"--rsh", sshPass,
		rs.server.SSHRemotePath,
		rs.server.LocalPath,
	)
	stdoutIn, _ := cmd.StdoutPipe()
	err := cmd.Start()
	if err != nil {
		log.Println(err)
		return
	}
	rs.process = cmd.Process
	rs.capture(stdoutIn)  // Capture the stdout of the process and log it real time.
	_ = cmd.Wait()  // Wait for the process to finish.

	rs.zipAndRotate()
}

// zipAndRotate will zip the current backup and delete the oldest one based on how many there are as specified in
// the config file. The zip rotation is based on the file name, if you rename the zip file, it will be ignored by
// the rotation.
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

// Stop kill the process.
func (rs *Instance) Stop() {
	_ = rs.process.Kill()
}

// capture gets the stdout of the process and detects when it has finished.
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
