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