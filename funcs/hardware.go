package funcs

import (
	"github.com/open-falcon/common/model"
	log "github.com/cihub/seelog"
	//"github.com/anchnet/hardware-agent/g"
	"os/exec"
	"fmt"
	"strings"
	"time"
	"io"
	//"strconv"
	"bytes"
	"strconv"
)

func HardwareMetrics() (L []*model.MetricValue) {
	log.Info("[INFO] start ipmitool at : ", time.Now().Format("15:04"))
	L = path_file_exec("./ipmitool.sh", L)
	log.Info("[INFO] end ipmitool at : ", time.Now().Format("15:04"))
	return L
}

func path_file_exec(fpath string, L []*model.MetricValue) ([]*model.MetricValue) {
	cmd := exec.Command(fpath)
	if (strings.Contains(fpath, " ")) {
		sep_index := strings.Index(fpath, " ")
		arg1 := fpath[0:sep_index]
		arg2 := fpath[sep_index + 1:len([]rune(fpath))]
		log.Info("[INFO] multi args , exec :", arg1, arg2)
		cmd = exec.Command(arg1, arg2)
	}

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Start()

	err_to, isTimeout := CmdRunWithTimeout(cmd, 10 * time.Minute)
	if isTimeout {
		// has be killed
		if err_to == nil {
			log.Info("[INFO] timeout and kill process", fpath, "successfully")
		}

		if err_to != nil {
			log.Info("[ERROR] kill process", fpath, "occur error:", err_to)
		}

		return L
	}

	// exec successfully
	for {
		buf, err := stdout.ReadString('\n')
		if err != nil && err != io.EOF {
			log.Info("[ERROR] stdout of", fpath, "error :", err)
			break
		}
		s := strings.Split(buf, "|")
		//fmt.Println(s)
		if (len(s) > 1) {
			//deal with s[3] , get Value
			if strings.Contains(s[3], "No Reading") || strings.Contains(s[3], "Disabled") {
				//如果值为No Reading或Disabled，则丢弃数据
				log.Info("[INFO] Drop Data : ", s[3])
				continue
			}
			value_arr := strings.Split(strings.Trim(s[3], " "), " ")
			value := strings.Replace(value_arr[0], "h\n", "", -1)
			//log.Info("[INFO] Value : ", value)

			// deal with s[0] , get Entity Number
			entity_arr := strings.Split(s[0], "(")
			entity_value := strings.TrimSpace(entity_arr[0])
			//log.Info("[INFO] Entity_ID : ", entity_value)

			entity_name_arr := strings.Split(entity_arr[1], ")")
			entity_name := strings.TrimSpace(entity_name_arr[0])
			entity_name = strings.Replace(entity_name, " ", "_", -1)
			//log.Info("[INFO] Entity_Name : ", entity_name)

			//deal with s[1] , get Metric
			metric_arr := strings.Split(s[1], "(")
			metric := strings.Replace(strings.Trim(metric_arr[0], " "), " ", "_", -1)
			metric = strings.Replace(metric, "_/_", "_", -1)
			//log.Info("[INFO] Metric : ", metric)

			//deal with s[2] , get Type
			type_arr := strings.Split(s[2], "(")
			sensor_type := strings.Replace(strings.Trim(type_arr[0], " "), " ", "_", -1)
			sensor_type = strings.Replace(sensor_type, "_/_", "_", -1)
			//log.Info("[INFO] Sensor_Type : ", sensor_type)

			tags := fmt.Sprintf("entity_name=%s,entity_id=%s,sensor_type=%s", entity_name, entity_value, sensor_type)
			if val, err := strconv.ParseFloat(value, 64); err == nil {
				L = append(L, GaugeValue(metric, val, tags))
			} else {
				log.Info("[ERROR] value parse float error , the value is ", value)
				log.Info("err : ", err.Error())
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