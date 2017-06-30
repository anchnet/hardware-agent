package cron

import (
	"time"

	"github.com/51idc/custom-agent/g"
	"github.com/open-falcon/common/model"
	"github.com/51idc/custom-agent/funcs"
)

func Collect() {

	if !g.Config().Transfer.Enabled {
		return
	}

	if len(g.Config().Transfer.Addrs) == 0 {
		return
	}

	for _, v := range funcs.Mappers {
		if (g.Config().Cycle > 0) {
			go collect(int64(g.Config().Cycle), v.Fs)
		} else {
			go collect(int64(v.Interval), v.Fs)
		}
		go collect(int64(v.Interval), v.FsAlive)
	}
}

func collect(sec int64, fns []func() []*model.MetricValue) {
	t := time.NewTicker(time.Second * time.Duration(sec)).C
	index_count := 0
	for {
		if index_count == 0 {
			hostname, err := g.Hostname()
			if err != nil {
				continue
			}

			mvs := []*model.MetricValue{}

			for _, fn := range fns {
				items := fn()
				if items == nil {
					continue
				}

				if len(items) == 0 {
					continue
				}

				for _, mv := range items {
					mvs = append(mvs, mv)
				}
			}

			now := time.Now().Unix()
			for j := 0; j < len(mvs); j++ {
				mvs[j].Step = sec
				mvs[j].Endpoint = hostname
				mvs[j].Timestamp = now
			}
			g.SendToTransfer(mvs)
		} else {
			<-t

			hostname, err := g.Hostname()
			if err != nil {
				continue
			}

			mvs := []*model.MetricValue{}

			for _, fn := range fns {
				items := fn()
				if items == nil {
					continue
				}

				if len(items) == 0 {
					continue
				}

				for _, mv := range items {
					mvs = append(mvs, mv)
				}
			}

			now := time.Now().Unix()
			for j := 0; j < len(mvs); j++ {
				mvs[j].Step = sec
				mvs[j].Endpoint = hostname
				mvs[j].Timestamp = now
			}
			g.SendToTransfer(mvs)
		}
		index_count++
	}
}
