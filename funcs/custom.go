package funcs

import (
	"github.com/open-falcon/common/model"
	"log"
	"github.com/51idc/custom-agent/g"
	"os/exec"
	"fmt"
	"strings"
	"time"
	"io"
	"strconv"
	"bytes"
)

func CustomMetrics() (L []*model.MetricValue) {
	path_list := g.Config().FilePath
	for _, fpath := range path_list {
		L = path_file_exec(fpath, L)
	}
	return L
}

func path_file_exec(fpath string, L []*model.MetricValue) ([]*model.MetricValue) {
	cmd := exec.Command(fpath)
	if (strings.Contains(fpath, " ")) {
		sep_index := strings.Index(fpath, " ")
		arg1 := fpath[0:sep_index]
		arg2 := fpath[sep_index + 1:len([]rune(fpath))]
		log.Println("[INFO] multi args , exec :", arg1, arg2)
		cmd = exec.Command(arg1, arg2)
	}

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Start()

	err_to, isTimeout := CmdRunWithTimeout(cmd, g.Config().ExecTimeout * time.Millisecond)
	if isTimeout {
		// has be killed
		if err_to == nil {
			log.Println("[INFO] timeout and kill process", fpath, "successfully")
		}

		if err_to != nil {
			log.Println("[ERROR] kill process", fpath, "occur error:", err_to)
		}

		return L
	}

	// exec successfully
	for {
		buf, err := stdout.ReadString('\n')
		if err != nil && err != io.EOF {
			log.Println("[ERROR] stdout of", fpath, "error :", err)
			break
		}
		s := strings.Split(buf, " ")
		if (len(s) > 1) {
			tag := s[0]
			value := s[1]
			value = strings.Replace(value, "\n", "", -1)
			value = strings.Replace(value, "\r", "", -1)
			tags := fmt.Sprintf("name=%s", tag)
			if val, err := strconv.ParseFloat(value, 64); err == nil {
				L = append(L, GaugeValue("custom.data", val, tags))
			} else {
				log.Println("[ERROR] value parse float error , the value is ", value)
				log.Println("err : ", err.Error())
			}
		}
		if err == io.EOF {
			break
		}
	}

	return L
}

func CmdRunWithTimeout(cmd *exec.Cmd, timeout time.Duration) (error, bool) {
	var err error

	//set group id
	//err = syscall.Setpgid(cmd.Process.Pid, cmd.Process.Pid)
	if err != nil {
		log.Println("Setpgid failed, error:", err)
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		log.Printf("timeout, process:%s will be killed", cmd.Path)

		go func() {
			<-done // allow goroutine to exit
		}()

	// cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} is necessary before cmd.Start()
		err = cmd.Process.Kill()
		if err != nil {
			log.Println("kill failed, error:", err)
		}

		return err, true
	case err = <-done:
		return err, false
	}
}