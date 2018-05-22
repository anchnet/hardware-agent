package main

import (
	"bytes"
	"encoding/json"
	"errors"
	log "github.com/cihub/seelog"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/anchnet/hardware-agent/g"
	"io"
)

type SysInfo struct {
	Endpoint string `json:"endpoint"`
	Brand    string `json:"brand"`
	Model    string `json:"model"`
}

func ReportSysInfo() {
	go func() {
		for {
			if err := reportSysInfo(); err == nil {
				log.Info("report sysinfo success")
				break
			}
			log.Info("report sysinfo fail")
			time.Sleep(time.Minute)
		}
	}()
}

func reportSysInfo() error {
	hostname, err := g.Hostname()
	if err != nil {
		log.Info(err)
		hostname = ""
	}

	cmd := exec.Command("./ipmitool_fru.sh")

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Start()

	err_to, isTimeout := CmdRunWithTimeout(cmd, g.Config().ExecTimeout * time.Second)
	if isTimeout {
		// has be killed
		if err_to == nil {
			log.Info("[INFO] timeout and kill process ipmitool_fru successfully")
		}

		if err_to != nil {
			log.Info("[ERROR] kill process ipmitool_fru occur error:", err_to)
		}

		return err_to
	}

	brand := ""
	sys_model := ""
	// exec successfully
	for {
		buf, err := stdout.ReadString('\n')
		if err != nil && err != io.EOF {
			log.Info("[ERROR] stdout of ipmitool_fru error :", err)
			break
		}
		s := strings.Split(buf, "|")
		//fmt.Println(s)
		if (len(s) > 1) {
			brand = strings.Trim(s[0], " ")
			sys_model = strings.Trim(s[1], " ")
		}
		if err == io.EOF {
			break
		}
	}

	sysinfo := SysInfo{
		Endpoint:  hostname,
		Brand:     brand,
		Model:     sys_model,
	}
	if g.Config().Debug {
		log.Info("sysinfo report: ", sysinfo)
	}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(sysinfo)

	res, err := http.Post(g.Config().SmartAPI, "application/json", b)
	if err != nil {
		return err
	}
	var message struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	json.NewDecoder(res.Body).Decode(&message)

	if message.Status == "error" {
		err = errors.New(message.Message)
	}
	return err
}

func CmdRunWithTimeout(cmd *exec.Cmd, timeout time.Duration) (error, bool) {
	var err error

	//set group id
	//err = syscall.Setpgid(cmd.Process.Pid, cmd.Process.Pid)
	if err != nil {
		log.Info("Setpgid failed, error:", err)
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		log.Infof("timeout, process:%s will be killed", cmd.Path)

		go func() {
			<-done // allow goroutine to exit
		}()

	// cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} is necessary before cmd.Start()
		err = cmd.Process.Kill()
		if err != nil {
			log.Info("kill failed, error:", err)
		}

		return err, true
	case err = <-done:
		return err, false
	}
}