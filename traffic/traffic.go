// Copyright (c) 2017 Northwestern Mutual.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

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
