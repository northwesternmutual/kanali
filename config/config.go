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

package config

import (
	"errors"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type flag struct {
	long  string
	short string
	value interface{}
	usage string
}

type flags []flag

var (
	// FlagTLSCaFile specifies the path to x509 certificate authority bundle for mutual TLS
	FlagTLSCaFile = flag{
		long:  "tls-ca-file",
		short: "a",
		value: "",
		usage: "Path to x509 certificate authority bundle for mutual TLS.",
	}
	// FlagBindAddress specifies the network address that Kanali will listen on for incoming requests
	FlagBindAddress = flag{
		long:  "bind-address",
		short: "b",
		value: "0.0.0.0",
		usage: "Network address that Kanali will listen on for incoming requests.",
	}
	// FlagTLSCertFile specifies the path to x509 certificate for HTTPS servers.
	FlagTLSCertFile = flag{
		long:  "tls-cert-file",
		short: "c",
		value: "",
		usage: "Path to x509 certificate for HTTPS servers.",
	}
	// FlagDecryptionKeyFile specifies the path to valid PEM-encoded private key that matches the public key used to encrypt API keys
	FlagDecryptionKeyFile = flag{
		long:  "decryption-key-file",
		short: "d",
		value: "",
		usage: "Path to valid PEM-encoded private key that matches the public key used to encrypt API keys.",
	}
	// FlagHeaderMaskValue sets the value to be used when omitting header values
	FlagHeaderMaskValue = flag{
		long:  "header-mask-value",
		short: "f",
		value: "ommitted",
		usage: "Sets the value to be used when omitting header values.",
	}
	// FlagPluginsLocation sets the location of custom plugins shared object (.so) files
	FlagPluginsLocation = flag{
		long:  "plugins-location",
		short: "g",
		value: "/",
		usage: "Location of custom plugins shared object (.so) files.",
	}
	// FlagInfluxdbAddr specifies the Influxdb address. Addr should be of the form 'http://host:port' or 'http://[ipv6-host%zone]:port'
	FlagInfluxdbAddr = flag{
		long:  "influxdb-addr",
		short: "i",
		value: "http://monitoring-influxdb.kube-system.svc.cluster.local",
		usage: "Influxdb address. Addr should be of the form 'http://host:port' or 'http://[ipv6-host%zone]:port'",
	}
	// FlagJaegerSamplerServerURL specifies the endpoint to the Jaeger sampler server
	FlagJaegerSamplerServerURL = flag{
		long:  "jaeger-sampler-server-url",
		short: "j",
		value: "jaeger-all-in-one-agent.default.svc.cluster.local",
		usage: "Endpoint to the Jaeger sampler server",
	}
	// FlagTLSPrivateKeyFile specifies the path to x509 private key matching --tls-cert-file
	FlagTLSPrivateKeyFile = flag{
		long:  "tls-private-key-file",
		short: "k",
		value: "",
		usage: "Path to x509 private key matching --tls-cert-file.",
	}
	// FlagLogLevel sets the logging level. Choose between 'debug', 'info', 'warn', 'error', 'fatal'
	FlagLogLevel = flag{
		long:  "log-level",
		short: "l",
		value: "info",
		usage: "Sets the logging level. Choose between 'debug', 'info', 'warn', 'error', 'fatal'.",
	}
	// FlagEnableMock enables Kanali's mock responses feature. Read the documentation for more information
	FlagEnableMock = flag{
		long:  "enable-mock",
		short: "m",
		value: false,
		usage: "Enables Kanali's mock responses feature. Read the documentation for more information.",
	}
	// FlagJaegerAgentURL specifies the endpoint to the Jaeger agent
	FlagJaegerAgentURL = flag{
		long:  "jaeger-agent-url",
		short: "n",
		value: "jaeger-all-in-one-agent.default.svc.cluster.local",
		usage: "Endpoint to the Jaeger agent",
	}
	// FlagPeerUDPPort sets the port that all Kanali instances will communicate to each other over
	FlagPeerUDPPort = flag{
		long:  "peer-udp-port",
		short: "o",
		value: 10001,
		usage: "Sets the port that all Kanali instances will communicate to each other over.",
	}
	// FlagKanaliPort sets the port that Kanali will listen on for incoming requests
	FlagKanaliPort = flag{
		long:  "kanali-port",
		short: "p",
		value: 0,
		usage: "Sets the port that Kanali will listen on for incoming requests.",
	}
	// FlagInfluxdbUsername specifies the Influxdb username
	FlagInfluxdbUsername = flag{
		long:  "influxdb-username",
		short: "q",
		value: "",
		usage: "Influxdb username",
	}
	// FlagInfluxdbPassword specifies the Influxdb password
	FlagInfluxdbPassword = flag{
		long:  "influxdb-password",
		short: "r",
		value: "",
		usage: "Influxdb password",
	}
	// FlagInfluxdbDatabase specifies the Influxdb database
	FlagInfluxdbDatabase = flag{
		long:  "influxdb-database",
		short: "s",
		value: "k8s",
		usage: "Influxdb database",
	}
	// FlagEnableProxyProtocol maintains the integrity of the remote client IP address when incoming traffic to Kanali includes the Proxy Protocol header
	FlagEnableProxyProtocol = flag{
		long:  "enable-proxy-protocol",
		short: "t",
		value: false,
		usage: "Maintain the integrity of the remote client IP address when incoming traffic to Kanali includes the Proxy Protocol header.",
	}
	// FlagEnableClusterIP enables to use of cluster ip as opposed to Kubernetes DNS for upstream routing
	FlagEnableClusterIP = flag{
		long:  "enable-cluster-ip",
		short: "u",
		value: false,
		usage: "Enables to use of cluster ip as opposed to Kubernetes DNS for upstream routing.",
	}
	// FlagDisableTLSCnValidation disables common name validate as part of an SSL handshake
	FlagDisableTLSCnValidation = flag{
		long:  "disable-tls-cn-validation",
		short: "v",
		value: false,
		usage: "Disable common name validate as part of an SSL handshake.",
	}
	// FlagUpstreamTimeout sets the length of upstream timeout
	FlagUpstreamTimeout = flag{
		long:  "upstream-timeout",
		short: "w",
		value: "0h0m10s",
		usage: "Set length of upstream timeout. Defaults to none",
	}
	// FlagApikeyHeaderKey specifies the name of the HTTP header holding the apikey
	FlagApikeyHeaderKey = flag{
		long:  "apikey-header-key",
		short: "y",
		value: "apikey",
		usage: "Name of the HTTP header holding the apikey.",
	}
	// FlagEtcdEndpoints specifies the list of etcd endpoints to connect with (scheme://ip:port), comma separated.
	FlagEtcdEndpoints = flag{
		long:  "etcd-endpoints",
		short: "",
		value: []string{},
		usage: "List of etcd endpoints to connect with (scheme://ip:port), comma separated.",
	}
	// FlagEtcdCAFile specifies the SSL Certificate Authority file used to secure etcd communication
	FlagEtcdCAFile = flag{
		long:  "etcd-cafile",
		short: "",
		value: "",
		usage: "SSL Certificate Authority file used to secure etcd communication.",
	}
	// FlagEtcdCertFile specifies the SSL certification file used to secure etcd communication.
	FlagEtcdCertFile = flag{
		long:  "etcd-certfile",
		short: "",
		value: "",
		usage: "SSL certification file used to secure etcd communication.",
	}
	// FlagEtcdKeyFile specifies the SSL key file used to secure etcd communication.
	FlagEtcdKeyFile = flag{
		long:  "etcd-keyfile",
		short: "",
		value: "",
		usage: "SSL key file used to secure etcd communication.",
	}
	// FlagEtcdPrefix specifies the prefix to prepend to all resource paths in etcd.
	FlagEtcdPrefix = flag{
		long:  "etcd-prefix",
		short: "",
		value: "kanali",
		usage: "The prefix to prepend to all resource paths in etcd.",
	}
)

// Flags represents the complete set of configuration options that Kanali can use
var Flags = flags{
	FlagTLSCaFile,
	FlagBindAddress,
	FlagTLSCertFile,
	FlagDecryptionKeyFile,
	FlagPluginsLocation,
	FlagInfluxdbAddr,
	FlagJaegerSamplerServerURL,
	FlagTLSPrivateKeyFile,
	FlagLogLevel,
	FlagEnableMock,
	FlagJaegerAgentURL,
	FlagPeerUDPPort,
	FlagKanaliPort,
	FlagInfluxdbUsername,
	FlagInfluxdbPassword,
	FlagInfluxdbDatabase,
	FlagEnableProxyProtocol,
	FlagEnableClusterIP,
	FlagDisableTLSCnValidation,
	FlagUpstreamTimeout,
	FlagHeaderMaskValue,
	FlagApikeyHeaderKey,
}

func (f flag) GetLong() string {
	return f.long
}

func (f flag) GetShort() string {
	return f.short
}

func (f flag) GetUsage() string {
	return f.usage
}

func (f flags) AddAll(cmd *cobra.Command) error {

	for _, currFlag := range f {
		switch v := currFlag.value.(type) {
		case int:
			cmd.Flags().IntP(currFlag.long, currFlag.short, v, currFlag.usage)
		case bool:
			cmd.Flags().BoolP(currFlag.long, currFlag.short, v, currFlag.usage)
		case string:
			cmd.Flags().StringP(currFlag.long, currFlag.short, v, currFlag.usage)
		case time.Duration:
			cmd.Flags().DurationP(currFlag.long, currFlag.short, v, currFlag.usage)
		case []string:
			cmd.Flags().StringSliceP(currFlag.long, currFlag.short, v, currFlag.usage)
		default:
			return errors.New("unsupported flag type")
		}
		if err := viper.BindPFlag(currFlag.long, cmd.Flags().Lookup(currFlag.long)); err != nil {
			return err
		}
		viper.SetDefault(currFlag.long, currFlag.value)
	}

	return nil

}
