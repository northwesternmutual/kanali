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
		FlagServerSecurePort,
		FlagServerInsecurePort,
		FlagServerInsecureBindAddress,
		FlagServerSecureBindAddress,
		FlagServerTLSCertFile,
		FlagServerTLSKeyFile,
		FlagServerTLSCaFile,
	)
}

var (
	// FlagServerSecurePort sets the port that Kanali will listen on for incoming requests
	FlagServerSecurePort = flags.Flag{
		Long:  "server.secure_port",
		Short: "",
		Value: 0,
		Usage: "Sets the port that Kanali will listen on for incoming requests.",
	}
	// FlagServerInsecurePort sets the port that Kanali will listen on for incoming requests
	FlagServerInsecurePort = flags.Flag{
		Long:  "server.insecure_port",
		Short: "",
		Value: 8080,
		Usage: "Sets the port that Kanali will listen on for incoming requests.",
	}
	// FlagServerInsecureBindAddress specifies the network address that Kanali will listen on for incoming requests
	FlagServerInsecureBindAddress = flags.Flag{
		Long:  "server.insecure_bind_address",
		Short: "",
		Value: "0.0.0.0",
		Usage: "Network address that Kanali will listen on for incoming requests.",
	}
	// FlagServerSecureBindAddress specifies the network address that Kanali will listen on for incoming requests
	FlagServerSecureBindAddress = flags.Flag{
		Long:  "server.secure_bind_address",
		Short: "",
		Value: "0.0.0.0",
		Usage: "Network address that Kanali will listen on for incoming requests.",
	}
	// FlagServerTLSCertFile specifies the path to x509 certificate for HTTPS servers.
	FlagServerTLSCertFile = flags.Flag{
		Long:  "server.tls.cert_file",
		Short: "c",
		Value: "",
		Usage: "Path to x509 certificate for HTTPS servers.",
	}
	// FlagServerTLSKeyFile pecifies the path to x509 private key matching --tls-cert-file
	FlagServerTLSKeyFile = flags.Flag{
		Long:  "server.tls.key_file",
		Short: "k",
		Value: "",
		Usage: "Path to x509 private key matching --server.tls.cert_file.",
	}
	// FlagServerTLSCaFile specifies the path to x509 certificate authority bundle for mutual TLS
	FlagServerTLSCaFile = flags.Flag{
		Long:  "server.tls.ca_file",
		Short: "",
		Value: "",
		Usage: "Path to x509 certificate authority bundle for mutual TLS.",
	}
)
