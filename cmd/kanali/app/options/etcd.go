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

package options

import (
	"github.com/northwesternmutual/kanali/pkg/flags"
)

func init() {
	KanaliGatewayOptions.Add(
		FlagEtcdCertFile,
		FlagEtcdKeyFile,
		FlagEtcdCaFile,
		FlagEtcdEndpoints,
		FlagEtcdPrefix,
	)
}

var (
	// FlagEtcdCertFile specifies the path to x509 certificate for ETCD servers.
	FlagEtcdCertFile = flags.Flag{
		Long:  "etcd.cert_file",
		Short: "",
		Value: "",
		Usage: "Path to x509 certificate for ETCD servers.",
	}
	// FlagEtcdKeyFile pecifies the path to x509 private key matching --etcd.cert-file
	FlagEtcdKeyFile = flags.Flag{
		Long:  "etcd.key_file",
		Short: "",
		Value: "",
		Usage: "Path to x509 private key matching --etcd.cert_file.",
	}
	// FlagEtcdCaFile specifies the path to x509 certificate authority bundle for mutual TLS
	FlagEtcdCaFile = flags.Flag{
		Long:  "etcd.ca_file",
		Short: "",
		Value: "",
		Usage: "Path to x509 certificate authority bundle for mutual TLS.",
	}
	// FlagEtcdEndpoints specifies a comma delimited list of ETCD hosts.
	FlagEtcdEndpoints = flags.Flag{
		Long:  "etcd.endpoints",
		Short: "",
		Value: []string{
			"http://127.0.0.1:2379",
		},
		Usage: "Comma delimited list of ETCD hosts.",
	}
	// FlagEtcdPrefix specifies the ETCD prefix.
	FlagEtcdPrefix = flags.Flag{
		Long:  "etcd.prefix",
		Short: "",
		Value: "kanali",
		Usage: "ETCD prefix.",
	}
)
