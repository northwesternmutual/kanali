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

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/coreos/etcd/pkg/transport"
	"github.com/golang/protobuf/proto"
	"github.com/northwesternmutual/kanali/cmd/kanali/app/options"
	"github.com/northwesternmutual/kanali/pkg/logging"
	store "github.com/northwesternmutual/kanali/pkg/store/kanali/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapgrpc"
)

// EtcdController contains an etcd client
// to be used when talking to etcd
type Controller struct {
	Client *clientv3.Client
}

// NewController create a new traffic controller
func NewController() (*Controller, error) {

	clientv3.SetLogger(zapgrpc.NewLogger(logging.WithContext(nil)))

	etcdConfig := clientv3.Config{
		Endpoints: viper.GetStringSlice(options.FlagEtcdEndpoints.GetLong()),
	}

	if isTLSDefined() {
		tlsInfo := transport.TLSInfo{
			CertFile:      viper.GetString(options.FlagEtcdCertFile.GetLong()),
			KeyFile:       viper.GetString(options.FlagEtcdKeyFile.GetLong()),
			TrustedCAFile: viper.GetString(options.FlagEtcdCaFile.GetLong()),
		}
		tlsConfig, err := tlsInfo.ClientConfig()
		if err != nil {
			return nil, err
		}
		etcdConfig.TLS = tlsConfig
	}

	etcdClient, err := clientv3.New(etcdConfig)
	if err != nil {
		return nil, err
	}

	return &Controller{
		Client: etcdClient,
	}, nil

}

func isTLSDefined() bool {
	return len(viper.GetString(options.FlagEtcdCertFile.GetLong())) > 0 && len(viper.GetString(options.FlagEtcdKeyFile.GetLong())) > 0
}

// Report reports a new traffic point to etcd
func (ctlr *Controller) Report(ctx context.Context, pt *store.TrafficPoint) {

	logger := logging.WithContext(ctx)

	data, err := proto.Marshal(pt)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if _, err := ctlr.Client.Put(ctx, options.FlagEtcdPrefix.GetLong(), string(data)); err != nil {
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
func (ctlr *Controller) MonitorTraffic(ctx context.Context) {

	rch := ctlr.Client.Watch(ctx, options.FlagEtcdPrefix.GetLong())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			go handleNewTrafficPoint(ev.Kv.Value)
		}
	}

}

func handleNewTrafficPoint(data []byte) {

	logger := logging.WithContext(nil)

	tp := &store.TrafficPoint{}
	err := proto.Unmarshal(data, tp)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	store.TrafficStore().Set(tp)
	logger.Debug("traffic point received and processed")

}
