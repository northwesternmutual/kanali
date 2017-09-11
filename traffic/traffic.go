package traffic

import (
	"context"

	"github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/golang/protobuf/proto"
	"github.com/northwesternmutual/kanali/config"
	"github.com/northwesternmutual/kanali/spec"
)

// ReportTraffic reports a new traffic point to etcd
func (ctlr *EtcdController) ReportTraffic(ctx context.Context, pt *spec.TrafficPoint) {

	data, err := proto.Marshal(pt)
	if err != nil {
		logrus.Errorf("could not marshal traffic point: %s", err.Error())
		return
	}

	if _, err := ctlr.Client.Put(ctx, config.FlagEtcdPrefix.GetLong(), string(data)); err != nil {
		switch err {
		case context.Canceled:
			logrus.Errorf("could not write traffic point - ctx is canceled by another routine: %s", err.Error())
		case context.DeadlineExceeded:
			logrus.Errorf("could not write traffic point - ctx is attached with a deadline is exceeded: %s", err.Error())
		case rpctypes.ErrEmptyKey:
			logrus.Errorf("could not write traffic point - client-side error: %s", err.Error())
		default:
			logrus.Errorf("could not write traffic point - bad cluster endpoints, which are not etcd servers: %s", err.Error())
		}
		return
	}

  logrus.Debug("traffic point reported")

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

	tp := &spec.TrafficPoint{}
	err := proto.Unmarshal(data, tp)
	if err != nil {
		logrus.Errorf("could not unmarshal traffic point: %s", err.Error())
		return
	}

	if err := spec.TrafficStore.Set(*tp); err != nil {
		logrus.Errorf("could not add traffic point to store: %s", err.Error())
	}

  logrus.Debug("traffic point received and processed")

}
