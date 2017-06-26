package main

import (
	"fmt"
	"flag"
	"os"
	"github.com/51idc/custom-agent/g"
	"github.com/51idc/custom-agent/funcs"
	"github.com/51idc/custom-agent/cron"
	"github.com/51idc/custom-agent/http"
	"time"
	"log"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	check := flag.Bool("check", false, "check collector")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	if *check {
		funcs.CheckCollector()
		os.Exit(0)
	}

	g.ParseConfig(*cfg)
	g.InitRootDir()
	g.InitLocalIps()
	g.InitRpcClients()
	if (g.Config().StartTime != "undefined") {
		log.Println("collecting will start at :", g.Config().StartTime)
		for {
			if (g.Config().StartTime == time.Now().Format("15:04")) {
				break;
			}
			time.Sleep(60)
		}
	}
	funcs.BuildMappers()
	cron.Collect()

	go http.Start()

	select {}

}



