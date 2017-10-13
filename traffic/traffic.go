package traffic

import (
	"context"

	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/golang/protobuf/proto"
	"github.com/northwesternmutual/kanali/config"
  "github.com/northwesternmutual/kanali/logging"
	"github.com/northwesternmutual/kanali/spec"
)

// ReportTraffic reports a new traffic point to etcd
func (ctlr *EtcdController) ReportTraffic(ctx context.Context, pt *spec.TrafficPoint) {

  logger := logging.WithContext(ctx)

	data, err := proto.Marshal(pt)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if _, err := ctlr.Client.Put(ctx, config.FlagEtcdPrefix.GetLong(), string(data)); err != nil {
		switch err {
		case context.Canceled:
			logger.Error(err.Error())
		case context.DeadlineExceeded:
			logger.Error(err.Error())
		case rpctypes.ErrEmptyKey:
			logger.Error(err.Error())
		default:
			logger.Error(err.Error())
		}
		return
	}

  logger.Debug("traffic point reported")

}

// MonitorTraffic watches for new traffic and adds to to the in memory traffic store
func (ctlr *EtcdController) MonitorTraffic() {

	rch := ctlr.Client.Watch(context.Background(), config.FlagEtcdPrefix.GetLong())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			go handleNewTrafficPoint(ev.Kv.Value)
		}
	}

}

func handleNewTrafficPoint(data []byte) {

  logger := logging.WithContext(nil)

	tp := &spec.TrafficPoint{}
	err := proto.Unmarshal(data, tp)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if err := spec.TrafficStore.Set(*tp); err != nil {
		logger.Error(err.Error())
	}

  logger.Debug("traffic point received and processed")

}