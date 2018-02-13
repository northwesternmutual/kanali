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
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/coreos/etcd/pkg/transport"
	"github.com/golang/protobuf/proto"
	"github.com/northwesternmutual/kanali/cmd/kanali/app/options"
	"github.com/northwesternmutual/kanali/pkg/log"
	store "github.com/northwesternmutual/kanali/pkg/store/kanali/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapgrpc"
)

// EtcdController contains an etcd client
// to be used when talking to etcd
type Controller struct {
	Client *clientv3.Client
}

var ctlr *Controller

// NewController create a new traffic controller
func NewController() (*Controller, error) {

	if ctlr != nil {
		return ctlr, nil
	}

	clientv3.SetLogger(zapgrpc.NewLogger(log.WithContext(nil)))

	etcdConfig := clientv3.Config{
		Endpoints:            viper.GetStringSlice(options.FlagEtcdEndpoints.GetLong()),
		DialTimeout:          time.Second * 5,
		AutoSyncInterval:     time.Second * 5,
		DialKeepAliveTime:    time.Second * 5,
		DialKeepAliveTimeout: time.Second * 5,
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

	ctlr = &Controller{
		Client: etcdClient,
	}

	return ctlr, nil

}

func isTLSDefined() bool {
	return len(viper.GetString(options.FlagEtcdCertFile.GetLong())) > 0 && len(viper.GetString(options.FlagEtcdKeyFile.GetLong())) > 0
}

// Report reports a new traffic point to etcd
func (ctlr *Controller) Report(ctx context.Context, pt *store.TrafficPoint) {
	logger := log.WithContext(ctx)

	data, err := proto.Marshal(pt)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	resp, err := ctlr.Client.Put(ctx, viper.GetString(options.FlagEtcdPrefix.GetLong()), string(data))
	if err != nil {
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
	logger.Debug("traffic point reported",
		zap.Uint64("ClusterId", resp.Header.ClusterId),
		zap.Uint64("MemberId", resp.Header.MemberId),
		zap.Int64("Revision", resp.Header.Revision),
		zap.Uint64("RaftTerm", resp.Header.RaftTerm),
	)
}

// Run begins monitoring traffic. When traffic is discovered,
// it is sent for processing. If the given context is cancelled,
// gracefull termination will commence and this method will return
// a nil error. If an error occurs while monitoring, then a non-nil
// error will be returned.
func (ctlr *Controller) Run(ctx context.Context) error {
	logger := log.WithContext(nil)

	respCh := ctlr.Client.Watch(ctx, viper.GetString(options.FlagEtcdPrefix.GetLong()), clientv3.WithPrefix())

	for watchResp := range respCh {
		logger.Debug("etcd info",
			zap.Uint64("ClusterId", watchResp.Header.ClusterId),
			zap.Uint64("MemberId", watchResp.Header.MemberId),
			zap.Int64("Revision", watchResp.Header.Revision),
			zap.Uint64("RaftTerm", watchResp.Header.RaftTerm),
		)
		if err := watchResp.Err(); err != nil {
			return err
		}
		for _, ev := range watchResp.Events {
			logger.Debug("new etcd event",
				zap.String("type", ev.Type.String()),
				zap.String("key", string(ev.Kv.Key)),
			)
			go processTraffic(ev.Kv.Value)
		}
	}

	logger.Debug("traffic controller will begin gracefull termination")
	if err := ctlr.Client.Close(); err != nil {
		logger.Error("traffic controller gracefull termination failed" + err.Error())
	}
	logger.Info("traffic controller gracefull termination successful")
	return nil
}

func processTraffic(data []byte) {

	logger := log.WithContext(nil)

	tp := &store.TrafficPoint{}
	err := proto.Unmarshal(data, tp)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	store.TrafficStore().Set(tp)
	logger.Debug("traffic point received and processed",
		zap.Int64("timestamp", tp.Time),
		zap.String("namespace", tp.Namespace),
		zap.String("namespace", tp.ProxyName),
		zap.String("namespace", tp.KeyName),
	)

}
