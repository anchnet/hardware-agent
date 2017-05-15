package main

import (
	"fmt"
	"flag"
	"os"
	"./g"
	"./funcs"
	"./cron"
	"./http"
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
	funcs.BuildMappers()
	cron.Collect()

	go http.Start()

	select {}

}



