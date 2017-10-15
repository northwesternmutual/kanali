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
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"
	"github.com/northwesternmutual/kanali/config"
	"github.com/northwesternmutual/kanali/logging"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapgrpc"
)

// EtcdController contains an etcd client
// to be used when talking to etcd
type EtcdController struct {
	Client *clientv3.Client
}

// NewController create a new etcd client
func NewController() (*EtcdController, error) {

	clientv3.SetLogger(zapgrpc.NewLogger(logging.WithContext(nil)))

	tlsInfo := transport.TLSInfo{
		CertFile:      viper.GetString(config.FlagEtcdCertFile.GetLong()),
		KeyFile:       viper.GetString(config.FlagEtcdKeyFile.GetLong()),
		TrustedCAFile: viper.GetString(config.FlagEtcdCaFile.GetLong()),
	}
	tlsConfig, err := tlsInfo.ClientConfig()
	if err != nil {
		return nil, err
	}

	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: viper.GetStringSlice(config.FlagEtcdEndpoints.GetLong()),
		TLS:       tlsConfig,
	})
	if err != nil {
		return nil, err
	}

	return &EtcdController{
		Client: etcdClient,
	}, nil

}
