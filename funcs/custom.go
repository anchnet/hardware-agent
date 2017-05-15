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
		fmt.Printf("out:", L)
	}
	return L
}

func path_file_exec(fpath string, L []*model.MetricValue) ([]*model.MetricValue) {
	cmd := exec.Command(fpath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("cmd.Output: ", err)
		return nil
	}
	cmd.Start()
	buff := bufio.NewReader(stdout)

	if err != nil {
		log.Println("[ERROR] exec custom file ", fpath, "fail. error:", err)
		return nil
	}
	// exec successfully
	var i = 0
	for {
		buf, err := buff.ReadString('\n')
		if err == io.EOF {
			if (i == 0) {
				log.Println("[DEBUG] stdout of", fpath, "is blank")
			}
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
			i++
		}
	}

	return L
}