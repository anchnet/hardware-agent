package funcs

import (
	"github.com/anchnet/hardware-agent/g"
	"github.com/open-falcon/common/model"
)

type FuncsAndInterval struct {
	Fs       []func() []*model.MetricValue
	FsAlive  []func() []*model.MetricValue
	Interval int
}

var Mappers []FuncsAndInterval

func BuildMappers() {
	interval := g.Config().Transfer.Interval
	Mappers = []FuncsAndInterval{
		FuncsAndInterval{
			Fs: []func() []*model.MetricValue{
				HardwareMetrics,
			},
			FsAlive: []func() []*model.MetricValue{
				AgentMetrics,
			},
			Interval: interval,
		},
	}
}
