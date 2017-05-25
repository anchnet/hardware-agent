package funcs

import (
	"github.com/open-falcon/common/model"
	"log"
	"github.com/51idc/custom-agent/g"
	"os/exec"
	"fmt"
	"bufio"
	"io"
	"strings"
	"strconv"
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
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("[ERROR] exec custom file ", fpath, "fail. error:", err)
		return nil
	}
	cmd.Start()
	buff := bufio.NewReader(stdout)
	//buff.Peek(1)
	//log.Println(buff.Buffered())

	//if (buff.Buffered() == 0) {
	//	log.Println("[DEBUG] stdout of", fpath, "is blank")
	//} else {
	// exec successfully
	for {
		buf, err := buff.ReadString('\n')
		if err != nil && err != io.EOF {
			log.Println("[ERROR] stdout of", fpath, "error :", err)
			break
		}
		s := strings.Split(buf, " ")
		if (len(s) > 1) {
			tag := s[0]
			value := s[1]
			value = strings.Replace(value, "\n", "", -1)
			tags := fmt.Sprintf("name=%s", tag)
			if val, err := strconv.ParseFloat(value, 64); err == nil {
				L = append(L, GaugeValue("custom.data", val, tags))
			} else {
				log.Println("[ERROR] value parse float error , the value is ", value)
			}
		}
		if err == io.EOF {
			break
		}
	}
	//}
	cmd.Wait()

	return L
}