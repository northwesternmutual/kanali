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

func init() {
	Flags.Add(
		FlagTLSCertFile,
		FlagTLSKeyFile,
		FlagTLSCaFile,
	)
}

var (
	// FlagTLSCertFile specifies the path to x509 certificate for HTTPS servers.
	FlagTLSCertFile = Flag{
		Long:  "tls.cert_file",
		Short: "c",
		Value: "",
		Usage: "Path to x509 certificate for HTTPS servers.",
	}
	// FlagTLSKeyFile pecifies the path to x509 private key matching --tls-cert-file
	FlagTLSKeyFile = Flag{
		Long:  "tls.key_file",
		Short: "k",
		Value: "",
		Usage: "Path to x509 private key matching --tls.cert_file.",
	}
	// FlagTLSCaFile specifies the path to x509 certificate authority bundle for mutual TLS
	FlagTLSCaFile = Flag{
		Long:  "tls.ca_file",
		Short: "",
		Value: "",
		Usage: "Path to x509 certificate authority bundle for mutual TLS.",
	}
)
